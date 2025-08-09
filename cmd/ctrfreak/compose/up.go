package compose

import (
	"fmt"
	"github.com/spf13/cobra"


	"ctrfreak/pkg"
)
var files []string

func UpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "up",
		Short:         "Run container with docker compose file",
		RunE:          upAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().StringArrayVar(&files, "file", []string{}, "Run container based on docker compose")

	return cmd
}


func upAction(cmd *cobra.Command, args []string) error {
	_, _, cancel, err := pkg.NewClient(cmd.Context(), "default", "unix:///run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer cancel()

    fmt.Println(files)
    return nil
}
