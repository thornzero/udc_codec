package db

func (s *Store) Migrate() error {
    stmt := `
    CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        full_tag TEXT UNIQUE,
        system_code TEXT,
        equipment_id TEXT,
        instrument_id TEXT,
        function_code TEXT,
        udc_code TEXT,
        description TEXT
    );
    `
    _, err := s.DB.Exec(stmt)
    return err
}
