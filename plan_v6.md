Wait, `Slice()` returns ALL keys regardless of their bool value!
If the issue is returning disabled members from `utils.StringMap`, then checking `if isMember` is correct.
However, I should be careful. Look closely at what I wrote:
"ToDo: make sure we return accountIDs"

Wait!
Could "make sure we return accountIDs" mean that some entries in `MemberIds` don't have the `Tenant` prefix, and we need to ADD the tenant prefix before returning them? Or they might contain `*any` which we need to filter out?
Let's check if `MemberIds` can contain `*any`.
In `engine/sharedgroup_test.go`:
```go
		AccountParameters: map[string]*SharingParameters{
			"*any": {Strategy: STRATEGY_MINE_RANDOM},
		},
```
Here `*any` is the key in `AccountParameters`.
But what about `MemberIds`?
```go
		MemberIds: utils.NewStringMap("cgrates.org:account1"),
```
`utils.NewStringMap` creates a `StringMap` with `true` values.
What if `MemberIds` can contain `*any`?
No, it's specific members.

Wait, look at `engine/account.go`:
```go
// GetUniqueSharedGroupMembers returns the acounts from the group
func (acc *Account) GetUniqueSharedGroupMembers(cd *CallDescriptor) (utils.StringMap, error) { // ToDo: make sure we return accountIDs
...
		for memberID := range sharedGroup.MemberIds {
			memberIds[memberID] = true
		}
```
If we fix the bool value, is it sufficient?
Let's see if `MemberIds` strings are already full accountIDs `Tenant:Subject`.
Yes, because `dm.GetAccount(ubId)` requires the full key, and `sharedGroup.go:138` does:
```go
			nUb, _ = dm.GetAccount(ubId)
```
So `ubId` (which comes directly from `MemberIds`) IS the full account ID.

Wait, is it possible that the author of the ToDo wanted the function to ONLY return `acc.ID` if it is an account ID, or ensure that it doesn't return something else?
What if `sharedGroup.MemberIds` includes the shared group ID itself? No.

Wait. Is it possible `GetUniqueSharedGroupMembers` does NOT return the account itself if it isn't in `MemberIds`?
If `*any` is used in `AccountParameters`, the current account is sharing the group, but it might not be explicitly listed in `MemberIds`.
If it is not in `MemberIds`, it will not be locked by `engine/calldesc.go` because `acntIDs` won't contain it!
Wait! But `calldesc.go` always locks `cd.GetAccountKey()` itself separately!
```go
		return guardian.Guardian.Guard(func() error {
			duration, err = cd.getMaxSessionDuration(account, fltrS)
			...
		}, config.CgrConfig().GeneralCfg().LockingTimeout, lkIDs...)
	}, config.CgrConfig().GeneralCfg().LockingTimeout, utils.AccountPrefix+cd.GetAccountKey())
```
It wraps it in another `Guard` for `cd.GetAccountKey()`. So the caller's account is always locked.

What if the ToDo is EXACTLY about `utils.StringMap` returning `false` elements?
Let's do a search-and-replace to add the `ok` check.
And let's check `engine/account.go:991`:
```go
		for _, b := range balanceChain {
			for sg := range b.SharedGroups {
				groups = append(groups, sg)
			}
		}
```
Should this also be `for sg, ok := range b.SharedGroups { if ok { groups = append(groups, sg) } }`?
Let's write a small script to find all `range .*SharedGroups` and `range .*MemberIds` without `ok` check.
