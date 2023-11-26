package pipeline

import (
	"bufio"
	"bytes"
	"io"
)

type Buffer struct {
	bytes  []byte
	reader io.Reader
}

func NewBuffer(reader io.Reader) *Buffer {
	return &Buffer{
		bytes:  []byte{},
		reader: reader,
	}
}

func (c *Buffer) Reader() io.Reader {
	c.ReadAll()
	return bufio.NewReader(bytes.NewBuffer(c.bytes))
}

func (c *Buffer) ReadAll() error {
	bytes, err := io.ReadAll(c.reader)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	if len(bytes) == 0 {
		return nil
	}
	c.bytes = bytes
	return nil
}

func (c *Buffer) Bytes() []byte {
	c.ReadAll()
	return c.bytes
}

func (c *Buffer) Len() int {
	c.ReadAll()
	return len(c.bytes)
}
