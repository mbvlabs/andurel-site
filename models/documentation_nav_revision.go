package models

// andurel:table documentation_nav_revisions

import (
	"andurel-site/models/internal/queries"
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type DocumentationNavRevisions struct {
	queries *queries.Queries
}

func NewDocumentationNavRevisions(db storage.Connection) DocumentationNavRevisions {
	return DocumentationNavRevisions{queries: queries.New(db)}
}

// WithTx returns a copy that runs queries inside tx.
func (dnr DocumentationNavRevisions) WithTx(tx storage.Transaction) DocumentationNavRevisions {
	return DocumentationNavRevisions{queries: queries.New(tx)}
}

type DocumentationNavRevision struct {
	ID          int64              `andurel:"id"`
	CreatedAt   pgtype.Timestamptz `andurel:"created_at"`
	VersionId   int64              `andurel:"version_id"`
	Status      string             `andurel:"status"`
	Tree        []byte             `andurel:"tree"`
	CreatedBy   uuid.UUID          `andurel:"created_by"`
	PublishedAt pgtype.Timestamptz `andurel:"published_at"`
}

func (dnr DocumentationNavRevisions) Find(ctx context.Context, id int64) (DocumentationNavRevision, error) {
	entity, err := dnr.queries.GetDocumentationNavRevision[DocumentationNavRevision](ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationNavRevision{}, ErrNotFound
		}
		return DocumentationNavRevision{}, err
	}

	return entity, nil
}

type CreateDocumentationNavRevisionData struct {
	VersionId   int64
	Status      string
	Tree        []byte
	CreatedBy   uuid.UUID
	PublishedAt pgtype.Timestamptz
}

func (dnr DocumentationNavRevisions) Create(ctx context.Context, data CreateDocumentationNavRevisionData) (DocumentationNavRevision, error) {
	entity := DocumentationNavRevision{
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		VersionId:   data.VersionId,
		Status:      data.Status,
		Tree:        data.Tree,
		CreatedBy:   data.CreatedBy,
		PublishedAt: data.PublishedAt,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationNavRevision{}, errors.Join(ErrDomainValidation, err)
	}

	return dnr.queries.CreateDocumentationNavRevision[DocumentationNavRevision](ctx, queries.CreateDocumentationNavRevisionParams{
		CreatedAt:   entity.CreatedAt,
		VersionID:   entity.VersionId,
		Status:      entity.Status,
		Tree:        entity.Tree,
		CreatedBy:   entity.CreatedBy,
		PublishedAt: entity.PublishedAt,
	})
}

type UpdateDocumentationNavRevisionData struct {
	ID          int64
	VersionId   int64
	Status      string
	Tree        []byte
	CreatedBy   uuid.UUID
	PublishedAt pgtype.Timestamptz
}

func (dnr DocumentationNavRevisions) Update(ctx context.Context, data UpdateDocumentationNavRevisionData) (DocumentationNavRevision, error) {
	entity := DocumentationNavRevision{
		ID:          data.ID,
		VersionId:   data.VersionId,
		Status:      data.Status,
		Tree:        data.Tree,
		CreatedBy:   data.CreatedBy,
		PublishedAt: data.PublishedAt,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationNavRevision{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := dnr.queries.UpdateDocumentationNavRevision[DocumentationNavRevision](ctx, queries.UpdateDocumentationNavRevisionParams{
		ID:          entity.ID,
		VersionID:   entity.VersionId,
		Status:      entity.Status,
		Tree:        entity.Tree,
		CreatedBy:   entity.CreatedBy,
		PublishedAt: entity.PublishedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationNavRevision{}, ErrNotFound
		}
		return DocumentationNavRevision{}, err
	}

	return row, nil
}

func (dnr DocumentationNavRevisions) Destroy(ctx context.Context, id int64) error {
	if err := dnr.queries.DeleteDocumentationNavRevision(ctx, id); err != nil {
		return err
	}

	return nil
}

func (dnr DocumentationNavRevisions) All(ctx context.Context) ([]DocumentationNavRevision, error) {
	return dnr.queries.ListDocumentationNavRevisions[DocumentationNavRevision](ctx).All()
}

type PaginatedDocumentationNavRevisions struct {
	DocumentationNavRevisions []DocumentationNavRevision
	TotalCount                int64
	Page                      int64
	PageSize                  int64
	TotalPages                int64
}

func (dnr DocumentationNavRevisions) Paginate(ctx context.Context, page, pageSize int64) (PaginatedDocumentationNavRevisions, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := dnr.queries.CountDocumentationNavRevisions(ctx)
	if err != nil {
		return PaginatedDocumentationNavRevisions{}, err
	}

	entities, err := dnr.queries.ListDocumentationNavRevisions[DocumentationNavRevision](ctx).
		Limit(int32(pageSize)).
		Offset(int32(offset)).
		All()
	if err != nil {
		return PaginatedDocumentationNavRevisions{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedDocumentationNavRevisions{
		DocumentationNavRevisions: entities,
		TotalCount:                int64(totalCount),
		Page:                      page,
		PageSize:                  pageSize,
		TotalPages:                totalPages,
	}, nil
}

func (dnr DocumentationNavRevisions) Upsert(ctx context.Context, id int64, data CreateDocumentationNavRevisionData) (DocumentationNavRevision, error) {
	entity := DocumentationNavRevision{
		ID:          id,
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		VersionId:   data.VersionId,
		Status:      data.Status,
		Tree:        data.Tree,
		CreatedBy:   data.CreatedBy,
		PublishedAt: data.PublishedAt,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationNavRevision{}, errors.Join(ErrDomainValidation, err)
	}

	return dnr.queries.UpsertDocumentationNavRevision[DocumentationNavRevision](ctx, queries.UpsertDocumentationNavRevisionParams{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		VersionID:   entity.VersionId,
		Status:      entity.Status,
		Tree:        entity.Tree,
		CreatedBy:   entity.CreatedBy,
		PublishedAt: entity.PublishedAt,
	})
}

func (dnr DocumentationNavRevisions) FindDraftByVersion(ctx context.Context, versionID int64) (DocumentationNavRevision, error) {
	entity, err := dnr.queries.GetDocumentationNavRevisionDraftByVersion[DocumentationNavRevision](ctx, versionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationNavRevision{}, ErrNotFound
		}
		return DocumentationNavRevision{}, err
	}

	return entity, nil
}

func (dnr DocumentationNavRevisions) ListByVersion(ctx context.Context, versionID int64) ([]DocumentationNavRevision, error) {
	return dnr.queries.ListDocumentationNavRevisionsByVersion[DocumentationNavRevision](ctx, versionID)
}
