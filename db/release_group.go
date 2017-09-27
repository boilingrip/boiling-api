package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type ReleaseGroup struct {
	ID          int
	Name        string // Album title
	Artists     []RoledArtist
	ReleaseDate time.Time
	Added       time.Time
	AddedBy     User
	Type        int       // Album/EP
	Releases    []Release // individual releases
	Tags        []string
}

type RoledArtist struct {
	Role   int
	Artist Artist
}

func (db *DB) AutocompleteReleaseGroups(s string) ([]ReleaseGroup, error) {
	if len(s) == 0 {
		return nil, errors.New("missing s")
	}

	rows, err := db.db.Query("SELECT rg.id,rg.name,rg.release_date,rg.added,u.id,u.username,rg.type FROM release_groups rg, users u WHERE rg.added_by = u.id AND rg.name ILIKE $1", fmt.Sprint("%", s, "%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []ReleaseGroup
	for rows.Next() {
		var group ReleaseGroup
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.ReleaseDate,
			&group.Added,
			&group.AddedBy.ID,
			&group.AddedBy.Username,
			&group.Type)
		if err != nil {
			return nil, err
		}

		err = db.populateReleaseGroupTags(&group)
		if err != nil {
			return nil, err
		}

		err = db.populateReleaseGroupArtists(&group)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (db *DB) AutocompleteReleaseGroupTags(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errors.New("missing s")
	}

	rows, err := db.db.Query("SELECT tag FROM release_group_tags WHERE tag LIKE $1", fmt.Sprint("%", s, "%"))
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

func (db *DB) populateReleaseGroupTags(group *ReleaseGroup) error {
	if group.ID < 0 {
		return errors.New("invalid group ID")
	}

	rows, err := db.db.Query("SELECT rgt.tag FROM release_group_tags rgt, release_group_tags_release_groups rgtrg WHERE rgt.id = rgtrg.tag AND rgtrg.release_group = $1", group.ID)
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

		group.Tags = append(group.Tags, tmp)
	}

	return nil
}

func (db *DB) populateReleaseGroupArtists(group *ReleaseGroup) error {
	if group.ID < 0 {
		return errors.New("invalid  group ID")
	}

	rows, err := db.db.Query("SELECT rga.role,a.id,a.name,a.bio FROM release_groups_artists rga, artists a WHERE rga.artist = a.id AND rga.release_group = $1", group.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var artist RoledArtist
		err = rows.Scan(
			&artist.Role,
			&artist.Artist.ID,
			&artist.Artist.Name,
			&artist.Artist.Bio)
		if err != nil {
			return err
		}

		// TODO populate artist tags, aliases, release groups?

		group.Artists = append(group.Artists, artist)
	}

	return nil
}

func (db *DB) GetReleaseGroup(id int) (*ReleaseGroup, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	res := db.db.QueryRow("SELECT rg.name,rg.release_date,rg.added,u.id,u.username,rg.type FROM release_groups rg, users u WHERE rg.added_by = u.id AND rg.id = $1", id)
	group := ReleaseGroup{ID: id}
	err := res.Scan(
		&group.Name,
		&group.ReleaseDate,
		&group.Added,
		&group.AddedBy.ID,
		&group.AddedBy.Username,
		&group.Type)
	if err != nil {
		return nil, err
	}

	err = db.populateReleaseGroupTags(&group)
	if err != nil {
		return nil, err
	}

	err = db.populateReleaseGroupArtists(&group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func insertReleaseGroupTagsTx(group ReleaseGroup, tx *sql.Tx) error {
	var res sql.Result
	for _, t := range group.Tags {
		var id int
		err := tx.QueryRow("INSERT INTO release_group_tags(tag) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", t).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == nil {
			// inserted
			res, err = tx.Exec("INSERT INTO release_group_tags_release_groups(release_group,tag) VALUES($1,$2)", group.ID, id)
		} else {
			// already present
			res, err = tx.Exec("INSERT INTO release_group_tags_release_groups(release_group,tag) VALUES($1,(SELECT id FROM release_group_tags WHERE tag=$2 LIMIT 1))", group.ID, t)
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

	return nil
}

func insertReleaseGroupArtistsTx(group ReleaseGroup, tx *sql.Tx) error {
	for _, a := range group.Artists {
		res, err := tx.Exec("INSERT INTO release_groups_artists(release_group,artist,role) VALUES($1,$2,$3)", group.ID, a.Artist.ID, a.Role)
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
	return nil
}

func insertReleaseGroupTx(group *ReleaseGroup, tx *sql.Tx) error {
	err := tx.QueryRow("INSERT INTO release_groups(name,release_date,type,added,added_by) VALUES ($1,$2,$3,$4,$5) RETURNING id", group.Name, group.ReleaseDate, group.Type, group.Added, group.AddedBy.ID).Scan(&group.ID)
	if err != nil {
		return err
	}

	err = insertReleaseGroupTagsTx(*group, tx)
	if err != nil {
		return err
	}

	err = insertReleaseGroupArtistsTx(*group, tx)
	if err != nil {
		return err
	}

	// TODO insert releases?
	return nil
}

func (db *DB) InsertReleaseGroup(group *ReleaseGroup) error {
	if group.AddedBy.ID < 0 {
		return errors.New("invalid user ID")
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = insertReleaseGroupTx(group, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) PopulateReleases(group *ReleaseGroup) error {
	if group.ID < 0 {
		return errors.New("invalid ID")
	}

	rows, err := db.db.Query("SELECT r.id,r.edition,r.medium,r.release_date,r.catalogue_number,l.id,l.name,r.added,u.id,u.username,r.original FROM releases r, record_labels l, users u WHERE r.record_label = l.id AND r.added_by = u.id AND r.release_group = $1", group.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp Release
		err = rows.Scan(
			&tmp.ID,
			&tmp.Edition,
			&tmp.Medium,
			&tmp.ReleaseDate,
			&tmp.CatalogueNumber,
			&tmp.RecordLabel.ID,
			&tmp.RecordLabel.Name,
			&tmp.Added,
			&tmp.AddedBy.ID,
			&tmp.AddedBy.Username,
			&tmp.Original)
		if err != nil {
			return err
		}

		err = db.populateReleaseTags(&tmp)
		if err != nil {
			return err
		}

		err = db.populateReleaseProperties(&tmp)
		if err != nil {
			return err
		}

		group.Releases = append(group.Releases, tmp)
	}

	return nil
}
