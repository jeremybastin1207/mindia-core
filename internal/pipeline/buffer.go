package pipeline

import (
	"bufio"
	"bytes"
	"io"
)

type Buffer struct {
	Body     []byte
	skipRead bool
	Reader   io.Reader
}

func (c *Buffer) MergeReader() io.Reader {
	c.ReadAll()
	return io.MultiReader(bufio.NewReader(bytes.NewBuffer(c.Body)), c.Reader)
}

func (c *Buffer) ReadAll() *bytes.Buffer {
	if !c.skipRead {
		body, _ := io.ReadAll(c.Reader)
		c.Body = body
		c.skipRead = true
	}
	return bytes.NewBuffer(c.Body)
}

func (c *Buffer) Len() int64 {
	c.ReadAll()
	return int64(len(c.Body))
}
