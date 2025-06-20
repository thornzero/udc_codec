package db

func (s *Store) Migrate() error {
	_, err := s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_name TEXT UNIQUE,
		full_bom_file TEXT,
		validated BOOLEAN
	);
	`)
	return err
}
