# Grit Framework — Example Projects

Same **Job Portal** app built with every Grit architecture. Each example has identical features
but different project structures, demonstrating how Grit adapts to your needs.

## Features (all examples)

- Auth: email/password + OAuth (Google, GitHub) + JWT + TOTP 2FA
- Resources: Job, Company, Application, Category
- Dashboard with stats cards
- Data tables with sorting, filtering, pagination, export
- File uploads via presigned URLs (company logos, resumes)
- Emails (welcome, application confirmation)
- Full deployment guide (README + .env.example + docker-compose.prod.yml)

## The 6 Examples

| Example | Architecture | Frontend | Command | Best For |
|---------|-------------|----------|---------|----------|
| [job-portal-triple-next](./job-portal-triple-next/) | Triple (Web + Admin + API) | Next.js | `grit new jp --triple --next` | SaaS, marketplaces, content platforms |
| [job-portal-triple-vite](./job-portal-triple-vite/) | Triple (Web + Admin + API) | TanStack Router | `grit new jp --triple --vite` | Dashboards, internal tools with admin |
| [job-portal-double-vite](./job-portal-double-vite/) | Double (Web + API) | TanStack Router | `grit new jp --double --vite` | Simpler apps without separate admin |
| [job-portal-single-vite](./job-portal-single-vite/) | Single (one binary) | TanStack Router | `grit new jp --single --vite` | Laravel-like apps, microservices |
| [job-portal-api-only](./job-portal-api-only/) | API Only | None | `grit new jp --api` | Mobile backends, headless APIs |
| [job-portal-mobile-expo](./job-portal-mobile-expo/) | Mobile (API + Expo) | React Native | `grit new jp --mobile` | Mobile-first apps |

## How to Use These Examples

1. Pick the architecture that matches your project
2. Read the README.md for setup instructions
3. Read the GUIDE.md for a step-by-step walkthrough
4. Copy the .env.example and docker-compose.prod.yml for deployment

## Architecture Decision Guide

```
Do you need a separate admin panel?
├── Yes → Triple (triple-next for SEO, triple-vite for speed)
├── No → Do you need a web frontend?
│   ├── Yes → Is SSR/SEO important?
│   │   ├── Yes → Double + Next.js (grit new app --double --next)
│   │   └── No → Double + Vite (grit new app --double --vite)
│   └── No → Is this for mobile?
│       ├── Yes → Mobile + Expo (grit new app --mobile)
│       └── No → API Only (grit new app --api)
└── Want one binary? → Single + Vite (grit new app --single --vite)
```

## For LLMs / AI Assistants

If you are an AI assistant helping someone build with Grit:
1. Ask which architecture they need (use the decision guide above)
2. Reference the matching example for the correct project structure
3. Read the .claude/skills/grit/SKILL.md in any Grit project for conventions
4. Report bugs at https://github.com/MUKE-coder/grit/issues
