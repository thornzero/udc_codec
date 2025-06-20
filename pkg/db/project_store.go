package db

import (
	_ "modernc.org/sqlite"
)

type ProjectRecord struct {
	ID          int64
	ProjectName string
	FullBOMFile string
	Validated   bool
}

func (s *Store) InsertProject(project ProjectRecord) (int64, error) {
	result, err := s.DB.Exec(`INSERT INTO projects (project_name, full_bom_file, validated) VALUES (?, ?, ?)`,
		project.ProjectName, project.FullBOMFile, project.Validated)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Store) GetAllProjects() ([]ProjectRecord, error) {
	rows, err := s.DB.Query(`SELECT id, project_name, full_bom_file, validated FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProjectRecord
	for rows.Next() {
		var p ProjectRecord
		if err := rows.Scan(&p.ID, &p.ProjectName, &p.FullBOMFile, &p.Validated); err != nil {
			return nil, err
		}
		records = append(records, p)
	}
	return records, nil
}
