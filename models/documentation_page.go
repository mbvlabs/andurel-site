package models

// andurel:table documentation_pages

import (
	"andurel-site/models/internal/queries"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type DocumentationPages struct {
	queries *queries.Queries
}

func NewDocumentationPages(db storage.Connection) DocumentationPages {
	return DocumentationPages{queries: queries.New(db)}
}

// WithTx returns a copy that runs queries inside tx.
func (dp DocumentationPages) WithTx(tx storage.Transaction) DocumentationPages {
	return DocumentationPages{queries: queries.New(tx)}
}

type DocumentationPage struct {
	ID                  int64              `andurel:"id"`
	CreatedAt           pgtype.Timestamptz `andurel:"created_at"`
	UpdatedAt           pgtype.Timestamptz `andurel:"updated_at"`
	VersionId           int64              `andurel:"version_id"`
	Slug                string             `andurel:"slug"`
	PublishedRevisionId pgtype.Int8        `andurel:"published_revision_id"`
}

func (dp DocumentationPages) Find(ctx context.Context, id int64) (DocumentationPage, error) {
	entity, err := dp.queries.GetDocumentationPage[DocumentationPage](ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationPage{}, ErrNotFound
		}
		return DocumentationPage{}, err
	}

	return entity, nil
}

type CreateDocumentationPageData struct {
	VersionId           int64
	Slug                string
	PublishedRevisionId pgtype.Int8
}

func (dp DocumentationPages) Create(ctx context.Context, data CreateDocumentationPageData) (DocumentationPage, error) {
	entity := DocumentationPage{
		CreatedAt:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
		VersionId:           data.VersionId,
		Slug:                data.Slug,
		PublishedRevisionId: data.PublishedRevisionId,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationPage{}, errors.Join(ErrDomainValidation, err)
	}

	return dp.queries.CreateDocumentationPage[DocumentationPage](ctx, queries.CreateDocumentationPageParams{
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
		VersionID:           entity.VersionId,
		Slug:                entity.Slug,
		PublishedRevisionID: entity.PublishedRevisionId,
	})
}

type UpdateDocumentationPageData struct {
	ID                  int64
	VersionId           int64
	Slug                string
	PublishedRevisionId pgtype.Int8
}

func (dp DocumentationPages) Update(ctx context.Context, data UpdateDocumentationPageData) (DocumentationPage, error) {
	entity := DocumentationPage{
		ID:                  data.ID,
		UpdatedAt:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
		VersionId:           data.VersionId,
		Slug:                data.Slug,
		PublishedRevisionId: data.PublishedRevisionId,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationPage{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := dp.queries.UpdateDocumentationPage[DocumentationPage](ctx, queries.UpdateDocumentationPageParams{
		ID:                  entity.ID,
		UpdatedAt:           entity.UpdatedAt,
		VersionID:           entity.VersionId,
		Slug:                entity.Slug,
		PublishedRevisionID: entity.PublishedRevisionId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationPage{}, ErrNotFound
		}
		return DocumentationPage{}, err
	}

	return row, nil
}

func (dp DocumentationPages) Destroy(ctx context.Context, id int64) error {
	if err := dp.queries.DeleteDocumentationPage(ctx, id); err != nil {
		return err
	}

	return nil
}

func (dp DocumentationPages) All(ctx context.Context) ([]DocumentationPage, error) {
	return dp.queries.ListDocumentationPages[DocumentationPage](ctx).All()
}

type PaginatedDocumentationPages struct {
	DocumentationPages []DocumentationPage
	TotalCount         int64
	Page               int64
	PageSize           int64
	TotalPages         int64
}

func (dp DocumentationPages) Paginate(ctx context.Context, page, pageSize int64) (PaginatedDocumentationPages, error) {
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

	totalCount, err := dp.queries.CountDocumentationPages(ctx)
	if err != nil {
		return PaginatedDocumentationPages{}, err
	}

	entities, err := dp.queries.ListDocumentationPages[DocumentationPage](ctx).
		Limit(int32(pageSize)).
		Offset(int32(offset)).
		All()
	if err != nil {
		return PaginatedDocumentationPages{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedDocumentationPages{
		DocumentationPages: entities,
		TotalCount:         int64(totalCount),
		Page:               page,
		PageSize:           pageSize,
		TotalPages:         totalPages,
	}, nil
}

func (dp DocumentationPages) Upsert(ctx context.Context, id int64, data CreateDocumentationPageData) (DocumentationPage, error) {
	entity := DocumentationPage{
		ID:                  id,
		CreatedAt:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
		VersionId:           data.VersionId,
		Slug:                data.Slug,
		PublishedRevisionId: data.PublishedRevisionId,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationPage{}, errors.Join(ErrDomainValidation, err)
	}

	return dp.queries.UpsertDocumentationPage[DocumentationPage](ctx, queries.UpsertDocumentationPageParams{
		ID:                  entity.ID,
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
		VersionID:           entity.VersionId,
		Slug:                entity.Slug,
		PublishedRevisionID: entity.PublishedRevisionId,
	})
}

func (dp DocumentationPages) FindByVersionSlug(ctx context.Context, versionID int64, slug string) (DocumentationPage, error) {
	entity, err := dp.queries.GetDocumentationPageByVersionSlug[DocumentationPage](ctx, queries.GetDocumentationPageByVersionSlugParams{
		VersionID: versionID,
		Slug:      slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationPage{}, ErrNotFound
		}
		return DocumentationPage{}, err
	}

	return entity, nil
}

func (dp DocumentationPages) ListByVersion(ctx context.Context, versionID int64) ([]DocumentationPage, error) {
	return dp.queries.ListDocumentationPagesByVersion[DocumentationPage](ctx, versionID)
}
