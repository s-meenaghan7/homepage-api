package visitor

import (
	"context"
	"testing"
)

type mockRepository struct {
	incrementFn func(ctx context.Context, pageID string) (int, error)
	retrieveFn  func(ctx context.Context, pageID string) (int, error)
}

func (m *mockRepository) IncrementCount(ctx context.Context, pageID string) (int, error) {
	return m.incrementFn(ctx, pageID)
}

func (m *mockRepository) RetrieveCount(ctx context.Context, pageID string) (int, error) {
	return m.retrieveFn(ctx, pageID)
}

func TestHandleVisit_Success(t *testing.T) {
	expected := 42
	repo := &mockRepository{
		incrementFn: func(ctx context.Context, pageID string) (int, error) {
			return expected, nil
		},
	}
	svc := NewService(repo)

	count, err := svc.HandleVisit(context.Background(), "/")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != expected {
		t.Fatalf("expected [%d], got [%d]", expected, count)
	}
}

func TestGetVisits_Success(t *testing.T) {
	expected := 42
	repo := &mockRepository{
		retrieveFn: func(ctx context.Context, pageID string) (int, error) {
			return expected, nil
		},
	}
	svc := NewService(repo)

	count, err := svc.GetVisits(context.Background(), "/")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != expected {
		t.Fatalf("expected [%d], got [%d]", expected, count)
	}
}
