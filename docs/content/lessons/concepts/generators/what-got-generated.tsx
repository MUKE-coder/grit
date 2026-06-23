import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Eight files appeared when you ran{' '}
        <code>grit generate resource Contact …</code>. This lesson walks
        through each one with the <strong>full source</strong> the
        generator wrote, so you stop seeing them as a black box and
        start treating them as your own code. Once you&apos;ve read each
        layer, you&apos;ll know exactly where to make any future change.
      </p>

      <h2>1. The Go model — <code>apps/api/internal/models/contact.go</code></h2>
      <p>
        The GORM struct that maps to the <code>contacts</code> table. Every
        field becomes a column; the struct tags tell GORM and the JSON
        encoder how to translate.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/models/contact.go"
        code={`package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Contact represents a contact in the system.
type Contact struct {
	ID        string         ` + '`' + `gorm:"primarykey;size:36" json:"id"` + '`' + `
	Name      string         ` + '`' + `gorm:"size:255" json:"name" binding:"required"` + '`' + `
	Email     string         ` + '`' + `gorm:"size:255;uniqueIndex" json:"email" binding:"required"` + '`' + `
	Phone     string         ` + '`' + `gorm:"size:255" json:"phone"` + '`' + `
	Version   int            ` + '`' + `gorm:"not null;default:1" json:"version"` + '`' + `
	CreatedAt time.Time      ` + '`' + `json:"created_at"` + '`' + `
	UpdatedAt time.Time      ` + '`' + `json:"updated_at"` + '`' + `
	DeletedAt gorm.DeletedAt ` + '`' + `gorm:"index" json:"-"` + '`' + `
}

// BeforeCreate generates a UUID before inserting.
func (m *Contact) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Contact) BeforeUpdate(tx *gorm.DB) error {
	m.Version++
	return nil
}`}
      />

      <p>Three things worth noticing:</p>
      <ul>
        <li>
          <strong>UUID string primary key.</strong> Grit uses{' '}
          <code>uuid.New().String()</code> set in a{' '}
          <code>BeforeCreate</code> hook — IDs can&apos;t be guessed
          (good IDOR defence) and the same value works across SQLite,
          Postgres, and MySQL.
        </li>
        <li>
          <strong>Soft delete is baked in.</strong>{' '}
          <code>gorm.DeletedAt</code> means <code>DELETE /api/contacts/:id</code>{' '}
          sets <code>deleted_at = now()</code> instead of actually
          dropping the row. Queries hide soft-deleted rows by default.
        </li>
        <li>
          <strong>Optimistic concurrency via Version.</strong> Every
          server-side update bumps <code>version</code>. Offline-first
          clients can detect &quot;someone else moved this&quot; and
          merge cleanly.
        </li>
      </ul>

      <h2>2. The service — <code>apps/api/internal/services/contact.go</code></h2>
      <p>
        All the actual logic. The service owns the database — handlers
        call into it; it never touches Gin or HTTP. (This is the
        convention that makes Grit handlers thin and services testable.)
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/services/contact.go (abridged)"
        code={`package services

import (
	"fmt"

	"gorm.io/gorm"
	"myapp/internal/models"
)

type ContactService struct{ db *gorm.DB }

func NewContactService(db *gorm.DB) *ContactService { return &ContactService{db: db} }

// CreateContactInput is what handlers bind requests into.
type CreateContactInput struct {
	Name  string ` + '`' + `json:"name"  binding:"required"` + '`' + `
	Email string ` + '`' + `json:"email" binding:"required,email"` + '`' + `
	Phone string ` + '`' + `json:"phone"` + '`' + `
}

func (s *ContactService) Create(in CreateContactInput) (*models.Contact, error) {
	c := &models.Contact{Name: in.Name, Email: in.Email, Phone: in.Phone}
	if err := s.db.Create(c).Error; err != nil {
		return nil, fmt.Errorf("creating contact: %w", err)
	}
	return c, nil
}

// List supports pagination + search across string-shaped fields.
func (s *ContactService) List(page, pageSize int, search string) ([]models.Contact, int64, error) {
	q := s.db.Model(&models.Contact{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}

	var total int64
	q.Count(&total)

	var contacts []models.Contact
	err := q.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at DESC").Find(&contacts).Error
	return contacts, total, err
}

// GetByID, Update, Delete — same pattern: just GORM calls with error wrapping.`}
      />

      <TipBox tone="info">
        Notice the auto-generated search clause covers{' '}
        <code>name</code>, <code>email</code>, and <code>phone</code> —
        every string-shaped column. That&apos;s why{' '}
        <code>GET /api/contacts?search=alice</code> just works after
        generation: the service already knows which columns to scan.
      </TipBox>

      <h2>3. The handler — <code>apps/api/internal/handlers/contact.go</code></h2>
      <p>
        Thin layer that bridges HTTP and the service. Five methods, one
        per CRUD verb.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/contact.go (Create only — others mirror it)"
        code={`package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"myapp/internal/respond"
	"myapp/internal/services"
)

type ContactHandler struct{ svc *services.ContactService }

func NewContactHandler(svc *services.ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

func (h *ContactHandler) Create(c *gin.Context) {
	var in services.CreateContactInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	out, err := h.svc.Create(in)
	if err != nil {
		respond.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respond.Created(c, out, "Contact created")
}`}
      />

      <p>
        <code>respond.Error</code> and <code>respond.Created</code> shape
        the JSON to match{' '}
        <a href="/docs/backend/response-format">Grit&apos;s response envelope</a>{' '}
        — so the frontend always gets <code>{`{ data, message }`}</code>{' '}
        or <code>{`{ error: { code, message } }`}</code>, no matter
        which handler it called.
      </p>

      <h2>4. The route injection — <code>apps/api/internal/routes/routes.go</code></h2>
      <p>
        The generator <em>edits</em> the existing routes file rather than
        creating a new one. It finds the marker comment{' '}
        <code>{`// grit:routes`}</code> and slots the new mounting block
        in before it:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (excerpt)"
        code={`// Contact resource — auto-generated, edit freely
contacts := api.Group("/contacts")
contacts.Use(middleware.Auth(cfg))
{
	contacts.GET("",        contactHandler.List)
	contacts.POST("",       contactHandler.Create)
	contacts.GET("/:id",    contactHandler.GetByID)
	contacts.PUT("/:id",    contactHandler.Update)
	contacts.DELETE("/:id", contactHandler.Delete)
}
// grit:routes`}
      />

      <p>
        The Services struct higher up the file also gets a new field
        (<code>Contact *ContactService</code>), and the handler is
        instantiated wherever <code>NewServices()</code> lives. Open the
        file and search for <code>Contact</code> — you&apos;ll see every
        site the generator touched.
      </p>

      <h2>5. The Zod schema — <code>packages/shared/src/schemas/contact.ts</code></h2>

      <CodeBlock
        language="ts"
        filename="packages/shared/src/schemas/contact.ts"
        code={`import { z } from "zod";

export const ContactSchema = z.object({
  name:  z.string().min(1, "Required"),
  email: z.string().min(1, "Required").email("Invalid email"),
  phone: z.string().optional(),
});

export const CreateContactSchema = ContactSchema;
export const UpdateContactSchema = ContactSchema.partial();

export type CreateContactInput = z.infer<typeof CreateContactSchema>;
export type UpdateContactInput = z.infer<typeof UpdateContactSchema>;`}
      />

      <p>
        Same shape lives in <em>two</em> places: as Go struct tags
        (validated by Gin&apos;s binding) and as a Zod schema (validated
        on the frontend before the request even leaves). One source of
        truth — fields described to the generator — emits both.
      </p>

      <h2>6. The TypeScript type — <code>packages/shared/src/types/contact.ts</code></h2>

      <CodeBlock
        language="ts"
        filename="packages/shared/src/types/contact.ts"
        code={`// Auto-generated by grit — DO NOT EDIT MANUALLY.
// Run \`grit sync\` to regenerate after changing the Go model.

export interface Contact {
  id: string;
  name: string;
  email: string;
  phone: string;
  version: number;
  created_at: string;
  updated_at: string;
}`}
      />

      <TipBox tone="warning">
        Hand-edits to this file are wiped on the next{' '}
        <code>grit sync</code>. If you want a custom property, add it to
        the Go model and let the sync flow propagate it — that&apos;s
        the next lesson&apos;s topic.
      </TipBox>

      <h2>7. The React Query hook — <code>apps/web/hooks/use-contacts.ts</code></h2>

      <CodeBlock
        language="ts"
        filename="apps/web/hooks/use-contacts.ts (abridged)"
        code={`import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import type { Contact } from "@repo/shared/types";

const KEY = ["contacts"] as const;

export function useContacts(params?: { page?: number; search?: string }) {
  return useQuery({
    queryKey: [...KEY, params],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: Contact[]; meta: { total: number } }>(
        "/api/contacts",
        { params },
      );
      return data;
    },
  });
}

export function useCreateContact() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: Partial<Contact>) =>
      apiClient.post<{ data: Contact }>("/api/contacts", input).then((r) => r.data.data),
    onSuccess: () => qc.invalidateQueries({ queryKey: KEY }),
  });
}

// useContact(id), useUpdateContact, useDeleteContact — same pattern.`}
      />

      <p>
        Pagination, search, and invalidation-after-mutate are wired up
        out of the box. Drop <code>useContacts()</code> into any page
        and you have a live list.
      </p>

      <h2>8. The admin resource page — <code>apps/admin/app/resources/contacts/page.tsx</code></h2>

      <CodeBlock
        language="tsx"
        filename="apps/admin/app/resources/contacts/page.tsx"
        code={`"use client";

import { defineResource } from "@/lib/admin/define-resource";

export default defineResource({
  model: "contacts",
  label: "Contacts",
  searchable: true,
  columns: [
    { key: "name",  label: "Name",  sortable: true, format: "text" },
    { key: "email", label: "Email", sortable: true, format: "text" },
    { key: "phone", label: "Phone", format: "text" },
    { key: "created_at", label: "Created", format: "relative", sortable: true },
  ],
  form: [
    { name: "name",  type: "text",  label: "Name",  required: true },
    { name: "email", type: "text",  label: "Email", required: true },
    { name: "phone", type: "text",  label: "Phone" },
  ],
});`}
      />

      <p>
        That declarative object is everything: the DataTable columns, the
        Create/Edit form, the search box, the pagination, the delete
        confirmation, the empty state. No HTML, no form wiring. The
        Filament-style runtime in <code>@/lib/admin/define-resource</code>{' '}
        renders it.
      </p>

      <TipBox tone="success">
        The generator is a <em>starting point</em>, not a final word.
        Once the files exist, edit them. Add custom service methods,
        extend the admin columns with a status badge, change the form to
        a multi-step wizard. The generator runs once; the code is yours
        from there.
      </TipBox>

      <h2>How the eight files talk to each other</h2>

      <CodeBlock
        language="text"
        code={`             Browser              ⇆            Go API
                                                        │
    apps/admin/.../contacts/page.tsx                    │
                 │  defineResource → DataTable + Form   │
                 ▼                                       │
       useContacts() hook                                │
       (apps/web/hooks/use-contacts.ts)                  │
                 │  axios GET /api/contacts             ─┤
                 │                                       │   ┌──────────────────────────────┐
                 │                                       ├──▶│ routes.go                    │
                 │                                       │   │   contactHandler.List        │
                 │                                       │   ├──────────────────────────────┤
                 │                                       │   │ handler → service.List()     │
                 │                                       │   ├──────────────────────────────┤
                 │                                       │   │ service → db.Model(Contact)  │
                 │                                       │   ├──────────────────────────────┤
                 │  validates response against           │   │ model → contacts table       │
                 ▼  packages/shared/types/contact.ts     │   └──────────────────────────────┘
       UI re-renders with typed data
       (Contact[] from @repo/shared/types)`}
      />

      <KnowledgeCheck
        question="You want the admin Contacts page to show a 'Total contacts' card at the top. Where do you add it?"
        choices={[
          {
            label: 'Re-run grit generate with a --summary flag',
            feedback:
              'Wrong — the generator runs once. Customization happens by editing the generated files directly.',
          },
          {
            label: 'Edit apps/admin/app/resources/contacts/page.tsx',
            correct: true,
            feedback:
              "Right — the generator's output is just a starting point. Wrap the defineResource call in your own component and render a summary card above it.",
          },
          {
            label: 'Edit the service to inject a "totalContacts" field on every contact',
            feedback:
              "Wrong direction — that's modelling a per-row field, not a page summary. Mixing concerns.",
          },
          {
            label: 'Add an admin middleware that injects summaries',
            feedback:
              'Way over-engineered. UI customisation belongs in the UI.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Open <strong>three</strong> of the eight generated files
              (your pick — try one Go, one TS, one TSX) and answer in{' '}
              <code>notes.md</code>:
            </p>
            <ul>
              <li>Which file?</li>
              <li>What do the first 10 lines do?</li>
              <li>
                If you were to extend it (e.g., add{' '}
                <code>is_active:bool</code> to Contact), what would you
                change in this file?
              </li>
            </ul>
          </>
        }
        hint={
          <>
            Pick one Go file, one TS file, one TSX file — covers the
            whole stack and the answer to &quot;where do I touch this?&quot;
            becomes muscle memory.
          </>
        }
        solution={
          <>
            <p>Example for the model file:</p>
            <ul>
              <li>
                <code>apps/api/internal/models/contact.go</code> — first
                10 lines: package declaration, imports, the Contact
                struct with GORM tags.
              </li>
              <li>
                Extension: add{' '}
                <code>
                  IsActive bool {`\`gorm:"default:true" json:"is_active"\``}
                </code>{' '}
                to the struct. Run <code>grit migrate</code> — GORM
                AutoMigrate adds the column. Run <code>grit sync</code> —
                TS types pick up the new field.
              </li>
            </ul>
            <p>
              You&apos;ve now seen how generation + manual edits compose.
              The generator does the boilerplate; you do the
              product-specific work.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You&apos;ve seen the file tour with three plain string fields.
        The next lesson goes wide — every field type, slug auto-gen, the
        image-upload trick hiding inside <code>string_array</code>, and
        the heuristic names that quietly upgrade your column storage.
      </p>
    </>
  )
}
