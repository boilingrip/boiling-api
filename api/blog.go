package api

import (
	"errors"
	"time"

	"github.com/kataras/iris"

	"github.com/boilingrip/boiling-api/db"
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
		ctx.Fail(userError(err, "invalid offset"), iris.StatusBadRequest)
		return
	}

	limit, err := ctx.URLParamInt("limit")
	if err != nil || limit < 0 {
		ctx.Fail(userError(err, "invalid limit"), iris.StatusBadRequest)
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

func (a *API) postBlog(ctx *context) {
	title := ctx.fields.mustGetString("title")
	content := ctx.fields.mustGetString("content")
	tags := ctx.fields.mustGetTags("tags")
	author, authorSet := ctx.fields.getInt("author")
	postedAt, postedAtSet := ctx.fields.getDate("posted_at")

	entry := db.BlogEntry{
		Title:    title,
		Content:  content,
		Tags:     tags,
		Author:   ctx.user,
		PostedAt: time.Now(),
	}

	if authorSet {
		entry.Author = db.User{ID: author}
	}
	if postedAtSet {
		entry.PostedAt = postedAt
	}

	err := a.db.InsertBlogEntry(&entry)
	if err != nil {
		ctx.Fail(userError(err, "unable to post blog"), iris.StatusBadRequest)
		return
	}

	ctx.Success(BlogResponse{Entry: fromBlogEntry(entry)})
}

func (a *API) updateBlog(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(userError(err, "invalid ID"), iris.StatusBadRequest)
		return
	}
	if id < 0 {
		ctx.Fail(errors.New("invalid ID"), iris.StatusBadRequest)
		return
	}

	title := ctx.fields.mustGetString("title")
	content := ctx.fields.mustGetString("content")
	tags := ctx.fields.mustGetTags("tags")
	author, authorSet := ctx.fields.getInt("author")
	postedAt, postedAtSet := ctx.fields.getDate("posted_at")

	canEditForeignPost, err := a.containsPrivilege(ctx.user.Privileges, "update_blog_not_owner")
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	original, err := a.db.GetBlogEntry(id)
	if err != nil {
		ctx.Fail(userError(err, "not found"), iris.StatusBadRequest)
		return
	}

	if !canEditForeignPost && original.Author.ID != ctx.user.ID {
		ctx.Fail(errors.New("can only edit own posts (missing update_blog_not_owner privilege)"), iris.StatusUnauthorized)
		return
	}

	// can update title, content and tags
	original.Title = title
	original.Content = content
	original.Tags = tags
	if authorSet {
		original.Author = db.User{ID: author}
	}
	if postedAtSet {
		original.PostedAt = postedAt
	}

	err = a.db.UpdateBlogEntry(*original)
	if err != nil {
		ctx.Fail(userError(err, "unable to update blog"), iris.StatusBadRequest)
		return
	}

	ctx.Success(BlogResponse{Entry: fromBlogEntry(*original)})
}

func (a *API) deleteBlog(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(userError(err, "invalid ID"), iris.StatusBadRequest)
		return
	}
	if id < 0 {
		ctx.Fail(errors.New("invalid ID"), iris.StatusBadRequest)
		return
	}

	canDeleteForeignPost, err := a.containsPrivilege(ctx.user.Privileges, "delete_blog_not_owner")
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	if !canDeleteForeignPost {
		original, err := a.db.GetBlogEntry(id)
		if err != nil {
			ctx.Fail(userError(err, "not found"), iris.StatusBadRequest)
			return
		}

		if original.Author.ID != ctx.user.ID {
			ctx.Fail(errors.New("can only delete own posts (missing delete_blog_not_owner privilege)"), iris.StatusUnauthorized)
			return
		}
	}

	err = a.db.DeleteBlogEntry(id)
	if err != nil {
		ctx.Fail(userError(err, "unable to delete blog"), iris.StatusBadRequest)
		return
	}

	ctx.Success(nil)
}
