package api

import (
	"testing"
	"time"

	"github.com/mutaborius/boiling-api/db"
	"github.com/stretchr/testify/require"
	"gopkg.in/gavv/httpexpect.v1"
)

func TestGetBlogs(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	entry := db.BlogEntry{
		Title:    "test title",
		Content:  "test content",
		Tags:     []string{"some", "tags"},
		Author:   db.User{ID: 1},
		PostedAt: time.Date(2001, 01, 01, 0, 0, 0, 0, time.FixedZone("", 0)),
	}

	err = tc.db.InsertBlogEntry(&entry)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.GET("/blogs").
		WithHeader("X-User-Token", tc.token).
		WithQuery("offset", 0).
		WithQuery("limit", 100).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Value("entries").Array().Length().Equal(1)
	post := obj.Value("data").Object().Value("entries").Array().Element(0).Object()

	post.Keys().ContainsOnly("id", "title", "content", "tags", "author", "posted_at")
	post.ValueEqual("id", entry.ID)
	post.ValueEqual("title", entry.Title)
	post.ValueEqual("content", entry.Content)
	post.ValueEqual("tags", entry.Tags)
	post.ValueEqual("posted_at", entry.PostedAt)
	post.Value("author").Object().ValueEqual("id", entry.Author.ID)
}

func TestInsertUpdateDeleteBlog(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	entry := BlogEntry{
		Title:   "Test Entry",
		Content: "some content",
		Tags:    []string{"some", "tags"},
	}

	e := httpexpect.New(t, "http://localhost:8080")

	// Post
	resp := e.POST("/blogs").
		WithHeader("X-User-Token", tc.token).
		WithJSON(entry).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Keys().ContainsOnly("entry")
	got := obj.Value("data").Object().Value("entry").Object()
	got.ValueEqual("title", entry.Title)
	got.ValueEqual("content", entry.Content)
	got.ValueEqual("tags", entry.Tags)
	got.Value("author").Object().ValueEqual("id", tc.user.ID)
	id := int(got.Value("id").Number().Raw())
	require.True(t, id > 0)

	// Update
	entry.Title = "some new title"
	entry.Content = "a new content"
	entry.Tags = []string{"different"}

	resp = e.POST("/blogs/{id}", id).
		WithHeader("X-User-Token", tc.token).
		WithJSON(entry).
		Expect().Status(200)

	obj = resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	post := obj.Value("data").Object().Value("entry").Object()

	post.Keys().ContainsOnly("id", "title", "content", "tags", "author", "posted_at")
	post.ValueEqual("id", id)
	post.ValueEqual("title", entry.Title)
	post.ValueEqual("content", entry.Content)
	post.ValueEqual("tags", entry.Tags)

	// Check update
	resp = e.GET("/blogs").
		WithHeader("X-User-Token", tc.token).
		WithQuery("offset", 0).
		WithQuery("limit", 100).
		Expect().Status(200)

	obj = resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Value("entries").Array().Length().Equal(1)
	post = obj.Value("data").Object().Value("entries").Array().Element(0).Object()

	post.Keys().ContainsOnly("id", "title", "content", "tags", "author", "posted_at")
	post.ValueEqual("id", id)
	post.ValueEqual("title", entry.Title)
	post.ValueEqual("content", entry.Content)
	post.ValueEqual("tags", entry.Tags)

	// Delete
	resp = e.DELETE("/blogs/{id}", id).
		WithHeader("X-User-Token", tc.token).
		Expect().Status(200)

	obj = resp.JSON().Object()
	obj.Keys().ContainsOnly("status")
	obj.ValueEqual("status", "success")

	// Check delete
	resp = e.GET("/blogs").
		WithHeader("X-User-Token", tc.token).
		WithQuery("offset", 0).
		WithQuery("limit", 100).
		Expect().Status(200)

	obj = resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Value("entries").Array().Length().Equal(0)
}
