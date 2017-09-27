package db

type Format struct {
	Format   string
	Encoding string
}

func (db *DB) GetAllFormats() (map[int]Format, error) {
	rows, err := db.db.Query("SELECT id,format,encoding FROM formats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int]Format)
	for rows.Next() {
		var (
			tmpI       int
			tmpF, tmpE string
		)
		err = rows.Scan(&tmpI, &tmpF, &tmpE)
		if err != nil {
			return nil, err
		}

		if _, ok := m[tmpI]; ok {
			panic("duplicate key in formats")
		}

		m[tmpI] = Format{tmpF, tmpE}
	}

	return m, nil
}
