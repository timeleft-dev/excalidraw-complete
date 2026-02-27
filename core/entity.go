package core

import (
	"bytes"
	"context"
)

type (
	Document struct {
		Data bytes.Buffer
	}

	DocumentStore interface {
		FindID(ctx context.Context, id string) (*Document, error)
		Create(ctx context.Context, document *Document) (string, error)
	}

	WorkspaceFileMeta struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		UpdatedAt string `json:"updated_at"`
		UpdatedBy string `json:"updated_by"`
	}

	WorkspaceFile struct {
		WorkspaceFileMeta
		Data string `json:"data"`
	}

	WorkspaceStore interface {
		List(ctx context.Context) ([]WorkspaceFileMeta, error)
		Find(ctx context.Context, id string) (*WorkspaceFile, error)
		// Save creates a new file when id is empty, otherwise upserts.
		Save(ctx context.Context, id, name string, data []byte, updatedBy string) (string, error)
		Delete(ctx context.Context, id string) error
	}
)
