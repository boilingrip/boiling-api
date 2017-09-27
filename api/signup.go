package api

import (
	"errors"
	"strings"

	"github.com/kataras/iris"
)

func (a *API) postSignup(ctx *context) {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	email := ctx.FormValue("email")
	if len(username) == 0 {
		ctx.Fail(errors.New("missing username"), iris.StatusBadRequest)
		return
	}
	if len(password) == 0 {
		ctx.Fail(errors.New("missing password"), iris.StatusBadRequest)
		return
	}
	if len(email) == 0 {
		ctx.Fail(errors.New("missing email"), iris.StatusBadRequest)
		return
	}
	if !strings.Contains(email, "@") || sanitizeString(email) != email {
		ctx.Fail(errors.New("invalid email"), iris.StatusBadRequest)
		return
	}
	if sanitizeString(username) != username {
		ctx.Fail(errors.New("invalid username"), iris.StatusBadRequest)
		return
	}
	if strings.TrimSpace(password) != password {
		ctx.Fail(errors.New("invalid password"), iris.StatusBadRequest)
		return
	}

	err := a.db.SignUpUser(username, password, email)
	if err != nil {
		ctx.Fail(err, iris.StatusBadRequest)
		return
	}

	ctx.Success(nil)
}
