---
title: "Your App Breaks in a Lift: How Offline-First Sync Actually Works"
subtitle: "Two devices, one row, no network, and both people pressed Save. Every offline app has to answer that, and the answers are mostly not about networking at all. Here is the whole engine Grit generates: the outbox that squashes, the keyset cursor that stops losing rows, the tombstone that makes deletes travel, and the conflict you must never resolve silently."
series: "The Daily Grit"
edition: 25
date: 2026-10-22
readingTime: "20 min"
author: "Muke JohnBaptist"
tags: [grit, go, offline-first, sync, distributed-systems, sqlite, correctness, desktop]
canonical: "https://gritframework.dev/blog/offline-first-sync-in-go"
---

A field engineer drives into a basement car park to inspect a building. No signal for forty minutes. She updates eleven records. Meanwhile her colleague, above ground, edits one of the same eleven.

She comes back up. Now what?

That question is the whole of offline sync, and almost none of it is about networking. It is about what you do with two versions of the truth and no authority to pick between them. This post walks through the engine Grit generates for that, both halves, with the code.

## The mental model: a branch, not a cache

The most common mistake is to think of offline data as a **cache**. A cache is a copy of someone else's truth that you throw away when it goes stale. If that were what this is, the problem would be easy: refetch and move on.

It is not a cache. It is a **branch**.

The moment the device goes offline and somebody edits a row, that device has forked. It has its own commits. The server has carried on and may have its own. When the network returns you are not refreshing a cache, you are **merging two branches**, and every hard part of offline sync is a hard part of merging: what changed, in what order, and what do we do when both sides touched the same line.

Hold that analogy. Everything below is a piece of machinery that a version control system also needs, and the resemblance is not a coincidence.

## The two tables on the device

```go
// Data model:
//   - records — local mirror of every row the user has touched. Reads
//     come from this table; pulls UPSERT into it.
//   - outbox  — pending local changes that haven't been pushed yet.
//     One entry per (table_name, entity_id) — multiple edits squash.
```

**`records` is the working copy.** Every read in the app comes from here, always, even when online. That single decision is what makes the app feel instant and is what "offline-first" actually means: the network is never on the path between a user and their data. It is a background process that reconciles.

**`outbox` is the stack of unpushed commits.** A local write does two things, and the order matters:

```go
// LocalCreate/Update/Delete writes both records (so reads see the
// change immediately) AND an outbox entry.
```

Write the mirror so the UI updates instantly. Write the outbox so the change survives the app being closed, the battery dying, or the user not reconnecting until Tuesday. An outbox that lives in memory is not an outbox; it is a plan to lose data at the first crash.

## Squashing: why the outbox is keyed, not appended

Here is the first design decision that is not obvious. The outbox is not a log of every change. It holds **at most one row per record**:

```go
type Outbox struct {
    ID       int64  `gorm:"primaryKey;autoIncrement"`
    Model    string `gorm:"size:50;uniqueIndex:idx_outbox_entity"`
    EntityID string `gorm:"size:36;uniqueIndex:idx_outbox_entity"`
    Op       string `gorm:"size:10"`
    Data     []byte `gorm:"type:blob"`
    Version  int
    // ...
}
```

That unique index on `(Model, EntityID)` is the whole policy. A user who edits the same inspection note nine times offline does not produce nine pushes. They produce one, holding the final state.

The merge rules are small enough to read in full:

```go
// enqueue applies squash semantics:
//   - new "create" + no entry            → INSERT create
//   - new "update" + existing create     → UPDATE the create's data
//   - new "update" + existing update     → UPDATE data
//   - new "delete" + existing create     → DELETE outbox row (never made it to server)
//   - new "delete" + existing update     → flip to delete, drop data
//   - new "delete" + no entry            → INSERT delete
```

Two of those are worth stopping on.

**`update` on top of an existing `create` keeps the op as `create`.** The server has never heard of this row. Sending it an update for a row it does not have would 404. The change of state is real but the *operation* is still "this row is new to you".

**`delete` on top of a `create` deletes the outbox row entirely.** The user made something offline and then changed their mind. The server never knew. There is nothing to tell it. Sending a create followed by a delete would be two round trips to arrive at the state of having done nothing.

This is `git rebase -i` with every commit on a file squashed into one. You lose the intermediate history, and for a sync protocol that is a feature: nobody wants the server's audit log to record eight keystrokes-worth of saves from a train.

## Pull before push, always

```go
// Sync runs Pull → Push. Pull first so the user pushes against the
// freshest server state and conflict surface area is minimized.
```

Pull first and some of your conflicts evaporate before they exist. If a colleague edited row 7 and you did not touch row 7, pulling takes their change quietly. If you *did* touch row 7, you now know you are in conflict before you push, which means the conflict can be presented while the user is still in the app, rather than as a failed push.

Push first and you are writing to a version of the server you have not looked at since yesterday. Every conflict that was going to happen still happens, and you find out about all of them afterwards.

Fetch before you merge. It is the same instinct.

## Pull: the cursor is the hard part

Pull asks "what changed since I last looked", which needs a bookmark. The obvious bookmark is a timestamp, and the obvious bookmark is wrong.

```go
// Keyset on (updated_at, id), so rows sharing a timestamp are not lost at a page
// boundary, as they were with a cursor of the time alone.
func (s *SyncService) Page(ctx context.Context, model, dest interface{}, since time.Time, afterID string, limit int) error {
    q := s.db(ctx).Unscoped().Model(model)
    if !since.IsZero() {
        if afterID == "" {
            q = q.Where("updated_at > ?", since)
        } else {
            q = q.Where("updated_at > ? OR (updated_at = ? AND id > ?)", since, since, afterID)
        }
    }
    return q.Order("updated_at asc, id asc").Limit(limit).Find(dest).Error
}
```

Here is the bug that clause exists to prevent, and it is a data-loss bug that looks like nothing.

Say 600 rows share a millisecond, because a batch import wrote them. The page size is 500. With a cursor of time alone, you send 500 rows and the cursor becomes that timestamp. Next page asks for `updated_at > T`. **The other 100 rows are at exactly `T`.** They are not greater than `T`. They are skipped, silently, forever. The client's mirror is missing a hundred rows and nothing anywhere reports an error.

The fix is to make the cursor a **compound key**: the timestamp *and* the last ID seen. `(updated_at = T AND id > lastID)` picks up the stragglers in the same millisecond, then `updated_at > T` moves on.

The wire format is unglamorous and good:

```go
func syncCursor(row interface{}) (string, bool) {
    // ...
    return t.Format(time.RFC3339Nano) + "~" + fmt.Sprint(id.Interface()), true
}
```

`2026-10-22T09:14:03.123456789Z~01a1...`. Nanosecond precision because millisecond precision creates more ties, and ties are what the ID half is there to break.

There is a second reason for keyset over `OFFSET`: rows are moving. A row edited during your paging gets a new `updated_at` and jumps to the end of the ordering. With `OFFSET 500` the rows behind it shift up by one and you skip one. Keyset pagination reads a stable window because it asks "after this key" rather than "skip this many".

## Deletes need a body

If a row is deleted on the server, the next pull returns everything that changed. A hard-deleted row does not appear in that list. It does not appear in *any* list. There is no way for the client to distinguish "deleted" from "unchanged", so the client keeps it. Forever.

So deletes must leave something behind to travel. A **tombstone**:

```go
// Delete soft-deletes a row, which is what leaves the tombstone a pull sends to
// every other device.
func (s *SyncService) Delete(ctx context.Context, obj interface{}, id string) error {
    return s.db(ctx).Delete(obj, "id = ?", id).Error
}
```

and pull reads **unscoped**, which is the line that makes it work:

```go
// Unscoped, because the soft-deleted rows are the tombstones: a client that never
// hears about a deletion keeps the row forever.
q := s.db(ctx).Unscoped().Model(model)
```

Normally you work hard to make soft-deleted rows invisible. Here you work hard to make sure they are visible to exactly one reader.

### The bug under the tombstone

This one is my favourite thing in the package, because it is subtle, it is a performance bug and a correctness bug at once, and it took a production failure on MySQL to find:

```go
// Sync pull pages on updated_at. A soft delete used to set only deleted_at, so
// pull ordered by the later of the two: an expression no index can serve, which
// sorted the whole table on every pull, and one MySQL has no function for, so
// pull failed there outright.
```

GORM's soft delete sets `deleted_at` and leaves `updated_at` alone. So a deleted row's "change time" was `deleted_at`, while an edited row's was `updated_at`, and pull had to order by whichever was later. `GREATEST(updated_at, deleted_at)` is a computed expression: **no index can serve it**, so every pull sorted the entire table. And MySQL's `MAX()` is an aggregate, not a two-argument function, so the query did not merely run slowly there. It did not run.

The fix is to make the two agree at the source, with a GORM callback that writes `updated_at` alongside `deleted_at` in the same statement:

```go
func touchUpdatedAt(tx *gorm.DB) {
    // ...
    set := stmt.Clauses["SET"]
    set.Builder = func(c clause.Clause, b clause.Builder) {
        assignments, _ := c.Expression.(clause.Set)
        if len(assignments) == 1 && assignments[0].Column.Name == deleted.DBName {
            assignments = append(assignments, clause.Assignment{
                Column: clause.Column{Name: updated.DBName},
                Value:  assignments[0].Value,
            })
        }
        b.WriteString("SET ")
        assignments.Build(b)
    }
    stmt.Clauses["SET"] = set
}
```

Now one indexed column carries edits and deletes alike, and pull is an index scan on every engine.

And because projects already had rows soft-deleted before this existed, whose deletions would never travel:

```go
// BackfillUpdatedAt moves updated_at up to deleted_at on rows soft-deleted
// before Install existed, so pull carries those deletes too. It changes nothing
// on a second run; grit migrate runs it.
```

A fix that only works for data created after the fix is half a fix.

## Push: where two branches meet

```go
type PushChange struct {
    Op      string                 `json:"op"`
    Model   string                 `json:"model"`
    ID      string                 `json:"id"`
    Version int                    `json:"version"`
    Data    map[string]interface{} `json:"data"`
}
```

`Version` is the whole protocol. It is the version the **client believed the server had** when it made the edit. It is the client saying: *I am changing version 7. If you are not on version 7, we have a problem.*

That is optimistic locking, the same mechanism as [If-Match on an update route](/blog/optimistic-locking-in-go-with-if-match), except the gap between read and write is not a few seconds of a form being filled in. It is forty minutes in a car park.

The server compares:

```go
serverVersion := getIntField(current, "Version")
if serverVersion != ch.Version {
    // The versions disagree. What happens next is the resource's
    // declared policy, decided here rather than in the client, because
    // a rule an old build can ignore is not a rule.
```

Read that comment twice. **Conflict policy is decided on the server.** Put it in the client and the policy is whatever version of the app is installed, which on desktop and mobile means a policy from eighteen months ago is still running somewhere. A rule an old build can ignore is not a rule.

### Each change stands alone

```go
// Push handles POST /api/sync/push. Each change is applied
// independently: one conflict does not abort the rest of the batch.
```

Forty changes, one conflicted: thirty-nine land. The alternative, a batch transaction that rolls back, means one disputed record blocks an entire day of unrelated work, and the user cannot clear it without resolving the conflict, which they may not be able to do from a basement.

The results come back positionally, one per change, so the client can mark thirty-nine outbox rows done and flag one.

### And the batch is bounded

```go
// MaxPushChanges is the most changes one push may carry. The sync clients Grit
// ships send their outbox in pushes of this size. An unbounded push of 20,000
// changes was one request running 60,000 queries, which the client gave up on
// and sent again.
const MaxPushChanges = 500
```

Note the failure mode: not an error, a **timeout**, which the client retried, producing the same enormous request again. An unbounded batch is a request that gets slower as the user's offline period gets longer, which is exactly backwards. The person who was offline longest is the person who most needs the sync to finish.

### The N+1 that would have been

A naive push loads the current row per change. Five hundred changes, five hundred SELECTs. So the handler loads them in one go first, grouped by model:

```go
func (h *SyncHandler) loadCurrent(c *gin.Context, changes []PushChange) currentRows {
    ids := map[string][]string{}
    for _, ch := range changes {
        if ch.Op == "update" || ch.Op == "delete" {
            ids[ch.Model] = append(ids[ch.Model], ch.ID)
        }
    }
    // ... one ByIDs query per model, into a map keyed by ID ...
}
```

Five hundred changes across three models is three queries instead of five hundred. This is the [N+1 that the profiler cannot see](/blog/finding-the-line-that-is-burning-your-cpu), on the one endpoint where the batch size is set by how long somebody was in a tunnel.

## Three conflict policies, and the one that must be the default

```go
const (
    // ConflictManual parks the change and hands both versions back so a human
    // decides. The default, because silently discarding somebody's work should
    // be something you opt into.
    ConflictManual = "manual"
    // ConflictServerWins discards the client's change and returns the server
    // row. For records the server is authoritative about: stock levels,
    // prices, anything a back office owns.
    ConflictServerWins = "server_wins"
    // ConflictClientWins applies the client's change over the server's. For
    // records whose author is the only person who edits them, where the
    // version check is protecting nothing.
    ConflictClientWins = "client_wins"
)
```

The defaulting is the argument. **Manual is the default because every automatic policy destroys somebody's work**, and that should be a choice a person made for a specific table with a reason, not something they inherited.

The reasons are per-resource and they are real. Stock levels: `server_wins`, because the warehouse system is authoritative and a device's forty-minute-old count is simply wrong. A personal draft only its author edits: `client_wins`, because the version check is protecting nothing. A shared inspection report two engineers both edit: `manual`, because there is no rule, only a judgement, and the system's job is to find a human to make it.

Notice `server_wins` is reported **distinctly from a conflict**:

```go
return PushResult{
    OK:            false,
    Code:          "SERVER_WINS",
    Message:       fmt.Sprintf("server v%d kept over client v%d", serverVersion, ch.Version),
    ServerVersion: serverVersion,
    ServerData:    current,
}
```

Not `VERSION_CONFLICT`. The client needs to tell "you lost, take this and move on" apart from "stop and ask a human", because one of those clears an outbox row automatically and the other must not.

And in every case, including the automatic ones, the server sends back `ServerData`. The client cannot build a merge dialog, or even an honest "your change was replaced by this" message, without the other version in hand.

### The client parks it

```go
// Conflicts surface as outbox rows with HasConflict=true so the UI
// can present a per-field merge dialog. ResolveConflict(...) writes
// the merged data back into the outbox and clears the flag.
```

The conflicted change stays in the outbox with both versions attached. It is not lost, not auto-applied, and not blocking anything else. The user resolves it when they have time, the merged result goes back into the outbox, and the next push sends it as a normal change.

That is a merge conflict, left in the working tree until somebody deals with it.

## Policies, and the fields that never leave the device

```go
type Policy struct {
    Mode     string `json:"mode"`
    Conflict string `json:"conflict"`

    // Fields, when set, is the allowlist sent to clients.
    Fields []string `json:"fields,omitempty"`

    // LocalOnly names fields that exist only on the device. They are stripped
    // from anything a client pushes, so a draft note or a scratch flag can
    // live in the mirror without ever reaching the server.
    LocalOnly []string `json:"local_only,omitempty"`

    // MaxOfflineAge is how long a client may go without a successful sync
    // before its data should be treated as stale. Zero means no limit.
    MaxOfflineAge time.Duration `json:"-"`
}
```

`Fields` is an allowlist, because a device that is lost or stolen contains whatever you mirrored onto it. Not every column needs to travel. Salary, internal notes, another customer's details: a table can be syncable without being syncable in full. Identity and bookkeeping columns are always included regardless, because without `id`, `version` and `updated_at` the client cannot version a row or know it was deleted.

`LocalOnly` is the opposite direction, and it is enforced where enforcement means something:

```go
// Enforced before the payload reaches a decoder: a field declared
// local_only is a promise it never leaves the device, and a promise kept
// only by well-behaved clients is not one.
ch.Data = policy.StripLocalOnly(ch.Data)
```

The server strips those fields on arrival. A client could send them, through a bug or a modified build, and the server declines to hear them.

`MaxOfflineAge` is the honest one:

```go
// This is advisory by design. The server cannot enforce it, because a
// client that has not synced is by definition not talking to the server.
// What it buys is a client that can say so instead of showing
// three-day-old stock levels as though they were current.
```

You cannot enforce a rule on a device that is not listening. What you can do is give the device enough information to be honest with its user, and a banner reading "last synced 3 days ago" is the entire feature. The alternative is three-day-old stock levels presented in the same typeface as live ones.

Also note what `online_only` does on pull:

```go
if policy.Mode == sync.ModeOnlineOnly {
    // Answering with an empty page would look like "nothing has changed",
    // and the client would mirror an empty table forever.
    respond.Fail(c, respond.CodeNotSyncable, model+" is registered online_only and is not mirrored")
}
```

An empty success is indistinguishable from "up to date". Say no properly.

## Two more things the server will not take from you

```go
// Fields the server owns: the row's id is the change's id, the version is the
// server's, and the timestamps are never the client's to set.
ch.Data = stripServerFields(ch.Data)
```

A device's clock is wrong. Not might be: is, by seconds at least, and by hours if someone changed it to get a free trial back. If the client could set `updated_at`, it could write a row that sorts before the cursor of every other device, and that row would never be pulled by anybody. **Sync ordering has to come from one clock**, and it is the server's.

And on save:

```go
// Associations are omitted, because a sync payload carries columns and not
// relations, and GORM would otherwise upsert whatever an association field
// happens to hold. created_at is omitted for the same reason: a client that
// round-trips a row should not be able to rewrite when it was made.
return s.db(ctx).Omit(clause.Associations, "CreatedAt").Save(obj).Error
```

A row round-trips through a client that deserialises it into a struct with an empty `Product` relation. Save it without `Omit(clause.Associations)` and GORM helpfully tries to upsert that empty relation.

## What you actually have to build

Strip away the Go and the list is short, and it is the same list whatever your stack:

1. **A local mirror every read comes from.** Not a cache in front of the network. The only source of truth the UI knows about.
2. **A durable outbox, keyed by record so edits squash.** On disk, not in memory.
3. **A compound cursor: `(timestamp, id)`.** A timestamp alone loses rows that share a millisecond, silently, forever.
4. **Tombstones.** Deletes need a body to travel, and your pull must read unscoped to carry them.
5. **A version per row, sent with every change**, so the server can tell a clean apply from a collision.
6. **A conflict policy chosen on the server, defaulting to asking a human.**
7. **Server-owned fields stripped on arrival**, because the device's clock is wrong and its ordering is not yours.
8. **A bounded batch**, or the user who was offline longest is the one whose sync never completes.

The engineer comes out of the car park. Her app pulls first, taking her colleague's eleven unrelated edits without a word. It pushes her eleven changes, squashed down from the thirty times she actually pressed Save. Ten apply. One comes back `VERSION_CONFLICT` with both versions attached, and she gets a dialog showing her text and his, and picks.

That last one is the only moment the software could not decide on its own, which is exactly the number of decisions you want to hand a person.

---

*Grit is a full-stack Go and Next.js framework. `grit new my-app --desktop` scaffolds the sync engine, the outbox and the conflict UI in this post, with `/api/sync/push` and `/api/sync/pull` on the server side. [Get started](/docs/start).*
