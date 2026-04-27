package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"mychainctl/pkg/registry"
)

const defaultRegistry = "cgr.dev"

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Image-related commands",
}

var imagesListCmd = &cobra.Command{
	Use:   "list <repository>",
	Short: "List image tags in a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFormat != "table" && outputFormat != "json" {
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}

		client := registry.NewClient(defaultRegistry)
		tags, err := client.ListTags(context.Background(), args[0])
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			payload := struct {
				Repository string   `json:"repository"`
				Tags       []string `json:"tags"`
			}{Repository: args[0], Tags: tags}

			return writeJSON(cmd, payload)
		}

		return writeTagsTable(cmd, tags)
	},
}

var imagesInspectCmd = &cobra.Command{
	Use:   "inspect <image:tag>",
	Short: "Inspect image metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFormat != "table" && outputFormat != "json" {
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}

		client := registry.NewClient(defaultRegistry)
		metadata, err := client.InspectImage(context.Background(), args[0])
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			payload := struct {
				Reference string                  `json:"reference"`
				Metadata  *registry.ImageMetadata `json:"metadata"`
			}{Reference: args[0], Metadata: metadata}

			return writeJSON(cmd, payload)
		}

		return writeMetadataTable(cmd, metadata)
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(imagesListCmd)
	imagesCmd.AddCommand(imagesInspectCmd)
}

func writeJSON(cmd *cobra.Command, payload interface{}) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

func writeTagsTable(cmd *cobra.Command, tags []string) error {
	writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(writer, "TAG"); err != nil {
		return err
	}
	for _, tag := range tags {
		if _, err := fmt.Fprintln(writer, tag); err != nil {
			return err
		}
	}
	return writer.Flush()
}

func writeMetadataTable(cmd *cobra.Command, metadata *registry.ImageMetadata) error {
	writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintf(writer, "Digest\t%s\n", metadata.Digest); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "Media Type\t%s\n", metadata.MediaType); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "Platform\t%s\n", metadata.Platform); err != nil {
		return err
	}
	return writer.Flush()
}
