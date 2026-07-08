package scaffold

// Mobile surface for full-database backups. Admins occasionally need to pull a
// backup while away from a laptop — this lists them, lets you take one on demand
// (rate-limited server-side), and opens the pre-signed download URL in the
// browser so the OS handles the file, not React Native.

// ExpoBackupsScreen is app/backups.tsx.
func ExpoBackupsScreen() string {
	return `import { View, Text, FlatList, Pressable, ActivityIndicator, RefreshControl, Alert, Linking } from "react-native";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { useTheme } from "@/lib/theme";
import { api } from "@/lib/api";

interface Backup {
  id: string;
  kind: "WEEKLY" | "MANUAL" | "CLI";
  status: "RUNNING" | "READY" | "FAILED" | "PURGED";
  size_bytes: number;
  table_count: number;
  row_count: number;
  error?: string;
  created_at: string;
}

const STATUS_TINT: Record<Backup["status"], string> = {
  RUNNING: "#74b9ff",
  READY: "#00b894",
  FAILED: "#ff6b6b",
  PURGED: "#9CA3AF",
};

function formatBytes(n: number): string {
  if (!n) return "—";
  if (n < 1024) return n + " B";
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + " KB";
  return (n / (1024 * 1024)).toFixed(1) + " MB";
}

export default function BackupsScreen() {
  const { palette } = useTheme();
  const qc = useQueryClient();

  // Poll every 3s while one is RUNNING, then go idle.
  const query = useQuery<Backup[]>({
    queryKey: ["backups"],
    queryFn: async () => (await api.get("/backups"))?.data ?? [],
    refetchInterval: (q) =>
      ((q.state.data as Backup[] | undefined) ?? []).some((b) => b.status === "RUNNING") ? 3000 : false,
  });

  const backups = query.data ?? [];
  const running = backups.some((b) => b.status === "RUNNING");

  const generate = useMutation({
    mutationFn: async () => api.post("/backups/generate", {}),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["backups"] }),
    onError: (e: any) => Alert.alert("Backup failed", e?.message || "Please try again"),
  });

  // The server mints a short-lived pre-signed URL; hand it to the OS browser so
  // a large archive downloads natively instead of through the JS bridge.
  const download = async (id: string) => {
    try {
      const res: any = await api.get("/backups/" + id + "/download");
      const url = res?.data?.url;
      if (url) await Linking.openURL(url);
    } catch (e: any) {
      Alert.alert("Download failed", e?.message || "Please try again");
    }
  };

  const renderItem = ({ item }: { item: Backup }) => (
    <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-4 mb-3">
      <View className="flex-row items-center justify-between">
        <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white">
          {new Date(item.created_at).toLocaleDateString()}
        </Text>
        <View style={{ backgroundColor: STATUS_TINT[item.status] + "22" }} className="px-2.5 py-1 rounded-full flex-row items-center">
          {item.status === "RUNNING" ? <ActivityIndicator size="small" color={STATUS_TINT[item.status]} /> : null}
          <Text style={{ color: STATUS_TINT[item.status] }} className="text-[11px] font-semibold ml-1">
            {item.status}
          </Text>
        </View>
      </View>

      <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-1">
        {item.kind} · {item.table_count || "—"} tables · {item.row_count ? item.row_count.toLocaleString() : "—"} rows ·{" "}
        {formatBytes(item.size_bytes)}
      </Text>

      {item.status === "FAILED" && item.error ? (
        <Text className="text-[12px] text-[#ff6b6b] mt-1" numberOfLines={2}>{item.error}</Text>
      ) : null}

      {item.status === "READY" ? (
        <Pressable onPress={() => download(item.id)} className="flex-row items-center mt-3">
          <Ionicons name="download-outline" size={16} color="#6c5ce7" />
          <Text className="text-[#6c5ce7] font-semibold text-[13px] ml-1.5">Download</Text>
        </Pressable>
      ) : null}
    </View>
  );

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader
        title="Backups"
        subtitle="Automatic every Sunday 02:00 UTC"
        showBack
        right={
          <Pressable onPress={() => generate.mutate()} disabled={generate.isPending || running} hitSlop={8}>
            {generate.isPending || running ? (
              <ActivityIndicator color="#6c5ce7" />
            ) : (
              <Ionicons name="add-circle" size={28} color="#6c5ce7" />
            )}
          </Pressable>
        }
      />

      <FlatList
        data={backups}
        keyExtractor={(item) => item.id}
        renderItem={renderItem}
        contentContainerStyle={{ paddingHorizontal: 24, paddingBottom: 40 }}
        refreshControl={
          <RefreshControl refreshing={query.isRefetching} onRefresh={query.refetch} tintColor={palette.refresh} />
        }
        ListEmptyComponent={
          query.isLoading ? (
            <ActivityIndicator color={palette.refresh} style={{ marginTop: 40 }} />
          ) : (
            <View className="items-center mt-16 px-8">
              <Ionicons name="server-outline" size={40} color={palette.inputIcon} />
              <Text className="text-[#6B7280] dark:text-[#9090a8] mt-3 text-center">No backups yet</Text>
              <Text className="text-[#9CA3AF] dark:text-[#606078] mt-1 text-[12px] text-center">
                The first one lands on Sunday, or take one now with +
              </Text>
            </View>
          )
        }
        ListFooterComponent={
          backups.length ? (
            <Text className="text-[11px] text-[#9CA3AF] dark:text-[#606078] text-center mt-2 px-4">
              Each archive is a ZIP of CSVs + dump.sql. Restore with: grit restore backup.zip
            </Text>
          ) : null
        }
      />
    </View>
  );
}
`
}
