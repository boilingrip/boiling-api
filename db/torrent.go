package db

import "time"

type Torrent struct {
	ID       int
	Uploader User
	Uploaded time.Time
	InfoHash [20]byte

	Format      string
	Encoding    string // Lossy/Lossless
	Size        int64
	Description string

	Leechers int
	Seeders  int
	Snatches int

	FileList []string

	Properties map[string]interface{}
	LeechType  string
}

func (db *DB) InsertTorrent(torrent *Torrent) error {
	return nil
}

func (db *DB) GetTorrent(id int) (*Torrent, error) {
	return &Torrent{}, nil
}

func (db *DB) RemoveTorrent(id int) error {
	return nil
}
