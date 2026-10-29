#!/usr/bin/env python3
"""Build the Grit Framework Discord server, repeatably.

Run it as often as you like: it looks everything up by name, creates what is
missing and leaves what already matches alone. That is the whole point of
putting the server in a script rather than clicking it together once. The
layout is reviewable in a diff, and rebuilding it somewhere else is one
command.

    python scripts/discord/setup.py            # dry run, prints the plan
    python scripts/discord/setup.py --apply    # does it
    python scripts/discord/setup.py --apply --prune   # also removes the
                                                      # default channels

Credentials come from ~/.grit/discord.env and never from this file:

    DISCORD_BOT_TOKEN=...
    DISCORD_GUILD_ID=...

That file is outside the repository on purpose. A bot token is a full login
for the bot, and the one way to be certain it is never committed is for it not
to be in the tree.
"""

import argparse
import json
import os
import sys
import time
import urllib.error
import urllib.request

API = "https://discord.com/api/v10"

# Grit's accent, so the server and the product are the same purple.
ACCENT = 0x6C5CE7

# Permission bits, named, because 2048 in a diff tells nobody anything.
VIEW_CHANNEL = 1 << 10
SEND_MESSAGES = 1 << 11
ADD_REACTIONS = 1 << 6
CREATE_PUBLIC_THREADS = 1 << 35

# Channel types.
TEXT, VOICE, CATEGORY = 0, 2, 4


# --------------------------------------------------------------- the structure

# Roles, highest first. The bot's own role sits above these because Discord
# will not let it hand out a role above itself.
ROLES = [
    {
        "name": "Maintainer",
        "color": ACCENT,
        "hoist": True,
        "mentionable": True,
        "why": "The people who ship Grit. Separated in the member list so a question can find them.",
    },
    {
        "name": "Contributor",
        "color": 0x00B894,
        "hoist": True,
        "mentionable": True,
        "why": "Anyone whose pull request has been merged.",
    },
    {
        "name": "Early Adopter",
        "color": 0x74B9FF,
        "hoist": False,
        "mentionable": False,
        "why": "Shipped something real with Grit. Worth being able to ask them things.",
    },
]

# Categories and their channels, in order. `readonly` denies @everyone
# SEND_MESSAGES: an announcement nobody can reply to under is the point of an
# announcement.
STRUCTURE = [
    {
        "category": "START HERE",
        "channels": [
            {
                "name": "welcome",
                "topic": "What Grit is, and the one command that gets you running.",
                "readonly": True,
            },
            {
                "name": "rules",
                "topic": "Four of them. Read once.",
                "readonly": True,
            },
            {
                "name": "announcements",
                "topic": "Releases, breaking changes and anything you need to act on.",
                "readonly": True,
            },
            {
                "name": "releases",
                "topic": "Every Grit release, posted automatically from GitHub.",
                "readonly": True,
            },
        ],
    },
    {
        "category": "HELP",
        "channels": [
            {
                "name": "install-and-setup",
                "topic": "grit new, Docker, Postgres, pnpm. Anything between you and a running project.",
            },
            {
                "name": "questions",
                "topic": "How do I do X in Grit. Paste the error, say what you expected.",
            },
            {
                "name": "bugs",
                "topic": "Something is wrong. A reproduction beats a description: the command, the output, the version from `grit version`.",
            },
        ],
    },
    {
        "category": "BUILD",
        "channels": [
            {
                "name": "general",
                "topic": "Go, Next.js, and whatever you are building.",
            },
            {
                "name": "showcase",
                "topic": "Things you shipped with Grit. Screenshots encouraged.",
            },
            {
                "name": "ideas-and-feedback",
                "topic": "What Grit should do next, and what it gets wrong today. The second one is more useful.",
            },
            {
                "name": "contributing",
                "topic": "Working on Grit itself: the repo, the tests, the release process.",
            },
        ],
    },
    {
        "category": "VOICE",
        "channels": [
            {"name": "Lounge", "type": VOICE},
            {"name": "Pairing", "type": VOICE},
        ],
    },
]

# The default channels Discord's template creates, which the structure above
# replaces. Removed only with --prune, and only when empty of anything the
# script did not expect.
DEFAULTS_TO_PRUNE = [
    ("Information", CATEGORY),
    ("Text Channels", CATEGORY),
    ("Voice Channels", CATEGORY),
    ("welcome-and-rules", TEXT),
    ("general", TEXT),
    ("events", TEXT),
    ("ideas-and-feedback", TEXT),
    ("announcements", TEXT),
    ("Lounge", VOICE),
    ("Community Hangout", VOICE),
    ("Stream Room", VOICE),
]

WELCOME = """# Grit Framework

**Go + React. Built with Grit.**

One command gives you a Go API, a Next.js front end, an admin panel that
builds itself from your models, and the parts every product needs anyway:
auth, roles, file storage, background jobs, a database browser and a firewall.

**Install** (no Go toolchain needed)
```
curl -fsSL https://gritframework.dev/install.sh | sh
```
On Windows PowerShell: `irm https://gritframework.dev/install.ps1 | iex`

**Then**
```
grit new myapp --triple
```

Docs <https://gritframework.dev>
Source <https://github.com/MUKE-coder/grit>
YouTube <https://www.youtube.com/@GritFramework>

New here? Say hello in #general, and ask in #questions the moment something
does not work. A question that turns out to be a framework bug is the most
useful thing you can post.
"""

RULES = """# Rules

**1. Be decent.** Disagree about code as much as you like. Not about people.

**2. Ask the real question.** Paste the error, the command and the output of
`grit version`. "It does not work" cannot be answered; a stack trace usually
answers itself.

**3. Search before posting.** The docs have a search, and so does this server.

**4. No spam, no DM advertising.** Showing what you built in #showcase is not
spam. Messaging twenty members about your agency is.

Bugs belong in #bugs here, or as a GitHub issue if you would rather it be
tracked: <https://github.com/MUKE-coder/grit/issues>
"""


# ---------------------------------------------------------------- the plumbing


class Discord:
    def __init__(self, token, dry_run):
        self.token = token
        self.dry_run = dry_run
        self.calls = 0

    def request(self, method, path, body=None, *, write=False):
        """One API call. Writes are skipped and described during a dry run."""
        if write and self.dry_run:
            return {"_dry_run": True}
        url = API + path
        data = json.dumps(body).encode() if body is not None else None
        req = urllib.request.Request(url, data=data, method=method)
        req.add_header("Authorization", "Bot " + self.token)
        req.add_header("Content-Type", "application/json")
        req.add_header("User-Agent", "GritSetup (https://github.com/MUKE-coder/grit, 1.0)")
        for attempt in range(6):
            try:
                self.calls += 1
                with urllib.request.urlopen(req, timeout=30) as res:
                    raw = res.read()
                    return json.loads(raw) if raw else {}
            except urllib.error.HTTPError as err:
                raw = err.read()
                # Discord asks you to wait rather than refusing outright, and
                # a script that ignores that gets the token rate limited for
                # everyone using it.
                if err.code == 429:
                    wait = 1.0
                    try:
                        wait = float(json.loads(raw).get("retry_after", 1.0))
                    except Exception:
                        pass
                    print(f"    rate limited, waiting {wait:.1f}s")
                    time.sleep(wait + 0.1)
                    continue
                detail = raw.decode(errors="replace")[:400]
                raise SystemExit(f"\n{method} {path} failed: HTTP {err.code}\n{detail}\n")
        raise SystemExit(f"{method} {path}: still rate limited after 6 attempts")

    def get(self, path):
        return self.request("GET", path)

    def post(self, path, body):
        return self.request("POST", path, body, write=True)

    def patch(self, path, body):
        return self.request("PATCH", path, body, write=True)

    def delete(self, path):
        return self.request("DELETE", path, write=True)


def load_credentials():
    path = os.path.join(os.path.expanduser("~"), ".grit", "discord.env")
    if not os.path.exists(path):
        raise SystemExit(
            f"No credentials at {path}.\n"
            "Create it with DISCORD_BOT_TOKEN and DISCORD_GUILD_ID, one per line."
        )
    values = {}
    with open(path, encoding="utf-8") as handle:
        for line in handle:
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, _, value = line.partition("=")
            values[key.strip()] = value.strip()
    token, guild = values.get("DISCORD_BOT_TOKEN"), values.get("DISCORD_GUILD_ID")
    if not token or not guild:
        raise SystemExit(f"{path} needs both DISCORD_BOT_TOKEN and DISCORD_GUILD_ID")
    return token, guild


# ------------------------------------------------------------------- the work


def readonly_overwrites(everyone_id, bot_id):
    """@everyone may read and react but not post, and the bot may still post.

    A channel-level deny beats a server-level allow, and the bot belongs to
    @everyone like any other member: denying SEND_MESSAGES there denies it to
    the bot too. The first run of this script created #welcome read-only and
    was then refused permission to post the welcome into it.
    """
    return [
        {
            "id": everyone_id,
            "type": 0,  # role
            "allow": str(VIEW_CHANNEL | ADD_REACTIONS),
            "deny": str(SEND_MESSAGES | CREATE_PUBLIC_THREADS),
        },
        {
            "id": bot_id,
            "type": 1,  # member
            "allow": str(VIEW_CHANNEL | SEND_MESSAGES),
            "deny": "0",
        },
    ]


def sync_roles(api, guild, existing):
    """Create the roles that are missing. Existing ones are left alone, because
    somebody may have changed a colour on purpose."""
    by_name = {r["name"]: r for r in existing}
    for spec in ROLES:
        if spec["name"] in by_name:
            print(f"  = role {spec['name']}  (already there)")
            continue
        print(f"  + role {spec['name']}  {spec['why']}")
        api.post(
            f"/guilds/{guild}/roles",
            {
                "name": spec["name"],
                "color": spec["color"],
                "hoist": spec["hoist"],
                "mentionable": spec["mentionable"],
                "permissions": "0",
            },
        )


def sync_channels(api, guild, everyone_id, bot_id, existing):
    """Create the categories and channels, and put each channel in its
    category. Returns the channels it knows about, by name."""
    by_name = {}
    for channel in existing:
        by_name.setdefault((channel["name"], channel["type"]), channel)

    created = {}
    for group in STRUCTURE:
        cat_key = (group["category"], CATEGORY)
        category = by_name.get(cat_key)
        if category:
            print(f"  = category {group['category']}  (already there)")
        else:
            print(f"  + category {group['category']}")
            category = api.post(
                f"/guilds/{guild}/channels",
                {"name": group["category"], "type": CATEGORY},
            )
        # On a dry run the category was not really created, so it has no id.
        # Left as None, every "is this channel in the right category" test
        # below compares against None and quietly answers "no move needed",
        # and the plan under-reports the moves it is about to make. A sentinel
        # id keeps the comparison meaningful without inventing a real one.
        parent_id = category.get("id") or ("(new) " + group["category"])

        for spec in group["channels"]:
            kind = spec.get("type", TEXT)
            key = (spec["name"], kind)
            if key in by_name:
                # A channel of this name already exists, so it is adopted
                # rather than duplicated. Adopting it means moving it: a
                # channel left in Discord's "Text Channels" while the rest of
                # the server is reorganised around it is the half-finished
                # result that makes people redo the whole thing by hand.
                found = by_name[key]
                created[spec["name"]] = found
                changes = {}
                if parent_id and found.get("parent_id") != parent_id:
                    changes["parent_id"] = parent_id
                if spec.get("topic") and found.get("topic") != spec["topic"]:
                    changes["topic"] = spec["topic"]
                if spec.get("readonly"):
                    changes["permission_overwrites"] = readonly_overwrites(everyone_id, bot_id)
                if not changes:
                    print(f"      = #{spec['name']}  (already correct)")
                    continue
                what = []
                if "parent_id" in changes:
                    what.append("moved here")
                if "topic" in changes:
                    what.append("new topic")
                if "permission_overwrites" in changes:
                    what.append("read only")
                print(f"      ~ #{spec['name']}  ({', '.join(what)})")
                api.patch(f"/channels/{found['id']}", changes)
                continue
            body = {"name": spec["name"], "type": kind}
            if parent_id:
                body["parent_id"] = parent_id
            if spec.get("topic"):
                body["topic"] = spec["topic"]
            if spec.get("readonly"):
                # @everyone may read and react, but not post. The role id of
                # @everyone is the guild id, which is a Discord quirk worth
                # knowing rather than looking up.
                body["permission_overwrites"] = readonly_overwrites(everyone_id, bot_id)
            note = "  (read only)" if spec.get("readonly") else ""
            print(f"      + #{spec['name']}{note}")
            made = api.post(f"/guilds/{guild}/channels", body)
            created[spec["name"]] = made
    return created


def post_once(api, channel, content, label):
    """Post a message, unless this bot has already posted one here.

    Re-running the script should not leave four copies of the rules.
    """
    if not channel or "id" not in channel:
        print(f"      (dry run) would post {label}")
        return
    history = api.get(f"/channels/{channel['id']}/messages?limit=50")
    for message in history:
        if message.get("author", {}).get("bot") and message.get("content", "").startswith(
            content.split("\n", 1)[0]
        ):
            print(f"      = {label} already posted")
            return
    print(f"      + posting {label}")
    api.post(f"/channels/{channel['id']}/messages", {"content": content})


def prune(api, existing, structure_names):
    """Remove Discord's template channels, and only those."""
    for name, kind in DEFAULTS_TO_PRUNE:
        for channel in existing:
            if channel["name"] != name or channel["type"] != kind:
                continue
            if (channel["name"], channel["type"]) in structure_names:
                # A name the new structure also uses: keep it, it was adopted.
                continue
            print(f"  - removing default {name}")
            api.delete(f"/channels/{channel['id']}")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--apply", action="store_true", help="make the changes")
    parser.add_argument("--prune", action="store_true", help="also delete the default channels")
    args = parser.parse_args()

    token, guild = load_credentials()
    api = Discord(token, dry_run=not args.apply)

    me = api.get("/users/@me")
    info = api.get(f"/guilds/{guild}")
    print(f"\n{'PLAN (nothing is changed)' if not args.apply else 'APPLYING'}")
    print(f"server: {info['name']}   bot: {me['username']}\n")

    roles = api.get(f"/guilds/{guild}/roles")
    channels = api.get(f"/guilds/{guild}/channels")
    everyone = guild  # @everyone's role id is the guild id

    print("Roles")
    sync_roles(api, guild, roles)

    print("\nChannels")
    made = sync_channels(api, guild, everyone, me["id"], channels)

    # Order the channels inside each category to match STRUCTURE. An adopted
    # channel keeps whatever position the template gave it, so #announcements
    # sat above #welcome in START HERE: the first thing a new member sees was
    # the second thing it should be.
    print("\nOrder")
    order = []
    for group in STRUCTURE:
        for i, spec in enumerate(group["channels"]):
            found = made.get(spec["name"])
            if found and "id" in found and found.get("position") != i:
                order.append({"id": found["id"], "position": i})
    if order:
        print(f"  setting the position of {len(order)} channel(s)")
        api.patch(f"/guilds/{guild}/channels", order)
    else:
        print("  already in order")

    print("\nPosts")
    post_once(api, made.get("welcome"), WELCOME, "the welcome message")
    post_once(api, made.get("rules"), RULES, "the rules")

    if args.prune:
        print("\nDefaults")
        keep = {(s["name"], s.get("type", TEXT)) for g in STRUCTURE for s in g["channels"]}
        keep |= {(g["category"], CATEGORY) for g in STRUCTURE}
        prune(api, channels, keep)
    else:
        print("\nDefaults: left alone. Pass --prune to remove Discord's template channels.")

    if not args.apply:
        print("\nNothing was changed. Run again with --apply.")
    print(f"\n{api.calls} API calls.\n")


if __name__ == "__main__":
    main()
