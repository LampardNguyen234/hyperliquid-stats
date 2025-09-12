package cmd

import (
	"fmt"
	"log"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/api"
	"github.com/spf13/cobra"
)

// largestCmd represents the largest command
var largestCmd = &cobra.Command{
	Use:     "largest-volume",
	Aliases: []string{"large-vol", "lvol"},
	Short:   "Fetch largest users by USD volume",
	Long: `Fetch and display the largest users by USD volume.
    
This command retrieves data from the largest_users_by_usd_volume endpoint
and displays it in the specified format.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(cfg.BaseURL, cfg.InfoURL)

		items, err := client.FetchLargestUsers()
		if err != nil {
			log.Fatalf("Error fetching largest users: %v", err)
		}

		count, _ := cmd.Flags().GetInt("count")
		fmt.Println(items.FormatString(count))
	},
}

func init() {
	rootCmd.AddCommand(largestCmd)

	largestCmd.Flags().IntP("count", "c", 25, "Number of largest users to display")
}
