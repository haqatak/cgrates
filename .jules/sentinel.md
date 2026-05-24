## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2026-05-10 - Dynamic Table Name SQL Injection
**Vulnerability:** In `engine/storage_sql.go` inside the `GetTpIds(colName string)` method, the `colName` argument was directly formatted into a SQL query string (`fmt.Sprintf(" (SELECT tpid FROM %s)", colName)`) without any validation or parameterization.
**Learning:** SQL parameterization (using `?` placeholders) only works for values, not for table or column names. When constructing queries dynamically with table names, directly formatting input strings creates a severe SQL injection vulnerability if the input is untrusted.
**Prevention:** Always validate dynamic table names against a strict allowlist of known, safe constants before using them in a query. Do not rely on ORM functions for table names unless they explicitly document safe handling, and avoid string formatting for query construction whenever possible.

## 2024-05-27 - [Sentinel] HTTP Client Timeout Configuration
**Vulnerability:** HTTP clients initialized without explicit timeouts (e.g., `&http.Client{}`) or with unbounded default configurations are vulnerable to hanging indefinitely if the remote server fails to respond, potentially leading to resource exhaustion (goroutine and connection leaks) and denial-of-service conditions.
**Learning:** The default behavior of `http.Client` in Go is to not impose any timeouts on network requests. This issue was found in `remoteSetAccount` in `engine/action.go` and `SureTaxProcessCdr` in `engine/suretax.go`.
**Prevention:** Always configure `Timeout` when initializing an `http.Client`. Use established system-wide timeout configurations like `config.CgrConfig().GeneralCfg().ReplyTimeout` to ensure consistent connection policies across the application.
