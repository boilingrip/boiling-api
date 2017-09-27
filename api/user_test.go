package api

import (
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/require"
)

func TestGetUserSelf(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.GET("/users").
		WithHeader("X-User-Token", tc.token).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Keys().ContainsOnly("user")
	user := obj.Value("data").Object().Value("user").Object()
	user.ValueEqual("id", tc.user.ID)
	user.ValueEqual("username", tc.user.Username)
	user.ValueEqual("email", tc.user.Email)
}

func TestGetUser(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.GET("/users/1").
		WithHeader("X-User-Token", tc.token).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Keys().ContainsOnly("user")
	user := obj.Value("data").Object().Value("user").Object()
	user.Keys().ContainsOnly("id", "username", "bio", "joined_at", "uploaded", "downloaded", "enabled")
	user.ValueEqual("id", 1)
}
