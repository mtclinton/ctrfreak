package container

import (
	"ctrfreak/pkg"
	"fmt"
	"errors"
	"strconv"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/oci"
	"github.com/containerd/errdefs"
	"github.com/containerd/log"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/cobra"
)


func RunCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:           "run",
		Args:          cobra.MinimumNArgs(1),
		Short:         "Run containers",
		RunE:          runAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().BoolP("detach", "d", false, "Run container in background")
	cmd.Flags().Bool("rm", false, "Automatically remove the container")

	return cmd
}

func runAction(cmd *cobra.Command, args []string) error {
	client, ctx, cancel, err := pkg.NewClient(cmd.Context(), "default", "unix:///run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer cancel()

	id := args[0]
	ref := args[0]
	var opts  []oci.SpecOpts
    var cOpts []containerd.NewContainerOpts
    var spec  containerd.NewContainerOpts
    fmt.Println(id)

    fmt.Println("fuck")

    var image containerd.Image
    i, err := client.ImageService().Get(ctx, ref)
    if err != nil {
        return err
    }
    fmt.Println("lol")
    image = containerd.NewImage(client, i)

    cOpts = append(cOpts,
        containerd.WithImage(image),
        containerd.WithImageConfigLabels(image),
        containerd.WithSnapshotter("overlayfs"))
    cOpts = append(cOpts, containerd.WithImageStopSignal(image, "SIGTERM"))

	var s specs.Spec
	spec = containerd.WithSpec(&s, opts...)

	cOpts = append(cOpts, spec)

	container, err := client.NewContainer(ctx, id, cOpts...)
    if err != nil {
        return err
    }
	rm, err := cmd.Flags().GetBool("rm")
	if err != nil {
		return err
	}

	detach, err := cmd.Flags().GetBool("detach")
	if err != nil {
		return err
	}

	if rm && !detach {
		defer func() {
			if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
				log.L.WithError(err).Error("failed to cleanup container")
			}
		}()
	}

	task, err := pkg.NewTask(ctx, client, container, nil, nil)
	if err != nil {
		return err
	}

	var statusC <-chan containerd.ExitStatus
	if !detach {
		defer func() {

			if _, err := task.Delete(ctx, containerd.WithProcessKill); err != nil && !errdefs.IsNotFound(err) {
				log.L.WithError(err).Error("failed to cleanup task")
			}
		}()

		if statusC, err = task.Wait(ctx); err != nil {
			return err
		}
	}

	if err := task.Start(ctx); err != nil {
		return err
	}
	if detach {
		return nil
	}

	sigc := ForwardAllSignals(ctx, task)
	defer StopCatch(sigc)

	status := <-statusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	if _, err := task.Delete(ctx); err != nil {
		return err
	}
	if code != 0 {
		return errors.New(strconv.Itoa(int(code)))
	}
	return nil
}
