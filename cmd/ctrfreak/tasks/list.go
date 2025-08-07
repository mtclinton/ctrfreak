package tasks

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"text/tabwriter"
	tasks "github.com/containerd/containerd/api/services/tasks/v1"

	"ctrfreak/pkg"
)


func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ls",
		Aliases:       []string{"list"},
		Short:         "List containerd tasks",
		RunE:          listAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}


func listAction(cmd *cobra.Command, args []string) error {
	client, ctx, cancel, err := pkg.NewClient(cmd.Context(), "default", "unix:///run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer cancel()
	s := client.TaskService()
		response, err := s.List(ctx, &tasks.ListTasksRequest{})
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "TASK\tPID\tSTATUS\t")
		for _, task := range response.Tasks {
			if _, err := fmt.Fprintf(w, "%s\t%d\t%s\n",
				task.ID,
				task.Pid,
				task.Status.String(),
			); err != nil {
				return err
			}
		}
		return w.Flush()
}
