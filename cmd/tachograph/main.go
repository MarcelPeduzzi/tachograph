package main

import (
	"context"
	"fmt"
	"image/color"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/tachograph-go"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCommand(),
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			base := c(lipgloss.Black, lipgloss.White)
			baseInverted := c(lipgloss.White, lipgloss.Black)
			return fang.ColorScheme{
				Base:         base,
				Title:        base,
				Description:  base,
				Comment:      base,
				Flag:         base,
				FlagDefault:  base,
				Command:      base,
				QuotedString: base,
				Argument:     base,
				Help:         base,
				Dash:         base,
				ErrorHeader:  [2]color.Color{baseInverted, base},
				ErrorDetails: base,
			}
		}),
	); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tachograph",
		Short: "Tachograph CLI",
	}
	cmd.AddGroup(&cobra.Group{ID: "ddd", Title: ".DDD Files"})
	cmd.AddCommand(newParseCommand())
	cmd.AddGroup(&cobra.Group{ID: "utils", Title: "Utils"})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

func newParseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse <file1> [file2] [...]",
		Short:   "Parse .DDD files",
		GroupID: "ddd",
		Args:    cobra.MinimumNArgs(1),
	}

	raw := cmd.Flags().Bool("raw", false, "Output raw intermediate format (skip semantic parsing)")
	authenticate := cmd.Flags().Bool("authenticate", false, "Authenticate signatures and certificates")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		for _, filename := range args {
			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", filename, err)
			}

			// Step 1: Always unmarshal to raw format first
			rawFile, err := tachograph.UnmarshalRawFile(data)
			if err != nil {
				return fmt.Errorf("error parsing raw %s: %w", filename, err)
			}

			// Step 2: Optionally authenticate (works on raw files)
			if *authenticate {
				if err := tachograph.AuthenticateRawFile(ctx, rawFile); err != nil {
					return fmt.Errorf("error authenticating %s: %w", filename, err)
				}
			}

			// Step 3: Output raw or parse to semantic format
			if *raw {
				// Output raw format (with or without authentication)
				fmt.Println(protojson.Format(rawFile))
			} else {
				// Parse to semantic format (authentication results are propagated)
				file, err := tachograph.ParseRawFile(rawFile)
				if err != nil {
					return fmt.Errorf("error parsing %s: %w", filename, err)
				}
				fmt.Println(protojson.Format(file))
			}
		}
		return nil
	}
	return cmd
}
