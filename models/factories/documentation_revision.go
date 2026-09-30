package factories

import (
	"context"
	"fmt"
	"uuid"

	"andurel-site/models"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
)

// DocumentationRevisionFactory wraps models.DocumentationRevision for testing
type DocumentationRevisionFactory struct {
	models.DocumentationRevision
}

type DocumentationRevisionOption func(*DocumentationRevisionFactory)

// BuildDocumentationRevision creates an in-memory DocumentationRevision with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateDocumentationRevision.
func BuildDocumentationRevision(pageId int64, createdBy uuid.UUID, opts ...DocumentationRevisionOption) models.DocumentationRevision {
	f := &DocumentationRevisionFactory{
		DocumentationRevision: models.DocumentationRevision{
			PageId:       pageId,
			Status:       faker.Word(),
			Title:        faker.Word(),
			MetaTitle:    faker.Word(),
			Description:  faker.Sentence(),
			BodyMarkdown: faker.Word(),
			Headings:     models.EmptyJSONArray,
			CreatedBy:    createdBy,
			PublishedAt:  pgtype.Timestamptz{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.DocumentationRevision
}

// CreateDocumentationRevision creates and persists a DocumentationRevision to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateDocumentationRevision(ctx context.Context, db storage.Connection, pageId int64, createdBy uuid.UUID, opts ...DocumentationRevisionOption) (models.DocumentationRevision, error) {
	built := BuildDocumentationRevision(pageId, createdBy, opts...)

	return models.NewDocumentationRevisions(db).Create(ctx, models.CreateDocumentationRevisionData{
		PageId:       built.PageId,
		Status:       built.Status,
		Title:        built.Title,
		MetaTitle:    built.MetaTitle,
		Description:  built.Description,
		BodyMarkdown: built.BodyMarkdown,
		Headings:     built.Headings,
		CreatedBy:    built.CreatedBy,
		PublishedAt:  built.PublishedAt,
	})
}

// CreateDocumentationRevisions creates multiple DocumentationRevision records at once
func CreateDocumentationRevisions(ctx context.Context, db storage.Connection, pageId int64, createdBy uuid.UUID, count int, opts ...DocumentationRevisionOption) ([]models.DocumentationRevision, error) {
	documentationRevisions := make([]models.DocumentationRevision, 0, count)

	for i := range count {
		entity, err := CreateDocumentationRevision(ctx, db, pageId, createdBy, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create documentationrevision %d: %w", i+1, err)
		}
		documentationRevisions = append(documentationRevisions, entity)
	}

	return documentationRevisions, nil
}

// Option functions

// WithDocumentationRevisionPageId sets the PageId field
func WithDocumentationRevisionPageId(value int64) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.PageId = value
	}
}

// WithDocumentationRevisionStatus sets the Status field
func WithDocumentationRevisionStatus(value string) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.Status = value
	}
}

// WithDocumentationRevisionTitle sets the Title field
func WithDocumentationRevisionTitle(value string) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.Title = value
	}
}

func WithDocumentationRevisionMetaTitle(value string) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.MetaTitle = value
	}
}

// WithDocumentationRevisionDescription sets the Description field
func WithDocumentationRevisionDescription(value string) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.Description = value
	}
}

// WithDocumentationRevisionBodyMarkdown sets the BodyMarkdown field
func WithDocumentationRevisionBodyMarkdown(value string) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.BodyMarkdown = value
	}
}

// WithDocumentationRevisionHeadings sets the Headings field
func WithDocumentationRevisionHeadings(value []byte) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.Headings = value
	}
}

// WithDocumentationRevisionCreatedBy sets the CreatedBy field
func WithDocumentationRevisionCreatedBy(value uuid.UUID) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.CreatedBy = value
	}
}

// WithDocumentationRevisionPublishedAt sets the PublishedAt field
func WithDocumentationRevisionPublishedAt(value pgtype.Timestamptz) DocumentationRevisionOption {
	return func(f *DocumentationRevisionFactory) {
		f.DocumentationRevision.PublishedAt = value
	}
}
