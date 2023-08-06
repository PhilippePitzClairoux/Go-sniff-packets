package database

import (
	"database/sql"
	"fmt"
	"github.com/google/gopacket"
	"log"
	"packet-sniffer/internal/internal"
)

const (
	InsertPacketQuery = `INSERT INTO packets (time, source, destination, content) VALUES (?, ?, ?, ?)`
	InsertLayerQuery  = `INSERT INTO packet_layers (packet_id,layer_type, content, content_blob, payload, payload_blob) VALUES (?, ?, ?, ?, ?, ?)`
	FetchLastRows     = `SELECT time, source, destination, content FROM packets ORDER BY time LIMIT 100`
	CountPackets      = `SELECT COUNT(*) FROM packets WHERE time >= ?`
)

func StorePackets(packets <-chan gopacket.Packet, preview chan<- string, stop chan struct{}) {
	for {
		select {
		case <-stop:
			log.Printf("Stopping internal consumption due to stop signal...")
			return
		case packet := <-packets:
			formattedPacket := internal.NewPacket(packet)
			if formattedPacket != nil {
				preview <- fmt.Sprintf("%s -> %s", formattedPacket.Source, formattedPacket.Destination)
				err := InsertPacket(formattedPacket)
				if err != nil {
					log.Printf("Could not insert internal to database : %s\n", err)
				}
			}

		}
	}
}

func InsertPacket(pck *internal.Packet) error {
	db, err := sql.Open("sqlite3", "packetdatabase.db")
	if err != nil {
		return err
	}
	defer db.Close()

	id, err := insertPacket(pck, db)
	if err != nil {
		return err
	}

	for key, value := range pck.Layers {
		err = insertLayer(id, key, value, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertLayer(packetId int64, layerType string, value internal.LayerInformation, db *sql.DB) error {
	prepare, err := db.Prepare(InsertLayerQuery)
	if err != nil {
		return err
	}

	_, err = prepare.Exec(packetId, layerType, value.Content, value.RawContent, value.Payload, value.RawPayload)
	if err != nil {
		return err
	}

	return nil
}

func insertPacket(pck *internal.Packet, db *sql.DB) (int64, error) {
	prepare, err := db.Prepare(InsertPacketQuery)
	if err != nil {
		return 0, err
	}

	res, err := prepare.Exec(pck.Time, pck.Source, pck.Destination, pck.Content)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func FetchLastInsertedPackets() ([]internal.Packet, error) {
	db, err := sql.Open("sqlite3", "packetdatabase.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(FetchLastRows)
	if err != nil {
		return nil, err
	}

	packets := make([]internal.Packet, 100)
	for rows.Next() {
		var time, source, destination, content string
		err = rows.Scan(&time, &source, &destination, &content)
		if err != nil {
			return nil, err
		}

		packets = append(packets, internal.Packet{
			Time:        time,
			Source:      source,
			Destination: destination,
			Content:     content,
		})
	}
	return packets, nil
}

func GetInsertedPackets(t string) (int, error) {
	db, err := sql.Open("sqlite3", "packetdatabase.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var val int
	err = db.QueryRow(CountPackets, t).Scan(&val)
	if err != nil {
		return 0, err
	}

	return val, nil
}
