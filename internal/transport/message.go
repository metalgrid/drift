package transport

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	mimeType       = "application/octet-stream"
	fieldSeparator = "|"
	endOfMessage   = '\n'
)

type Message struct {
	Type string
}

type Offer struct {
	Message
	Filename string
	Mimetype string
	Size     int64
}

func (o Offer) MarshalMessage() []byte {
	return []byte(
		strings.Join(
			[]string{
				o.Type,
				o.Filename,
				o.Mimetype,
				strconv.FormatInt(o.Size, 10),
			},
			fieldSeparator,
		) + string(endOfMessage),
	)
}

type Answer struct {
	Message
	Kind string
}

func (a Answer) Accepted() bool {
	return a.Kind == "ACCEPT"
}

func (a Answer) MarshalMessage() []byte {
	return []byte(
		strings.Join(
			[]string{
				a.Type,
				a.Kind,
			},
			fieldSeparator,
		) + string(endOfMessage),
	)
}

func UnmarshalMessage(msg string) any {
	var err error
	msg, _ = strings.CutSuffix(msg, string(endOfMessage))
	switch {
	case strings.HasPrefix(msg, "OFFER"):
		parts := strings.Split(msg, fieldSeparator)
		if len(parts) != 4 {
			break
		}
		var size int64
		size, err = strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			break
		}
		return Offer{
			Message{parts[0]},
			parts[1],
			parts[2],
			size,
		}

	case strings.HasPrefix(msg, "ANSWER"):
		parts := strings.Split(msg, fieldSeparator)
		if len(parts) != 2 {
			break
		}
		return Answer{
			Message{parts[0]},
			parts[1],
		}
	}
	return err
}

func MakeOffer(filename string) (Offer, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return Offer{}, fmt.Errorf("failed offering a file transfer: %w", err)
	}

	return Offer{
		Message{"OFFER"},
		fileInfo.Name(),
		mimeType,
		fileInfo.Size(),
	}, nil
}

func Accept() Answer {
	return Answer{
		Message{"ANSWER"},
		"ACCEPT",
	}
}

func Decline() Answer {
	return Answer{
		Message{"ANSWER"},
		"DECLINE",
	}
}

func formatSize(size int64) string {
	const (
		kib = 1024
		mib = kib * 1024
		gib = mib * 1024
		tib = gib * 1024
	)

	switch {
	case size >= tib:
		return fmt.Sprintf("%.2f TiB", float64(size)/float64(tib))
	case size >= gib:
		return fmt.Sprintf("%.2f GiB", float64(size)/float64(gib))
	case size >= mib:
		return fmt.Sprintf("%.2f MiB", float64(size)/float64(mib))
	case size >= kib:
		return fmt.Sprintf("%.2f KiB", float64(size)/float64(kib))
	default:
		return fmt.Sprintf("%d Bytes", size)
	}
}
