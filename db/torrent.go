package db

import "time"

type Torrent struct {
	ID         int
	Uploaded   time.Time
	UploadedBy User
	InfoHash   [20]byte

	Format      int
	Size        int64
	Description string

	Leechers int
	Seeders  int
	Snatches int

	FileList []string

	Properties map[string]string
	LeechType  string
}

func (db *DB) InsertTorrent(torrent *Torrent) error {
	return nil
}

func (db *DB) GetTorrent(id int) (*Torrent, error) {
	return &Torrent{}, nil
}

func (db *DB) DeleteTorrent(id int) error {
	return nil
}
