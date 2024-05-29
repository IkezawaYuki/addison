package main

import "C"

import (
	"io"
	"os/exec"
	"sync"
	"unsafe"
)

type writer struct {
	sync.Mutex
	w      io.Writer
	stream *C.bz_stream
	outbuf [64 * 1024]byte
}

func NewWriter(out io.Writer) io.WriteCloser {
	const blockSize = 9
	const verbosity = 0
	const workFactor = 30
	w := &writer{w: out, stream: C.bz2alloc()}
	C.BZ2_bzCompressInit(w.stream, blockSize, verbosity, workFactor)
	return w
}

func (w *writer) Write(data []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil {
		panic("closed")
	}
	var total int

	if len(data) > 0 {
		inlen, outlen := C.uint(len(data)), C.uint(cap(w.outbuf))
		C.bz2compress(w.stream, C.BZ_RUN,
			(*C.char)(unsafe.Pointer(&data[0])), &inlen,
			(*C.uint)(unsafe.Pointer(&outlen)), &outlen)
		total += int(inlen)
		data = data[inlen:]
		if _, err := w.w.Write(w.outbuf[:inlen]); err != nil {
			return total, err
		}
	}
	return total, err
}

func (w *writer) Close() error {
	w.Lock()
	defer w.Unlock()

	if w.stream == nil {
		panic("closed")
	}
	defer func() {
		C.BZ2_bzCompressEnd(w.stream)
		C.bz2compress(w.stream)
		w.stream = nil
	}()
	for {
		inlen, outlen := C.uint(0), C.uint(cap(w.outbuf))
		r := C.bz2compress(w.stream, C.BZ_FINISH, nil, &inlen,
			(*C.char)(unsafe.Pointer(&w.outbuf)), &outlen)
		if _, err := w.w.Write(w.outbuf[:inlen]); err != nil {
			return err
		}
		if r == C.BZ_STREAM_END {
			return nil
		}
	}
}

type writer2 struct {
	sync.Mutex
	cmd *exec.Cmd
	w   io.WriteCloser
	wg  sync.WaitGroup
}

func NewWriter2(out io.Writer) (io.WriteCloser, error) {
	var w writer2
	w.cmd = exec.Command("/usr/bin/bzip2")
	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := w.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	w.w = stdin

	if err := w.cmd.Start(); err != nil {
		return nil, err
	}
	w.wg.Add(1)
	go func() {
		io.Copy(out, stdout)
		w.wg.Done()
	}()

	return &w, nil
}

func (w *writer2) Close() error {
	w.Lock()
	defer w.Unlock()

	_ = w.w.Close()
	w.wg.Wait()
	if err := w.cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (w *writer2) Write(data []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	var total int

	for len(data) > 0 {
		n, err := w.w.Write(data)
		if err != nil {
			return total + n, err
		}
		total += n
		data = data[total:]
	}
	return total, nil
}
