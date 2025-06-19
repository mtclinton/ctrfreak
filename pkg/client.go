package pkg

import (
	"context"
	"strings"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/namespaces"
)

// NewClient returns a new containerd client
func NewClient(ctx context.Context, namespace, address string, opts ...containerd.Opt) (*containerd.Client, context.Context, context.CancelFunc, error) {
	ctx = namespaces.WithNamespace(ctx, namespace)
	address = strings.TrimPrefix(address, "unix://")
	client, err := containerd.New(address, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	return client, ctx, cancel, nil
}
