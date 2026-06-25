package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/echotik/cli/internal/client"
	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var paramsJSON string
	var dataJSON string
	cmd := &cobra.Command{
		Use:   "api <method> <path>",
		Short: "Call a raw EchoTik API endpoint",
		Args:  cobra.ExactArgs(2),
		Example: `  echotik api GET /api/v3/echotik/product/list --params '{"page":"1","page_size":"20"}'
  echotik api GET /api/v3/realtime/product/search --params '{"keyword":"sunscreen"}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params, err := parseStringMap(paramsJSON)
			if err != nil {
				return output.NewError("validation_error", "invalid --params: "+err.Error(), "pass a JSON object, e.g. --params '{\"page\":\"1\"}'")
			}
			body, err := parseBody(dataJSON)
			if err != nil {
				return output.NewError("validation_error", "invalid --data: "+err.Error(), "pass a JSON object or array")
			}
			return call(cmd, strings.ToUpper(args[0]), args[1], params, body)
		},
	}
	cmd.Flags().StringVar(&paramsJSON, "params", "", "query parameters as JSON object")
	cmd.Flags().StringVar(&dataJSON, "data", "", "request body as JSON")
	return cmd
}

func call(cmd *cobra.Command, method, path string, params map[string]string, body any) error {
	baseURL, username, password, err := core.ResolveCredential()
	if err != nil {
		return output.NewError("authentication_error", err.Error(), "run: echotik config set-credential")
	}
	c := client.New(baseURL, username, password)
	resp, err := c.Do(cmd.Context(), client.Request{Method: method, Path: path, Params: params, Body: body})
	if err != nil {
		return writeAPIError(cmd, err)
	}
	data := resp.JSON
	if data == nil {
		data = string(resp.Raw)
	}
	return output.Success(cmd.OutOrStdout(), data, map[string]interface{}{
		"statusCode": resp.StatusCode,
		"path":       path,
	})
}

func writeAPIError(cmd *cobra.Command, err error) error {
	var typ = "api_error"
	var hint string
	if strings.Contains(err.Error(), "HTTP 401") {
		typ = "authentication_error"
		hint = "run: echotik config set-credential"
	}
	if strings.Contains(err.Error(), "HTTP 429") {
		typ = "rate_limit"
		hint = "retry with backoff"
	}
	if strings.Contains(err.Error(), "HTTP 500") {
		typ = "server_error"
		hint = "retry with backoff; realtime EchoTik endpoints can be risk-controlled"
	}
	return output.NewError(typ, err.Error(), hint)
}

func parseStringMap(raw string) (map[string]string, error) {
	out := map[string]string{}
	if strings.TrimSpace(raw) == "" {
		return out, nil
	}
	var anyMap map[string]any
	if err := json.Unmarshal([]byte(raw), &anyMap); err != nil {
		return nil, err
	}
	for k, v := range anyMap {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

func parseBody(raw string) (any, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var out any
	err := json.Unmarshal([]byte(raw), &out)
	return out, err
}
