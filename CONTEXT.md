# CONTEXT — sub2api (agent glossary)

Short, durable terms for agents. Prefer ADRs for decisions; this file is the
vocabulary map.

## Grok / xAI

| Term | Meaning |
|------|---------|
| **Grok tool root** | Function tool `parameters` root must be a pure JSON Schema `object`. xAI rejects `anyOf`/`oneOf` roots that include non-object branches (e.g. `null`). |
| **Grok max tools** | Hard limit **250** tools per request. |
| **xAI error shape** | `{"code":"...","error":"<string>"}` — not OpenAI’s nested `error.message`. |
| **Tool schema sanitization** | Pre-forward collapse of non-object parameter roots + tool cap. When collapsing anyOf, **parent `$defs` must move onto the chosen branch** so `$ref` still resolves. See [ADR 0001](docs/adr/0001-grok-tool-schema-sanitization.md). |

## Client path (local Mark setup)

| Term | Meaning |
|------|---------|
| **CC Switch** | Desktop tray proxy on `127.0.0.1:15721`. Codex provider `ohmyapi` points base_url at local sub2api (`http://localhost:8080`). |
| **Codex `/responses`** | What Codex sends to CC Switch; may be rewritten to chat/completions before Grok. |
| **ohmyapi** | CC Switch provider name for this sub2api instance (not a separate SaaS). |

## Hotfix deploy (local)

| Term | Meaning |
|------|---------|
| **hotfix image** | `sub2api-local:0.1.155-grok-tools-*` built from `sub2api-deploy/Dockerfile.hotfix`. |
| **Cross-compile** | Binary must be `GOOS=linux GOARCH=arm64` for Docker Desktop on Apple Silicon. macOS Mach-O fails at runtime with shell parse errors. |

## Where to look

- Grok forward + sanitizers: `backend/internal/service/openai_gateway_grok.go`
- Raw chat path: `backend/internal/service/openai_gateway_chat_completions_raw.go`
- Upstream error message parse: `backend/internal/service/gateway_upstream_response.go`
- Local compose override: sibling repo `sub2api-deploy/` (see `README-HOTFIX.md` there)
