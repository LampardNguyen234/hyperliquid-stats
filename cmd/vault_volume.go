package cmd

import (
	"fmt"
	"log"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/api"
	"github.com/spf13/cobra"
)

// vaultVolumeCmd represents the vault-volume command
var vaultVolumeCmd = &cobra.Command{
	Use:     "vault-volume",
	Aliases: []string{"vault-vol", "vvol"},
	Short:   "Fetch vault volume information",
	Long: `Fetch and display volume information for vaults.
    
This command retrieves volume data for all vaults or a specific vault by address.
Volume data includes day, week, month, and all-time statistics for both regular and perpetual trading.
Results are sorted by HLP priority first, then by the specified field descending.
Use --hlp flag to show only HLP vault volumes.
Use --count flag to limit the number of vaults displayed.
Use --workers flag to control concurrent fetching (default: 5 workers).
Use --sort-by flag to sort by tvl, day, week, month, or all-time (default: tvl).
Use --summary flag to display aggregated totals and top 10 vaults by TVL.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(cfg.BaseURL, cfg.InfoURL)

		address, _ := cmd.Flags().GetString("address")

		if address != "" {
			// Fetch specific vault volume
			volume, err := client.FetchVaultVolume(address)
			if err != nil {
				log.Fatalf("Error fetching vault volume for address %s: %v", address, err)
			}

			// Get vault name from vault list (optional, fallback to address)
			vaults, _ := client.FetchAllVault()
			vaultName := address
			for _, vault := range vaults {
				if vault.Data.Address == address {
					vaultName = vault.Data.Name
					break
				}
			}

			fmt.Println(volume.FormatSingle(vaultName, address))
		} else {
			// Get count for display limiting
			count, _ := cmd.Flags().GetInt("count")

			// Apply HLP filtering if requested
			hlpOnly, _ := cmd.Flags().GetBool("hlp")

			// Get number of workers
			workers, _ := cmd.Flags().GetInt("workers")

			// Fetch all vault volumes concurrently
			volumes, err := client.FetchAllVaultVolumesConcurrent(hlpOnly, count, workers)
			if err != nil {
				log.Fatalf("Error fetching all vault volumes: %v", err)
			}

			// Check if summary mode is requested
			summaryMode, _ := cmd.Flags().GetBool("summary")
			if summaryMode {
				fmt.Println(volumes.FormatSummary())
			} else {
				// Apply sorting by specified field
				sortBy, _ := cmd.Flags().GetString("sort-by")
				volumes = volumes.SortByField(sortBy)

				fmt.Println(volumes.FormatString())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultVolumeCmd)
	vaultVolumeCmd.Flags().String("address", "", "Specific vault address to fetch volume for")
	vaultVolumeCmd.Flags().Bool("hlp", false, "Show only HLP vaults")
	vaultVolumeCmd.Flags().IntP("count", "c", 0, "Number of vaults to display (0 for all)")
	vaultVolumeCmd.Flags().IntP("workers", "w", 5, "Number of concurrent workers for fetching vault volumes")
	vaultVolumeCmd.Flags().String("sort-by", "tvl", "Sort results by: tvl, day, week, month, all-time (HLP vaults always first)")
	vaultVolumeCmd.Flags().Bool("summary", false, "Display summary of vault volumes (totals by HLP/non-HLP and top 10 TVL)")
}
