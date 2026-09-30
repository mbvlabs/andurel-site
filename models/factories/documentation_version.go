package factories

import (
	"context"
	"fmt"

	"andurel-site/models"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
)

// DocumentationVersionFactory wraps models.DocumentationVersion for testing
type DocumentationVersionFactory struct {
	models.DocumentationVersion
}

type DocumentationVersionOption func(*DocumentationVersionFactory)

// BuildDocumentationVersion creates an in-memory DocumentationVersion with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateDocumentationVersion.
func BuildDocumentationVersion(opts ...DocumentationVersionOption) models.DocumentationVersion {
	f := &DocumentationVersionFactory{
		DocumentationVersion: models.DocumentationVersion{
			Slug:           faker.Word(),
			Label:          faker.Word(),
			IsLatest:       randomBool(),
			Position:       randomInt(1, 1000, 100),
			PublishedNavId: pgtype.Int8{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.DocumentationVersion
}

// CreateDocumentationVersion creates and persists a DocumentationVersion to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateDocumentationVersion(ctx context.Context, db storage.Connection, opts ...DocumentationVersionOption) (models.DocumentationVersion, error) {
	built := BuildDocumentationVersion(opts...)

	return models.NewDocumentationVersions(db).Create(ctx, models.CreateDocumentationVersionData{
		Slug:           built.Slug,
		Label:          built.Label,
		IsLatest:       built.IsLatest,
		Position:       built.Position,
		PublishedNavId: built.PublishedNavId,
	})
}

// CreateDocumentationVersions creates multiple DocumentationVersion records at once
func CreateDocumentationVersions(ctx context.Context, db storage.Connection, count int, opts ...DocumentationVersionOption) ([]models.DocumentationVersion, error) {
	documentationVersions := make([]models.DocumentationVersion, 0, count)

	for i := range count {
		entity, err := CreateDocumentationVersion(ctx, db, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create documentationversion %d: %w", i+1, err)
		}
		documentationVersions = append(documentationVersions, entity)
	}

	return documentationVersions, nil
}

// Option functions

// WithDocumentationVersionSlug sets the Slug field
func WithDocumentationVersionSlug(value string) DocumentationVersionOption {
	return func(f *DocumentationVersionFactory) {
		f.DocumentationVersion.Slug = value
	}
}

// WithDocumentationVersionLabel sets the Label field
func WithDocumentationVersionLabel(value string) DocumentationVersionOption {
	return func(f *DocumentationVersionFactory) {
		f.DocumentationVersion.Label = value
	}
}

// WithDocumentationVersionIsLatest sets the IsLatest field
func WithDocumentationVersionIsLatest(value bool) DocumentationVersionOption {
	return func(f *DocumentationVersionFactory) {
		f.DocumentationVersion.IsLatest = value
	}
}

// WithDocumentationVersionPosition sets the Position field
func WithDocumentationVersionPosition(value int32) DocumentationVersionOption {
	return func(f *DocumentationVersionFactory) {
		f.DocumentationVersion.Position = value
	}
}

// WithDocumentationVersionPublishedNavId sets the PublishedNavId field
func WithDocumentationVersionPublishedNavId(value pgtype.Int8) DocumentationVersionOption {
	return func(f *DocumentationVersionFactory) {
		f.DocumentationVersion.PublishedNavId = value
	}
}
