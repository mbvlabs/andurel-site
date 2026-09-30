package factories

import (
	"context"
	"fmt"

	"andurel-site/models"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
)

// DocumentationPageFactory wraps models.DocumentationPage for testing
type DocumentationPageFactory struct {
	models.DocumentationPage
}

type DocumentationPageOption func(*DocumentationPageFactory)

// BuildDocumentationPage creates an in-memory DocumentationPage with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateDocumentationPage.
func BuildDocumentationPage(versionId int64, opts ...DocumentationPageOption) models.DocumentationPage {
	f := &DocumentationPageFactory{
		DocumentationPage: models.DocumentationPage{
			VersionId:           versionId,
			Slug:                faker.Word(),
			PublishedRevisionId: pgtype.Int8{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.DocumentationPage
}

// CreateDocumentationPage creates and persists a DocumentationPage to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateDocumentationPage(ctx context.Context, db storage.Connection, versionId int64, opts ...DocumentationPageOption) (models.DocumentationPage, error) {
	built := BuildDocumentationPage(versionId, opts...)

	return models.NewDocumentationPages(db).Create(ctx, models.CreateDocumentationPageData{
		VersionId:           built.VersionId,
		Slug:                built.Slug,
		PublishedRevisionId: built.PublishedRevisionId,
	})
}

// CreateDocumentationPages creates multiple DocumentationPage records at once
func CreateDocumentationPages(ctx context.Context, db storage.Connection, versionId int64, count int, opts ...DocumentationPageOption) ([]models.DocumentationPage, error) {
	documentationPages := make([]models.DocumentationPage, 0, count)

	for i := range count {
		entity, err := CreateDocumentationPage(ctx, db, versionId, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create documentationpage %d: %w", i+1, err)
		}
		documentationPages = append(documentationPages, entity)
	}

	return documentationPages, nil
}

// Option functions

// WithDocumentationPageVersionId sets the VersionId field
func WithDocumentationPageVersionId(value int64) DocumentationPageOption {
	return func(f *DocumentationPageFactory) {
		f.DocumentationPage.VersionId = value
	}
}

// WithDocumentationPageSlug sets the Slug field
func WithDocumentationPageSlug(value string) DocumentationPageOption {
	return func(f *DocumentationPageFactory) {
		f.DocumentationPage.Slug = value
	}
}

// WithDocumentationPagePublishedRevisionId sets the PublishedRevisionId field
func WithDocumentationPagePublishedRevisionId(value pgtype.Int8) DocumentationPageOption {
	return func(f *DocumentationPageFactory) {
		f.DocumentationPage.PublishedRevisionId = value
	}
}
