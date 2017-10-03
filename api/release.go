package api

import (
	"time"

	"github.com/boilingrip/boiling-api/db"
)

type Release struct {
	ID              int          `json:"id"`
	ReleaseGroup    ReleaseGroup `json:"release_group"`
	Edition         *string      `json:"edition,omitempty"`
	Medium          string       `json:"medium"`
	ReleaseDate     time.Time    `json:"release_date"`
	CatalogueNumber *string      `json:"catalogue_number,omitempty"`
	//RecordLabel     RecordLabel
	//Torrents        []Torrent
	Added    time.Time `json:"added"`
	AddedBy  BaseUser  `json:"added_by"`
	Original bool      `json:"original"`
	Tags     []string  `json:"tags,omitempty"`

	// Properties lists "official" properties of releases, for example
	// "LossyMasterApproved", to be set by trusted users or staff.
	Properties map[string]string `json:"properties"`
}

func (a *API) releaseFromDBRelease(dbR *db.Release) Release {
	r := Release{
		ID:          dbR.ID,
		Medium:      a.c.media.MustReverseLookUp(dbR.Medium),
		ReleaseDate: dbR.ReleaseDate,
		Added:       dbR.Added,
		AddedBy:     baseUserFromDBUser(dbR.AddedBy),
		Original:    dbR.Original,
		Tags:        dbR.Tags,
		Properties:  dbR.Properties,
	}
	if dbR.Edition.Valid {
		r.Edition = &dbR.Edition.String
	}
	if dbR.CatalogueNumber.Valid {
		r.CatalogueNumber = &dbR.CatalogueNumber.String
	}

	return r
}

type ReleaseResponse struct {
	Release Release `json:"release"`
}
