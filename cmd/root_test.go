package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestOldProductCommandsAreRejected(t *testing.T) {
	root := NewRootCommand()
	root.SetArgs([]string{"product", "+search"})
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)

	err := root.Execute()
	if err == nil {
		t.Fatal("expected old +search command to be rejected")
	}
	if !strings.Contains(err.Error(), `unknown command "+search"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProductSearchValidationBeforeAuth(t *testing.T) {
	root := NewRootCommand()
	root.SetArgs([]string{"product", "basic", "search", "--region", "US"})
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)

	err := root.Execute()
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "--keyword is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMediaResolveOutputContract(t *testing.T) {
	root := NewRootCommand()
	root.SetArgs([]string{"media", "resolve", "--url", "https://example.com/video.mp4"})
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)

	if err := root.Execute(); err != nil {
		t.Fatalf("media resolve failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid json output: %v\n%s", err, out.String())
	}
	if envelope["ok"] != true {
		t.Fatalf("expected ok=true, got %v", envelope["ok"])
	}
	data, ok := envelope["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", envelope["data"])
	}
	for _, key := range []string{"records", "entities", "relations", "artifacts", "raw"} {
		if _, ok := data[key]; !ok {
			t.Fatalf("missing data.%s in output", key)
		}
	}
	meta, ok := envelope["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected meta object, got %T", envelope["meta"])
	}
	if meta["entity"] != "media" || meta["operation"] != "resolve" {
		t.Fatalf("unexpected meta: %#v", meta)
	}
}
