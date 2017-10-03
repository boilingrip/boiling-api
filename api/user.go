package api

import (
	"time"

	"github.com/kataras/iris"

	"github.com/boilingrip/boiling-api/db"
)

type BaseUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func baseUserFromDBUser(dbU db.User) BaseUser {
	return BaseUser{
		ID:       dbU.ID,
		Username: dbU.Username,
	}
}

func dbUserFromBaseUser(u BaseUser) db.User {
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
	Bio          *string    `json:"bio,omitempty"`
	Enabled      bool       `json:"enabled"`
	CanLogin     bool       `json:"can_login,omitempty"`
	JoinedAt     time.Time  `json:"joined_at"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	LastAccess   *time.Time `json:"last_access,omitempty"`
	Uploaded     int64      `json:"uploaded"`
	Downloaded   int64      `json:"downloaded"`
}

func userFromDBUser(dbU db.User) User {
	u := User{
		ID:           dbU.ID,
		Username:     dbU.Username,
		Email:        dbU.Email,
		PasswordHash: dbU.PasswordHash,
		Enabled:      dbU.Enabled,
		CanLogin:     dbU.CanLogin,
		JoinedAt:     dbU.JoinedAt,
		Uploaded:     dbU.Uploaded,
		Downloaded:   dbU.Downloaded,
	}
	if dbU.LastLogin.Valid {
		u.LastLogin = &dbU.LastLogin.Time
	}
	if dbU.LastAccess.Valid {
		u.LastAccess = &dbU.LastAccess.Time
	}
	if dbU.Bio.Valid {
		u.Bio = &dbU.Bio.String
	}
	return u
}

func dbUserFromUser(u User) db.User {
	dbU := db.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Enabled:      u.Enabled,
		CanLogin:     u.CanLogin,
		JoinedAt:     u.JoinedAt,
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
	if u.Bio != nil {
		dbU.Bio.Valid = true
		dbU.Bio.String = *u.Bio
	}
	return dbU
}

type UserResponse struct {
	User User `json:"user"`
}

func (a *API) getUserSelf(ctx *context) {
	u, err := a.db.GetUser(ctx.user.ID)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}
	u.PasswordHash = ""

	ctx.Success(UserResponse{userFromDBUser(*u)})
}

func (a *API) getUser(ctx *context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Fail(userError(err, "invalid ID"), iris.StatusBadRequest)
		return
	}

	if id == ctx.user.ID {
		a.getUserSelf(ctx)
		return
	}

	u, err := a.db.GetUser(id)
	if err != nil {
		ctx.Fail(userError(err, "not found"), iris.StatusBadRequest)
		return
	}
	// remove confidential stuff TODO paranoia?
	u.PasswordHash = ""
	u.Email = ""
	u.LastAccess.Valid = false
	u.LastLogin.Valid = false
	u.CanLogin = false

	ctx.Success(UserResponse{userFromDBUser(*u)})
}
