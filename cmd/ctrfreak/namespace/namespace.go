package namespace

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"ctrfreak/pkg"
)

func NamespaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "namespace",
		Aliases:       []string{"ns"},
		Short:         "Manage containerd namespaces",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(listCommand())
	cmd.AddCommand(createCommand())
	return cmd
}

func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ls",
		Aliases:       []string{"list"},
		Short:         "List containerd namespaces",
		RunE:          listAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "c",
		Aliases:       []string{"create"},
		Short:         "create containerd namespace",
		Args:          cobra.MinimumNArgs(1),
		RunE:          createAction,
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
	namespaces := client.NamespaceService()
	nss, err := namespaces.List(ctx)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 1, 8, 1, ' ', 0)
	if _, err := fmt.Fprintln(tw, "NAME\tLABELS\t"); err != nil {
		return err
	}
	for _, ns := range nss {
		labels, err := namespaces.Labels(ctx, ns)
		if err != nil {
			return err
		}

		var labelStrings []string
		for k, v := range labels {
			labelStrings = append(labelStrings, strings.Join([]string{k, v}, "="))
		}
		sort.Strings(labelStrings)

		if _, err := fmt.Fprintf(tw, "%v\t%v\t\n", ns, strings.Join(labelStrings, ",")); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func createAction(cmd *cobra.Command, args []string) error {
	client, ctx, cancel, err := pkg.NewClient(cmd.Context(), "default", "unix:///run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer cancel()
	namespaces := client.NamespaceService()
	return namespaces.Create(ctx, args[0], nil)
}