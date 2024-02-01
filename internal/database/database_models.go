package database

import "time"

type DB_Short struct {
	Short      string    `db:"Short"`
	Type       string    `db:"Type"`
	Author     string    `db:"Author"`
	Timestamp  time.Time `db:"Timestamp"`
	Ip         string    `db:"Ip"`
	Restricted bool      `db:"Restricted"`
	Downloads  int       `db:"Downloads"`
}

type DB_File struct {
	Short    string `db:"Short"`
	Filename string `db:"Filename"`
}

type DB_Paste struct {
	Short string `db:"Short"`
	Text  string `db:"Text"`
}
