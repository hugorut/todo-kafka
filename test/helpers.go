package test

import "bytes"

type Writer struct {
	msg bytes.Buffer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.msg.Write(p)
}

func (w *Writer) Msg() string {
	return w.msg.String()
}

