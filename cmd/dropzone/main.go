package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/metalgrid/dropzone/internal/notification"
	"github.com/metalgrid/dropzone/internal/secret"
	"github.com/metalgrid/dropzone/internal/server"
	"github.com/metalgrid/dropzone/internal/zeroconf"
	"github.com/rs/zerolog/log"
)

func main() {
	n := notification.NewNotifier()
	n.SendNotification()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("failed determining local machine's hostname")
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("failed determining local user")
	}

	username := user.Username
	if user.Name != "" {
		username = user.Name
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

	zcSvc, err := zeroconf.Advertise(
		servicePort,
		username,
		hostname,
		fmt.Sprintf("%x", *pubkey),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed registering to zeroconf")
	}

	defer zcSvc.Shutdown()

	peers, err := zeroconf.Discover(appCtx)
	if err != nil {
		panic(err)
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
				log.Info().Str("system", "connection_processor").Msg("shutting down")
				return
			case conn := <-connections:
				defer conn.Close()
				peer := peers.GetByAddr(conn.RemoteAddr())
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

				go handleEncryptedConnection(conn, &peerPublicKey, privkey)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		log.Debug().Msg("Attempting to send a file to the first service we have")
		for _, peer := range peers.All() {
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

			sc.Write([]byte("OFFER|Something, blah-blah...\n"))

			rdr := bufio.NewReader(sc)
			for {
				msg, err := rdr.ReadString('\n')
				if err != nil {
					panic(err)
				}
				fmt.Print(msg)
			}
		}
	}()

	wg.Wait()
}

func handleEncryptedConnection(conn net.Conn, pubkey, privkey *[32]byte) {
	defer conn.Close()

	sc, err := secret.SecureConnection(conn, pubkey, privkey)
	if err != nil {
		log.Error().Err(err).Msg("failed securing connection")
		return
	}

	reader := bufio.NewReader(sc)

	invalidMessages := 0
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Error().Err(err).Msg("failed reading from encrypted connection")
		}

		switch {
		case strings.HasPrefix(msg, "OFFER"):
			log.Debug().Str("message", msg).Msg("received a file transfer offer")
			sc.Write([]byte("ACCEPT|this should actually fly when it gets a bit bigger i hope...\n"))
		default:
			log.Warn().Str("message", msg).Msg("received invalid message")
			invalidMessages++
			if invalidMessages > 5 {
				return
			}
		}
	}
}
