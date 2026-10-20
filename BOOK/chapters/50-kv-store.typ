#import "../template.typ": *

= Build your own key-value store: a write-ahead log and crash recovery

A map in memory is a key-value store until the process restarts. This chapter
builds one that survives a crash, which turns out to be a question about
exactly one system call.

Everything here uses only the standard library, and everything you need is in
Parts I to III.

== The problem, stated precisely

```go
m := map[string]string{}
m["user:1"] = "Ada"
// power cut
```

The data is gone. The obvious fix, writing the map to a file on every change,
is both slow and *still wrong*: a crash halfway through writing leaves a
corrupt file, and you have traded losing everything for losing everything in a
less obvious way.

The standard answer is a *write-ahead log*. Before changing memory, append the
change to a file. On startup, replay the file. Appending is fast because it is
sequential, and a partially written record at the end is detectable and
discardable.

#bridge[
  This is what Postgres, SQLite, Kafka and every serious database do. The WAL
  is the durable thing; the in-memory structure is a cache of it that can be
  rebuilt. Once you have written one, a great deal of database documentation
  stops being mysterious.
]

== The record format

A log is a sequence of records. Each needs a length, a checksum and a payload:

```
+--------+----------+--------+---------+-------+---------+
| crc32  | keyLen   | valLen |   op    |  key  |  value  |
| 4 byte | 4 byte   | 4 byte | 1 byte  |  ...  |   ...   |
+--------+----------+--------+---------+-------+---------+
```

The checksum is what makes a half-written record detectable. Without it a
truncated write looks like a valid record with garbage in it.

```go
package kv

import (
    "bufio"
    "encoding/binary"
    "errors"
    "fmt"
    "hash/crc32"
    "io"
    "os"
    "sync"
)

type op byte

const (
    opSet op = 1
    opDel op = 2
)

const headerSize = 4 + 4 + 4 + 1   // crc, keyLen, valLen, op
```

== Encoding a record

```go
func encode(o op, key, val []byte) []byte {
    buf := make([]byte, headerSize+len(key)+len(val))

    binary.BigEndian.PutUint32(buf[4:8], uint32(len(key)))
    binary.BigEndian.PutUint32(buf[8:12], uint32(len(val)))
    buf[12] = byte(o)
    copy(buf[headerSize:], key)
    copy(buf[headerSize+len(key):], val)

    // The checksum covers everything after itself.
    sum := crc32.ChecksumIEEE(buf[4:])
    binary.BigEndian.PutUint32(buf[0:4], sum)
    return buf
}
```

`encoding/binary` writes fixed-width integers in a defined byte order.
Big-endian is conventional for file formats; the only rule that matters is that
the reader and writer agree.

#gotcha[
  Do not use `fmt.Sprintf` or JSON for a log format. Both are variable-width,
  so you cannot know how many bytes to read without scanning for a delimiter,
  and a delimiter can appear inside a value. A length prefix removes the
  question entirely: read the header, then read exactly that many bytes.
]

== The store

```go
type Store struct {
    mu   sync.RWMutex
    data map[string][]byte

    log  *os.File
    w    *bufio.Writer
    path string
}

func Open(path string) (*Store, error) {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
    if err != nil {
        return nil, fmt.Errorf("opening log %s: %w", path, err)
    }

    s := &Store{
        data: map[string][]byte{},
        log:  f,
        w:    bufio.NewWriter(f),
        path: path,
    }

    // Rebuild memory from the log before accepting any writes.
    if err := s.replay(); err != nil {
        f.Close()
        return nil, fmt.Errorf("replaying log: %w", err)
    }
    return s, nil
}

func (s *Store) Close() error {
    if err := s.w.Flush(); err != nil {
        s.log.Close()
        return fmt.Errorf("flushing: %w", err)
    }
    return s.log.Close()
}
```

`O_APPEND` matters: every write goes to the end regardless of where the file
offset happens to be, which is what makes concurrent appends safe at the
operating system level.

== Replay, and tolerating a torn tail

```go
func (s *Store) replay() error {
    if _, err := s.log.Seek(0, io.SeekStart); err != nil {
        return err
    }
    r := bufio.NewReader(s.log)
    header := make([]byte, headerSize)
    offset := int64(0)

    for {
        if _, err := io.ReadFull(r, header); err != nil {
            if errors.Is(err, io.EOF) {
                break                      // clean end
            }
            if errors.Is(err, io.ErrUnexpectedEOF) {
                break                      // torn header: stop here
            }
            return err
        }

        sum := binary.BigEndian.Uint32(header[0:4])
        keyLen := binary.BigEndian.Uint32(header[4:8])
        valLen := binary.BigEndian.Uint32(header[8:12])
        o := op(header[12])

        body := make([]byte, keyLen+valLen)
        if _, err := io.ReadFull(r, body); err != nil {
            break                          // torn body: stop here
        }

        // Verify before trusting any of it.
        check := crc32.ChecksumIEEE(append(header[4:], body...))
        if check != sum {
            break                          // corrupt: the log ends here
        }

        key := string(body[:keyLen])
        switch o {
        case opSet:
            s.data[key] = body[keyLen:]
        case opDel:
            delete(s.data, key)
        }
        offset += int64(headerSize) + int64(len(body))
    }

    // Discard anything after the last good record, so the next append does
    // not sit behind garbage.
    if err := s.log.Truncate(offset); err != nil {
        return fmt.Errorf("truncating to %d: %w", offset, err)
    }
    _, err := s.log.Seek(0, io.SeekEnd)
    return err
}
```

The three `break`s are the whole crash-recovery story. A torn header, a torn
body and a failed checksum all mean the same thing: *the log ends at the last
good record*, and everything after it was never acknowledged to a caller, so
discarding it is correct.

#gotcha[
  `append(header[4:], body...)` is the slice trap from chapter six wearing a
  disguise. `header[4:]` has spare capacity, so `append` may write into
  `header` rather than allocating. It happens to be harmless here because
  `header` is re-read next iteration, but it is exactly the pattern that
  corrupts data elsewhere. The careful version:

  ```go
  h := crc32.NewIEEE()
  h.Write(header[4:])
  h.Write(body)
  check := h.Sum32()
  ```

  Which is also faster, because it allocates nothing.
]

== Writing

```go
func (s *Store) Set(key string, val []byte) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, err := s.w.Write(encode(opSet, []byte(key), val)); err != nil {
        return fmt.Errorf("appending set: %w", err)
    }
    if err := s.sync(); err != nil {
        return err
    }

    // Memory is updated only after the log is durable. The other order
    // means a reader can see a value that a crash will erase.
    s.data[key] = val
    return nil
}

func (s *Store) Delete(key string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, err := s.w.Write(encode(opDel, []byte(key), nil)); err != nil {
        return fmt.Errorf("appending delete: %w", err)
    }
    if err := s.sync(); err != nil {
        return err
    }
    delete(s.data, key)
    return nil
}

func (s *Store) Get(key string) ([]byte, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.data[key]
    return v, ok
}
```

The ordering in `Set` is the correctness argument, and it is worth saying out
loud: *log first, then memory*. If memory were updated first, a reader could
observe a value that the crash removes, and the store would have lied.

== The one system call that matters

```go
func (s *Store) sync() error {
    if err := s.w.Flush(); err != nil {       // bufio -> kernel
        return fmt.Errorf("flushing: %w", err)
    }
    if err := s.log.Sync(); err != nil {      // kernel -> disk
        return fmt.Errorf("syncing: %w", err)
    }
    return nil
}
```

Two steps, and people routinely do only the first.

`Flush` moves bytes out of Go's buffer into the operating system. At that point
the data survives your *process* crashing, because the kernel has it.

`Sync` is `fsync`, and it asks the kernel to put the bytes on the physical
device. Only after this does the data survive the *machine* losing power.

#gotcha[
  `fsync` is slow: on spinning disks milliseconds, on an SSD still far slower
  than anything else in this program. Calling it on every write makes a store
  that does maybe a few thousand writes a second.

  That is the durability dial, and every database exposes it. Postgres calls it
  `synchronous_commit`, Redis calls it `appendfsync`, and the options are
  always the same three:

  #table(
    columns: (auto, 1fr),
    stroke: 0.4pt + rgb("#e3e6ea"),
    inset: 5pt,
    [every write], [Slowest, loses nothing],
    [every N ms], [Fast, loses up to N milliseconds on a power cut],
    [never], [Fastest, the operating system decides, can lose a lot],
  )

  There is no correct answer, only a choice about what you can afford to lose.
  A store that does not tell you which one it made is the one to distrust.
]

== Compaction

The log grows forever. Setting one key a million times means a million records
and a one-entry map.

```go
// Compact writes a fresh log containing only the current state, then swaps
// it in atomically.
func (s *Store) Compact() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    tmp := s.path + ".compact"
    f, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
    if err != nil {
        return fmt.Errorf("creating %s: %w", tmp, err)
    }

    w := bufio.NewWriter(f)
    for k, v := range s.data {
        if _, err := w.Write(encode(opSet, []byte(k), v)); err != nil {
            f.Close()
            os.Remove(tmp)
            return fmt.Errorf("writing compacted record: %w", err)
        }
    }
    if err := w.Flush(); err != nil {
        f.Close()
        os.Remove(tmp)
        return err
    }
    if err := f.Sync(); err != nil {        // durable before the swap
        f.Close()
        os.Remove(tmp)
        return err
    }
    f.Close()

    // Rename is atomic on the same filesystem: the log is either the old
    // one or the new one, never a half-written mixture.
    if err := os.Rename(tmp, s.path); err != nil {
        os.Remove(tmp)
        return fmt.Errorf("swapping log: %w", err)
    }

    s.log.Close()
    nf, err := os.OpenFile(s.path, os.O_RDWR|os.O_APPEND, 0o600)
    if err != nil {
        return fmt.Errorf("reopening log: %w", err)
    }
    s.log, s.w = nf, bufio.NewWriter(nf)
    return nil
}
```

`os.Rename` on the same filesystem is atomic, which is the trick that makes
this safe. Write the new file completely, sync it, then swap. A crash at any
point leaves either the old log or the new one, and both are valid.

#bridge[
  Write to a temp file, fsync, rename is *the* durable-replace pattern, and it
  applies far beyond databases. Any time you overwrite a file that matters, a
  configuration file or an export, this is the correct way to do it. Writing in
  place can leave a half-written file that nothing can read.
]

#project(title: "Project 34: make it crash and watch it recover")[
  The tests are the chapter. Each proves a property the code was written for.

  ```go
  package kv

  import (
      "os"
      "path/filepath"
      "testing"
  )

  func TestSurvivesReopen(t *testing.T) {
      path := filepath.Join(t.TempDir(), "test.log")

      s, err := Open(path)
      if err != nil {
          t.Fatal(err)
      }
      if err := s.Set("user:1", []byte("Ada")); err != nil {
          t.Fatal(err)
      }
      if err := s.Set("user:2", []byte("Grace")); err != nil {
          t.Fatal(err)
      }
      if err := s.Delete("user:1"); err != nil {
          t.Fatal(err)
      }
      s.Close()

      // Reopen: everything comes from the log.
      s2, err := Open(path)
      if err != nil {
          t.Fatal(err)
      }
      defer s2.Close()

      if _, ok := s2.Get("user:1"); ok {
          t.Error("deleted key came back")
      }
      v, ok := s2.Get("user:2")
      if !ok || string(v) != "Grace" {
          t.Errorf("user:2 = %q %v, want Grace true", v, ok)
      }
  }

  // The important one: a log cut off mid record must not break the store.
  func TestTornTailIsDiscarded(t *testing.T) {
      path := filepath.Join(t.TempDir(), "torn.log")

      s, _ := Open(path)
      s.Set("a", []byte("one"))
      s.Set("b", []byte("two"))
      s.Close()

      // Simulate a crash during the third write by appending a partial
      // record: a header with no body behind it.
      f, _ := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o600)
      f.Write(encode(opSet, []byte("c"), []byte("three"))[:8])
      f.Close()

      s2, err := Open(path)
      if err != nil {
          t.Fatalf("a torn tail must not fail the open: %v", err)
      }
      defer s2.Close()

      if _, ok := s2.Get("c"); ok {
          t.Error("a half written record was replayed")
      }
      if v, ok := s2.Get("b"); !ok || string(v) != "two" {
          t.Error("records before the tear were lost")
      }

      // And the store must still be writable afterwards.
      if err := s2.Set("d", []byte("four")); err != nil {
          t.Fatalf("writing after recovery: %v", err)
      }
  }

  func TestCompactionKeepsState(t *testing.T) {
      path := filepath.Join(t.TempDir(), "compact.log")
      s, _ := Open(path)

      for i := 0; i < 1000; i++ {
          s.Set("counter", []byte("value"))
      }
      before, _ := os.Stat(path)

      if err := s.Compact(); err != nil {
          t.Fatal(err)
      }
      after, _ := os.Stat(path)

      if after.Size() >= before.Size() {
          t.Errorf("compaction did not shrink: %d -> %d",
              before.Size(), after.Size())
      }
      if v, ok := s.Get("counter"); !ok || string(v) != "value" {
          t.Error("compaction lost the data")
      }
      s.Close()

      s2, _ := Open(path)
      defer s2.Close()
      if _, ok := s2.Get("counter"); !ok {
          t.Error("compacted log did not replay")
      }
  }
  ```

  ```bash
  go test -race -v ./...
  ```

  `TestTornTailIsDiscarded` is the one to run first. It writes eight bytes of a
  record that should have been twenty, which is exactly what a power cut mid
  write leaves behind, and asserts three things: the open succeeds, the partial
  record is not replayed, and the store still accepts writes. A store that
  fails any of those will one day refuse to start, and refusing to start is
  how a durable store loses data.
]

#tryit[
  Three extensions, each teaching something different.

  Add a `syncEvery` option so a caller can choose the durability dial, and
  benchmark `Set` with sync on every write against sync every 100ms. The ratio
  you measure is the reason databases have that setting.

  Add automatic compaction when the log exceeds a multiple of the live data
  size, and think about what happens if the process dies during compaction.
  The rename makes it safe; prove it with a test that leaves a stale
  `.compact` file behind.

  Then replace the map with the B-tree from the next chapter and watch the
  store gain ordered range queries, which is the step from a key-value store
  to something you could call a database.
]
