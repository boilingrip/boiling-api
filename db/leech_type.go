package db

func (db *DB) GetAllLeechTypes() (map[int]string, error) {
	rows, err := db.db.Query("SELECT id,type FROM leech_types")
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
