package models

import (
	"encoding/json"
	"strings"
)

const (
	DocumentationStatusDraft      = "draft"
	DocumentationStatusPublished  = "published"
	DocumentationStatusSuperseded = "superseded"
)

var EmptyJSONArray = []byte("[]")

type DocumentationNavNode struct {
	Title    string                 `json:"title"`
	PageID   *int64                 `json:"pageId"`
	Children []DocumentationNavNode `json:"children"`
}

func ParseDocumentationNavTree(raw []byte) ([]DocumentationNavNode, error) {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return []DocumentationNavNode{}, nil
	}

	var nodes []DocumentationNavNode
	if err := json.Unmarshal(raw, &nodes); err != nil {
		return nil, err
	}
	if nodes == nil {
		return []DocumentationNavNode{}, nil
	}

	return nodes, nil
}

func MarshalDocumentationNavTree(nodes []DocumentationNavNode) ([]byte, error) {
	if nodes == nil {
		nodes = []DocumentationNavNode{}
	}

	return json.Marshal(nodes)
}

func NormalizeJSONArray(raw []byte) []byte {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return append([]byte(nil), EmptyJSONArray...)
	}

	return raw
}
