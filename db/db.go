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
	Close() error

	SignUpUser(username, password, email string) error
	LoginAndGetUser(username, password string) (*User, error)
	GetUser(id int) (*User, error)
	UpdateUserDeltaUpDown(id, deltaUp, deltaDown int) error
	UpdateUserSetLastAccess(id int, lastAccess time.Time) error
	UpdateUserSetLastLogin(id int, lastLogin time.Time) error
	UpdateUserAddPrivileges(id int, privileges []int) error
	PopulateUserPrivileges(u *User) error

	GetPasskeyForUser(id int) (*Passkey, error)
	GetAllPasskeysForUser(id int) ([]Passkey, error)
	GenerateNewPasskeyForUser(id int) (string, error)

	InsertTokenForUser(u User) (*APIToken, error)
	GetToken(token string) (*APIToken, error)

	InsertBlogEntry(post *BlogEntry) error
	GetBlogEntry(id int) (*BlogEntry, error)
	UpdateBlogEntry(post BlogEntry) error
	DeleteBlogEntry(id int) error
	GetBlogEntries(limit, offset int) ([]BlogEntry, error)

	AutocompleteArtists(s string) ([]Artist, error)
	AutocompleteArtistTags(s string) ([]string, error)
	GetArtist(id int) (*Artist, error)
	PopulateReleaseGroups(artist *Artist) error
	InsertArtist(artist *Artist) error

	GetAllPrivileges() (map[int]string, error)

	GetAllFormats() (map[int]Format, error)

	GetAllMedia() (map[int]string, error)

	GetAllReleaseGroupRoles() (map[int]string, error)

	GetAllLeechTypes() (map[int]string, error)

	GetAllReleaseProperties() (map[int]string, error)
	AddReleaseProperty(key string) error

	AutocompleteRecordLabels(s string) ([]RecordLabel, error)
	GetRecordLabel(id int) (*RecordLabel, error)
	InsertRecordLabel(label *RecordLabel) error

	AutocompleteReleaseTags(s string) ([]string, error)
	InsertRelease(release *Release) error
	SetReleaseProperty(id int, k, v string) error
	GetRelease(id int) (*Release, error)
	DeleteRelease(id int) error
	PopulateTorrents(release *Release) error

	AutocompleteReleaseGroups(s string) ([]ReleaseGroup, error)
	AutocompleteReleaseGroupTags(s string) ([]string, error)
	GetAllReleaseGroupTypes() (map[int]string, error)
	GetReleaseGroup(id int) (*ReleaseGroup, error)
	InsertReleaseGroup(group *ReleaseGroup) error
	PopulateReleases(group *ReleaseGroup) error
	SearchReleaseGroups(q *Query, offset, limit int) ([]ReleaseGroup, error)
}

func New(database, user, pass string) (BoilingDB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", user, pass, database))
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func compileMarkdown(input []byte) []byte {
	unsafe := blackfriday.Run(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return html
}

func (db *DB) Close() error {
	return db.db.Close()
}
