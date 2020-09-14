package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scrapey",
	Short: "GoScrapey is a website image scraper",
	Long: `An image scraper build with love by Brizinger in Go
	More information can be found on http://github.com/brizinger/GoScrapey`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}

// Execute - executes the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
