"use client";

import { useState } from "react";
import Link from "next/link";
import {
  ExternalLink,
  Globe,
  Database,
  Layers,
  Box,
  Zap,
  X,
} from "lucide-react";
import { Button } from "@/components/ui/button";

// ── Types ───────────────────────────────────────────────────────

export interface ShowcaseProject {
  name: string;
  url: string;
  description: string;
  longDescription: string;
  image: string;
  tags: string[];
  techStack: string[];
  stats: {
    tables: string;
    models: string;
    goroutines: string;
    modules: string;
    highlights: string[];
  };
  featured?: boolean;
}

// ── Grid card ───────────────────────────────────────────────────

function GridCard({
  project,
  onClick,
}: {
  project: ShowcaseProject;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="group rounded-xl border border-border/60 bg-card/50 overflow-hidden transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5 text-left w-full cursor-pointer"
    >
      <div className="aspect-[16/10] bg-accent/30 overflow-hidden">
        <img
          src={project.image}
          alt={`${project.name} screenshot`}
          className="w-full h-full object-cover object-top transition-transform duration-500 group-hover:scale-105"
        />
      </div>
      <div className="p-4">
        <div className="flex items-center justify-between mb-1.5">
          <h3 className="font-semibold text-foreground group-hover:text-primary transition-colors">
            {project.name}
          </h3>
          <ExternalLink className="h-3.5 w-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
        <p className="text-sm text-muted-foreground mb-3 line-clamp-2">
          {project.description}
        </p>
        <div className="flex flex-wrap items-center gap-2 mb-3">
          {project.tags.slice(0, 3).map((tag) => (
            <span
              key={tag}
              className="inline-flex items-center rounded-md bg-accent/60 border border-border/40 px-2 py-0.5 text-[11px] font-medium text-muted-foreground"
            >
              {tag}
            </span>
          ))}
        </div>
        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <span className="flex items-center gap-1.5">
            <Database className="h-3 w-3" />
            {project.stats.tables} tables
          </span>
          <span className="flex items-center gap-1.5">
            <Layers className="h-3 w-3" />
            {project.stats.modules} modules
          </span>
        </div>
      </div>
    </button>
  );
}

// ── Detail bottom sheet ─────────────────────────────────────────

function ProjectDetailSheet({
  project,
  onClose,
}: {
  project: ShowcaseProject;
  onClose: () => void;
}) {
  return (
    <div className="fixed inset-0 z-50">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/70 backdrop-blur-sm animate-in fade-in-0 duration-300"
        onClick={onClose}
      />

      {/* Sheet */}
      <div className="absolute inset-x-0 bottom-0 max-h-[90vh] overflow-y-auto rounded-t-2xl border-t border-border/60 bg-background shadow-2xl animate-in slide-in-from-bottom duration-300">
        {/* Drag handle */}
        <div className="sticky top-0 z-10 flex justify-center pt-3 pb-2 bg-background rounded-t-2xl">
          <div className="h-1 w-10 rounded-full bg-border/60" />
        </div>

        {/* Close button */}
        <button
          type="button"
          onClick={onClose}
          className="absolute right-4 top-4 z-20 flex h-8 w-8 items-center justify-center rounded-full bg-accent/60 border border-border/40 text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
        >
          <X className="h-4 w-4" />
        </button>

        {/* Screenshot */}
        <div className="relative aspect-[16/7] bg-accent/30 overflow-hidden">
          <img
            src={project.image}
            alt={`${project.name} screenshot`}
            className="w-full h-full object-cover object-top"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-transparent to-transparent" />
          <div className="absolute bottom-6 left-6 right-6 flex items-end justify-between">
            <div>
              <div className="flex items-center gap-3 mb-2">
                <h2 className="text-2xl md:text-3xl font-bold text-foreground">
                  {project.name}
                </h2>
                {project.featured && (
                  <span className="inline-flex items-center rounded-full bg-primary/20 border border-primary/30 px-2.5 py-0.5 text-[10px] font-mono font-medium text-primary uppercase tracking-wider">
                    Featured
                  </span>
                )}
              </div>
              <p className="text-sm text-foreground/70">
                {project.description}
              </p>
            </div>
            <Button
              size="sm"
              className="hidden md:inline-flex h-8 px-3 text-xs bg-primary/90 hover:bg-primary shrink-0"
              asChild
            >
              <Link href={project.url} target="_blank" rel="noreferrer">
                Visit Site
                <ExternalLink className="ml-1.5 h-3 w-3" />
              </Link>
            </Button>
          </div>
        </div>

        {/* Details */}
        <div className="p-6 md:p-8">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Left: Description + Tags + Modules */}
            <div className="lg:col-span-2 space-y-6">
              <p className="text-muted-foreground leading-relaxed">
                {project.longDescription}
              </p>

              <div className="flex flex-wrap gap-2">
                {project.tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center rounded-md bg-accent/60 border border-border/40 px-2.5 py-1 text-xs font-medium text-muted-foreground"
                  >
                    {tag}
                  </span>
                ))}
              </div>

              <div>
                <h3 className="text-sm font-semibold text-foreground mb-3">
                  Modules &amp; Capabilities
                </h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {project.stats.highlights.map((item) => (
                    <div
                      key={item}
                      className="flex items-center gap-2 text-sm text-muted-foreground"
                    >
                      <div className="h-1.5 w-1.5 rounded-full bg-primary/60 shrink-0" />
                      {item}
                    </div>
                  ))}
                </div>
              </div>

              {/* Mobile visit button */}
              <div className="md:hidden">
                <Button
                  size="sm"
                  className="h-9 px-4 text-sm bg-primary/90 hover:bg-primary w-full"
                  asChild
                >
                  <Link href={project.url} target="_blank" rel="noreferrer">
                    Visit {project.name}
                    <ExternalLink className="ml-1.5 h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>
            </div>

            {/* Right: Stats + Tech Stack */}
            <div className="space-y-4">
              <div className="rounded-xl border border-border/60 bg-accent/30 p-5">
                <h3 className="text-sm font-semibold text-foreground mb-4">
                  Project Scale
                </h3>
                <div className="space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                      <Database className="h-5 w-5 text-primary" />
                    </div>
                    <div>
                      <div className="text-2xl font-bold text-foreground">
                        {project.stats.tables}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Database Tables
                      </div>
                    </div>
                  </div>
                  <div className="h-px bg-border/40" />
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                      <Box className="h-5 w-5 text-primary" />
                    </div>
                    <div>
                      <div className="text-2xl font-bold text-foreground">
                        {project.stats.models}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Registered Models
                      </div>
                    </div>
                  </div>
                  {project.stats.goroutines && (
                    <>
                      <div className="h-px bg-border/40" />
                      <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-500/10 border border-amber-500/20">
                          <Zap className="h-5 w-5 text-amber-500" />
                        </div>
                        <div>
                          <div className="text-2xl font-bold text-foreground">
                            {project.stats.goroutines}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Active Goroutines
                          </div>
                        </div>
                      </div>
                    </>
                  )}
                  <div className="h-px bg-border/40" />
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                      <Layers className="h-5 w-5 text-primary" />
                    </div>
                    <div>
                      <div className="text-2xl font-bold text-foreground">
                        {project.stats.modules}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Feature Modules
                      </div>
                    </div>
                  </div>
                  <div className="h-px bg-border/40" />
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/10 border border-emerald-500/20">
                      <Globe className="h-5 w-5 text-emerald-500" />
                    </div>
                    <div>
                      <div className="text-sm font-semibold text-foreground">
                        Live in Production
                      </div>
                      <div className="text-xs text-muted-foreground">
                        <Link
                          href={project.url}
                          target="_blank"
                          rel="noreferrer"
                          className="text-primary hover:underline"
                        >
                          {project.url.replace("https://", "")}
                        </Link>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="rounded-xl border border-border/60 bg-accent/30 p-5">
                <h3 className="text-sm font-semibold text-foreground mb-2">
                  Tech Stack
                </h3>
                <div className="flex flex-wrap gap-1.5">
                  {project.techStack.map((tech) => (
                    <span
                      key={tech}
                      className="inline-flex items-center rounded-md bg-background/60 border border-border/40 px-2 py-0.5 text-[11px] font-mono text-muted-foreground"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// ── Main grid component ─────────────────────────────────────────

export function ShowcaseGrid({ projects }: { projects: ShowcaseProject[] }) {
  const [selectedProject, setSelectedProject] =
    useState<ShowcaseProject | null>(null);

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {projects.map((project) => (
          <GridCard
            key={project.name}
            project={project}
            onClick={() => setSelectedProject(project)}
          />
        ))}
      </div>

      {selectedProject && (
        <ProjectDetailSheet
          project={selectedProject}
          onClose={() => setSelectedProject(null)}
        />
      )}
    </>
  );
}
