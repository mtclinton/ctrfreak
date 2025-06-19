package ctrfreak

import (
	"fmt"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"

	"ctrfreak/pkg"
)

func PsCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:           "ps",
		Args:          cobra.NoArgs,
		Short:         "List containers",
		RunE:          psAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	return cmd
}

func psAction(cmd *cobra.Command, args []string) error {
	client, ctx, cancel, err := pkg.NewClient(cmd.Context(), "default", "unix:///run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer cancel()
	containers, err := client.Containers(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
	_, err = fmt.Fprintln(w, "CONTAINER\tIMAGE\tRUNTIME\t")
	if err != nil {
		return err
	}
	for _, c := range containers {
		info, err := c.Info(ctx, containerd.WithoutRefreshedMetadata)
		if err != nil {
			return err
		}
		imageName := info.Image
		if imageName == "" {
			imageName = "-"
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t\n",
			c.ID(),
			imageName,
			info.Runtime.Name,
		); err != nil {
			return err
		}
	}
	return w.Flush()
}
