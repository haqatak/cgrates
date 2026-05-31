Wait, look at `engine/sharedgroup.go:133`:
```go
	for ubId := range sg.MemberIds {
		var nUb *Account
```
If `engine/sharedgroup.go` also does `for ubId := range sg.MemberIds` without checking the bool value, then maybe the bool value is ALWAYS `true` when a member is present, and absent (`delete`) when not?
If `delete` is used, then `range` only yields present keys.
Let's see how `MemberIds` is populated.
In `engine/datamanager_test.go`:
```go
		MemberIds: utils.StringMap{
			"string1": true,
			"string2": false,
		},
```
This is a test. Is `false` ever used in production code?
Wait. If we look at the ToDo again: `// ToDo: make sure we return accountIDs`.
What else could it mean?
Could `accountIDs` mean a slice of strings?
The return type is `(utils.StringMap, error)`. `utils.StringMap` IS `map[string]bool`. So returning a map of account IDs as keys is exactly what `utils.StringMap` is.

Is it possible that `MemberIds` can contain `*any`?
No, `*any` is used in `AccountParameters`, not `MemberIds`. `MemberIds` is specific.

Could "accountIDs" mean we need to make sure we return the actual IDs, not something else like the `Tenant` alone?
What if `GetUniqueSharedGroupMembers` needs to return account IDs without the tenant, or with the tenant?
What does `cd.GetAccountKey()` return? `utils.ConcatenatedKey(cd.Tenant, subj)`. This returns `tenant:subject`.
So `acntIDs` MUST contain `tenant:subject`.

Wait... is there a chance that `sharedGroup.MemberIds` does NOT contain the full `tenant:subject` but just `subject`?
Let's check how `SetSharedGroup` handles `MemberIds`.
