---
title: "The save that quietly ate your colleague's edit"
subtitle: "Two people open the same record, both press Save, and the second one wins. Nothing errors, nothing is logged, and the first person's work is gone. The fix is one WHERE clause and an HTTP header that has been in the spec since 1999, and most APIs still do not use it."
series: "The Daily Grit"
edition: 23
date: 2026-10-20
readingTime: "13 min"
author: "Muke JohnBaptist"
tags: [grit, go, concurrency, http, api-design, optimistic-locking, correctness]
canonical: "https://gritframework.dev/blog/the-save-that-ate-your-colleagues-edit"
---

Two support agents open the same customer record. One fixes the phone number. The other, a minute later, fixes the address. Both press Save.

The phone number is gone.

Nobody is told. No error, no log line, no conflict. The second save wrote the whole record from a form that was loaded before the first save happened, so it wrote back the old phone number along with the new address. From the database's point of view, two perfectly valid updates were applied in order. From the first agent's point of view, their work silently evaporated.

This is the **lost update** problem, and it is the most common data-integrity bug in CRUD applications, because the obvious implementation has it and it takes real work to even notice.

## Why the obvious handler has this bug

```go
func (h *CustomerHandler) Update(c *gin.Context) {
    var req UpdateCustomerRequest
    c.ShouldBindJSON(&req)
    h.DB.Model(&customer).Where("id = ?", id).Updates(req)
}
```

That is correct in the sense of doing what it says. The problem is what the client sent: a complete record, read at some earlier moment, with one field changed. `Updates` writes all of it.

Note that a transaction does not help. Both updates are individually atomic; the agents are minutes apart, not microseconds. Neither does a stricter isolation level, because there is no concurrent transaction to serialise against. The first write committed long before the second began. This is not a race in the database. It is a race between two humans and a stale form.

The ingredient missing is that the second request never says what it thinks it is overwriting.

## Pessimistic locking, and why it is not the answer here

The heavyweight fix is to lock the record when someone opens it. `SELECT ... FOR UPDATE`, or an application-level "Alice is editing this" flag.

It works, and it is the wrong tool for a web application:

- **A lock has to be released by someone.** A user who opens a record and closes their laptop holds it until a timeout. Now you need the timeout, and a way to break a stale lock, and a story for what happens when the timeout is wrong.
- **It holds a database transaction open for human time.** Seconds to minutes, across a think-time pause, consuming a connection the whole while.
- **It is pessimistic in the literal sense.** It pays the cost on every edit to prevent a conflict that is rare. Most records are edited by one person at a time.

Optimistic locking inverts it: let everyone write, detect the collision, and make the loser deal with it. You pay only when it actually happens.

## A version column, and the WHERE clause that gives it meaning

Every model Grit generates carries one:

```go
type Product struct {
    ID      string `gorm:"primarykey;size:36" json:"id"`
    Name    string `gorm:"size:255" json:"name"`
    Version int    `gorm:"not null;default:1" json:"version"`
    // ...
}

func (m *Product) BeforeUpdate(tx *gorm.DB) error {
    tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
    return nil
}
```

Every update increments it, in the database, as part of the same statement. There is no read-then-write, so two concurrent updates cannot both compute version 8 from version 7.

On its own a version column does nothing. It is a number that goes up. As the package comment says:

> Generated models already carry a version column that every update increments; this is the WHERE clause that makes it mean something.

The clause is this:

```sql
UPDATE products SET name = ?, version = version + 1
WHERE id = ? AND version = 7
```

If the record is still at version 7, one row is affected and the write lands. If somebody else got there first, it is at 8, **zero rows are affected**, and nothing is overwritten. The database tells you about the conflict by doing nothing, and the number of affected rows is how you hear about it.

## The client has to say what it is overwriting

The server cannot infer the version the client read. The client has to send it, and HTTP has had the header for this since 1999: `If-Match`.

```go
res := db.Model(&item).Scopes(concurrency.IfMatch(c)).Updates(updates)
if concurrency.Conflicted(c, res) {
    // 409, with the current version
}
```

`IfMatch` reads the header and turns it into the WHERE clause:

```go
raw := strings.TrimSpace(c.GetHeader("If-Match"))
if raw == "" || raw == "*" {
    return func(db *gorm.DB) *gorm.DB { return db }
}
c.Set(askedKey, true)
version, err := strconv.Atoi(strings.Trim(strings.TrimPrefix(raw, "W/"), `"`))
```

Three things in that small function are deliberate.

**No header means no change in behaviour.** An existing client that knows nothing about any of this keeps working exactly as before, last-write-wins. That is what makes this shippable into an API with clients you do not control. The guarantee is opt-in per request, and a client opts in by being able to lose.

**`*` is also a pass.** `If-Match: *` means "as long as it exists", which is the HTTP spec's meaning, and it is the right escape hatch for a client that wants to force a write.

**It accepts three spellings.** A bare `7`, a quoted `"7"`, and the weak ETag form `W/"7"` the API sends back. Clients echo back what they were given, and ETags are quoted by the spec, so an implementation accepting only bare integers rejects every well-behaved HTTP client. Accepting the form you emit is the minimum bar.

And a value that parses as none of those matches nothing, so a malformed header fails closed. The write does not land. Given the alternative is silently overwriting with an unverified assumption, refusing is correct.

## 409 is only useful if it says what happened

```go
func WriteConflict(c *gin.Context, current int) {
    c.Header("ETag", Tag(current))
    c.JSON(http.StatusConflict, gin.H{
        "error": gin.H{
            "code":    "VERSION_CONFLICT",
            "message": "This record changed after you loaded it. Reload it and try again.",
            "details": gin.H{"current_version": current},
        },
    })
}
```

A bare 409 is nearly useless to a client: it knows it lost and nothing else. Returning the current version lets the client re-fetch, show a real diff, and offer "keep mine / keep theirs / merge", which is the interface users deserve for this. The response also carries the new `ETag`, so a client that wants to retry already has the value to send back.

There is a comment above that function worth repeating, because it is the kind of detail that only shows up after it has bitten someone:

> Built by hand rather than through `respond.Fail`: Fail's details are strings, and `current_version` is a number on the wire. v3.285.0 sent it as `"2"`, which broke every client comparing it with the version it holds.

A generic error helper that stringifies everything is fine until one of the fields is meant to be compared rather than displayed. `"9" !== 9` in JavaScript, so a client doing the obvious check against the version it holds silently never matches, and the retry logic quietly never fires. Worth knowing if your error envelope has a `details` map typed as `map[string]string`.

## The honest limits

**It is per-record.** Two edits to the same record conflict. Two edits to related records that are only valid together will not. That needs a transaction with both rows, or a version on the aggregate.

**The window is the client's.** The guarantee holds from the version the client read. A client that reads, waits an hour, and sends `If-Match: 7` gets a conflict, which is right. A client that reads twice and sends the version from the first read has lied to you, and nothing can catch that.

**It only covers writes that go through it.** A background job, an admin script or a bulk import that writes without the version clause can still clobber. Optimistic locking is a property of every write path, not of the table.

**Counters are the wrong use.** For a field only ever incremented, `SET count = count + 1` is atomic, correct, and never conflicts. Optimistic locking is for whole-record edits where a human composed the new state from the old.

## Why this is in the framework and not in your handler

Lost updates are a bug class, not a bug. Every update route has it by default, and you discover it only when a user tells you their work vanished, if they even realise it did rather than assuming they forgot to press Save.

So it belongs where routes are generated. Every model gets the version column and the `BeforeUpdate` hook; every generated update route runs the scope and checks for a conflict; the ETag goes out on reads so clients have something to echo. The default is safe and the only thing a client has to do to get the guarantee is send back what it was given.

The broader pattern is worth naming, because it shows up well beyond this. **When two parties can act on the same state, make each one declare what it believes that state to be, and reject the ones that are wrong.** Compare-and-swap is the same idea in a CPU instruction, and ETags are the same idea for HTTP caching. Each replaces "I assume nothing changed" with "this only applies if nothing changed", which is the difference between a system that loses data quietly and one that tells you.

---

*Grit is a full-stack Go and Next.js framework. `grit generate resource` writes the version column, the If-Match scope and the 409 in this post for every resource. [Get started](/docs/start).*
