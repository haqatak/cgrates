## 2026-04-26 - [Fix SQL Injection in GetTpTableIds]
**Vulnerability:** SQL injection vulnerability via query string formatting in `engine/storage_sql.go`. User input could be injected into `WHERE tpid='%s'` and `AND %s='%s'` parameters.
**Learning:** `SQLStorage` supports both raw DB operations `sqls.Db` and GORM operations `sqls.db`. GORM allows parameterized queries effectively, although its `.Where` interface natively supports raw parameterized queries for values but can be further secured using `.Where(map[string]interface{}{key: value})` for dynamic column names. Note that distinct keys from `utils.TPDistinctIds` should not be backticked because it breaks PostgreSQL.
**Prevention:** Always use GORM's `.Table()` and `.Where()` methods for query building with parameterized inputs (`?` or maps) rather than raw string concatenation `fmt.Sprintf` for SQL query construction.

## 2026-04-26 - [Fix SQL Injection in GetCDRs filters]
**Vulnerability:** SQL injection vulnerability via query string formatting in `engine/storage_sql.go` inside the `GetCDRs` method. Unsanitized `destinationPrefix` user input could be injected into `LIKE` clauses (e.g., `fmt.Sprintf(" destination LIKE '%s%%'", destPrefix)`).
**Learning:** Even when building dynamic SQL blocks (like `OR` or `AND` condition strings), parameterized variables should always be injected via `?` slice arguments inside GORM `Where()` methods.
**Prevention:** Instead of string formatting loops and `bytes.Buffer`, build an array of `conditions` ("destination LIKE ?") and `values` ("prefix%"), then apply them as `q.Where(strings.Join(conditions, " OR "), values...)`.
