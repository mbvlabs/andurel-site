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

// DocumentationNavRevisionFactory wraps models.DocumentationNavRevision for testing
type DocumentationNavRevisionFactory struct {
	models.DocumentationNavRevision
}

type DocumentationNavRevisionOption func(*DocumentationNavRevisionFactory)

// BuildDocumentationNavRevision creates an in-memory DocumentationNavRevision with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateDocumentationNavRevision.
func BuildDocumentationNavRevision(versionId int64, createdBy uuid.UUID, opts ...DocumentationNavRevisionOption) models.DocumentationNavRevision {
	f := &DocumentationNavRevisionFactory{
		DocumentationNavRevision: models.DocumentationNavRevision{
			VersionId:   versionId,
			Status:      faker.Word(),
			Tree:        models.EmptyJSONArray,
			CreatedBy:   createdBy,
			PublishedAt: pgtype.Timestamptz{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.DocumentationNavRevision
}

// CreateDocumentationNavRevision creates and persists a DocumentationNavRevision to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateDocumentationNavRevision(ctx context.Context, db storage.Connection, versionId int64, createdBy uuid.UUID, opts ...DocumentationNavRevisionOption) (models.DocumentationNavRevision, error) {
	built := BuildDocumentationNavRevision(versionId, createdBy, opts...)

	return models.NewDocumentationNavRevisions(db).Create(ctx, models.CreateDocumentationNavRevisionData{
		VersionId:   built.VersionId,
		Status:      built.Status,
		Tree:        built.Tree,
		CreatedBy:   built.CreatedBy,
		PublishedAt: built.PublishedAt,
	})
}

// CreateDocumentationNavRevisions creates multiple DocumentationNavRevision records at once
func CreateDocumentationNavRevisions(ctx context.Context, db storage.Connection, versionId int64, createdBy uuid.UUID, count int, opts ...DocumentationNavRevisionOption) ([]models.DocumentationNavRevision, error) {
	documentationNavRevisions := make([]models.DocumentationNavRevision, 0, count)

	for i := range count {
		entity, err := CreateDocumentationNavRevision(ctx, db, versionId, createdBy, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create documentationnavrevision %d: %w", i+1, err)
		}
		documentationNavRevisions = append(documentationNavRevisions, entity)
	}

	return documentationNavRevisions, nil
}

// Option functions

// WithDocumentationNavRevisionVersionId sets the VersionId field
func WithDocumentationNavRevisionVersionId(value int64) DocumentationNavRevisionOption {
	return func(f *DocumentationNavRevisionFactory) {
		f.DocumentationNavRevision.VersionId = value
	}
}

// WithDocumentationNavRevisionStatus sets the Status field
func WithDocumentationNavRevisionStatus(value string) DocumentationNavRevisionOption {
	return func(f *DocumentationNavRevisionFactory) {
		f.DocumentationNavRevision.Status = value
	}
}

// WithDocumentationNavRevisionTree sets the Tree field
func WithDocumentationNavRevisionTree(value []byte) DocumentationNavRevisionOption {
	return func(f *DocumentationNavRevisionFactory) {
		f.DocumentationNavRevision.Tree = value
	}
}

// WithDocumentationNavRevisionCreatedBy sets the CreatedBy field
func WithDocumentationNavRevisionCreatedBy(value uuid.UUID) DocumentationNavRevisionOption {
	return func(f *DocumentationNavRevisionFactory) {
		f.DocumentationNavRevision.CreatedBy = value
	}
}

// WithDocumentationNavRevisionPublishedAt sets the PublishedAt field
func WithDocumentationNavRevisionPublishedAt(value pgtype.Timestamptz) DocumentationNavRevisionOption {
	return func(f *DocumentationNavRevisionFactory) {
		f.DocumentationNavRevision.PublishedAt = value
	}
}
