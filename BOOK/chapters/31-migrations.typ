#import "../template.typ": *

= Models, migrations and seeding at scale

D10 from chapter twenty six was "changing the schema is the scary part: the
migration that half ran, with no way back". This chapter is how that is
answered, and it is more interesting than it sounds, because the usual answer
is a lie.

== The model

A generated model is a struct with tags:

```go
type Invoice struct {
    ID        uint         `gorm:"primarykey" json:"id"`
    Number    string       `gorm:"size:50;uniqueIndex" json:"number"`
    Total     money.Amount `gorm:"not null" json:"total"`
    Due       jsontime.Date `json:"due"`
    ClientID  uint         `gorm:"index" json:"client_id"`
    Client    *Client      `json:"client,omitempty"`
    Version   int          `gorm:"default:1" json:"version"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

Three of those fields are the framework's rather than yours.

`Version` is optimistic locking, and chapter thirty three explains the 409 it
produces. `DeletedAt` is soft delete: a `DELETE` sets the timestamp and every
query excludes it automatically, which is how the Trash page works.
`CreatedAt` and `UpdatedAt` are maintained by GORM.

#gotcha[
  Soft delete has a consequence people meet late: a `uniqueIndex` on a column
  still applies to deleted rows. Delete a user with the email `ada@example.com`
  and you cannot create another with that email, because the old row is still
  there.

  The usual fix is a composite unique index including `deleted_at`, so the same
  value can exist once among live rows and many times among deleted ones. It is
  worth knowing before a user reports it.
]

== Migrations, and the honest problem

```bash
grit migrate
```

It reads your models and makes the database match. No migration files to write.

That convenience has a well-known cost and this is where most frameworks are
quiet: *GORM's AutoMigrate has no down scripts*. It can add a column. It cannot
tell you what it added last Tuesday, and it cannot undo it.

So the answer here is not a pretend rollback. It is a record:

```bash
grit migrate            # apply, and record what this run changed
grit migrate --status   # what previous runs did
grit migrate --down     # undo the last run, using that record
```

Each run writes down what it actually added. `--down` reads that record and
drops exactly those things. It is not a general rollback, it is an undo of
*recorded* runs, and the distinction matters: a change made by someone typing
SQL into a console was never recorded, so it cannot be undone this way.

#bridge[
  Prisma and Rails generate a migration file per change, with an explicit `up`
  and `down`, which is more ceremony and gives you a reviewable artifact. This
  trades the file for automatic detection plus a recorded history.

  Neither is free. The file-based approach costs a step per change and gives
  you a precise diff in code review. This approach costs you reviewability and
  gives you speed. Know which you are using, and never assume a framework that
  offers `--down` can undo something it did not do.
]

== The refusal that protects you

```bash
grit migrate --down
```

Dropping a column destroys data. The command asks, and a non-interactive run
needs `--yes`, deliberately, so a CI pipeline cannot quietly delete a column
because somebody renamed a struct field.

A rename is the case to understand. Changing `Name` to `FullName` in a struct
is not a rename to AutoMigrate: it is an unknown column appearing and a known
one no longer being mentioned. The safe sequence is to add the new column,
backfill it, then remove the old one in a later change, which is the same
expand-and-contract dance every schema tool requires.

== Seeding, and why it is not optional

D9 was "nothing to develop against".

```bash
grit seed                        # every seeder
grit seed Invoice --count 50000  # one, at scale
```

The generator writes a factory per resource that produces realistic values for
each field type: a real-looking name for `string`, a valid address for `email`,
an amount in range for `money`.

At ten rows everything works. The sort is instant, pagination has one page, and
an N+1 costs ten queries instead of ten thousand. Every list-page bug you have
is invisible until the data arrives, which is why the habit is to seed *before*
you look at a screen rather than after a user complains.

== A million rows, with the heap flat

The seeder is the pipeline from chapter fourteen, in production:

```
build rows in chunks  ->  bounded channel  ->  batch writers  ->  database
```

Measured on a laptop: *999,800 rows in 59 seconds, with the heap flat at 10 MB
the whole time*.

The number that matters is the third one. Sixteen thousand rows a second is a
fine figure and it is not the point. The point is that memory is a property of
the *pipeline* rather than of how many rows you asked for: the channel holds a
fixed number of batches in flight, so ten thousand rows and ten million look
the same to the heap.

#gotcha[
  The naive versions both fail, in opposite directions. One row per `INSERT`
  keeps memory perfect and takes an afternoon, because you pay query parsing, a
  transaction and a round trip a million times. Building all million rows into
  a slice first is fast until the process is holding a million structs plus the
  driver's copy, and then it swaps or is killed.

  *Batching the writes is the win, not the goroutines.* Going from one insert
  per row to one insert per batch moves the number by an order of magnitude.
  Parallelism on top is worth perhaps five times more. Parallelism without
  batching is worth almost nothing, because a million round trips is a million
  round trips however many goroutines are waiting on them.
]

The seeder is also *resumable*: it counts what is in the table and inserts only
the difference. Run it twice and the second run does nothing. Kill it halfway
through a million rows and run it again, and it continues.

== Backups, and the half that counts

```bash
grit backup           # write an archive
grit restore <file>   # replay one
```

`cmd/restore` existing beside `cmd/backup` is the point. B5 asked "is any of
this backed up", and the honest reading of that question is "has anyone ever
restored one".

A backup you have never restored is a rumour. The archive carries a row count
per table, so a restore can be checked rather than assumed, and the restore is
a command rather than a runbook page nobody has followed.

#project(title: "Project 28: break a migration on purpose")[
  The point of this project is to experience the failure modes while they are
  cheap.

  ```bash
  grit new migrations-lab --api
  cd migrations-lab
  grit generate resource Product \
    --fields "name:string, price:money, sku:string"
  grit migrate
  grit seed Product --count 1000
  ```

  *One: watch a run get recorded.*

  ```bash
  grit migrate --status
  ```

  Note what the last run says it changed.

  *Two: add a column and undo it.*

  ```bash
  grit generate field Product weight:float
  grit migrate
  grit migrate --status      # a new run, with the new column
  grit migrate --down        # and gone again
  ```

  Check in `/studio` that the column really left.

  *Three: the rename trap.* Edit `models/product.go` by hand and change `SKU`
  to `Code`. Run `grit migrate --status` after migrating and look carefully: a
  column was added, and the old one is still there with the data in it. Nothing
  was renamed, because nothing could know that was your intention.

  Now do it properly: add the new field, write a tiny `cmd/` program that
  copies the values across, and only then remove the old one. That is
  expand-and-contract, and doing it once by hand is worth more than reading
  about it.

  *Four: the scale that matters.*

  ```bash
  time grit seed Product --count 200000
  ```

  Watch the progress line report rows per second and the heap. Then open the
  admin Products page, sort by price and search. If anything is slow, open
  `/pulse` and look at the query count: this is the moment that a missing index
  announces itself, and it is the entire argument for seeding before you demo.

  *Five: prove the restore.*

  ```bash
  grit backup
  grit seed Product --count 10          # change the data
  grit restore <the archive>
  ```

  Then count the rows. Restoring once, in a project where it does not matter,
  is how you find out your backup process works before the night it has to.
]

#tryit[
  Add a `uniqueIndex` to the product SKU, create a product, soft delete it, and
  try to create another with the same SKU. Read the error, then fix it with a
  composite index including `deleted_at` and confirm both cases now work: one
  live row with that SKU, and any number of deleted ones.

  That is a bug most applications ship, discovered by a user who deleted
  something and could not re-add it.
]

#assignment(name: "Migrations - Assignment One")[
  Make the data real, then prove you could get it back.

  *One.* Add three fields to an existing model with `grit generate field`, then
  migrate. Record what the migration refused to do and why.

  *Two.* Seed a hundred thousand rows and watch the memory while it runs:

  ```bash
  grit seed Reading --count 100000
  ```

  Put the peak in `DATA.md`.

  *Three.* Now the half that counts. Take a backup, drop a table, restore it,
  and confirm the rows are back. Write the exact commands in `DATA.md`.

  *Four.* Write one paragraph on what you would still lose if the machine died
  between the backup and the drop. Every honest backup story has that
  paragraph.]
