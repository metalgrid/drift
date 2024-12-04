package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"
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
	var opts *zeroconf.ZeroconfOptions = nil

	if len(os.Args) > 1 {
		opts = &zeroconf.ZeroconfOptions{
			Identity: os.Args[1],
		}
	}

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

	zcSvc, err := zeroconf.NewZeroconfService(servicePort, fmt.Sprintf("%x", *pubkey), opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating zeroconf service")
	}

	err = zcSvc.Start(appCtx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting zeroconf service")
	}
	defer zcSvc.Shutdown()

	transferRequests := make(chan platform.Request)
	platformGateway := platform.NewGateway(zcSvc.Peers(), transferRequests)
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting transfer gateway")
	}

	// Connection error handling
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
					log.Warn().Str("peer", peer.Instance).Msg("public key not found")
					_ = conn.Close()
					continue
				}

				decodedKey, err := hex.DecodeString(pk)
				if err != nil {
					log.Warn().Str("peer", peer.Instance).Str("pk", pk).Err(err).Msg("invalid public key")
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
					continue
				}

				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.Addresses[0], peer.Port))
				if err != nil {
					platformGateway.Notify(fmt.Sprintf("Unable to connect to peer: %s", err))
					return
				}

				pk, err := hex.DecodeString(peer.GetRecord("pk"))
				if err != nil {
					platformGateway.Notify(fmt.Sprintf("Unable to retrieve peer's public key: %s", err))
					return
				}

				var peerpk [32]byte
				copy(peerpk[:], pk)

				sc, err := secret.SecureConnection(conn, &peerpk, privkey)
				if err != nil {
					platformGateway.Notify(fmt.Sprintf("Unable to secure connection with peer: %s", err))
				}

				go transport.HandleConnection(context.WithValue(appCtx, "filename", request.File), sc, platformGateway)
				transport.SendFile(request.File, sc)
			}
		}
	}()

	err = platformGateway.Run(appCtx)
	if err != nil {
		log.Error().Err(err).Msg("platform gateway failed")
		shutdown()
	}
	wg.Wait()
}
