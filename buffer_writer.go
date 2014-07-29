package logserver

import (
	"os"
	"time"
	"bufio"
)

/*
type bufWriter struct{
	name string
	buf  []byte
	wr *os.File
}

func (o *bufWriter) Write(b []byte) {
	o.buf = append(o.buf, b...)
}

func (o *bufWriter) Flush() {
	o.wr.Write(o.buf)
	o.buf = o.buf[:0]
}

func (o *bufWriter) Close() {
	o.Flush()
	o.wr.Close()
}

func (o *bufWriter) Rotate(now time.Time) {

	o.Close()

	o.wr = newLogFile(o.name, now)
}

func newBufferWriter(name string) *bufWriter {
	buf := make([]byte, 5*bufSize)
	w := &bufWriter{
		name:name,
		buf:buf[:0],
	}

	w.Rotate(time.Now())
	return w
}
*/

type bufWriter struct{
	name string
	f *os.File
	wr *bufio.Writer
}

func (o *bufWriter) Write(b []byte) {
	o.wr.Write(b)
}

func (o *bufWriter) Flush() {
	o.wr.Flush()
}

func (o *bufWriter) Close() {
	o.Flush()
	o.f.Close()
}

func (o *bufWriter) Reset(now time.Time) *os.File {

	o.f = newLogFile(o.name, now)

	o.wr.Reset(o.f)

	return o.f
}

func newBufferWriter(name string) *bufWriter {

	w := &bufWriter{
		name:name,
	}

	w.f = newLogFile(name, time.Now())

	w.wr = bufio.NewWriterSize(w.f, 8*bufSize)

	return w
}

type BufferWriter struct {
	lastHour int
	wr    map[string]*bufWriter
}

func (o *BufferWriter) Write(k string, b []byte) {
	o.wr[k].Write(b)
}

func (o *BufferWriter) Rotate(now time.Time) {
	h := now.Hour()
	if o.lastHour == h {
		return
	}

	//o.lastHour = h
	//print(h)

	for k, _ := range logType {
		o.wr[k].Close()
		o.wr[k].Reset(now)
	}
}

func (o *BufferWriter) Flush() {
	for k, _ := range logType {
		o.wr[k].Flush()
	}
}

func NewBufferWriter() *BufferWriter {

	wr := make(map[string]*bufWriter, len(logType))
	for k, v := range logType {
		wr[k] = newBufferWriter(v)
	}

	fw := &BufferWriter{-1, wr}

	return fw
}

