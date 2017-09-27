package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Release struct {
	ID              int
	ReleaseGroup    ReleaseGroup
	Edition         string
	Medium          int
	ReleaseDate     time.Time
	CatalogueNumber string
	RecordLabel     RecordLabel
	Torrents        []Torrent
	Added           time.Time
	AddedBy         User
	Original        bool
	Tags            []string

	// Properties lists "official" properties of releases, for example
	// "LossyMasterApproved", to be set by trusted users or staff.
	Properties map[string]string
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

func insertReleasePropertiesTx(release Release, tx *sql.Tx) error {
	for k, v := range release.Properties {
		_, err := tx.Exec("INSERT INTO release_properties_releases(release, property, value) VALUES ($1,(SELECT id from release_properties WHERE release_properties.property = $2 LIMIT 1),$3)", release.ID, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertReleaseTagsTx(release Release, tx *sql.Tx) error {
	var res sql.Result
	var err error
	for _, t := range release.Tags {
		var id int
		err = tx.QueryRow("INSERT INTO release_tags(tag) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", t).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == nil {
			// inserted
			res, err = tx.Exec("INSERT INTO release_tags_releases(release,tag) VALUES($1,$2)", release.ID, id)
		} else {
			// already present
			res, err = tx.Exec("INSERT INTO release_tags_releases(release,tag) VALUES($1,(SELECT id FROM release_tags WHERE tag=$2 LIMIT 1))", release.ID, t)
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

func insertReleaseTx(release *Release, tx *sql.Tx) error {
	err := tx.QueryRow("INSERT INTO releases(edition,medium,release_group,record_label,added,added_by,release_date,catalogue_number,original) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id",
		release.Edition,
		release.Medium,
		release.ReleaseGroup.ID,
		release.RecordLabel.ID,
		release.Added,
		release.AddedBy.ID,
		release.ReleaseDate,
		release.CatalogueNumber,
		release.Original).Scan(&release.ID)
	if err != nil {
		return err
	}

	err = insertReleaseTagsTx(*release, tx)
	if err != nil {
		return err
	}

	err = insertReleasePropertiesTx(*release, tx)

	return err
}

func (db *DB) InsertRelease(release *Release) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = insertReleaseTx(release, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) SetReleaseProperty(id int, k, v string) error {
	if id < 0 {
		return errors.New("invalid ID")
	}
	var value *string
	if len(v) != 0 {
		value = &v
	}

	res, err := db.db.Exec("INSERT INTO release_properties_releases(release, property, value) VALUES ($1,(SELECT id from release_properties WHERE release_properties.property = $2 LIMIT 1),$3) ON CONFLICT (release,property) DO UPDATE SET value = $3", id, k, value)
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

	return nil
}

func (db *DB) populateReleaseTags(r *Release) error {
	if r.ID < 0 {
		return errors.New("invalid ID")
	}

	rows, err := db.db.Query("SELECT t.tag FROM release_tags t, release_tags_releases rtr WHERE rtr.tag = t.id AND rtr.release = $1", r.ID)
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
		r.Tags = append(r.Tags, tmp)
	}

	return nil
}

func (db *DB) populateReleaseProperties(r *Release) error {
	if r.ID < 0 {
		return errors.New("invalid ID")
	}

	rows, err := db.db.Query("SELECT p.property,rpr.value FROM release_properties p, release_properties_releases rpr WHERE rpr.property = p.ID AND rpr.release = $1", r.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var (
			tmpK string
			tmpV sql.NullString
		)
		err = rows.Scan(&tmpK, &tmpV)
		if err != nil {
			return err
		}
		if _, ok := m[tmpK]; ok {
			return errors.New("duplicate release property?")
		}
		// If !tmpV.Valid, then tmpV.String == ""
		m[tmpK] = tmpV.String
	}

	r.Properties = m
	return nil
}

func (db *DB) GetRelease(id int) (*Release, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	row := db.db.QueryRow("SELECT r.edition,r.medium,g.id,g.name,l.id,l.name,r.added,u.id,u.username,r.release_date,r.catalogue_number,r.original FROM releases r, release_groups g, record_labels l, users u WHERE r.release_group = g.id AND r.record_label = l.id AND r.added_by = u.id AND r.id = $1", id)
	r := Release{ID: id}
	err := row.Scan(
		&r.Edition,
		&r.Medium,
		&r.ReleaseGroup.ID,
		&r.ReleaseGroup.Name,
		&r.RecordLabel.ID,
		&r.RecordLabel.Name,
		&r.Added,
		&r.AddedBy.ID,
		&r.AddedBy.Username,
		&r.ReleaseDate,
		&r.CatalogueNumber,
		&r.Original)
	if err != nil {
		return nil, err
	}

	err = db.populateReleaseTags(&r)
	if err != nil {
		return nil, err
	}

	err = db.populateReleaseProperties(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func deleteReleasePropertiesTx(id int, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM release_properties_releases WHERE release=$1", id)
	return err
}

func deleteReleaseTagsTx(id int, tx *sql.Tx) error {
	res, err := tx.Exec("DELETE FROM release_tags_releases WHERE release = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected > 0 {
		_, err := tx.Exec("DELETE FROM release_tags t WHERE (SELECT COUNT(*) FROM release_tags_releases rtr WHERE rtr.tag=t.id) = 0")
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteReleaseTx(id int, tx *sql.Tx) error {
	err := deleteReleaseTagsTx(id, tx)
	if err != nil {
		return err
	}

	err = deleteReleasePropertiesTx(id, tx)
	if err != nil {
		return err
	}

	res, err := tx.Exec("DELETE FROM releases WHERE id=$1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("release not found")
	}

	return nil
}

func (db *DB) DeleteRelease(id int) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = deleteReleaseTx(id, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) PopulateTorrents(release *Release) error {
	return nil
}
