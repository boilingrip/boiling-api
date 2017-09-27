package db

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type Passkey struct {
	Uid       int
	Passkey   string
	CreatedAt time.Time
	Valid     bool
}

const passkeyLength = 64

func (db *DB) GetPasskeyForUser(id int) (*Passkey, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	passkey := Passkey{
		Uid:   id,
		Valid: true,
	}
	err := db.db.QueryRow("SELECT passkey,created_at FROM user_passkeys WHERE uid = $1 AND valid=TRUE", id).Scan(
		&passkey.Passkey,
		&passkey.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &passkey, nil
}

func (db *DB) GetAllPasskeysForUser(id int) ([]Passkey, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	rows, err := db.db.Query("SELECT passkey,created_at,valid FROM user_passkeys WHERE uid = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passkeys []Passkey
	for rows.Next() {
		tmp := Passkey{
			Uid: id,
		}
		err = rows.Scan(
			&tmp.Passkey,
			&tmp.CreatedAt,
			&tmp.Valid)
		if err != nil {
			return nil, err
		}

		passkeys = append(passkeys, tmp)
	}

	return passkeys, nil
}

func generateNewPasskeyForUserTx(id int, passkey string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE user_passkeys SET valid=FALSE WHERE uid=$1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO user_passkeys(uid,passkey,created_at,valid) VALUES ($1,$2,now(),TRUE)", id, passkey)
	return err
}

func (db *DB) GenerateNewPasskeyForUser(id int) (string, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return "", err
	}

	passkey := generateRandomAlphanumeric(passkeyLength)

	err = generateNewPasskeyForUserTx(id, passkey, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return passkey, nil
}
