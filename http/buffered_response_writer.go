package http

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

type BufferedResponseWriter struct {
	parent http.ResponseWriter
	buffer *bytes.Buffer
}

func NewBufferedResponseWriter(parent http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		parent: parent,
		buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
	}
}

func (b BufferedResponseWriter) Header() http.Header {
	return b.parent.Header()
}

func (b BufferedResponseWriter) Write(bytes []byte) (int, error) {
	return b.buffer.Write(bytes)
}

func (b BufferedResponseWriter) WriteHeader(statusCode int) {
	b.parent.Header().Add("Content-Length", strconv.Itoa(b.buffer.Len()))
	b.parent.WriteHeader(statusCode)
}

func (b BufferedResponseWriter) Close() error {
	_, err := io.Copy(b.parent, b.buffer)
	return err
}
