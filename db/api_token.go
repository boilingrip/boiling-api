package db

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type APIToken struct {
	Token     string
	CreatedAt time.Time
	User      User
}

func generateRandomKey(length int) string {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}

	return hex.EncodeToString(buf)
}

// tokenLength defines the length of an API token.
// Note that this is the number of random bytes generated - they're then base16
// encoded, so the string representation is actually 128 characters long.
const tokenLength = 64

func (db *DB) InsertTokenForUser(u User) (*APIToken, error) {
	s := generateRandomKey(tokenLength)
	var t time.Time

	err := db.db.QueryRow("INSERT INTO api_tokens(token,created_at,uid) VALUES ($1,NOW(),$2) RETURNING created_at", s, u.ID).Scan(&t)
	if err != nil {
		return nil, err
	}

	return &APIToken{
		Token:     s,
		CreatedAt: t,
		User:      u,
	}, nil
}

func (db *DB) GetToken(token string) (*APIToken, error) {
	if len(token) == 0 {
		return nil, errors.New("invalid token")
	}

	t := APIToken{Token: token}
	res := db.db.QueryRow("SELECT t.created_at,t.uid,u.username,u.email,u.last_login,u.last_access,u.enabled,u.can_login,u.uploaded,u.downloaded FROM api_tokens t, users u WHERE t.uid = u.id AND t.token = $1", token)

	err := res.Scan(
		&t.CreatedAt,
		&t.User.ID,
		&t.User.Username,
		&t.User.Email,
		&t.User.LastLogin,
		&t.User.LastAccess,
		&t.User.Enabled,
		&t.User.CanLogin,
		&t.User.Uploaded,
		&t.User.Downloaded)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
