// Copyright © 2017 Sven Agneessens <sven.agneessens@gmail.com>
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

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/apex/log"
	cliHandler "github.com/apex/log/handlers/cli"
	multiHandler "github.com/apex/log/handlers/multi"
	textHandler "github.com/apex/log/handlers/text"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logFile *os.File
	verbose bool
	debug bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lora-logger",
	Short: "A logger for traffic from lora packet forwarders",
	Long: `A logger for traffic from lora packet forwarders.

This logger can be configured to monitor network traffic and filter out
the traffic from an active packet forwarder running on the same device.
It will log the protocol messages to a log file and/or standard output.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var logLevel = log.InfoLevel
		var logHandlers []log.Handler

		if verbose {
			logHandlers = append(logHandlers, cliHandler.Default)
		}

		if debug {
			logLevel = log.DebugLevel
		}

		absLogFileLocation, err := filepath.Abs("lora.log")
		if err != nil {
			panic(err)
		}
		logFile, err = os.OpenFile(absLogFileLocation, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		if err == nil {
			logHandlers = append(logHandlers, textHandler.New(logFile))
		}

		log.SetHandler(multiHandler.New(logHandlers...))
		log.SetLevel(logLevel)
	},
	//Run: func(cmd *cobra.Command, args []string) {
	//},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if logFile != nil {
			logFile.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lora-logger.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print everything to standard output")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logs")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("version", "V", false, "print build and version info")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configName := ".lora-logger"

		// Search config in home directory with name ".lora-logger" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(configName)

		// Add standard path to config file
		cfgFile = path.Join(home, string(append([]byte(configName), []byte(".yaml")...)))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
