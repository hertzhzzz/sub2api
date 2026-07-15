package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractUpstreamErrorMessage_XAIStringErrorField(t *testing.T) {
	t.Parallel()

	// xAI Grok returns flat {"code","error"} where error is a string, not
	// OpenAI's nested {"error":{"message":"..."}}.
	body := []byte(`{"code":"invalid-argument","error":"codex_app__automation_update: tool parameter root must be an object type (root schema is an anyOf/oneOf union with a non-object branch)"}`)
	msg := extractUpstreamErrorMessage(body)
	require.Equal(t,
		"codex_app__automation_update: tool parameter root must be an object type (root schema is an anyOf/oneOf union with a non-object branch)",
		msg,
	)
}

func TestExtractUpstreamErrorMessage_PrefersNestedOpenAIShape(t *testing.T) {
	t.Parallel()

	body := []byte(`{"error":{"message":"nested message","type":"invalid_request_error"}}`)
	require.Equal(t, "nested message", extractUpstreamErrorMessage(body))
}

func TestExtractUpstreamErrorMessage_MaximumToolsLimit(t *testing.T) {
	t.Parallel()

	body := []byte(`{"code":"invalid-argument","error":"Maximum tools limit reached. 261 tools have been provided but the maximum is 250."}`)
	require.Equal(t,
		"Maximum tools limit reached. 261 tools have been provided but the maximum is 250.",
		extractUpstreamErrorMessage(body),
	)
}
