# ADR 0001: Sanitize Grok tool schemas before upstream forward

**Status:** Accepted  
**Date:** 2026-07-15  
**Branch / fix:** `fix/grok-oauth-429-bounded-failover` (local hotfix also deployed as `sub2api-local:0.1.155-grok-tools-*`)

## Context

Clients (Codex + CC Switch local proxy) call sub2api with OpenAI-compatible
`/v1/chat/completions` or `/v1/responses`, targeting Grok OAuth accounts.

Codex app tools (e.g. `codex_app__automation_update`) often ship function
`parameters` whose **root** is an `anyOf` / `oneOf` union that includes a
non-object branch (`null`). xAI rejects these with HTTP 400:

```text
tool parameter root must be an object type
(root schema is an anyOf/oneOf union with a non-object branch)
```

xAI also enforces a hard cap of **250** tools per request.

Separately, xAI error bodies use a flat shape:

```json
{"code":"invalid-argument","error":"human readable string"}
```

sub2api previously only extracted `error.message` (OpenAI nested shape), so
clients only saw the useless string `Upstream error: 400`. CC Switch then
wrapped that as “local proxy failed while handling Codex endpoint /responses”.

## Decision

1. **Before** forwarding to Grok, normalize tools:
   - Collapse parameter roots that are `anyOf` / `oneOf` / `allOf` to the first
     object-like branch (or an empty object schema if none).
   - Cap tools arrays at **250** (keep first 250).
2. Apply on both paths:
   - Responses: `patchGrokResponsesBody` → `sanitizeGrokResponsesTools`
   - Chat Completions raw: `sanitizeGrokChatCompletionsTools`
3. Extend `extractUpstreamErrorMessage` to read string-valued `error` fields
   (xAI shape) so client/ops messages preserve the real upstream text.

## Consequences

- **Positive:** Codex/CC Switch + Grok no longer fails solely on union tool
  schemas; error surfaces become debuggable.
- **Trade-off:** Union alternatives beyond the first object branch are dropped;
  tools beyond index 249 are dropped. Acceptable for proxy compatibility;
  not a full JSON Schema rewrite engine.
- **Follow-up (2026-07-15):** Collapsing anyOf without copying parent `$defs`
  caused xAI `unresolvable $ref '#/$defs/__schema0'`. Collapse now merges
  parent `$defs` / `$definitions` onto the chosen object branch.
- **Tests:** unit coverage in `openai_gateway_grok_test.go` and
  `gateway_upstream_error_message_test.go` (including defs-preserving collapse).

## References

- Implementation: `backend/internal/service/openai_gateway_grok.go`
- Chat path hook: `backend/internal/service/openai_gateway_chat_completions_raw.go`
- Error extract: `backend/internal/service/gateway_upstream_response.go`
