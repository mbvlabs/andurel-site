package services

import (
	"context"
	"errors"
	"strings"
	"time"
	"uuid"

	"andurel-site/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

var ErrNavPageMissing = errors.New("nav references an unknown page slug")

type DocumentationPageSnapshot struct {
	Slug         string
	Title        string
	MetaTitle    string
	Description  string
	BodyMarkdown string
}

type DocumentationNavSnapshotNode struct {
	Title    string
	Slug     string
	Children []DocumentationNavSnapshotNode
}

type DocumentationSnapshot struct {
	Slug     string
	Label    string
	IsLatest bool
	Position int32
	Pages    []DocumentationPageSnapshot
	Nav      []DocumentationNavSnapshotNode
}

type DocumentationListItem struct {
	Slug      string
	Label     string
	IsLatest  bool
	Position  int32
	PageCount int
}

type ReplaceDocumentationInput struct {
	Slug     string
	Label    string
	IsLatest bool
	Position int32
	Publish  bool
	Pages    []DocumentationPageSnapshot
	Nav      []DocumentationNavSnapshotNode
}

func (d Documentation) ListDocumentations(ctx context.Context) ([]DocumentationListItem, error) {
	versions, err := d.versions.AllByPosition(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]DocumentationListItem, 0, len(versions))
	for _, version := range versions {
		pages, err := d.pages.ListByVersion(ctx, version.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, DocumentationListItem{
			Slug:      version.Slug,
			Label:     version.Label,
			IsLatest:  version.IsLatest,
			Position:  version.Position,
			PageCount: len(pages),
		})
	}

	return items, nil
}

func (d Documentation) GetDocumentation(ctx context.Context, slug string) (DocumentationSnapshot, error) {
	version, err := d.versions.FindBySlug(ctx, strings.ToLower(strings.TrimSpace(slug)))
	if err != nil {
		return DocumentationSnapshot{}, err
	}

	return d.snapshotForVersion(ctx, version)
}

func (d Documentation) CreateDocumentation(
	ctx context.Context,
	createdBy uuid.UUID,
	input ReplaceDocumentationInput,
) (DocumentationSnapshot, error) {
	return d.replaceDocumentation(ctx, createdBy, "", input)
}

func (d Documentation) UpdateDocumentation(
	ctx context.Context,
	createdBy uuid.UUID,
	slug string,
	input ReplaceDocumentationInput,
) (DocumentationSnapshot, error) {
	return d.replaceDocumentation(ctx, createdBy, slug, input)
}

func (d Documentation) replaceDocumentation(
	ctx context.Context,
	createdBy uuid.UUID,
	existingSlug string,
	input ReplaceDocumentationInput,
) (DocumentationSnapshot, error) {
	if err := validateDocumentationSlug("slug", input.Slug); err != nil {
		return DocumentationSnapshot{}, err
	}
	if err := validateRequiredLabel(input.Label); err != nil {
		return DocumentationSnapshot{}, err
	}
	if err := validateImportedPages(input.Pages); err != nil {
		return DocumentationSnapshot{}, err
	}

	var snapshot DocumentationSnapshot
	err := storage.RunInTransaction(ctx, d.db, func(ctx context.Context, tx storage.Transaction) error {
		versions := d.versions.WithTx(tx)
		pages := d.pages.WithTx(tx)
		revisions := d.revisions.WithTx(tx)
		navs := d.navRevisions.WithTx(tx)

		var version models.DocumentationVersion
		var err error
		if existingSlug == "" {
			if err := ensureVersionSlugAvailable(ctx, versions, strings.ToLower(input.Slug), 0); err != nil {
				return err
			}
			version, err = createImportedVersion(ctx, versions, navs, createdBy, input)
			if err != nil {
				return err
			}
		} else {
			version, err = versions.FindBySlug(ctx, strings.ToLower(strings.TrimSpace(existingSlug)))
			if err != nil {
				return err
			}
			if err := ensureVersionSlugAvailable(ctx, versions, strings.ToLower(input.Slug), version.ID); err != nil {
				return err
			}
			if input.IsLatest {
				if err := versions.ClearLatest(ctx); err != nil {
					return err
				}
			}
			version, err = versions.Update(ctx, models.UpdateDocumentationVersionData{
				ID:             version.ID,
				Slug:           strings.ToLower(input.Slug),
				Label:          strings.TrimSpace(input.Label),
				IsLatest:       input.IsLatest,
				Position:       input.Position,
				PublishedNavId: version.PublishedNavId,
			})
			if err != nil {
				return err
			}
		}

		for _, page := range input.Pages {
			if _, err := upsertImportedPage(ctx, pages, revisions, version.ID, createdBy, page, input.Publish); err != nil {
				return err
			}
		}

		pageRows, err := pages.ListByVersion(ctx, version.ID)
		if err != nil {
			return err
		}
		slugs := make(map[string]int64, len(pageRows))
		for _, page := range pageRows {
			slugs[page.Slug] = page.ID
		}

		tree, err := navTreeFromSnapshot(input.Nav, slugs)
		if err != nil {
			return err
		}

		if _, err := saveImportedNav(ctx, versions, navs, version, createdBy, tree, input.Publish); err != nil {
			return err
		}

		version, err = versions.Find(ctx, version.ID)
		if err != nil {
			return err
		}

		snapshot, err = documentationSnapshot(ctx, pages, revisions, navs, version)
		return err
	})
	if err != nil {
		return DocumentationSnapshot{}, err
	}

	return snapshot, nil
}

func (d Documentation) snapshotForVersion(ctx context.Context, version models.DocumentationVersion) (DocumentationSnapshot, error) {
	return documentationSnapshot(ctx, d.pages, d.revisions, d.navRevisions, version)
}

func createImportedVersion(
	ctx context.Context,
	versions models.DocumentationVersions,
	navs models.DocumentationNavRevisions,
	createdBy uuid.UUID,
	input ReplaceDocumentationInput,
) (models.DocumentationVersion, error) {
	if input.IsLatest {
		if err := versions.ClearLatest(ctx); err != nil {
			return models.DocumentationVersion{}, err
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
		return models.DocumentationVersion{}, err
	}

	if _, err := navs.Create(ctx, models.CreateDocumentationNavRevisionData{
		VersionId:   created.ID,
		Status:      models.DocumentationStatusDraft,
		Tree:        models.EmptyJSONArray,
		CreatedBy:   createdBy,
		PublishedAt: pgtype.Timestamptz{},
	}); err != nil {
		return models.DocumentationVersion{}, err
	}

	return created, nil
}

func upsertImportedPage(
	ctx context.Context,
	pages models.DocumentationPages,
	revisions models.DocumentationRevisions,
	versionID int64,
	createdBy uuid.UUID,
	input DocumentationPageSnapshot,
	publish bool,
) (models.DocumentationPage, error) {
	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	page, err := pages.FindByVersionSlug(ctx, versionID, slug)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return models.DocumentationPage{}, err
		}
		page, err = pages.Create(ctx, models.CreateDocumentationPageData{
			VersionId:           versionID,
			Slug:                slug,
			PublishedRevisionId: pgtype.Int8{},
		})
		if err != nil {
			return models.DocumentationPage{}, err
		}
	}

	headings, err := headingsJSON(input.BodyMarkdown)
	if err != nil {
		return models.DocumentationPage{}, err
	}

	draft, err := forkPageDraft(ctx, revisions, page, createdBy)
	if err != nil {
		return models.DocumentationPage{}, err
	}

	title := strings.TrimSpace(input.Title)
	updated, err := revisions.Update(ctx, models.UpdateDocumentationRevisionData{
		ID:           draft.ID,
		PageId:       page.ID,
		Status:       models.DocumentationStatusDraft,
		Title:        title,
		MetaTitle:    strings.TrimSpace(input.MetaTitle),
		Description:  strings.TrimSpace(input.Description),
		BodyMarkdown: input.BodyMarkdown,
		Headings:     headings,
		CreatedBy:    createdBy,
		PublishedAt:  pgtype.Timestamptz{},
	})
	if err != nil {
		return models.DocumentationPage{}, err
	}

	if !publish {
		return page, nil
	}

	if page.PublishedRevisionId.Valid && page.PublishedRevisionId.Int64 != updated.ID {
		live, err := revisions.Find(ctx, page.PublishedRevisionId.Int64)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return models.DocumentationPage{}, err
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
				return models.DocumentationPage{}, err
			}
		}
	}

	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	released, err := revisions.Update(ctx, models.UpdateDocumentationRevisionData{
		ID:           updated.ID,
		PageId:       page.ID,
		Status:       models.DocumentationStatusPublished,
		Title:        updated.Title,
		MetaTitle:    updated.MetaTitle,
		Description:  updated.Description,
		BodyMarkdown: updated.BodyMarkdown,
		Headings:     headings,
		CreatedBy:    createdBy,
		PublishedAt:  now,
	})
	if err != nil {
		return models.DocumentationPage{}, err
	}

	return pages.Update(ctx, models.UpdateDocumentationPageData{
		ID:                  page.ID,
		VersionId:           page.VersionId,
		Slug:                page.Slug,
		PublishedRevisionId: pgtype.Int8{Int64: released.ID, Valid: true},
	})
}

func saveImportedNav(
	ctx context.Context,
	versions models.DocumentationVersions,
	navs models.DocumentationNavRevisions,
	version models.DocumentationVersion,
	createdBy uuid.UUID,
	tree []models.DocumentationNavNode,
	publish bool,
) (models.DocumentationNavRevision, error) {
	treeJSON, err := models.MarshalDocumentationNavTree(tree)
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}

	draft, err := forkNavDraft(ctx, navs, version, createdBy)
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}

	updated, err := navs.Update(ctx, models.UpdateDocumentationNavRevisionData{
		ID:          draft.ID,
		VersionId:   version.ID,
		Status:      models.DocumentationStatusDraft,
		Tree:        treeJSON,
		CreatedBy:   createdBy,
		PublishedAt: pgtype.Timestamptz{},
	})
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}
	if !publish {
		return updated, nil
	}

	if version.PublishedNavId.Valid && version.PublishedNavId.Int64 != updated.ID {
		live, err := navs.Find(ctx, version.PublishedNavId.Int64)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return models.DocumentationNavRevision{}, err
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
				return models.DocumentationNavRevision{}, err
			}
		}
	}

	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	released, err := navs.Update(ctx, models.UpdateDocumentationNavRevisionData{
		ID:          updated.ID,
		VersionId:   version.ID,
		Status:      models.DocumentationStatusPublished,
		Tree:        models.NormalizeJSONArray(updated.Tree),
		CreatedBy:   createdBy,
		PublishedAt: now,
	})
	if err != nil {
		return models.DocumentationNavRevision{}, err
	}

	if _, err := versions.Update(ctx, models.UpdateDocumentationVersionData{
		ID:             version.ID,
		Slug:           version.Slug,
		Label:          version.Label,
		IsLatest:       version.IsLatest,
		Position:       version.Position,
		PublishedNavId: pgtype.Int8{Int64: released.ID, Valid: true},
	}); err != nil {
		return models.DocumentationNavRevision{}, err
	}

	return released, nil
}

func documentationSnapshot(
	ctx context.Context,
	pages models.DocumentationPages,
	revisions models.DocumentationRevisions,
	navs models.DocumentationNavRevisions,
	version models.DocumentationVersion,
) (DocumentationSnapshot, error) {
	pageRows, err := pages.ListByVersion(ctx, version.ID)
	if err != nil {
		return DocumentationSnapshot{}, err
	}

	idsToSlug := make(map[int64]string, len(pageRows))
	exportedPages := make([]DocumentationPageSnapshot, 0, len(pageRows))
	for _, page := range pageRows {
		idsToSlug[page.ID] = page.Slug
		exported, err := pageSnapshot(ctx, revisions, page)
		if err != nil {
			return DocumentationSnapshot{}, err
		}
		exportedPages = append(exportedPages, exported)
	}

	navRevision, err := currentNavRevision(ctx, navs, version)
	if err != nil {
		return DocumentationSnapshot{}, err
	}
	tree, err := models.ParseDocumentationNavTree(navRevision.Tree)
	if err != nil {
		return DocumentationSnapshot{}, err
	}

	return DocumentationSnapshot{
		Slug:     version.Slug,
		Label:    version.Label,
		IsLatest: version.IsLatest,
		Position: version.Position,
		Pages:    exportedPages,
		Nav:      navSnapshotFromTree(tree, idsToSlug),
	}, nil
}

func pageSnapshot(
	ctx context.Context,
	revisions models.DocumentationRevisions,
	page models.DocumentationPage,
) (DocumentationPageSnapshot, error) {
	var revision models.DocumentationRevision
	var err error
	if page.PublishedRevisionId.Valid {
		revision, err = revisions.Find(ctx, page.PublishedRevisionId.Int64)
		if err != nil && !errors.Is(err, models.ErrNotFound) {
			return DocumentationPageSnapshot{}, err
		}
	}
	if revision.ID == 0 {
		revision, err = revisions.FindDraftByPage(ctx, page.ID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return DocumentationPageSnapshot{
					Slug:  page.Slug,
					Title: page.Slug,
				}, nil
			}
			return DocumentationPageSnapshot{}, err
		}
	}

	return DocumentationPageSnapshot{
		Slug:         page.Slug,
		Title:        revision.Title,
		MetaTitle:    revision.MetaTitle,
		Description:  revision.Description,
		BodyMarkdown: revision.BodyMarkdown,
	}, nil
}

func currentNavRevision(
	ctx context.Context,
	navs models.DocumentationNavRevisions,
	version models.DocumentationVersion,
) (models.DocumentationNavRevision, error) {
	if version.PublishedNavId.Valid {
		live, err := navs.Find(ctx, version.PublishedNavId.Int64)
		if err == nil {
			return live, nil
		}
		if !errors.Is(err, models.ErrNotFound) {
			return models.DocumentationNavRevision{}, err
		}
	}

	draft, err := navs.FindDraftByVersion(ctx, version.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.DocumentationNavRevision{Tree: models.EmptyJSONArray}, nil
		}
		return models.DocumentationNavRevision{}, err
	}

	return draft, nil
}

func validateImportedPages(pages []DocumentationPageSnapshot) error {
	b := validation.NewBuilder()
	seen := make(map[string]struct{}, len(pages))
	for _, page := range pages {
		field := "pages"
		if err := validateDocumentationSlug("slug", page.Slug); err != nil {
			b.Add(field, "slug", err.Error())
		}
		title := strings.TrimSpace(page.Title)
		if title == "" {
			b.Add(field, "title", "is required")
		}
		if len(title) > 255 {
			b.Add(field, "title", "is too long")
		}
		if len(strings.TrimSpace(page.MetaTitle)) > 255 {
			b.Add(field, "metaTitle", "is too long")
		}
		slug := strings.ToLower(strings.TrimSpace(page.Slug))
		if slug != "" {
			if _, exists := seen[slug]; exists {
				b.Add(field, "slug", "must be unique in this documentation")
			}
			seen[slug] = struct{}{}
		}
	}
	if !b.Errors().Empty() {
		return b.Errors()
	}

	return nil
}

func navTreeFromSnapshot(nodes []DocumentationNavSnapshotNode, slugs map[string]int64) ([]models.DocumentationNavNode, error) {
	tree := make([]models.DocumentationNavNode, 0, len(nodes))
	for _, node := range nodes {
		item := models.DocumentationNavNode{
			Title: strings.TrimSpace(node.Title),
		}
		if item.Title == "" {
			b := validation.NewBuilder()
			b.Required("nav.title", item.Title)
			return nil, b.Errors()
		}
		if slug := strings.ToLower(strings.TrimSpace(node.Slug)); slug != "" {
			id, ok := slugs[slug]
			if !ok {
				return nil, ErrNavPageMissing
			}
			item.PageID = &id
		}
		children, err := navTreeFromSnapshot(node.Children, slugs)
		if err != nil {
			return nil, err
		}
		item.Children = children
		tree = append(tree, item)
	}

	return tree, nil
}

func navSnapshotFromTree(nodes []models.DocumentationNavNode, ids map[int64]string) []DocumentationNavSnapshotNode {
	out := make([]DocumentationNavSnapshotNode, 0, len(nodes))
	for _, node := range nodes {
		item := DocumentationNavSnapshotNode{
			Title:    node.Title,
			Children: navSnapshotFromTree(node.Children, ids),
		}
		if node.PageID != nil {
			item.Slug = ids[*node.PageID]
		}
		out = append(out, item)
	}

	return out
}
