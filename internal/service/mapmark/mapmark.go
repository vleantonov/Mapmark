package mapmark

import (
	"context"
	"echoFramework/internal/domain"
)

type MarkRepository interface {
	Save(ctx context.Context, mark domain.Mark) (int64, error)
	Update(ctx context.Context, mark domain.Mark) error
	Get(ctx context.Context) ([]domain.Mark, error)
	GetById(ctx context.Context, id int) (domain.Mark, error)
}

type Mark struct {
	mr MarkRepository
}

func New(mr MarkRepository) *Mark {
	return &Mark{mr: mr}
}

func (m *Mark) Get(ctx context.Context) ([]domain.Mark, error) {
	return m.mr.Get(ctx)
}

func (m *Mark) GetById(ctx context.Context, id int) (domain.Mark, error) {
	return m.mr.GetById(ctx, id)
}

func (m *Mark) Save(ctx context.Context, mark domain.Mark) (int64, error) {

	id, err := m.mr.Save(ctx, mark)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Mark) Update(ctx context.Context, mark domain.Mark) error {
	return m.mr.Update(ctx, mark)
}
