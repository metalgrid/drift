package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/metalgrid/drift/internal/platform"
	"github.com/metalgrid/drift/internal/secret"
	"github.com/metalgrid/drift/internal/server"
	"github.com/metalgrid/drift/internal/transport"
	"github.com/metalgrid/drift/internal/zeroconf"
	"github.com/rs/zerolog/log"
)

func main() {
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

	zcSvc, err := zeroconf.NewZeroconfService(servicePort, fmt.Sprintf("%x", *pubkey))
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

	platformGateway := platform.NewGateway(zcSvc.Peers())
	transferRequests, err := platformGateway.Run(appCtx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting transfer gateway")
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

	// Incoming connection handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-appCtx.Done():
				log.Info().Str("system", "inbound_connection_processor").Msg("stopping")
				return
			case conn := <-connections:
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
				go transport.HandleConnection(appCtx, sc, platformGateway)
			}
		}
	}()

	// Outgoing connection handler
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-appCtx.Done():
				log.Info().Str("system", "outbound_connection_processor").Msg("stopping")
				return
			case request := <-transferRequests:
				peer := zcSvc.Peers().GetByInstance(request.To)
				if peer == nil {
					platformGateway.Notify(fmt.Sprintf("User %s not found", request.To))
					log.Warn().Str("peer", request.To).Msg("could not find peer")
					continue
				}
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.AddrIPv4[0], peer.Port))
				if err != nil {
					log.Error().Err(err).Str("peer", peer.Instance).Msg("failed to connect to peer")
					return
				}
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

				go transport.HandleConnection(context.WithValue(appCtx, "filename", request.File), sc, platformGateway)
				transport.SendFile(request.File, sc)
			}
		}
	}()

	wg.Wait()
}
