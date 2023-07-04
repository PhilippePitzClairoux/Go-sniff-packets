package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"math"
	"os"
	"packet-sniffer/goprettypackets"
	"packet-sniffer/goresolve"
	"strings"
)

var (
	deviceName         = flag.String("interface", "enp0s31f6", "interface name (ip addr)")
	showLayers         = flag.Bool("show-layers", true, "display layers and their content")
	resolveIpAddresses = flag.Bool("resolve-ip", true, "try finding host name if possible")
	websiteFilter      = flag.String("website-filter", "", "Only show packets that have this website")
)

func main() {
	flag.Parse()
	packets := make(chan gopacket.Packet)
	defer close(packets)

	fmt.Println("Got the following values : ", *deviceName, *showLayers, *resolveIpAddresses, *websiteFilter)
	bufio.NewReader(os.Stdin).ReadLine()

	listener, err := pcap.OpenLive(*deviceName, math.MaxInt32, true, pcap.BlockForever)

	assertError(err, "Could not open interface in live mode.")
	go capturePackets(listener, packets)

	defer listener.Close()
	for packet := range packets {

		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			fmt.Println("Could not get network layer")
			continue
		}

		flow := networkLayer.NetworkFlow()

		go func(packet gopacket.Packet, flow gopacket.Flow) {
			var srcString = flow.Src().String()
			var dstString = flow.Dst().String()
			var displayPacket = len(*websiteFilter) > 0
			resolveIpAddress(&displayPacket, &srcString, &dstString)

			if (len(*websiteFilter) > 0 && displayPacket) || len(*websiteFilter) == 0 {
				fmt.Printf("Got a packet : source => %s , destination => %s\n", srcString, dstString)

				if *showLayers {
					for _, layer := range packet.Layers() {
						displayLayer(&layer)
					}
				}
			}
		}(packet, flow)
	}
}

func resolveIpAddress(displayPacket *bool, srcString *string, dstString *string) {
	if *resolveIpAddresses || *displayPacket {
		src := goresolve.Ip(*srcString)
		dst := goresolve.Ip(*dstString)

		if strings.Contains(src[0], *websiteFilter) || strings.Contains(dst[0], *websiteFilter) {
			*displayPacket = true
		} else {
			*displayPacket = false
		}

		*srcString = fmt.Sprintf("(%s) %v", *srcString, src)
		*dstString = fmt.Sprintf("(%s) %v", *dstString, dst)
	}
}

func displayLayer(layer *gopacket.Layer) {
	layerType := (*layer).LayerType()
	layerContent := (*layer).LayerContents()
	layerPayload := (*layer).LayerPayload()

	layerContentFormatted := goprettypackets.FormatRawPacket(layerContent)
	layerPayloadFormatted := goprettypackets.FormatRawPacket(layerPayload)

	fmt.Printf("\tlayer type => %s\n\tlayer content => %s\n\tlayer payload => %s\n",
		layerType.String(), layerContentFormatted, layerPayloadFormatted)
}

func capturePackets(handle *pcap.Handle, output chan gopacket.Packet) {
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		output <- packet
	}
}

func assertError(_error error, message string) {
	if _error != nil {
		log.Fatal(message)
	}
}
