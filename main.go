package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"math"
	"net"
	"packet-sniffer/pkg/goprettypackets"
	"sync"
)

var (
	deviceName         = flag.String("interface", "enp0s31f6", "interface name (ip addr)")
	showLayers         = flag.Bool("show-layers", true, "display layers and their content")
	resolveIpAddresses = flag.Bool("resolve-ip", true, "try finding host name if possible")

	resolvedIp     = make(map[string][]string, 0)
	resolveIpMutex sync.Mutex
)

func main() {
	flag.Parse()
	packets := make(chan gopacket.Packet)
	defer close(packets)

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
			if *resolveIpAddresses {
				src := resolveIp(flow.Src().String())
				dst := resolveIp(flow.Dst().String())

				fmt.Printf("Got a packet : source => (%s) %v , destination => (%s) %v\n",
					flow.Src(), src, flow.Dst(), dst)

			} else {
				fmt.Printf("Got a packet : source => %s , destination => %s\n", flow.Src(), flow.Dst())
			}

			if *showLayers {
				for _, layer := range packet.Layers() {
					displayLayer(&layer, packet)
				}
			}
		}(packet, flow)
	}
}

func displayLayer(layer *gopacket.Layer, packet gopacket.Packet) {
	layerType := (*layer).LayerType()
	layerContent := (*layer).LayerContents()
	layerPayload := (*layer).LayerPayload()

	layerContentFormatted := goprettypackets.FormatRawPacket(layerContent)
	layerPayloadFormatted := goprettypackets.FormatRawPacket(layerPayload)

	fmt.Printf("\tlayer type => %s\n\tlayer content => %s\n\tlayer payload => %s\n",
		layerType.String(), layerContentFormatted, layerPayloadFormatted)
}

func resolveIp(ip string) []string {
	resolveIpMutex.Lock()
	defer resolveIpMutex.Unlock()

	if resolvedIp[ip] != nil {
		return resolvedIp[ip]
	}

	ips, ok := net.LookupAddr(ip)
	if ok != nil {
		fmt.Println(ok)
		ips = []string{ip}
	}

	resolvedIp[ip] = ips
	return ips
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
