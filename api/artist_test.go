package api

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/boilingrip/boiling-api/db"
)

func TestGetArtist(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	a, err := getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)
	err = givePrivileges(a, tc.user.ID, "get_artist")
	require.Nil(t, err)

	a1 := db.Artist{
		Name: "test1",
		Bio:  sql.NullString{String: "Some bio"},
		Aliases: []db.ArtistAlias{{
			Alias:   "best1",
			Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
			AddedBy: db.User{ID: 1}}},
		Tags:    []string{"tag1", "tag2"},
		Added:   time.Date(2010, 03, 02, 12, 34, 0, 0, time.FixedZone("", 0)),
		AddedBy: db.User{ID: 1},
	}

	err = tc.db.InsertArtist(&a1)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.GET("/artists/{id}", a1.ID).
		WithHeader("X-User-Token", tc.token).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Keys().ContainsOnly("artist")
	artist := obj.Value("data").Object().Value("artist").Object()

	artist.Keys().ContainsOnly("id", "name", "bio", "aliases", "tags", "added", "added_by")
	artist.ValueEqual("id", a1.ID)
	artist.ValueEqual("name", a1.Name)
	artist.Value("aliases").Array().Length().Equal(1)
	artist.ValueEqual("bio", a1.Bio.String)
	artist.ValueEqual("tags", a1.Tags)
	artist.ValueEqual("added", a1.Added)
	artist.Value("added_by").Object().ValueEqual("id", a1.AddedBy.ID)

	alias := artist.Value("aliases").Array().Element(0).Object()
	alias.Keys().ContainsOnly("alias", "added", "added_by")
	alias.ValueEqual("alias", a1.Aliases[0].Alias)
	alias.ValueEqual("added", a1.Aliases[0].Added)
	alias.Value("added_by").Object().ValueEqual("id", a1.Aliases[0].AddedBy.ID)
}
