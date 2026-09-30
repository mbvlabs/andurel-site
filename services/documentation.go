package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"uuid"

	"andurel-site/docs"
	"andurel-site/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

var (
	ErrSlugTaken    = errors.New("slug already in use")
	ErrDraftMissing = errors.New("no draft to publish")
	ErrPageMismatch = errors.New("page does not belong to this version")
)

var documentationSlugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

type Documentation struct {
	db           storage.Connection
	versions     models.DocumentationVersions
	pages        models.DocumentationPages
	revisions    models.DocumentationRevisions
	navRevisions models.DocumentationNavRevisions
}

func NewDocumentation(
	db storage.Connection,
	versions models.DocumentationVersions,
	pages models.DocumentationPages,
	revisions models.DocumentationRevisions,
	navRevisions models.DocumentationNavRevisions,
) Documentation {
	return Documentation{
		db:           db,
		versions:     versions,
		pages:        pages,
		revisions:    revisions,
		navRevisions: navRevisions,
	}
}

type CreateDocumentationVersionInput struct {
	Slug     string
	Label    string
	IsLatest bool
	Position int32
}

type UpdateDocumentationVersionInput struct {
	ID       int64
	Slug     string
	Label    string
	IsLatest bool
	Position int32
}

type CreateDocumentationPageInput struct {
	VersionID int64
	Slug      string
	Title     string
	CreatedBy uuid.UUID
}

type SavePageDraftInput struct {
	PageID       int64
	Slug         string
	Title        string
	MetaTitle    string
	Description  string
	BodyMarkdown string
	CreatedBy    uuid.UUID
}

type SaveNavDraftInput struct {
	VersionID int64
	Tree      []models.DocumentationNavNode
	CreatedBy uuid.UUID
}

type DocumentationPageSummary struct {
	Page        models.DocumentationPage
	Title       string
	HasDraft    bool
	IsPublished bool
}

type DocumentationVersionWorkspace struct {
	Version  models.DocumentationVersion
	Pages    []DocumentationPageSummary
	NavDraft models.DocumentationNavRevision
}

type DocumentationPageEditor struct {
	Version   models.DocumentationVersion
	Page      models.DocumentationPage
	Draft     models.DocumentationRevision
	Published *models.DocumentationRevision
	Revisions []models.DocumentationRevision
}

func (d Documentation) ListVersions(ctx context.Context) ([]models.DocumentationVersion, error) {
	return d.versions.AllByPosition(ctx)
}

func (d Documentation) CreateVersion(
	ctx context.Context,
	createdBy uuid.UUID,
	input CreateDocumentationVersionInput,
) (models.DocumentationVersion, error) {
	if err := validateDocumentationSlug("slug", input.Slug); err != nil {
		return models.DocumentationVersion{}, err
	}
	if err := validateRequiredLabel(input.Label); err != nil {
		return models.DocumentationVersion{}, err
	}

	var version models.DocumentationVersion
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		navs := d.navRevisions.WithTx(tx)

		if err := ensureVersionSlugAvailable(ctx, versions, strings.ToLower(input.Slug), 0); err != nil {
			return err
		}
		if input.IsLatest {
			if err := versions.ClearLatest(ctx); err != nil {
				return err
			}
		}

		created, err := versions.Create(ctx, models.CreateDocumentationVersionData{
			Slug:           strings.ToLower(input.Slug),
			Label:          strings.TrimSpace(input.Label),
			IsLatest:       input.IsLatest,
			Position:       input.Position,
			PublishedNavId: pgtype.Int8{},
		})
		if err != nil {
			return err
		}

		if _, err := navs.Create(ctx, models.CreateDocumentationNavRevisionData{
			VersionId:   created.ID,
			Status:      models.DocumentationStatusDraft,
			Tree:        models.EmptyJSONArray,
			CreatedBy:   createdBy,
			PublishedAt: pgtype.Timestamptz{},
		}); err != nil {
			return err
		}

		version = created
		return nil
	})
	if err != nil {
		return models.DocumentationVersion{}, err
	}

	return version, nil
}

func (d Documentation) UpdateVersion(
	ctx context.Context,
	input UpdateDocumentationVersionInput,
) (models.DocumentationVersion, error) {
	if err := validateDocumentationSlug("slug", input.Slug); err != nil {
		return models.DocumentationVersion{}, err
	}
	if err := validateRequiredLabel(input.Label); err != nil {
		return models.DocumentationVersion{}, err
	}

	var version models.DocumentationVersion
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)

		existing, err := versions.Find(ctx, input.ID)
		if err != nil {
			return err
		}
		if err := ensureVersionSlugAvailable(ctx, versions, strings.ToLower(input.Slug), existing.ID); err != nil {
			return err
		}
		if input.IsLatest {
			if err := versions.ClearLatest(ctx); err != nil {
				return err
			}
		}

		updated, err := versions.Update(ctx, models.UpdateDocumentationVersionData{
			ID:             existing.ID,
			Slug:           strings.ToLower(input.Slug),
			Label:          strings.TrimSpace(input.Label),
			IsLatest:       input.IsLatest,
			Position:       input.Position,
			PublishedNavId: existing.PublishedNavId,
		})
		if err != nil {
			return err
		}

		version = updated
		return nil
	})
	if err != nil {
		return models.DocumentationVersion{}, err
	}

	return version, nil
}

func (d Documentation) VersionWorkspace(
	ctx context.Context,
	versionID int64,
	userID uuid.UUID,
) (DocumentationVersionWorkspace, error) {
	var workspace DocumentationVersionWorkspace
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)
		navs := d.navRevisions.WithTx(tx)

		version, err := versions.Find(ctx, versionID)
		if err != nil {
			return err
		}

		navDraft, err := forkNavDraft(ctx, navs, version, userID)
		if err != nil {
			return err
		}

		pageRows, err := pages.ListByVersion(ctx, version.ID)
		if err != nil {
			return err
		}

		summaries := make([]DocumentationPageSummary, 0, len(pageRows))
		for _, page := range pageRows {
			summary := DocumentationPageSummary{
				Page:        page,
				IsPublished: page.PublishedRevisionId.Valid,
			}

			draft, err := revisions.FindDraftByPage(ctx, page.ID)
			if err == nil {
				summary.HasDraft = true
				summary.Title = draft.Title
			} else if !errors.Is(err, models.ErrNotFound) {
				return err
			}

			if summary.Title == "" && page.PublishedRevisionId.Valid {
				published, err := revisions.Find(ctx, page.PublishedRevisionId.Int64)
				if err != nil && !errors.Is(err, models.ErrNotFound) {
					return err
				}
				if err == nil {
					summary.Title = published.Title
				}
			}

			if summary.Title == "" {
				summary.Title = page.Slug
			}

			summaries = append(summaries, summary)
		}

		workspace = DocumentationVersionWorkspace{
			Version:  version,
			Pages:    summaries,
			NavDraft: navDraft,
		}
		return nil
	})
	if err != nil {
		return DocumentationVersionWorkspace{}, err
	}

	return workspace, nil
}

func (d Documentation) CreatePage(
	ctx context.Context,
	input CreateDocumentationPageInput,
) (models.DocumentationPage, error) {
	if err := validateDocumentationSlug("slug", input.Slug); err != nil {
		return models.DocumentationPage{}, err
	}

	title := strings.TrimSpace(input.Title)
	b := validation.NewBuilder()
	b.Required("title", title)
	b.MaxLen("title", title, 255)
	if !b.Errors().Empty() {
		return models.DocumentationPage{}, b.Errors()
	}

	var page models.DocumentationPage
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)

		if _, err := versions.Find(ctx, input.VersionID); err != nil {
			return err
		}

		slug := strings.ToLower(input.Slug)
		if err := ensurePageSlugAvailable(ctx, pages, input.VersionID, slug, 0); err != nil {
			return err
		}

		created, err := pages.Create(ctx, models.CreateDocumentationPageData{
			VersionId:           input.VersionID,
			Slug:                slug,
			PublishedRevisionId: pgtype.Int8{},
		})
		if err != nil {
			return err
		}

		headings, err := headingsJSON("# " + title + "\n")
		if err != nil {
			return err
		}

		if _, err := revisions.Create(ctx, models.CreateDocumentationRevisionData{
			PageId:       created.ID,
			Status:       models.DocumentationStatusDraft,
			Title:        title,
			MetaTitle:    "",
			Description:  "",
			BodyMarkdown: "# " + title + "\n\n",
			Headings:     headings,
			CreatedBy:    input.CreatedBy,
			PublishedAt:  pgtype.Timestamptz{},
		}); err != nil {
			return err
		}

		page = created
		return nil
	})
	if err != nil {
		return models.DocumentationPage{}, err
	}

	return page, nil
}

func (d Documentation) SaveNavDraft(
	ctx context.Context,
	input SaveNavDraftInput,
) (models.DocumentationNavRevision, error) {
	treeJSON, err := models.MarshalDocumentationNavTree(input.Tree)
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}

	var draft models.DocumentationNavRevision
	err = storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		navs := d.navRevisions.WithTx(tx)

		version, err := versions.Find(ctx, input.VersionID)
		if err != nil {
			return err
		}

		current, err := forkNavDraft(ctx, navs, version, input.CreatedBy)
		if err != nil {
			return err
		}

		updated, err := navs.Update(ctx, models.UpdateDocumentationNavRevisionData{
			ID:          current.ID,
			VersionId:   version.ID,
			Status:      models.DocumentationStatusDraft,
			Tree:        treeJSON,
			CreatedBy:   input.CreatedBy,
			PublishedAt: pgtype.Timestamptz{},
		})
		if err != nil {
			return err
		}

		draft = updated
		return nil
	})
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}

	return draft, nil
}

func (d Documentation) PublishNav(
	ctx context.Context,
	versionID int64,
	userID uuid.UUID,
) (models.DocumentationNavRevision, error) {
	var published models.DocumentationNavRevision
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		navs := d.navRevisions.WithTx(tx)

		version, err := versions.Find(ctx, versionID)
		if err != nil {
			return err
		}

		draft, err := navs.FindDraftByVersion(ctx, version.ID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return ErrDraftMissing
			}
			return err
		}

		if version.PublishedNavId.Valid && version.PublishedNavId.Int64 != draft.ID {
			live, err := navs.Find(ctx, version.PublishedNavId.Int64)
			if err != nil && !errors.Is(err, models.ErrNotFound) {
				return err
			}
			if err == nil && live.Status == models.DocumentationStatusPublished {
				if _, err := navs.Update(ctx, models.UpdateDocumentationNavRevisionData{
					ID:          live.ID,
					VersionId:   live.VersionId,
					Status:      models.DocumentationStatusSuperseded,
					Tree:        live.Tree,
					CreatedBy:   live.CreatedBy,
					PublishedAt: live.PublishedAt,
				}); err != nil {
					return err
				}
			}
		}

		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		released, err := navs.Update(ctx, models.UpdateDocumentationNavRevisionData{
			ID:          draft.ID,
			VersionId:   version.ID,
			Status:      models.DocumentationStatusPublished,
			Tree:        models.NormalizeJSONArray(draft.Tree),
			CreatedBy:   userID,
			PublishedAt: now,
		})
		if err != nil {
			return err
		}

		if _, err := versions.Update(ctx, models.UpdateDocumentationVersionData{
			ID:             version.ID,
			Slug:           version.Slug,
			Label:          version.Label,
			IsLatest:       version.IsLatest,
			Position:       version.Position,
			PublishedNavId: pgtype.Int8{Int64: released.ID, Valid: true},
		}); err != nil {
			return err
		}

		published = released
		return nil
	})
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}

	return published, nil
}

func (d Documentation) PageEditor(
	ctx context.Context,
	pageID int64,
	userID uuid.UUID,
) (DocumentationPageEditor, error) {
	var editor DocumentationPageEditor
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)

		page, err := pages.Find(ctx, pageID)
		if err != nil {
			return err
		}

		version, err := versions.Find(ctx, page.VersionId)
		if err != nil {
			return err
		}

		draft, err := forkPageDraft(ctx, revisions, page, userID)
		if err != nil {
			return err
		}

		history, err := revisions.ListByPage(ctx, page.ID)
		if err != nil {
			return err
		}

		var published *models.DocumentationRevision
		if page.PublishedRevisionId.Valid {
			live, err := revisions.Find(ctx, page.PublishedRevisionId.Int64)
			if err != nil && !errors.Is(err, models.ErrNotFound) {
				return err
			}
			if err == nil {
				published = &live
			}
		}

		editor = DocumentationPageEditor{
			Version:   version,
			Page:      page,
			Draft:     draft,
			Published: published,
			Revisions: history,
		}
		return nil
	})
	if err != nil {
		return DocumentationPageEditor{}, err
	}

	return editor, nil
}

func (d Documentation) SavePageDraft(
	ctx context.Context,
	input SavePageDraftInput,
) (models.DocumentationRevision, error) {
	if err := validateDocumentationSlug("slug", input.Slug); err != nil {
		return models.DocumentationRevision{}, err
	}

	title := strings.TrimSpace(input.Title)
	metaTitle := strings.TrimSpace(input.MetaTitle)
	b := validation.NewBuilder()
	b.Required("title", title)
	b.MaxLen("title", title, 255)
	b.MaxLen("metaTitle", metaTitle, 255)
	if !b.Errors().Empty() {
		return models.DocumentationRevision{}, b.Errors()
	}

	headings, err := headingsJSON(input.BodyMarkdown)
	if err != nil {
		return models.DocumentationRevision{}, err
	}

	var draft models.DocumentationRevision
	err = storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)

		page, err := pages.Find(ctx, input.PageID)
		if err != nil {
			return err
		}

		slug := strings.ToLower(input.Slug)
		if err := ensurePageSlugAvailable(ctx, pages, page.VersionId, slug, page.ID); err != nil {
			return err
		}

		if _, err := pages.Update(ctx, models.UpdateDocumentationPageData{
			ID:                  page.ID,
			VersionId:           page.VersionId,
			Slug:                slug,
			PublishedRevisionId: page.PublishedRevisionId,
		}); err != nil {
			return err
		}

		current, err := forkPageDraft(ctx, revisions, page, input.CreatedBy)
		if err != nil {
			return err
		}

		updated, err := revisions.Update(ctx, models.UpdateDocumentationRevisionData{
			ID:           current.ID,
			PageId:       page.ID,
			Status:       models.DocumentationStatusDraft,
			Title:        title,
			MetaTitle:    metaTitle,
			Description:  strings.TrimSpace(input.Description),
			BodyMarkdown: input.BodyMarkdown,
			Headings:     headings,
			CreatedBy:    input.CreatedBy,
			PublishedAt:  pgtype.Timestamptz{},
		})
		if err != nil {
			return err
		}

		draft = updated
		return nil
	})
	if err != nil {
		return models.DocumentationRevision{}, err
	}

	return draft, nil
}

func (d Documentation) PublishPage(
	ctx context.Context,
	pageID int64,
	userID uuid.UUID,
) (models.DocumentationRevision, error) {
	var published models.DocumentationRevision
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)

		page, err := pages.Find(ctx, pageID)
		if err != nil {
			return err
		}

		draft, err := revisions.FindDraftByPage(ctx, page.ID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return ErrDraftMissing
			}
			return err
		}

		headings, err := headingsJSON(draft.BodyMarkdown)
		if err != nil {
			return err
		}

		if page.PublishedRevisionId.Valid && page.PublishedRevisionId.Int64 != draft.ID {
			live, err := revisions.Find(ctx, page.PublishedRevisionId.Int64)
			if err != nil && !errors.Is(err, models.ErrNotFound) {
				return err
			}
			if err == nil && live.Status == models.DocumentationStatusPublished {
				if _, err := revisions.Update(ctx, models.UpdateDocumentationRevisionData{
					ID:           live.ID,
					PageId:       live.PageId,
					Status:       models.DocumentationStatusSuperseded,
					Title:        live.Title,
					MetaTitle:    live.MetaTitle,
					Description:  live.Description,
					BodyMarkdown: live.BodyMarkdown,
					Headings:     live.Headings,
					CreatedBy:    live.CreatedBy,
					PublishedAt:  live.PublishedAt,
				}); err != nil {
					return err
				}
			}
		}

		now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
		released, err := revisions.Update(ctx, models.UpdateDocumentationRevisionData{
			ID:           draft.ID,
			PageId:       page.ID,
			Status:       models.DocumentationStatusPublished,
			Title:        draft.Title,
			MetaTitle:    draft.MetaTitle,
			Description:  draft.Description,
			BodyMarkdown: draft.BodyMarkdown,
			Headings:     headings,
			CreatedBy:    userID,
			PublishedAt:  now,
		})
		if err != nil {
			return err
		}

		if _, err := pages.Update(ctx, models.UpdateDocumentationPageData{
			ID:                  page.ID,
			VersionId:           page.VersionId,
			Slug:                page.Slug,
			PublishedRevisionId: pgtype.Int8{Int64: released.ID, Valid: true},
		}); err != nil {
			return err
		}

		published = released
		return nil
	})
	if err != nil {
		return models.DocumentationRevision{}, err
	}

	return published, nil
}

func (d Documentation) RestoreRevision(
	ctx context.Context,
	pageID, revisionID int64,
	userID uuid.UUID,
) (models.DocumentationRevision, error) {
	var draft models.DocumentationRevision
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)

		page, err := pages.Find(ctx, pageID)
		if err != nil {
			return err
		}

		source, err := revisions.Find(ctx, revisionID)
		if err != nil {
			return err
		}
		if source.PageId != page.ID {
			return ErrPageMismatch
		}

		current, err := forkPageDraft(ctx, revisions, page, userID)
		if err != nil {
			return err
		}

		updated, err := revisions.Update(ctx, models.UpdateDocumentationRevisionData{
			ID:           current.ID,
			PageId:       page.ID,
			Status:       models.DocumentationStatusDraft,
			Title:        source.Title,
			MetaTitle:    source.MetaTitle,
			Description:  source.Description,
			BodyMarkdown: source.BodyMarkdown,
			Headings:     models.NormalizeJSONArray(source.Headings),
			CreatedBy:    userID,
			PublishedAt:  pgtype.Timestamptz{},
		})
		if err != nil {
			return err
		}

		draft = updated
		return nil
	})
	if err != nil {
		return models.DocumentationRevision{}, err
	}

	return draft, nil
}

func (d Documentation) RenderPreview(ctx context.Context, pageID int64, source string) (string, []docs.Heading, error) {
	if _, err := d.pages.Find(ctx, pageID); err != nil {
		return "", nil, err
	}

	html, err := docs.RenderHTML([]byte(source))
	if err != nil {
		return "", nil, err
	}

	headings, err := docs.ExtractHeadings([]byte(source))
	if err != nil {
		return "", nil, err
	}
	if headings == nil {
		headings = []docs.Heading{}
	}

	return html, headings, nil
}

func forkPageDraft(
	ctx context.Context,
	revisions models.DocumentationRevisions,
	page models.DocumentationPage,
	userID uuid.UUID,
) (models.DocumentationRevision, error) {
	draft, err := revisions.FindDraftByPage(ctx, page.ID)
	if err == nil {
		return draft, nil
	}
	if !errors.Is(err, models.ErrNotFound) {
		return models.DocumentationRevision{}, err
	}

	title := page.Slug
	metaTitle := ""
	description := ""
	body := "# " + page.Slug + "\n\n"
	headings := models.EmptyJSONArray

	if page.PublishedRevisionId.Valid {
		live, err := revisions.Find(ctx, page.PublishedRevisionId.Int64)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return models.DocumentationRevision{}, err
		}
		if err == nil {
			title = live.Title
			metaTitle = live.MetaTitle
			description = live.Description
			body = live.BodyMarkdown
			headings = models.NormalizeJSONArray(live.Headings)
		}
	}

	created, err := headingsJSON(body)
	if err == nil {
		headings = created
	}

	return revisions.Create(ctx, models.CreateDocumentationRevisionData{
		PageId:       page.ID,
		Status:       models.DocumentationStatusDraft,
		Title:        title,
		MetaTitle:    metaTitle,
		Description:  description,
		BodyMarkdown: body,
		Headings:     headings,
		CreatedBy:    userID,
		PublishedAt:  pgtype.Timestamptz{},
	})
}

func forkNavDraft(
	ctx context.Context,
	navs models.DocumentationNavRevisions,
	version models.DocumentationVersion,
	userID uuid.UUID,
) (models.DocumentationNavRevision, error) {
	draft, err := navs.FindDraftByVersion(ctx, version.ID)
	if err == nil {
		return draft, nil
	}
	if !errors.Is(err, models.ErrNotFound) {
		return models.DocumentationNavRevision{}, err
	}

	tree := models.EmptyJSONArray
	if version.PublishedNavId.Valid {
		live, err := navs.Find(ctx, version.PublishedNavId.Int64)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return models.DocumentationNavRevision{}, err
		}
		if err == nil {
			tree = models.NormalizeJSONArray(live.Tree)
		}
	}

	return navs.Create(ctx, models.CreateDocumentationNavRevisionData{
		VersionId:   version.ID,
		Status:      models.DocumentationStatusDraft,
		Tree:        tree,
		CreatedBy:   userID,
		PublishedAt: pgtype.Timestamptz{},
	})
}

func headingsJSON(source string) ([]byte, error) {
	headings, err := docs.ExtractHeadings([]byte(source))
	if err != nil {
		return nil, fmt.Errorf("extract headings: %w", err)
	}
	if headings == nil {
		headings = []docs.Heading{}
	}

	encoded, err := json.Marshal(headings)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func validateDocumentationSlug(field, slug string) error {
	b := validation.NewBuilder()
	trimmed := strings.TrimSpace(slug)
	b.Required(field, trimmed)
	b.MaxLen(field, trimmed, 64)
	if trimmed != "" && !documentationSlugPattern.MatchString(strings.ToLower(trimmed)) {
		b.Add(field, "format", "must start with a letter or number and use lowercase letters, numbers, dots, underscores, or hyphens")
	}
	if !b.Errors().Empty() {
		return b.Errors()
	}

	return nil
}

func validateRequiredLabel(label string) error {
	b := validation.NewBuilder()
	trimmed := strings.TrimSpace(label)
	b.Required("label", trimmed)
	b.MaxLen("label", trimmed, 255)
	if !b.Errors().Empty() {
		return b.Errors()
	}

	return nil
}

func ensureVersionSlugAvailable(
	ctx context.Context,
	versions models.DocumentationVersions,
	slug string,
	ignoreID int64,
) error {
	existing, err := versions.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.ID == ignoreID {
		return nil
	}

	return ErrSlugTaken
}

func ensurePageSlugAvailable(
	ctx context.Context,
	pages models.DocumentationPages,
	versionID int64,
	slug string,
	ignoreID int64,
) error {
	existing, err := pages.FindByVersionSlug(ctx, versionID, slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.ID == ignoreID {
		return nil
	}

	return ErrSlugTaken
}
