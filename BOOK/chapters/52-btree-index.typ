#import "../template.typ": *

= Build your own on-disk index: pages and a B-tree

The key-value store in chapter fifty one has a limit it does not advertise: on
startup it replays the whole log into a map, so the index lives in memory and
the startup time grows with the data. At a few hundred megabytes that is fine.
At fifty gigabytes it is not a database, it is a long wait.

The answer is an index that lives on disk and is read a page at a time. That is
what every real database does, and the structure is nearly always a B-tree.

== Why a tree and not a hash

A hash map finds one key in constant time, and that is all it does. A database
needs more:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Range scan], [`WHERE created_at > ?` reads a contiguous run of keys],
  [Ordering], [`ORDER BY` with no sort step],
  [Prefix], [`LIKE 'inv-2026%'` is a range],
  [Locality], [Nearby keys are on the same page, so one read serves many],
)

A B-tree gives all four because it keeps keys in order, and it is shaped for
disks rather than for memory: wide and shallow, so finding a key costs very few
reads.

#bridge[
  The arithmetic is the whole argument. With a 4 KB page and keys around 30
  bytes, one page holds roughly a hundred entries.

  One hundred to the power of three is a million. So a million keys is *three
  page reads deep*, and a hundred million is four. A binary tree with the same
  million keys is twenty levels, which is twenty seeks.

  Depth is the number of disk reads, and fanout is how you buy it down.
]

== Pages, the unit of everything

```go
const PageSize = 4096

type PageID uint32

type Page struct {
    ID   PageID
    Data [PageSize]byte
}
```

Everything is read and written in whole pages, at offset `int64(id) *
PageSize`. The size matches what the operating system and the disk move anyway,
so a partial read costs the same as a full one.

A page starts with a header, then entries:

```
offset  bytes  meaning
0       1      node type: 1 leaf, 2 internal
1       2      number of keys
3       4      next leaf (leaves only), 0 for none
7       ...    entries
```

Fixed offsets, little endian, no reflection, no JSON. The point of a page
format is that reading it is a few `binary.LittleEndian.Uint16` calls, not a
decoder.

== Leaves hold data, internal nodes hold directions

This is a B+tree, which is what databases actually use:

```
                 [ 50 | 120 ]              internal: separators
                /     |      \
        [10,30,45] [55,90] [130,180]       leaves: keys and values
             |         |         |
             +-------->+-------->+          leaves are linked
```

Two properties follow, and both matter.

*All values are in the leaves.* An internal node is only directions, so it fits
more separators, so the fanout is higher and the tree is shallower.

*Leaves are linked.* A range scan descends once and then walks sideways,
reading consecutive pages without touching the tree again.

== Searching

```go
func (t *Tree) Get(key []byte) ([]byte, error) {
    pid := t.root
    for {
        node, err := t.read(pid)
        if err != nil {
            return nil, err
        }
        if node.isLeaf() {
            i, found := node.search(key)
            if !found {
                return nil, ErrNotFound
            }
            return node.value(i), nil
        }
        pid = node.child(node.childIndex(key))
    }
}
```

`search` is a binary search *inside the page*, because the keys in a page are
sorted:

```go
func (n *node) search(key []byte) (int, bool) {
    lo, hi := 0, n.count()
    for lo < hi {
        mid := (lo + hi) / 2
        switch bytes.Compare(n.key(mid), key) {
        case 0:
            return mid, true
        case -1:
            lo = mid + 1
        default:
            hi = mid
        }
    }
    return lo, false
}
```

Note the loop in `Get` rather than recursion. Descending a tree iteratively
costs nothing extra and removes any question about stack depth.

== Inserting, and the split

Inserting is "find the leaf, put it in". The interesting half is what happens
when the leaf is full.

```
leaf is full:
  [10, 30, 45, 55]  + 60

split:
  [10, 30]  [45, 55, 60]
       \       /
   parent gains separator 45
```

The split propagates: if the parent is now full, it splits too. If the *root*
splits, a new root is made above it and *the tree gets one level deeper*. That
is the only way a B-tree grows in height, which is why its depth is uniform and
why every lookup costs the same.

```go
func (t *Tree) insert(pid PageID, key, val []byte) (
    split *splitResult, err error) {

    node, err := t.read(pid)
    if err != nil {
        return nil, err
    }
    if node.isLeaf() {
        node.put(key, val)
        if node.fits() {
            return nil, t.write(node)
        }
        return t.splitLeaf(node)
    }

    i := node.childIndex(key)
    sp, err := t.insert(node.child(i), key, val)
    if err != nil || sp == nil {
        return nil, err
    }
    node.insertSeparator(i, sp.key, sp.right)
    if node.fits() {
        return nil, t.write(node)
    }
    return t.splitInternal(node)
}
```

#gotcha[
  `fits()` rather than "count < max". Keys and values are variable length, so a
  page is full when the bytes no longer fit, not when some count is reached.

  Writing this as a count works perfectly until a long key arrives, and then
  corrupts a page by overwriting the one after it. Size the check in bytes and
  test it with a key that is deliberately enormous.
]

== Deleting, honestly

Textbook B-trees merge underfull nodes and rebalance. It is fiddly, it is where
the subtle bugs live, and many real systems do something simpler:

*Mark the entry deleted and leave the page alone.* Space is reclaimed when the
page is rewritten or during a compaction pass.

The cost is pages that are emptier than they could be, which costs reads. The
benefit is a much smaller surface for corruption. Start here, and implement
merging only when you can measure that it matters.

== The free list

Pages become free when a compaction empties them. Without a list of them, the
file only grows:

```go
type freeList struct {
    ids []PageID
}

func (t *Tree) allocate() PageID {
    if id, ok := t.free.pop(); ok {
        return id
    }
    t.nextPage++
    return t.nextPage - 1
}
```

The free list lives in a page of its own, usually near the start, and it is the
second thing to write after the root pointer.

== Durability, and the page that was half written

A page write is not atomic. The power goes out, and 4096 bytes are partly old
and partly new: a *torn page*. A B-tree with a torn internal node is not a tree
with one bad row, it is a tree that cannot be traversed.

Two defences, and real databases pick one:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Write-ahead log], [Write the change to a log and `fsync` it, then update the
    page. Recovery replays the log. This is chapter fifty one's machinery,
    reused],
  [Copy on write], [Never overwrite a live page. Write a new one, then flip the
    root pointer with one small atomic write],
)

Copy on write is the one to build first here, because it fits the structure:

```go
// Every modified page is written to a NEW page id. The parents
// are rewritten to point at it, all the way to the root. Then
// one write updates the header's root pointer, and only that
// write has to be atomic.
```

A header small enough to live in a single sector gets atomicity from the
hardware, which is why database file headers are tiny.

#gotcha[
  `fsync` is the whole of durability, and forgetting it makes everything look
  fast and work perfectly until the first real power cut.

  Order matters too: write the new pages, `fsync`, *then* write the root
  pointer, `fsync` again. Flipping the root before the pages it points at are
  durable gives you a root pointing at garbage, which is worse than losing the
  write.
]

== A buffer pool, because reading the same page is silly

```go
type pool struct {
    mu     sync.Mutex
    pages  map[PageID]*Page
    order  *list.List            // LRU
    max    int
}
```

The top two levels of the tree are read on every single lookup, so caching
turns three reads into one. A few thousand pages of cache covers the internal
nodes of a surprisingly large tree, and that is why databases want memory.

== Testing a data structure you cannot see

Three techniques, in increasing order of how much they find.

*Differential testing against a sorted map.* Any sequence of operations must
produce the same answers as `map` plus a sort:

```go
func TestMatchesReference(t *testing.T) {
    tree := openTemp(t)
    ref := map[string]string{}

    for i := 0; i < 50_000; i++ {
        k, v := randKey(), randValue()
        tree.Put([]byte(k), []byte(v))
        ref[k] = v
    }
    for k, want := range ref {
        got, err := tree.Get([]byte(k))
        if err != nil || string(got) != want {
            t.Fatalf("Get(%q) = %q, %v; want %q", k, got, err, want)
        }
    }
}
```

*An invariant check.* Walk the tree and assert: keys sorted within each page,
every key in a child within its parent's separators, all leaves at the same
depth, pages at least half full except the root. Run it after every thousand
inserts in a test.

*A crash test.* Kill the process mid-write, reopen, and check the invariants
plus every acknowledged key. This is the one that finds the `fsync` ordering
bug, and nothing else will.

#project(title: "Project 48: the index, measured")[
  In `part6-from-scratch/52-btree`, build it and then find out what you built.

  *One: pages.* A pager with `read`, `write`, `allocate`, a free list and a
  header page. Test that reopening finds the root.

  *Two: the tree.* `Get`, `Put`, `Delete` with tombstones, splits for leaves
  and internal nodes, and a root split that deepens the tree.

  *Three: the invariant checker*, run after every thousand operations in a
  test.

  *Four: the differential test* against a reference map, with fifty thousand
  random keys.

  *Five: ranges.* `Scan(from, to)` using the leaf links. Prove with a timing
  that scanning ten thousand consecutive keys is far cheaper than ten thousand
  `Get` calls.

  *Six: the depth.* Insert one million keys and report the tree height and the
  fanout you achieved. Compare it with the arithmetic at the top of this
  chapter.

  *Seven: durability.* Copy on write with the root flip last. Then kill the
  process in a loop during writes, reopen each time, and assert the invariants
  and every acknowledged key.

  *Eight: the cache.* Add an LRU buffer pool and measure lookups per second
  with it at 0, 100 and 10000 pages. Put the three numbers in `RESULTS.md`.
]

#tryit[
  Open your tree file in a hex viewer after inserting a few keys.

  ```bash
  xxd -l 256 index.db
  ```

  Finding your header, your node type byte and your first key in the bytes is
  the moment a file format stops being abstract. Every database file you will
  ever meet is this, with more fields.
]

#assignment(name: "B-Tree Index - Assignment One")[
  Build the index and report it like an engineer rather than a tutorial.

  *One.* The pager, the B+tree with splits, range scans over linked leaves, and
  the free list.

  *Two.* The invariant checker and the differential test, both in CI.

  *Three.* `RESULTS.md` with: tree height at 1 thousand, 100 thousand and 1
  million keys; lookups per second cold and warm; the scan comparison; and the
  file size against the raw data size, which is your space amplification.

  *Four.* The crash test, run at least fifty times, with the number of runs and
  any failure you found.

  *Five.* One honest paragraph: what does SQLite do that your index does not?
  Name three things. Knowing the list is more valuable than closing it.
]
