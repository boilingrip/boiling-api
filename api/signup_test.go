package api

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/gavv/httpexpect.v1"
)

func TestSignup(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	resp := e.POST("/signup").
		WithFormField("username", "abc").
		WithFormField("password", "pass123pass123").
		WithFormField("email", "abc@some.example.org").
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status")
	obj.ValueEqual("status", "success")

	resp = e.POST("/login").
		WithFormField("username", "abc").
		WithFormField("password", "pass123pass123").
		Expect().Status(200)

	obj = resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	data := obj.Value("data").Object()
	data.Keys().ContainsOnly("user", "token")
	user := data.Value("user").Object()
	user.ValueEqual("username", "abc")
	user.ValueEqual("email", "abc@some.example.org")
}
