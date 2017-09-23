package db

import "time"

type ReleaseGroup struct {
	ID          int
	Name        string // Album title
	ReleaseDate time.Time
	ReleaseType string    // Album/EP
	Releases    []Release // individual releases
	Tags        []string
}

func (db *DB) AutocompleteReleaseGroup(s string) ([]ReleaseGroup, error) {
	return []ReleaseGroup{}, nil
}

func (db *DB) AutocompleteReleaseGroupTags(s string) ([]string, error) {
	return []string{}, nil
}

func (db *DB) GetAllReleaseGroupTypes() ([]string, error) {
	return []string{}, nil
}

func (db *DB) GetReleaseGroup(id int) (*ReleaseGroup, error) {
	return &ReleaseGroup{ID: id}, nil
}

func (db *DB) InsertReleaseGroup(group *ReleaseGroup) error {
	// TODO set ID to generated ID
	return nil
}

func (db *DB) PopulateReleases(group *ReleaseGroup) error {
	return nil
}
