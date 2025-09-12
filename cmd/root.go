package cmd

import (
	"fmt"
	"os"

	"github.com/LampardNguyen234/hyperliquid-stats/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "hyperliquid-stats",
	Aliases: []string{"hype-stats"},
	Short:   "A CLI tool for fetching and analyzing volume data",
	Long: `Hyperliquid Stats is a command-line tool for fetching and analyzing
cryptocurrency volume data from various endpoints.

It provides multiple output formats (table, JSON, CSV) and supports
multiple data sources for comprehensive volume analysis.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hype-stats.yaml)")
	rootCmd.PersistentFlags().StringP("base-url", "b", config.DefaultBaseURL, "Base URL for the API")
	rootCmd.PersistentFlags().StringP("info-url", "i", config.DefaultInfoURL, "Info URL for the API")

	// Bind flags to viper
	viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("info_url", rootCmd.PersistentFlags().Lookup("info-url"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".volume-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Initialize config
	cfg = config.New()
}
