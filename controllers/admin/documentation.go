package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uuid"

	"andurel-site/models"
	"andurel-site/router/cookies"
	"andurel-site/services"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type DocumentationVersionData struct {
	ID             int64  `json:"id"`
	Slug           string `json:"slug"`
	Label          string `json:"label"`
	IsLatest       bool   `json:"isLatest"`
	Position       int32  `json:"position"`
	PublishedNavID *int64 `json:"publishedNavId"`
	UpdatedAt      string `json:"updatedAt"`
}

type DocumentationPageSummaryData struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	HasDraft    bool   `json:"hasDraft"`
	IsPublished bool   `json:"isPublished"`
}

type DocumentationNavNodeData struct {
	Title    string                     `json:"title"`
	PageID   *int64                     `json:"pageId"`
	Children []DocumentationNavNodeData `json:"children"`
}

type DocumentationNavDraftData struct {
	ID     int64                      `json:"id"`
	Status string                     `json:"status"`
	Tree   []DocumentationNavNodeData `json:"tree"`
}

type DocumentationHeadingData struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Level int    `json:"level"`
}

type DocumentationRevisionData struct {
	ID           int64                      `json:"id"`
	Status       string                     `json:"status"`
	Title        string                     `json:"title"`
	MetaTitle    string                     `json:"metaTitle"`
	Description  string                     `json:"description"`
	BodyMarkdown string                     `json:"bodyMarkdown"`
	Headings     []DocumentationHeadingData `json:"headings"`
	CreatedAt    string                     `json:"createdAt"`
	PublishedAt  string                     `json:"publishedAt"`
}

type DocumentationPageData struct {
	ID                  int64  `json:"id"`
	VersionID           int64  `json:"versionId"`
	Slug                string `json:"slug"`
	PublishedRevisionID *int64 `json:"publishedRevisionId"`
}

type DocumentationVersionIndexProps struct {
	Versions []DocumentationVersionData `json:"versions"`
}

type DocumentationVersionItemProps struct {
	Version  DocumentationVersionData       `json:"version"`
	Pages    []DocumentationPageSummaryData `json:"pages"`
	NavDraft DocumentationNavDraftData      `json:"navDraft"`
}

type DocumentationPageItemProps struct {
	Version   DocumentationVersionData    `json:"version"`
	Page      DocumentationPageData       `json:"page"`
	Draft     DocumentationRevisionData   `json:"draft"`
	Published *DocumentationRevisionData  `json:"published"`
	Revisions []DocumentationRevisionData `json:"revisions"`
}

type CreateDocumentationVersionFormPayload struct {
	Slug     string `json:"slug"`
	Label    string `json:"label"`
	IsLatest bool   `json:"isLatest"`
}

type UpdateDocumentationVersionFormPayload struct {
	Slug     string `json:"slug"`
	Label    string `json:"label"`
	IsLatest bool   `json:"isLatest"`
}

type ReorderDocumentationVersionsFormPayload struct {
	IDs []int64 `json:"ids"`
}

type CreateDocumentationPageFormPayload struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type UpdateDocumentationPageFormPayload struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	MetaTitle    string `json:"metaTitle"`
	Description  string `json:"description"`
	BodyMarkdown string `json:"bodyMarkdown"`
}

type PreviewDocumentationPageFormPayload struct {
	BodyMarkdown string `json:"bodyMarkdown"`
}

type UpdateDocumentationNavFormPayload struct {
	Tree []DocumentationNavNodeData `json:"tree"`
}

func currentAdminUserID(etx *echo.Context) (uuid.UUID, error) {
	app, err := kiks.Get[*cookies.App](etx.Request().Context())
	if err != nil {
		return uuid.UUID{}, err
	}
	if app == nil || app.UserID == "" {
		return uuid.UUID{}, errors.New("unauthenticated")
	}

	return uuid.Parse(app.UserID)
}

func flashAdminError(etx *echo.Context, err error) {
	kiks.AddFlash(etx.Request().Context(), kiks.FlashError, adminErrorMessage(err))
}

func adminErrorMessage(err error) string {
	if ve, ok := validation.As(err); ok && !ve.Empty() {
		if ve[0].Field != "" {
			return fmt.Sprintf("%s %s", ve[0].Field, ve[0].Message)
		}
		return ve[0].Message
	}
	switch {
	case errors.Is(err, services.ErrSlugTaken):
		return "That slug is already in use"
	case errors.Is(err, services.ErrVersionOrderMismatch):
		return "The version list is out of date. Refresh and try again"
	case errors.Is(err, services.ErrAPITokenTaken):
		return "That token value is already in use"
	case errors.Is(err, services.ErrDraftMissing):
		return "Save a draft before releasing"
	case errors.Is(err, services.ErrPageMismatch):
		return "That revision does not belong to this page"
	case errors.Is(err, models.ErrNotFound):
		return "Not found"
	default:
		return "Something went wrong"
	}
}

func newDocumentationVersionData(entity models.DocumentationVersion) DocumentationVersionData {
	return DocumentationVersionData{
		ID:             entity.ID,
		Slug:           entity.Slug,
		Label:          entity.Label,
		IsLatest:       entity.IsLatest,
		Position:       entity.Position,
		PublishedNavID: optionalInt64(entity.PublishedNavId),
		UpdatedAt:      timestamptzString(entity.UpdatedAt),
	}
}

func newDocumentationPageSummaryData(summary services.DocumentationPageSummary) DocumentationPageSummaryData {
	return DocumentationPageSummaryData{
		ID:          summary.Page.ID,
		Slug:        summary.Page.Slug,
		Title:       summary.Title,
		HasDraft:    summary.HasDraft,
		IsPublished: summary.IsPublished,
	}
}

func newDocumentationPageData(page models.DocumentationPage) DocumentationPageData {
	return DocumentationPageData{
		ID:                  page.ID,
		VersionID:           page.VersionId,
		Slug:                page.Slug,
		PublishedRevisionID: optionalInt64(page.PublishedRevisionId),
	}
}

func newDocumentationRevisionData(revision models.DocumentationRevision) DocumentationRevisionData {
	return DocumentationRevisionData{
		ID:           revision.ID,
		Status:       revision.Status,
		Title:        revision.Title,
		MetaTitle:    revision.MetaTitle,
		Description:  revision.Description,
		BodyMarkdown: revision.BodyMarkdown,
		Headings:     decodeHeadings(revision.Headings),
		CreatedAt:    timestamptzString(revision.CreatedAt),
		PublishedAt:  timestamptzString(revision.PublishedAt),
	}
}

func newDocumentationNavDraftData(revision models.DocumentationNavRevision) DocumentationNavDraftData {
	return DocumentationNavDraftData{
		ID:     revision.ID,
		Status: revision.Status,
		Tree:   navTreeData(revision.Tree),
	}
}

func navTreeData(raw []byte) []DocumentationNavNodeData {
	nodes, err := models.ParseDocumentationNavTree(raw)
	if err != nil {
		return []DocumentationNavNodeData{}
	}

	return mapNavNodes(nodes)
}

func mapNavNodes(nodes []models.DocumentationNavNode) []DocumentationNavNodeData {
	items := make([]DocumentationNavNodeData, 0, len(nodes))
	for _, node := range nodes {
		item := DocumentationNavNodeData{
			Title:    node.Title,
			PageID:   node.PageID,
			Children: mapNavNodes(node.Children),
		}
		if item.Children == nil {
			item.Children = []DocumentationNavNodeData{}
		}
		items = append(items, item)
	}

	return items
}

func mapNavNodesToModels(nodes []DocumentationNavNodeData) []models.DocumentationNavNode {
	items := make([]models.DocumentationNavNode, 0, len(nodes))
	for _, node := range nodes {
		item := models.DocumentationNavNode{
			Title:    node.Title,
			PageID:   node.PageID,
			Children: mapNavNodesToModels(node.Children),
		}
		if item.Children == nil {
			item.Children = []models.DocumentationNavNode{}
		}
		items = append(items, item)
	}

	return items
}

func decodeHeadings(raw []byte) []DocumentationHeadingData {
	if len(raw) == 0 {
		return []DocumentationHeadingData{}
	}

	var headings []DocumentationHeadingData
	if err := json.Unmarshal(raw, &headings); err != nil || headings == nil {
		return []DocumentationHeadingData{}
	}

	return headings
}

func optionalInt64(value pgtype.Int8) *int64 {
	if !value.Valid {
		return nil
	}
	id := value.Int64
	return &id
}

func timestamptzString(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}

	return value.Time.UTC().Format(time.RFC3339)
}
