package cmd

import (
	"fmt"
	"log"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/api"
	"github.com/spf13/cobra"
)

// largestTradeCountCmd represents the largest-trade-count command
var largestTradeCountCmd = &cobra.Command{
	Use:     "largest-trade-count",
	Aliases: []string{"largest-trades", "ltc"},
	Short:   "Fetch largest users by trade count",
	Long: `Fetch and display the largest users by trade count.
    
This command retrieves data from the largest_users_by_trade_count endpoint
and displays it in the specified format.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(cfg.BaseURL, cfg.InfoURL)

		items, err := client.FetchLargestTradeCounts()
		if err != nil {
			log.Fatalf("Error fetching largest trade counts: %v", err)
		}

		count, _ := cmd.Flags().GetInt("count")
		fmt.Println(items.FormatString(count))
	},
}

func init() {
	rootCmd.AddCommand(largestTradeCountCmd)
	largestTradeCountCmd.Flags().IntP("count", "c", 25, "Number of largest trade count users to display")
}
