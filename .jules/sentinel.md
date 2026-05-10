## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2024-05-24 - Missing Timeout in HTTP Clients
**Vulnerability:** HTTP clients in `ees/s3.go`, `ees/sqs.go`, `engine/action.go`, and `engine/suretax.go` were initialized without an explicit `Timeout` setting.
**Learning:** Default `http.Client` instances in Go have no timeout, meaning a slow or unresponsive server can cause the client connection to hang indefinitely, leading to resource exhaustion (e.g., file descriptor or goroutine leaks) and potential Denial of Service (DoS).
**Prevention:** Always configure the `Timeout` field when creating an `http.Client`. Use the application's global configuration where available, such as `config.CgrConfig().GeneralCfg().ReplyTimeout`.
