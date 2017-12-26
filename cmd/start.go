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
	"bytes"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/bullettime/lora-logger/protocol"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start logging",
	Long: `lora-logger start filters the network traffic with the predefined settings (or default)
from the active packet forwarder and logs this traffic to a log file and/or standard output.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			device            = viper.GetString("device")
			host              = viper.GetString("host")
			port              = viper.GetInt("port")
			snapshotLen int32 = 65535
			promiscuous       = viper.GetBool("promiscuous")
			timeout           = time.Duration(viper.GetInt("timeout")) * time.Second
			handle      *pcap.Handle
		)
		log.WithFields(log.Fields{
			"device":      device,
			"host":        host,
			"port":        port,
			"promiscuous": promiscuous,
			"timeout":     timeout,
		}).Debug("loaded settings")

		// Open device
		handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
		if err != nil {
			log.WithError(err).Fatal("open device failed")
		}
		defer handle.Close()

		// Set filter
		var buffer bytes.Buffer
		buffer.WriteString("udp")
		if len(host) > 0 {
			buffer.WriteString(" and host ")
			buffer.WriteString(host)
		}
		if port != 0 {
			buffer.WriteString(" and port ")
			buffer.WriteString(strconv.Itoa(port))
		}
		filter := buffer.String()
		log.WithField("filter", filter).Debug("constructed filter")
		err = handle.SetBPFFilter(filter)
		if err != nil {
			log.WithError(err).Fatal("filter failed")
		}

		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			data := packet.TransportLayer().LayerPayload()
			packet, err := protocol.HandlePacket(data)
			if err != nil {
				ctx := log.WithField("data", data)
				ctx.WithError(err).Error("protocol error")
			} else {
				packet.Log(log.Log)
			}

			//switch p := packet.(type) {
			//case protocol.PushDataPacket:
			//}
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.SetDefault("device", "eth0")
	viper.SetDefault("promiscuous", false)
	viper.SetDefault("timeout", -1)
}
