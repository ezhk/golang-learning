package cmd

import (
	"fmt"
	"os"

	config "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Configuration
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Calendar management tools",
	Long: `Calednar tools that allows different actions:
- calendar: server-side app, that provide gPRC API;
- scheduler: read event from database and notify about upcoming events;
- sender: process notifications and send messages to users.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/calendar.yaml", "Configuration filepath")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	if cfgFile != "" {
		// Process config file and store them into struct.
		cfg = config.NewConfig(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".main" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".main")
	}
}
