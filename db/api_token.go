package db

import (
	"errors"
	"math/rand"
	"time"
)

type APIToken struct {
	Token     string
	CreatedAt time.Time
	User      User
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomAlphanumeric(length int) string {
	s := make([]byte, length)
	for i := 0; i < len(s); i++ {
		s[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(s)
}

const tokenLength = 128

func (db *DB) InsertTokenForUser(u User) (*APIToken, error) {
	s := generateRandomAlphanumeric(tokenLength)
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
