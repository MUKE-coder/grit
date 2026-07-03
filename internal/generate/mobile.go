package generate

import (
	"path/filepath"
	"strings"
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
	if err := g.writeMobileHook(names); err != nil {
		return err
	}
	if err := g.writeMobileListScreen(names); err != nil {
		return err
	}
	return g.writeMobileDetailScreen(names)
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
export function use__PLURAL_PASCAL__(search = "", filters: Record<string, string> = {}) {
  return useInfiniteQuery({
    queryKey: ["__PLURAL__", { search, filters }],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => {
      const qs = new URLSearchParams({ page: String(pageParam), page_size: "20" });
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

func (g *Generator) writeMobileListScreen(names Names) error {
	title := g.mobileTitleField()
	imageExpr := g.mobileImageExpr("item")

	letterTile := `<View className="w-12 h-12 rounded-xl bg-[#6c5ce7]/12 mr-3 items-center justify-center">
          <Text className="text-[#6c5ce7] font-bold text-[16px]">{String(item.__TITLE__ || "?").charAt(0).toUpperCase()}</Text>
        </View>`
	rowImage := letterTile
	if imageExpr != "" {
		rowImage = "{" + imageExpr + " ? (\n" +
			`          <Image source={{ uri: ` + imageExpr + ` }} style={{ width: 48, height: 48, borderRadius: 12, marginRight: 12 }} contentFit="cover" />` +
			"\n        ) : (\n          " + letterTile + "\n        )}"
	}

	subtitle := g.mobileSubtitleField()
	rowSubtitle := ""
	if subtitle != "" {
		rowSubtitle = `<Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-0.5" numberOfLines={1}>
          {item.` + subtitle + ` ? String(item.` + subtitle + `) : ""}
        </Text>`
	}

	imageImport := ""
	if imageExpr != "" {
		imageImport = "import { Image } from \"expo-image\";\n"
	}

	tmpl := `import { useState } from "react";
import { View, Text, TextInput, FlatList, Pressable, ActivityIndicator, RefreshControl } from "react-native";
__IMAGE_IMPORT__import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { use__PLURAL_PASCAL__, type __PASCAL__ } from "@/hooks/use-__KEBAB__";

export default function __PLURAL_PASCAL__Screen() {
  const router = useRouter();
  const { palette } = useTheme();
  const [search, setSearch] = useState("");
  const query = use__PLURAL_PASCAL__(search);
  const items = query.data?.pages.flatMap((p) => p.data) ?? [];

  const renderItem = ({ item }: { item: __PASCAL__ }) => (
    <Pressable
      onPress={() => router.push("/__KEBAB__/" + item.id)}
      className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 mb-3"
    >
      __ROW_IMAGE__
      <View className="flex-1">
        <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white" numberOfLines={1}>
          {item.__TITLE__ ? String(item.__TITLE__) : "Untitled"}
        </Text>
        __ROW_SUBTITLE__
      </View>
      <Ionicons name="chevron-forward" size={18} color={palette.inputIcon} />
    </Pressable>
  );

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="__PLURAL_TITLE__" subtitle="Browse all __PLURAL_LOWER__" showBack />
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
      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        contentContainerStyle={{ paddingHorizontal: 24, paddingBottom: 40 }}
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
            <View className="items-center mt-16">
              <Ionicons name="folder-open-outline" size={40} color={palette.inputIcon} />
              <Text className="text-[#6B7280] dark:text-[#9090a8] mt-3">No __PLURAL_LOWER__ yet</Text>
            </View>
          )
        }
        ListFooterComponent={
          query.isFetchingNextPage ? (
            <ActivityIndicator color={palette.refresh} style={{ marginVertical: 16 }} />
          ) : null
        }
      />
    </View>
  );
}
`
	content := g.applyMobileTokens(tmpl, names)
	content = strings.ReplaceAll(content, "__IMAGE_IMPORT__", imageImport)
	content = strings.ReplaceAll(content, "__ROW_IMAGE__", rowImage)
	content = strings.ReplaceAll(content, "__ROW_SUBTITLE__", rowSubtitle)
	content = strings.ReplaceAll(content, "__TITLE__", title)

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

	tmpl := `import { View, Text, ScrollView, ActivityIndicator } from "react-native";
__IMAGE_IMPORT__import { useLocalSearchParams } from "expo-router";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { use__PASCAL__ } from "@/hooks/use-__KEBAB__";

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
  const { id } = useLocalSearchParams<{ id: string }>();
  const { palette } = useTheme();
  const { data: item, isLoading } = use__PASCAL__(id);

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
	content = strings.ReplaceAll(content, "__TITLE__", title)

	path := filepath.Join(g.mobileRoot(), "app", names.PluralKebab, "[id].tsx")
	return writeFileWithDirs(path, content)
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
