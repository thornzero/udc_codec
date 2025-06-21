package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thornzero/udc_codec/pkg/udc"
)

func main() {
	var rootCmd = &cobra.Command{Use: "udccli"}

	var scrapeCmd = &cobra.Command{
		Use:   "scrape",
		Short: "Production-grade full UDC recursive scrape",
		Run: func(cmd *cobra.Command, args []string) {
			err := udc.ScrapeFullHierarchy("data/udc_full.yaml")
			if err != nil {
				panic(err)
			}
			fmt.Println("✅ UDC recursive production scrape complete!")
		},
	}

	var lookupCmd = &cobra.Command{
		Use:   "lookup [code]",
		Short: "Lookup a UDC code",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			codec, err := udc.LoadCodec("data/udc_full.yaml")
			if err != nil {
				fmt.Println("Error loading codec:", err)
				os.Exit(1)
			}
			title, ok := codec.Lookup(args[0])
			if ok {
				fmt.Printf("%s => %s\n", args[0], title)
			} else {
				fmt.Println("Code not found.")
			}
		},
	}

	var addendumCmd = &cobra.Command{
		Use:   "addendum",
		Short: "Manage UDC addendum files",
	}

	var listAddendumsCmd = &cobra.Command{
		Use:   "list",
		Short: "List all addendum files",
		Run: func(cmd *cobra.Command, args []string) {
			am := udc.NewAddendumManager("data")
			addendums, err := am.ListAddendums()
			if err != nil {
				fmt.Println("Error listing addendums:", err)
				os.Exit(1)
			}
			if len(addendums) == 0 {
				fmt.Println("No addendum files found.")
				return
			}
			fmt.Println("Addendum files:")
			for _, addendum := range addendums {
				fmt.Printf("  - %s\n", addendum)
			}
		},
	}

	var addAddendumCmd = &cobra.Command{
		Use:   "add [code] [title] [filename]",
		Short: "Add a classification to an addendum file (filename is optional)",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			code := args[0]
			title := args[1]
			var filename string
			if len(args) == 3 {
				filename = args[2]
			}

			node := &udc.Node{
				Code:  code,
				Title: title,
			}

			am := udc.NewAddendumManager("data")
			err := am.Add(filename, []*udc.Node{node})
			if err != nil {
				fmt.Println("Error adding to addendum:", err)
				os.Exit(1)
			}

			// Determine the actual filename used
			actualFilename := filename
			if actualFilename == "" {
				actualFilename = "default"
			}
			if !strings.HasPrefix(actualFilename, "udc_addendum_") {
				actualFilename = "udc_addendum_" + actualFilename
			}
			if !strings.HasSuffix(actualFilename, ".yaml") {
				actualFilename = actualFilename + ".yaml"
			}

			fmt.Printf("✅ Added classification to addendum file: %s\n", actualFilename)
		},
	}

	var deleteAddendumCmd = &cobra.Command{
		Use:   "delete [filename]",
		Short: "Delete an addendum file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			// Ensure filename has correct format
			if !strings.HasPrefix(filename, "udc_addendum_") {
				filename = "udc_addendum_" + filename
			}
			if !strings.HasSuffix(filename, ".yaml") {
				filename = filename + ".yaml"
			}

			am := udc.NewAddendumManager("data")
			err := am.DeleteAddendum(filename)
			if err != nil {
				fmt.Println("Error deleting addendum:", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Deleted addendum file: %s\n", filename)
		},
	}

	rootCmd.AddCommand(scrapeCmd)
	rootCmd.AddCommand(lookupCmd)

	addendumCmd.AddCommand(listAddendumsCmd)
	addendumCmd.AddCommand(addAddendumCmd)
	addendumCmd.AddCommand(deleteAddendumCmd)
	rootCmd.AddCommand(addendumCmd)

	rootCmd.Execute()
}
