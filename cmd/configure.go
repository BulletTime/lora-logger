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
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/google/gopacket/pcap"
	"github.com/segmentio/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Device      string `yaml:"device"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Promiscuous bool   `yaml:"promiscuous"`
	Timeout     int    `yaml:"timeout"`
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure lora-logger",
	Long: `lora-logger configure creates a yaml configuration file for logging
the traffic from the active packet forwarder.

Various different values for settings that are needed to acquire the traffic from
the packet forwarder are asked to choose.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			newDevice      string
			newHost        string
			newPort        int
			newPromiscuous bool
			newTimeout     int
			err            error
		)

		// Find all devices
		devices, err := pcap.FindAllDevs()
		if err == nil {
			var deviceList []string
			for _, device := range devices {
				deviceList = append(deviceList, device.Name)
			}
			newDevice = deviceList[prompt.Choose("device", deviceList)]
		} else {
			log.WithError(err).Error("failed setting device")
		}

		newHost = prompt.String("server hostname/ip [empty for any]")

		portS := prompt.StringRequired("server port [0 for any]")
		newPort, err = strconv.Atoi(portS)
		if err != nil {
			log.WithField("new port", portS).WithError(err).Warn("failed setting port (is it an integer?)")
		}

		newPromiscuous = prompt.Confirm("enable promiscuous mode (enable only to experiment) [yes/no]")

		timeoutS := prompt.StringRequired("capture packets every x seconds [-1 for continuous]")
		newTimeout, err = strconv.Atoi(timeoutS)
		if err != nil {
			log.WithField("new timeout", timeoutS).WithError(err).Warn("failed setting timeout (is it an integer?)")
		}

		newConfig := &yamlConfig{
			Device:      newDevice,
			Host:        newHost,
			Port:        newPort,
			Promiscuous: newPromiscuous,
			Timeout:     newTimeout,
		}

		output, err := yaml.Marshal(newConfig)
		if err != nil {
			log.WithError(err).Error("failed generating yaml config")
			os.Exit(1)
		}

		if len(viper.ConfigFileUsed()) == 0 {
			viper.SetConfigFile(cfgFile)
		}

		f, err := os.Create(viper.ConfigFileUsed())
		if err != nil {
			log.WithError(err).Error("failed creating log file")
			os.Exit(1)
		}

		defer f.Close()

		f.Write(output)
		log.WithField("path", viper.ConfigFileUsed()).Debug("new configuration file saved")
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
