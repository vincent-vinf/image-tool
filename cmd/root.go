package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vincent-vinf/image-tool/pkg/utils"
)

var (
	logger        = utils.GetLogger()
	outputDir     string
	imageListPath string
	platform      string
	autoZip       bool

	username string
	password string

	//srcUsername string
	//srcPassword string
	//
	//dstUsername string
	//dstPassword string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "image-tool",
	Short: "A brief description of your application",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&outputDir, "dir", "d", "images", "output dir")
	rootCmd.PersistentFlags().StringVarP(&imageListPath, "images", "i", "", "images.txt path")
	rootCmd.PersistentFlags().StringVar(&platform, "platform", "linux/amd64", "image platform")
	rootCmd.PersistentFlags().BoolVar(&autoZip, "zip", true, "automatically compress image tar using pigz")

	//rootCmd.PersistentFlags().StringVar(&srcUsername, "src-username", "", "source username")
	//rootCmd.PersistentFlags().StringVar(&srcPassword, "src-password", "", "source password")
	//
	//rootCmd.PersistentFlags().StringVar(&dstUsername, "dst-username", "", "destination username")
	//rootCmd.PersistentFlags().StringVar(&dstPassword, "dst-password", "", "destination password")

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
