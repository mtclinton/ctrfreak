package pkg

import (
	"context"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/platforms"
	"github.com/containerd/containerd/v2/defaults"
	"github.com/containerd/containerd/v2/core/remotes/docker"
	refdocker "github.com/distribution/reference"
	"github.com/pkg/errors"

)


type EnsuredImage struct {
	Ref         string
	Image       containerd.Image
	Snapshotter string
}

// PullMode is either one of "always", "missing", "never"
type PullMode = string

func EnsureImage(ctx context.Context, client *containerd.Client, rawRef string, mode PullMode) (*EnsuredImage, error) {
	named, err := refdocker.ParseDockerRef(rawRef)
	if err != nil {
		return nil, err
	}
	ref := named.String()

	sn := defaults.DefaultSnapshotter

	if mode != "always" {
		if i, err := client.ImageService().Get(ctx, ref); err == nil {
			image := containerd.NewImage(client, i)
			res := &EnsuredImage{
				Ref:         ref,
				Image:       image,
				Snapshotter: sn,
			}
			if unpacked, err := image.IsUnpacked(ctx, sn); err == nil && !unpacked {
				if err := image.Unpack(ctx, sn); err != nil {
					return nil, err
				}
			}
			return res, nil
		}
	}

	if mode == "never" {
		return nil, errors.Errorf("image %q is not available", rawRef)
	}

	ctx, done, err := client.WithLease(ctx)
	if err != nil {
		return nil, err
	}
	defer done(ctx)

	resovlerOpts := docker.ResolverOptions{}
	resolver := docker.NewResolver(resovlerOpts)

	config := &FetchConfig{
		Resolver:        resolver,
		PlatformMatcher: platforms.Default(),
	}

	img, err := Fetch(ctx, client, ref, config)
	if err != nil {
		return nil, err
	}
	i := containerd.NewImageWithPlatform(client, img, config.PlatformMatcher)
	if err = i.Unpack(ctx, sn); err != nil {
		return nil, err
	}
	res := &EnsuredImage{
		Ref:         ref,
		Image:       i,
		Snapshotter: sn,
	}
	return res, nil
}
