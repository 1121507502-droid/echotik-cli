package runner

import (
	"errors"

	"github.com/echotik/cli/internal/client"
	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/normalize"
	"github.com/echotik/cli/internal/output"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

type Request struct {
	Entity       string
	Capability   string
	Operation    string
	Freshness    string
	Method       string
	Path         string
	Params       map[string]string
	Body         any
	RecordsAs    string
	RelationFrom *normalize.RelationNode
	RelationType string
}

func Run(cmd *cobra.Command, req Request) error {
	if req.Method == "" {
		req.Method = "GET"
	}
	baseURL, username, password, err := core.ResolveCredential()
	if err != nil {
		return output.NewError("authentication_error", err.Error(), "run: echotik config set-credential")
	}
	c := client.New(baseURL, username, password)
	resp, err := c.Do(cmd.Context(), client.Request{
		Method: req.Method,
		Path:   req.Path,
		Params: schema.Clean(req.Params),
		Body:   req.Body,
	})
	if err != nil {
		return output.NewError("api_error", err.Error(), "check parameters or retry with backoff")
	}
	raw := resp.JSON
	if raw == nil {
		raw = string(resp.Raw)
	}
	data := normalize.FromRaw(raw, normalize.Options{
		Entity:       choose(req.RecordsAs, req.Entity),
		RelationFrom: req.RelationFrom,
		RelationType: req.RelationType,
	})
	return output.Success(cmd.OutOrStdout(), data, map[string]interface{}{
		"entity":     req.Entity,
		"capability": req.Capability,
		"operation":  req.Operation,
		"freshness":  req.Freshness,
		"path":       req.Path,
		"params":     schema.Clean(req.Params),
		"statusCode": resp.StatusCode,
	})
}

func Unsupported(cmd *cobra.Command, entity, capability, operation, hint string) error {
	return output.NewError(
		"unsupported_operation",
		entity+" "+capability+" "+operation+" is not mapped to an EchoTik API endpoint yet",
		hint,
	)
}

func Validate(err error) error {
	if err == nil {
		return nil
	}
	var validation *schema.ValidationError
	if errors.As(err, &validation) {
		return output.NewError("validation_error", validation.Message, validation.Hint)
	}
	return err
}

func choose(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
