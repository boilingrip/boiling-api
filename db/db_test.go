package db

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func cleanDB() (BoilingDB, error) {
	bdb, err := New("boilingtest", "boilingtest", "boilingtest")
	if err != nil {
		return nil, err
	}

	db := bdb.(*DB)

	file, err := ioutil.ReadFile("create.sql")
	if err != nil {
		return nil, err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := db.db.Exec(request)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestCleanDB(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)
	require.NotNil(t, db)
}
