package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/metalgrid/dropzone/internal/notification"
	"github.com/metalgrid/dropzone/internal/secret"
	"github.com/metalgrid/dropzone/internal/server"
	"github.com/metalgrid/dropzone/internal/transport"
	"github.com/metalgrid/dropzone/internal/zeroconf"
	"github.com/rs/zerolog/log"
)

func main() {
	n := notification.NewNotifier()
	n.SendNotification()

	privkey, pubkey, err := secret.GenerateX25519KeyPair()
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating encryption keys")
	}

	appCtx, shutdown := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer shutdown()
	wg := &sync.WaitGroup{}

	servicePort, connections, connectionErrors, err := server.Start(appCtx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed listening for connections")
	}

	zcSvc, err := zeroconf.NewZeroconfService(servicePort, fmt.Sprintf("%x", &pubkey))
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating zeroconf service")
	}

	defer zcSvc.Shutdown()

	err = zcSvc.Discover(appCtx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed discovering local services")
	}

	err = zcSvc.Advertise()
	if err != nil {
		log.Fatal().Err(err).Msg("failed advertising ourselves")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-appCtx.Done():
				log.Info().Str("system", "connection_errors_processor").Msg("stopping")
				return
			case err := <-connectionErrors:
				log.Error().Err(err).Msg("incoming connection error")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-appCtx.Done():
				log.Info().Str("system", "connection_processor").Msg("stopping")
				return
			case conn := <-connections:
				defer conn.Close() // TODO: too many deferred ?
				peer := zcSvc.Peers().GetByAddr(conn.RemoteAddr())
				if peer == nil {
					log.Warn().Stringer("address", conn.RemoteAddr()).Msg("unknown peer")
					_ = conn.Close()
					continue
				}

				pk := peer.GetRecord("pk")
				if pk == "" {
					log.Warn().Str("peer", peer.GetInstance()).Msg("public key not found")
					_ = conn.Close()
					continue
				}

				decodedKey, err := hex.DecodeString(pk)
				if err != nil {
					log.Warn().Str("peer", peer.GetInstance()).Str("pk", pk).Err(err).Msg("invalid public key")
					_ = conn.Close()
					continue
				}

				var peerPublicKey [32]byte
				copy(peerPublicKey[:], decodedKey)

				// Secure the connection
				sc, err := secret.SecureConnection(conn, &peerPublicKey, privkey)

				go transport.HandleConnection(sc)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		log.Debug().Msg("Attempting to send a file to the first service we have")
		for _, peer := range zcSvc.Peers().All() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.AddrIPv4[0], peer.Port))
			if err != nil {
				log.Error().Err(err).Str("peer", peer.Instance).Msg("failed to connect to peer")
				return
			}
			defer conn.Close()
			pk, err := hex.DecodeString(peer.GetRecord("pk"))
			if err != nil {
				panic(err)
			}

			var peerpk [32]byte
			copy(peerpk[:], pk)
			sc, err := secret.SecureConnection(conn, &peerpk, privkey)
			if err != nil {
				panic(err)
			}

			go transport.HandleConnection(sc)
			transport.SendFile("/etc/os-release", sc)

		}
	}()

	wg.Wait()
}
