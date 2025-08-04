package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"ctrfreak/cmd/ctrfreak"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#A1FAFA")).
	Background(lipgloss.Color("#1D56A4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(32)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ctrfreak",
	Short: "Docker-compatible CLI for containerd",
}

// init function adds the testCmd to the rootCmd
func init() {
	rootCmd.AddCommand(
	    ctrfreak.PsCommand(),
	    ctrfreak.ContainerCommand(),
	    ctrfreak.NamespaceCommand(),
	)
}

// main function executes the rootCmd
func main() {
	if err := fang.Execute(context.TODO(), rootCmd); err != nil {
		os.Exit(1)
	}
}
