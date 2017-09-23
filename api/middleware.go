package api

import (
	"errors"

	"github.com/kataras/iris"
)

func (a *API) withLogin(ctx *context) {
	tokenString := ctx.GetHeader("X-User-Token")
	if tokenString == "" {
		ctx.Fail(errors.New("missing token"), iris.StatusUnauthorized)
		return
	}

	token, err := a.db.GetToken(tokenString)
	if err != nil {
		ctx.Fail(errors.New("invalid token"), iris.StatusUnauthorized)
		return
	}

	ctx.user = token.User
	ctx.loggedIn = true

	ctx.Next()
}
