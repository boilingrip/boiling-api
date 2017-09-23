package api

import (
	"time"

	"github.com/mutaborius/boiling-api/db"
)

type PublicUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func fromPublicUser(dbU db.User) PublicUser {
	return PublicUser{
		ID:       dbU.ID,
		Username: dbU.Username,
	}
}

func toPublicUser(u PublicUser) db.User {
	return db.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email,omitempty"`
	PasswordHash string     `json:"password_hash,omitempty"`
	Enabled      bool       `json:"enabled"`
	CanLogin     bool       `json:"can_login"`
	JoinDate     time.Time  `json:"join_date"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	LastAccess   *time.Time `json:"last_access,omitempty"`
	Uploaded     int64      `json:"uploaded"`
	Downloaded   int64      `json:"downloaded"`
}

func fromUser(dbU db.User) User {
	u := User{
		ID:           dbU.ID,
		Username:     dbU.Username,
		Email:        dbU.Email,
		PasswordHash: dbU.PasswordHash,
		Enabled:      dbU.Enabled,
		CanLogin:     dbU.CanLogin,
		JoinDate:     dbU.JoinDate,
		Uploaded:     dbU.Uploaded,
		Downloaded:   dbU.Downloaded,
	}
	if dbU.LastLogin.Valid {
		u.LastLogin = &dbU.LastLogin.Time
	}
	if dbU.LastAccess.Valid {
		u.LastAccess = &dbU.LastAccess.Time
	}
	return u
}

func toUser(u User) db.User {
	dbU := db.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Enabled:      u.Enabled,
		CanLogin:     u.CanLogin,
		JoinDate:     u.JoinDate,
		Uploaded:     u.Uploaded,
		Downloaded:   u.Downloaded,
	}
	if u.LastLogin != nil {
		dbU.LastLogin.Valid = true
		dbU.LastLogin.Time = *u.LastLogin
	}
	if u.LastAccess != nil {
		dbU.LastAccess.Valid = true
		dbU.LastAccess.Time = *u.LastAccess
	}
	return dbU
}
