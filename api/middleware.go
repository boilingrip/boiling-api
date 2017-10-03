package api

import (
	"errors"
	"fmt"
	"sort"
	"time"

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
		// TODO distinguish errors
		ctx.Fail(errors.New("invalid token"), iris.StatusUnauthorized)
		return
	}

	err = a.db.UpdateUserSetLastAccess(token.User.ID, time.Now())
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	err = a.db.PopulateUserPrivileges(&token.User)
	if err != nil {
		ctx.Error(err, iris.StatusInternalServerError)
		return
	}

	ctx.user = token.User

	ctx.Next()
}

func (a *API) containsPrivilege(userPrivileges []int, privilege string) (bool, error) {
	a.c.privileges.RLock()
	p, err := a.c.privileges.l.LookUp(privilege)
	a.c.privileges.RUnlock()
	if err != nil {
		return false, err
	}

	i := sort.SearchInts(userPrivileges, p)

	return i < len(userPrivileges) && userPrivileges[i] == p, nil
}

func (a *API) withPrivilege(privileges []string) func(*context) {
	a.c.privileges.RLock()
	defer a.c.privileges.RUnlock()
	for _, p := range privileges {
		has := a.c.privileges.l.Has(p)
		if !has {
			panic(fmt.Sprintf("withPrivilege: cache doesn't know privilege %s", p))
		}
	}

	return func(ctx *context) {
		for _, p := range privileges {
			allowed, err := a.containsPrivilege(ctx.user.Privileges, p)
			if err != nil {
				ctx.Error(err, iris.StatusInternalServerError)
				return
			}
			if !allowed {
				ctx.Application().Logger().Warn(fmt.Sprintf("user %d (%s) tried to access %s, but is missing privilege %s", ctx.user.ID, ctx.user.Username, ctx.GetCurrentRoute().Name(), p))
				ctx.Fail(fmt.Errorf("missing privilege %s", p), iris.StatusForbidden)
				return
			}
		}

		ctx.Next()
	}
}
