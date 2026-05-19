## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2026-05-10 - Dynamic Table Name SQL Injection
**Vulnerability:** In `engine/storage_sql.go` inside the `GetTpIds(colName string)` method, the `colName` argument was directly formatted into a SQL query string (`fmt.Sprintf(" (SELECT tpid FROM %s)", colName)`) without any validation or parameterization.
**Learning:** SQL parameterization (using `?` placeholders) only works for values, not for table or column names. When constructing queries dynamically with table names, directly formatting input strings creates a severe SQL injection vulnerability if the input is untrusted.
**Prevention:** Always validate dynamic table names against a strict allowlist of known, safe constants before using them in a query. Do not rely on ORM functions for table names unless they explicitly document safe handling, and avoid string formatting for query construction whenever possible.
## 2025-02-17 - Missing Timeouts on HTTP Clients

**Vulnerability:** Found `http.Client` initializations missing explicit `Timeout` configurations in `remoteSetAccount` (engine/action.go) and `SureTaxProcessCdr` (engine/suretax.go). Missing timeouts leave the application vulnerable to resource exhaustion (Denial of Service) via Slowloris attacks or indefinitely hanging upstream API dependencies.

**Learning:** When interacting with external systems via HTTP APIs, relying on default `http.Client` configurations leaves the system susceptible to hangs. Timeouts must be explicitly bound. Within the CGRateS codebase, `config.CgrConfig().GeneralCfg().ReplyTimeout` is the standardized way to enforce bounds on replies.

**Prevention:** Ensure that all instantiated `http.Client` and `http.Server` objects define strict operational timeouts. Code reviews should explicitly block any default client usages lacking predefined deadline/timeout contexts.
