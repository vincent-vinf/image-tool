package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"image-tool/pkg/image"
	"image-tool/pkg/utils"
)

var (
	logger        = utils.GetLogger()
	outputDir     string
	imageListPath string
	imageHandler  image.Handler
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
	_ = viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))

	rootCmd.PersistentFlags().StringVarP(&imageListPath, "images", "i", "images.txt", "images.txt path")
	_ = viper.BindPFlag("images", rootCmd.Flags().Lookup("images"))

	imageHandler = image.DockerHandler{}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
