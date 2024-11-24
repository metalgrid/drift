package transport

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func HandleConnection(conn net.Conn) {
	fmt.Println("handling connection", conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr))
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		raw, err := reader.ReadString(byte(endOfMessage))
		if err != nil {
			fmt.Println("error reading from remote", err)
			return
		}

		msg := UnmarshalMessage(raw)

		switch m := msg.(type) {
		case error:
			fmt.Printf("error: %+v", m)
			return
		case Offer:
			fmt.Printf("offer: %+v", m)
			_, err = conn.Write(Accept().MarshalMessage())
			if err != nil {
				return
			}
			err = storeFile("dropzone-incoming", m.Filename, reader)
			if err != nil {
				fmt.Println("failed storing:", err)
				return
			}
		case Answer:
			fmt.Printf("answer: %+v", m)
			if m.Accepted() {
				err = sendFile("/etc/os-release", writer)
				if err != nil {
					fmt.Println("failed sending:", err)
					return
				}
			}
		}
	}
}

func storeFile(incoming, file string, reader *bufio.Reader) error {
	err := os.MkdirAll(incoming, 0777)

	if err != nil {
		return err
	}

	f, err := os.Create(incoming + string(os.PathSeparator) + file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)

	return err
}

func sendFile(file string, writer *bufio.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(writer, f)

	return err
}

func SendFile(filename string, conn net.Conn) {
	offer, err := MakeOffer(filename)
	if err != nil {
		fmt.Println("failed creating file offer:", err)
		return
	}

	_, err = conn.Write(offer.MarshalMessage())
	if err != nil {
		fmt.Println("failed sending file:", err)
	}
}
