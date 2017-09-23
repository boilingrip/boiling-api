package api

import (
	"errors"
	"time"

	"strings"

	"github.com/kataras/iris"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mutaborius/boiling-api/db"
)

type BlogEntry struct {
	ID       int        `json:"id"`
	Author   PublicUser `json:"author"`
	Title    string     `json:"title"`
	PostedAt time.Time  `json:"posted_at"`
	Content  string     `json:"content"`
	Tags     []string   `json:"tags"`
}

func fromBlogEntry(dbE db.BlogEntry) BlogEntry {
	return BlogEntry{
		ID:       dbE.ID,
		Author:   fromPublicUser(dbE.Author),
		Title:    dbE.Title,
		PostedAt: dbE.PostedAt,
		Content:  dbE.Content,
		Tags:     dbE.Tags,
	}
}

func toBlogEntry(e BlogEntry) db.BlogEntry {
	return db.BlogEntry{
		ID:       e.ID,
		Author:   toPublicUser(e.Author),
		Title:    e.Title,
		PostedAt: e.PostedAt,
		Content:  e.Content,
		Tags:     e.Tags,
	}
}

type BlogEntriesResponse struct {
	Entries []BlogEntry `json:"entries"`
}

func (a *API) getBlogs(ctx *context) {
	offset, err := ctx.URLParamInt("offset")
	if err != nil || offset < 0 {
		ctx.Fail(errors.New("invalid offset"), iris.StatusBadRequest)
		return
	}

	limit, err := ctx.URLParamInt("limit")
	if err != nil || limit < 0 {
		ctx.Fail(errors.New("invalid limit"), iris.StatusBadRequest)
		return
	}
	if limit > 50 {
		limit = 50
	}

	posts, err := a.db.GetBlogEntries(limit, offset)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	entries := make([]BlogEntry, 0, len(posts))
	for i := range posts {
		entries = append(entries, fromBlogEntry(posts[i]))
	}

	ctx.Success(BlogEntriesResponse{Entries: entries})
}

type BlogResponse struct {
	Entry BlogEntry `json:"entry"`
}

func (e BlogEntry) validate() error {
	if len(e.Title) == 0 {
		return errors.New("missing Title")
	}
	if len(e.Content) == 0 {
		return errors.New("missing Content")
	}
	return nil
}

func (e *BlogEntry) sanitize() {
	e.Title = bluemonday.StrictPolicy().Sanitize(e.Title)
	e.Content = bluemonday.UGCPolicy().Sanitize(e.Content)
	for i := range e.Tags {
		e.Tags[i] = strings.ToLower(bluemonday.StrictPolicy().Sanitize(e.Tags[i]))
	}
}

func (a *API) postBlog(ctx *context) {
	var entry BlogEntry
	err := ctx.ReadJSON(&entry)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}

	// post as logged-in user (TODO maybe admin can override this?)
	entry.Author = fromPublicUser(ctx.user)
	entry.PostedAt = time.Now() // TODO maybe admin can override this?

	err = entry.validate()
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
	}

	entry.sanitize()

	dbE := toBlogEntry(entry)
	err = a.db.InsertBlogEntry(&dbE)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}

	ctx.Success(BlogResponse{Entry: fromBlogEntry(dbE)})
}

func (a *API) updateBlog(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}
	if id < 0 {
		ctx.Fail(errors.New("invalid ID"), iris.StatusBadRequest)
		return
	}

	var entry BlogEntry
	err = ctx.ReadJSON(&entry)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}
	entry.ID = id

	// post as logged-in user (TODO maybe admin can override this?)
	entry.Author = fromPublicUser(ctx.user)
	entry.PostedAt = time.Now() // TODO maybe admin can override this?

	err = entry.validate()
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
	}

	entry.sanitize()

	dbE := toBlogEntry(entry)
	err = a.db.UpdateBlogEntry(dbE)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}

	ctx.Success(nil)
}

func (a *API) deleteBlog(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}
	if id < 0 {
		ctx.Fail(errors.New("invalid ID"), iris.StatusBadRequest)
		return
	}

	err = a.db.DeleteBlogEntry(id)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}

	ctx.Success(nil)
}
