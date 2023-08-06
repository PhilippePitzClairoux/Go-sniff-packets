package database

import (
	"database/sql"
	_ "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func init() {
	database, err := sql.Open("sqlite3", "packetdatabase.db")
	defer database.Close()
	if err != nil {
		log.Fatal("Could not open/create database : ", err)
	}

	//create various tables
	_, err = database.Exec(`
	CREATE TABLE IF NOT EXISTS packets (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    time TIMESTAMP,
	    source VARCHAR(4096),
	    destination VARCHAR(4096),
	    content TEXT
	);
	
	CREATE TABLE IF NOT EXISTS packet_layers (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    packet_id INTEGER,
	    layer_type VARCHAR(100),
	    content TEXT,
	    content_blob BLOB,
	    payload TEXT,
	    payload_blob BLOB,
	    FOREIGN KEY (packet_id) REFERENCES packets(id)
	)

`)
	if err != nil {
		log.Fatal("Could not create/update database ", err)
	}
}
