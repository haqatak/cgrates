## 2024-05-24 - Missing ReadHeaderTimeout in HTTP Server
**Vulnerability:** The HTTP and HTTPS servers in `cores/server.go` were initialized using `http.ListenAndServe` or `http.Server` without setting explicitly any timeout configurations such as `ReadHeaderTimeout` or `ReadTimeout`.
**Learning:** In Go, default `http.Server` instances do not enforce timeouts for reading request headers. This leaves the application vulnerable to resource exhaustion from slow-client attacks (like Slowloris), where attackers hold connections open by sending headers very slowly.
**Prevention:** Always use explicit `http.Server` instantiation and configure timeout fields, particularly `ReadHeaderTimeout`, to safe defaults (e.g., 10 seconds) when exposing HTTP/HTTPS endpoints.

## 2026-05-10 - Dynamic Table Name SQL Injection
**Vulnerability:** In `engine/storage_sql.go` inside the `GetTpIds(colName string)` method, the `colName` argument was directly formatted into a SQL query string (`fmt.Sprintf(" (SELECT tpid FROM %s)", colName)`) without any validation or parameterization.
**Learning:** SQL parameterization (using `?` placeholders) only works for values, not for table or column names. When constructing queries dynamically with table names, directly formatting input strings creates a severe SQL injection vulnerability if the input is untrusted.
**Prevention:** Always validate dynamic table names against a strict allowlist of known, safe constants before using them in a query. Do not rely on ORM functions for table names unless they explicitly document safe handling, and avoid string formatting for query construction whenever possible.

## 2024-05-24 - Missing Timeout in HTTP Clients
**Vulnerability:** Several `http.Client` instances across the codebase (e.g., in `engine/action.go`, `engine/suretax.go`, `ees/s3.go`, `ees/sqs.go`) were initialized without a `Timeout` value.
**Learning:** Default `http.Client` initializations in Go have no timeout. This can lead to indefinite hangs, memory leaks, and resource exhaustion if the remote server is unresponsive or slow to send data.
**Prevention:** Always explicitly configure the `Timeout` field when initializing `http.Client` instances. In production code, use a safe global standard like `config.CgrConfig().GeneralCfg().ReplyTimeout`.
## 2024-05-24 - Missing Timeout in `http.Client` Instantiations
**Vulnerability:** Several `http.Client` instances in `engine/suretax.go` and `engine/action.go` were instantiated without configuring timeouts.
**Learning:** In Go, default `http.Client` instances do not enforce a timeout, which means they can hang indefinitely if the external service fails to respond. This can lead to goroutine leaks and resource exhaustion attacks (a potential Denial of Service risk).
**Prevention:** Always configure the `Timeout` field when initializing `http.Client` (e.g., using `config.CgrConfig().GeneralCfg().ReplyTimeout` where available) to enforce a strict upper bound on request durations.

## 2024-05-24 - Missing HTTP Client Timeouts
**Vulnerability:** Several `http.Client` instances in `engine/suretax.go`, `engine/action.go`, `ees/s3.go`, and `ees/sqs.go` were initialized without explicitly setting a `Timeout`.
**Learning:** Default `http.Client` instances in Go do not enforce any request timeouts. This exposes the application to resource exhaustion vulnerabilities (denial-of-service) because network calls can hang indefinitely if the external service is slow or unreachable.
**Prevention:** Always explicitly configure the `Timeout` field when initializing `http.Client`. In this application, a reusable pattern is to use the global configuration standard via `config.CgrConfig().GeneralCfg().ReplyTimeout`.
## 2024-05-25 - Missing Timeout in HTTP Client
**Vulnerability:** The HTTP clients in `engine/action.go` and `engine/suretax.go` were initialized using `&http.Client{}` without setting explicitly any timeout configurations such as `Timeout`.
**Learning:** In Go, default `http.Client` instances do not enforce timeouts for requests. This leaves the application vulnerable to resource exhaustion from slow servers or network hangs, where the client will wait indefinitely for a response.
**Prevention:** Always use explicit `http.Client` instantiation and configure the `Timeout` field to a safe default (e.g., `config.CgrConfig().GeneralCfg().ReplyTimeout`) when making outbound HTTP requests, except when interacting with systems that require long-polling or large transfers (like AWS S3/SQS).

## 2024-05-24 - Use of Weak Pseudo-Random Number Generator
**Vulnerability:** The function `RandomInteger` in `utils/coreutils.go` generated random numbers using `math/rand.Int63n`, which relies on a deterministic and predictable weak pseudo-random number generator (PRNG). If an attacker can deduce the PRNG seed or sequence, they may predict the generated outputs, compromising anything relying on the unpredictability of these integers (like session IDs, tokens, or backoffs).
**Learning:** `math/rand` should not be used for anything where unpredictability is important for security or collision resistance. Standard practice is to default to `crypto/rand`, which utilizes the operating system's cryptographically secure PRNG.
**Prevention:** Always use `crypto/rand` when generating numbers, strings, or bytes that require unpredictability. Use `math/rand` only when performance is absolutely critical and cryptographic security is strictly not required, or when determinism is intentionally desired (e.g., in controlled tests).

## 2026-05-30 - SQL Injection via PostgreSQL search_path
**Vulnerability:** In `engine/storage_postgres.go`, the `pgSchema` variable (representing the PostgreSQL schema name) was directly interpolated into a `SET search_path` SQL command (`fmt.Sprintf("set search_path='%s'", pgSchema)`) during initialization. If `pgSchema` originated from untrusted user input, it could contain unescaped single quotes, leading to a SQL injection vulnerability since SET commands cannot be parameterized natively using `?` in GORM or standard SQL interfaces.
**Learning:** PostgreSQL `SET` configuration commands (such as `SET search_path`) do not support bound parameters (placeholders like `$1` or `?`). When dynamically configuring these options, developers are often forced to use string concatenation or formatting, which opens up injection risks.
**Prevention:** When parameterization is not supported (such as in `SET` commands for configurations or dynamic table names), prevent SQL injection by rigorously validating or sanitizing the input. For single-quoted literals in PostgreSQL, any single quote in the input must be replaced with two single quotes (`strings.ReplaceAll(val, "'", "''")`) to escape it properly, or validate against an explicit allowlist.
## 2025-05-29 - [Fix SQL Injection in PostgreSQL Search Path]
**Vulnerability:** A SQL Injection vulnerability existed in `engine/storage_postgres.go` where the `pgSchema` configuration parameter was directly interpolated into a `SET search_path='%s'` statement without sanitization. An attacker able to control this parameter could inject malicious SQL commands by utilizing single quotes and semicolons.
**Learning:** In PostgreSQL, `SET` configuration commands do not support parameterized variables. When interpolating dynamic values into single-quoted SQL string literals for these commands, the proper way to sanitize is by escaping single quotes.
**Prevention:** Prevent SQL injection in non-parameterizable single-quoted string literals by replacing single quotes with two single quotes (e.g., `strings.ReplaceAll(val, "'", "''")`).

## 2026-07-06 - Use of Weak Pseudo-Random Number Generator in Data Converters
**Vulnerability:** The `RandomConverter.Convert` method in `utils/dataconverter.go` used the predictable and insecure `math/rand` package (specifically `rand.Int()` and `rand.Intn()`) to generate random integers. Since these general-purpose converters might be used in security-sensitive contexts (like generating tokens or session IDs), utilizing a weak PRNG creates a vulnerability where attackers could predict the generated outputs.
**Learning:** General-purpose data conversion functions that output randomness should assume they might be used for secure generation. A cryptographically secure PRNG must be used to ensure unpredictability. `math/rand` is insufficient for such purposes due to its determinism and predictable sequence.
**Prevention:** Always default to cryptographically secure PRNGs, such as `utils.RandomInteger` (which wraps `crypto/rand`), when generating random numbers that might have security implications. Reserve `math/rand` for non-security-critical paths where performance is the sole concern and the output predictability cannot lead to exploits.
