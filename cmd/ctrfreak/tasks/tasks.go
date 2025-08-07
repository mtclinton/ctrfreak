package tasks

import (
	"github.com/spf13/cobra"
)

func TasksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tasks",
		Aliases:       []string{"t"},
		Short:         "Manage containerd tasks",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(ListCommand())
	return cmd
}
