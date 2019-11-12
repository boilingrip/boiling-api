package api

import (
	"time"

	"github.com/kataras/iris/v12"
)

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (a *API) postLogin(ctx *context) {
	username := ctx.fields.mustGetString("username")
	password := ctx.fields.mustGetString("password")

	u, err := a.db.LoginAndGetUser(username, password)
	if err != nil {
		ctx.Fail(userError(err, "unable to log in"), iris.StatusBadRequest)
		return
	}
	u.PasswordHash = "" // just to be sure

	err = a.db.UpdateUserSetLastLogin(u.ID, time.Now())
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	tok, err := a.db.InsertTokenForUser(*u)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	ctx.Success(LoginResponse{
		User:  userFromDBUser(*u),
		Token: tok.Token,
	})
}
