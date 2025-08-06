package pkg

import (
	"context"
	"os"
	"io"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/cio"
)

type stdinCloser struct {
	stdin  *os.File
	closer func()
}

func (s *stdinCloser) Read(p []byte) (int, error) {
	n, err := s.stdin.Read(p)
	if err == io.EOF {
		if s.closer != nil {
			s.closer()
		}
	}
	return n, err
}

// NewTask creates a new task
func NewTask(ctx context.Context, client *containerd.Client, container containerd.Container, ioOpts []cio.Opt, opts ...containerd.NewTaskOpts) (containerd.Task, error) {
	stdinC := &stdinCloser{
		stdin: os.Stdin,
	}

	spec, err := container.Spec(ctx)
	if err != nil {
		return nil, err
	}

	if spec.Linux != nil {
		if len(spec.Linux.UIDMappings) != 0 {
			opts = append(opts, containerd.WithUIDOwner(spec.Linux.UIDMappings[0].HostID))
		}
		if len(spec.Linux.GIDMappings) != 0 {
			opts = append(opts, containerd.WithGIDOwner(spec.Linux.GIDMappings[0].HostID))
		}
	}

	var ioCreator cio.Creator
	ioCreator = cio.NewCreator(append([]cio.Opt{cio.WithStreams(stdinC, os.Stdout, os.Stderr)}, ioOpts...)...)

	t, err := container.NewTask(ctx, ioCreator, opts...)

	if err != nil {
		return nil, err
	}
	stdinC.closer = func() {
		t.CloseIO(ctx, containerd.WithStdinCloser)
	}
	return t, nil
}