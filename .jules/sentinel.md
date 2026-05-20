## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2026-05-10 - Dynamic Table Name SQL Injection
**Vulnerability:** In `engine/storage_sql.go` inside the `GetTpIds(colName string)` method, the `colName` argument was directly formatted into a SQL query string (`fmt.Sprintf(" (SELECT tpid FROM %s)", colName)`) without any validation or parameterization.
**Learning:** SQL parameterization (using `?` placeholders) only works for values, not for table or column names. When constructing queries dynamically with table names, directly formatting input strings creates a severe SQL injection vulnerability if the input is untrusted.
**Prevention:** Always validate dynamic table names against a strict allowlist of known, safe constants before using them in a query. Do not rely on ORM functions for table names unless they explicitly document safe handling, and avoid string formatting for query construction whenever possible.
## 2024-05-24 - Missing Timeout in `http.Client` Instantiations
**Vulnerability:** Several `http.Client` instances in `engine/suretax.go` and `engine/action.go` were instantiated without configuring timeouts.
**Learning:** In Go, default `http.Client` instances do not enforce a timeout, which means they can hang indefinitely if the external service fails to respond. This can lead to goroutine leaks and resource exhaustion attacks (a potential Denial of Service risk).
**Prevention:** Always configure the `Timeout` field when initializing `http.Client` (e.g., using `config.CgrConfig().GeneralCfg().ReplyTimeout` where available) to enforce a strict upper bound on request durations.
