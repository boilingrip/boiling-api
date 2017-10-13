package api

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/gavv/httpexpect.v1"
)

func TestCacheEndpoints(t *testing.T) {
	tc, err := cleanDBWithLogin()
	require.Nil(t, err)
	_, err = getDefaultAPIWithDB(tc.db)
	require.Nil(t, err)

	e := httpexpect.New(t, "http://localhost:8080")

	testCacheEndpoint(e, tc.token, "formats")
	testCacheEndpoint(e, tc.token, "leech_types")
	testCacheEndpoint(e, tc.token, "media")
	testCacheEndpoint(e, tc.token, "release_group_types")
	testCacheEndpoint(e, tc.token, "release_properties")
	testCacheEndpoint(e, tc.token, "release_roles")
	testCacheEndpoint(e, tc.token, "privileges")
}

func testCacheEndpoint(e *httpexpect.Expect, token, endpoint string) {
	resp := e.GET("/{s}", endpoint).
		WithHeader("X-User-Token", token).
		Expect().Status(200)

	obj := resp.JSON().Object()
	obj.Keys().ContainsOnly("status", "data")
	obj.ValueEqual("status", "success")
	obj.Value("data").Object().Value(endpoint).Array().Length().Gt(0)
}
