package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type RecordLabel struct {
	ID          int
	Name        string
	Description sql.NullString
	Founded     pq.NullTime
	Added       time.Time
	AddedBy     User
}

func (db *DB) AutocompleteRecordLabels(s string) ([]RecordLabel, error) {
	if len(s) == 0 {
		return nil, errors.New("misssing s")
	}

	rows, err := db.db.Query("SELECT l.id,l.name,l.description,l.founded,l.added,l.added_by,u.username FROM record_labels l, users u WHERE u.id = l.added_by AND l.name ILIKE $1", fmt.Sprint("%", s, "%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []RecordLabel
	for rows.Next() {
		var tmp RecordLabel
		err = rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Description,
			&tmp.Founded,
			&tmp.Added,
			&tmp.AddedBy.ID,
			&tmp.AddedBy.Username)
		if err != nil {
			return nil, err
		}
		labels = append(labels, tmp)
	}

	return labels, nil
}

func (db *DB) GetRecordLabel(id int) (*RecordLabel, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	row := db.db.QueryRow("SELECT l.name,l.description,l.founded,l.added,l.added_by,u.username FROM record_labels l, users u WHERE u.id = l.added_by AND l.id = $1", id)

	label := RecordLabel{ID: id}
	err := row.Scan(
		&label.Name,
		&label.Description,
		&label.Founded,
		&label.Added,
		&label.AddedBy.ID,
		&label.AddedBy.Username)
	if err != nil {
		return nil, err
	}

	return &label, nil
}

func (db *DB) InsertRecordLabel(label *RecordLabel) error {
	if label == nil {
		return errors.New("missing label")
	}
	var desc *string
	if label.Description.String != "" {
		desc = &label.Description.String
	}
	var founded *time.Time
	var t time.Time
	if label.Founded.Time != t {
		founded = &label.Founded.Time
	}

	err := db.db.QueryRow("INSERT INTO record_labels (name,description,founded,added,added_by) VALUES ($1,$2,$3,now(),$4) RETURNING id",
		label.Name, desc, founded, label.AddedBy.ID).Scan(&label.ID)

	return err
}
