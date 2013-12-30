package models

import (
	"time"
)

type Network struct {
	HashRate    string    "hashRate"
	Difficulty  string    "difficulty"
	Mined       string    "mined"
	BlockCount  string    "blockCount"
	GeneratedAt time.Time "generatedAt"
}

var networkCollection = "network"

// Insert saves a new WDC network snapshot to the database.
func (n *Network) Insert(conn *MgoConnection) error {
	return conn.DB.C(networkCollection).Insert(n)
}
