package api

import (
	"errors"
	"strings"

	"github.com/kataras/iris/v12"
)

func (a *API) postSignup(ctx *context) {
	username := ctx.fields.mustGetString("username")
	password := ctx.fields.mustGetString("password")
	email := ctx.fields.mustGetString("email")

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
		ctx.Fail(userError(err, "unable to sign up"), iris.StatusBadRequest)
		return
	}

	ctx.Success(nil)
}
