package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thornzero/udc_codec/pkg/udc"
)

func main() {
	var rootCmd = &cobra.Command{Use: "udccli"}

	var scrapeCmd = &cobra.Command{
		Use:   "scrape",
		Short: "Production-grade full UDC recursive scrape",
		Run: func(cmd *cobra.Command, args []string) {
			err := udc.ScrapeFullHierarchyChromedpProduction("data/udc_full.yaml")
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

	rootCmd.AddCommand(scrapeCmd)
	rootCmd.AddCommand(lookupCmd)
	rootCmd.Execute()
}
