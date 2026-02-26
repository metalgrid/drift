package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	"github.com/metalgrid/drift/internal/platform"
)

const (
	missingOutboundTransferStateToken     = "missing_outbound_transfer_state"
	outboundAnswerNoSupportedPendingToken = "outbound_answer_no_supported_pending_offer"
	unsolicitedAnswerIgnoredToken         = "unsolicited_answer_ignored"
)

type OutboundTransferState struct {
	mu           sync.Mutex
	pendingFiles []string
}

func NewOutboundTransferState() *OutboundTransferState {
	return &OutboundTransferState{}
}

func (s *OutboundTransferState) SetPendingFiles(files []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pendingFiles = append([]string(nil), files...)
}

func (s *OutboundTransferState) ConsumePendingFiles() ([]string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pendingFiles) == 0 {
		return nil, false
	}
	files := append([]string(nil), s.pendingFiles...)
	s.pendingFiles = nil
	return files, true
}

func (s *OutboundTransferState) ClearPendingFiles() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pendingFiles = nil
}

func HandleConnection(ctx context.Context, conn net.Conn, gw platform.Gateway, outbound *OutboundTransferState) {
	fmt.Println("handling connection", conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr))
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		raw, err := reader.ReadString(byte(endOfMessage))
		if err != nil {
			fmt.Println("error reading from remote:", err)
			return
		}

		msg := UnmarshalMessage(raw)

		switch m := msg.(type) {
		case error:
			gw.Notify(fmt.Sprintf("Error: %s", m))
			return
		case BatchOffer:
			var totalSize int64
			fileInfos := make([]platform.FileInfo, len(m.Files))
			for i, file := range m.Files {
				totalSize += file.Size
				fileInfos[i] = platform.FileInfo{
					Filename: file.Filename,
					Size:     file.Size,
				}
			}

			var answer string
			if bg, ok := gw.(platform.BatchGateway); ok {
				answer = bg.AskBatch(conn.RemoteAddr().String(), fileInfos)
			} else {
				answer = gw.Ask(fmt.Sprintf("Incoming batch: %d files (%s)", len(m.Files), formatSize(totalSize)))
			}

			if answer == "ACCEPT" {
				_, err = conn.Write(Accept().MarshalMessage())
				if err != nil {
					return
				}
			}

			if answer == "DECLINE" || answer == "" {
				_, _ = conn.Write(Decline().MarshalMessage())
				return
			}

			fp := filepath.Join(xdg.UserDirs.Download, "Drift")
			for _, file := range m.Files {
				err = storeFile(fp, file.Filename, file.Size, conn, nil)
				if err != nil {
					gw.Notify(fmt.Sprintf("Failed storing file %s: %s", file.Filename, err))
					return
				}
			}
			gw.Notify(fmt.Sprintf("Batch received: %d files", len(m.Files)))
		case Offer:
			answer := gw.Ask(fmt.Sprintf("Incoming file: %s (%s)", m.Filename, formatSize(m.Size)))
			if answer == "ACCEPT" {
				_, err = conn.Write(Accept().MarshalMessage())
				if err != nil {
					return
				}
			}

			// empty string means waiting for an action from the local user has timed out, so we decline by default
			if answer == "DECLINE" || answer == "" {
				_, _ = conn.Write(Decline().MarshalMessage())
				return
			}

			fp := filepath.Join(xdg.UserDirs.Download, "Drift")
			err = storeFile(fp, m.Filename, m.Size, conn, nil)
			if err != nil {
				gw.Notify(fmt.Sprintf("Failed storing file: %s", err))
				return
			}
			gw.Notify(fmt.Sprintf("File received: %s", m.Filename))
		case Answer:
			if m.Accepted() {
				if outbound == nil {
					gw.Notify(missingOutboundTransferStateToken)
					gw.Notify(outboundAnswerNoSupportedPendingToken)
					return
				}

				files, ok := outbound.ConsumePendingFiles()
				if !ok {
					gw.Notify(unsolicitedAnswerIgnoredToken)
					return
				}

				if len(files) == 0 {
					gw.Notify(outboundAnswerNoSupportedPendingToken)
					return
				}

				for _, file := range files {
					err = sendFile(file, conn, nil)
					if err != nil {
						gw.Notify(fmt.Sprintf("Failed sending %s: %s", file, err))
						return
					}
				}

				if len(files) == 1 {
					gw.Notify(fmt.Sprintf("File sent: %s", files[0]))
					continue
				}
				gw.Notify(fmt.Sprintf("Batch sent: %d files", len(files)))
				continue
			}
			if outbound != nil {
				outbound.ClearPendingFiles()
			}
			return
		}
	}
}

func storeFile(incoming, file string, size int64, reader io.Reader, progress ProgressFunc) error {
	err := os.MkdirAll(incoming, 0777)

	if err != nil {
		return err
	}

	fp := filepath.Join(incoming, file)
	// f, err := os.Create(filepath)
	f, err := os.CreateTemp(incoming, file+"*.drift")
	if err != nil {
		return err
	}

	lr := io.LimitReader(reader, size)
	pr := NewProgressReader(lr, size, progress)
	bytes, err := f.ReadFrom(pr)
	_ = bytes
	_ = f.Close()
	if err != nil {
		// if we're able to create the file, we should be able to remove it as well
		_ = os.Remove(f.Name())
		return err
	}
	return os.Rename(f.Name(), fp)
}

func sendFile(file string, writer io.Writer, progress ProgressFunc) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	pw := NewProgressWriter(writer, fi.Size(), progress)
	bytes, err := f.WriteTo(pw)
	_ = bytes
	return err
}

func SendFile(filename string, conn net.Conn, outbound *OutboundTransferState) error {
	fmt.Println("Sending file", filename)
	if outbound == nil {
		return fmt.Errorf("%s", missingOutboundTransferStateToken)
	}
	outbound.SetPendingFiles([]string{filename})

	offer, err := MakeOffer(filename)
	if err != nil {
		outbound.ClearPendingFiles()
		fmt.Println("failed creating file offer:", err)
		return err
	}

	_, err = conn.Write(offer.MarshalMessage())
	if err != nil {
		outbound.ClearPendingFiles()
		fmt.Println("failed sending file:", err)
		return err
	}

	return nil
}

func SendBatch(filenames []string, conn net.Conn, outbound *OutboundTransferState) error {
	fmt.Println("Sending batch of", len(filenames), "files")
	if outbound == nil {
		return fmt.Errorf("%s", missingOutboundTransferStateToken)
	}
	outbound.SetPendingFiles(filenames)

	batch, err := MakeBatchOffer(filenames)
	if err != nil {
		outbound.ClearPendingFiles()
		fmt.Println("failed creating batch offer:", err)
		return err
	}

	_, err = conn.Write(batch.MarshalMessage())
	if err != nil {
		outbound.ClearPendingFiles()
		fmt.Println("failed sending batch:", err)
		return err
	}

	return nil
}
