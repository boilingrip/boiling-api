package api

import (
	"errors"

	"github.com/kataras/iris"
)

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (a *API) postLogin(ctx *context) {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	if len(username) == 0 {
		ctx.Fail(errors.New("missing username"), iris.StatusBadRequest)
		return
	}
	if len(password) == 0 {
		ctx.Fail(errors.New("missing password"), iris.StatusBadRequest)
		return
	}

	u, err := a.db.LoginAndGetUser(username, password)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}
	u.PasswordHash = "" // just to be sure

	tok, err := a.db.InsertTokenForUser(*u)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	ctx.Success(LoginResponse{
		User:  fromUser(*u),
		Token: tok.Token,
	})
}
