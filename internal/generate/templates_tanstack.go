package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// writeReactQueryHooksTanStack writes hooks to an arbitrary directory (for TanStack projects).
func (g *Generator) writeReactQueryHooksTanStack(names Names, hooksDir string) error {
	content := fmt.Sprintf(`import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface %s {
  id: number
%s
  created_at: string
  updated_at: string
}

export function use%s() {
  return useQuery({
    queryKey: ['%s'],
    queryFn: async () => {
      const res = await api.get('/api/%s')
      return res.data.data as %s[]
    },
  })
}

export function use%s(id: number) {
  return useQuery({
    queryKey: ['%s', id],
    queryFn: async () => {
      const res = await api.get('/api/%s/' + id)
      return res.data.data as %s
    },
    enabled: !!id,
  })
}

export function useCreate%s() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: Partial<%s>) => {
      const res = await api.post('/api/%s', data)
      return res.data.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['%s'] })
    },
  })
}

export function useUpdate%s() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, data }: { id: number; data: Partial<%s> }) => {
      const res = await api.put('/api/%s/' + id, data)
      return res.data.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['%s'] })
    },
  })
}

export function useDelete%s() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete('/api/%s/' + id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['%s'] })
    },
  })
}
`,
		names.Pascal,
		g.buildTSFields(names),
		names.PluralPascal, names.Plural, names.Plural, names.Pascal,
		names.Pascal, names.Plural, names.Plural, names.Pascal,
		names.Pascal, names.Pascal, names.Plural, names.Plural,
		names.Pascal, names.Pascal, names.Plural, names.Plural,
		names.Pascal, names.Plural, names.Plural,
	)

	return os.WriteFile(filepath.Join(hooksDir, fmt.Sprintf("use-%s.ts", names.PluralKebab)), []byte(content), 0644)
}

// buildTSFields generates TypeScript interface fields from the resource definition.
func (g *Generator) buildTSFields(names Names) string {
	var b strings.Builder
	for _, f := range g.Definition.Fields {
		jsonName := toSnakeCase(f.Name)
		tsType := goTypeToTS(f.GoType())
		b.WriteString(fmt.Sprintf("  %s: %s\n", jsonName, tsType))
	}
	return b.String()
}

// writeResourceDefinitionTanStack writes a resource definition to apps/admin/src/resources/.
func (g *Generator) writeResourceDefinitionTanStack(names Names) error {
	content := g.resourceDefinitionContent(names)
	dir := filepath.Join(g.Root, "apps", "admin", "src", "resources")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, names.PluralKebab+".ts"), []byte(content), 0644)
}

// writeResourcePageTanStack writes a TanStack Router resource page.
func (g *Generator) writeResourcePageTanStack(names Names) error {
	content := fmt.Sprintf(`import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { %sResource } from '@/resources/%s'

export const Route = createFileRoute('/_dashboard/resources/%s')({
  component: () => <ResourcePage resource={%sResource} />,
})
`, names.Camel, names.PluralKebab, names.PluralKebab, names.Camel)

	dir := filepath.Join(g.Root, "apps", "admin", "src", "routes", "_dashboard", "resources")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, names.PluralKebab+".tsx"), []byte(content), 0644)
}

// resourceDefinitionContent generates the resource definition (shared between Next.js and TanStack).
func (g *Generator) resourceDefinitionContent(names Names) string {
	// Build column definitions
	var columns strings.Builder
	for _, f := range g.Definition.Fields {
		jsonName := toSnakeCase(f.Name)
		label := toPascalCase(f.Name)
		columns.WriteString(fmt.Sprintf(`    { key: '%s', label: '%s', sortable: true },
`, jsonName, label))
	}

	// Build form field definitions
	var fields strings.Builder
	for _, f := range g.Definition.Fields {
		jsonName := toSnakeCase(f.Name)
		if f.IsSlug() {
			continue
		}
		label := toPascalCase(f.Name)
		fieldType := f.FormFieldType()
		fields.WriteString(fmt.Sprintf(`    { name: '%s', label: '%s', type: '%s' },
`, jsonName, label, fieldType))
	}

	return fmt.Sprintf(`import { defineResource } from '@/lib/resource'

export const %sResource = defineResource({
  name: '%s',
  plural: '%s',
  apiEndpoint: '/api/%s',
  columns: [
%s  ],
  fields: [
%s  ],
})
`, names.Camel, names.Pascal, names.PluralPascal, names.Plural, columns.String(), fields.String())
}
