package cmd

import (
	"fmt"
	"log"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/api"
	"github.com/spf13/cobra"
)

// getVaultCmd represents the get-vault command
var getVaultCmd = &cobra.Command{
	Use:     "get-vault",
	Aliases: []string{"vaults", "vault"},
	Short:   "Fetch open vaults with HLP priority",
	Long: `Fetch and display open vaults with HLP vaults shown first.
    
This command retrieves only open vault data, showing HLP vaults first,
then other vaults, with TVL sorting within each category (descending by default).
Vaults with TVL below the minimum threshold (default: 50,000) are filtered out.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(cfg.BaseURL, cfg.InfoURL)

		items, err := client.FetchAllVault()
		if err != nil {
			log.Fatalf("Error fetching vaults: %v", err)
		}

		// Always filter to show only open vaults
		items = items.FilterOpenVaults()

		// Apply TVL threshold filtering
		minTVL, _ := cmd.Flags().GetFloat64("min-tvl")
		items = items.FilterByMinTVL(minTVL)

		// Apply HLP prioritization with TVL sorting within each category
		sortDesc, _ := cmd.Flags().GetBool("desc")
		items = items.SortWithHLPPriority(!sortDesc)

		count, _ := cmd.Flags().GetInt("count")
		fmt.Println(items.FormatString(count))
	},
}

func init() {
	rootCmd.AddCommand(getVaultCmd)
	getVaultCmd.Flags().IntP("count", "c", 100, "Number of vaults to display")
	getVaultCmd.Flags().Bool("desc", true, "Sort TVL in descending order (within HLP/non-HLP categories)")
	getVaultCmd.Flags().Float64("min-tvl", 50000, "Minimum TVL threshold for vault filtering")
}
