package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"gosniff/internal/database"
	tui2 "gosniff/internal/tui"
	"log"
	"math"
	"os"
)

var (
	captureOnly = flag.Bool("tui", false, "Enable/Disable TUI (if tui == true, packet viewing and searching will be possible)")
)

func main() {
	flag.Parse()
	packets := make(chan gopacket.Packet)
	stopChannel := make(chan struct{})
	pckPreview := make(chan string)

	defer close(packets)
	defer close(stopChannel)
	defer close(pckPreview)

	// Setup pcap pcapHandler
	pcapHandler, err := pcap.OpenLive("any", math.MaxInt32, true, pcap.BlockForever)
	assertError(err, "Could not open interface in live mode.")

	// Setup logfile after we have the confirmation we can run command as SUDO/admin/root
	logFile, err := os.OpenFile("go-sniff.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	assertError(err, "Could not open/create/append to log file")

	defer logFile.Close()
	log.SetOutput(logFile)

	go capturePackets(pcapHandler, packets, stopChannel)
	go database.StorePackets(packets, pckPreview, stopChannel)

	var tui tea.Model
	if *captureOnly {
		tui = tui2.NewBrowsePacketTui()
	} else {
		tui = tui2.NewPacketInfinitSpinner(pckPreview)
	}

	var p tea.Model

	if p, err = tea.NewProgram(tui).Run(); err != nil {
		log.Printf("Error running program: %s\n", err)
		fmt.Println("Could run program properly. Please view go-sniff.log for more details")
	}

	fmt.Println("Closing channels that consume/generate packets...")
	fmt.Printf(
		"We've received and stored %d packets! Please view packetdatabase.db in order to search/view transmitted data.",
		p.(tui2.SpinnerModel).GetCount(),
	)
}

func capturePackets(handle *pcap.Handle, output chan<- gopacket.Packet, stop chan struct{}) {
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetChannel := packetSource.Packets()

	for {
		select {
		case <-stop:
			log.Printf("Stopping internal sending due to stop signal...")
			return
		case packet := <-packetChannel:
			output <- packet
		}
	}
}

func assertError(_error error, message string) {
	if _error != nil {
		log.Fatal(message, _error)
	}
}
