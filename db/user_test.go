package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignUpUser(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	err = db.SignUpUser("testuser", "testtest12345", "test@example.com")
	require.Nil(t, err)

	u, err := db.LoginAndGetUser("testuser", "testtest12345")
	require.Nil(t, err)
	require.NotNil(t, u)

	require.Equal(t, "test@example.com", u.Email)
	require.Equal(t, true, u.Enabled)
	require.Equal(t, true, u.CanLogin)
	require.False(t, u.LastAccess.Valid)
	require.False(t, u.LastLogin.Valid)

	u2, err := db.GetUser(u.ID)
	require.Nil(t, err)
	require.Equal(t, u, u2)
}

func TestUpdateUserDeltaUpDown(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	err = db.SignUpUser("testuser", "testpwtest1234", "test@example.com")
	require.Nil(t, err)

	u, err := db.LoginAndGetUser("testuser", "testpwtest1234")
	require.Nil(t, err)
	require.NotNil(t, u)

	err = db.UpdateUserDeltaUpDown(u.ID, 513, -234)
	require.Nil(t, err)

	u, err = db.GetUser(u.ID)
	require.Nil(t, err)
	require.Equal(t, int64(513), u.Uploaded)
	require.Equal(t, int64(-234), u.Downloaded)
}
