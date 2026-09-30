package models

// andurel:table documentation_revisions

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

type DocumentationRevisions struct {
	queries *queries.Queries
}

func NewDocumentationRevisions(db storage.Connection) DocumentationRevisions {
	return DocumentationRevisions{queries: queries.New(db)}
}

// WithTx returns a copy that runs queries inside tx.
func (dr DocumentationRevisions) WithTx(tx storage.Transaction) DocumentationRevisions {
	return DocumentationRevisions{queries: queries.New(tx)}
}

type DocumentationRevision struct {
	ID           int64              `andurel:"id"`
	CreatedAt    pgtype.Timestamptz `andurel:"created_at"`
	PageId       int64              `andurel:"page_id"`
	Status       string             `andurel:"status"`
	Title        string             `andurel:"title"`
	MetaTitle    string             `andurel:"meta_title"`
	Description  string             `andurel:"description"`
	BodyMarkdown string             `andurel:"body_markdown"`
	Headings     []byte             `andurel:"headings"`
	CreatedBy    uuid.UUID          `andurel:"created_by"`
	PublishedAt  pgtype.Timestamptz `andurel:"published_at"`
}

func (dr DocumentationRevisions) Find(ctx context.Context, id int64) (DocumentationRevision, error) {
	entity, err := dr.queries.GetDocumentationRevision[DocumentationRevision](ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationRevision{}, ErrNotFound
		}
		return DocumentationRevision{}, err
	}

	return entity, nil
}

type CreateDocumentationRevisionData struct {
	PageId       int64
	Status       string
	Title        string
	MetaTitle    string
	Description  string
	BodyMarkdown string
	Headings     []byte
	CreatedBy    uuid.UUID
	PublishedAt  pgtype.Timestamptz
}

func (dr DocumentationRevisions) Create(ctx context.Context, data CreateDocumentationRevisionData) (DocumentationRevision, error) {
	entity := DocumentationRevision{
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		PageId:       data.PageId,
		Status:       data.Status,
		Title:        data.Title,
		MetaTitle:    data.MetaTitle,
		Description:  data.Description,
		BodyMarkdown: data.BodyMarkdown,
		Headings:     data.Headings,
		CreatedBy:    data.CreatedBy,
		PublishedAt:  data.PublishedAt,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationRevision{}, errors.Join(ErrDomainValidation, err)
	}

	return dr.queries.CreateDocumentationRevision[DocumentationRevision](ctx, queries.CreateDocumentationRevisionParams{
		CreatedAt:    entity.CreatedAt,
		PageID:       entity.PageId,
		Status:       entity.Status,
		Title:        entity.Title,
		MetaTitle:    entity.MetaTitle,
		Description:  entity.Description,
		BodyMarkdown: entity.BodyMarkdown,
		Headings:     entity.Headings,
		CreatedBy:    entity.CreatedBy,
		PublishedAt:  entity.PublishedAt,
	})
}

type UpdateDocumentationRevisionData struct {
	ID           int64
	PageId       int64
	Status       string
	Title        string
	MetaTitle    string
	Description  string
	BodyMarkdown string
	Headings     []byte
	CreatedBy    uuid.UUID
	PublishedAt  pgtype.Timestamptz
}

func (dr DocumentationRevisions) Update(ctx context.Context, data UpdateDocumentationRevisionData) (DocumentationRevision, error) {
	entity := DocumentationRevision{
		ID:           data.ID,
		PageId:       data.PageId,
		Status:       data.Status,
		Title:        data.Title,
		MetaTitle:    data.MetaTitle,
		Description:  data.Description,
		BodyMarkdown: data.BodyMarkdown,
		Headings:     data.Headings,
		CreatedBy:    data.CreatedBy,
		PublishedAt:  data.PublishedAt,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationRevision{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := dr.queries.UpdateDocumentationRevision[DocumentationRevision](ctx, queries.UpdateDocumentationRevisionParams{
		ID:           entity.ID,
		PageID:       entity.PageId,
		Status:       entity.Status,
		Title:        entity.Title,
		MetaTitle:    entity.MetaTitle,
		Description:  entity.Description,
		BodyMarkdown: entity.BodyMarkdown,
		Headings:     entity.Headings,
		CreatedBy:    entity.CreatedBy,
		PublishedAt:  entity.PublishedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationRevision{}, ErrNotFound
		}
		return DocumentationRevision{}, err
	}

	return row, nil
}

func (dr DocumentationRevisions) Destroy(ctx context.Context, id int64) error {
	if err := dr.queries.DeleteDocumentationRevision(ctx, id); err != nil {
		return err
	}

	return nil
}

func (dr DocumentationRevisions) All(ctx context.Context) ([]DocumentationRevision, error) {
	return dr.queries.ListDocumentationRevisions[DocumentationRevision](ctx).All()
}

type PaginatedDocumentationRevisions struct {
	DocumentationRevisions []DocumentationRevision
	TotalCount             int64
	Page                   int64
	PageSize               int64
	TotalPages             int64
}

func (dr DocumentationRevisions) Paginate(ctx context.Context, page, pageSize int64) (PaginatedDocumentationRevisions, error) {
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

	totalCount, err := dr.queries.CountDocumentationRevisions(ctx)
	if err != nil {
		return PaginatedDocumentationRevisions{}, err
	}

	entities, err := dr.queries.ListDocumentationRevisions[DocumentationRevision](ctx).
		Limit(int32(pageSize)).
		Offset(int32(offset)).
		All()
	if err != nil {
		return PaginatedDocumentationRevisions{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedDocumentationRevisions{
		DocumentationRevisions: entities,
		TotalCount:             int64(totalCount),
		Page:                   page,
		PageSize:               pageSize,
		TotalPages:             totalPages,
	}, nil
}

func (dr DocumentationRevisions) Upsert(ctx context.Context, id int64, data CreateDocumentationRevisionData) (DocumentationRevision, error) {
	entity := DocumentationRevision{
		ID:           id,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		PageId:       data.PageId,
		Status:       data.Status,
		Title:        data.Title,
		MetaTitle:    data.MetaTitle,
		Description:  data.Description,
		BodyMarkdown: data.BodyMarkdown,
		Headings:     data.Headings,
		CreatedBy:    data.CreatedBy,
		PublishedAt:  data.PublishedAt,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationRevision{}, errors.Join(ErrDomainValidation, err)
	}

	return dr.queries.UpsertDocumentationRevision[DocumentationRevision](ctx, queries.UpsertDocumentationRevisionParams{
		ID:           entity.ID,
		CreatedAt:    entity.CreatedAt,
		PageID:       entity.PageId,
		Status:       entity.Status,
		Title:        entity.Title,
		MetaTitle:    entity.MetaTitle,
		Description:  entity.Description,
		BodyMarkdown: entity.BodyMarkdown,
		Headings:     entity.Headings,
		CreatedBy:    entity.CreatedBy,
		PublishedAt:  entity.PublishedAt,
	})
}

func (dr DocumentationRevisions) FindDraftByPage(ctx context.Context, pageID int64) (DocumentationRevision, error) {
	entity, err := dr.queries.GetDocumentationRevisionDraftByPage[DocumentationRevision](ctx, pageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationRevision{}, ErrNotFound
		}
		return DocumentationRevision{}, err
	}

	return entity, nil
}

func (dr DocumentationRevisions) ListByPage(ctx context.Context, pageID int64) ([]DocumentationRevision, error) {
	return dr.queries.ListDocumentationRevisionsByPage[DocumentationRevision](ctx, pageID)
}
