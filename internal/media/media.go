package media

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/echotik/cli/internal/normalize"
)

type DownloadResult struct {
	Artifacts []normalize.Artifact `json:"artifacts"`
	Manifest  string               `json:"manifest"`
}

func ArtifactFromURL(rawURL string) normalize.Artifact {
	typ := "media"
	lower := strings.ToLower(rawURL)
	switch {
	case strings.Contains(lower, ".mp4") || strings.Contains(lower, "video"):
		typ = "video"
	case strings.Contains(lower, ".jpg") || strings.Contains(lower, ".jpeg") || strings.Contains(lower, ".png") || strings.Contains(lower, ".webp") || strings.Contains(lower, "image"):
		typ = "image"
	}
	return normalize.Artifact{Type: typ, SourceURL: rawURL}
}

func Download(rawURL, outputDir string) (*DownloadResult, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("url is required")
	}
	if outputDir == "" {
		outputDir = "echotik-assets"
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}
	name := fileName(rawURL, resp.Header.Get("Content-Type"))
	path := filepath.Join(outputDir, name)
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, hash), resp.Body); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return nil, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return nil, err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return nil, err
	}
	sum := hex.EncodeToString(hash.Sum(nil))
	artifact := ArtifactFromURL(rawURL)
	artifact.LocalPath = path
	artifact.SHA256 = sum
	manifestPath := filepath.Join(outputDir, "manifest.json")
	result := &DownloadResult{
		Artifacts: []normalize.Artifact{artifact},
		Manifest:  manifestPath,
	}
	_ = writeManifest(manifestPath, result)
	return result, nil
}

func fileName(rawURL, contentType string) string {
	parsed, err := url.Parse(rawURL)
	if err == nil {
		base := filepath.Base(parsed.Path)
		if base != "." && base != "/" && base != "" {
			return sanitize(base)
		}
	}
	ext := ".bin"
	if strings.Contains(contentType, "mp4") {
		ext = ".mp4"
	} else if strings.Contains(contentType, "jpeg") {
		ext = ".jpg"
	} else if strings.Contains(contentType, "png") {
		ext = ".png"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	}
	return fmt.Sprintf("echotik-%d%s", time.Now().UnixNano(), ext)
}

func sanitize(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", "?", "_", "&", "_", "=", "_", ":", "_")
	return replacer.Replace(name)
}

func writeManifest(path string, result *DownloadResult) error {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
