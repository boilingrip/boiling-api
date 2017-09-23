package db

func (db *DB) GetAllMedia() (map[int]string, error) {
	rows, err := db.db.Query("SELECT id,medium FROM media")
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
		if _, ok := m[tmpI]; ok {
			panic("duplicate key for media?")
		}
		m[tmpI] = tmpS
	}

	return m, nil
}