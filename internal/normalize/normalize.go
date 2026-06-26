package normalize

import (
	"fmt"
	"strings"
)

type Data struct {
	Records   []any      `json:"records"`
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
	Artifacts []Artifact `json:"artifacts"`
	Raw       any        `json:"raw"`
}

type Entity struct {
	Type string         `json:"type"`
	ID   string         `json:"id"`
	Name string         `json:"name,omitempty"`
	Raw  map[string]any `json:"raw,omitempty"`
}

type Relation struct {
	From RelationNode   `json:"from"`
	To   RelationNode   `json:"to"`
	Type string         `json:"type"`
	Raw  map[string]any `json:"raw,omitempty"`
}

type RelationNode struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Artifact struct {
	Type        string `json:"type"`
	SourceURL   string `json:"sourceUrl,omitempty"`
	DownloadURL string `json:"downloadUrl,omitempty"`
	LocalPath   string `json:"localPath,omitempty"`
	ExpiresAt   string `json:"expiresAt,omitempty"`
	SHA256      string `json:"sha256,omitempty"`
}

type Options struct {
	Entity       string
	RelationFrom *RelationNode
	RelationType string
}

func FromRaw(raw any, opts Options) Data {
	records := ExtractRecords(raw)
	data := Data{
		Records:   records,
		Entities:  []Entity{},
		Relations: []Relation{},
		Artifacts: []Artifact{},
		Raw:       raw,
	}
	seenEntities := map[string]bool{}
	seenArtifacts := map[string]bool{}
	for _, record := range records {
		m, ok := record.(map[string]any)
		if !ok {
			continue
		}
		for _, entity := range ExtractEntities(m) {
			key := entity.Type + ":" + entity.ID
			if entity.ID != "" && !seenEntities[key] {
				data.Entities = append(data.Entities, entity)
				seenEntities[key] = true
			}
			if opts.RelationFrom != nil && entity.ID != "" && entity.Type != opts.RelationFrom.Type {
				relType := opts.RelationType
				if relType == "" {
					relType = opts.RelationFrom.Type + "_related_to_" + entity.Type
				}
				data.Relations = append(data.Relations, Relation{
					From: *opts.RelationFrom,
					To: RelationNode{
						Type: entity.Type,
						ID:   entity.ID,
					},
					Type: relType,
					Raw:  m,
				})
			}
		}
		for _, artifact := range ExtractArtifacts(m) {
			key := artifact.Type + ":" + artifact.SourceURL + ":" + artifact.DownloadURL
			if key != "::" && !seenArtifacts[key] {
				data.Artifacts = append(data.Artifacts, artifact)
				seenArtifacts[key] = true
			}
		}
	}
	if opts.Entity != "" {
		for _, record := range records {
			m, ok := record.(map[string]any)
			if !ok {
				continue
			}
			id := idForEntity(opts.Entity, m)
			if id == "" {
				continue
			}
			key := opts.Entity + ":" + id
			if !seenEntities[key] {
				data.Entities = append(data.Entities, Entity{
					Type: opts.Entity,
					ID:   id,
					Name: nameForEntity(opts.Entity, m),
					Raw:  m,
				})
				seenEntities[key] = true
			}
		}
	}
	return data
}

func ExtractRecords(raw any) []any {
	if raw == nil {
		return []any{}
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return []any{raw}
	}
	data, ok := m["data"]
	if !ok || data == nil {
		return []any{raw}
	}
	switch v := data.(type) {
	case []any:
		return v
	case map[string]any:
		for _, key := range []string{"list", "items", "records", "results", "e_com_items", "comments"} {
			if nested, ok := v[key].([]any); ok {
				return nested
			}
		}
		return []any{v}
	default:
		return []any{v}
	}
}

func ExtractEntities(m map[string]any) []Entity {
	var entities []Entity
	for _, spec := range []struct {
		entity string
		idKeys []string
	}{
		{"product", []string{"product_id", "productId"}},
		{"shop", []string{"seller_id", "sellerId", "shop_id"}},
		{"creator", []string{"user_id", "userId", "unique_id", "uniqueId"}},
		{"video", []string{"video_id", "videoId"}},
		{"live", []string{"room_id", "roomId"}},
		{"category", []string{"category_id", "category_l2_id", "category_l3_id"}},
	} {
		id := firstString(m, spec.idKeys...)
		if id == "" {
			continue
		}
		entities = append(entities, Entity{
			Type: spec.entity,
			ID:   id,
			Name: nameForEntity(spec.entity, m),
			Raw:  m,
		})
	}
	return entities
}

func ExtractArtifacts(m map[string]any) []Artifact {
	var artifacts []Artifact
	for _, spec := range []struct {
		key string
		typ string
	}{
		{"cover_url", "image"},
		{"reflow_cover", "image"},
		{"avatar", "image"},
		{"play_addr", "video"},
		{"download_url", "video"},
		{"downloadUrl", "video"},
		{"play_url", "video"},
		{"url", "media"},
	} {
		if value := firstString(m, spec.key); strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
			artifact := Artifact{Type: spec.typ, SourceURL: value}
			if spec.key == "download_url" || spec.key == "downloadUrl" {
				artifact.DownloadURL = value
				artifact.SourceURL = ""
			}
			artifacts = append(artifacts, artifact)
		}
	}
	return artifacts
}

func idForEntity(entity string, m map[string]any) string {
	switch entity {
	case "product":
		return firstString(m, "product_id", "productId")
	case "shop":
		return firstString(m, "seller_id", "sellerId", "shop_id")
	case "creator":
		return firstString(m, "user_id", "userId", "unique_id", "uniqueId")
	case "video":
		return firstString(m, "video_id", "videoId")
	case "live":
		return firstString(m, "room_id", "roomId")
	default:
		return ""
	}
}

func nameForEntity(entity string, m map[string]any) string {
	switch entity {
	case "product":
		return firstString(m, "product_name", "title", "name")
	case "shop":
		return firstString(m, "seller_name", "shop_name", "name")
	case "creator":
		return firstString(m, "nick_name", "unique_id", "name")
	case "video":
		return firstString(m, "video_desc", "desc", "title")
	default:
		return firstString(m, "name", "title")
	}
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch typed := v.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return typed
			}
		case fmt.Stringer:
			return typed.String()
		case float64:
			return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%.0f", typed), ".0"), ".")
		case int:
			return fmt.Sprintf("%d", typed)
		case int64:
			return fmt.Sprintf("%d", typed)
		}
	}
	return ""
}
