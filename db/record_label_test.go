package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestInsertGetRecordLabel(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	l1 := RecordLabel{
		Name:    "NONESUCH",
		AddedBy: User{ID: 1},
	}

	l2 := RecordLabel{
		Name:        "mau5trap",
		AddedBy:     User{ID: 0},
		Description: sql.NullString{String: "Label for electronic music, founded by deadmau5"},
		Founded:     pq.NullTime{Time: time.Date(2007, 01, 01, 0, 0, 0, 0, time.FixedZone("", 0))},
	}

	err = db.InsertRecordLabel(&l1)
	require.Nil(t, err)

	err = db.InsertRecordLabel(&l2)
	require.Nil(t, err)

	l1r, err := db.GetRecordLabel(l1.ID)
	require.Nil(t, err)
	require.Equal(t, l1.Name, l1r.Name)
	require.Equal(t, l1.AddedBy.ID, l1r.AddedBy.ID)
	require.False(t, l1r.Founded.Valid)
	require.False(t, l1r.Description.Valid)

	l2r, err := db.GetRecordLabel(l2.ID)
	require.Nil(t, err)
	require.Equal(t, l2.Name, l2r.Name)
	require.Equal(t, l2.AddedBy.ID, l2r.AddedBy.ID)
	require.True(t, l2r.Founded.Valid)
	require.Equal(t, l2.Founded.Time, l2r.Founded.Time)
	require.True(t, l2r.Description.Valid)
	require.Equal(t, l2.Description.String, l2r.Description.String)
}

func TestAutocompleteRecordLabels(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	l1 := RecordLabel{
		Name:    "NONESUCH",
		AddedBy: User{ID: 1},
	}

	l2 := RecordLabel{
		Name:        "mau5trap",
		AddedBy:     User{ID: 0},
		Description: sql.NullString{String: "Label for electronic music, founded by deadmau5"},
		Founded:     pq.NullTime{Time: time.Date(2007, 01, 01, 0, 0, 0, 0, time.FixedZone("", 0))},
	}

	err = db.InsertRecordLabel(&l1)
	require.Nil(t, err)

	err = db.InsertRecordLabel(&l2)
	require.Nil(t, err)

	labels, err := db.AutocompleteRecordLabels("such")
	require.Nil(t, err)
	require.Equal(t, 1, len(labels))
	require.Equal(t, l1.Name, labels[0].Name)

	labels, err = db.AutocompleteRecordLabels("u")
	require.Nil(t, err)
	require.Equal(t, 2, len(labels))
}
