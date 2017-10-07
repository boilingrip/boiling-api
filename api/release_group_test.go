package api

import (
	"database/sql"
	"testing"
	"time"

	"github.com/boilingrip/boiling-api/db"
	"github.com/stretchr/testify/require"
	"gopkg.in/gavv/httpexpect.v1"
)

func TestGetReleaseGroup(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	a, err := getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)
	err = givePrivileges(a, tc.user.ID, "get_release_group")
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

	g := db.ReleaseGroup{
		Name: "4x4=12",
		Artists: []db.RoledArtist{
			{
				Role:   0,
				Artist: db.Artist{ID: a1.ID},
			},
		},
		ReleaseDate: time.Date(2010, 12, 13, 13, 14, 15, 0, time.FixedZone("", 0)),
		Added:       time.Date(2012, 2, 2, 2, 2, 2, 0, time.FixedZone("", 0)),
		AddedBy:     db.User{ID: 1},
		Type:        0,
		Tags:        []string{"electronic", "canadian"},
	}

	err = tc.db.InsertReleaseGroup(&g)
	require.Nil(t, err)

	l := db.RecordLabel{
		Name:    "NONESUCH",
		AddedBy: db.User{ID: 1},
	}

	err = tc.db.InsertRecordLabel(&l)
	require.Nil(t, err)

	r := db.Release{
		ReleaseGroup:    db.ReleaseGroup{ID: g.ID},
		Medium:          0,
		ReleaseDate:     time.Date(2012, 3, 2, 0, 0, 0, 0, time.FixedZone("", 0)),
		CatalogueNumber: sql.NullString{String: "NS00004"},
		RecordLabel:     db.RecordLabel{ID: l.ID},
		Added:           time.Date(2012, 3, 3, 0, 0, 2, 0, time.FixedZone("", 0)),
		AddedBy:         db.User{ID: 1},
		Original:        true,
		Tags:            []string{"special.k.edition", "some.tag"},
		Properties:      map[string]string{"LossyWebApproved": "", "LossyMasterApproved": "true"},
	}

	err = tc.db.InsertRelease(&r)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.GET("/release_groups/{id}", a1.ID).
		WithHeader("X-User-Token", tc.token).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Keys().ContainsOnly("release_group")
	group := obj.Value("data").Object().Value("release_group").Object()

	group.Keys().ContainsOnly("id", "name", "artists", "release_date", "added", "added_by", "type", "tags", "releases")
	group.ValueEqual("id", g.ID)
	group.ValueEqual("name", g.Name)
	group.ValueEqual("release_date", g.ReleaseDate)
	group.ValueEqual("added", g.Added)
	group.Value("added_by").Object().ValueEqual("id", g.AddedBy.ID)
	group.ValueEqual("type", "Album")
	group.ValueEqual("tags", g.Tags)
	group.Value("artists").Array().Length().Equal(1)
	group.Value("releases").Array().Length().Equal(1)

	artist := group.Value("artists").Array().Element(0).Object()
	artist.Keys().ContainsOnly("role", "artist")
	artist.ValueEqual("role", "Main")
	artist.Value("artist").Object().ValueEqual("id", a1.ID)
	artist.Value("artist").Object().ValueEqual("name", a1.Name)

	release := group.Value("releases").Array().Element(0).Object()
	//release.Keys().ContainsOnly("id")
	release.ValueEqual("id", r.ID)

}
