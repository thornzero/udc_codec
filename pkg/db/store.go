package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func OpenDB(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) InsertTag(t *TagRecord) error {
	_, err := s.DB.Exec(`
        INSERT INTO tags 
        (full_tag, system_code, equipment_id, instrument_id, function_code, udc_code, description) 
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.FullTag, t.SystemCode, t.EquipmentID, t.InstrumentID, t.FunctionCode, t.UDCCode, t.Description)
	return err
}

func (s *Store) LookupTag(fulltag string) (*TagRecord, error) {
	row := s.DB.QueryRow(`SELECT id, full_tag, system_code, equipment_id, instrument_id, function_code, udc_code, description FROM tags WHERE full_tag = ?`, fulltag)
	var t TagRecord
	if err := row.Scan(&t.ID, &t.FullTag, &t.SystemCode, &t.EquipmentID, &t.InstrumentID, &t.FunctionCode, &t.UDCCode, &t.Description); err != nil {
		return nil, err
	}
	return &t, nil
}
