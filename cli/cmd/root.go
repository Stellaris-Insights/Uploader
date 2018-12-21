// Copyright © 2018 C45tr0 <william.the.developer+stellaris@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package cmd contains all of the commands available to the stellaris-insights binary.
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/stellaris-insights/uploader"
	"github.com/stellaris-insights/uploader/api"
)

var home string
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stellaris-insights",
	Short: "Game data uploader for Stellaris Insights",
	Long:  `This tool provides a way to easily upload your saves as your playing to stellarisinsights.com`,
	Run: func(cmd *cobra.Command, args []string) {
		userdataDir := ""

		if viper.IsSet("userdata") {
			userdataDir = viper.GetString("userdata")
		} else {
			userdataDir = path.Join(home, "Documents", "Paradox Interactive", "Stellaris")
		}

		fmt.Println(userdataDir)

		err := survey.AskOne(&survey.Input{
			Message: "Stellaris Userdata path: ",
			Default: userdataDir,
		}, &userdataDir, survey.Required)

		if err != nil {
			fmt.Println("Failed to ask question")
			fmt.Println(err)
			os.Exit(1)
		}

		userdataDir, err = filepath.Abs(userdataDir)
		if err != nil {
			fmt.Println("Failed to find absolute path")
			fmt.Println(err)
			os.Exit(1)
		}

		viper.Set("userdata", userdataDir)

		if _, err = os.Stat(path.Join(userdataDir, "save games")); err != nil {
			fmt.Println("Can not find or load save games folder in specified folder")
			fmt.Println(err)
			os.Exit(1)
		}

		isContinuation := false

		err = survey.AskOne(&survey.Confirm{
			Message: "Is this a continuation of a previous upload session?",
		}, &isContinuation, survey.Required)

		if err != nil {
			fmt.Println("Failed to ask question")
			fmt.Println(err)
			os.Exit(1)
		}

		uploadSessionID := ""
		uploadSessionSecret := ""

		if isContinuation {
			err = survey.AskOne(&survey.Input{
				Message: "Upload Session Id",
			}, &uploadSessionID, survey.Required)

			if err != nil {
				fmt.Println("Failed to ask question")
				fmt.Println(err)
				os.Exit(1)
			}

			err = survey.AskOne(&survey.Password{
				Message: "Upload Session Secret",
			}, &uploadSessionSecret, survey.Required)

			if err != nil {
				fmt.Println("Failed to ask question")
				fmt.Println(err)
				os.Exit(1)
			}
		} // TODO: Gen session id and session secret if it is not a continuation

		fmt.Println("****************************************************************")
		fmt.Println("Starting watcher")
		fmt.Printf("Watching directory: %#v\n", path.Join(userdataDir, "save games"))
		fmt.Printf("Upload Session Id: %#v\n", uploadSessionID)
		fmt.Printf("Upload Session Secret: %#v\n", uploadSessionSecret)
		fmt.Println("****************************************************************")

		w := uploader.NewSaveGameWatcher(
			uploader.NewFSNotifyWrapper(),
			uploader.NewS3Uploader(
				api.NewS3ApiService(&http.Client{}, "https://api.stellarisinsights.com/"),
				uploadSessionID,
				uploadSessionSecret,
				userdataDir,
			),
		)
		w.Start(userdataDir)
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stellaris-insights.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		var err error

		// Find home directory.
		home, err = homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".stellaris-insights" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".stellaris-insights")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
