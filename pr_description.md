Title: 🔒 Fix SQL injection vulnerability in Postgres storage configuration

Description:

🎯 **What:**
This PR fixes a potential SQL injection vulnerability in `engine/storage_postgres.go` related to setting the `search_path` dynamically during database connection initialization.

⚠️ **Risk:**
If an attacker is able to control or manipulate the `pgSchema` configuration parameter, they could potentially execute arbitrary SQL statements by injecting quotes and semicolons. For instance, setting `pgSchema` to `public'; DROP TABLE users;--` would execute the trailing malicious statement using the database connection credentials, leading to data loss or tampering.

🛡️ **Solution:**
The vulnerability was fixed by properly escaping the `pgSchema` variable before interpolating it into the `SET search_path` statement. Single quotes in the provided string are now replaced with double single quotes (`''`) utilizing `strings.ReplaceAll(pgSchema, "'", "''")` to neutralize the possibility of terminating the string literal and injecting additional queries.