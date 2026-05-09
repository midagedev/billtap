package diagnostics

import (
	"context"
	"strings"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) RecordRequestTrace(ctx context.Context, trace RequestTrace) error {
	if s == nil || s.repo == nil {
		return nil
	}
	if strings.TrimSpace(trace.ID) == "" {
		trace.ID = "rt_" + strings.ReplaceAll(s.now().Format("20060102150405.000000000"), ".", "")
	}
	trace.Object = ObjectRequestTrace
	if trace.CreatedAt.IsZero() {
		trace.CreatedAt = s.now()
	}
	return s.repo.RecordRequestTrace(ctx, trace)
}

func (s *Service) ListRequestTraces(ctx context.Context, filter RequestTraceFilter) ([]RequestTrace, error) {
	if s == nil || s.repo == nil {
		return []RequestTrace{}, nil
	}
	if filter.Limit == 0 {
		filter.Limit = 100
	}
	return s.repo.ListRequestTraces(ctx, filter)
}
