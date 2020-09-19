//Package cmd contains all of the commands for the goss CLI tool.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// for flags
	asJSON bool

	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "goss",
		Short: "goss in an AWS SSM Paramter Store manager",
		Long: `goss is used to interact with the AWS SSM Parameter Store in a
variety of helpful ways.

You can interact in bulk through the 'import' sub-command to import parameters
directly from a local file.

You can also interact with paths individually to list, put and delete
parameters.
		`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}
)

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
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goss.toml)")
	rootCmd.PersistentFlags().BoolVar(&asJSON, "json", false, "output as json")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	// Search config in home directory with name ".goss" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".goss")
	// 	viper.SetConfigType("toml")
	// }

	// viper.AutomaticEnv() // read in environment variables that match

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
