package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Artist struct {
	ID            int
	Name          string
	Aliases       []ArtistAlias
	ReleaseGroups []RoledReleaseGroup
	Added         time.Time
	AddedBy       User
	Bio           sql.NullString
	Tags          []string
}

type ArtistAlias struct {
	Alias   string
	Added   time.Time
	AddedBy User
}

type RoledReleaseGroup struct {
	Role         int
	ReleaseGroup ReleaseGroup
}

func (db *DB) AutocompleteArtists(s string) ([]Artist, error) {
	if len(s) == 0 {
		return nil, errors.New("missing s")
	}
	rows, err := db.db.Query("SELECT DISTINCT a.id,a.name,a.bio,a.added,u.id,u.username FROM artists a, artist_aliases al, users u WHERE a.added_by = u.id AND  (a.name LIKE $1 OR (al.alias LIKE $1 AND al.artist = a.id))", fmt.Sprint("%", s, "%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []Artist
	for rows.Next() {
		var tmp Artist
		err = rows.Scan(&tmp.ID,
			&tmp.Name,
			&tmp.Bio,
			&tmp.Added,
			&tmp.AddedBy.ID,
			&tmp.AddedBy.Username)
		if err != nil {
			return nil, err
		}

		err = db.populateArtistAliases(&tmp)
		if err != nil {
			return nil, err
		}

		err = db.populateArtistTags(&tmp)
		if err != nil {
			return nil, err
		}

		artists = append(artists, tmp)
	}

	return artists, nil
}

func (db *DB) populateArtistTags(a *Artist) error {
	if a == nil {
		return errors.New("missing artist")
	}

	rows, err := db.db.Query("SELECT t.tag FROM artist_tags t,artist_tags_artists a WHERE a.artist = $1 AND a.tag = t.id ", a.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp string
		err = rows.Scan(&tmp)
		if err != nil {
			return err
		}

		a.Tags = append(a.Tags, tmp)
	}

	return nil
}

func (db *DB) populateArtistAliases(a *Artist) error {
	if a == nil {
		return errors.New("missing artist")
	}

	rows, err := db.db.Query("SELECT a.alias,a.added,u.id,u.username FROM artist_aliases a, users u WHERE a.added_by = u.id AND artist=$1", a.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp ArtistAlias
		err = rows.Scan(
			&tmp.Alias,
			&tmp.Added,
			&tmp.AddedBy.ID,
			&tmp.AddedBy.Username)
		if err != nil {
			return err
		}
		a.Aliases = append(a.Aliases, tmp)
	}

	return nil
}

func (db *DB) AutocompleteArtistTags(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errors.New("missing s")
	}

	rows, err := db.db.Query("SELECT tag FROM artist_tags WHERE tag LIKE $1", fmt.Sprint("%", s, "%"))
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

func (db *DB) GetArtist(id int) (*Artist, error) {
	if id < 0 {
		return nil, errors.New("invalid id")
	}
	row := db.db.QueryRow("SELECT a.name,a.bio,a.added,u.id,u.username FROM artists a, users u WHERE a.added_by = u.id AND a.id = $1", id)

	artist := Artist{ID: id}
	err := row.Scan(
		&artist.Name,
		&artist.Bio,
		&artist.Added,
		&artist.AddedBy.ID,
		&artist.AddedBy.Username)
	if err != nil {
		return nil, err
	}

	err = db.populateArtistAliases(&artist)
	if err != nil {
		return nil, err
	}

	err = db.populateArtistTags(&artist)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func (db *DB) PopulateReleaseGroups(artist *Artist) error {
	if artist.ID < 0 {
		return errors.New("invalid artist ID")
	}

	rows, err := db.db.Query("SELECT rga.role,rg.id,rg.name,rg.type,rg.release_date FROM release_groups rg, release_groups_artists rga WHERE rg.id = rga.release_group AND rga.artist = $1", artist.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp RoledReleaseGroup
		err = rows.Scan(
			&tmp.Role,
			&tmp.ReleaseGroup.ID,
			&tmp.ReleaseGroup.Name,
			&tmp.ReleaseGroup.Type,
			&tmp.ReleaseGroup.ReleaseDate)
		if err != nil {
			return err
		}

		// TODO populate tags, possibly artists?

		artist.ReleaseGroups = append(artist.ReleaseGroups, tmp)
	}

	return nil
}

func insertArtistTx(artist *Artist, tx *sql.Tx) error {
	if artist.AddedBy.ID < 0 {
		return errors.New("invalid user ID")
	}

	var bio *string
	if artist.Bio.String != "" {
		bio = &artist.Bio.String
	}

	err := tx.QueryRow("INSERT INTO artists(name,bio,added,added_by) VALUES ($1,$2,$3,$4) RETURNING id", artist.Name, bio, artist.Added, artist.AddedBy.ID).Scan(&artist.ID)
	if err != nil {
		return err
	}

	var res sql.Result
	for _, t := range artist.Tags {
		var id int
		err = tx.QueryRow("INSERT INTO artist_tags(tag) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", t).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == nil {
			// inserted
			res, err = tx.Exec("INSERT INTO artist_tags_artists(artist,tag) VALUES($1,$2)", artist.ID, id)
		} else {
			// already present
			res, err = tx.Exec("INSERT INTO artist_tags_artists(artist,tag) VALUES($1,(SELECT id FROM artist_tags WHERE tag=$2 LIMIT 1))", artist.ID, t)
		}
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return errors.New("did not insert")
		}
	}

	for _, a := range artist.Aliases {
		res, err = tx.Exec("INSERT INTO artist_aliases(artist,alias,added,added_by) VALUES($1,$2,$3,$4)", artist.ID, a.Alias, a.Added, a.AddedBy.ID)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return errors.New("did not insert")
		}
	}

	// TODO insert release groups?
	return nil
}

func (db *DB) InsertArtist(artist *Artist) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = insertArtistTx(artist, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
