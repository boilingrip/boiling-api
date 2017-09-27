package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAutocompleteReleaseGroupTags(t *testing.T) {
	db, err := cleanDB()
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

	tags, err := db.AutocompleteReleaseGroupTags("elec")
	require.Nil(t, err)
	require.Equal(t, 1, len(tags))
	require.Equal(t, "electronic", tags[0])

	tags, err = db.AutocompleteReleaseGroupTags("i")
	require.Nil(t, err)
	require.Equal(t, 2, len(tags))
	require.Contains(t, tags, "electronic")
	require.Contains(t, tags, "canadian")
}

func TestAutocompleteReleaseGroups(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a := Artist{
		Name:    "deadmau5",
		Added:   time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	err = db.InsertArtist(&a)
	require.Nil(t, err)

	g1 := ReleaseGroup{
		Name: "Some old title",
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

	g2 := ReleaseGroup{
		Name: "Some new title",
		Artists: []RoledArtist{
			{
				Role:   0,
				Artist: Artist{ID: a.ID},
			},
		},
		ReleaseDate: time.Date(2011, 12, 13, 13, 14, 15, 0, time.FixedZone("", 0)),
		Added:       time.Date(2013, 2, 2, 2, 2, 2, 0, time.FixedZone("", 0)),
		AddedBy:     User{ID: 1},
		Type:        0,
		Tags:        []string{"electronic", "canadian", "house"},
	}

	err = db.InsertReleaseGroup(&g1)
	require.Nil(t, err)

	err = db.InsertReleaseGroup(&g2)
	require.Nil(t, err)

	groups, err := db.AutocompleteReleaseGroups("old")
	require.Nil(t, err)
	require.Equal(t, 1, len(groups))
	require.Equal(t, g1.Name, groups[0].Name)
	require.Equal(t, 1, len(groups[0].Artists))
	require.Equal(t, g1.Artists[0].Role, groups[0].Artists[0].Role)
	require.Equal(t, g1.Artists[0].Artist.ID, groups[0].Artists[0].Artist.ID)
	require.Equal(t, g1.ReleaseDate, groups[0].ReleaseDate)
	require.Equal(t, g1.Added, groups[0].Added)
	require.Equal(t, g1.AddedBy.ID, groups[0].AddedBy.ID)
	require.Equal(t, g1.Type, groups[0].Type)
	require.Equal(t, g1.Tags, groups[0].Tags)

	groups, err = db.AutocompleteReleaseGroups("new")
	require.Nil(t, err)
	require.Equal(t, 1, len(groups))
	require.Equal(t, g2.Name, groups[0].Name)
	require.Equal(t, 1, len(groups[0].Artists))
	require.Equal(t, g2.Artists[0].Role, groups[0].Artists[0].Role)
	require.Equal(t, g2.Artists[0].Artist.ID, groups[0].Artists[0].Artist.ID)
	require.Equal(t, g2.ReleaseDate, groups[0].ReleaseDate)
	require.Equal(t, g2.Added, groups[0].Added)
	require.Equal(t, g2.AddedBy.ID, groups[0].AddedBy.ID)
	require.Equal(t, g2.Type, groups[0].Type)
	require.Equal(t, g2.Tags, groups[0].Tags)

	groups, err = db.AutocompleteReleaseGroups("title")
	require.Nil(t, err)
	require.Equal(t, 2, len(groups))
}

func TestInsertGetReleaseGroup(t *testing.T) {
	db, err := cleanDB()
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

	got, err := db.GetReleaseGroup(g.ID)
	require.Nil(t, err)
	require.Equal(t, g.Name, got.Name)
	require.Equal(t, len(g.Artists), len(got.Artists))
	require.Equal(t, g.Artists[0].Role, got.Artists[0].Role)
	require.Equal(t, g.Artists[0].Artist.ID, got.Artists[0].Artist.ID)
	require.Equal(t, g.ReleaseDate, got.ReleaseDate)
	require.Equal(t, g.Added, got.Added)
	require.Equal(t, g.AddedBy.ID, got.AddedBy.ID)
	require.Equal(t, g.Type, got.Type)
	require.Equal(t, g.Tags, got.Tags)
}
