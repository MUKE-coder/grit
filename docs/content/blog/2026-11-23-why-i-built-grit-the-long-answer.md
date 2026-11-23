---
title: "Why I built Grit: the long answer"
subtitle: "Somebody asked me why Grit exists. The short version is that I am obsessed with three things, performance, interfaces and security, and every decision in the framework comes out of one of them. Here is the whole chain, including the parts that are just me being stubborn."
series: "The Daily Grit"
date: 2026-11-23
readingTime: "9 min"
author: "Muke JohnBaptist"
tags: [grit, golang, nextjs, full-stack, meta-framework, founder-story]
canonical: "https://gritframework.dev/blog/why-i-built-grit-the-long-answer"
---

I have written about [the pain that made me start](/blog/why-i-built-grit): rebuilding auth, roles, audit trails, rate limiting, storage, email and jobs before touching the actual product, on every project. That post is about what wore me down.

This one is the question underneath it. Somebody asked me why I built Grit rather than using something that already exists, and the honest answer starts a long way before the first commit.

## Three things I am obsessed with

Everything below follows from these, so they come first.

**Performance.** I want my apps to be extremely fast, and I want them to still be fast once they are on a server. That sounds obvious and it is not: plenty of stacks are quick on a laptop and quite different under load. Chasing that took me through Node and the MERN stack, Next.js doing the whole thing, Hono, Encore.ts and Laravel. None of them gave me what I was after. Go and Rust are the two languages that did. I went with Go because it is the one a whole team can read.

**Interfaces.** I was a graphic designer before I was a programmer. I care, probably more than is reasonable, what everything looks like, and clients come back to me because of it. This matters for a framework more than it sounds: a generated admin panel that looks generated is not finished. If the thing a tool produces is something you immediately want to redesign, the tool did not save you the work.

**Security.** I like working with large companies, and the thing they worry about most is security. They have millions of customers and cannot afford a mistake. They want control of all of it: the server, the application, the firewall, the logs, the audit trail, roles and permissions. Not a checkbox. The actual controls, in their hands.

## Why Go, and why I kept Next.js

I have built more than a hundred full-stack applications. Most of them were Next.js with API routes, and that is genuinely amazing for a lot of things.

I moved to Go because an enterprise project needs a backend that carries real load, and that is what Go does. But I never wanted to leave Next.js behind. Its routing is good, its SEO story is good, and it is still the most pleasant way I know to write React. It also behaves perfectly well as a client-side app when that is what you need.

So the question became: what is the easiest possible way to build with Next.js and Go at the same time?

## The problems are always the same

If you have built full-stack apps you already know them. Anyone who has done MERN knows exactly what I am talking about.

You have a server app and a client app. Two codebases that have to agree. You have to keep the types in sync, and they drift the moment you stop paying attention. Then somebody asks for a mobile app, or a desktop app, and now you are thinking about what the repository looks like and how any of it deploys.

The answer to all of that is one repository: an API, a web app, mobile and desktop all in the same place.

And then there is a thing about the web half that took me a while to notice. For a serious product it is usually **two apps, not one**. Think of a bank. There is `mybank.com`, where customers browse and open accounts, and there is `admin.mybank.com`, where staff run the place. Two clients, one API. That shape is so common in the work I do that when I built Grit I made it the default:

```bash
grit new my-app --triple    # API + customer web app + admin panel
```

## The boring part, which is most of the work

Here is the part nobody puts in a conference talk.

A full-stack app is MVC, and for every single table you build the same things twice. On the client: a form, a table, search, export, pagination, edit, delete and a detail page. On the server: a schema, validation, a model, a service, a handler.

For every table. So what we all actually do is copy the last one and edit it. It is tedious, it is error-prone, and it is where the hours go.

But look at what is true about that work: it is **predictable**. You already know the table is called `Product`. You already know it holds a name, a price and an image. Everything in that first version can be worked out from those two facts. That is not intelligence, it is arithmetic, and a machine should do it.

```bash
grit generate resource Product --fields "name:string,price:money,image:file"
```

That writes the Go model, the service, the handler, the routes, the Zod schema, the TypeScript types, the React Query hooks and the admin screen. Not as a guess at your product: as the starting point the shape of the data already implies. The business logic is yours, and it is the only part that was ever interesting.

This is also, to me, the strongest thing about the CLI. I do not know another framework that gives you the whole vertical slice from a table name and a field list.

## Everything else arrived because I needed it twice

None of the rest was planned. Each one is a thing I wanted on a real project, built properly because I knew I would want it again.

- **[Sentinel](/docs/security)**, because clients ask about security and I wanted real answers: a firewall, blocking, audit.
- **[Pulse](/docs/backend/pulse)**, because I wanted to see which line was burning the CPU instead of guessing.
- **[GORM Studio](/docs/infrastructure/studio)**, because I came from the JavaScript world and missed Prisma Studio. I wanted the database in a browser.
- **Wails**, because clients kept wanting a desktop app.
- **Expo**, because they kept wanting a phone app.
- **API documentation**, because I was writing Swagger by hand on every project.

## The rule I will not break

Everything above could have been built as a framework that hides the details. I did not want that.

Grit generates code into your repository and gets out of the way. You can open any file it wrote. You can change any file it wrote. If a generated handler is not what you want, edit it, and [the upgrade knows you did](/docs/upgrades). Your production binary does not import Grit at all. If I disappeared tomorrow you would still have a working Go application and three ordinary React apps.

That is not modesty. It is the condition that makes the rest of it safe to use.

## Built to be used with AI

I use AI constantly, so Grit is built to be worked on by it. There is an [MCP server](/docs/plugins/mcp), there are [expert skills](/docs/expert-skills), and the CLI does enough in a single command that an agent spends a fraction of the tokens it would spend writing those same files by hand. One command instead of forty tool calls is cheaper, faster, and far less likely to go wrong halfway.

## And plugins, because nothing ships everything

No framework can include everything, and one that tries becomes a thing nobody can read. So the rest is [plugins](/plugins): install what you need, write your own, publish it for everyone else.

---

That is the whole answer. Grit is the thing that makes my own work faster, built by someone who could not stop caring about how fast it runs, how it looks, and whether it is safe to put in front of a million people.

If that sounds like your three things too, it will probably suit you.

```bash
curl -fsSL https://gritframework.dev/install.sh | sh
grit new my-app --triple --next
```

*Describe your data. Get the whole app.*
