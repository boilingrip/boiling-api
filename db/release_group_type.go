package db

func (db *DB) GetAllReleaseGroupTypes() (map[int]string, error) {
	rows, err := db.db.Query("SELECT id,type FROM release_group_types")
	if err != nil {
		return nil, err
	}

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
		if _, ok := m[tmpI]; ok {
			panic("duplicate key for release group types?")
		}
		m[tmpI] = tmpS
	}

	return m, nil
}
