package container

import (
	"github.com/spf13/cobra"
)

func ContainerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "container",
		Short:         "Manage containerd containers",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(PsCommand())
	cmd.AddCommand(RunCommand())
	return cmd
}



