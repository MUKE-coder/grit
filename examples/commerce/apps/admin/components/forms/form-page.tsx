"use client";

import { useLocalizedResource, useT } from "@/lib/i18n";
import { useRouter, useSearchParams } from "next/navigation";
import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "@/components/forms/form-builder";
import { useCreateResource, useUpdateResource, useResourceItem } from "@/hooks/use-resource";
import { ChevronLeft } from "@/lib/icons";

interface FormPageProps {
  resource: ResourceDefinition;
}

export function FormPage({ resource: definition }: FormPageProps) {
  const resource = useLocalizedResource(definition);
  const t = useT();
  const router = useRouter();
  const searchParams = useSearchParams();
  const editId = searchParams.get("edit");
  const isEdit = editId !== null;

  const { data: item, isLoading } = useResourceItem(
    resource.endpoint,
    editId ?? "",
    { enabled: isEdit }
  );

  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint, resource.label?.singular ?? resource.name);

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const handleSubmit = (data: Record<string, unknown>) => {
    if (isEdit && editId) {
      update(
        { id: editId, body: data },
        { onSuccess: () => router.push(`/resources/${resource.slug}`) }
      );
    } else {
      create(data, { onSuccess: () => router.push(`/resources/${resource.slug}`) });
    }
  };

  if (isEdit && isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
          >
            <ChevronLeft className="h-4 w-4" />
            {t("form.backTo", "Back to {name}", { name: pluralName })}
          </button>
        </div>
        <div className="rounded-xl border border-border bg-bg-secondary p-8">
          <div className="animate-pulse space-y-4">
            <div className="h-6 w-48 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.back()}
          className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
          {t("form.backTo", "Back to {name}", { name: pluralName })}
        </button>
      </div>

      <div>
        <h1 className="text-2xl font-bold text-foreground">
          {isEdit
            ? t("form.editTitle", "Edit {name}", { name: singularName })
            : t("form.createTitle", "Create {name}", { name: singularName })}
        </h1>
        <p className="text-text-secondary mt-1">
          {isEdit ? `Update this ${singularName.toLowerCase()}'s details` : `Add a new ${singularName.toLowerCase()} to your application`}
        </p>
      </div>

      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <FormBuilder
          form={resource.form}
          defaultValues={isEdit && item?.data ? (item.data as Record<string, unknown>) : undefined}
          onSubmit={handleSubmit}
          onCancel={() => router.back()}
          isSubmitting={isCreating || isUpdating}
          submitLabel={isEdit ? t("form.update", "Update") : t("form.create", "Create")}
        />
      </div>
    </div>
  );
}
