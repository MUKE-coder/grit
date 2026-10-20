#import "../template.typ": *

= Field types: what each one decides

Thirty one field types ship. This chapter is not a reference, it is an argument
about why a field type is the right place to put a decision.

== The argument

When you write `total:money`, you are not choosing a storage format. You are
making five decisions at once:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Database], [an integer column in minor units, never a float],
  [Go], [a `money.Amount`, with arithmetic that cannot drift],
  [Validation], [a Zod rule that accepts `12.34` and rejects `12.345`],
  [Form], [a money input with a currency prefix and two decimal places],
  [Table], [right aligned, formatted for the locale],
)

Say `float` instead and all five change, each in a direction that is slightly
wrong, and the one that costs you is the first. `0.1 + 0.2` is not `0.3` in any
language with IEEE floats, and an invoice total that is out by a penny is found
by an accountant rather than by a test.

#bridge[
  You already wrote this fix by hand in chapter nineteen, implementing
  `MarshalJSON` on a `Money int64`. That was the right instinct. The difference
  here is that the correct decision is the *default*, applied in five places
  consistently, rather than something each developer remembers or does not.
]

== The full list

```
string  text  richtext  int  uint  float  bool  toggle  select  radio
check  datetime  date  money  slug  belongs_to  one_to_one  many_to_many
string_array  file  files  email  url  domain  tel  country  color
percent  rating  time  json
```

Most are unsurprising. The ones worth knowing properly are below.

== Text in three sizes

`string`, `text` and `richtext` differ in the column type and the input:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`string`], [A varchar, a single line input. Names, codes, references],
  [`text`], [A text column, a textarea. Descriptions, notes],
  [`richtext`], [A text column holding HTML, and an editor in the form],
)

`richtext` is the one with a hidden cost: it stores HTML that a user supplied,
so it has to be sanitised on the way in or you have shipped stored cross-site
scripting. The generated service does that, which is the reason to use the
field type rather than a `text` plus an editor you wire yourself.

== Choices, in three shapes

`select`, `radio` and `check` all constrain a value to a set, and differ only
in the control:

```bash
grit generate resource Ticket \
  --fields "status:select:open=Open|closed=Closed|held=On hold"
```

The values before the `=` go in the database; the labels after it go on the
screen. This produces a constraint, a Go constant set, a Zod enum and a
dropdown, from one string.

#gotcha[
  Use `select` for one choice from many, `radio` when there are three or fewer
  and the user should see them all at once, and `check` for several choices at
  once, which stores a `string_array`.

  The trap is choosing `string` and validating in your head. A free `status`
  column collects `open`, `Open`, `OPEN` and `opened` within a year, and no
  report ever agrees again.
]

== Relationships

Three, and they are the field types that save the most work.

`belongs_to` is one word that produces a foreign key column with a constraint,
a preloaded struct field, a nested TypeScript type, and a searchable dropdown
in the form that shows the related record's *label column* rather than a UUID.

```bash
--fields "client:belongs_to"
```

```go
ClientID uint    `json:"client_id"`
Client   *Client `json:"client"`     // preloaded on read
```

`one_to_one` is the same with a uniqueness constraint. `many_to_many` creates
the join table and a multi-select picker.

#bridge[
  The dropdown showing a name rather than an ID is the detail worth noticing,
  because it is the canonical hand-written bug. It happens because the frontend
  has no way to know which column is the human label, so it falls back to the
  primary key. The generator knows, because `belongs_to` told it to look.
]

== Identity and files

`slug` derives a URL-safe string from another field and keeps it unique, which
is three decisions people usually get wrong separately: how to slugify, when to
regenerate, and what to do about collisions.

`file` and `files` store an upload reference rather than bytes in a column, and
wire the upload pipeline from chapter forty two: resize, re-encode, strip EXIF,
generate a thumbnail. Writing `photo:file` is how B8 gets answered.

== Validated strings

`email`, `url`, `domain`, `tel` and `country` are all strings with a format.
They look cosmetic and are not: each produces a matching rule in Go *and* in
Zod, so the browser and the API agree about what is valid.

#gotcha[
  The point is the agreement, not the validation. Validating only in the
  frontend is a user-experience feature, not a security control, because
  anybody can post directly to your API. Validating only in the backend means
  the user discovers the problem after a round trip.

  One field type produces both, from one definition, which is D3 again.
]

== Numbers with meaning

`int`, `uint` and `float` are the raw ones. `percent`, `rating` and `money` are
numbers that carry their own rules: a range, a step, a format. Prefer the
specific one, because the constraint travels with it.

`json` stores arbitrary structure. Use it for genuinely unstructured data such
as a webhook payload, and resist using it to avoid designing a schema, because
nothing can query it well and nothing validates it.

== Time, in three precisions

`date` is a calendar day with no time and no zone: a birthday, an invoice due
date. `datetime` is an instant, stored UTC. `time` is a clock time with no
date, for opening hours.

#gotcha[
  Using `datetime` for a birthday is a real bug, not a style preference. Store
  it as an instant and somebody in a different timezone sees the day before.
  The rule: *if the answer should be the same everywhere in the world, it is a
  `date`.*
]

== Choosing well, in one table

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Money], [`money`, never `float`],
  [A fixed set], [`select`, never `string`],
  [A calendar day], [`date`, never `datetime`],
  [A relationship], [`belongs_to`, never an `int` you join by hand],
  [An upload], [`file`, never a `string` holding a path],
  [An email], [`email`, so both sides validate it],
  [Unstructured], [`json`, and only when it really is],
)

Every row of that table is a decision that is cheap now and expensive in six
months, which is the whole case for field types being a type system rather than
a convenience.

#project(title: "Project 27: one resource with every interesting type")[
  Generate a resource that exercises the types that carry decisions, then look
  at what each produced.

  ```bash
  grit generate resource Property \
    --fields "title:string, slug:slug, description:richtext, \
  price:money, bedrooms:int, rating:rating, available_from:date, \
  viewing_time:time, status:select:draft=Draft|live=Live|sold=Sold, \
  contact_email:email, website:url, photos:files, owner:belongs_to"
  grit migrate
  grit seed Property --count 200
  grit start
  ```

  Then check each claim rather than taking it on trust.

  *The money column.* Open `models/property.go` and find `price`. What Go type,
  and what does the GORM tag say? Then open the admin form and type `12.345`
  into it. Then post `{"price": 12.345}` directly to the API with curl. Both
  should refuse, and that is the agreement this chapter is about.

  *The date against the time.* Find `available_from` and `viewing_time` in the
  generated TypeScript. They are different types, because one is a day and one
  is a clock reading.

  *The select.* Find the status constants in Go and the enum in the Zod schema.
  Then try posting `{"status": "pending"}` and read the error.

  *The relationship.* Open the Property form and find the owner dropdown. It
  shows names. Now look at `services/property.go` and find the preload that
  makes that possible without an N+1.

  *The slug.* Create a property called "Garden Flat", then another with the
  same name, and look at the two slugs.

  *The files.* Upload two photos and then look at what landed in storage: the
  original is not what is served.
]

#tryit[
  Generate the same resource twice, once with `price:money` and once with
  `price:float`, into throwaway projects. Diff the two generated models, Zod
  schemas and admin forms. Count how many files differ from that one word.

  Then, in the float version, set a price of `0.1`, add `0.2` to it in a Go
  test, and assert it equals `0.3`. Watch it fail. That failing test is the
  entire argument of this chapter in four lines.
]
