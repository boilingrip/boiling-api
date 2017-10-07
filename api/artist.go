package api

import (
	"errors"
	"time"

	"github.com/kataras/iris"

	"github.com/boilingrip/boiling-api/db"
)

type Artist struct {
	ID            int                 `json:"id"`
	Name          string              `json:"name"`
	Aliases       []ArtistAlias       `json:"aliases,omitempty"`
	ReleaseGroups []RoledReleaseGroup `json:"release_groups,omitempty"`
	Added         time.Time           `json:"added"`
	AddedBy       BaseUser            `json:"added_by"`
	Bio           *string             `json:"bio,omitempty"`
	Tags          []string            `json:"tags,omitempty"`
}

func (a *API) artistFromDBArtist(dbA *db.Artist) Artist {
	artist := Artist{
		ID:      dbA.ID,
		Name:    dbA.Name,
		Added:   dbA.Added,
		AddedBy: baseUserFromDBUser(dbA.AddedBy),
		Tags:    dbA.Tags,
	}
	if dbA.Bio.Valid {
		artist.Bio = &dbA.Bio.String
	}
	for _, dbAlias := range dbA.Aliases {
		alias := artistAliasFromDBArtistAlias(dbAlias)
		artist.Aliases = append(artist.Aliases, alias)
	}
	if len(dbA.ReleaseGroups) > 0 {
		for _, dbg := range dbA.ReleaseGroups {
			rg := a.roledReleaseGroupFromDBRoledReleaseGroup(&dbg)
			artist.ReleaseGroups = append(artist.ReleaseGroups, rg)
		}
	}

	return artist
}

type RoledReleaseGroup struct {
	Role         string           `json:"role"`
	ReleaseGroup BaseReleaseGroup `json:"release_group"`
}

func (a *API) roledReleaseGroupFromDBRoledReleaseGroup(dbrg *db.RoledReleaseGroup) RoledReleaseGroup {
	return RoledReleaseGroup{
		Role:         a.c.releaseRoles.MustReverseLookUp(dbrg.Role),
		ReleaseGroup: a.baseReleaseGroupFromDBReleaseGroup(&dbrg.ReleaseGroup),
	}
}

type ArtistAlias struct {
	Alias   string    `json:"alias"`
	Added   time.Time `json:"added"`
	AddedBy BaseUser  `json:"added_by"`
}

func artistAliasFromDBArtistAlias(dbA db.ArtistAlias) ArtistAlias {
	return ArtistAlias{
		Alias:   dbA.Alias,
		Added:   dbA.Added,
		AddedBy: baseUserFromDBUser(dbA.AddedBy),
	}
}

type BaseArtist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func baseArtistFromDBArtist(dbA *db.Artist) BaseArtist {
	return BaseArtist{
		ID:   dbA.ID,
		Name: dbA.Name,
	}
}

type ArtistResponse struct {
	Artist Artist `json:"artist"`
}

func (a *API) getArtist(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(userError(err, "invalid ID"), iris.StatusBadRequest)
		return
	}
	if id < 0 {
		ctx.Fail(errors.New("invalid ID"), iris.StatusBadRequest)
		return
	}

	artist, err := a.db.GetArtist(id)
	if err != nil {
		ctx.Fail(userError(err, "not found"), iris.StatusNotFound)
		return
	}

	err = a.db.PopulateReleaseGroups(artist)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	ctx.Success(ArtistResponse{Artist: a.artistFromDBArtist(artist)})
}
