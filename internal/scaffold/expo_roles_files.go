package scaffold

// Roles & permissions for the mobile app.
//
// The admin panel has had a permission editor since v3.69.0, but mobile had
// no notion of roles at all — it did not even call /api/auth/permissions. That
// left two visible failures: the home screen fetched an ADMIN-only route for
// every user and rendered the resulting 403 as "0", and the create-user screen
// offered a hardcoded ["ADMIN","EDITOR","USER"] list, so a role you defined in
// the admin never appeared on the phone.
//
// These screens mirror the admin's model exactly — same endpoints, same
// wildcard semantics — rather than inventing a second, simpler one that would
// drift. Selection is seeded from the server's `expanded` list and sent back as
// explicit keys; the client never resolves wildcards itself.
//
// As everywhere else, this is a UX layer. Every route is enforced server-side.

// expoRolesListScreen emits app/roles/index.tsx.
func expoRolesListScreen() string {
	return `import { View, Text, ScrollView, Pressable, ActivityIndicator, RefreshControl } from "react-native";
import { useRouter } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { usePermissions } from "@/hooks/use-permissions";
import { api } from "@/lib/api";

interface Role {
  id: string;
  name: string;
  description: string;
  grants: string[];
  expanded: string[];
  is_system: boolean;
  user_count: number;
}

export default function RolesScreen() {
  const router = useRouter();
  const { can, isLoading: permsLoading } = usePermissions();

  const { data: roles, isLoading, refetch, isRefetching } = useQuery<Role[]>({
    queryKey: ["roles"],
    enabled: can("roles.view"),
    queryFn: async () => {
      const res = await api.get("/roles");
      return (res.data ?? []) as Role[];
    },
  });

  if (permsLoading || isLoading) {
    return (
      <View className="flex-1 bg-[#F7F7F9] dark:bg-[#0a0a0f]">
        <ScreenHeader title="Roles & permissions" showBack />
        <View className="flex-1 items-center justify-center">
          <ActivityIndicator color="#6c5ce7" />
        </View>
      </View>
    );
  }

  // Hidden rather than empty: a user without roles.view has no business
  // seeing the shape of the permission system.
  if (!can("roles.view")) {
    return (
      <View className="flex-1 bg-[#F7F7F9] dark:bg-[#0a0a0f]">
        <ScreenHeader title="Roles & permissions" showBack />
        <View className="flex-1 items-center justify-center px-8">
          <Ionicons name="lock-closed-outline" size={40} color="#9090a8" />
          <Text className="text-[#6B7280] dark:text-[#9090a8] text-[15px] text-center mt-4">
            You do not have permission to view roles.
          </Text>
        </View>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-[#F7F7F9] dark:bg-[#0a0a0f]">
      <ScreenHeader title="Roles & permissions" showBack />
      <ScrollView
        className="flex-1 px-5"
        refreshControl={<RefreshControl refreshing={isRefetching} onRefresh={refetch} tintColor="#6c5ce7" />}
      >
        <Text className="text-[13.5px] text-[#6B7280] dark:text-[#9090a8] mb-4">
          Define what each role can do. A role granted a whole resource keeps any actions added to it later.
        </Text>

        {can("roles.create") && (
          <Pressable
            onPress={() => router.push("/roles/new")}
            className="flex-row items-center justify-center rounded-2xl bg-[#6c5ce7] py-3.5 mb-5"
          >
            <Ionicons name="add" size={18} color="#FFFFFF" />
            <Text className="text-white font-semibold text-[15px] ml-2">New role</Text>
          </Pressable>
        )}

        {roles?.map((role) => (
          <Pressable
            key={role.id}
            onPress={() => router.push("/roles/" + role.id)}
            className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl p-4 mb-3"
          >
            <View className="flex-row items-center mb-1">
              <Text className="text-[16px] font-semibold text-[#0F1018] dark:text-white">{role.name}</Text>
              {role.is_system && (
                <View className="ml-2 px-2 py-0.5 rounded-full bg-[#F4F4F6] dark:bg-[#1a1a24]">
                  <Text className="text-[11px] text-[#6B7280] dark:text-[#9090a8]">Built-in</Text>
                </View>
              )}
            </View>
            <Text className="text-[13.5px] text-[#6B7280] dark:text-[#9090a8] mb-3">{role.description}</Text>
            <View className="flex-row items-center">
              <Text className="text-[12.5px] text-[#6B7280] dark:text-[#9090a8]">
                {role.user_count} {role.user_count === 1 ? "user" : "users"}
              </Text>
              <Text className="text-[12.5px] text-[#6c5ce7] ml-3 font-medium">
                {role.grants?.includes("*") ? "all permissions" : (role.expanded?.length ?? 0) + " granted"}
              </Text>
            </View>
          </Pressable>
        ))}

        <View className="h-10" />
      </ScrollView>
    </View>
  );
}
`
}

// expoRoleEditorScreen emits app/roles/[id].tsx — create and edit in one
// screen. "new" is treated as the create case so there is a single permission
// tree to keep correct instead of two.
func expoRoleEditorScreen() string {
	return `import { useState, useEffect } from "react";
import { View, Text, TextInput, ScrollView, Pressable, ActivityIndicator, Alert } from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";
import { usePermissions } from "@/hooks/use-permissions";
import { api } from "@/lib/api";

interface Feature {
  key: string;
  name: string;
  actions: string[];
}
interface Group {
  key: string;
  name: string;
  features: Feature[];
}
interface Module {
  key: string;
  name: string;
  groups: Group[];
}
interface Catalog {
  keys: string[];
  modules: Module[];
}
interface Role {
  id: string;
  name: string;
  description: string;
  grants: string[];
  expanded: string[];
  is_system: boolean;
}

export default function RoleEditorScreen() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { id } = useLocalSearchParams<{ id: string }>();
  const isNew = id === "new";
  const { can } = usePermissions();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [selected, setSelected] = useState<Set<string>>(new Set());
  // Permission modules collapse by default; expand one at a time.
  const [openModules, setOpenModules] = useState<Set<string>>(new Set());
  const toggleModule = (key: string) =>
    setOpenModules((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const { data: catalog } = useQuery<Catalog>({
    queryKey: ["permission-catalog"],
    queryFn: async () => {
      const res = await api.get("/permissions");
      return res.data as Catalog;
    },
  });

  const { data: role, isLoading } = useQuery<Role>({
    queryKey: ["role", id],
    enabled: !isNew,
    queryFn: async () => {
      const res = await api.get("/roles/" + id);
      return res.data as Role;
    },
  });

  // Seeded from the server-expanded list so wildcards resolve exactly as the
  // API resolves them.
  useEffect(() => {
    if (role) {
      setName(role.name);
      setDescription(role.description);
      setSelected(new Set(role.expanded ?? []));
    }
  }, [role]);

  const locked = !isNew && role?.is_system;

  function toggle(key: string) {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  }

  function toggleFeature(feature: Feature) {
    const keys = feature.actions.map((a) => feature.key + "." + a);
    const allOn = keys.every((k) => selected.has(k));
    setSelected((prev) => {
      const next = new Set(prev);
      for (const k of keys) {
        if (allOn) next.delete(k);
        else next.add(k);
      }
      return next;
    });
  }

  const onSave = async () => {
    setError("");
    if (!name.trim()) return setError("Name is required");
    setSaving(true);
    try {
      const body = {
        name: name.trim(),
        description: description.trim(),
        grants: Array.from(selected),
      };
      if (isNew) await api.post("/roles", body);
      else await api.put("/roles/" + id, body);
      await queryClient.invalidateQueries({ queryKey: ["roles"] });
      router.back();
    } catch (e: any) {
      setError(e.message || "Failed to save role");
    } finally {
      setSaving(false);
    }
  };

  const onDelete = () => {
    Alert.alert("Delete role", "Users with only this role lose its permissions.", [
      { text: "Cancel", style: "cancel" },
      {
        text: "Delete",
        style: "destructive",
        onPress: async () => {
          try {
            await api.delete("/roles/" + id);
            await queryClient.invalidateQueries({ queryKey: ["roles"] });
            router.back();
          } catch (e: any) {
            setError(e.message || "Failed to delete role");
          }
        },
      },
    ]);
  };

  if (!isNew && isLoading) {
    return (
      <View className="flex-1 bg-[#F7F7F9] dark:bg-[#0a0a0f]">
        <ScreenHeader title="Role" showBack />
        <View className="flex-1 items-center justify-center">
          <ActivityIndicator color="#6c5ce7" />
        </View>
      </View>
    );
  }

  const inputClass =
    "bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-3.5 text-[#0F1018] dark:text-white text-[15px] mb-4";
  const labelClass = "text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2";

  return (
    <View className="flex-1 bg-[#F7F7F9] dark:bg-[#0a0a0f]">
      <ScreenHeader title={isNew ? "New role" : name || "Role"} showBack />
      <ScrollView className="flex-1 px-5">
        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-2xl px-4 py-3 mb-4">
            <Text className="text-[#ff6b6b] text-[13.5px]">{error}</Text>
          </View>
        ) : null}

        <Text className={labelClass}>Name</Text>
        <TextInput
          className={inputClass}
          placeholder="EDITOR"
          placeholderTextColor="#9CA3AF"
          value={name}
          onChangeText={setName}
          autoCapitalize="characters"
          editable={!locked}
        />

        <Text className={labelClass}>Description</Text>
        <TextInput
          className={inputClass}
          placeholder="What this role is for"
          placeholderTextColor="#9CA3AF"
          value={description}
          onChangeText={setDescription}
        />

        {locked ? (
          <Text className="text-[12.5px] text-[#6B7280] dark:text-[#9090a8] mb-4">
            Built-in role — the name is fixed, but its permissions can still be changed.
          </Text>
        ) : null}

        <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white mb-1 mt-2">Permissions</Text>
        <Text className="text-[12.5px] text-[#6B7280] dark:text-[#9090a8] mb-4">
          {selected.size} of {catalog?.keys?.length ?? 0} granted
        </Text>

        {catalog?.modules?.map((mod) => {
          const modKeys = mod.groups.flatMap((g) =>
            g.features.flatMap((f) => f.actions.map((a) => f.key + "." + a))
          );
          const modOn = modKeys.filter((k) => selected.has(k)).length;
          const isOpen = openModules.has(mod.key);
          return (
          <View key={mod.key} className="mb-5">
            {/* Collapsed by default — tap a module to expand it. */}
            <Pressable onPress={() => toggleModule(mod.key)} className="flex-row items-center justify-between mb-2">
              <View className="flex-row items-center">
                <Ionicons name={isOpen ? "chevron-down" : "chevron-forward"} size={14} color="#9090a8" />
                <Text className="text-[12px] font-semibold uppercase tracking-wide text-[#9090a8] ml-1">{mod.name}</Text>
              </View>
              <Text className="text-[12px] text-[#6c5ce7]">{modOn}/{modKeys.length}</Text>
            </Pressable>
            {isOpen && mod.groups.map((group) => (
              <View key={group.key} className="mb-3">
                {group.features.map((feature) => {
                  const keys = feature.actions.map((a) => feature.key + "." + a);
                  const onCount = keys.filter((k) => selected.has(k)).length;
                  return (
                    <View
                      key={feature.key}
                      className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl p-4 mb-2"
                    >
                      <Pressable onPress={() => toggleFeature(feature)} className="flex-row items-center justify-between mb-3">
                        <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white">{feature.name}</Text>
                        <Text className="text-[12.5px] text-[#6c5ce7]">
                          {onCount}/{keys.length}
                        </Text>
                      </Pressable>
                      <View className="flex-row flex-wrap">
                        {feature.actions.map((action) => {
                          const key = feature.key + "." + action;
                          const on = selected.has(key);
                          return (
                            <Pressable
                              key={key}
                              onPress={() => toggle(key)}
                              className={
                                on
                                  ? "flex-row items-center px-3 py-2 mr-2 mb-2 rounded-full bg-[#6c5ce7]"
                                  : "flex-row items-center px-3 py-2 mr-2 mb-2 rounded-full bg-[#F4F4F6] dark:bg-[#1a1a24]"
                              }
                            >
                              {on && <Ionicons name="checkmark" size={13} color="#FFFFFF" style={{ marginRight: 4 }} />}
                              <Text className={on ? "text-white text-[13px] capitalize" : "text-[#6B7280] dark:text-[#9090a8] text-[13px] capitalize"}>
                                {action}
                              </Text>
                            </Pressable>
                          );
                        })}
                      </View>
                    </View>
                  );
                })}
              </View>
            ))}
          </View>
          );
        })}

        <Pressable
          onPress={onSave}
          disabled={saving}
          className="rounded-2xl bg-[#6c5ce7] py-4 items-center mb-3"
        >
          {saving ? <ActivityIndicator color="#fff" /> : <Text className="text-white font-semibold text-[15px]">Save role</Text>}
        </Pressable>

        {!isNew && !role?.is_system && can("roles.delete") && (
          <Pressable onPress={onDelete} className="rounded-2xl border border-[#ff6b6b]/40 py-4 items-center mb-3">
            <Text className="text-[#ff6b6b] font-semibold text-[15px]">Delete role</Text>
          </Pressable>
        )}

        <View className="h-10" />
      </ScrollView>
    </View>
  );
}
`
}
