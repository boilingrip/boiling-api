package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertGetToken(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	u := User{
		ID: 1,
	}

	token, err := db.InsertTokenForUser(u)
	require.Nil(t, err)
	require.NotNil(t, token)

	require.Equal(t, u.ID, token.User.ID)
	require.NotEmpty(t, token.Token)
	require.NotEmpty(t, token.CreatedAt)
	require.Equal(t, tokenLength*2, len(token.Token))

	token2, err := db.GetToken(token.Token)
	require.Nil(t, err)
	require.NotNil(t, token2)
	require.Equal(t, token.Token, token2.Token)
	require.Equal(t, token.CreatedAt, token2.CreatedAt)
	require.Equal(t, token.User.ID, token2.User.ID)
}
