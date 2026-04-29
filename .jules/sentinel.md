## 2026-04-26 - [Fix SQL Injection in GetTpTableIds]
**Vulnerability:** SQL injection vulnerability via query string formatting in `engine/storage_sql.go`. User input could be injected into `WHERE tpid='%s'` and `AND %s='%s'` parameters.
**Learning:** `SQLStorage` supports both raw DB operations `sqls.Db` and GORM operations `sqls.db`. GORM allows parameterized queries effectively, although its `.Where` interface natively supports raw parameterized queries for values but can be further secured using `.Where(map[string]interface{}{key: value})` for dynamic column names. Note that distinct keys from `utils.TPDistinctIds` should not be backticked because it breaks PostgreSQL.
**Prevention:** Always use GORM's `.Table()` and `.Where()` methods for query building with parameterized inputs (`?` or maps) rather than raw string concatenation `fmt.Sprintf` for SQL query construction.

## 2026-04-26 - [Fix SQL Injection in GetCDRs filters]
**Vulnerability:** SQL injection vulnerability via query string formatting in `engine/storage_sql.go` inside the `GetCDRs` method. Unsanitized `destinationPrefix` user input could be injected into `LIKE` clauses (e.g., `fmt.Sprintf(" destination LIKE '%s%%'", destPrefix)`).
**Learning:** Even when building dynamic SQL blocks (like `OR` or `AND` condition strings), parameterized variables should always be injected via `?` slice arguments inside GORM `Where()` methods.
**Prevention:** Instead of string formatting loops and `bytes.Buffer`, build an array of `conditions` ("destination LIKE ?") and `values` ("prefix%"), then apply them as `q.Where(strings.Join(conditions, " OR "), values...)`.

## 2026-04-26 - [Fix SQL Injection in SQLEe PrepareMap and PrepareOrderMap]
**Vulnerability:** SQL injection vulnerability via query string formatting in `ees/sql.go`. User input from `cgrEv.Event` (map keys) could be injected into `INSERT INTO` or `UPDATE` queries using `fmt.Sprintf` without properly escaping the keys, allowing arbitrary SQL execution if map keys contained malicious backticks or strings.
**Learning:** `PrepareMap` and `PrepareOrderMap` built queries via `fmt.Sprintf` directly instead of passing map values natively to GORM `Create` or `Updates`. GORM can handle `map[string]interface{}` natively for updates and inserts which escapes both column names and values automatically. `sqlPosterRequest` struct was updated to store `Create`, `Update`, and `Where` maps natively to rely on GORM's `.Table(sqlEe.tableName).Create()` and `.Updates()`.
**Prevention:** Always use GORM's built-in parameterization and map creation features rather than manually concatenating strings and slices of `?` parameters to defend against user-supplied column names.
