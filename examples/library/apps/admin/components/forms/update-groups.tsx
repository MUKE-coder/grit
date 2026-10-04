"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import type { ResourceDefinition, FieldDefinition, GroupDefinition } from "@/lib/resource";
import { writableValues } from "@/lib/form-values";
import { FieldRenderer } from "@/components/forms/form-builder";
import { useResourceItem, usePatchResource } from "@/hooks/use-resource";
import { ChevronLeft, Loader2 } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";

interface UpdateGroupsProps {
  resource: ResourceDefinition;
  id: string;
}

export function UpdateGroups({ resource, id }: UpdateGroupsProps) {
  const router = useRouter();
  const { data: item, isLoading } = useResourceItem(resource.endpoint, id);
  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  // Only render update-applicable groups.
  const updateGroups = (resource.form.groups ?? []).filter(
    (g) => !g.scope || g.scope === "update" || g.scope === "both"
  );

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <div className="animate-pulse space-y-3">
            <div className="h-6 w-48 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
          </div>
        </div>
      </div>
    );
  }

  const record = (item?.data ?? {}) as Record<string, unknown>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.push(`/resources/${resource.slug}`)}
          className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
          Back to {pluralName}
        </button>
      </div>

      <div>
        <h1 className="text-2xl font-bold text-foreground">Edit {singularName}</h1>
        <p className="text-text-secondary mt-1">
          Each section saves independently — change what you need without affecting the rest.
        </p>
      </div>

      <div className="space-y-4">
        {updateGroups.map((group) => (
          <GroupCard
            key={group.title}
            resource={resource}
            group={group}
            record={record}
            id={id}
          />
        ))}
      </div>
    </div>
  );
}

interface GroupCardProps {
  resource: ResourceDefinition;
  group: GroupDefinition;
  record: Record<string, unknown>;
  id: string;
}

function GroupCard({ resource, group, record, id }: GroupCardProps) {
  const { mutate: patch, isPending } = usePatchResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const [isDirty, setIsDirty] = useState(false);

  // Build defaults from the record limited to this group's fields.
  const groupFields: FieldDefinition[] = resource.form.fields.filter((f) =>
    group.fields.includes(f.key)
  );
  const defaults: Record<string, unknown> = {};
  for (const f of groupFields) {
    defaults[f.key] = record[f.key] ?? "";
  }

  const {
    control,
    handleSubmit,
    watch,
    getValues,
    formState: { errors },
  } = useForm({ defaultValues: defaults });

  // Watch for changes so the Save button stays subdued until something edited.
  watch(() => {
    if (!isDirty) setIsDirty(true);
  });

  const onSave = handleSubmit((values) => {
    // Send only the values belonging to this group — that's the whole
    // point of PATCH-per-group.
    patch({ id, body: writableValues(groupFields, values) }, { onSuccess: () => setIsDirty(false) });
  });

  return (
    <section className="rounded-xl border border-border bg-bg-secondary p-6">
      <header className="mb-4">
        <h2 className="text-lg font-semibold text-foreground">{group.title}</h2>
        {group.description && (
          <p className="text-sm text-text-secondary mt-1">{group.description}</p>
        )}
      </header>

      <form onSubmit={onSave} className="space-y-4">
        {groupFields.map((field) => (
          <FieldRenderer key={field.key} field={field} control={control} errors={errors} getValues={getValues} />
        ))}

        <div className="flex justify-end pt-2">
          <button
            type="submit"
            disabled={!isDirty || isPending}
            className={buttonClasses()}
          >
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            Save {group.title.toLowerCase()}
          </button>
        </div>
      </form>
    </section>
  );
}
