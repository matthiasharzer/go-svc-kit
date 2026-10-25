package ioutil

import (
	"context"
	"errors"
	"io"
)

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (cr *contextReader) Read(p []byte) (int, error) {
	err := cr.ctx.Err()
	if err != nil {
		return 0, err
	}
	return cr.reader.Read(p)
}

// CopyWithContext works like io.Copy, but interrupts writing if the context.Context is canceled. Returns the number
// of bytes copied. If the context is canceled before the entire write completed, no error will be returned.
func CopyWithContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	ctxReader := &contextReader{ctx: ctx, reader: src}

	written, err := io.Copy(dst, ctxReader)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			// expected when client disconnects while writing
			return written, nil
		}
		return written, err
	}
	return written, nil
}
