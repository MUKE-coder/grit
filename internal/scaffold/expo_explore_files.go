package scaffold

// Explore destination screens + profile helpers for the scaffolded Expo app.
// Each screen uses the shared ScreenHeader (safe-area + back button). Users,
// Notifications and Storage are wired to real API endpoints; Analytics reads
// live counts; Content and Integrations are polished starting points.

// expoUploadHelper picks an image from the library and uploads it to the
// POST /uploads endpoint, returning the stored file URL (used for avatars).
func expoUploadHelper() string {
	return `import * as ImagePicker from "expo-image-picker";
import * as SecureStore from "expo-secure-store";
import * as FileSystem from "expo-file-system/legacy";
import { API_URL } from "./api";

// Pick an image and upload it to POST /uploads (multipart). Returns the
// public URL of the stored file, or null if the user cancelled.
export async function pickAndUploadImage(): Promise<string | null> {
  const perm = await ImagePicker.requestMediaLibraryPermissionsAsync();
  if (!perm.granted) throw new Error("Photo library permission is required");

  const result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ["images"],
    allowsEditing: true,
    aspect: [1, 1],
    quality: 0.8,
  });
  if (result.canceled || !result.assets?.length) return null;

  const asset = result.assets[0];

  // Derive a MIME type from the file extension when the picker omits it, so
  // the multipart part is tagged with a type the server's allowlist accepts.
  const name = asset.fileName || asset.uri.split("/").pop() || "photo.jpg";
  const ext = (name.split(".").pop() || "jpg").toLowerCase();
  const extToMime: Record<string, string> = {
    jpg: "image/jpeg",
    jpeg: "image/jpeg",
    png: "image/png",
    gif: "image/gif",
    webp: "image/webp",
    heic: "image/heic",
  };
  const mimeType = asset.mimeType || extToMime[ext] || "image/jpeg";

  const token = await SecureStore.getItemAsync("access_token");

  // Upload the local file natively as multipart/form-data. This is far more
  // reliable than fetch + FormData on React Native (RN 0.81 / Expo SDK 54),
  // where the file part can be dropped entirely — the server then reports
  // "No file provided". expo-file-system streams the file itself.
  const res = await FileSystem.uploadAsync(API_URL + "/uploads", asset.uri, {
    httpMethod: "POST",
    uploadType: FileSystem.FileSystemUploadType.MULTIPART,
    fieldName: "file",
    mimeType,
    headers: token ? { Authorization: "Bearer " + token } : {},
  });

  const json = res.body ? JSON.parse(res.body) : null;
  if (res.status < 200 || res.status >= 300) {
    // Surface the server's reason (e.g. "File type not allowed") instead of a
    // generic failure, so upload problems are diagnosable from the UI.
    throw new Error(json?.error?.message || "Upload failed (" + res.status + ")");
  }
  return json?.data?.url ?? null;
}
`
}

func expoExploreUsers() string {
	return `import { View, Text, FlatList, Pressable, ActivityIndicator, RefreshControl } from "react-native";
import { useInfiniteQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { api } from "@/lib/api";

interface UserRow {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  avatar?: string;
}

const ROLE_TINT: Record<string, string> = {
  ADMIN: "#6c5ce7",
  EDITOR: "#00b894",
  USER: "#74b9ff",
};

export default function UsersScreen() {
  const router = useRouter();
  // GET /users is an admin route — the seeded admin account can browse it.
  const query = useInfiniteQuery({
    queryKey: ["explore-users"],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => api.get("/users?page=" + pageParam + "&page_size=20"),
    getNextPageParam: (last: any) =>
      last?.meta && last.meta.page < last.meta.pages ? last.meta.page + 1 : undefined,
  });
  const users: UserRow[] = query.data?.pages.flatMap((p: any) => p.data) ?? [];

  const renderItem = ({ item }: { item: UserRow }) => {
    const initials = ((item.first_name?.[0] || "") + (item.last_name?.[0] || "")).toUpperCase() || "?";
    const tint = ROLE_TINT[item.role] || "#6c5ce7";
    return (
      <View className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 mb-3">
        <View className="w-11 h-11 rounded-full bg-[#6c5ce7]/12 items-center justify-center mr-3">
          <Text className="text-[#6c5ce7] font-bold">{initials}</Text>
        </View>
        <View className="flex-1">
          <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white" numberOfLines={1}>
            {item.first_name} {item.last_name}
          </Text>
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8]" numberOfLines={1}>{item.email}</Text>
        </View>
        <View style={{ backgroundColor: tint + "22" }} className="px-2.5 py-1 rounded-full">
          <Text style={{ color: tint }} className="text-[11px] font-semibold capitalize">
            {(item.role || "user").toLowerCase()}
          </Text>
        </View>
      </View>
    );
  };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader
        title="Users"
        subtitle="All user accounts"
        showBack
        right={
          <Pressable onPress={() => router.push("/users/new")} hitSlop={8}>
            <Ionicons name="add-circle" size={28} color="#6c5ce7" />
          </Pressable>
        }
      />
      <FlatList
        data={users}
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
            <Text className="text-center text-[#6B7280] dark:text-[#9090a8] mt-16">No users found</Text>
          )
        }
        ListFooterComponent={
          query.isFetchingNextPage ? <ActivityIndicator color="#6c5ce7" style={{ marginVertical: 16 }} /> : null
        }
      />
    </View>
  );
}
`
}

func expoExploreNotifications() string {
	return `import { View, Text, FlatList, ActivityIndicator, RefreshControl } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { api } from "@/lib/api";

interface Notification {
  id: string;
  title?: string;
  message?: string;
  body?: string;
  read?: boolean;
  read_at?: string | null;
  created_at?: string;
}

export default function NotificationsScreen() {
  const query = useQuery({
    queryKey: ["explore-notifications"],
    queryFn: async () => api.get("/notifications"),
  });
  const items: Notification[] = (query.data as any)?.data ?? [];

  const renderItem = ({ item }: { item: Notification }) => {
    const unread = !item.read && !item.read_at;
    return (
      <View className="flex-row items-start bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 mb-3">
        <View className="w-10 h-10 rounded-full bg-[#6c5ce7]/12 items-center justify-center mr-3">
          <Ionicons name="notifications" size={18} color="#6c5ce7" />
        </View>
        <View className="flex-1">
          <Text className="text-[14px] font-semibold text-[#0F1018] dark:text-white">
            {item.title || "Notification"}
          </Text>
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">
            {item.message || item.body || ""}
          </Text>
        </View>
        {unread ? <View className="w-2.5 h-2.5 rounded-full bg-[#6c5ce7] mt-1" /> : null}
      </View>
    );
  };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Notifications" subtitle="Alerts and messages" showBack />
      <FlatList
        data={items}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        contentContainerStyle={{ paddingHorizontal: 24, paddingBottom: 40 }}
        refreshControl={<RefreshControl refreshing={query.isRefetching} onRefresh={query.refetch} tintColor="#6c5ce7" />}
        ListEmptyComponent={
          query.isLoading ? (
            <ActivityIndicator color="#6c5ce7" style={{ marginTop: 40 }} />
          ) : (
            <View className="items-center mt-16">
              <Ionicons name="notifications-off-outline" size={40} color="#9CA3AF" />
              <Text className="text-[#6B7280] dark:text-[#9090a8] mt-3">You're all caught up</Text>
            </View>
          )
        }
      />
    </View>
  );
}
`
}

func expoExploreStorage() string {
	return `import { View, Text, FlatList, ActivityIndicator, RefreshControl } from "react-native";
import { Image } from "expo-image";
import { useInfiniteQuery } from "@tanstack/react-query";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { api } from "@/lib/api";

interface Upload {
  id: string;
  original_name: string;
  mime_type: string;
  size: number;
  url: string;
  thumbnail_url?: string;
}

function humanSize(bytes: number): string {
  if (!bytes) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return (bytes / Math.pow(1024, i)).toFixed(i ? 1 : 0) + " " + units[i];
}

export default function StorageScreen() {
  const query = useInfiniteQuery({
    queryKey: ["explore-uploads"],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => api.get("/uploads?page=" + pageParam + "&page_size=20"),
    getNextPageParam: (last: any) =>
      last?.meta && last.meta.page < last.meta.pages ? last.meta.page + 1 : undefined,
  });
  const files: Upload[] = query.data?.pages.flatMap((p: any) => p.data) ?? [];
  const total = (query.data?.pages[0] as any)?.meta?.total ?? files.length;

  const renderItem = ({ item }: { item: Upload }) => {
    const isImage = item.mime_type?.startsWith("image/");
    const thumb = item.thumbnail_url || item.url;
    return (
      <View className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-3 mb-3">
        {isImage && thumb ? (
          <Image source={{ uri: thumb }} style={{ width: 44, height: 44, borderRadius: 10, marginRight: 12 }} contentFit="cover" />
        ) : (
          <View className="w-11 h-11 rounded-[10px] bg-[#6c5ce7]/12 items-center justify-center mr-3">
            <Ionicons name="document-outline" size={20} color="#6c5ce7" />
          </View>
        )}
        <View className="flex-1">
          <Text className="text-[14px] font-medium text-[#0F1018] dark:text-white" numberOfLines={1}>
            {item.original_name}
          </Text>
          <Text className="text-[12px] text-[#6B7280] dark:text-[#9090a8]">{humanSize(item.size)}</Text>
        </View>
      </View>
    );
  };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Storage" subtitle={total + " file" + (total === 1 ? "" : "s")} showBack />
      <FlatList
        data={files}
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
              <Ionicons name="cloud-outline" size={40} color="#9CA3AF" />
              <Text className="text-[#6B7280] dark:text-[#9090a8] mt-3">No files uploaded yet</Text>
            </View>
          )
        }
      />
    </View>
  );
}
`
}

func expoExploreAnalytics() string {
	return `import { View, Text, ScrollView, RefreshControl } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { api } from "@/lib/api";

function Metric({ icon, label, value, tint }: { icon: string; label: string; value: string; tint: string }) {
  return (
    <View className="flex-1 bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 m-1.5">
      <View style={{ backgroundColor: tint + "20" }} className="w-9 h-9 rounded-xl items-center justify-center mb-3">
        <Ionicons name={icon as any} size={18} color={tint} />
      </View>
      <Text className="text-[24px] font-bold text-[#0F1018] dark:text-white">{value}</Text>
      <Text className="text-[12px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">{label}</Text>
    </View>
  );
}

export default function AnalyticsScreen() {
  const usersQ = useQuery({ queryKey: ["an-users"], queryFn: async () => api.get("/users?page=1&page_size=1") });
  const filesQ = useQuery({ queryKey: ["an-uploads"], queryFn: async () => api.get("/uploads?page=1&page_size=1") });

  const userCount = (usersQ.data as any)?.meta?.total ?? 0;
  const fileCount = (filesQ.data as any)?.meta?.total ?? 0;
  const refreshing = usersQ.isRefetching || filesQ.isRefetching;
  const onRefresh = () => { usersQ.refetch(); filesQ.refetch(); };

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Analytics" subtitle="Usage at a glance" showBack />
      <ScrollView
        contentContainerStyle={{ padding: 22, paddingBottom: 40 }}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor="#6c5ce7" />}
      >
        <View className="flex-row">
          <Metric icon="people-outline" label="Total users" value={String(userCount)} tint="#6c5ce7" />
          <Metric icon="cloud-outline" label="Files stored" value={String(fileCount)} tint="#00b894" />
        </View>
        <View className="flex-row">
          <Metric icon="pulse-outline" label="Active today" value={String(userCount)} tint="#74b9ff" />
          <Metric icon="trending-up-outline" label="Growth" value="+0%" tint="#fdcb6e" />
        </View>

        <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-5 mt-3">
          <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white mb-1">Wire up your metrics</Text>
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] leading-5">
            These cards read live counts from your API. Add resource-specific endpoints
            and drop more Metric cards here as your product grows.
          </Text>
        </View>
      </ScrollView>
    </View>
  );
}
`
}

func expoExploreContent() string {
	return `import { View, Text, ScrollView } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";

const SECTIONS = [
  { icon: "document-text-outline", title: "Pages", description: "Static marketing and info pages", tint: "#6c5ce7" },
  { icon: "newspaper-outline", title: "Posts", description: "Blog posts and announcements", tint: "#00b894" },
  { icon: "images-outline", title: "Media", description: "Images, video and documents", tint: "#74b9ff" },
];

export default function ContentScreen() {
  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Content" subtitle="Posts, pages and media" showBack />
      <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 40 }}>
        {SECTIONS.map((s) => (
          <View
            key={s.title}
            className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-5 mb-3"
          >
            <View style={{ backgroundColor: s.tint + "20" }} className="w-11 h-11 rounded-xl items-center justify-center mr-4">
              <Ionicons name={s.icon as any} size={22} color={s.tint} />
            </View>
            <View className="flex-1">
              <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white">{s.title}</Text>
              <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">{s.description}</Text>
            </View>
          </View>
        ))}
        <View className="items-center mt-8">
          <Text className="text-[13px] text-[#9CA3AF] dark:text-[#606078] text-center px-6">
            Generate a content resource with{"\n"}
            <Text className="font-semibold text-[#6c5ce7]">grit generate resource Post</Text>
            {"\n"}to power this screen with real data.
          </Text>
        </View>
      </ScrollView>
    </View>
  );
}
`
}

func expoExploreIntegrations() string {
	return `import { useState } from "react";
import { View, Text, ScrollView, Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";

interface Integration {
  id: string;
  name: string;
  description: string;
  icon: string;
  tint: string;
}

const INTEGRATIONS: Integration[] = [
  { id: "stripe", name: "Stripe", description: "Payments and billing", icon: "card-outline", tint: "#6c5ce7" },
  { id: "resend", name: "Resend", description: "Transactional email", icon: "mail-outline", tint: "#00b894" },
  { id: "slack", name: "Slack", description: "Team notifications", icon: "chatbubbles-outline", tint: "#74b9ff" },
  { id: "github", name: "GitHub", description: "Sync issues and PRs", icon: "logo-github", tint: "#9090a8" },
  { id: "openai", name: "OpenAI", description: "AI features", icon: "sparkles-outline", tint: "#fdcb6e" },
];

export default function IntegrationsScreen() {
  const [connected, setConnected] = useState<Record<string, boolean>>({});

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Integrations" subtitle="Connected services" showBack />
      <ScrollView contentContainerStyle={{ padding: 24, paddingBottom: 40 }}>
        {INTEGRATIONS.map((it) => {
          const on = connected[it.id];
          return (
            <View
              key={it.id}
              className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 mb-3"
            >
              <View style={{ backgroundColor: it.tint + "20" }} className="w-11 h-11 rounded-xl items-center justify-center mr-4">
                <Ionicons name={it.icon as any} size={22} color={it.tint} />
              </View>
              <View className="flex-1">
                <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white">{it.name}</Text>
                <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">{it.description}</Text>
              </View>
              <Pressable
                onPress={() => setConnected((c) => ({ ...c, [it.id]: !c[it.id] }))}
                className={on ? "px-4 py-2 rounded-full bg-[#6c5ce7]/12" : "px-4 py-2 rounded-full bg-[#6c5ce7]"}
              >
                <Text className={on ? "text-[13px] font-semibold text-[#6c5ce7]" : "text-[13px] font-semibold text-white"}>
                  {on ? "Connected" : "Connect"}
                </Text>
              </Pressable>
            </View>
          );
        })}
      </ScrollView>
    </View>
  );
}
`
}

func expoChangePasswordScreen() string {
	return `import { useState } from "react";
import { View, Text, TextInput, ActivityIndicator, Pressable, Alert } from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import * as Haptics from "expo-haptics";
import { ScreenHeader } from "@/components/ui/screen-header";
import { api } from "@/lib/api";

export default function ChangePasswordScreen() {
  const router = useRouter();
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [show, setShow] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const onSubmit = async () => {
    setError("");
    if (password.length < 8) return setError("Password must be at least 8 characters");
    if (password !== confirm) return setError("Passwords do not match");

    setLoading(true);
    try {
      await api.put("/profile", { password });
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success).catch(() => {});
      Alert.alert("Password changed", "Your password has been updated.", [
        { text: "OK", onPress: () => router.back() },
      ]);
    } catch (err: any) {
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Error).catch(() => {});
      setError(err.message || "Failed to change password");
    } finally {
      setLoading(false);
    }
  };

  const field = "bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl flex-row items-center px-4";

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Change Password" subtitle="Choose a new password" showBack />
      <View className="px-6 pt-2">
        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center">
            <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
            <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{error}</Text>
          </View>
        ) : null}

        <Text className="text-[12.5px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2">New password</Text>
        <View className={field} style={{ height: 52 }}>
          <Ionicons name="lock-closed-outline" size={17} color="#9CA3AF" />
          <TextInput
            className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
            placeholder="Min. 8 characters"
            placeholderTextColor="#9CA3AF"
            value={password}
            onChangeText={setPassword}
            secureTextEntry={!show}
            autoCapitalize="none"
          />
          <Pressable onPress={() => setShow((s) => !s)} hitSlop={10} className="p-1">
            <Ionicons name={show ? "eye-off-outline" : "eye-outline"} size={19} color={show ? "#6c5ce7" : "#9CA3AF"} />
          </Pressable>
        </View>

        <Text className="text-[12.5px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2 mt-4">Confirm password</Text>
        <View className={field} style={{ height: 52 }}>
          <Ionicons name="lock-closed-outline" size={17} color="#9CA3AF" />
          <TextInput
            className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
            placeholder="Repeat password"
            placeholderTextColor="#9CA3AF"
            value={confirm}
            onChangeText={setConfirm}
            secureTextEntry={!show}
            autoCapitalize="none"
          />
        </View>

        <Pressable
          onPress={onSubmit}
          disabled={loading}
          className="bg-[#6c5ce7] rounded-full py-4 items-center mt-6"
          style={{ opacity: loading ? 0.7 : 1 }}
        >
          {loading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text className="text-white font-semibold text-[15px]">Update password</Text>
          )}
        </Pressable>
      </View>
    </View>
  );
}
`
}
