package db

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type BlogEntry struct {
	ID       int
	Author   User
	Title    string
	PostedAt time.Time
	Content  string
	Tags     []string
}

func insertBlogTagsTx(post BlogEntry, tx *sql.Tx) error {
	var res sql.Result
	var err error
	for _, t := range post.Tags {
		var id int
		err = tx.QueryRow("INSERT INTO blog_tags(tag) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", t).Scan(&id)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == nil {
			// inserted
			res, err = tx.Exec("INSERT INTO blog_tags_blogs(blog,tag) VALUES($1,$2)", post.ID, id)
		} else {
			// already present
			res, err = tx.Exec("INSERT INTO blog_tags_blogs(blog,tag) VALUES($1,(SELECT id FROM blog_tags WHERE tag=$2 LIMIT 1))", post.ID, t)
		}
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return errors.New("did not insert")
		}
	}

	return nil
}

func addBlogPostTx(post *BlogEntry, tx *sql.Tx) error {
	err := tx.QueryRow("INSERT INTO blogs(author,title,content,posted_at) VALUES($1,$2,$3,$4) RETURNING id", post.Author.ID, post.Title, post.Content, post.PostedAt).Scan(&post.ID)
	if err != nil {
		return err
	}

	return insertBlogTagsTx(*post, tx)
}

func (db *DB) InsertBlogEntry(post *BlogEntry) error {
	if post == nil {
		return errors.New("no post provided")
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = addBlogPostTx(post, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func deletePostTagsTx(id int, tx *sql.Tx) error {
	res, err := tx.Exec("DELETE FROM blog_tags_blogs WHERE blog = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 0 {
		_, err = tx.Exec("DELETE FROM blog_tags t WHERE NOT EXISTS (SELECT * FROM blog_tags_blogs b WHERE b.tag = t.id);")
		return err
	}

	return nil
}

func deleteBlogPostTx(id int, tx *sql.Tx) error {
	err := deletePostTagsTx(id, tx)
	if err != nil {
		return err
	}

	res, err := tx.Exec("DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("did not delete")
	}

	return nil
}

func (db *DB) DeleteBlogEntry(id int) error {
	if id < 0 {
		return errors.New("invalid ID")
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = deleteBlogPostTx(id, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func updateBlogPostTx(post BlogEntry, tx *sql.Tx) error {
	err := deletePostTagsTx(post.ID, tx)
	if err != nil {
		return err
	}

	err = insertBlogTagsTx(post, tx)
	if err != nil {
		return err
	}

	res, err := tx.Exec("UPDATE blogs SET author=$1,title=$2,content=$3,posted_at=$4 WHERE id=$5", post.Author.ID, post.Title, post.Content, post.PostedAt, post.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("entry not found")
	}

	return nil
}

func (db *DB) UpdateBlogEntry(post BlogEntry) error {
	if post.ID < 0 {
		return errors.New("invalid ID")
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = updateBlogPostTx(post, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) populateBlogTags(entry *BlogEntry) error {
	inner, err := db.db.Query("SELECT t.tag FROM blog_tags t, blog_tags_blogs b WHERE b.blog = $1 AND b.tag = t.id", entry.ID)
	if err != nil {
		return err
	}
	defer inner.Close()

	for inner.Next() {
		var tag string
		err = inner.Scan(&tag)
		if err != nil {
			return err
		}
		entry.Tags = append(entry.Tags, tag)
	}

	return nil
}

func (db *DB) GetBlogEntry(id int) (*BlogEntry, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	row := db.db.QueryRow("SELECT b.id,b.author,u.username,b.title,b.content,b.posted_at FROM blogs b,users u WHERE b.author = u.id AND b.id = $1", id)
	var entry BlogEntry
	err := row.Scan(
		&entry.ID,
		&entry.Author.ID,
		&entry.Author.Username,
		&entry.Title,
		&entry.Content,
		&entry.PostedAt)

	if err != nil {
		return nil, err
	}

	err = db.populateBlogTags(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (db *DB) GetBlogEntries(limit, offset int) ([]BlogEntry, error) {
	if limit < 1 {
		return nil, errors.New("invalid limit")
	}
	if offset < 0 {
		return nil, errors.New("invalid offset")
	}

	rows, err := db.db.Query("SELECT b.id,b.author,u.username,b.title,b.content,b.posted_at FROM blogs b,users u WHERE b.author = u.id ORDER BY b.posted_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]BlogEntry, 0)

	for rows.Next() {
		var entry BlogEntry
		err = rows.Scan(
			&entry.ID,
			&entry.Author.ID,
			&entry.Author.Username,
			&entry.Title,
			&entry.Content,
			&entry.PostedAt)
		if err != nil {
			return nil, err
		}

		err = db.populateBlogTags(&entry)
		if err != nil {
			return nil, err
		}

		result = append(result, entry)
	}

	return result, nil
}
