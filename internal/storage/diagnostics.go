package storage

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hckim/billtap/internal/diagnostics"
)

var _ diagnostics.Repository = (*SQLiteStore)(nil)

func (s *SQLiteStore) RecordRequestTrace(ctx context.Context, trace diagnostics.RequestTrace) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO request_traces (
		id, method, path, query, status, duration_ms, request_id, idempotency_key,
		request_headers, request_body, response_body, response_object, response_object_id,
		error_type, error_code, error_param, related_ids, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		trace.ID,
		trace.Method,
		trace.Path,
		trace.Query,
		trace.Status,
		trace.DurationMS,
		trace.RequestID,
		trace.IdempotencyKey,
		encodeMap(trace.RequestHeaders),
		trace.RequestBody,
		trace.ResponseBody,
		trace.ResponseObject,
		trace.ResponseObjectID,
		trace.ErrorType,
		trace.ErrorCode,
		trace.ErrorParam,
		encodeStringSlices(trace.RelatedIDs),
		encodeTime(trace.CreatedAt),
	)
	return err
}

func (s *SQLiteStore) ListRequestTraces(ctx context.Context, filter diagnostics.RequestTraceFilter) ([]diagnostics.RequestTrace, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.Method != "" {
		clauses = append(clauses, "method = ?")
		args = append(args, strings.ToUpper(filter.Method))
	}
	if filter.Path != "" {
		clauses = append(clauses, "path LIKE ?")
		args = append(args, "%"+filter.Path+"%")
	}
	if filter.Status != 0 {
		clauses = append(clauses, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.RequestID != "" {
		clauses = append(clauses, "request_id = ?")
		args = append(args, filter.RequestID)
	}
	if filter.IdempotencyKey != "" {
		clauses = append(clauses, "idempotency_key = ?")
		args = append(args, filter.IdempotencyKey)
	}
	if filter.ObjectID != "" {
		clauses = append(clauses, "(response_object_id = ? OR path LIKE ? OR related_ids LIKE ?)")
		args = append(args, filter.ObjectID, "%"+filter.ObjectID+"%", "%"+filter.ObjectID+"%")
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, `SELECT id, method, path, query, status, duration_ms, request_id, idempotency_key,
		request_headers, request_body, response_body, response_object, response_object_id,
		error_type, error_code, error_param, related_ids, created_at
		FROM request_traces WHERE `+strings.Join(clauses, " AND ")+`
		ORDER BY created_at DESC, id DESC LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []diagnostics.RequestTrace
	for rows.Next() {
		trace, err := scanRequestTrace(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, trace)
	}
	return out, rows.Err()
}

func scanRequestTrace(row scanner) (diagnostics.RequestTrace, error) {
	var trace diagnostics.RequestTrace
	var headers, relatedIDs, createdAt string
	if err := row.Scan(
		&trace.ID,
		&trace.Method,
		&trace.Path,
		&trace.Query,
		&trace.Status,
		&trace.DurationMS,
		&trace.RequestID,
		&trace.IdempotencyKey,
		&headers,
		&trace.RequestBody,
		&trace.ResponseBody,
		&trace.ResponseObject,
		&trace.ResponseObjectID,
		&trace.ErrorType,
		&trace.ErrorCode,
		&trace.ErrorParam,
		&relatedIDs,
		&createdAt,
	); err != nil {
		return trace, err
	}
	trace.Object = diagnostics.ObjectRequestTrace
	trace.RequestHeaders = decodeMap(headers)
	trace.RelatedIDs = decodeStringSlices(relatedIDs)
	trace.CreatedAt = decodeTime(createdAt)
	return trace, nil
}

func encodeStringSlices(values map[string][]string) string {
	if values == nil {
		return "{}"
	}
	raw, err := json.Marshal(values)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func decodeStringSlices(raw string) map[string][]string {
	if raw == "" {
		return nil
	}
	var out map[string][]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
