package scaffold

// Built-in mobile resource screens for the Blog feature and User creation.
// These mirror the generated-resource screens (card list + page form) but are
// backed by the scaffold's built-in Blog and User admin endpoints, so a
// --mobile project can browse seeded blogs and add users/blogs out of the box.

func expoBlogHook() string {
	return `import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

export interface Blog {
  id: string;
  title: string;
  slug: string;
  content: string;
  image: string;
  excerpt: string;
  published: boolean;
  created_at: string;
  updated_at: string;
}

export interface BlogsPage {
  data: Blog[];
  meta: { total: number; page: number; page_size: number; pages: number };
}

// Admin blog list — the seeded admin account can browse every post (drafts
// included). Public published posts are also available at GET /blogs.
export function useBlogs(search = "") {
  return useInfiniteQuery({
    queryKey: ["blogs", { search }],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => {
      const qs = new URLSearchParams({ page: String(pageParam), page_size: "20" });
      if (search) qs.set("search", search);
      return (await api.get("/admin/blogs?" + qs.toString())) as BlogsPage;
    },
    getNextPageParam: (last) =>
      last.meta && last.meta.page < last.meta.pages ? last.meta.page + 1 : undefined,
  });
}

export function useCreateBlog() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: Record<string, unknown>) => api.post("/admin/blogs", input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["blogs"] }),
  });
}
`
}

func expoBlogListScreen() string {
	return `import { useState } from "react";
import { View, Text, TextInput, FlatList, Pressable, ActivityIndicator, RefreshControl } from "react-native";
import { Image } from "expo-image";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useBlogs, type Blog } from "@/hooks/use-blogs";

export default function BlogsScreen() {
  const router = useRouter();
  const [search, setSearch] = useState("");
  const query = useBlogs(search);
  const items = query.data?.pages.flatMap((p) => p.data) ?? [];

  const renderItem = ({ item }: { item: Blog }) => (
    <View className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 mb-3">
      {item.image ? (
        <Image source={{ uri: item.image }} style={{ width: 48, height: 48, borderRadius: 12, marginRight: 12 }} contentFit="cover" />
      ) : (
        <View className="w-12 h-12 rounded-xl bg-[#6c5ce7]/12 mr-3 items-center justify-center">
          <Ionicons name="newspaper-outline" size={20} color="#6c5ce7" />
        </View>
      )}
      <View className="flex-1">
        <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white" numberOfLines={1}>{item.title}</Text>
        <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-0.5" numberOfLines={1}>{item.excerpt}</Text>
      </View>
      <View style={{ backgroundColor: (item.published ? "#00b894" : "#9CA3AF") + "22" }} className="px-2.5 py-1 rounded-full">
        <Text style={{ color: item.published ? "#00b894" : "#9CA3AF" }} className="text-[11px] font-semibold">
          {item.published ? "Live" : "Draft"}
        </Text>
      </View>
    </View>
  );

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader
        title="Blogs"
        subtitle="Posts and articles"
        showBack
        right={
          <Pressable onPress={() => router.push("/blogs/new")} hitSlop={8}>
            <Ionicons name="add-circle" size={28} color="#6c5ce7" />
          </Pressable>
        }
      />
      <View className="px-6 pb-3">
        <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl flex-row items-center px-4" style={{ height: 48 }}>
          <Ionicons name="search-outline" size={18} color="#9CA3AF" />
          <TextInput
            className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
            placeholder="Search posts..."
            placeholderTextColor="#9CA3AF"
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
        onEndReached={() => { if (query.hasNextPage && !query.isFetchingNextPage) query.fetchNextPage(); }}
        onEndReachedThreshold={0.4}
        refreshControl={<RefreshControl refreshing={query.isRefetching} onRefresh={query.refetch} tintColor="#6c5ce7" />}
        ListEmptyComponent={
          query.isLoading ? (
            <ActivityIndicator color="#6c5ce7" style={{ marginTop: 40 }} />
          ) : (
            <View className="items-center mt-16">
              <Ionicons name="newspaper-outline" size={40} color="#9CA3AF" />
              <Text className="text-[#6B7280] dark:text-[#9090a8] mt-3">No posts yet — run grit seed</Text>
            </View>
          )
        }
        ListFooterComponent={query.isFetchingNextPage ? <ActivityIndicator color="#6c5ce7" style={{ marginVertical: 16 }} /> : null}
      />
    </View>
  );
}
`
}

func expoBlogCreateScreen() string {
	return `import { useState } from "react";
import { View, Text, TextInput, ScrollView, Pressable, Switch, ActivityIndicator, Alert } from "react-native";
import { Image } from "expo-image";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { pickAndUploadImage } from "@/lib/upload";
import { useCreateBlog } from "@/hooks/use-blogs";

export default function CreateBlogScreen() {
  const router = useRouter();
  const create = useCreateBlog();
  const [title, setTitle] = useState("");
  const [excerpt, setExcerpt] = useState("");
  const [content, setContent] = useState("");
  const [image, setImage] = useState<string | null>(null);
  const [published, setPublished] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const inputClass =
    "bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-3.5 text-[#0F1018] dark:text-white text-[15px] mb-4";
  const labelClass = "text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2";

  const onPickImage = async () => {
    try {
      const url = await pickAndUploadImage();
      if (url) setImage(url);
    } catch (e: any) {
      Alert.alert("Upload failed", e.message || "Please try again");
    }
  };

  const onSubmit = async () => {
    setError("");
    if (!title.trim()) return setError("Title is required");
    setSaving(true);
    try {
      await create.mutateAsync({ title, excerpt, content, image: image || "", published });
      router.back();
    } catch (e: any) {
      setError(e.message || "Failed to create post");
    } finally {
      setSaving(false);
    }
  };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="New Post" showBack />
      <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 48 }} keyboardShouldPersistTaps="handled">
        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center">
            <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
            <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{error}</Text>
          </View>
        ) : null}

        <Text className={labelClass}>Cover image</Text>
        <Pressable onPress={onPickImage} className="mb-4 h-40 rounded-2xl border border-dashed border-[#E5E7EB] dark:border-[#2a2a3a] items-center justify-center overflow-hidden bg-white dark:bg-[#111118]">
          {image ? (
            <Image source={{ uri: image }} style={{ width: "100%", height: "100%" }} contentFit="cover" />
          ) : (
            <View className="items-center">
              <Ionicons name="cloud-upload-outline" size={28} color="#9CA3AF" />
              <Text className="text-[#6B7280] dark:text-[#9090a8] mt-2 text-[13px]">Tap to upload</Text>
            </View>
          )}
        </Pressable>

        <Text className={labelClass}>Title</Text>
        <TextInput className={inputClass} placeholder="Title" placeholderTextColor="#9CA3AF" value={title} onChangeText={setTitle} />
        <Text className={labelClass}>Excerpt</Text>
        <TextInput className={inputClass} placeholder="Short summary" placeholderTextColor="#9CA3AF" value={excerpt} onChangeText={setExcerpt} />
        <Text className={labelClass}>Content</Text>
        <TextInput className={inputClass} placeholder="Write your post..." placeholderTextColor="#9CA3AF" value={content} onChangeText={setContent} multiline numberOfLines={6} style={{ minHeight: 140, textAlignVertical: "top" }} />

        <View className="flex-row items-center justify-between mb-4">
          <Text className={labelClass} style={{ marginBottom: 0 }}>Published</Text>
          <Switch value={published} onValueChange={setPublished} trackColor={{ false: "#D1D5DB", true: "#6c5ce7" }} thumbColor="#ffffff" />
        </View>

        <Pressable onPress={onSubmit} disabled={saving} className="bg-[#6c5ce7] rounded-full py-4 items-center mt-2" style={{ opacity: saving ? 0.7 : 1 }}>
          {saving ? <ActivityIndicator color="#fff" /> : <Text className="text-white font-semibold text-[15px]">Publish post</Text>}
        </Pressable>
      </ScrollView>
    </View>
  );
}
`
}

func expoUserCreateScreen() string {
	return `import { useState } from "react";
import { View, Text, TextInput, ScrollView, Pressable, Switch, ActivityIndicator } from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { api } from "@/lib/api";

const ROLES = ["ADMIN", "EDITOR", "USER"];

export default function CreateUserScreen() {
  const router = useRouter();
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("USER");
  const [active, setActive] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const inputClass =
    "bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-3.5 text-[#0F1018] dark:text-white text-[15px] mb-4";
  const labelClass = "text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2";

  const onSubmit = async () => {
    setError("");
    if (!firstName.trim() || !lastName.trim()) return setError("Name is required");
    if (!email.includes("@")) return setError("Enter a valid email");
    if (password.length < 6) return setError("Password must be at least 6 characters");
    setSaving(true);
    try {
      await api.post("/admin/users", {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        role,
        active,
      });
      router.back();
    } catch (e: any) {
      setError(e.message || "Failed to create user");
    } finally {
      setSaving(false);
    }
  };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="New User" showBack />
      <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 48 }} keyboardShouldPersistTaps="handled">
        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center">
            <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
            <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{error}</Text>
          </View>
        ) : null}

        <Text className={labelClass}>First name</Text>
        <TextInput className={inputClass} placeholder="First name" placeholderTextColor="#9CA3AF" value={firstName} onChangeText={setFirstName} />
        <Text className={labelClass}>Last name</Text>
        <TextInput className={inputClass} placeholder="Last name" placeholderTextColor="#9CA3AF" value={lastName} onChangeText={setLastName} />
        <Text className={labelClass}>Email</Text>
        <TextInput className={inputClass} placeholder="you@example.com" placeholderTextColor="#9CA3AF" value={email} onChangeText={setEmail} keyboardType="email-address" autoCapitalize="none" />
        <Text className={labelClass}>Password</Text>
        <TextInput className={inputClass} placeholder="Min. 6 characters" placeholderTextColor="#9CA3AF" value={password} onChangeText={setPassword} secureTextEntry />

        <Text className={labelClass}>Role</Text>
        <View className="flex-row mb-4">
          {ROLES.map((r) => (
            <Pressable key={r} onPress={() => setRole(r)} className={role === r ? "px-4 py-2 mr-2 rounded-full bg-[#6c5ce7]" : "px-4 py-2 mr-2 rounded-full bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a]"}>
              <Text className={role === r ? "text-white font-medium capitalize" : "text-[#0F1018] dark:text-white capitalize"}>{r.toLowerCase()}</Text>
            </Pressable>
          ))}
        </View>

        <View className="flex-row items-center justify-between mb-4">
          <Text className={labelClass} style={{ marginBottom: 0 }}>Active</Text>
          <Switch value={active} onValueChange={setActive} trackColor={{ false: "#D1D5DB", true: "#6c5ce7" }} thumbColor="#ffffff" />
        </View>

        <Pressable onPress={onSubmit} disabled={saving} className="bg-[#6c5ce7] rounded-full py-4 items-center mt-2" style={{ opacity: saving ? 0.7 : 1 }}>
          {saving ? <ActivityIndicator color="#fff" /> : <Text className="text-white font-semibold text-[15px]">Create user</Text>}
        </Pressable>
      </ScrollView>
    </View>
  );
}
`
}
