package db

import (
	"errors"
	"fmt"
	"time"
)

type Release struct {
	ID              int
	Name            string
	Edition         string
	Media           string
	ReleaseDate     time.Time
	CatalogueNumber string
	RecordLabel     RecordLabel
	Torrents        []Torrent
	Added           time.Time
	Original        bool
	Properties      map[string]interface{}
}

func (db *DB) AutocompleteReleaseTags(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errors.New("missing s")
	}

	rows, err := db.db.Query("SELECT tag FROM release_tags WHERE tag LIKE $1", fmt.Sprint("%", s, "%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tmp string
		err = rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tmp)
	}

	return tags, nil
}

func (db *DB) InsertRelease(release *Release) error {
	// TODO set ID to generated ID
	return nil
}

func (db *DB) GetRelease(id int) (*Release, error) {
	return &Release{}, nil
}

func (db *DB) DeleteRelease(id int) error {
	return nil
}

func (db *DB) PopulateFormats(release *Release) error {
	return nil
}
