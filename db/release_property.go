package db

import "errors"

func (db *DB) GetAllReleaseProperties() (map[int]string, error) {
	rows, err := db.db.Query("SELECT id,property FROM release_properties")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int]string)
	for rows.Next() {
		var (
			tmpI int
			tmpS string
		)
		err = rows.Scan(&tmpI, &tmpS)
		if err != nil {
			return nil, err
		}

		m[tmpI] = tmpS
	}

	return m, nil
}

func (db *DB) AddReleaseProperty(key string) error {
	if len(key) == 0 {
		return errors.New("missing key")
	}

	res, err := db.db.Exec("INSERT INTO release_properties(property) VALUES ($1)", key)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("did not insert")
	}

	return nil
}
