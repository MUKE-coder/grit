---
page: /docs/batteries/backups
title: grit backup writes an archive without object storage
---

The backups page says `--output` writes a local archive and needs no storage
credentials, which is the form a person reaches for first: before any bucket
exists, and when they want a copy on the machine in front of them.

That path is easy to break from the other side. The default uploads to object
storage, so a change to the uploader can leave `--output` requiring credentials
it should never have asked for, and the error arrives as a failed backup rather
than as a wrong flag.

The archive is also checked for what the page says is in it: a CSV per table, a
`dump.sql` of INSERTs in parent-to-child order, and a `metadata.json` manifest.

```bash
grit generate resource Invoice --fields "total:money,reference:string"
grit migrate
grit seed

grit backup --output /tmp/grit-backup.zip
test -s /tmp/grit-backup.zip
unzip -l /tmp/grit-backup.zip | grep -q 'metadata.json'
unzip -l /tmp/grit-backup.zip | grep -q 'dump.sql'
unzip -l /tmp/grit-backup.zip | grep -q 'invoices.csv'
```
