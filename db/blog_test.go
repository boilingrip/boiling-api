package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsertGetDeleteBlogEntry(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	original := BlogEntry{
		ID:       0,
		Author:   User{ID: 1},
		Title:    "Test Post Do Not Read",
		PostedAt: time.Date(2001, 01, 01, 0, 0, 0, 0, time.FixedZone("", 0)),
		Content:  "This is a test post.",
		Tags:     []string{"test"},
	}
	e := original // make a copy

	// test insert
	err = db.InsertBlogEntry(&e)
	require.Nil(t, err)

	// check nothing has been changed
	require.Equal(t, original.Author.ID, e.Author.ID)
	require.Equal(t, original.Title, e.Title)
	require.Equal(t, original.PostedAt, e.PostedAt)
	require.Equal(t, original.Content, e.Content)
	require.Equal(t, original.Tags, e.Tags)

	// test get post
	p, err := db.GetBlogEntries(100, 0)
	require.Nil(t, err)
	require.Equal(t, 1, len(p))

	// check ID was set correctly
	require.Equal(t, p[0].ID, e.ID)

	// check nothing has changed
	e = p[0]
	require.Equal(t, original.Author.ID, e.Author.ID)
	require.Equal(t, original.Title, e.Title)
	require.Equal(t, original.PostedAt, e.PostedAt)
	require.Equal(t, original.Content, e.Content)
	require.Equal(t, original.Tags, e.Tags)

	// check author username
	require.NotEmpty(t, e.Author.Username)

	// test delete
	err = db.DeleteBlogEntry(e.ID)
	require.Nil(t, err)

	p, err = db.GetBlogEntries(100, 0)
	require.Nil(t, err)
	require.Equal(t, 0, len(p))
}

func TestGetBlogEntriesOrdering(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	original := BlogEntry{
		ID:       0,
		Author:   User{ID: 1},
		Title:    "Test Post Do Not Read",
		PostedAt: time.Date(2001, 01, 01, 0, 0, 0, 0, time.FixedZone("", 0)),
		Content:  "This is a test post.",
		Tags:     []string{"test"},
	}

	err = db.InsertBlogEntry(&original)
	require.Nil(t, err)

	original.PostedAt = time.Date(2001, 01, 02, 0, 0, 0, 0, time.FixedZone("", 0))
	original.Tags = []string{"test", "second"}

	err = db.InsertBlogEntry(&original)
	require.Nil(t, err)

	p, err := db.GetBlogEntries(100, 0)
	require.Nil(t, err)
	require.Equal(t, 2, len(p))

	require.Equal(t, original.Tags, p[0].Tags)
	require.Equal(t, []string{"test"}, p[1].Tags)

	require.True(t, p[0].PostedAt.After(p[1].PostedAt))
}

func TestUpdateBlogEntry(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	post := BlogEntry{
		ID:       0,
		Author:   User{ID: 1},
		Title:    "Test Post Do Not Read",
		PostedAt: time.Date(2001, 01, 01, 0, 0, 0, 0, time.FixedZone("", 0)),
		Content:  "This is a test post.",
		Tags:     []string{"test"},
	}

	err = db.InsertBlogEntry(&post)
	require.Nil(t, err)

	post.Author = User{ID: 0}
	post.Title = "Some new title"
	post.PostedAt = time.Date(2001, 01, 02, 0, 0, 0, 0, time.FixedZone("", 0))
	post.Content = "Updated content"
	post.Tags = []string{"something", "else"}

	err = db.UpdateBlogEntry(post)
	require.Nil(t, err)

	p, err := db.GetBlogEntries(100, 0)
	require.Nil(t, err)
	require.Equal(t, 1, len(p))
	require.Equal(t, post.ID, p[0].ID)
	require.Equal(t, post.Author.ID, p[0].Author.ID)
	require.Equal(t, post.Title, p[0].Title)
	require.Equal(t, post.PostedAt, p[0].PostedAt)
	require.Equal(t, post.Content, p[0].Content)
	require.Equal(t, post.Tags, p[0].Tags)
}
