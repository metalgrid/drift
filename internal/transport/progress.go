package transport

import "io"

// ProgressFunc is called after each Write/Read operation with cumulative bytes transferred and total bytes
type ProgressFunc func(bytesTransferred int64, totalBytes int64)

// ProgressWriter wraps an io.Writer and calls a callback after each Write operation
type ProgressWriter struct {
	writer           io.Writer
	totalBytes       int64
	bytesTransferred int64
	callback         ProgressFunc
}

// NewProgressWriter creates a new ProgressWriter that tracks write progress
func NewProgressWriter(w io.Writer, total int64, fn ProgressFunc) *ProgressWriter {
	return &ProgressWriter{
		writer:     w,
		totalBytes: total,
		callback:   fn,
	}
}

// Write writes data to the underlying writer and calls the progress callback
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.bytesTransferred += int64(n)
	if pw.callback != nil {
		pw.callback(pw.bytesTransferred, pw.totalBytes)
	}
	return n, err
}

// ProgressReader wraps an io.Reader and calls a callback after each Read operation
type ProgressReader struct {
	reader           io.Reader
	totalBytes       int64
	bytesTransferred int64
	callback         ProgressFunc
}

// NewProgressReader creates a new ProgressReader that tracks read progress
func NewProgressReader(r io.Reader, total int64, fn ProgressFunc) *ProgressReader {
	return &ProgressReader{
		reader:     r,
		totalBytes: total,
		callback:   fn,
	}
}

// Read reads data from the underlying reader and calls the progress callback
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.bytesTransferred += int64(n)
	if pr.callback != nil {
		pr.callback(pr.bytesTransferred, pr.totalBytes)
	}
	return n, err
}
