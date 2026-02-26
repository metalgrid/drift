package app

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"sync"

	"github.com/metalgrid/drift/internal/config"
	"github.com/metalgrid/drift/internal/platform"
	"github.com/metalgrid/drift/internal/secret"
	"github.com/metalgrid/drift/internal/server"
	"github.com/metalgrid/drift/internal/transport"
	"github.com/metalgrid/drift/internal/zeroconf"
	"github.com/rs/zerolog/log"
)

func Run(ctx context.Context, identity string) error {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		cfg = config.DefaultConfig()
	}

	if identity == "" && cfg.Identity != "" {
		identity = cfg.Identity
	}

	var opts *zeroconf.ZeroconfOptions
	if identity != "" {
		opts = &zeroconf.ZeroconfOptions{Identity: identity}
	}

	privkey, pubkey, err := secret.GenerateX25519KeyPair()
	if err != nil {
		return fmt.Errorf("failed creating encryption keys: %w", err)
	}

	wg := &sync.WaitGroup{}

	servicePort, connections, connectionErrors, err := server.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed listening for connections: %w", err)
	}

	zcSvc, err := zeroconf.NewZeroconfService(servicePort, fmt.Sprintf("%x", *pubkey), opts)
	if err != nil {
		return fmt.Errorf("failed creating zeroconf service: %w", err)
	}

	if err := zcSvc.Start(ctx); err != nil {
		return fmt.Errorf("failed starting zeroconf service: %w", err)
	}
	defer zcSvc.Shutdown()

	transferRequests := make(chan platform.Request)
	platformGateway := platform.NewGateway(zcSvc.Peers(), transferRequests)
	if err != nil {
		return fmt.Errorf("failed starting transfer gateway: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
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
			case <-ctx.Done():
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

				peerPublicKey, err := parsePeerPublicKeyHex(pk)
				if err != nil {
					log.Warn().Str("peer", peer.Instance).Str("pk", pk).Err(err).Msg("invalid public key")
					_ = conn.Close()
					continue
				}

				sc, err := secret.SecureConnection(conn, peerPublicKey, privkey)
				if err != nil {
					log.Warn().Err(err).Msg("failed securing connection")
					_ = conn.Close()
					continue
				}
				go transport.HandleConnection(ctx, sc, platformGateway, nil)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("system", "outbound_connection_processor").Msg("stopping")
				return
			case request := <-transferRequests:
				peer := zcSvc.Peers().GetByInstance(request.To)
				if peer == nil {
					platformGateway.Notify(fmt.Sprintf("User %s not found", request.To))
					continue
				}

				targetHost, ok := firstPeerAddress(peer.Addresses)
				if !ok {
					platformGateway.Notify("Unable to connect: peer has no addresses")
					continue
				}

				target := net.JoinHostPort(targetHost, strconv.Itoa(peer.Port))
				conn, err := net.Dial("tcp", target)
				if err != nil {
					platformGateway.Notify(fmt.Sprintf("Unable to connect to peer: %s", err))
					continue
				}

				peerpk, err := parsePeerPublicKeyHex(peer.GetRecord("pk"))
				if err != nil {
					platformGateway.Notify(fmt.Sprintf("Unable to retrieve peer's public key: %s", err))
					_ = conn.Close()
					continue
				}

				sc, err := secret.SecureConnection(conn, peerpk, privkey)
				if err != nil {
					platformGateway.Notify(fmt.Sprintf("Unable to secure connection with peer: %s", err))
					_ = conn.Close()
					continue
				}

				if len(request.Files) > 1 {
					outbound := transport.NewOutboundTransferState()
					go transport.HandleConnection(ctx, sc, platformGateway, outbound)
					if err := transport.SendBatch(request.Files, sc, outbound); err != nil {
						platformGateway.Notify(fmt.Sprintf("Unable to send batch offer: %s", err))
						_ = sc.Close()
					}
				} else if len(request.Files) == 1 {
					outbound := transport.NewOutboundTransferState()
					go transport.HandleConnection(ctx, sc, platformGateway, outbound)
					if err := transport.SendFile(request.Files[0], sc, outbound); err != nil {
						platformGateway.Notify(fmt.Sprintf("Unable to send file offer: %s", err))
						_ = sc.Close()
					}
				}
			}
		}
	}()

	if err := platformGateway.Run(ctx); err != nil {
		return fmt.Errorf("platform gateway failed: %w", err)
	}

	wg.Wait()
	return nil
}

func parsePeerPublicKeyHex(pkHex string) (*[32]byte, error) {
	decodedKey, err := hex.DecodeString(pkHex)
	if err != nil {
		return nil, err
	}
	if len(decodedKey) != 32 {
		return nil, fmt.Errorf("invalid public key length %d", len(decodedKey))
	}

	var peerPublicKey [32]byte
	copy(peerPublicKey[:], decodedKey)
	return &peerPublicKey, nil
}

func firstPeerAddress(addrs []netip.Addr) (string, bool) {
	if len(addrs) == 0 {
		return "", false
	}
	return addrs[0].String(), true
}
