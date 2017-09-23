package api

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/gavv/httpexpect.v1"
)

func TestLogin(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.POST("/login").
		WithFormField("username", tc.user.Username).
		WithFormField("password", tc.password).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	data := obj.Value("data").Object()
	data.Keys().ContainsOnly("user", "token")
	user := data.Value("user").Object()
	user.ValueEqual("username", tc.user.Username)
	user.ValueEqual("id", tc.user.ID)
	user.ValueEqual("email", tc.user.Email)

	// TODO use the token for something
}
