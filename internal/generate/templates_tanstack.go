package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"strings"
)

// writeReactQueryHooksTanStack writes hooks to an arbitrary directory (for TanStack projects).
func (g *Generator) writeReactQueryHooksTanStack(names Names, hooksDir string) error {
	content := fmt.Sprintf(`import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface %s {
  id: string
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

export function use%s(id: string) {
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
    mutationFn: async ({ id, data }: { id: string; data: Partial<%s> }) => {
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
    mutationFn: async (id: string) => {
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
//
// Uses the same builder as the Next.js admin — both consume the identical
// lib/resource.ts defineResource(), so the content must not diverge.
func (g *Generator) writeResourceDefinitionTanStack(names Names) error {
	content := g.resourceDefinitionFileContent(names)
	root := g.tanStackResourcesRoot()

	// A folder per resource, matching the Next generator and the built-in
	// blogs and users. It used to write a flat src/resources/<kebab>.ts while
	// both route writers imported '@/resources/<kebab>/<kebab>', so every
	// generated resource failed to typecheck on the line reaching its own
	// definition.
	//
	// A project still holding a flat file keeps it where it is: writing a
	// second copy in a folder would leave two definitions and a registry
	// pointing at whichever one it found first.
	if flat := filepath.Join(root, names.PluralKebab+".ts"); fileExists(flat) {
		if err := writeFileWithDirs(flat, content); err != nil {
			return err
		}
		return writeResourceCustomStub(root, names)
	}

	dir, path := scaffold.ResourceDefPath(root, names.PluralKebab)
	if err := writeFileWithDirs(path, content); err != nil {
		return err
	}
	return writeResourceCustomStub(dir, names)
}

// tanStackResourcesRoot is where the TanStack admin keeps resource
// definitions.
func (g *Generator) tanStackResourcesRoot() string {
	return filepath.Join(g.Root, "apps", "admin", "src", "resources")
}

// tanStackResourceImport is the module specifier a route uses to reach a
// resource definition.
//
// Asked rather than assumed, because the two layouts need different specifiers
// and a project can be in either: the folder form for anything generated now,
// the flat form for one written before that was fixed. Guessing is what broke
// this in the first place.
func (g *Generator) tanStackResourceImport(names Names) string {
	root := g.tanStackResourcesRoot()
	if fileExists(filepath.Join(root, names.PluralKebab+".ts")) {
		return "@/resources/" + names.PluralKebab
	}
	return "@/resources/" + names.PluralKebab + "/" + names.PluralKebab
}

// writeResourcePageTanStack writes a TanStack Router resource list page. It uses
// the <kebab>/index.tsx form (not a flat <kebab>.tsx) so a sibling $id detail
// route can live under the same directory — mirroring the built-in blogs.
func (g *Generator) writeResourcePageTanStack(names Names) error {
	content := fmt.Sprintf(`import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { %sResource } from '%s'

export const Route = createFileRoute('/_dashboard/resources/%s/')({
  component: () => <ResourcePage resource={%sResource} />,
})
`, names.Camel, g.tanStackResourceImport(names), names.PluralKebab, names.Camel)

	dir := filepath.Join(g.Root, "apps", "admin", "src", "routes", "_dashboard", "resources", names.PluralKebab)
	os.MkdirAll(dir, 0755)
	// Remove any stale flat route from a pre-detail-page generation so the two
	// don't collide on the same path.
	os.Remove(filepath.Join(g.Root, "apps", "admin", "src", "routes", "_dashboard", "resources", names.PluralKebab+".tsx"))
	return os.WriteFile(filepath.Join(dir, "index.tsx"), []byte(content), 0644)
}

// writeResourceDetailPageTanStack writes the per-resource $id detail route.
func (g *Generator) writeResourceDetailPageTanStack(names Names) error {
	content := fmt.Sprintf(`import { createFileRoute } from '@tanstack/react-router'
import { ResourceDetailPage } from '@/components/resource/resource-detail-page'
import { %sResource } from '%s'

export const Route = createFileRoute('/_dashboard/resources/%s/$id')({
  component: RouteComponent,
})

function RouteComponent() {
  const { id } = Route.useParams()
  return <ResourceDetailPage resource={%sResource} id={id} />
}
`, names.Camel, g.tanStackResourceImport(names), names.PluralKebab, names.Camel)

	dir := filepath.Join(g.Root, "apps", "admin", "src", "routes", "_dashboard", "resources", names.PluralKebab)
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "$id.tsx"), []byte(content), 0644)
}
