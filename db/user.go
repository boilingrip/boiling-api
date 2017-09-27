package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	Bio          sql.NullString
	Enabled      bool
	CanLogin     bool
	JoinedAt     time.Time
	LastLogin    pq.NullTime
	LastAccess   pq.NullTime
	Uploaded     int64
	Downloaded   int64
	Privileges   []int
}

func (db *DB) UpdateUserSetLastLogin(id int, lastLogin time.Time) error {
	if id < 0 {
		return errors.New("invalid id")
	}

	res, err := db.db.Exec("UPDATE users SET last_login = $1, last_access=$1 WHERE id=$2", lastLogin, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("user not found")
	}

	return nil
}

func (db *DB) UpdateUserSetLastAccess(id int, lastAccess time.Time) error {
	if id < 0 {
		return errors.New("invalid id")
	}

	res, err := db.db.Exec("UPDATE users SET last_access=$1 WHERE id=$2", lastAccess, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("user not found")
	}

	return nil
}

func updateUserAddPrivilegesTx(id int, privileges []int, tx *sql.Tx) error {
	for _, p := range privileges {
		res, err := tx.Exec("INSERT INTO users_privileges(uid,privilege) VALUES($1,$2) ON CONFLICT (uid,privilege) DO UPDATE SET privilege=$2", id, p)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return errors.New("user not found? duplicate privilege?")
		}
	}

	return nil
}

func (db *DB) UpdateUserAddPrivileges(id int, privileges []int) error {
	if id < 0 {
		return errors.New("invalid ID")
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	err = updateUserAddPrivilegesTx(id, privileges, tx)
	if err != nil {
		log.Warnln("Rolling back transaction due to error", log.Fields{"err": err})
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) UpdateUserDeltaUpDown(id, deltaUp, deltaDown int) error {
	if id < 0 {
		return errors.New("invalid id")
	}

	res, err := db.db.Exec("UPDATE users SET uploaded = uploaded + $1, downloaded = downloaded + $2 WHERE id = $3", deltaUp, deltaDown, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("user not found")
	}

	return nil
}

func (db *DB) SignUpUser(username, password, email string) error {
	if len(username) == 0 || len(password) == 0 || len(email) == 0 {
		return errors.New("missing username/password/email")
	}

	if !checkPasswordRequirements(password) {
		return errors.New("password does not meet the requirements")
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	res, err := db.db.Exec("INSERT INTO users(username, email, password, enabled, can_login, joined_at) VALUES ($1,$2,$3,TRUE,TRUE,NOW())", username, email, pwHash)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("did not insert")
	}

	return nil
}

func (db *DB) PopulateUserPrivileges(u *User) error {
	if u.ID < 0 {
		return errors.New("invalid ID")
	}

	rows, err := db.db.Query("SELECT privilege FROM users_privileges WHERE uid =$1 ORDER BY privilege ASC", u.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			return err
		}

		u.Privileges = append(u.Privileges, tmp)
	}

	return nil
}

func (db *DB) GetUser(id int) (*User, error) {
	if id < 0 {
		return nil, errors.New("invalid ID")
	}

	row := db.db.QueryRow("SELECT email,username,password,bio,enabled,can_login,joined_at,last_login,last_access,uploaded,downloaded FROM users WHERE id=$1", id)

	user := User{ID: id}
	err := row.Scan(
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Bio,
		&user.Enabled,
		&user.CanLogin,
		&user.JoinedAt,
		&user.LastLogin,
		&user.LastAccess,
		&user.Uploaded,
		&user.Downloaded,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func checkPasswordRequirements(password string) bool {
	// TODO maybe figure out something better, also maybe move this to the API layer
	return len(password) >= 12
}

func (db *DB) LoginAndGetUser(username, password string) (*User, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("missing username/password")
	}

	row := db.db.QueryRow("SELECT id,email,password,bio,enabled,can_login,joined_at,last_access,last_login,uploaded,downloaded FROM users WHERE username = $1", username)

	user := User{Username: username}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.Enabled,
		&user.CanLogin,
		&user.JoinedAt,
		&user.LastAccess,
		&user.LastLogin,
		&user.Uploaded,
		&user.Downloaded)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !user.Enabled {
		return nil, errors.New("user disabled")
	}
	if !user.CanLogin {
		return nil, errors.New("login disabled")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Warn("Bcrypt error: %s", err)
		}
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
