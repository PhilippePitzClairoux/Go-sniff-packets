package internal

import (
	"fmt"
	"github.com/google/gopacket"
	layers2 "github.com/google/gopacket/layers"
	"log"
	"packet-sniffer/pkg/goprettypackets"
	"packet-sniffer/pkg/goresolve"
	"strconv"
	"strings"
)

type Packet struct {
	Source      string
	Destination string
	Content     string
	Time        string
	Layers      map[string]LayerInformation
}

type LayerInformation struct {
	Payload    string
	RawPayload []byte
	Content    string
	RawContent []byte
}

func NewPacket(source gopacket.Packet) *Packet {
	networkLayer := source.NetworkLayer()
	arpLayer, ok := source.Layer(layers2.LayerTypeARP).(*layers2.ARP)
	var srcString string
	var dstString string

	if networkLayer == nil && ok {
		srcString = fmt.Sprintf("%s|%s",
			byteAddrToString(arpLayer.SourceProtAddress, "."),
			byteAddrToString(arpLayer.SourceHwAddress, ":"),
		)
		dstString = fmt.Sprintf("%s|%s",
			byteAddrToString(arpLayer.DstProtAddress, "."),
			byteAddrToString(arpLayer.DstHwAddress, ":"),
		)
		log.Printf("Using ARP layer instead of network layer : src=%s, dst=%s\n", srcString, dstString)
	} else if networkLayer != nil {
		flow := networkLayer.NetworkFlow()
		srcString = flow.Src().String()
		dstString = flow.Dst().String()
	} else {
		log.Printf("Could not find networkLayer or ARPLayer... %+v\n", source)
	}

	resolveIpAddress(&srcString, &dstString)
	formattedLayers := make([]string, 0)
	layers := make(map[string]LayerInformation)

	for _, layer := range source.Layers() {
		formattedLayers = append(formattedLayers, displayLayer(&layer))
		layers[layer.LayerType().String()] = LayerInformation{
			Payload:    showPrintableCharactersOnly(layer.LayerPayload()),
			RawPayload: layer.LayerPayload(),

			Content:    showPrintableCharactersOnly(layer.LayerContents()),
			RawContent: layer.LayerContents(),
		}
	}

	return &Packet{
		Time:        source.Metadata().Timestamp.String(),
		Source:      srcString,
		Destination: dstString,
		Content: strings.Join(
			formattedLayers,
			"\n",
		),
		Layers: layers,
	}
}

func resolveIpAddress(srcString *string, dstString *string) {
	src := goresolve.Ip(*srcString)
	dst := goresolve.Ip(*dstString)

	*srcString = fmt.Sprintf("(%s) %v", *srcString, src)
	*dstString = fmt.Sprintf("(%s) %v", *dstString, dst)
}

func displayLayer(layer *gopacket.Layer) string {
	layerType := (*layer).LayerType()
	layerContent := (*layer).LayerContents()
	layerPayload := (*layer).LayerPayload()

	layerContentFormatted := goprettypackets.FormatRawPacket(layerContent)
	layerPayloadFormatted := goprettypackets.FormatRawPacket(layerPayload)

	return fmt.Sprintf("\tlayer type => %s\n\tlayer Content =>\n %s\n\tlayer Payload => %s\n",
		layerType.String(), layerContentFormatted, layerPayloadFormatted)
}

func byteAddrToString(bytes []byte, separator string) string {
	var output string
	for i, b := range bytes {
		output += strconv.Itoa(int(b))
		if i != len(bytes)-1 {
			output += separator
		}
	}

	return output
}

func showPrintableCharactersOnly(b []byte) string {
	var buffer string
	for _, char := range b {
		if (char >= 0x21 && char <= 0x7e) || char >= 0xA1 {
			buffer += fmt.Sprintf("%c", char)
		}
	}

	return buffer
}
