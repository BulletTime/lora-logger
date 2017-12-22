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
	"time"

	"github.com/apex/log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

var (
	device       string = "eth0"
	snapshot_len int32  = 65535
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")

		//// Find all devices
		//devices, err := pcap.FindAllDevs()
		//if err != nil {
		//	log.WithError(err).Fatal("finding devices")
		//}
		//
		//// Print device information
		//fmt.Println("Devices found:")
		//for _, device := range devices {
		//	fmt.Println("\nName: ", device.Name)
		//	fmt.Println("Description: ", device.Description)
		//	fmt.Println("Devices addresses: ", device.Description)
		//	for _, address := range device.Addresses {
		//		fmt.Println("- IP address: ", address.IP)
		//		fmt.Println("- Subnet mask: ", address.Netmask)
		//	}
		//}

		// Open device
		handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
		if err != nil {
			log.WithError(err).Fatal("open device failed")
		}
		defer handle.Close()

		// Set filter
		var filter string = "udp and port 1700"
		err = handle.SetBPFFilter(filter)
		if err != nil {
			log.WithError(err).Fatal("filter failed")
		}
		fmt.Println("Only capturing UDP port 1700 packets.")

		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			// Process packet here
			fmt.Println(packet)
			data := packet.TransportLayer().LayerPayload()
			fmt.Println(data)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//func udpServer() {
//	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:1700")
//	if err != nil {
//		log.WithError(err).Fatal("failed to resolve udp address")
//	}
//	log.WithField("addr", addr).Info("starting gateway udp listener")
//
//	pc, err := net.ListenUDP("udp", addr)
//	if err != nil {
//		log.WithError(err).Fatal("failed to start listener")
//	}
//	defer pc.Close()
//
//	buf := make([]byte, 65507)
//
//	for {
//		i, addr, err := pc.ReadFromUDP(buf)
//		if err != nil {
//			log.WithError(err).Error("udp read error")
//		}
//		data := make([]byte, i)
//		copy(data, buf[:i])
//		log.WithFields(log.Fields{
//			"addr":         addr,
//			"data":         data,
//			"string(data)": string(data),
//		}).Info("received udp message")
//	}
//}
