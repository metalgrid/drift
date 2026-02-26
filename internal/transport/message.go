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

type FileEntry struct {
	Filename string
	Mimetype string
	Size     int64
}

type BatchOffer struct {
	Message
	Files []FileEntry
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

func (b BatchOffer) MarshalMessage() []byte {
	parts := []string{b.Type, strconv.Itoa(len(b.Files))}
	for _, file := range b.Files {
		parts = append(parts, file.Filename, file.Mimetype, strconv.FormatInt(file.Size, 10))
	}
	return []byte(strings.Join(parts, fieldSeparator) + string(endOfMessage))
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
	case strings.HasPrefix(msg, "BATCH_OFFER"):
		parts := strings.Split(msg, fieldSeparator)
		if len(parts) < 2 {
			break
		}
		count, parseErr := strconv.Atoi(parts[1])
		if parseErr != nil {
			err = parseErr
			break
		}
		if count <= 0 {
			err = fmt.Errorf("batch count must be positive")
			break
		}
		expectedParts := 2 + (count * 3)
		if len(parts) != expectedParts {
			err = fmt.Errorf("batch count mismatch: expected %d parts, got %d", expectedParts, len(parts))
			break
		}
		files := make([]FileEntry, count)
		for i := 0; i < count; i++ {
			idx := 2 + (i * 3)
			if err = validateProtocolFilename(parts[idx]); err != nil {
				break
			}
			size, parseErr := strconv.ParseInt(parts[idx+2], 10, 64)
			if parseErr != nil {
				err = parseErr
				break
			}
			files[i] = FileEntry{
				Filename: parts[idx],
				Mimetype: parts[idx+1],
				Size:     size,
			}
		}
		if err != nil {
			break
		}
		return BatchOffer{
			Message{parts[0]},
			files,
		}

	case strings.HasPrefix(msg, "OFFER"):
		parts := strings.Split(msg, fieldSeparator)
		if len(parts) != 4 {
			break
		}
		if err = validateProtocolFilename(parts[1]); err != nil {
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
	if err := validateProtocolFilename(fileInfo.Name()); err != nil {
		return Offer{}, err
	}

	return Offer{
		Message{"OFFER"},
		fileInfo.Name(),
		mimeType,
		fileInfo.Size(),
	}, nil
}

func MakeBatchOffer(filenames []string) (BatchOffer, error) {
	if len(filenames) == 0 {
		return BatchOffer{}, fmt.Errorf("batch offer requires at least one file")
	}

	files := make([]FileEntry, 0, len(filenames))
	for _, filename := range filenames {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			return BatchOffer{}, fmt.Errorf("failed offering file %s: %w", filename, err)
		}
		if err := validateProtocolFilename(fileInfo.Name()); err != nil {
			return BatchOffer{}, err
		}
		files = append(files, FileEntry{
			Filename: fileInfo.Name(),
			Mimetype: mimeType,
			Size:     fileInfo.Size(),
		})
	}

	return BatchOffer{
		Message{"BATCH_OFFER"},
		files,
	}, nil
}

func validateProtocolFilename(name string) error {
	if name == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if strings.Contains(name, fieldSeparator) || strings.ContainsRune(name, endOfMessage) {
		return fmt.Errorf("filename contains reserved protocol delimiters")
	}
	return nil
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
