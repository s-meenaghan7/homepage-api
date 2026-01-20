package visitor

import (
	"context"
)

type VisitorService interface {
	GetVisits(ctx context.Context, pageID string) (int, error)
	HandleVisit(ctx context.Context, pageID string) (int, error)
}

type Service struct {
	repo VisitorRepository
}

func NewService(repo VisitorRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetVisits(ctx context.Context, pageID string) (int, error) {
	return s.repo.RetrieveCount(ctx, pageID)
}

func (s *Service) HandleVisit(ctx context.Context, pageID string) (int, error) {
	return s.repo.IncrementCount(ctx, pageID)
}
