## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2026-05-10 - Dynamic Table Name SQL Injection
**Vulnerability:** In `engine/storage_sql.go` inside the `GetTpIds(colName string)` method, the `colName` argument was directly formatted into a SQL query string (`fmt.Sprintf(" (SELECT tpid FROM %s)", colName)`) without any validation or parameterization.
**Learning:** SQL parameterization (using `?` placeholders) only works for values, not for table or column names. When constructing queries dynamically with table names, directly formatting input strings creates a severe SQL injection vulnerability if the input is untrusted.
**Prevention:** Always validate dynamic table names against a strict allowlist of known, safe constants before using them in a query. Do not rely on ORM functions for table names unless they explicitly document safe handling, and avoid string formatting for query construction whenever possible.

## 2026-05-25 - Missing Timeout in HTTP Client
**Vulnerability:** In `engine/action.go` and `engine/suretax.go`, the `http.Client` instances used to interact with external APIs (`remoteSetAccount` and `SureTaxProcessCdr`) did not configure a `Timeout` field.
**Learning:** In Go, the default `http.Client` (and instances created without an explicit `Timeout`) has no timeout. This can lead to the client waiting indefinitely for a response if the external server hangs or experiences network issues. This vulnerability can result in resource exhaustion (e.g., hanging goroutines) leading to a Denial of Service (DoS).
**Prevention:** Always set an explicit `Timeout` when initializing an `http.Client`. In production code like the `engine` package, use a configurable, global standard such as `config.CgrConfig().GeneralCfg().ReplyTimeout` to ensure all remote calls eventually return.
