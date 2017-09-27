package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAutocompleteArtists(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a1 := Artist{
		Name: "test1",
		Bio:  sql.NullString{String: "Some bio"},
		Aliases: []ArtistAlias{{
			Alias:   "best1",
			Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
			AddedBy: User{ID: 1}}},
		Tags:    []string{"tag1", "tag2"},
		Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	a2 := Artist{
		Name: "test2",
		Aliases: []ArtistAlias{{
			Alias:   "best2",
			Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
			AddedBy: User{ID: 1}}},
		Tags:    []string{"tag1", "tag3"},
		Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 0},
	}

	err = db.InsertArtist(&a1)
	require.Nil(t, err)

	err = db.InsertArtist(&a2)
	require.Nil(t, err)

	a, err := db.AutocompleteArtists("est1")
	require.Nil(t, err)
	require.Equal(t, 1, len(a))
	require.Equal(t, a1.Name, a[0].Name)
	require.True(t, a[0].Bio.Valid)
	require.Equal(t, a1.Bio.String, a[0].Bio.String)
	require.Equal(t, 1, len(a[0].Aliases))
	require.Equal(t, a1.Aliases[0].Alias, a[0].Aliases[0].Alias)
	require.Equal(t, a1.Aliases[0].Added, a[0].Aliases[0].Added)
	require.Equal(t, a1.Aliases[0].AddedBy.ID, a[0].Aliases[0].AddedBy.ID)
	require.Equal(t, a1.Tags, a[0].Tags)
	require.Equal(t, a1.Added, a[0].Added)
	require.Equal(t, a1.AddedBy.ID, a[0].AddedBy.ID)

	a, err = db.AutocompleteArtists("test")
	require.Nil(t, err)
	require.Equal(t, 2, len(a))

	a, err = db.AutocompleteArtists("best2")
	require.Nil(t, err)
	require.Equal(t, 1, len(a))
	require.Equal(t, a2.Name, a[0].Name)
	require.False(t, a[0].Bio.Valid)
	require.Equal(t, 1, len(a[0].Aliases))
	require.Equal(t, a2.Aliases[0].Alias, a[0].Aliases[0].Alias)
	require.Equal(t, a2.Aliases[0].Added, a[0].Aliases[0].Added)
	require.Equal(t, a2.Aliases[0].AddedBy.ID, a[0].Aliases[0].AddedBy.ID)
	require.Equal(t, a2.Tags, a[0].Tags)
	require.Equal(t, a2.Added, a[0].Added)
	require.Equal(t, a2.AddedBy.ID, a[0].AddedBy.ID)
}

func TestInsertGetArtist(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a1 := Artist{
		Name: "test1",
		Bio:  sql.NullString{String: "Some bio"},
		Aliases: []ArtistAlias{{
			Alias:   "best1",
			Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
			AddedBy: User{ID: 1}}},
		Tags:    []string{"tag1", "tag2"},
		Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	err = db.InsertArtist(&a1)
	require.Nil(t, err)

	a2, err := db.GetArtist(a1.ID)
	require.Nil(t, err)
	require.NotNil(t, a2)
	require.Equal(t, a1.Name, a2.Name)
	require.True(t, a2.Bio.Valid)
	require.Equal(t, a1.Bio.String, a2.Bio.String)
	require.Equal(t, len(a1.Aliases), len(a2.Aliases))
	require.Equal(t, a1.Aliases[0].Alias, a2.Aliases[0].Alias)
	require.Equal(t, a1.Aliases[0].Added, a2.Aliases[0].Added)
	require.Equal(t, a1.Aliases[0].AddedBy.ID, a2.Aliases[0].AddedBy.ID)
	require.Equal(t, a1.Tags, a2.Tags)
	require.Equal(t, a1.Added, a2.Added)
	require.Equal(t, a1.AddedBy.ID, a2.AddedBy.ID)
}

func TestAutocompleteArtistTags(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a1 := Artist{
		Name: "test1",
		Bio:  sql.NullString{String: "Some bio"},
		Aliases: []ArtistAlias{{
			Alias:   "best1",
			Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
			AddedBy: User{ID: 1}}},
		Tags:    []string{"tag1", "tag2"},
		Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 1},
	}

	a2 := Artist{
		Name: "test2",
		Aliases: []ArtistAlias{{
			Alias:   "best2",
			Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
			AddedBy: User{ID: 1}}},
		Tags:    []string{"tag1", "tag3"},
		Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
		AddedBy: User{ID: 0},
	}

	err = db.InsertArtist(&a1)
	require.Nil(t, err)

	err = db.InsertArtist(&a2)
	require.Nil(t, err)

	tags, err := db.AutocompleteArtistTags("a")
	require.Nil(t, err)
	require.Equal(t, 3, len(tags))
	require.Contains(t, tags, "tag1")
	require.Contains(t, tags, "tag2")
	require.Contains(t, tags, "tag3")

	tags, err = db.AutocompleteArtistTags("ag3")
	require.Nil(t, err)
	require.Equal(t, 1, len(tags))
	require.Equal(t, "tag3", tags[0])
}
