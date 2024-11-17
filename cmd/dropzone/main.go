package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/metalgrid/dropzone/internal/secret"
	"github.com/metalgrid/dropzone/internal/server"
	"github.com/metalgrid/dropzone/internal/zeroconf"
	"github.com/rs/zerolog/log"
)

func main() {

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

	appCtx, shutdownApp := context.WithCancel(context.Background())
	defer shutdownApp()

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

	go func() {
		for err := range connectionErrors {
			log.Error().Err(err).Msg("incoming connection error")
		}
	}()

	go func() {
		for conn := range connections {

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
	}()

	time.Sleep(time.Second * 2)
	log.Debug().Msg("Attempting to send a file to the first service we have")
	for _, peer := range peers.All() {
		conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.AddrIPv4[0], peer.Port))
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

		b := make([]byte, 1024)
		for {
			_, err := sc.Read(b)
			if err != nil {
				panic(err)
			}
			fmt.Print(string(b))
		}
	}

	select {}
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
