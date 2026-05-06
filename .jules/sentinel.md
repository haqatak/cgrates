## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2024-05-24 - Missing Timeout in HTTP Client
**Vulnerability:** Several HTTP clients were initialized without explicit `Timeout` configurations, such as in `engine/action.go` and `engine/suretax.go`.
**Learning:** In Go, `http.Client` defaults to no timeout. This leaves the application susceptible to indefinite hangs or resource exhaustion if external servers fail to respond in a timely manner.
**Prevention:** Always set an explicit `Timeout` on `http.Client` instantiations. In the `cgrates` project, use the centralized `config.CgrConfig().GeneralCfg().ReplyTimeout` configuration to enforce system-wide request thresholds rather than hardcoding values.
