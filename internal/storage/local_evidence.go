package storage

import (
	"context"
	"errors"
	"fmt"
)

// SaveLocalEvidence upserts one evidence document.
func (s *SQLiteStore) SaveLocalEvidence(ctx context.Context, kind, id, data string) error {
	if s == nil || s.db == nil {
		return errors.New("sqlite store is not open")
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO local_evidence (kind, id, data)
		VALUES (?, ?, ?)
		ON CONFLICT(kind, id) DO UPDATE SET data = excluded.data`, kind, id, data); err != nil {
		return fmt.Errorf("save local evidence %s/%s: %w", kind, id, err)
	}
	return nil
}

// DeleteLocalEvidence removes one evidence document. Deleting an absent row is not an error.
func (s *SQLiteStore) DeleteLocalEvidence(ctx context.Context, kind, id string) error {
	if s == nil || s.db == nil {
		return errors.New("sqlite store is not open")
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM local_evidence WHERE kind = ? AND id = ?`, kind, id); err != nil {
		return fmt.Errorf("delete local evidence %s/%s: %w", kind, id, err)
	}
	return nil
}

// LoadLocalEvidence returns every stored evidence document as kind -> id -> JSON.
func (s *SQLiteStore) LoadLocalEvidence(ctx context.Context) (map[string]map[string]string, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("sqlite store is not open")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT kind, id, data FROM local_evidence ORDER BY kind, id`)
	if err != nil {
		return nil, fmt.Errorf("load local evidence: %w", err)
	}
	defer rows.Close()

	out := map[string]map[string]string{}
	for rows.Next() {
		var kind, id, data string
		if err := rows.Scan(&kind, &id, &data); err != nil {
			return nil, fmt.Errorf("scan local evidence: %w", err)
		}
		if out[kind] == nil {
			out[kind] = map[string]string{}
		}
		out[kind][id] = data
	}
	return out, rows.Err()
}
