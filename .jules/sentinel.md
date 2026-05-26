## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2026-05-10 - Dynamic Table Name SQL Injection
**Vulnerability:** In `engine/storage_sql.go` inside the `GetTpIds(colName string)` method, the `colName` argument was directly formatted into a SQL query string (`fmt.Sprintf(" (SELECT tpid FROM %s)", colName)`) without any validation or parameterization.
**Learning:** SQL parameterization (using `?` placeholders) only works for values, not for table or column names. When constructing queries dynamically with table names, directly formatting input strings creates a severe SQL injection vulnerability if the input is untrusted.
**Prevention:** Always validate dynamic table names against a strict allowlist of known, safe constants before using them in a query. Do not rely on ORM functions for table names unless they explicitly document safe handling, and avoid string formatting for query construction whenever possible.
## 2026-05-26 - Missing HTTP Client Timeouts
**Vulnerability:** Several `http.Client` instances in the `engine` package (`engine/action.go` and `engine/suretax.go`) were instantiated without an explicit `Timeout` configuration.
**Learning:** In Go, default `http.Client` instances without timeouts can cause goroutines to hang indefinitely if the remote server is slow or unresponsive. Since `engine/action.go` uses `a.ExtraParameters` to connect to dynamic arbitrary endpoints (SSRF vectors), the absence of a timeout can be explicitly leveraged for resource exhaustion.
**Prevention:** Always use explicit `Timeout` fields when instantiating `http.Client`, particularly `Timeout: config.CgrConfig().GeneralCfg().ReplyTimeout`, to enforce safe operational boundaries across all external API connections.
