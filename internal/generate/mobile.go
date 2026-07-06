package generate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// This file adds Expo (mobile) code generation to `grit generate resource`.
// When the project has an apps/expo mobile app, generating a resource also
// scaffolds a typed React Query hook plus a paginated list screen and a
// read detail screen — the mobile equivalent of the web hooks + admin pages.
// expo-router is file-based, so creating the files under app/<plural>/ is all
// that's needed to register the /<plural> and /<plural>/:id routes.

func (g *Generator) mobileRoot() string {
	return filepath.Join(g.Root, "apps", "expo")
}

// hasMobileApp reports whether an Expo app exists to generate into.
func (g *Generator) hasMobileApp() bool {
	return dirExists(g.mobileRoot())
}

// writeMobileFiles generates the hook + list + detail screens for the mobile
// app. It also ensures the shared ScreenHeader component exists (older
// scaffolds predate it). Callers should guard with hasMobileApp().
func (g *Generator) writeMobileFiles(names Names) error {
	if err := g.ensureMobileScreenHeader(); err != nil {
		return err
	}
	if err := g.ensureMobileFormSheet(); err != nil {
		return err
	}
	if err := g.ensureMobileExportHelper(); err != nil {
		return err
	}
	if err := g.ensureMobileUploadHelper(); err != nil {
		return err
	}
	if err := g.ensureMobileImagePickerSheet(); err != nil {
		return err
	}
	if err := g.ensureMobileImportHelper(); err != nil {
		return err
	}
	if err := g.ensureMobileImportStore(); err != nil {
		return err
	}
	if err := g.ensureMobileImportSheet(); err != nil {
		return err
	}
	if err := g.ensureMobileImportBanner(); err != nil {
		return err
	}
	if err := g.injectMobileBanner(); err != nil {
		return err
	}
	if err := g.writeMobileHook(names); err != nil {
		return err
	}
	if err := g.writeMobileListScreen(names); err != nil {
		return err
	}
	if err := g.writeMobileDetailScreen(names); err != nil {
		return err
	}
	if err := g.writeMobileFormComponent(names); err != nil {
		return err
	}
	if err := g.writeMobileCreateScreen(names); err != nil {
		return err
	}
	if err := g.writeMobileEditScreen(names); err != nil {
		return err
	}
	return g.injectMobileResourceLink(names)
}

// injectMobileResourceLink adds a card for this resource to the More tab's
// Resources section, so its list screen is reachable from the app. It's a
// no-op if the More screen or its marker is missing (older scaffold), or if
// the link was already injected — keeping repeat generations idempotent.
func (g *Generator) injectMobileResourceLink(names Names) error {
	path := filepath.Join(g.mobileRoot(), "app", "(tabs)", "explore.tsx")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	const marker = "// grit:mobile-resources"
	route := "/" + names.PluralKebab
	if !strings.Contains(content, marker) || strings.Contains(content, `route: "`+route+`"`) {
		return nil
	}
	entry := "  { title: \"" + names.PluralPascal + "\", description: \"Browse and manage " +
		strings.ToLower(names.PluralPascal) + "\", icon: \"cube-outline\", color: \"#6c5ce7\", route: \"" + route + "\" },"
	content = strings.Replace(content, marker, marker+"\n"+entry, 1)
	return os.WriteFile(path, []byte(content), 0o644)
}

// lowerCamel converts any field name to lowerCamelCase for a JS identifier:
// "unit_price" → "unitPrice", "category_id" → "categoryId".
func lowerCamel(name string) string {
	p := toPascalCase(name)
	if p == "" {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

// mobileTitleField returns the snake_case field best used as a row/detail
// title — prefers a field literally named name/title, then the first string
// or slug field, and falls back to the id.
func (g *Generator) mobileTitleField() string {
	for _, f := range g.Definition.Fields {
		n := toSnakeCase(f.Name)
		if n == "name" || n == "title" {
			return n
		}
	}
	for _, f := range g.Definition.Fields {
		switch FieldType(f.Type) {
		case FieldString, FieldSlug:
			return toSnakeCase(f.Name)
		}
	}
	return "id"
}

// mobileSubtitleField returns a snake_case field for the list row subtitle,
// or "" when nothing suitable exists. Skips the title field and only uses
// plain text-ish fields (never relations or files).
func (g *Generator) mobileSubtitleField() string {
	title := g.mobileTitleField()
	for _, f := range g.Definition.Fields {
		n := toSnakeCase(f.Name)
		if n == title {
			continue
		}
		switch FieldType(f.Type) {
		case FieldText, FieldString, FieldRichtext:
			return n
		}
	}
	return ""
}

// mobileImageExpr returns a JS expression (relative to `item`) that resolves
// to an image URL for the resource, or "" when the resource has no image.
// Handles file/files FileRef fields and URL-shaped string fields.
func (g *Generator) mobileImageExpr(itemVar string) string {
	for _, f := range g.Definition.Fields {
		n := toSnakeCase(f.Name)
		if f.IsFile() {
			return itemVar + "." + n + "?.url"
		}
		if f.IsFiles() {
			return itemVar + "." + n + "?.[0]?.url"
		}
	}
	for _, f := range g.Definition.Fields {
		if FieldType(f.Type) == FieldString && isURLField(strings.ToLower(toSnakeCase(f.Name))) {
			return itemVar + "." + toSnakeCase(f.Name)
		}
	}
	return ""
}

func (g *Generator) mobileHasFileField() bool {
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			return true
		}
	}
	return false
}

// humanizeLabel turns a field name into a display label: "unit_price" →
// "Unit Price", "category_id" → "Category".
func humanizeLabel(name string) string {
	name = strings.TrimSuffix(toSnakeCase(name), "_id")
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		if commonInitialisms[p] {
			parts[i] = strings.ToUpper(p)
		} else {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// mobileDetailRows builds the <Row/> JSX for the detail screen, one per field.
func (g *Generator) mobileDetailRows() string {
	var b strings.Builder
	for _, f := range g.Definition.Fields {
		n := toSnakeCase(f.Name)
		var label, valueExpr string
		switch {
		case f.IsBelongsTo():
			base := strings.TrimSuffix(n, "_id")
			label = f.RelatedModelName()
			// Prefer the related record's name/title (preloaded by the API),
			// fall back to the raw foreign key.
			valueExpr = "(item." + base + " && (item." + base + ".name || item." + base + ".title)) || item." + f.FKColumnName()
		case f.IsManyToMany():
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n + "?.length ? item." + n + ".length + \" linked\" : \"—\""
		case f.IsFiles():
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n + "?.length ? item." + n + ".length + \" file(s)\" : \"—\""
		case f.IsFile():
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n + "?.name || (item." + n + " ? \"1 file\" : \"—\")"
		case f.IsStringArray():
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n + "?.join(\", \")"
		case FieldType(f.Type) == FieldBool:
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n + " ? \"Yes\" : \"No\""
		case FieldType(f.Type) == FieldDatetime || FieldType(f.Type) == FieldDate:
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n + " ? new Date(item." + n + ").toLocaleString() : \"—\""
		default:
			label = humanizeLabel(f.Name)
			valueExpr = "item." + n
		}
		b.WriteString("            <Row label=\"" + label + "\" value={" + valueExpr + "} />\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// ---- templates ---------------------------------------------------------

func (g *Generator) writeMobileHook(names Names) error {
	fileRef := ""
	if g.mobileHasFileField() {
		fileRef = "type FileRef = { url: string; name?: string; size?: number; type?: string };\n\n"
	}

	tmpl := `import { useInfiniteQuery, useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

__FILEREF__export interface __PASCAL__ {
  id: string;
__FIELDS__  created_at: string;
  updated_at: string;
}

export interface __PLURAL_PASCAL__Page {
  data: __PASCAL__[];
  meta: { total: number; page: number; page_size: number; pages: number };
}

// Paginated, infinite-scroll list. Accumulates pages — call fetchNextPage()
// when the list reaches its end. Pass equality filters (e.g. a belongs_to
// foreign key) to scope the list: use__PLURAL_PASCAL__("", { category_id: id }).
export function use__PLURAL_PASCAL__(
  search = "",
  filters: Record<string, string> = {},
  sortBy = "created_at",
  sortOrder: "asc" | "desc" = "desc",
) {
  return useInfiniteQuery({
    queryKey: ["__PLURAL__", { search, filters, sortBy, sortOrder }],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => {
      const qs = new URLSearchParams({
        page: String(pageParam),
        page_size: "20",
        sort_by: sortBy,
        sort_order: sortOrder,
      });
      if (search) qs.set("search", search);
      for (const [k, v] of Object.entries(filters)) if (v) qs.set(k, v);
      return (await api.get("/__PLURAL__?" + qs.toString())) as __PLURAL_PASCAL__Page;
    },
    getNextPageParam: (last) =>
      last.meta.page < last.meta.pages ? last.meta.page + 1 : undefined,
  });
}

export function use__PASCAL__(id: string) {
  return useQuery<__PASCAL__>({
    queryKey: ["__PLURAL__", id],
    queryFn: async () => (await api.get("/__PLURAL__/" + id)).data,
    enabled: !!id,
  });
}

export function useCreate__PASCAL__() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => api.post("/__PLURAL__", input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["__PLURAL__"] }),
  });
}

export function useUpdate__PASCAL__() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, ...input }: { id: string } & Record<string, unknown>) =>
      api.put("/__PLURAL__/" + id, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["__PLURAL__"] }),
  });
}

export function useDelete__PASCAL__() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => api.delete("/__PLURAL__/" + id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["__PLURAL__"] }),
  });
}
`
	content := g.applyMobileTokens(tmpl, names)
	content = strings.ReplaceAll(content, "__FILEREF__", fileRef)
	content = strings.ReplaceAll(content, "__FIELDS__", g.buildTSInterfaceFields())

	path := filepath.Join(g.mobileRoot(), "hooks", "use-"+names.PluralKebab+".ts")
	return writeFileWithDirs(path, content)
}

// mobileColumn describes one column of the resource table.
type mobileColumn struct {
	label    string
	cellExpr string // JS expression (relative to `item`) for the cell text
	width    int
	sortKey  string // API sort_by key; "" = not sortable
	bold     bool
}

// mobileColumnFor maps a field to its table cell expression / width / sort key.
func mobileColumnFor(f Field) (cellExpr string, width int, sortKey string) {
	n := toSnakeCase(f.Name)
	switch {
	case f.IsBelongsTo():
		base := strings.TrimSuffix(n, "_id")
		return "(item." + base + " && (item." + base + ".name || item." + base + ".title)) || item." + f.FKColumnName() + " || \"\"", 150, ""
	case FieldType(f.Type) == FieldBool:
		return "item." + n + " ? \"Yes\" : \"No\"", 90, ""
	case FieldType(f.Type) == FieldInt || FieldType(f.Type) == FieldUint || FieldType(f.Type) == FieldFloat:
		return "String(item." + n + " ?? \"\")", 110, n
	case FieldType(f.Type) == FieldDatetime || FieldType(f.Type) == FieldDate:
		return "item." + n + " ? new Date(item." + n + ").toLocaleDateString() : \"\"", 140, n
	default: // string, slug, text, richtext
		sk := ""
		if f.IsSortable() {
			sk = n
		}
		return "String(item." + n + " ?? \"\")", 150, sk
	}
}

// mobileTableColumns builds the ordered column set: the title field first
// (wider, bold), then every scalar/relation field. file/files, many_to_many
// and string_array are omitted (not cell-friendly).
func (g *Generator) mobileTableColumns() []mobileColumn {
	titleN := g.mobileTitleField()
	var cols []mobileColumn

	var titleField *Field
	for i := range g.Definition.Fields {
		if toSnakeCase(g.Definition.Fields[i].Name) == titleN {
			titleField = &g.Definition.Fields[i]
			break
		}
	}
	if titleField != nil {
		expr, _, sk := mobileColumnFor(*titleField)
		cols = append(cols, mobileColumn{label: humanizeLabel(titleField.Name), cellExpr: expr, width: 180, sortKey: sk, bold: true})
	} else {
		cols = append(cols, mobileColumn{label: "ID", cellExpr: "String(item.id).slice(0, 8)", width: 120})
	}

	for _, f := range g.Definition.Fields {
		if toSnakeCase(f.Name) == titleN {
			continue
		}
		if f.IsFileField() || f.IsManyToMany() || f.IsStringArray() {
			continue
		}
		label := humanizeLabel(f.Name)
		if f.IsBelongsTo() {
			label = f.RelatedModelName()
		}
		expr, w, sk := mobileColumnFor(f)
		cols = append(cols, mobileColumn{label: label, cellExpr: expr, width: w, sortKey: sk})
	}
	return cols
}

func (g *Generator) writeMobileListScreen(names Names) error {
	cols := g.mobileTableColumns()

	var header, row strings.Builder
	total := 0
	for _, c := range cols {
		total += c.width
		w := strconv.Itoa(c.width)
		// header cell
		if c.sortKey != "" {
			header.WriteString("            <Pressable onPress={() => onSort(\"" + c.sortKey + "\")} style={{ width: " + w + " }} className=\"px-3 py-3 flex-row items-center\">\n")
			header.WriteString("              <Text className=\"text-[12px] font-semibold text-[#6B7280] dark:text-[#9090a8]\" numberOfLines={1}>" + c.label + "</Text>\n")
			header.WriteString("              {sortBy === \"" + c.sortKey + "\" ? <Ionicons name={sortOrder === \"asc\" ? \"arrow-up\" : \"arrow-down\"} size={12} color=\"#6c5ce7\" style={{ marginLeft: 4 }} /> : null}\n")
			header.WriteString("            </Pressable>\n")
		} else {
			header.WriteString("            <View style={{ width: " + w + " }} className=\"px-3 py-3\">\n")
			header.WriteString("              <Text className=\"text-[12px] font-semibold text-[#6B7280] dark:text-[#9090a8]\" numberOfLines={1}>" + c.label + "</Text>\n")
			header.WriteString("            </View>\n")
		}
		// row cell
		textClass := "text-[14px] text-[#0F1018] dark:text-white"
		if c.bold {
			textClass = "text-[14px] font-semibold text-[#0F1018] dark:text-white"
		}
		row.WriteString("      <View style={{ width: " + w + " }} className=\"px-3 py-3\">\n")
		row.WriteString("        <Text numberOfLines={1} className=\"" + textClass + "\">{" + c.cellExpr + "}</Text>\n")
		row.WriteString("      </View>\n")
	}

	// belongs_to filters: a filter sheet with a relationship picker per FK. The
	// API already supports ?<fk>=<id>, and the list hook takes those filters.
	var filterImports, filterQueries, filterJSX strings.Builder
	seenFImport := map[string]bool{}
	hasFilters := false
	for _, f := range g.Definition.Fields {
		if !f.IsBelongsTo() {
			continue
		}
		hasFilters = true
		relNames := MakeNames(f.RelatedModelName())
		hook := "use" + relNames.PluralPascal
		imp := "@/hooks/use-" + relNames.PluralKebab
		if !seenFImport[imp] {
			filterImports.WriteString("import { " + hook + " } from \"" + imp + "\";\n")
			seenFImport[imp] = true
		}
		fk := f.FKColumnName()
		optsVar := "f" + toPascalCase(fk) + "Opts"
		qVar := "f" + toPascalCase(fk) + "Query"
		filterQueries.WriteString("  const " + qVar + " = " + hook + "();\n")
		filterQueries.WriteString("  const " + optsVar + " = " + qVar + ".data?.pages.flatMap((p: any) => p.data) ?? [];\n")
		unsel := "\"px-4 py-2 mr-2 rounded-full bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a]\""
		sel := "\"px-4 py-2 mr-2 rounded-full bg-[#6c5ce7]\""
		filterJSX.WriteString("        <Text className=\"text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2\">" + f.RelatedModelName() + "</Text>\n")
		filterJSX.WriteString("        <ScrollView horizontal showsHorizontalScrollIndicator={false} className=\"mb-4\">\n")
		filterJSX.WriteString("          <Pressable onPress={() => setFilters((f) => { const n = { ...f }; delete n." + fk + "; return n; })} className={!filters." + fk + " ? " + sel + " : " + unsel + "}>\n")
		filterJSX.WriteString("            <Text className={!filters." + fk + " ? \"text-white font-medium\" : \"text-[#0F1018] dark:text-white\"}>All</Text>\n")
		filterJSX.WriteString("          </Pressable>\n")
		filterJSX.WriteString("          {" + optsVar + ".map((opt: any) => (\n")
		filterJSX.WriteString("            <Pressable key={opt.id} onPress={() => setFilters((f) => ({ ...f, " + fk + ": opt.id }))} className={filters." + fk + " === opt.id ? " + sel + " : " + unsel + "}>\n")
		filterJSX.WriteString("              <Text className={filters." + fk + " === opt.id ? \"text-white font-medium\" : \"text-[#0F1018] dark:text-white\"}>{opt.name || opt.title || opt.id}</Text>\n")
		filterJSX.WriteString("            </Pressable>\n")
		filterJSX.WriteString("          ))}\n")
		filterJSX.WriteString("        </ScrollView>\n")
	}

	filtersArg := "{}"
	filterState, filterIcon, filterSheet := "", "", ""
	if hasFilters {
		filtersArg = "filters"
		filterState = "  const [filters, setFilters] = useState<Record<string, string>>({});\n  const [filterOpen, setFilterOpen] = useState(false);\n"
		filterIcon = "            <Pressable onPress={() => setFilterOpen(true)} hitSlop={8} className=\"mr-4\">\n" +
			"              <View>\n" +
			"                <Ionicons name=\"funnel-outline\" size={21} color=\"#6c5ce7\" />\n" +
			"                {Object.keys(filters).length > 0 ? <View className=\"absolute -top-1 -right-1 w-2.5 h-2.5 rounded-full bg-[#ff6b6b]\" /> : null}\n" +
			"              </View>\n" +
			"            </Pressable>\n"
		filterSheet = "      <FormSheet visible={filterOpen} onClose={() => setFilterOpen(false)} title=\"Filters\">\n" +
			filterJSX.String() +
			"        <Pressable onPress={() => setFilters({})} className=\"border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-full py-3 items-center mb-3\">\n" +
			"          <Text className=\"text-[#6B7280] dark:text-[#9090a8] font-semibold\">Clear all</Text>\n" +
			"        </Pressable>\n" +
			"        <Pressable onPress={() => setFilterOpen(false)} className=\"bg-[#6c5ce7] rounded-full py-4 items-center\">\n" +
			"          <Text className=\"text-white font-semibold text-[15px]\">Done</Text>\n" +
			"        </Pressable>\n" +
			"      </FormSheet>\n"
	}

	tmpl := `import { useState } from "react";
import { View, Text, TextInput, ScrollView, FlatList, Pressable, ActivityIndicator, RefreshControl, Alert } from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { FormSheet } from "@/components/ui/form-sheet";
import { useTheme } from "@/lib/theme";
import { use__PLURAL_PASCAL__, useCreate__PASCAL__, type __PASCAL__ } from "@/hooks/use-__KEBAB__";
import { __PASCAL__Form } from "@/components/resource-forms/__KEBAB__-form";
import { exportResourceCsv } from "@/lib/export";
import { ImportSheet } from "@/components/ui/import-sheet";
__FILTER_IMPORTS__
const TABLE_WIDTH = __TABLE_WIDTH__;

export default function __PLURAL_PASCAL__Screen() {
  const router = useRouter();
  const { palette } = useTheme();
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState("created_at");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");
  const [sheetOpen, setSheetOpen] = useState(false);
  const [importOpen, setImportOpen] = useState(false);
__FILTER_STATE__  const create = useCreate__PASCAL__();
__FILTER_QUERIES__  const query = use__PLURAL_PASCAL__(search, __FILTERS_ARG__, sortBy, sortOrder);
  const items = query.data?.pages.flatMap((p) => p.data) ?? [];

  const onSort = (key: string) => {
    if (sortBy === key) {
      setSortOrder((o) => (o === "asc" ? "desc" : "asc"));
    } else {
      setSortBy(key);
      setSortOrder("asc");
    }
  };

  const onExport = async () => {
    try {
      await exportResourceCsv("__PLURAL__", search ? "search=" + encodeURIComponent(search) : "");
    } catch (e: any) {
      Alert.alert("Export failed", e.message || "Please try again");
    }
  };

  const renderItem = ({ item }: { item: __PASCAL__ }) => (
    <Pressable
      onPress={() => router.push("/__KEBAB__/" + item.id)}
      className="flex-row items-center border-b border-[#E5E7EB] dark:border-[#1f1f2b] bg-white dark:bg-[#111118]"
      style={{ width: TABLE_WIDTH }}
    >
__COLUMNS_ROW__    </Pressable>
  );

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader
        title="__PLURAL_TITLE__"
        subtitle="Browse all __PLURAL_LOWER__"
        showBack
        right={
          <View className="flex-row items-center">
__FILTER_ICON__            <Pressable onPress={onExport} hitSlop={8} className="mr-4">
              <Ionicons name="download-outline" size={23} color="#6c5ce7" />
            </Pressable>
            <Pressable onPress={() => setImportOpen(true)} hitSlop={8} className="mr-4">
              <Ionicons name="cloud-upload-outline" size={23} color="#6c5ce7" />
            </Pressable>
            <Pressable onPress={() => setSheetOpen(true)} hitSlop={8}>
              <Ionicons name="add-circle" size={28} color="#6c5ce7" />
            </Pressable>
          </View>
        }
      />
      <View className="px-6 pb-3">
        <View
          className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl flex-row items-center px-4"
          style={{ height: 48 }}
        >
          <Ionicons name="search-outline" size={18} color={palette.inputIcon} />
          <TextInput
            className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
            placeholder="Search __PLURAL_LOWER__..."
            placeholderTextColor={palette.placeholder}
            value={search}
            onChangeText={setSearch}
            autoCapitalize="none"
          />
        </View>
      </View>

      <ScrollView horizontal showsHorizontalScrollIndicator={true} style={{ flex: 1 }} contentContainerStyle={{ flexGrow: 1 }}>
        <View style={{ width: TABLE_WIDTH }}>
          <View className="flex-row border-b-2 border-[#E5E7EB] dark:border-[#2a2a3a]" style={{ width: TABLE_WIDTH }}>
__COLUMNS_HEADER__          </View>
          <FlatList
            data={items}
            keyExtractor={(item) => item.id}
            renderItem={renderItem}
            style={{ flex: 1 }}
            onEndReached={() => {
              if (query.hasNextPage && !query.isFetchingNextPage) query.fetchNextPage();
            }}
            onEndReachedThreshold={0.4}
            refreshControl={
              <RefreshControl refreshing={query.isRefetching} onRefresh={query.refetch} tintColor={palette.refresh} />
            }
            ListEmptyComponent={
              query.isLoading ? (
                <ActivityIndicator color={palette.refresh} style={{ marginTop: 40 }} />
              ) : (
                <Text className="text-[#6B7280] dark:text-[#9090a8] p-6">No __PLURAL_LOWER__ yet</Text>
              )
            }
            ListFooterComponent={
              query.isFetchingNextPage ? (
                <ActivityIndicator color={palette.refresh} style={{ marginVertical: 16 }} />
              ) : null
            }
          />
        </View>
      </ScrollView>

      <FormSheet visible={sheetOpen} onClose={() => setSheetOpen(false)} title="New __SINGULAR_TITLE__">
        <__PASCAL__Form
          submitting={create.isPending}
          submitLabel="Create __SINGULAR_TITLE__"
          onSubmit={async (values) => {
            await create.mutateAsync(values);
            setSheetOpen(false);
          }}
        />
      </FormSheet>
__FILTER_SHEET__
      <ImportSheet plural="__PLURAL__" visible={importOpen} onClose={() => setImportOpen(false)} onImported={() => query.refetch()} />
    </View>
  );
}
`
	content := g.applyMobileTokens(tmpl, names)
	content = strings.ReplaceAll(content, "__TABLE_WIDTH__", strconv.Itoa(total))
	content = strings.ReplaceAll(content, "__COLUMNS_HEADER__", header.String())
	content = strings.ReplaceAll(content, "__COLUMNS_ROW__", row.String())
	content = strings.ReplaceAll(content, "__FILTER_IMPORTS__", strings.TrimRight(filterImports.String(), "\n"))
	content = strings.ReplaceAll(content, "__FILTER_STATE__", filterState)
	content = strings.ReplaceAll(content, "__FILTER_QUERIES__", filterQueries.String())
	content = strings.ReplaceAll(content, "__FILTERS_ARG__", filtersArg)
	content = strings.ReplaceAll(content, "__FILTER_ICON__", filterIcon)
	content = strings.ReplaceAll(content, "__FILTER_SHEET__", filterSheet)

	path := filepath.Join(g.mobileRoot(), "app", names.PluralKebab, "index.tsx")
	return writeFileWithDirs(path, content)
}

func (g *Generator) writeMobileDetailScreen(names Names) error {
	title := g.mobileTitleField()
	imageExpr := g.mobileImageExpr("item")

	heroImage := ""
	imageImport := ""
	if imageExpr != "" {
		imageImport = "import { Image } from \"expo-image\";\n"
		heroImage = "{" + imageExpr + " ? (\n" +
			`            <Image source={{ uri: ` + imageExpr + ` }} style={{ width: "100%", height: 200, borderRadius: 20, marginBottom: 16 }} contentFit="cover" />` +
			"\n          ) : null}"
	}

	tmpl := `import { View, Text, ScrollView, ActivityIndicator, Pressable, Alert } from "react-native";
__IMAGE_IMPORT__import { useLocalSearchParams, useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { use__PASCAL__, useDelete__PASCAL__ } from "@/hooks/use-__KEBAB__";

function Row({ label, value }: { label: string; value?: string | number | null }) {
  return (
    <View className="flex-row items-start justify-between px-5 py-4 border-b border-[#E5E7EB] dark:border-[#2a2a3a]">
      <Text className="text-[14px] text-[#6B7280] dark:text-[#9090a8]">{label}</Text>
      <Text
        className="text-[14px] text-[#0F1018] dark:text-white font-medium flex-1 text-right ml-4"
        numberOfLines={4}
      >
        {value === null || value === undefined || value === "" ? "—" : String(value)}
      </Text>
    </View>
  );
}

export default function __PASCAL__DetailScreen() {
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { palette } = useTheme();
  const { data: item, isLoading } = use__PASCAL__(id);
  const del = useDelete__PASCAL__();

  const onDelete = () => {
    Alert.alert("Delete __SINGULAR_LOWER__", "This can't be undone.", [
      { text: "Cancel", style: "cancel" },
      {
        text: "Delete",
        style: "destructive",
        onPress: async () => {
          try {
            await del.mutateAsync(id);
            router.back();
          } catch (e: any) {
            Alert.alert("Delete failed", e.message || "Please try again");
          }
        },
      },
    ]);
  };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="__SINGULAR_TITLE__" showBack />
      {isLoading || !item ? (
        <ActivityIndicator color={palette.refresh} style={{ marginTop: 40 }} />
      ) : (
        <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 48 }} showsVerticalScrollIndicator={false}>
          __HERO_IMAGE__
          <Text className="text-[22px] font-bold text-[#0F1018] dark:text-white mb-4">
            {item.__TITLE__ ? String(item.__TITLE__) : "Untitled"}
          </Text>
          <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl overflow-hidden">
__DETAIL_ROWS__
          </View>

          <Pressable
            onPress={() => router.push({ pathname: "/__KEBAB__/edit/[id]", params: { id } })}
            className="bg-[#6c5ce7] rounded-full py-4 items-center mt-6 flex-row justify-center"
          >
            <Ionicons name="create-outline" size={18} color="#fff" />
            <Text className="text-white font-semibold text-[15px] ml-2">Edit</Text>
          </Pressable>
          <Pressable
            onPress={onDelete}
            className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-full py-4 items-center mt-3 flex-row justify-center"
          >
            <Ionicons name="trash-outline" size={18} color="#ff6b6b" />
            <Text className="text-[#ff6b6b] font-semibold text-[15px] ml-2">Delete</Text>
          </Pressable>
        </ScrollView>
      )}
    </View>
  );
}
`
	content := g.applyMobileTokens(tmpl, names)
	content = strings.ReplaceAll(content, "__IMAGE_IMPORT__", imageImport)
	content = strings.ReplaceAll(content, "__HERO_IMAGE__", heroImage)
	content = strings.ReplaceAll(content, "__DETAIL_ROWS__", g.mobileDetailRows())
	content = strings.ReplaceAll(content, "__SINGULAR_LOWER__", strings.ToLower(names.Pascal))
	content = strings.ReplaceAll(content, "__TITLE__", title)

	path := filepath.Join(g.mobileRoot(), "app", names.PluralKebab, "[id].tsx")
	return writeFileWithDirs(path, content)
}

// writeMobileFormComponent generates a shared <X>Form component that both the
// create and edit screens render. It pre-fills from an optional `initial`
// record and calls onSubmit(values) — so it works unchanged inside a full page
// or (later) a bottom sheet. Fields: text/number/textarea/toggle for scalars,
// image upload for file fields, chip picker for belongs_to.
func (g *Generator) writeMobileFormComponent(names Names) error {
	var (
		extraImports strings.Builder
		stateLines   strings.Builder
		optionLines  strings.Builder
		validations  strings.Builder
		payload      strings.Builder
		fieldsJSX    strings.Builder
	)
	hasFile := false
	seenImport := map[string]bool{}

	for _, f := range g.Definition.Fields {
		t := FieldType(f.Type)
		if t == FieldSlug || t == FieldManyToMany || t == FieldStringArray {
			continue
		}
		n := toSnakeCase(f.Name)
		camel := lowerCamel(f.Name)
		pascal := toPascalCase(f.Name)
		label := humanizeLabel(f.Name)

		switch {
		case f.IsBelongsTo():
			rel := f.RelatedModelName()
			relNames := MakeNames(rel)
			hook := "use" + relNames.PluralPascal
			imp := "@/hooks/use-" + relNames.PluralKebab
			if !seenImport[imp] {
				extraImports.WriteString("import { " + hook + " } from \"" + imp + "\";\n")
				seenImport[imp] = true
			}
			fk := f.FKColumnName()
			fkCamel := lowerCamel(fk)
			fkSetter := "set" + toPascalCase(fk)
			optsVar := lowerCamel(relNames.Plural) + "Opts"
			queryVar := lowerCamel(relNames.Plural) + "Query"
			stateLines.WriteString("  const [" + fkCamel + ", " + fkSetter + "] = useState(i." + fk + " ?? \"\");\n")
			optionLines.WriteString("  const " + queryVar + " = " + hook + "();\n")
			optionLines.WriteString("  const " + optsVar + " = " + queryVar + ".data?.pages.flatMap((p: any) => p.data) ?? [];\n")
			validations.WriteString("    if (!" + fkCamel + ") return setError(\"" + label + " is required\");\n")
			payload.WriteString("        " + fk + ": " + fkCamel + ",\n")
			fieldsJSX.WriteString("      <Text className={labelClass}>" + label + "</Text>\n")
			fieldsJSX.WriteString("      <ScrollView horizontal showsHorizontalScrollIndicator={false} className=\"mb-4\">\n")
			fieldsJSX.WriteString("        {" + optsVar + ".map((opt: any) => (\n")
			fieldsJSX.WriteString("          <Pressable key={opt.id} onPress={() => " + fkSetter + "(opt.id)} className={" + fkCamel + " === opt.id ? \"px-4 py-2 mr-2 rounded-full bg-[#6c5ce7]\" : \"px-4 py-2 mr-2 rounded-full bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a]\"}>\n")
			fieldsJSX.WriteString("            <Text className={" + fkCamel + " === opt.id ? \"text-white font-medium\" : \"text-[#0F1018] dark:text-white\"}>{opt.name || opt.title || opt.id}</Text>\n")
			fieldsJSX.WriteString("          </Pressable>\n")
			fieldsJSX.WriteString("        ))}\n")
			fieldsJSX.WriteString("      </ScrollView>\n")

		case f.IsFileField():
			hasFile = true
			urlVar := camel + "Url"
			urlSetter := "set" + pascal + "Url"
			stateLines.WriteString("  const [" + urlVar + ", " + urlSetter + "] = useState<string | null>(i." + n + "?.url ?? null);\n")
			payload.WriteString("        " + n + ": " + urlVar + " ? { url: " + urlVar + " } : undefined,\n")
			fieldsJSX.WriteString("      <Text className={labelClass}>" + label + "</Text>\n")
			fieldsJSX.WriteString("      <Pressable onPress={() => openPicker(" + urlSetter + ")} className=\"mb-4 h-40 rounded-2xl border border-dashed border-[#E5E7EB] dark:border-[#2a2a3a] items-center justify-center overflow-hidden bg-white dark:bg-[#111118]\">\n")
			fieldsJSX.WriteString("        {uploading ? (\n")
			fieldsJSX.WriteString("          <ActivityIndicator color=\"#6c5ce7\" />\n")
			fieldsJSX.WriteString("        ) : " + urlVar + " ? (\n")
			fieldsJSX.WriteString("          <Image source={{ uri: " + urlVar + " }} style={{ width: \"100%\", height: \"100%\" }} contentFit=\"cover\" />\n")
			fieldsJSX.WriteString("        ) : (\n")
			fieldsJSX.WriteString("          <View className=\"items-center\">\n")
			fieldsJSX.WriteString("            <Ionicons name=\"image-outline\" size={28} color=\"#9CA3AF\" />\n")
			fieldsJSX.WriteString("            <Text className=\"text-[#6B7280] dark:text-[#9090a8] mt-2 text-[13px]\">Tap to add a photo</Text>\n")
			fieldsJSX.WriteString("          </View>\n")
			fieldsJSX.WriteString("        )}\n")
			fieldsJSX.WriteString("      </Pressable>\n")

		case t == FieldBool:
			stateLines.WriteString("  const [" + camel + ", set" + pascal + "] = useState(i." + n + " ?? false);\n")
			payload.WriteString("        " + n + ": " + camel + ",\n")
			fieldsJSX.WriteString("      <View className=\"flex-row items-center justify-between mb-4\">\n")
			fieldsJSX.WriteString("        <Text className={labelClass} style={{ marginBottom: 0 }}>" + label + "</Text>\n")
			fieldsJSX.WriteString("        <Switch value={" + camel + "} onValueChange={set" + pascal + "} trackColor={{ false: \"#D1D5DB\", true: \"#6c5ce7\" }} thumbColor=\"#ffffff\" />\n")
			fieldsJSX.WriteString("      </View>\n")

		case t == FieldInt || t == FieldUint:
			stateLines.WriteString("  const [" + camel + ", set" + pascal + "] = useState(i." + n + " != null ? String(i." + n + ") : \"\");\n")
			payload.WriteString("        " + n + ": Number(" + camel + ") || 0,\n")
			fieldsJSX.WriteString(mobileTextInput(label, camel, "set"+pascal, "numeric", false))

		case t == FieldFloat:
			stateLines.WriteString("  const [" + camel + ", set" + pascal + "] = useState(i." + n + " != null ? String(i." + n + ") : \"\");\n")
			payload.WriteString("        " + n + ": parseFloat(" + camel + ") || 0,\n")
			fieldsJSX.WriteString(mobileTextInput(label, camel, "set"+pascal, "decimal-pad", false))

		case t == FieldText || t == FieldRichtext:
			stateLines.WriteString("  const [" + camel + ", set" + pascal + "] = useState(i." + n + " ?? \"\");\n")
			payload.WriteString("        " + n + ": " + camel + ",\n")
			fieldsJSX.WriteString(mobileTextInput(label, camel, "set"+pascal, "default", true))

		case t == FieldDatetime || t == FieldDate:
			stateLines.WriteString("  const [" + camel + ", set" + pascal + "] = useState(i." + n + " ?? \"\");\n")
			payload.WriteString("        " + n + ": " + camel + " || undefined,\n")
			fieldsJSX.WriteString(mobileTextInput(label, camel, "set"+pascal, "default", false))

		default: // string
			stateLines.WriteString("  const [" + camel + ", set" + pascal + "] = useState(i." + n + " ?? \"\");\n")
			payload.WriteString("        " + n + ": " + camel + ",\n")
			fieldsJSX.WriteString(mobileTextInput(label, camel, "set"+pascal, "default", false))
		}
	}

	fileImports := ""
	pickHandler := ""
	pickSheet := ""
	if hasFile {
		fileImports = "import { useRef } from \"react\";\nimport { Image } from \"expo-image\";\nimport { uploadLocalFile } from \"@/lib/upload\";\nimport { ImagePickerSheet } from \"@/components/ui/image-picker-sheet\";\n"
		pickHandler = `
  // A single picker sheet serves every image field: openPicker points it at the
  // tapped field's setter, then the chosen local image is uploaded and stored.
  const [pickerOpen, setPickerOpen] = useState(false);
  const [uploading, setUploading] = useState(false);
  const pickerTarget = useRef<((u: string) => void) | null>(null);

  const openPicker = (setter: (u: string) => void) => {
    pickerTarget.current = setter;
    setPickerOpen(true);
  };

  const onImagesSelected = async (uris: string[]) => {
    setPickerOpen(false);
    const setter = pickerTarget.current;
    if (!uris.length || !setter) return;
    setUploading(true);
    try {
      const url = await uploadLocalFile(uris[0]);
      setter(url);
    } catch (e: any) {
      Alert.alert("Upload failed", e.message || "Please try again");
    } finally {
      setUploading(false);
    }
  };
`
		pickSheet = "      <ImagePickerSheet visible={pickerOpen} onClose={() => setPickerOpen(false)} onSelect={onImagesSelected} />\n"
	}

	tmpl := `import { useState } from "react";
import { View, Text, TextInput, ScrollView, Pressable, Switch, ActivityIndicator, Alert } from "react-native";
import { Ionicons } from "@expo/vector-icons";
__FILE_IMPORTS____EXTRA_IMPORTS__
export interface __PASCAL__FormProps {
  initial?: Record<string, any>;
  onSubmit: (values: Record<string, unknown>) => Promise<void> | void;
  submitting?: boolean;
  submitLabel?: string;
}

// Shared create/edit form. Renders inside a page or a bottom sheet; the parent
// owns the mutation and navigation via onSubmit.
export function __PASCAL__Form({ initial, onSubmit, submitting, submitLabel }: __PASCAL__FormProps) {
  const i: any = initial || {};
  const [error, setError] = useState("");
__STATE__
__OPTIONS__
  const inputClass =
    "bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-3.5 text-[#0F1018] dark:text-white text-[15px] mb-4";
  const labelClass = "text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2";
__PICK_HANDLER__
  const submit = async () => {
    setError("");
__VALIDATIONS__
    try {
      await onSubmit({
__PAYLOAD__      });
    } catch (e: any) {
      setError(e.message || "Something went wrong");
    }
  };

  return (
    <View>
      {error ? (
        <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center">
          <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
          <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{error}</Text>
        </View>
      ) : null}

__FIELDS__
      <Pressable
        onPress={submit}
        disabled={submitting}
        className="bg-[#6c5ce7] rounded-full py-4 items-center mt-2"
        style={{ opacity: submitting ? 0.7 : 1 }}
      >
        {submitting ? (
          <ActivityIndicator color="#fff" />
        ) : (
          <Text className="text-white font-semibold text-[15px]">{submitLabel || "Save"}</Text>
        )}
      </Pressable>
__PICK_SHEET__    </View>
  );
}
`
	content := g.applyMobileTokens(tmpl, names)
	content = strings.ReplaceAll(content, "__PICK_SHEET__", pickSheet)
	content = strings.ReplaceAll(content, "__FILE_IMPORTS__", fileImports)
	content = strings.ReplaceAll(content, "__EXTRA_IMPORTS__", extraImports.String())
	content = strings.ReplaceAll(content, "__STATE__", strings.TrimRight(stateLines.String(), "\n"))
	content = strings.ReplaceAll(content, "__OPTIONS__", strings.TrimRight(optionLines.String(), "\n"))
	content = strings.ReplaceAll(content, "__PICK_HANDLER__", pickHandler)
	content = strings.ReplaceAll(content, "__VALIDATIONS__", strings.TrimRight(validations.String(), "\n"))
	content = strings.ReplaceAll(content, "__PAYLOAD__", payload.String())
	content = strings.ReplaceAll(content, "__FIELDS__", strings.TrimRight(fieldsJSX.String(), "\n"))

	path := filepath.Join(g.mobileRoot(), "components", "resource-forms", names.PluralKebab+"-form.tsx")
	return writeFileWithDirs(path, content)
}

func (g *Generator) writeMobileCreateScreen(names Names) error {
	tmpl := `import { View, ScrollView } from "react-native";
import { useRouter } from "expo-router";
import { ScreenHeader } from "@/components/ui/screen-header";
import { __PASCAL__Form } from "@/components/resource-forms/__KEBAB__-form";
import { useCreate__PASCAL__ } from "@/hooks/use-__KEBAB__";

export default function Create__PASCAL__Screen() {
  const router = useRouter();
  const create = useCreate__PASCAL__();
  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="New __SINGULAR_TITLE__" showBack />
      <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 48 }} keyboardShouldPersistTaps="handled">
        <__PASCAL__Form
          submitting={create.isPending}
          submitLabel="Create __SINGULAR_TITLE__"
          onSubmit={async (values) => {
            await create.mutateAsync(values);
            router.back();
          }}
        />
      </ScrollView>
    </View>
  );
}
`
	content := g.applyMobileTokens(tmpl, names)
	path := filepath.Join(g.mobileRoot(), "app", names.PluralKebab, "new.tsx")
	return writeFileWithDirs(path, content)
}

func (g *Generator) writeMobileEditScreen(names Names) error {
	tmpl := `import { View, ScrollView, ActivityIndicator } from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { ScreenHeader } from "@/components/ui/screen-header";
import { __PASCAL__Form } from "@/components/resource-forms/__KEBAB__-form";
import { use__PASCAL__, useUpdate__PASCAL__ } from "@/hooks/use-__KEBAB__";

export default function Edit__PASCAL__Screen() {
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: item, isLoading } = use__PASCAL__(id);
  const update = useUpdate__PASCAL__();

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Edit __SINGULAR_TITLE__" showBack />
      {isLoading || !item ? (
        <ActivityIndicator color="#6c5ce7" style={{ marginTop: 40 }} />
      ) : (
        <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 48 }} keyboardShouldPersistTaps="handled">
          <__PASCAL__Form
            initial={item}
            submitting={update.isPending}
            submitLabel="Save changes"
            onSubmit={async (values) => {
              await update.mutateAsync({ id, ...values });
              router.back();
            }}
          />
        </ScrollView>
      )}
    </View>
  );
}
`
	content := g.applyMobileTokens(tmpl, names)
	path := filepath.Join(g.mobileRoot(), "app", names.PluralKebab, "edit", "[id].tsx")
	return writeFileWithDirs(path, content)
}

// mobileTextInput builds a labelled TextInput block for the create form.
func mobileTextInput(label, valueVar, setter, keyboard string, multiline bool) string {
	extra := ""
	if keyboard != "default" {
		extra += " keyboardType=\"" + keyboard + "\""
	}
	if multiline {
		extra += " multiline numberOfLines={4} style={{ minHeight: 96, textAlignVertical: \"top\" }}"
	}
	return "        <Text className={labelClass}>" + label + "</Text>\n" +
		"        <TextInput className={inputClass} placeholder=\"" + label + "\" placeholderTextColor=\"#9CA3AF\" value={" + valueVar + "} onChangeText={" + setter + "}" + extra + " />\n"
}

// applyMobileTokens replaces the shared name tokens used by every mobile
// template. Field/image tokens are handled per-template by the callers.
func (g *Generator) applyMobileTokens(tmpl string, names Names) string {
	r := strings.NewReplacer(
		"__PLURAL_PASCAL__", names.PluralPascal,
		"__PLURAL_TITLE__", names.PluralPascal,
		"__PLURAL_LOWER__", strings.ToLower(names.PluralPascal),
		"__PLURAL__", names.Plural,
		"__KEBAB__", names.PluralKebab,
		"__SINGULAR_TITLE__", names.Pascal,
		"__PASCAL__", names.Pascal,
	)
	return r.Replace(tmpl)
}

// ensureMobileScreenHeader writes the shared ScreenHeader component if it is
// missing. Newer scaffolds ship it in the base app; this covers projects
// generated before it existed. It never overwrites a customised one.
func (g *Generator) ensureMobileScreenHeader() error {
	path := filepath.Join(g.mobileRoot(), "components", "ui", "screen-header.tsx")
	if fileExists(path) {
		return nil
	}
	return writeFileWithDirs(path, mobileScreenHeaderContent())
}

// ensureMobileExportHelper writes lib/export.ts (CSV export → share sheet) if
// missing. Written once; never overwrites a customised copy.
func (g *Generator) ensureMobileExportHelper() error {
	path := filepath.Join(g.mobileRoot(), "lib", "export.ts")
	if fileExists(path) {
		return nil
	}
	return writeFileWithDirs(path, mobileExportHelperContent())
}

// mobileExportHelperContent is the source for lib/export.ts. Kept in sync with
// the base scaffold's expoExportHelper.
func mobileExportHelperContent() string {
	return `import * as FileSystem from "expo-file-system/legacy";
import * as Sharing from "expo-sharing";
import * as SecureStore from "expo-secure-store";
import { API_URL } from "./api";

// Download the resource's CSV export and open the share sheet.
export async function exportResourceCsv(plural: string, query = ""): Promise<string> {
  const token = await SecureStore.getItemAsync("access_token");
  const url = API_URL + "/" + plural + "/export" + (query ? "?" + query : "");
  const fileUri = FileSystem.cacheDirectory + plural + "-export.csv";
  const res = await FileSystem.downloadAsync(url, fileUri, {
    headers: token ? { Authorization: "Bearer " + token } : {},
  });
  if (res.status < 200 || res.status >= 300) {
    throw new Error("Export failed (" + res.status + ")");
  }
  if (await Sharing.isAvailableAsync()) {
    await Sharing.shareAsync(res.uri, { mimeType: "text/csv", dialogTitle: "Export " + plural });
  }
  return res.uri;
}
`
}

// The CSV import stack — lib/import.ts (upload + poll), lib/import-progress.ts
// (background store), components/ui/import-sheet.tsx (pick→preview→progress) and
// components/ui/import-progress-banner.tsx (persistent banner) — is framework
// plumbing that must move together: the sheet imports the store, the store polls
// via the helper. So unlike the ensure-once components, these are OVERWRITTEN to
// the current scaffold source on every generate, which upgrades older projects
// to the background-import flow and never leaves a half-old, non-building set.
func (g *Generator) ensureMobileImportHelper() error {
	path := filepath.Join(g.mobileRoot(), "lib", "import.ts")
	return writeFileWithDirs(path, scaffold.ExpoImportHelper())
}

// ensureMobileUploadHelper / ensureMobileImagePickerSheet keep the image-upload
// stack current: the generated form imports uploadLocalFile from lib/upload.ts
// and the ImagePickerSheet component, so both must exist at the versions the
// form expects. Overwritten on every generate (framework plumbing) so older
// projects gain the fetch+FormData upload and the permission-aware picker.
func (g *Generator) ensureMobileUploadHelper() error {
	path := filepath.Join(g.mobileRoot(), "lib", "upload.ts")
	return writeFileWithDirs(path, scaffold.ExpoUploadHelper())
}

func (g *Generator) ensureMobileImagePickerSheet() error {
	path := filepath.Join(g.mobileRoot(), "components", "ui", "image-picker-sheet.tsx")
	return writeFileWithDirs(path, scaffold.ExpoImagePickerSheet())
}

func (g *Generator) ensureMobileImportStore() error {
	path := filepath.Join(g.mobileRoot(), "lib", "import-progress.ts")
	return writeFileWithDirs(path, scaffold.ExpoImportProgressStore())
}

func (g *Generator) ensureMobileImportSheet() error {
	path := filepath.Join(g.mobileRoot(), "components", "ui", "import-sheet.tsx")
	return writeFileWithDirs(path, scaffold.ExpoImportSheet())
}

func (g *Generator) ensureMobileImportBanner() error {
	path := filepath.Join(g.mobileRoot(), "components", "ui", "import-progress-banner.tsx")
	return writeFileWithDirs(path, scaffold.ExpoImportBanner())
}

// injectMobileBanner mounts <ImportProgressBanner /> in the root layout so the
// background-import progress persists across navigation. Idempotent: a no-op if
// the layout is missing (older scaffold with a different shape) or already wired.
func (g *Generator) injectMobileBanner() error {
	path := filepath.Join(g.mobileRoot(), "app", "_layout.tsx")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)
	if strings.Contains(content, "ImportProgressBanner") {
		return nil
	}
	const anchor = "      </Stack>"
	if !strings.Contains(content, anchor) {
		return nil
	}
	imp := "import { ImportProgressBanner } from \"@/components/ui/import-progress-banner\";\n"
	// Add the import after the last existing import line at the top of the file.
	if idx := strings.LastIndex(content, "\nimport "); idx != -1 {
		end := strings.Index(content[idx+1:], "\n")
		if end != -1 {
			pos := idx + 1 + end + 1
			content = content[:pos] + imp + content[pos:]
		}
	}
	mount := anchor + "\n      {/* grit:mobile-banner */}\n      <ImportProgressBanner />"
	content = strings.Replace(content, anchor, mount, 1)
	return os.WriteFile(path, []byte(content), 0o644)
}

// ensureMobileFormSheet writes the shared FormSheet bottom-sheet component if
// it is missing (the list screen's quick-create sheet renders inside it).
// Written once; never overwrites a customised copy.
func (g *Generator) ensureMobileFormSheet() error {
	path := filepath.Join(g.mobileRoot(), "components", "ui", "form-sheet.tsx")
	if fileExists(path) {
		return nil
	}
	return writeFileWithDirs(path, mobileFormSheetContent())
}

// mobileFormSheetContent is the source for the shared bottom sheet. Built on
// React Native's Modal (no extra deps) — a themed, keyboard-aware sheet that
// slides up with a backdrop and drag handle. Kept in sync with the base
// scaffold's expoFormSheet.
func mobileFormSheetContent() string {
	return `import type { ReactNode } from "react";
import { Modal, View, Text, Pressable, ScrollView, KeyboardAvoidingView, Platform } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useSafeAreaInsets } from "react-native-safe-area-context";

interface FormSheetProps {
  visible: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
}

// Bottom sheet for quick create/edit. Renders the shared resource form inside
// a slide-up sheet; a full-page form is used for longer records.
export function FormSheet({ visible, onClose, title, children }: FormSheetProps) {
  const insets = useSafeAreaInsets();
  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose} statusBarTranslucent>
      <View className="flex-1 justify-end" style={{ backgroundColor: "rgba(0,0,0,0.4)" }}>
        <Pressable className="flex-1" onPress={onClose} />
        <KeyboardAvoidingView behavior={Platform.OS === "ios" ? "padding" : undefined}>
          <View className="bg-[#F4F4F6] dark:bg-[#0a0a0f] rounded-t-[28px] overflow-hidden" style={{ maxHeight: "88%" }}>
            <View className="items-center pt-3 pb-1">
              <View className="w-10 h-1.5 rounded-full bg-[#D1D5DB] dark:bg-[#2a2a3a]" />
            </View>
            <View className="flex-row items-center justify-between px-6 pb-3">
              <Text className="text-[20px] font-bold text-[#0F1018] dark:text-white">{title}</Text>
              <Pressable
                onPress={onClose}
                hitSlop={10}
                className="w-8 h-8 rounded-full items-center justify-center bg-white dark:bg-[#1a1a24] border border-[#E5E7EB] dark:border-[#2a2a3a]"
              >
                <Ionicons name="close" size={18} color="#9CA3AF" />
              </Pressable>
            </View>
            <ScrollView
              contentContainerStyle={{ paddingHorizontal: 24, paddingBottom: insets.bottom + 24 }}
              keyboardShouldPersistTaps="handled"
              showsVerticalScrollIndicator={false}
            >
              {children}
            </ScrollView>
          </View>
        </KeyboardAvoidingView>
      </View>
    </Modal>
  );
}
`
}

// mobileScreenHeaderContent is the source for the shared safe-area header
// with an optional back button. Kept in sync with the copy shipped by the
// base scaffold (internal/scaffold expoScreenHeader).
func mobileScreenHeaderContent() string {
	return `import { View, Text, Pressable } from "react-native";
import type { ReactNode } from "react";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useTheme } from "@/lib/theme";

interface ScreenHeaderProps {
  title: string;
  subtitle?: string;
  showBack?: boolean;
  right?: ReactNode;
}

// Safe-area-aware page header. Non-tab screens pass showBack to get a back
// button; tab screens can use it without one for a consistent large title.
export function ScreenHeader({ title, subtitle, showBack = false, right }: ScreenHeaderProps) {
  const insets = useSafeAreaInsets();
  const router = useRouter();
  const { palette } = useTheme();
  return (
    <View style={{ paddingTop: insets.top + 8 }} className="px-6 pb-3 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <View className="flex-row items-center">
        {showBack ? (
          <Pressable
            onPress={() => router.back()}
            hitSlop={10}
            className="mr-3 w-9 h-9 rounded-full items-center justify-center bg-white dark:bg-[#1a1a24] border border-[#E5E7EB] dark:border-[#2a2a3a]"
          >
            <Ionicons name="chevron-back" size={20} color={palette.inputIcon} />
          </Pressable>
        ) : null}
        <View className="flex-1">
          <Text className="text-[26px] font-bold text-[#0F1018] dark:text-white tracking-tight">{title}</Text>
          {subtitle ? (
            <Text className="text-[14px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">{subtitle}</Text>
          ) : null}
        </View>
        {right ? <View className="ml-3">{right}</View> : null}
      </View>
    </View>
  );
}
`
}
