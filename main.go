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

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test cli",
	Short: "Tests that cli is working correctly",
	Args:  cobra.MaximumNArgs(1), // Allows 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		name := "World"
		if len(args) > 0 {
			name = args[0]
		}
		fmt.Println(style.Render("TESTING: Hello,", name, "!\n"))
	},
}

// init function adds the testCmd to the rootCmd
func init() {
	rootCmd.AddCommand(testCmd, ctrfreak.PsCommand(), ctrfreak.NamespaceCommand())
}

// main function executes the rootCmd
func main() {
	if err := fang.Execute(context.TODO(), rootCmd); err != nil {
		os.Exit(1)
	}
}
