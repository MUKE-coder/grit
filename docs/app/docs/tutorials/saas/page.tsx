import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials/saas')

export default function TutorialSaaSPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Tutorial
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Build a Project Management SaaS
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build a complete project management tool with projects, tasks,
                comments, status workflows, email notifications, role-based
                access, and a dashboard with live task stats. This tutorial
                covers relationships, background jobs, the mailer, middleware
                customization, and admin panel badges.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Prerequisites
              </h2>
              <ul className="space-y-2.5">
                {[
                  "Go 1.21+ installed",
                  "Node.js 18+ and pnpm installed",
                  "Docker and Docker Compose installed",
                  "Grit CLI installed — see /docs/getting-started/installation (one-line install script)",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-1">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* ============================================================ */}
            {/* STEP 1 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create the project
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Scaffold a new Grit project called <code>taskflow</code>. This
                  generates the complete monorepo with Go API, Next.js web app,
                  admin panel, shared packages, and Docker configuration.
                </p>

                <CodeBlock terminal code={`grit new taskflow\ncd taskflow`} className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 2 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start Docker services
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Launch PostgreSQL, Redis, MinIO, and Mailhog. Redis is
                  especially important for this project because we will use it
                  for the background job queue that sends email notifications.
                </p>

                <CodeBlock terminal code="docker compose up -d" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 3 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Project resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Every task belongs to a project. Generate a Project resource
                  with name, description, and status fields.
                </p>

                <CodeBlock terminal code={`grit generate resource Project --fields "name:string,description:text,status:string"`} className="glow-purple-sm mb-4" />

                <CodeBlock filename="apps/api/internal/models/project.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Project struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Status      string         \`gorm:"size:50;default:active" json:"status"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 4 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Task resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Tasks are the core of the application. Generate a Task
                  resource with title, description, status, priority, and a due
                  date.
                </p>

                <CodeBlock terminal code={`grit generate resource Task --fields "title:string,description:text,status:string,priority:string,dueDate:date"`} className="glow-purple-sm mb-4" />

                <CodeBlock filename="apps/api/internal/models/task.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Task struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Title       string         \`gorm:"size:255;not null" json:"title" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Status      string         \`gorm:"size:50;default:todo" json:"status"\`
    Priority    string         \`gorm:"size:50;default:medium" json:"priority"\`
    DueDate     *time.Time     \`json:"due_date"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 5 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                5
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Comment resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Team members need to discuss tasks. Generate a Comment
                  resource with a content field. We will wire up relationships
                  to Task and User in the next step.
                </p>

                <CodeBlock terminal code={`grit generate resource Comment --fields "content:text"`} className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 6 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                6
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Set up relationships
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now wire the models together. A Task belongs to a Project and
                  is assigned to a User. A Comment belongs to a Task and is
                  authored by a User. A Project has many Tasks.
                </p>

                <CodeBlock filename="apps/api/internal/models/task.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Task struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Title       string         \`gorm:"size:255;not null" json:"title" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Status      string         \`gorm:"size:50;default:todo" json:"status"\`
    Priority    string         \`gorm:"size:50;default:medium" json:"priority"\`
    DueDate     *time.Time     \`json:"due_date"\`

    // Belongs to Project
    ProjectID   uint           \`gorm:"index;not null" json:"project_id" binding:"required"\`
    Project     Project        \`gorm:"foreignKey:ProjectID" json:"project,omitempty"\`

    // Assigned to User
    AssigneeID  *uint          \`gorm:"index" json:"assignee_id"\`
    Assignee    *User          \`gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"\`

    // Has many comments
    Comments    []Comment      \`gorm:"foreignKey:TaskID" json:"comments,omitempty"\`

    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />

                <CodeBlock filename="apps/api/internal/models/comment.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Comment struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Content   string         \`gorm:"type:text;not null" json:"content" binding:"required"\`

    // Belongs to Task
    TaskID    uint           \`gorm:"index;not null" json:"task_id" binding:"required"\`
    Task      Task           \`gorm:"foreignKey:TaskID" json:"task,omitempty"\`

    // Authored by User
    UserID    uint           \`gorm:"index;not null" json:"user_id"\`
    User      User           \`gorm:"foreignKey:UserID" json:"user,omitempty"\`

    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />

                <CodeBlock filename="apps/api/internal/models/project.go &mdash; add Tasks
                      relation" code={`type Project struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Status      string         \`gorm:"size:50;default:active" json:"status"\`
    Tasks       []Task         \`gorm:"foreignKey:ProjectID" json:"tasks,omitempty"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Run <code>grit sync</code> to regenerate the TypeScript types
                  with the new relationships.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 7 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                7
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Add a status workflow
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Tasks follow a workflow: <strong>Todo</strong> &rarr;{" "}
                  <strong>In Progress</strong> &rarr; <strong>Review</strong>{" "}
                  &rarr; <strong>Done</strong>. Define status constants and add
                  a validation method that enforces valid transitions.
                </p>

                <CodeBlock filename="apps/api/internal/models/task.go &mdash; add status
                      constants and validation" code={`// Task status constants
const (
    TaskStatusTodo       = "todo"
    TaskStatusInProgress = "in_progress"
    TaskStatusReview     = "review"
    TaskStatusDone       = "done"
)

// Task priority constants
const (
    TaskPriorityLow    = "low"
    TaskPriorityMedium = "medium"
    TaskPriorityHigh   = "high"
    TaskPriorityUrgent = "urgent"
)

// ValidStatuses returns all valid task statuses.
var ValidStatuses = []string{
    TaskStatusTodo,
    TaskStatusInProgress,
    TaskStatusReview,
    TaskStatusDone,
}

// ValidTransitions defines which status transitions are allowed.
var ValidTransitions = map[string][]string{
    TaskStatusTodo:       {TaskStatusInProgress},
    TaskStatusInProgress: {TaskStatusReview, TaskStatusTodo},
    TaskStatusReview:     {TaskStatusDone, TaskStatusInProgress},
    TaskStatusDone:       {TaskStatusTodo}, // reopen
}

// CanTransitionTo checks whether the task can move to the given status.
func (t *Task) CanTransitionTo(newStatus string) bool {
    allowed, ok := ValidTransitions[t.Status]
    if !ok {
        return false
    }
    for _, s := range allowed {
        if s == newStatus {
            return true
        }
    }
    return false
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now enforce the workflow in the Task service&apos;s update
                  method:
                </p>

                <CodeBlock filename="apps/api/internal/services/task.go &mdash; Update method" code={`func (s *TaskService) Update(id uint, input map[string]interface{}) (*models.Task, error) {
    task, err := s.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("task not found: %w", err)
    }

    // Validate status transition if status is being changed
    if newStatus, ok := input["status"].(string); ok && newStatus != task.Status {
        if !task.CanTransitionTo(newStatus) {
            return nil, fmt.Errorf(
                "invalid status transition: cannot move from %q to %q",
                task.Status, newStatus,
            )
        }
    }

    result := s.db.Model(task).Updates(input)
    if result.Error != nil {
        return nil, fmt.Errorf("failed to update task: %w", result.Error)
    }

    // Reload with relationships
    return s.GetByID(id)
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 8 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                8
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create a project dashboard widget
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Add a custom API endpoint that returns task statistics per
                  project, then display it as a dashboard widget in the admin
                  panel.
                </p>

                <CodeBlock filename="apps/api/internal/handlers/project.go &mdash; add stats
                      endpoint" code={`// GetStats returns task statistics grouped by status.
func (h *ProjectHandler) GetStats(c *gin.Context) {
    type StatusCount struct {
        Status string \`json:"status"\`
        Count  int64  \`json:"count"\`
    }

    var stats []StatusCount
    result := h.db.Model(&models.Task{}).
        Select("status, COUNT(*) as count").
        Group("status").
        Find(&stats)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to fetch task stats",
            },
        })
        return
    }

    // Also get total counts
    var totalTasks int64
    var totalProjects int64
    var overdueTasks int64

    h.db.Model(&models.Task{}).Count(&totalTasks)
    h.db.Model(&models.Project{}).Count(&totalProjects)
    h.db.Model(&models.Task{}).
        Where("due_date < ? AND status != ?", time.Now(), models.TaskStatusDone).
        Count(&overdueTasks)

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "total_tasks":    totalTasks,
            "total_projects": totalProjects,
            "overdue_tasks":  overdueTasks,
            "by_status":      stats,
        },
    })
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Register the route:
                </p>

                <CodeBlock filename="apps/api/internal/routes/routes.go" code={`// Inside the authenticated group
projects.GET("/stats", projectHandler.GetStats)`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now create a React component that fetches and displays these
                  stats on the admin dashboard:
                </p>

                <CodeBlock language="tsx" filename="apps/admin/components/widgets/project-stats.tsx" code={`'use client'

import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'
import { FolderKanban, CheckCircle2, AlertTriangle, ListTodo } from 'lucide-react'

interface ProjectStats {
  total_tasks: number
  total_projects: number
  overdue_tasks: number
  by_status: { status: string; count: number }[]
}

export function ProjectStatsWidget() {
  const { data } = useQuery<{ data: ProjectStats }>({
    queryKey: ['project-stats'],
    queryFn: async () => {
      const { data } = await apiClient.get('/api/projects/stats')
      return data
    },
    refetchInterval: 30000, // refresh every 30 seconds
  })

  const stats = data?.data

  const cards = [
    {
      label: 'Total Projects',
      value: stats?.total_projects ?? 0,
      icon: FolderKanban,
      color: 'text-primary',
    },
    {
      label: 'Total Tasks',
      value: stats?.total_tasks ?? 0,
      icon: ListTodo,
      color: 'text-blue-400',
    },
    {
      label: 'Completed',
      value: stats?.by_status.find((s) => s.status === 'done')?.count ?? 0,
      icon: CheckCircle2,
      color: 'text-emerald-400',
    },
    {
      label: 'Overdue',
      value: stats?.overdue_tasks ?? 0,
      icon: AlertTriangle,
      color: 'text-red-400',
    },
  ]

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card) => (
        <div
          key={card.label}
          className="rounded-xl border border-border/40 bg-card/50 p-5"
        >
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-medium text-muted-foreground/60 uppercase tracking-wider">
              {card.label}
            </span>
            <card.icon className={\`h-4 w-4 \${card.color}\`} />
          </div>
          <p className="text-2xl font-bold tracking-tight">{card.value}</p>
        </div>
      ))}
    </div>
  )
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 9 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                9
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Add email notifications when tasks are assigned
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  When a task is assigned to a user, send them an email
                  notification. This uses Grit&apos;s background job queue
                  (Redis + asynq) so the API response is not blocked by email
                  sending.
                </p>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  First, define the job type and its payload:
                </p>

                <CodeBlock filename="apps/api/internal/jobs/task_assigned.go" code={`package jobs

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hibiken/asynq"
    "taskflow/apps/api/internal/mail"
)

const TypeTaskAssigned = "task:assigned"

type TaskAssignedPayload struct {
    TaskTitle     string \`json:"task_title"\`
    ProjectName   string \`json:"project_name"\`
    AssigneeEmail string \`json:"assignee_email"\`
    AssigneeName  string \`json:"assignee_name"\`
}

// EnqueueTaskAssigned creates a new task assignment notification job.
func (c *Client) EnqueueTaskAssigned(payload TaskAssignedPayload) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    task := asynq.NewTask(TypeTaskAssigned, data)
    _, err = c.client.Enqueue(task, asynq.MaxRetry(3))
    return err
}

// HandleTaskAssigned processes the task assignment email.
func HandleTaskAssigned(mailer *mail.Mailer) asynq.HandlerFunc {
    return func(ctx context.Context, t *asynq.Task) error {
        var payload TaskAssignedPayload
        if err := json.Unmarshal(t.Payload(), &payload); err != nil {
            return fmt.Errorf("failed to unmarshal payload: %w", err)
        }

        return mailer.Send(ctx, mail.SendOptions{
            To:       payload.AssigneeEmail,
            Subject:  fmt.Sprintf("You've been assigned: %s", payload.TaskTitle),
            Template: "task-assigned",
            Data: map[string]string{
                "Name":    payload.AssigneeName,
                "Task":    payload.TaskTitle,
                "Project": payload.ProjectName,
            },
        })
    }
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now dispatch the job from the Task service whenever a task is
                  assigned:
                </p>

                <CodeBlock filename="apps/api/internal/services/task.go &mdash; add to Update
                      method" code={`// After the task is updated, check if assignee changed
if newAssigneeID, ok := input["assignee_id"]; ok && newAssigneeID != nil {
    // Reload the updated task with relationships
    updated, _ := s.GetByID(id)
    if updated.Assignee != nil {
        err := s.jobClient.EnqueueTaskAssigned(jobs.TaskAssignedPayload{
            TaskTitle:     updated.Title,
            ProjectName:   updated.Project.Name,
            AssigneeEmail: updated.Assignee.Email,
            AssigneeName:  updated.Assignee.Name,
        })
        if err != nil {
            // Log the error but don't fail the update
            log.Printf("failed to enqueue task assignment email: %v", err)
        }
    }
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Create the email template:
                </p>

                <CodeBlock language="markup" filename="apps/api/internal/mail/templates/task-assigned.html" code={`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: 'Onest', sans-serif; background: #0a0a0f; color: #e8e8f0; padding: 40px;">
  <div style="max-width: 500px; margin: 0 auto; background: #111118; border-radius: 12px; padding: 32px; border: 1px solid #2a2a3a;">
    <h2 style="color: #6c5ce7; margin-top: 0;">New Task Assigned</h2>
    <p>Hi {{.Name}},</p>
    <p>You have been assigned a new task:</p>
    <div style="background: #1a1a24; border-radius: 8px; padding: 16px; margin: 16px 0; border: 1px solid #2a2a3a;">
      <p style="margin: 0; font-weight: 600;">{{.Task}}</p>
      <p style="margin: 4px 0 0; color: #9090a8; font-size: 14px;">Project: {{.Project}}</p>
    </div>
    <p style="color: #9090a8; font-size: 14px;">Log in to Taskflow to view the details and get started.</p>
  </div>
</body>
</html>`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 10 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                10
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Customize the admin panel with badges and colors
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Update the Task resource definition to show status badges with
                  workflow colors, priority badges, assignee names, and project
                  names. Add filters for status, priority, and project.
                </p>

                <CodeBlock language="typescript" filename="apps/admin/resources/tasks.ts" code={`import { defineResource } from '@grit/admin'

export default defineResource({
  name: 'Task',
  endpoint: '/api/tasks',
  icon: 'CheckSquare',

  table: {
    columns: [
      { key: 'id', label: 'ID', sortable: true },
      { key: 'title', label: 'Title', sortable: true, searchable: true },
      { key: 'project.name', label: 'Project', relation: 'project' },
      { key: 'assignee.name', label: 'Assignee', relation: 'assignee' },
      { key: 'status', label: 'Status', sortable: true, badge: {
        todo:        { color: 'slate', label: 'Todo' },
        in_progress: { color: 'blue', label: 'In Progress' },
        review:      { color: 'yellow', label: 'Review' },
        done:        { color: 'green', label: 'Done' },
      }},
      { key: 'priority', label: 'Priority', badge: {
        low:    { color: 'slate', label: 'Low' },
        medium: { color: 'blue', label: 'Medium' },
        high:   { color: 'orange', label: 'High' },
        urgent: { color: 'red', label: 'Urgent' },
      }},
      { key: 'due_date', label: 'Due Date', format: 'date' },
      { key: 'created_at', label: 'Created', format: 'relative' },
    ],
    filters: [
      { key: 'status', type: 'select', options: [
        { label: 'Todo', value: 'todo' },
        { label: 'In Progress', value: 'in_progress' },
        { label: 'Review', value: 'review' },
        { label: 'Done', value: 'done' },
      ]},
      { key: 'priority', type: 'select', options: [
        { label: 'Low', value: 'low' },
        { label: 'Medium', value: 'medium' },
        { label: 'High', value: 'high' },
        { label: 'Urgent', value: 'urgent' },
      ]},
      { key: 'project_id', type: 'select', resource: 'projects',
        displayKey: 'name', label: 'Project' },
      { key: 'due_date', type: 'date-range' },
    ],
    actions: ['create', 'edit', 'delete'],
    bulkActions: ['delete'],
  },

  form: {
    fields: [
      { key: 'title', label: 'Title', type: 'text', required: true },
      { key: 'description', label: 'Description', type: 'textarea' },
      { key: 'project_id', label: 'Project', type: 'relation',
        resource: 'projects', displayKey: 'name', required: true },
      { key: 'assignee_id', label: 'Assignee', type: 'relation',
        resource: 'users', displayKey: 'name' },
      { key: 'status', label: 'Status', type: 'select',
        options: ['todo', 'in_progress', 'review', 'done'], default: 'todo' },
      { key: 'priority', label: 'Priority', type: 'select',
        options: ['low', 'medium', 'high', 'urgent'], default: 'medium' },
      { key: 'due_date', label: 'Due Date', type: 'date' },
    ],
  },
})`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 11 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                11
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Role-based access control
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Admins should see all tasks across all projects. Regular users
                  should only see tasks assigned to them or tasks in projects
                  they own. Add a middleware-based scope and apply it in the
                  task handler.
                </p>

                <CodeBlock filename="apps/api/internal/services/task.go &mdash; scoped query" code={`// GetAllScoped returns tasks scoped to the current user's role.
func (s *TaskService) GetAllScoped(
    userID uint, role string,
    page, pageSize int, sort, order, search string,
) ([]models.Task, int64, error) {
    var tasks []models.Task
    var total int64

    query := s.db.Model(&models.Task{})

    // Non-admin users only see their own tasks
    if role != "admin" {
        query = query.Where("assignee_id = ?", userID)
    }

    // Search
    if search != "" {
        query = query.Where("title ILIKE ?", "%"+search+"%")
    }

    query.Count(&total)

    // Sorting
    if sort != "" {
        direction := "ASC"
        if order == "desc" {
            direction = "DESC"
        }
        query = query.Order(sort + " " + direction)
    } else {
        query = query.Order("created_at DESC")
    }

    // Preload relationships
    query = query.Preload("Project").Preload("Assignee")

    // Pagination
    offset := (page - 1) * pageSize
    result := query.Offset(offset).Limit(pageSize).Find(&tasks)

    return tasks, total, result.Error
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Use the scoped query in the handler by extracting the current
                  user from the Gin context (set by the auth middleware):
                </p>

                <CodeBlock filename="apps/api/internal/handlers/task.go &mdash; GetAll handler" code={`func (h *TaskHandler) GetAll(c *gin.Context) {
    // Extract the authenticated user from context
    user, _ := c.Get("user")
    currentUser := user.(*models.User)

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    sort := c.DefaultQuery("sort", "")
    order := c.DefaultQuery("order", "")
    search := c.DefaultQuery("search", "")

    tasks, total, err := h.service.GetAllScoped(
        currentUser.ID, currentUser.Role,
        page, pageSize, sort, order, search,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
        })
        return
    }

    pages := int(total) / pageSize
    if int(total)%pageSize != 0 {
        pages++
    }

    c.JSON(http.StatusOK, gin.H{
        "data": tasks,
        "meta": gin.H{
            "total": total, "page": page,
            "page_size": pageSize, "pages": pages,
        },
    })
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 12 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                12
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Run and test everything
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Start all services and verify the complete workflow.
                </p>

                <CodeBlock terminal code="grit dev" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Test the workflow:
                </p>

                <ol className="space-y-2 text-[13px] text-muted-foreground/70 list-decimal pl-5">
                  <li>
                    Open the admin panel at <code>http://localhost:3001</code>{" "}
                    and register an admin user.
                  </li>
                  <li>
                    Create a project called &quot;Website Redesign&quot; with
                    status &quot;active&quot;.
                  </li>
                  <li>
                    Create a task titled &quot;Design homepage&quot;, assign it
                    to your user, set priority to &quot;high&quot;.
                  </li>
                  <li>
                    Check Mailhog at <code>http://localhost:8025</code> &mdash;
                    you should see the assignment notification email.
                  </li>
                  <li>
                    Update the task status from &quot;todo&quot; to
                    &quot;in_progress&quot; &mdash; it should succeed.
                  </li>
                  <li>
                    Try changing it directly to &quot;done&quot; &mdash; the API
                    should reject the invalid transition.
                  </li>
                  <li>
                    Add a comment on the task to test the commenting system.
                  </li>
                  <li>
                    Check the dashboard stats widget &mdash; it should show 1
                    project, 1 task, and the correct status breakdown.
                  </li>
                  <li>
                    Browse the database in GORM Studio at{" "}
                    <code>http://localhost:8080/studio</code>.
                  </li>
                </ol>
              </div>
            </div>

            {/* ============================================================ */}
            {/* SUMMARY */}
            {/* ============================================================ */}
            <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
              <h2 className="text-xl font-semibold tracking-tight mb-4">
                What you&apos;ve built
              </h2>
              <ul className="space-y-2.5">
                {[
                  "A project management SaaS with Projects, Tasks, and Comments",
                  "Three-level relationships: Project → Task → Comment, Task → User (assignee)",
                  "A status workflow with enforced transitions (Todo → In Progress → Review → Done)",
                  "Background job queue that sends email notifications on task assignment",
                  "HTML email templates styled to match the Grit dark theme",
                  "An admin dashboard with live task statistics (total, completed, overdue)",
                  "Admin panel with status and priority badges in color-coded columns",
                  "Role-based access: admins see all tasks, users see only their own",
                  "Server-side pagination, sorting, and filtering across all resources",
                  "Docker-based infrastructure with PostgreSQL, Redis, MinIO, and Mailhog",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-0.5">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="prose-grit mb-8">
              <h2>Next steps</h2>
              <p>
                You now have a solid project management foundation. Here are
                ideas to extend it:
              </p>
              <ul>
                <li>
                  <strong>Kanban board</strong> &mdash; use the Grit Pro
                  tier&apos;s kanban widget to display tasks in a drag-and-drop
                  board view.
                </li>
                <li>
                  <strong>Activity timeline</strong> &mdash; log every status
                  change and assignment as an &quot;Activity&quot; model and
                  display it on the task detail page.
                </li>
                <li>
                  <strong>File attachments</strong> &mdash; use the Grit storage
                  service to attach files to tasks and comments.
                </li>
                <li>
                  <strong>Recurring tasks</strong> &mdash; use the cron
                  scheduler to auto-create tasks on a schedule.
                </li>
                <li>
                  <strong>WebSocket updates</strong> &mdash; push real-time task
                  status changes to all connected team members.
                </li>
                <li>
                  <strong>Time tracking</strong> &mdash; add a TimeEntry
                  resource linked to tasks to log hours worked.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/blog" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Build a Blog
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/ecommerce" className="gap-1.5">
                  Build an E-Commerce Store
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
