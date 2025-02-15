package file

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
)

type FileUpdate interface {
	io.WriteCloser
	io.StringWriter
}

type fileUpdate struct {
	*bytes.Buffer
	path     string
	f        *os.File
	complete bool
	allow    func(b []byte) bool
}

func (f *fileUpdate) Close() (err error) {
	if f.complete {
		return nil
	}
	defer func() {
		f.complete = true
		if f.f != nil {
			_ = f.f.Close()
		}
	}()
	var c []byte
	if f.f != nil {
		b := &bytes.Buffer{}
		if _, err = io.Copy(b, f.f); err != nil {
			return errors.Wrap(err, "failed to read file")
		}
		_ = f.f.Close()
		f.f = nil
		c = b.Bytes()
		if f.allow != nil && !f.allow(c) {
			return
		}
	}
	if !bytes.Equal(c, f.Buffer.Bytes()) {
		if f.f, err = os.Create(f.path); err != nil {
			return errors.Wrap(err, "failed to open file for update")
		}
		if _, err = io.Copy(f.f, f.Buffer); err != nil {
			return errors.Wrap(err, "failed to write file")
		}
	}
	return
}

func NewFileUpdate(path string, allow func(b []byte) bool) (u FileUpdate, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "failed to open file for update")
	}
	return &fileUpdate{
		Buffer: &bytes.Buffer{},
		path:   path,
		f:      f,
		allow:  allow,
	}, nil
}

var _ io.WriteCloser = (*fileUpdate)(nil)

func OnlyGenerated(b []byte) bool {
	return bytes.Contains(b, []byte("Code generated by"))
}
