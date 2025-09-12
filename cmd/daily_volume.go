package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/api"
	"github.com/spf13/cobra"
)

// dailyVolumeCmd represents the daily-volume command
var dailyVolumeCmd = &cobra.Command{
	Use:     "daily-volume",
	Aliases: []string{"daily", "dvol"},
	Short:   "Fetch daily USD volume data",
	Long: `Fetch and display daily USD volume data.
    
This command retrieves data from the daily_usd_volume endpoint
and displays it in the specified format.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(cfg.BaseURL, cfg.InfoURL)

		// Parse date flags if provided
		var fromDate, toDate *time.Time
		now := time.Now()

		// Check for range flag
		rangeFlag, _ := cmd.Flags().GetString("range")

		// Check for explicit date flags
		fromDateStr, _ := cmd.Flags().GetString("from-date")
		toDateStr, _ := cmd.Flags().GetString("to-date")
		hasExplicitDates := fromDateStr != "" || toDateStr != ""

		// Validate conflicting flags
		if rangeFlag != "" && hasExplicitDates {
			log.Fatal("Error: Cannot use --range flag with explicit --from-date/--to-date flags")
		}

		// Parse range flag if provided
		if rangeFlag != "" {
			var err error
			fromDate, toDate, err = parseRangeFlag(rangeFlag, now)
			if err != nil {
				log.Fatalf("Error parsing range: %v", err)
			}
		} else {
			// Parse explicit date flags if no range flag is used
			if fromDateStr != "" {
				if parsed, err := time.Parse("2006-01-02", fromDateStr); err != nil {
					log.Fatalf("Error parsing from-date: %v. Expected format: YYYY-MM-DD", err)
				} else {
					fromDate = &parsed
				}
			}

			if toDateStr != "" {
				if parsed, err := time.Parse("2006-01-02", toDateStr); err != nil {
					log.Fatalf("Error parsing to-date: %v. Expected format: YYYY-MM-DD", err)
				} else {
					toDate = &parsed
				}
			}

			// Validate date range for explicit dates
			if fromDate != nil && toDate != nil && fromDate.After(*toDate) {
				log.Fatal("Error: from-date must be before or equal to to-date")
			}
		}

		items, err := client.FetchDailyVolume(fromDate, toDate)
		if err != nil {
			log.Fatalf("Error fetching daily volume: %v", err)
		}

		// Apply sorting
		sortOrder, _ := cmd.Flags().GetString("sort")
		var sortDescending bool
		switch strings.ToLower(sortOrder) {
		case "desc", "descending", "":
			sortDescending = true // Default to descending
		case "asc", "ascending":
			sortDescending = false
		default:
			log.Fatalf("Error: Invalid sort order '%s'. Valid options: asc, desc", sortOrder)
		}

		items = items.SortByTime(sortDescending)

		count, _ := cmd.Flags().GetInt("count")
		if fromDate != nil || toDate != nil {
			count = 0
		}
		fmt.Println(items.FormatString(count))
	},
}

func init() {
	rootCmd.AddCommand(dailyVolumeCmd)
	dailyVolumeCmd.Flags().IntP("count", "c", 25, "Number of daily volume entries to display")
	dailyVolumeCmd.Flags().String("from-date", "", "Start date for filtering (YYYY-MM-DD format)")
	dailyVolumeCmd.Flags().String("to-date", "", "End date for filtering (YYYY-MM-DD format)")
	dailyVolumeCmd.Flags().StringP("range", "r", "", "Time range for filtering (e.g., 7D, 30D, 3M, 1Y)")
	dailyVolumeCmd.Flags().StringP("sort", "s", "desc", "Sort order for time: asc (ascending) or desc (descending)")
}
