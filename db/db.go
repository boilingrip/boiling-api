package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type DB struct {
	db *sql.DB
}

type BoilingDB interface {
	SignUpUser(username, password, email string) error
	LoginAndGetUser(username, password string) (*User, error)
	GetUser(id int) (*User, error)
	UpdateUserDeltaUpDown(id, deltaUp, deltaDown int) error
	UpdateUserSetLastAccess(id int, lastAccess time.Time) error
	UpdateUserSetLastLogin(id int, lastLogin time.Time) error

	InsertTokenForUser(u User) (*APIToken, error)
	GetToken(token string) (*APIToken, error)

	InsertBlogEntry(post *BlogEntry) error
	UpdateBlogEntry(post BlogEntry) error
	DeleteBlogEntry(id int) error
	GetBlogEntries(limit, offset int) ([]BlogEntry, error)

	AutocompleteArtist(s string) ([]Artist, error)
	AutocompleteArtistTags(s string) ([]string, error)
	GetArtist(id int) (*Artist, error)
	PopulateReleaseGroups(artist *Artist) error
	InsertArtist(artist *Artist) error

	GetAllMedia() (map[int]string, error)

	AutocompleteRecordLabel(s string) ([]RecordLabel, error)
	GetRecordLabel(id int) (*RecordLabel, error)
	InsertRecordLabel(label *RecordLabel) error

	AutocompleteReleaseTags(s string) ([]string, error)
	InsertRelease(release *Release) error
	GetRelease(id int) (*Release, error)
	DeleteRelease(id int) error
	PopulateFormats(release *Release) error

	AutocompleteReleaseGroup(s string) ([]ReleaseGroup, error)
	AutocompleteReleaseGroupTags(s string) ([]string, error)
	GetAllReleaseGroupTypes() ([]string, error)
	GetReleaseGroup(id int) (*ReleaseGroup, error)
	InsertReleaseGroup(group *ReleaseGroup) error
	PopulateReleases(group *ReleaseGroup) error
}

func New(database, user, pass string) (BoilingDB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", user, pass, database))
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func compileMarkdown(input []byte) []byte {
	unsafe := blackfriday.MarkdownCommon(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return html
}
