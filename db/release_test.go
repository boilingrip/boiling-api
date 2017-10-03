package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsertGetDeleteRelease(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	l := RecordLabel{
		Name:    "NONESUCH",
		AddedBy: User{ID: 1},
	}

	err = db.InsertRecordLabel(&l)
	require.Nil(t, err)

	a := Artist{
		Name:    "deadmau5",
		Added:   time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	err = db.InsertArtist(&a)
	require.Nil(t, err)

	g := ReleaseGroup{
		Name: "4x4=12",
		Artists: []RoledArtist{
			{
				Role:   0,
				Artist: Artist{ID: a.ID},
			},
		},
		ReleaseDate: time.Date(2010, 12, 13, 13, 14, 15, 0, time.FixedZone("", 0)),
		Added:       time.Date(2012, 2, 2, 2, 2, 2, 0, time.FixedZone("", 0)),
		AddedBy:     User{ID: 1},
		Type:        0,
		Tags:        []string{"electronic", "canadian"},
	}

	err = db.InsertReleaseGroup(&g)
	require.Nil(t, err)

	r := Release{
		ReleaseGroup:    ReleaseGroup{ID: g.ID},
		Medium:          0,
		ReleaseDate:     time.Date(2012, 3, 2, 0, 0, 0, 0, time.FixedZone("", 0)),
		CatalogueNumber: sql.NullString{String: "NS00004"},
		RecordLabel:     RecordLabel{ID: l.ID},
		Added:           time.Date(2012, 3, 3, 0, 0, 2, 0, time.FixedZone("", 0)),
		AddedBy:         User{ID: 1},
		Original:        true,
		Tags:            []string{"special.k.edition", "some.tag"},
		Properties:      map[string]string{"LossyWebApproved": "", "LossyMasterApproved": "true"},
	}

	err = db.InsertRelease(&r)
	require.Nil(t, err)

	got, err := db.GetRelease(r.ID)
	require.Nil(t, err)
	require.False(t, got.Edition.Valid)
	require.Equal(t, r.Medium, got.Medium)
	require.Equal(t, r.ReleaseDate, got.ReleaseDate)
	require.True(t, got.CatalogueNumber.Valid)
	require.Equal(t, r.CatalogueNumber.String, got.CatalogueNumber.String)
	require.Equal(t, r.RecordLabel.ID, got.RecordLabel.ID)
	require.Equal(t, r.Added, got.Added)
	require.Equal(t, r.AddedBy.ID, got.AddedBy.ID)
	require.Equal(t, r.Original, got.Original)
	require.Equal(t, r.Tags, got.Tags)
	require.Equal(t, r.Properties, got.Properties)

	err = db.DeleteRelease(r.ID)
	require.Nil(t, err)

	_, err = db.GetRelease(r.ID)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestAutocompleteReleaseTags(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	l := RecordLabel{
		Name:    "NONESUCH",
		AddedBy: User{ID: 1},
	}

	err = db.InsertRecordLabel(&l)
	require.Nil(t, err)

	a := Artist{
		Name:    "deadmau5",
		Added:   time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	err = db.InsertArtist(&a)
	require.Nil(t, err)

	g := ReleaseGroup{
		Name: "4x4=12",
		Artists: []RoledArtist{
			{
				Role:   0,
				Artist: Artist{ID: a.ID},
			},
		},
		ReleaseDate: time.Date(2010, 12, 13, 13, 14, 15, 0, time.FixedZone("", 0)),
		Added:       time.Date(2012, 2, 2, 2, 2, 2, 0, time.FixedZone("", 0)),
		AddedBy:     User{ID: 1},
		Type:        0,
		Tags:        []string{"electronic", "canadian"},
	}

	err = db.InsertReleaseGroup(&g)
	require.Nil(t, err)

	r1 := Release{
		ReleaseGroup: ReleaseGroup{ID: g.ID},
		Medium:       0,
		ReleaseDate:  time.Date(2012, 3, 2, 0, 0, 0, 0, time.FixedZone("", 0)),
		RecordLabel:  RecordLabel{ID: l.ID},
		Added:        time.Date(2012, 3, 3, 0, 0, 2, 0, time.FixedZone("", 0)),
		AddedBy:      User{ID: 1},
		Original:     true,
		Tags:         []string{"special.k.edition", "some.tag"},
		Properties:   map[string]string{"LossyWebApproved": "", "LossyMasterApproved": "true"},
	}

	r2 := Release{
		ReleaseGroup: ReleaseGroup{ID: g.ID},
		Medium:       2,
		ReleaseDate:  time.Date(2012, 3, 2, 0, 0, 0, 0, time.FixedZone("", 0)),
		RecordLabel:  RecordLabel{ID: l.ID},
		Added:        time.Date(2012, 3, 3, 0, 0, 2, 0, time.FixedZone("", 0)),
		AddedBy:      User{ID: 1},
		Original:     true,
		Tags:         []string{"special.v.edition", "some.other.tag"},
		Properties:   map[string]string{"LossyWebApproved": "", "LossyMasterApproved": "true"},
	}

	err = db.InsertRelease(&r1)
	require.Nil(t, err)

	err = db.InsertRelease(&r2)
	require.Nil(t, err)

	tags, err := db.AutocompleteReleaseTags("k")
	require.Nil(t, err)
	require.Equal(t, 1, len(tags))
	require.Equal(t, "special.k.edition", tags[0])

	tags, err = db.AutocompleteReleaseTags("tag")
	require.Nil(t, err)
	require.Equal(t, 2, len(tags))
	require.Contains(t, tags, "some.tag")
	require.Contains(t, tags, "some.other.tag")
}

func TestReleaseSetProperty(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	l := RecordLabel{
		Name:    "NONESUCH",
		AddedBy: User{ID: 1},
	}

	err = db.InsertRecordLabel(&l)
	require.Nil(t, err)

	a := Artist{
		Name:    "deadmau5",
		Added:   time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	err = db.InsertArtist(&a)
	require.Nil(t, err)

	g := ReleaseGroup{
		Name: "4x4=12",
		Artists: []RoledArtist{
			{
				Role:   0,
				Artist: Artist{ID: a.ID},
			},
		},
		ReleaseDate: time.Date(2010, 12, 13, 13, 14, 15, 0, time.FixedZone("", 0)),
		Added:       time.Date(2012, 2, 2, 2, 2, 2, 0, time.FixedZone("", 0)),
		AddedBy:     User{ID: 1},
		Type:        0,
		Tags:        []string{"electronic", "canadian"},
	}

	err = db.InsertReleaseGroup(&g)
	require.Nil(t, err)

	r := Release{
		ReleaseGroup: ReleaseGroup{ID: g.ID},
		Medium:       0,
		ReleaseDate:  time.Date(2012, 3, 2, 0, 0, 0, 0, time.FixedZone("", 0)),
		RecordLabel:  RecordLabel{ID: l.ID},
		Added:        time.Date(2012, 3, 3, 0, 0, 2, 0, time.FixedZone("", 0)),
		AddedBy:      User{ID: 1},
		Original:     true,
		Tags:         []string{"special.k.edition", "some.tag"},
		Properties:   map[string]string{"LossyWebApproved": "", "LossyMasterApproved": "true"},
	}

	err = db.InsertRelease(&r)
	require.Nil(t, err)

	got, err := db.GetRelease(r.ID)
	require.Nil(t, err)
	require.Equal(t, r.Properties, got.Properties)

	err = db.SetReleaseProperty(r.ID, "CassetteApproved", "blah")
	require.Nil(t, err)

	err = db.SetReleaseProperty(r.ID, "LossyWebApproved", "true")
	require.Nil(t, err)

	got, err = db.GetRelease(r.ID)
	require.Nil(t, err)
	require.Equal(t, "blah", got.Properties["CassetteApproved"])
	require.Equal(t, "true", got.Properties["LossyWebApproved"])
}
