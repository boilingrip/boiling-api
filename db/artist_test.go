package db

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutocompleteArtist(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a1 := Artist{
		Name:    "test1",
		Bio:     sql.NullString{String: "Some bio"},
		Aliases: []string{"best1"},
		Tags:    []string{"tag1", "tag2"},
	}

	a2 := Artist{
		Name:    "test2",
		Aliases: []string{"best2"},
		Tags:    []string{"tag1", "tag3"},
	}

	err = db.InsertArtist(&a1)
	require.Nil(t, err)

	err = db.InsertArtist(&a2)
	require.Nil(t, err)

	a, err := db.AutocompleteArtist("est1")
	require.Nil(t, err)
	require.Equal(t, 1, len(a))
	require.Equal(t, a1.Name, a[0].Name)
	require.True(t, a[0].Bio.Valid)
	require.Equal(t, a1.Bio.String, a[0].Bio.String)
	require.Equal(t, a1.Aliases, a[0].Aliases)
	require.Equal(t, a1.Tags, a[0].Tags)

	a, err = db.AutocompleteArtist("test")
	require.Nil(t, err)
	require.Equal(t, 2, len(a))

	a, err = db.AutocompleteArtist("best2")
	require.Nil(t, err)
	require.Equal(t, 1, len(a))
	require.Equal(t, a2.Name, a[0].Name)
	require.False(t, a[0].Bio.Valid)
	require.Equal(t, a2.Aliases, a[0].Aliases)
	require.Equal(t, a2.Tags, a[0].Tags)
}

func TestInsertGetArtist(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a1 := Artist{
		Name:    "test1",
		Bio:     sql.NullString{String: "Some bio"},
		Aliases: []string{"best1"},
		Tags:    []string{"tag1", "tag2"},
	}

	err = db.InsertArtist(&a1)
	require.Nil(t, err)

	a2, err := db.GetArtist(a1.ID)
	require.Nil(t, err)
	require.NotNil(t, a2)
	require.Equal(t, a1.Name, a2.Name)
	require.True(t, a2.Bio.Valid)
	require.Equal(t, a1.Bio.String, a2.Bio.String)
	require.Equal(t, a1.Aliases, a2.Aliases)
	require.Equal(t, a1.Tags, a2.Tags)
}

func TestAutocompleteArtistTags(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	a1 := Artist{
		Name:    "test1",
		Bio:     sql.NullString{String: "Some bio"},
		Aliases: []string{"best1"},
		Tags:    []string{"tag1", "tag2"},
	}

	a2 := Artist{
		Name:    "test2",
		Aliases: []string{"best2"},
		Tags:    []string{"tag1", "tag3"},
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
