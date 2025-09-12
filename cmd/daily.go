package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/api"
	"github.com/spf13/cobra"
)

// parseRangeFlag parses range strings like "7D", "1M", "30D", "1Y" etc.
// Returns the fromDate and toDate for the specified range
func parseRangeFlag(rangeStr string, now time.Time) (*time.Time, *time.Time, error) {
	if rangeStr == "" {
		return nil, nil, nil
	}

	// Normalize to uppercase
	rangeStr = strings.ToUpper(strings.TrimSpace(rangeStr))

	// Regex to match patterns like 7D, 1M, 30D, 1Y
	re := regexp.MustCompile(`^(\d+)([DMYW])$`)
	matches := re.FindStringSubmatch(rangeStr)

	if len(matches) != 3 {
		return nil, nil, fmt.Errorf("invalid range format '%s'. Expected formats: 1D, 7D, 30D, 1M, 3M, 1Y", rangeStr)
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid number in range '%s': %v", rangeStr, err)
	}

	unit := matches[2]
	var fromTime time.Time

	switch unit {
	case "D": // Days
		fromTime = now.AddDate(0, 0, -num)
	case "M": // Months
		fromTime = now.AddDate(0, -num, 0)
	case "Y": // Years
		fromTime = now.AddDate(-num, 0, 0)
	case "W": // Weeks
		fromTime = now.AddDate(0, 0, -num*7)
	default:
		return nil, nil, fmt.Errorf("unsupported time unit '%s'. Supported units: D (days), W (weeks), M (months), Y (years)", unit)
	}

	return &fromTime, &now, nil
}

// dailyCmd represents the daily command
var dailyCmd = &cobra.Command{
	Use:     "daily-volume-by-user",
	Aliases: []string{"daily-by-user", "duvol"},
	Short:   "Fetch daily USD volume for a specific user",
	Long: `Fetch and display daily USD volume data for a specific user.
    
This command retrieves data from the daily_usd_volume_by_user endpoint
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

		// Get user filter if provided
		userFilter, _ := cmd.Flags().GetString("user")

		items, err := client.FetchDailyVolumeByUser(fromDate, toDate, userFilter)
		if err != nil {
			log.Fatalf("Error fetching daily volume for user: %v", err)
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
	rootCmd.AddCommand(dailyCmd)
	dailyCmd.Flags().IntP("count", "c", 25, "Number of daily volume entries to display")
	dailyCmd.Flags().String("from-date", "", "Start date for filtering (YYYY-MM-DD format)")
	dailyCmd.Flags().String("to-date", "", "End date for filtering (YYYY-MM-DD format)")
	dailyCmd.Flags().StringP("range", "r", "", "Time range for filtering (e.g., 7D, 30D, 3M, 1Y)")
	dailyCmd.Flags().StringP("sort", "s", "desc", "Sort order for time: asc (ascending) or desc (descending)")
	dailyCmd.Flags().StringP("user", "u", "", "Filter data for a specific user")
}
