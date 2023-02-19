package movie

import (
	"time"

	"github.com/turao/topics/metadata"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type Movie interface {
	ID() ID
	Title() string

	Delete()
	metadata.Auditable
	metadata.MultiTenant
}

type movie struct {
	cfg config
}

var _ Movie = (*movie)(nil)

func NewMovie(cfg config) *movie {
	return &movie{
		cfg: cfg,
	}
}

func (m *movie) ID() ID {
	return m.cfg.id
}

func (m *movie) Title() string {
	return m.cfg.title
}

func (m *movie) Tenancy() metadata.Tenancy {
	return m.cfg.tenancy
}

func (m *movie) CreatedAt() time.Time {
	return m.cfg.createdAt
}

func (m *movie) DeletedAt() *time.Time {
	return m.cfg.deletedAt
}

func (m *movie) Delete() {
	if m.DeletedAt() == nil {
		now := time.Now()
		m.cfg.deletedAt = &now
	}
}
