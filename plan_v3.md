Wait, let's look at another file to see if `accountIDs` usually refers to the ID WITHOUT the tenant prefix!
If `accountID` means strictly the `subj` and not the `tenant:subj`, I need to be sure.

If I look at `utils/coreutils.go` or `engine/models.go`, how are Account IDs defined?
```go
// AccountID  string  // reference for shared balances
```
Is `AccountID` the full key?
Yes, usually `Account` keys in CGRateS maps are the full `tenant:account` string.
`cd.GetAccountKey()` returns `utils.ConcatenatedKey(cd.Tenant, subj)`, which is the full string.
In `calldesc.go:688`:
```go
		var lkIDs []string
		for acntID := range acntIDs {
			if acntID != cd.GetAccountKey() {
				lkIDs = append(lkIDs, utils.AccountPrefix+acntID)
			}
		}
```
If `acntID` were just the `subj`, then `acntID != cd.GetAccountKey()` would be `subj != "tenant:subj"`, which is always TRUE! So `lkIDs` would incorrectly include the caller's account locking key (with missing tenant maybe).
And `utils.AccountPrefix+acntID` would become `*account:subj` instead of `*account:tenant:subj`.
But wait! Does `utils.AccountPrefix+acntID` expect the full key?
Yes. Locking keys in CGRateS are `*account:tenant:subj`. So `acntID` MUST be the full key.

So what does "ToDo: make sure we return accountIDs" mean?
Let's reconsider the `false` values in `MemberIds`.
If `isMember` is false, it shouldn't be returned. This is definitely a bug in the current code, but is it what the ToDo refers to?
Look at the ToDo again: `// ToDo: make sure we return accountIDs`.
Maybe there is a risk that `sharedGroup.MemberIds` contains `*any`?
Wait. `SharedGroup` struct has:
```go
type SharedGroup struct {
	Id                string
	AccountParameters map[string]*SharingParameters
	MemberIds         utils.StringMap
	//members           []*Account // accounts caching
}
```
`AccountParameters` keys are `AccountID`s or `*any`.
Wait, `MemberIds` is `utils.StringMap`. It contains the explicit member IDs.

Wait, look at `engine/sharedgroup.go`:
```go
func (sg *SharedGroup) GetBalances(destination, category, balanceType string, ub *Account, aTime time.Time) (bc Balances) {
	//	if len(sg.members) == 0 {
	for ubId := range sg.MemberIds {
		var nUb *Account
		if ubId == ub.ID { // skip the initiating user
			nUb = ub
		} else {
			nUb, _ = dm.GetAccount(ubId)
			if nUb == nil || nUb.Disabled {
				continue
			}
		}
        // ...
```
Wait! `for ubId := range sg.MemberIds`
In `GetBalances`, it ignores the boolean value of `sg.MemberIds` as well!
Wait... If `GetBalances` does this, then maybe `utils.StringMap` doesn't use `false`?
But we know it can: `MemberIds: utils.StringMap{"string1": true, "string2": false}` in tests.
Ah, `utils.StringMap` is generally just `map[string]bool`. Why would it have `false`?
If you remove a key from `utils.StringMap`, do you `delete(m, k)` or `m[k] = false`?
Usually `delete` is better. But sometimes it is `m[k] = false`.
Let's see if there is a difference between `for k := range m` and `for k, v := range m { if v { ... } }`.

Wait! Look closely at `account.go` lines 1019-1021:
```go
		for memberID := range sharedGroup.MemberIds {
			memberIds[memberID] = true
		}
```
If we change it to:
```go
		for memberID, isMember := range sharedGroup.MemberIds {
			if isMember {
				memberIds[memberID] = true
			}
		}
```
Does it return "accountIDs"?
Wait. Is there a possibility that `sharedGroupIds` are returned instead of `memberID`?
No, the code says:
```go
	memberIds := make(utils.StringMap)
	for _, sgID := range sharedGroupIds {
        ...
		for memberID := range sharedGroup.MemberIds {
			memberIds[memberID] = true
		}
```

Wait, what if `memberID` in `sharedGroup.MemberIds` can be something other than an accountID?
What if `AccountID` refers to the `Account` property of `SharedGroupMdl`?
No, `SharedGroupMdl` in `models.go`:
```go
type SharedGroupMdl struct {
	Id            int64
	Tpid          string
	Tag           string `index:"0" re:".*"`
	Account       string `index:"1" re:".*"`
	Strategy      string `index:"2" re:".*"`
	RatingSubject string `index:"3" re:".*"`
}
```
Wait! In `models.go`, `Account` is a regex, like `.*` or `*any`.
If `Account` is a regex or `*any`, then `MemberIds` might be populated with actual account IDs dynamically?
No, `MemberIds` is an exact map of account IDs. The API allows adding specific accounts to a shared group.

Let's read the ToDo again: `// ToDo: make sure we return accountIDs`.
Maybe it's just a reminder to verify that `sharedGroup.MemberIds` actually holds `accountID`s? Yes, they do.
But wait! What if the `memberID` needs to have the tenant prefixed if it doesn't?
If `memberID` doesn't have a tenant, `dm.GetAccount(ubId)` handles adding default tenant if configured.

Let's look at `GetUniqueSharedGroupMembers` definition again:
```go
// GetUniqueSharedGroupMembers returns the acounts from the group
func (acc *Account) GetUniqueSharedGroupMembers(cd *CallDescriptor) (utils.StringMap, error) { // ToDo: make sure we return accountIDs
```
Notice the typo: `returns the acounts`.
Maybe we should just implement the `isMember` check as it's the most correct way to handle `utils.StringMap`, and also check if we should be returning `acc.ID` as well? Wait. If the account is in the shared group, its ID is in `MemberIds`. But what if the account has a shared group but its ID is NOT in `MemberIds`? (e.g. it was dynamically applied via `*any`).
If an account has a balance with a shared group, but the account is NOT explicitly in `sharedGroup.MemberIds`?
Ah!!!
If a shared group has `*any` in `AccountParameters`, then ANY account can use it.
So the account using it might NOT be in `sharedGroup.MemberIds`!
If so, `acc.ID` might not be in the returned `memberIds` map.
But the function is `GetUniqueSharedGroupMembers`, which gets ALL members of the shared groups this account is interacting with.
Wait, if it's meant to return *account IDs* that will be locked, should it also lock the accounts that are NOT in `MemberIds` but are sharing the group? We don't know them! We only know the explicit `MemberIds`.

Let's think about `utils.StringMap`. It's defined as `map[string]bool`.
```go
		for memberID, ok := range sharedGroup.MemberIds {
			if ok {
				memberIds[memberID] = true
			}
		}
```
Is this what the ToDo meant? Let's check `engine/account.go` and see if `ok` check is common.
