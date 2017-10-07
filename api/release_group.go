package api

import (
	"errors"
	"time"

	"github.com/kataras/iris"

	"github.com/boilingrip/boiling-api/db"
)

type BaseReleaseGroup struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ReleaseDate time.Time `json:"release_date"`
}

func (a *API) baseReleaseGroupFromDBReleaseGroup(dbRG *db.ReleaseGroup) BaseReleaseGroup {
	return BaseReleaseGroup{
		ID:          dbRG.ID,
		Name:        dbRG.Name,
		Type:        a.c.releaseGroupTypes.MustReverseLookUp(dbRG.Type),
		ReleaseDate: dbRG.ReleaseDate,
	}
}

type ReleaseGroup struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"` // Album title
	Artists     []RoledArtist `json:"artists"`
	ReleaseDate time.Time     `json:"release_date"`
	Added       time.Time     `json:"added"`
	AddedBy     User          `json:"added_by"`
	Type        string        `json:"type"`               // Album/EP
	Releases    []Release     `json:"releases,omitempty"` // individual releases
	Tags        []string      `json:"tags,omitempty"`
}

func (a *API) releaseGroupFromDBReleaseGroup(dbRG *db.ReleaseGroup) ReleaseGroup {
	rg := ReleaseGroup{
		ID:          dbRG.ID,
		Name:        dbRG.Name,
		ReleaseDate: dbRG.ReleaseDate,
		Added:       dbRG.Added,
		AddedBy:     userFromDBUser(dbRG.AddedBy),
		Tags:        dbRG.Tags,
		Type:        a.c.releaseGroupTypes.MustReverseLookUp(dbRG.Type),
	}

	for _, dba := range dbRG.Artists {
		artist := a.roledArtistFromDBRoledArtist(&dba)
		rg.Artists = append(rg.Artists, artist)
	}

	for _, dbr := range dbRG.Releases {
		r := a.releaseFromDBRelease(&dbr)
		rg.Releases = append(rg.Releases, r)
	}

	return rg
}

type RoledArtist struct {
	Role   string     `json:"role"`
	Artist BaseArtist `json:"artist"`
}

func (a *API) roledArtistFromDBRoledArtist(dbra *db.RoledArtist) RoledArtist {
	return RoledArtist{
		Role:   a.c.releaseRoles.MustReverseLookUp(dbra.Role),
		Artist: baseArtistFromDBArtist(&dbra.Artist),
	}
}

type ReleaseGroupResponse struct {
	ReleaseGroup ReleaseGroup `json:"release_group"`
}

func (a *API) getReleaseGroup(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(userError(err, "invalid ID"), iris.StatusBadRequest)
		return
	}
	if id < 0 {
		ctx.Fail(errors.New("invalid ID"), iris.StatusBadRequest)
		return
	}

	group, err := a.db.GetReleaseGroup(id)
	if err != nil {
		ctx.Fail(userError(err, "not found"), iris.StatusNotFound)
		return
	}

	err = a.db.PopulateReleases(group)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	ctx.Success(ReleaseGroupResponse{ReleaseGroup: a.releaseGroupFromDBReleaseGroup(group)})
}
