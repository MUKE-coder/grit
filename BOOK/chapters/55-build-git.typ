#import "../template.typ": *

= Build your own Git

Git is the tool you use most and understand least. Its interface is famously
awkward and its data model is famously elegant, and the gap between those two
facts is why most people learn commands rather than concepts.

The model is small enough to implement in an afternoon: a content-addressed
store with four object types. By the end of this chapter your program will
write objects that the real `git` command reads, and `git log` will show
commits your code created.

== Everything is an object, named by its hash

```
blob    a file's contents. No name, no permissions, just bytes
tree    a directory: names, modes, and the hashes they point at
commit  a tree, parents, author, message
tag     an annotated tag
```

An object is stored as:

```
<type> <size>\0<contents>
```

compressed with zlib, and its name is the SHA-1 of *that whole thing*,
uncompressed, in hexadecimal. The path is the first two hex characters as a
directory and the rest as the filename:

```
.git/objects/e6/9de29bb2d1d6434b8b29ae775ad8c2e48c5391
```

```go
func hashObject(typ string, content []byte) (string, []byte) {
    header := fmt.Sprintf("%s %d\x00", typ, len(content))
    full := append([]byte(header), content...)

    sum := sha1.Sum(full)
    return hex.EncodeToString(sum[:]), full
}
```

Fifteen lines, and that function is the whole of Git's identity model.

#bridge[
  "Content addressed" means the name *is* the content, which has three
  consequences that explain most of Git's behaviour.

  *Identical files are stored once*, however many directories or branches
  contain them.

  *Nothing can be modified.* Changing a file produces a different hash, so it
  is a new object. The old one is still there.

  *A commit's hash covers its entire history*, because it names its tree and
  its parent, which name theirs. That is why rewriting one old commit changes
  the hash of every commit after it, and why a force push is a different
  history rather than an edited one.
]

== init, and the directory

```go
func cmdInit(dir string) error {
    for _, d := range []string{
        "objects", "refs/heads", "refs/tags",
    } {
        if err := os.MkdirAll(filepath.Join(dir, ".git", d),
            0o755); err != nil {
            return err
        }
    }
    return os.WriteFile(filepath.Join(dir, ".git", "HEAD"),
        []byte("ref: refs/heads/main\n"), 0o644)
}
```

That is a valid repository. `git status` will run in it.

`HEAD` is a *symbolic reference*: a file containing the name of another
reference. That one indirection is the entire branch mechanism. Switching
branches rewrites that line; a detached HEAD is the same file containing a hash
instead.

== Writing and reading objects

```go
func writeObject(root, typ string, content []byte) (string, error) {
    sum, full := hashObject(typ, content)
    path := objectPath(root, sum)

    if _, err := os.Stat(path); err == nil {
        return sum, nil          // already have it, by definition
    }
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return "", err
    }

    var buf bytes.Buffer
    zw := zlib.NewWriter(&buf)
    if _, err := zw.Write(full); err != nil {
        return "", err
    }
    if err := zw.Close(); err != nil {
        return "", err
    }
    return sum, os.WriteFile(path, buf.Bytes(), 0o444)
}
```

Two details that are not accidents. *If the file exists, you already have the
content*, because the name is the content, so there is nothing to compare.
And the mode is `0444`, read-only: objects are immutable, and the filesystem
is told so.

== Trees, and the format that catches everybody

A tree is a sequence of entries, with no separators and no JSON:

```
<mode> <name>\0<20 raw bytes of sha>
```

```go
func writeTree(root, dir string) (string, error) {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return "", err
    }
    var buf bytes.Buffer

    for _, e := range sortTreeEntries(entries) {
        if e.Name() == ".git" {
            continue
        }
        var mode, sum string

        if e.IsDir() {
            mode = "40000"
            sum, err = writeTree(root, filepath.Join(dir, e.Name()))
        } else {
            mode = fileMode(e)          // 100644 or 100755
            data, rerr := os.ReadFile(filepath.Join(dir, e.Name()))
            if rerr != nil {
                return "", rerr
            }
            sum, err = writeObject(root, "blob", data)
        }
        if err != nil {
            return "", err
        }
        raw, _ := hex.DecodeString(sum)
        fmt.Fprintf(&buf, "%s %s\x00", mode, e.Name())
        buf.Write(raw)
    }
    return writeObject(root, "tree", buf.Bytes())
}
```

The modes are a short list: `100644` a file, `100755` an executable, `40000` a
directory, `120000` a symlink, `160000` a submodule. They look like Unix
permissions and are not: Git stores only whether a file is executable.

#gotcha[
  Two things here will produce objects real Git rejects, and both are silent.

  *The hash is twenty raw bytes, not forty hex characters.* It is the one place
  in the format where the hash is not text.

  *Entries must be sorted*, and the sort is not plain byte order on the name: a
  directory sorts as though its name ended with `/`. So `lib.rs` comes before
  `lib/`, which is the opposite of what a naive sort gives.

  Get either wrong and your objects hash differently from Git's for the same
  content, which means your repository and a real one disagree about history
  that is identical. Check with `git cat-file -p` early rather than after
  building three more layers.
]

== Commits

```go
func commitTree(root, tree string, parents []string,
    author string, when time.Time, msg string) (string, error) {

    var b strings.Builder
    fmt.Fprintf(&b, "tree %s\n", tree)
    for _, p := range parents {
        fmt.Fprintf(&b, "parent %s\n", p)
    }
    stamp := fmt.Sprintf("%d %s", when.Unix(), when.Format("-0700"))
    fmt.Fprintf(&b, "author %s %s\n", author, stamp)
    fmt.Fprintf(&b, "committer %s %s\n", author, stamp)
    fmt.Fprintf(&b, "\n%s\n", msg)

    return writeObject(root, "commit", []byte(b.String()))
}
```

Plain text, which is why `git cat-file -p HEAD` is readable.

Zero parents is the first commit. One parent is ordinary. *Two or more parents
is a merge*, and that is the entire representation of merging: a commit that
points at both histories. There is no other record of a merge anywhere in Git.

== References, and what a branch is

```go
func updateRef(root, ref, sum string) error {
    path := filepath.Join(root, ".git", ref)
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return err
    }
    return os.WriteFile(path, []byte(sum+"\n"), 0o644)
}
```

A branch is a file containing forty one bytes. That is the whole implementation.

It explains things that otherwise seem arbitrary. Creating a branch is instant
because it writes one small file. Deleting a branch deletes a file and *not a
single commit*, which is why the commits are still there and `git reflog` can
find them. Moving a branch is overwriting those forty one bytes, which is what
`reset --hard` does.

== Walking the log

```go
func cmdLog(root, start string) error {
    for sum := start; sum != ""; {
        obj, err := readObject(root, sum)
        if err != nil {
            return err
        }
        c := parseCommit(obj.content)

        fmt.Printf("commit %s\n%s\n\n    %s\n\n",
            sum, c.Author, c.Message)

        if len(c.Parents) == 0 {
            return nil
        }
        sum = c.Parents[0]       // first parent only
    }
    return nil
}
```

`git log` with no flags does exactly this: follow the first parent. The
difference between that and `--graph` is whether you also walk the other
parents of merges, which needs a queue ordered by commit time and is where
log's real complexity lives.

== The index, and why it exists

The staging area is a file, `.git/index`, holding the paths and hashes that
will go into the next tree, plus the stat information Git uses to decide
whether a file changed without reading it.

That last part is the real purpose. `git status` in a repository with a hundred
thousand files is fast because it compares size and modification time from the
index rather than hashing every file.

For this project, a simple text index is enough:

```
100644 e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 src/main.go
```

Real Git's index is binary and versioned. Replacing it is a chapter of its own,
and it is the right thing to leave out.

== What you are not building

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Packfiles], [Thousands of objects delta-compressed into one file. This is
    how a repository stays small, and it is the biggest omission here],
  [The network], [`git fetch` and `git push` speak their own protocol, with
    negotiation about which objects the other side lacks],
  [Merging], [Creating a merge commit is easy. Computing the merged *content*
    is a three-way diff algorithm],
  [Rebase], [Replaying commits onto a new base, which is merge plus
    bookkeeping],
)

Loose objects are correct and are what Git itself writes before `gc` runs. Your
repository will simply be larger than it needs to be, which for a learning
project is not a cost at all.

#gotcha[
  SHA-1 is broken for collision resistance, and Git still uses it.

  That is survivable because Git is not using the hash as a security control
  against a motivated attacker with your repository, and because Git adds a
  collision detection step on top. The project has a SHA-256 mode and a long
  transition plan, because changing the hash changes every object name in every
  repository that exists.

  The lesson worth taking: *an identifier is a decision you cannot reverse*. If
  you build a content-addressed anything today, use SHA-256, and store the
  algorithm next to the hash so the next change is possible.
]

#project(title: "Project 51: a repository real Git can read")[
  In `part6-from-scratch/55-git`, build it, and validate against the real
  thing at every step.

  *One: `init` and `hash-object`.* Hash a file with your tool and with real
  Git:

  ```bash
  echo "hello" > a.txt
  ./mygit hash-object -w a.txt
  git hash-object a.txt
  ```

  The two hashes must be identical. If they are not, stop here: your header
  format is wrong and nothing later will work.

  *Two: `cat-file`.* Read any object back, printing its type, size and
  contents. Check against `git cat-file -p`.

  *Three: `write-tree`.* Build a tree from a directory. Then verify with real
  Git:

  ```bash
  ./mygit write-tree
  git cat-file -p <the hash it printed>
  ```

  Real Git must be able to print your tree. This is where the sorting rule and
  the raw-bytes rule get caught.

  *Four: `commit-tree` and `log`.* Make a commit, update `refs/heads/main`,
  and then run *real* `git log`. Your commit must appear with its message and
  author.

  *Five: `status`.* Compare the working directory against the index and the
  HEAD tree. Report added, modified and deleted files.

  *Six: `checkout`.* Restore the working directory from a tree. Then make a
  commit, change files, and check out the earlier commit with your own tool.

  *Seven: the proof.* Create a repository entirely with your tool, then run
  `git fsck`. It checks every object's hash and every reference. A clean run
  is the strongest statement you can make about this project.
]

#tryit[
  In any repository you have, run:

  ```bash
  git cat-file -p HEAD
  git cat-file -p HEAD^{tree}
  git cat-file -p HEAD:README.md
  ```

  Three commands, and you have walked from a commit to a tree to a file's
  contents by hand. Everything Git does is built on that walk, and after this
  chapter you can see it.
]

#assignment(name: "Build Git - Assignment One")[
  Build the tool and prove interoperability rather than claiming it.

  *One.* `init`, `hash-object`, `cat-file`, `write-tree`, `commit-tree`,
  `log`, `status`, `checkout`.

  *Two.* `INTEROP.md` with paired output: your hash next to `git hash-object`,
  your tree printed by `git cat-file -p`, and your commit shown by real
  `git log`.

  *Three.* A clean `git fsck` on a repository your tool created, with the
  output pasted.

  *Four.* One deliberate break: sort the tree entries naively, show the hash
  differing from Git's for the same directory, then fix it and show them
  matching.

  *Five.* A paragraph explaining to a colleague, in plain language, why
  rewriting an old commit changes every commit after it. If you can write that
  from the data model rather than from memory, you have got the chapter.
]
