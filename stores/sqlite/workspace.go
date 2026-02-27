package sqlite

import (
	"context"
	"database/sql"
	"excalidraw-complete/core"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type workspaceStore struct {
	db *sql.DB
}

func NewWorkspaceStore(dataSourceName string) core.WorkspaceStore {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS workspace_files (
		id         TEXT PRIMARY KEY,
		name       TEXT NOT NULL,
		data       BLOB NOT NULL,
		updated_at DATETIME NOT NULL,
		updated_by TEXT NOT NULL DEFAULT ''
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return &workspaceStore{db: db}
}

func (s *workspaceStore) List(ctx context.Context) ([]core.WorkspaceFileMeta, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, updated_at, updated_by FROM workspace_files ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []core.WorkspaceFileMeta{}
	for rows.Next() {
		var f core.WorkspaceFileMeta
		if err := rows.Scan(&f.ID, &f.Name, &f.UpdatedAt, &f.UpdatedBy); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}

func (s *workspaceStore) Find(ctx context.Context, id string) (*core.WorkspaceFile, error) {
	var f core.WorkspaceFile
	var data []byte
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, data, updated_at, updated_by FROM workspace_files WHERE id = ?", id).
		Scan(&f.ID, &f.Name, &data, &f.UpdatedAt, &f.UpdatedBy)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workspace file %s not found", id)
	}
	if err != nil {
		return nil, err
	}
	f.Data = string(data)
	return &f, nil
}

func (s *workspaceStore) Save(ctx context.Context, id, name string, data []byte, updatedBy string) (string, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	if id == "" {
		id = ulid.Make().String()
		logrus.WithFields(logrus.Fields{"id": id, "name": name}).Info("Creating workspace file")
		_, err := s.db.ExecContext(ctx,
			"INSERT INTO workspace_files (id, name, data, updated_at, updated_by) VALUES (?, ?, ?, ?, ?)",
			id, name, data, now, updatedBy)
		return id, err
	}

	logrus.WithFields(logrus.Fields{"id": id, "name": name}).Info("Updating workspace file")
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO workspace_files (id, name, data, updated_at, updated_by) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name       = excluded.name,
			data       = excluded.data,
			updated_at = excluded.updated_at,
			updated_by = excluded.updated_by`,
		id, name, data, now, updatedBy)
	return id, err
}

func (s *workspaceStore) Delete(ctx context.Context, id string) error {
	logrus.WithField("id", id).Info("Deleting workspace file")
	_, err := s.db.ExecContext(ctx, "DELETE FROM workspace_files WHERE id = ?", id)
	return err
}
