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

type ArtistsResponse struct {
	Artists []Artist `json:"artists"`
}

type TagsResponse struct {
	Tags []string `json:"tags"`
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

func (a *API) autocompleteArtist(ctx *context) {
	s := ctx.Params().Get("s")
	if len(s) == 0 {
		ctx.Fail(errors.New("missing fragment"), iris.StatusBadRequest)
		return
	}

	artists, err := a.db.AutocompleteArtists(s)
	if err != nil {
		ctx.Fail(userError(err, "not found"), iris.StatusNotFound)
		return
	}

	var toReturn []Artist
	for _, artist := range artists {
		toReturn = append(toReturn, a.artistFromDBArtist(&artist))
	}

	ctx.Success(ArtistsResponse{Artists: toReturn})
}

func (a *API) autocompleteArtistTags(ctx *context) {
	s := ctx.Params().Get("s")
	if len(s) == 0 {
		ctx.Fail(errors.New("missing fragment"), iris.StatusBadRequest)
		return
	}

	tags, err := a.db.AutocompleteArtistTags(s)
	if err != nil {
		ctx.Fail(userError(err, "not found"), iris.StatusNotFound)
		return
	}

	ctx.Success(TagsResponse{Tags: tags})
}
