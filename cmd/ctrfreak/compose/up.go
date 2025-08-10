package compose

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/compose-spec/compose-go/v2/cli"
	"log"
	"context"


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

	projectName := "my_project"
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{files[0]},
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithName(projectName),
	)
	if err != nil {
		log.Fatal(err)
	}

	project, err := options.LoadProject(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Use the MarshalYAML method to get YAML representation
	projectYAML, err := project.MarshalYAML()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(projectYAML))
    return nil
}
