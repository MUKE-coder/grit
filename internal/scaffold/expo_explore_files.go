package scaffold

// Explore destination screens + profile helpers for the scaffolded Expo app.
// Each screen uses the shared ScreenHeader (safe-area + back button). Users,
// Notifications and Storage are wired to real API endpoints; Analytics reads
// live counts; Content and Integrations are polished starting points.

// expoUploadHelper picks an image from the library and uploads it to the
// POST /uploads endpoint, returning the stored file URL (used for avatars).
// expoExportHelper downloads a resource's CSV export (honouring the current
// search) and opens the native share sheet. Kept in sync with
// internal/generate mobileExportHelperContent.
func expoExportHelper() string {
	return `import * as FileSystem from "expo-file-system/legacy";
import * as Sharing from "expo-sharing";
import * as SecureStore from "expo-secure-store";
import { API_URL } from "./api";

// Download the resource's CSV export and open the share sheet. ` + "`" + `query` + "`" + ` is an
// optional querystring (e.g. the current search) appended to /<plural>/export.
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

// expoImportHelper is lib/import.ts — parse a CSV for preview and upload it to
// the resource's /import endpoint with progress. Kept in sync with
// internal/generate mobileImportHelperContent.
func ExpoImportHelper() string {
	return `import * as FileSystem from "expo-file-system/legacy";
import * as SecureStore from "expo-secure-store";
import * as Sharing from "expo-sharing";
import { API_URL } from "./api";

export interface CsvPreview {
  headers: string[];
  rows: string[][];
  total: number;
}

// Download a ready-to-fill CSV template for the resource and open the share
// sheet so the user can fill it in and re-import.
export async function downloadResourceTemplate(plural: string): Promise<string> {
  const token = await SecureStore.getItemAsync("access_token");
  const url = API_URL + "/" + plural + "/import/template";
  const fileUri = FileSystem.cacheDirectory + plural + "-template.csv";
  const res = await FileSystem.downloadAsync(url, fileUri, {
    headers: token ? { Authorization: "Bearer " + token } : {},
  });
  if (res.status < 200 || res.status >= 300) {
    throw new Error("Template download failed (" + res.status + ")");
  }
  if (await Sharing.isAvailableAsync()) {
    await Sharing.shareAsync(res.uri, { mimeType: "text/csv", dialogTitle: plural + " template" });
  }
  return res.uri;
}

// Read + naively parse a CSV for a preview (first ` + "`" + `limit` + "`" + ` data rows).
export async function parseCsvPreview(fileUri: string, limit = 8): Promise<CsvPreview> {
  const text = await FileSystem.readAsStringAsync(fileUri);
  const lines = text.split(/\r?\n/).filter((l) => l.trim().length > 0);
  const parseLine = (l: string) => l.split(",").map((c) => c.trim().replace(/^"|"$/g, ""));
  const headers = lines.length ? parseLine(lines[0]) : [];
  const dataLines = lines.slice(1);
  return { headers, rows: dataLines.slice(0, limit).map(parseLine), total: dataLines.length };
}

export interface ImportResult {
  created: number;
  skipped: number;
  failed: number;
  errors: { row: number; message: string }[];
}

export interface ImportJobStatus extends ImportResult {
  id: string;
  status: "processing" | "completed" | "failed";
  total: number;
  processed: number;
  message: string;
}

// Upload the CSV to /<plural>/import. The server processes it in the BACKGROUND
// and responds 202 with a job id, so a large file never blocks the request.
// We then poll /imports/:id until it finishes, reporting progress (0..1) so the
// caller can drive a progress bar. Because the polling loop lives here (and is
// kicked off from a module-level store), it keeps running even if the screen
// that started it unmounts.
export async function importResourceCsv(
  plural: string,
  fileUri: string,
  onProgress?: (fraction: number) => void,
): Promise<ImportResult> {
  const token = await SecureStore.getItemAsync("access_token");
  const res = await FileSystem.uploadAsync(API_URL + "/" + plural + "/import", fileUri, {
    httpMethod: "POST",
    uploadType: FileSystem.FileSystemUploadType.MULTIPART,
    fieldName: "file",
    mimeType: "text/csv",
    headers: token ? { Authorization: "Bearer " + token } : {},
  });
  if (!res || res.status < 200 || res.status >= 300) {
    const body = res?.body ? JSON.parse(res.body) : null;
    throw new Error(body?.error?.message || "Import failed (" + (res?.status ?? "?") + ")");
  }
  const started = JSON.parse(res.body).data as { job_id: string; total: number };
  if (onProgress && started.total === 0) onProgress(1);
  return pollImportJob(started.job_id, onProgress);
}

// Poll a background import until it reaches a terminal state, reporting
// processed/total progress on the way. Resolves with the final counts + errors.
export async function pollImportJob(
  jobId: string,
  onProgress?: (fraction: number) => void,
): Promise<ImportResult> {
  const token = await SecureStore.getItemAsync("access_token");
  const url = API_URL + "/imports/" + jobId;
  for (;;) {
    const res = await fetch(url, {
      headers: token ? { Authorization: "Bearer " + token } : {},
    });
    if (!res.ok) throw new Error("Could not check import status (" + res.status + ")");
    const job = (await res.json()).data as ImportJobStatus;
    if (onProgress && job.total > 0) onProgress(job.processed / job.total);
    if (job.status === "completed" || job.status === "failed") {
      return {
        created: job.created,
        skipped: job.skipped,
        failed: job.failed,
        errors: job.errors || [],
      };
    }
    await new Promise((r) => setTimeout(r, 700));
  }
}
`
}

// ExpoImportProgressStore is lib/import-progress.ts — a tiny module-level store
// that owns in-flight background imports. Because it lives outside React, an
// import keeps uploading + polling even after the sheet that started it closes,
// which is what makes imports feel like they "run in the background". The banner
// and the import sheet both subscribe to it.
func ExpoImportProgressStore() string {
	return `import { useSyncExternalStore } from "react";
import { importResourceCsv, type ImportResult } from "./import";

export interface ImportProgress {
  id: string;
  plural: string;
  label: string;
  status: "processing" | "completed" | "failed";
  fraction: number;
  result?: ImportResult;
  error?: string;
}

let jobs: ImportProgress[] = [];
const listeners = new Set<() => void>();
let counter = 0;

function emit() {
  listeners.forEach((l) => l());
}
function patch(id: string, p: Partial<ImportProgress>) {
  jobs = jobs.map((j) => (j.id === id ? { ...j, ...p } : j));
  emit();
}

export function subscribeImports(l: () => void) {
  listeners.add(l);
  return () => {
    listeners.delete(l);
  };
}
export function getImports() {
  return jobs;
}
export function dismissImport(id: string) {
  jobs = jobs.filter((j) => j.id !== id);
  emit();
}

// Kick off a background import. Uploading + polling continue independently of
// any screen, so the user can close the sheet and navigate away. Returns the
// local job id so a still-open sheet can render its own progress/summary.
export function startImport(
  plural: string,
  fileUri: string,
  label: string,
  onDone?: () => void,
): string {
  const id = "imp-" + counter++;
  jobs = [...jobs, { id, plural, label, status: "processing", fraction: 0 }];
  emit();
  (async () => {
    try {
      const result = await importResourceCsv(plural, fileUri, (f) => patch(id, { fraction: f }));
      patch(id, { status: "completed", fraction: 1, result });
      onDone?.();
    } catch (e: any) {
      patch(id, { status: "failed", error: e?.message || "Import failed" });
    }
  })();
  return id;
}

// useImports subscribes a component to the live list of background imports.
export function useImports(): ImportProgress[] {
  return useSyncExternalStore(subscribeImports, getImports, getImports);
}
`
}

// ExpoImportBanner is components/ui/import-progress-banner.tsx — a persistent
// bottom banner that shows every in-flight/recent background import. Mounted
// once at the root so it survives navigation; the import sheet can be closed
// while an import keeps running here.
func ExpoImportBanner() string {
	return `import { View, Text, Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useImports, dismissImport, type ImportProgress } from "@/lib/import-progress";

export function ImportProgressBanner() {
  const jobs = useImports();
  const insets = useSafeAreaInsets();
  if (!jobs.length) return null;
  return (
    <View
      pointerEvents="box-none"
      style={{ position: "absolute", left: 12, right: 12, bottom: insets.bottom + 84 }}
    >
      {jobs.map((j) => (
        <BannerRow key={j.id} job={j} />
      ))}
    </View>
  );
}

function BannerRow({ job }: { job: ImportProgress }) {
  const done = job.status === "completed";
  const failed = job.status === "failed";
  const icon = failed ? "alert-circle" : done ? "checkmark-circle" : "cloud-upload-outline";
  const tint = failed ? "#ff6b6b" : done ? "#00b894" : "#6c5ce7";
  return (
    <View className="bg-white dark:bg-[#1a1a24] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-3 mb-2 shadow-lg">
      <View className="flex-row items-center">
        <Ionicons name={icon as any} size={18} color={tint} />
        <Text className="flex-1 text-[13px] font-semibold text-[#0F1018] dark:text-white ml-2" numberOfLines={1}>
          {failed
            ? "Import failed"
            : done
              ? "Import complete — " + (job.result?.created ?? 0) + " added"
              : "Importing " + job.label}
        </Text>
        {done || failed ? (
          <Pressable onPress={() => dismissImport(job.id)} hitSlop={8}>
            <Ionicons name="close" size={16} color="#9CA3AF" />
          </Pressable>
        ) : (
          <Text className="text-[12px] text-[#6B7280] dark:text-[#9090a8]">{Math.round(job.fraction * 100)}%</Text>
        )}
      </View>
      {!done && !failed ? (
        <View className="w-full h-1.5 rounded-full bg-[#E5E7EB] dark:bg-[#2a2a3a] overflow-hidden mt-2">
          <View style={{ width: Math.round(job.fraction * 100) + "%" }} className="h-1.5 bg-[#6c5ce7]" />
        </View>
      ) : null}
      {failed && job.error ? (
        <Text className="text-[12px] text-[#ff6b6b] mt-1" numberOfLines={2}>{job.error}</Text>
      ) : null}
      {done && job.result ? (
        <Text className="text-[12px] text-[#6B7280] dark:text-[#9090a8] mt-1">
          {job.result.created} created · {job.result.skipped} skipped · {job.result.failed} failed
        </Text>
      ) : null}
    </View>
  );
}
`
}

// expoImportSheet is components/ui/import-sheet.tsx — the pick → preview →
// progress → summary import flow. Kept in sync with internal/generate
// mobileImportSheetContent.
func ExpoImportSheet() string {
	return `import { useState } from "react";
import { View, Text, Pressable, ActivityIndicator, ScrollView } from "react-native";
import * as DocumentPicker from "expo-document-picker";
import { Ionicons } from "@expo/vector-icons";
import { FormSheet } from "./form-sheet";
import { parseCsvPreview, downloadResourceTemplate, type CsvPreview } from "@/lib/import";
import { startImport, useImports } from "@/lib/import-progress";

interface ImportSheetProps {
  plural: string;
  visible: boolean;
  onClose: () => void;
  onImported?: () => void;
}

export function ImportSheet({ plural, visible, onClose, onImported }: ImportSheetProps) {
  const [fileUri, setFileUri] = useState<string | null>(null);
  const [fileName, setFileName] = useState("");
  const [preview, setPreview] = useState<CsvPreview | null>(null);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [error, setError] = useState("");

  // The import runs in a module-level store, so it keeps going after this sheet
  // closes — we just read its live state back to show progress while open.
  const jobs = useImports();
  const job = jobs.find((j) => j.id === activeId) || null;

  const reset = () => {
    setFileUri(null);
    setFileName("");
    setPreview(null);
    setActiveId(null);
    setError("");
  };
  const close = () => {
    reset();
    onClose();
  };

  const pick = async () => {
    setError("");
    const res = await DocumentPicker.getDocumentAsync({
      type: ["text/csv", "text/comma-separated-values", "application/csv", "*/*"],
      copyToCacheDirectory: true,
    });
    if (res.canceled || !res.assets?.length) return;
    const asset = res.assets[0];
    setFileUri(asset.uri);
    setFileName(asset.name);
    setActiveId(null);
    try {
      setPreview(await parseCsvPreview(asset.uri));
    } catch (e: any) {
      setError(e.message || "Could not read the file");
    }
  };

  const runImport = () => {
    if (!fileUri) return;
    setError("");
    // Hand off to the background store; the sheet can be closed immediately and
    // the persistent banner will keep showing progress + the result.
    setActiveId(startImport(plural, fileUri, fileName || plural, onImported));
  };

  const onTemplate = async () => {
    setError("");
    try {
      await downloadResourceTemplate(plural);
    } catch (e: any) {
      setError(e.message || "Could not download template");
    }
  };

  const cw = 130;
  const labelClass = "text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8]";
  const displayError = error || (job?.status === "failed" ? job.error || "Import failed" : "");

  return (
    <FormSheet visible={visible} onClose={close} title="Import CSV">
      {displayError ? (
        <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center">
          <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
          <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{displayError}</Text>
        </View>
      ) : null}

      {job?.status === "completed" && job.result ? (
        <View>
          <View className="items-center mb-5">
            <Ionicons name="checkmark-circle" size={44} color="#00b894" />
            <Text className="text-[17px] font-bold text-[#0F1018] dark:text-white mt-2">Import complete</Text>
          </View>
          <SummaryRow label="Created" value={job.result.created} color="#00b894" />
          <SummaryRow label="Skipped (duplicates)" value={job.result.skipped} color="#fdcb6e" />
          <SummaryRow label="Failed" value={job.result.failed} color="#ff6b6b" />
          {job.result.errors.length > 0 ? (
            <View className="mt-3">
              <Text className={labelClass + " mb-2"}>Errors</Text>
              {job.result.errors.slice(0, 20).map((e, idx) => (
                <Text key={idx} className="text-[12px] text-[#6B7280] dark:text-[#9090a8] mb-1">
                  Row {e.row}: {e.message}
                </Text>
              ))}
            </View>
          ) : null}
          <Pressable onPress={close} className="bg-[#6c5ce7] rounded-full py-4 items-center mt-5">
            <Text className="text-white font-semibold text-[15px]">Done</Text>
          </Pressable>
        </View>
      ) : job?.status === "processing" ? (
        <View className="py-8 items-center">
          <Text className="text-[15px] text-[#0F1018] dark:text-white mb-4">Importing {fileName}…</Text>
          <View className="w-full h-2.5 rounded-full bg-[#E5E7EB] dark:bg-[#2a2a3a] overflow-hidden">
            <View style={{ width: Math.round(job.fraction * 100) + "%" }} className="h-2.5 bg-[#6c5ce7]" />
          </View>
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-2">{Math.round(job.fraction * 100)}%</Text>
          <Pressable onPress={close} className="rounded-full py-3 px-8 items-center mt-4">
            <Text className="text-[#6c5ce7] font-semibold text-[14px]">Continue in background</Text>
          </Pressable>
        </View>
      ) : preview ? (
        <View>
          <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white mb-1">{fileName}</Text>
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mb-3">{preview.total} rows to import</Text>
          <ScrollView horizontal showsHorizontalScrollIndicator className="mb-2" style={{ maxHeight: 260 }}>
            <View>
              <View className="flex-row border-b-2 border-[#E5E7EB] dark:border-[#2a2a3a]">
                {preview.headers.map((h, i) => (
                  <View key={i} style={{ width: cw }} className="px-2 py-2">
                    <Text className={labelClass} numberOfLines={1}>{h}</Text>
                  </View>
                ))}
              </View>
              {preview.rows.map((r, ri) => (
                <View key={ri} className="flex-row border-b border-[#E5E7EB] dark:border-[#1f1f2b]">
                  {preview.headers.map((_, ci) => (
                    <View key={ci} style={{ width: cw }} className="px-2 py-2">
                      <Text className="text-[13px] text-[#0F1018] dark:text-white" numberOfLines={1}>{r[ci] ?? ""}</Text>
                    </View>
                  ))}
                </View>
              ))}
            </View>
          </ScrollView>
          {preview.total > preview.rows.length ? (
            <Text className="text-[12px] text-[#9CA3AF] dark:text-[#606078] mb-3">
              +{preview.total - preview.rows.length} more rows
            </Text>
          ) : null}
          <Pressable onPress={runImport} className="bg-[#6c5ce7] rounded-full py-4 items-center">
            <Text className="text-white font-semibold text-[15px]">Import {preview.total} rows</Text>
          </Pressable>
          <Pressable onPress={pick} className="rounded-full py-3 items-center mt-2">
            <Text className="text-[#6c5ce7] font-semibold text-[14px]">Choose a different file</Text>
          </Pressable>
        </View>
      ) : (
        <View className="py-8 items-center">
          <Ionicons name="document-attach-outline" size={44} color="#9CA3AF" />
          <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white mt-3">Import from CSV</Text>
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] text-center mt-1 px-6">
            The first row should be column headers that match the field names.
          </Text>
          <Pressable onPress={pick} className="bg-[#6c5ce7] rounded-full py-4 px-8 items-center mt-5">
            <Text className="text-white font-semibold text-[15px]">Choose CSV file</Text>
          </Pressable>
          <Pressable onPress={onTemplate} className="rounded-full py-3 px-8 items-center mt-2 flex-row">
            <Ionicons name="download-outline" size={16} color="#6c5ce7" />
            <Text className="text-[#6c5ce7] font-semibold text-[14px] ml-2">Download template</Text>
          </Pressable>
        </View>
      )}
    </FormSheet>
  );
}

function SummaryRow({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <View className="flex-row items-center justify-between py-2.5 border-b border-[#E5E7EB] dark:border-[#1f1f2b]">
      <Text className="text-[14px] text-[#6B7280] dark:text-[#9090a8]">{label}</Text>
      <Text style={{ color }} className="text-[16px] font-bold">{value}</Text>
    </View>
  );
}
`
}

// ExpoRelationSelect is components/ui/relation-select.tsx — a searchable
// single-select for belongs_to fields. Pills don't scale past a handful of
// options, so this opens a bottom sheet with a pinned search box and a
// scrollable, filterable list. Used by generated resource forms for every
// belongs_to relationship.
func ExpoRelationSelect() string {
	return `import { useState } from "react";
import { View, Text, Pressable, TextInput, FlatList, Modal } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useSafeAreaInsets } from "react-native-safe-area-context";

interface RelationOption {
  id: string;
  [key: string]: any;
}

interface RelationSelectProps {
  label: string;
  value: string;
  onChange: (id: string) => void;
  options: RelationOption[];
  placeholder?: string;
}

function optionLabel(o: RelationOption): string {
  return String(o.name ?? o.title ?? o.label ?? o.id);
}

export function RelationSelect({ label, value, onChange, options, placeholder }: RelationSelectProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const insets = useSafeAreaInsets();

  const selected = options.find((o) => o.id === value);
  const q = search.trim().toLowerCase();
  const filtered = q ? options.filter((o) => optionLabel(o).toLowerCase().includes(q)) : options;

  const close = () => {
    setOpen(false);
    setSearch("");
  };

  return (
    <View className="mb-4">
      <Text className="text-[13px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2">{label}</Text>
      <Pressable
        onPress={() => setOpen(true)}
        className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-3.5 flex-row items-center justify-between"
      >
        <Text
          className={selected ? "text-[15px] text-[#0F1018] dark:text-white" : "text-[15px] text-[#9CA3AF]"}
          numberOfLines={1}
        >
          {selected ? optionLabel(selected) : placeholder || "Select " + label}
        </Text>
        <Ionicons name="chevron-down" size={18} color="#9CA3AF" />
      </Pressable>

      <Modal visible={open} transparent animationType="slide" onRequestClose={close}>
        <Pressable className="flex-1 bg-black/40" onPress={close} />
        <View
          className="bg-[#F4F4F6] dark:bg-[#0a0a0f] rounded-t-3xl"
          style={{ maxHeight: "72%", paddingBottom: insets.bottom + 12 }}
        >
          <View className="items-center pt-3 pb-1">
            <View className="w-10 h-1 rounded-full bg-[#D1D5DB] dark:bg-[#2a2a3a]" />
          </View>
          <View className="flex-row items-center justify-between px-5 pb-3">
            <Text className="text-[17px] font-bold text-[#0F1018] dark:text-white">Select {label}</Text>
            <Pressable onPress={close} hitSlop={8}>
              <Ionicons name="close" size={22} color="#9CA3AF" />
            </Pressable>
          </View>
          <View className="px-5 pb-2">
            <View
              className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl flex-row items-center px-4"
              style={{ height: 46 }}
            >
              <Ionicons name="search-outline" size={18} color="#9CA3AF" />
              <TextInput
                className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
                placeholder={"Search " + label + "…"}
                placeholderTextColor="#9CA3AF"
                value={search}
                onChangeText={setSearch}
                autoCapitalize="none"
                autoFocus
              />
            </View>
          </View>
          <FlatList
            data={filtered}
            keyExtractor={(o) => String(o.id)}
            keyboardShouldPersistTaps="handled"
            contentContainerStyle={{ paddingHorizontal: 20, paddingBottom: 8 }}
            renderItem={({ item }) => {
              const isSel = item.id === value;
              return (
                <Pressable
                  onPress={() => {
                    onChange(item.id);
                    close();
                  }}
                  className="flex-row items-center justify-between px-4 py-3.5 mb-2 rounded-2xl bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b]"
                >
                  <Text className="text-[15px] text-[#0F1018] dark:text-white flex-1" numberOfLines={1}>
                    {optionLabel(item)}
                  </Text>
                  {isSel ? <Ionicons name="checkmark-circle" size={20} color="#6c5ce7" /> : null}
                </Pressable>
              );
            }}
            ListEmptyComponent={<Text className="text-center text-[#9CA3AF] mt-8">No matches</Text>}
          />
        </View>
      </Modal>
    </View>
  );
}
`
}

// ExpoNumberFormat is lib/format.ts — thousands-separator helpers for numeric
// inputs. Typing 1000 shows "1,000" while the payload still carries the plain
// number 1000. Values are NOT scaled (no cents conversion): what you type is
// what's stored.
func ExpoNumberFormat() string {
	return `// Format a numeric text input with thousands separators as the user types:
// "1000" -> "1,000". Set allowDecimal for float fields (keeps up to 2 decimals).
// The value is never scaled — 100 means 100, not cents.
export function formatNumberInput(value: string, allowDecimal = false): string {
  if (value === null || value === undefined || value === "") return "";
  let cleaned = allowDecimal
    ? String(value).replace(/[^0-9.]/g, "")
    : String(value).replace(/[^0-9]/g, "");

  if (allowDecimal) {
    const parts = cleaned.split(".");
    if (parts.length > 1) {
      cleaned = parts[0] + "." + parts.slice(1).join("").slice(0, 2);
    }
  }

  const [intPart, decPart] = cleaned.split(".");
  const withCommas = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  return decPart !== undefined ? withCommas + "." + decPart : withCommas;
}

// Strip the separators back to a plain number for the API payload.
export function parseNumberInput(value: string): number {
  const n = Number(String(value ?? "").replace(/,/g, ""));
  return Number.isFinite(n) ? n : 0;
}
`
}

// ExpoImageResolver is lib/images.ts — resolveImageUrl(). Dev storage (MinIO)
// hands back URLs like http://localhost:9002/... which only resolve on the
// machine running the server; a device/emulator can't reach "localhost". This
// rewrites the host to the same dev host the app already uses for the API so
// stored images actually load in lists, detail screens and previews. In
// production (real S3/R2 public URLs) it's a no-op. Used everywhere an <Image>
// renders a stored file URL.
func ExpoImageResolver() string {
	return `import Constants from "expo-constants";

// The dev host the app already uses to reach Metro / the API. A device or
// emulator can't reach "localhost"/"127.0.0.1" — those point at the device
// itself — so we reuse this host to rewrite storage URLs below.
function devHost(): string | undefined {
  const hostUri =
    Constants.expoConfig?.hostUri ||
    (Constants.expoGoConfig as any)?.debuggerHost ||
    (Constants.manifest2 as any)?.extra?.expoGo?.debuggerHost;
  const host = hostUri ? String(hostUri).split(":")[0] : undefined;
  return host && host !== "localhost" && host !== "127.0.0.1" ? host : undefined;
}

// resolveImageUrl makes a stored file URL loadable on a real device/emulator.
// Dev storage returns http://localhost:<port>/... which only resolves on the
// server machine; we rewrite the host to the dev host the app talks to for the
// API, keeping the storage port. Production public URLs pass through unchanged.
export function resolveImageUrl(url?: string | null): string | undefined {
  if (!url) return undefined;
  const host = devHost();
  if (host) {
    return url.replace(/^(https?:\/\/)(localhost|127\.0\.0\.1)/i, ` + "`" + `$1${host}` + "`" + `);
  }
  return url;
}
`
}

func ExpoUploadHelper() string {
	return `import * as ImagePicker from "expo-image-picker";
import * as SecureStore from "expo-secure-store";
import { API_URL } from "./api";

// Pick an image and upload it to POST /uploads (multipart). Returns the public
// URL of the stored file, or null if the user cancelled.
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

  return uploadLocalFile(result.assets[0].uri, result.assets[0].fileName, result.assets[0].mimeType);
}

// uploadLocalFile POSTs a local file:// URI to /uploads as multipart/form-data.
//
// It uses fetch + FormData with a React Native file descriptor ({ uri, name,
// type }). Two rules make this reliable on Expo SDK 54 / New Architecture:
//   1. NEVER set a Content-Type header — fetch derives "multipart/form-data"
//      WITH the boundary from the FormData. Setting it by hand drops the
//      boundary and the server can't find the file part ("No file provided").
//   2. Give the descriptor a real name + type so the server's MIME allowlist
//      accepts it.
// This replaces expo-file-system's uploadAsync, which sends an EMPTY body under
// the New Architecture on SDK 54 (the request lands in a few ms with no file).
export async function uploadLocalFile(
  uri: string,
  fileName?: string | null,
  mimeType?: string | null,
): Promise<string> {
  const name = fileName || uri.split("/").pop() || "photo.jpg";
  const ext = (name.split(".").pop() || "jpg").toLowerCase();
  const extToMime: Record<string, string> = {
    jpg: "image/jpeg",
    jpeg: "image/jpeg",
    png: "image/png",
    gif: "image/gif",
    webp: "image/webp",
    heic: "image/heic",
  };
  const type = mimeType || extToMime[ext] || "image/jpeg";

  const token = await SecureStore.getItemAsync("access_token");

  const formData = new FormData();
  // RN's FormData accepts a { uri, name, type } file descriptor — the cast is
  // needed because the DOM FormData types don't model it.
  formData.append("file", { uri, name, type } as any);

  const res = await fetch(API_URL + "/uploads", {
    method: "POST",
    headers: token ? { Authorization: "Bearer " + token } : undefined,
    body: formData,
  });

  const text = await res.text();
  const json = text ? JSON.parse(text) : null;
  if (!res.ok) {
    // Surface the server's reason (e.g. "File type not allowed") instead of a
    // generic failure, so upload problems are diagnosable from the UI.
    throw new Error(json?.error?.message || "Upload failed (" + res.status + ")");
  }
  return json?.data?.url ?? null;
}
`
}

// ExpoImagePickerSheet is components/ui/image-picker-sheet.tsx — a themed
// bottom sheet for choosing an image. It replaces the old "tap the dropzone →
// jump straight into Android's system crop screen" flow, which was confusing
// (the only button was CROP and there was no explicit "use this photo"). Now
// the user gets: a clean sheet → grant permission → pick from library or camera
// → PREVIEW → decide to Use, Crop, or choose a different photo. Cropping (the
// native editor) only happens when the user explicitly taps Crop. Supports
// multi-select for "files" fields. onSelect returns the chosen local URIs; the
// caller uploads them via uploadLocalFile.
func ExpoImagePickerSheet() string {
	return `import { useState } from "react";
import { View, Text, Pressable, Image, ActivityIndicator, ScrollView, Linking } from "react-native";
import * as ImagePicker from "expo-image-picker";
import { Ionicons } from "@expo/vector-icons";
import { FormSheet } from "./form-sheet";

interface ImagePickerSheetProps {
  visible: boolean;
  onClose: () => void;
  multiple?: boolean;
  onSelect: (uris: string[]) => void;
}

export function ImagePickerSheet({ visible, onClose, multiple, onSelect }: ImagePickerSheetProps) {
  const [assets, setAssets] = useState<string[]>([]);
  const [permError, setPermError] = useState("");
  const [busy, setBusy] = useState(false);

  const reset = () => {
    setAssets([]);
    setPermError("");
    setBusy(false);
  };
  const close = () => {
    reset();
    onClose();
  };

  const fromLibrary = async (crop = false) => {
    setPermError("");
    const perm = await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (!perm.granted) {
      setPermError("Photo library access is off. Turn it on in Settings to choose photos.");
      return;
    }
    setBusy(true);
    try {
      const res = await ImagePicker.launchImageLibraryAsync({
        mediaTypes: ["images"],
        allowsEditing: crop,
        aspect: [1, 1],
        quality: 0.8,
        allowsMultipleSelection: !crop && !!multiple,
        selectionLimit: multiple ? 10 : 1,
      });
      if (!res.canceled && res.assets?.length) {
        setAssets(res.assets.map((a) => a.uri));
      }
    } finally {
      setBusy(false);
    }
  };

  const fromCamera = async () => {
    setPermError("");
    const perm = await ImagePicker.requestCameraPermissionsAsync();
    if (!perm.granted) {
      setPermError("Camera access is off. Turn it on in Settings to take a photo.");
      return;
    }
    setBusy(true);
    try {
      const res = await ImagePicker.launchCameraAsync({
        mediaTypes: ["images"],
        allowsEditing: true,
        aspect: [1, 1],
        quality: 0.8,
      });
      if (!res.canceled && res.assets?.length) setAssets(res.assets.map((a) => a.uri));
    } finally {
      setBusy(false);
    }
  };

  const use = () => {
    const chosen = assets;
    onSelect(chosen);
    close();
  };

  return (
    <FormSheet visible={visible} onClose={close} title={multiple ? "Add photos" : "Add photo"}>
      {permError ? (
        <View className="bg-[#fdcb6e]/10 border border-[#fdcb6e]/30 rounded-2xl px-4 py-3 mb-4">
          <Text className="text-[13px] text-[#0F1018] dark:text-white mb-2">{permError}</Text>
          <Pressable onPress={() => Linking.openSettings()}>
            <Text className="text-[#6c5ce7] font-semibold text-[13px]">Open Settings</Text>
          </Pressable>
        </View>
      ) : null}

      {assets.length ? (
        <View>
          {assets.length === 1 ? (
            <Image source={{ uri: assets[0] }} style={{ width: "100%", height: 220, borderRadius: 16 }} />
          ) : (
            <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mb-1">
              {assets.map((u, i) => (
                <Image key={i} source={{ uri: u }} style={{ width: 96, height: 96, borderRadius: 12, marginRight: 8 }} />
              ))}
            </ScrollView>
          )}
          <Pressable onPress={use} className="bg-[#6c5ce7] rounded-full py-4 items-center mt-4">
            <Text className="text-white font-semibold text-[15px]">
              {multiple ? "Use " + assets.length + " photo" + (assets.length > 1 ? "s" : "") : "Use photo"}
            </Text>
          </Pressable>
          {!multiple ? (
            <Pressable onPress={() => fromLibrary(true)} className="rounded-full py-3 items-center mt-2 flex-row justify-center">
              <Ionicons name="crop-outline" size={16} color="#6c5ce7" />
              <Text className="text-[#6c5ce7] font-semibold text-[14px] ml-2">Crop instead</Text>
            </Pressable>
          ) : null}
          <Pressable onPress={() => setAssets([])} className="rounded-full py-3 items-center">
            <Text className="text-[#6B7280] dark:text-[#9090a8] font-semibold text-[14px]">Choose different</Text>
          </Pressable>
        </View>
      ) : busy ? (
        <View className="py-12 items-center">
          <ActivityIndicator color="#6c5ce7" />
          <Text className="text-[13px] text-[#6B7280] dark:text-[#9090a8] mt-3">Opening…</Text>
        </View>
      ) : (
        <View className="py-1">
          <SourceButton
            icon="images-outline"
            label="Choose from library"
            hint={multiple ? "Select one or more photos" : "Pick a photo from your gallery"}
            onPress={() => fromLibrary(false)}
          />
          <SourceButton icon="camera-outline" label="Take a photo" hint="Use your camera" onPress={fromCamera} />
        </View>
      )}
    </FormSheet>
  );
}

function SourceButton({
  icon,
  label,
  hint,
  onPress,
}: {
  icon: string;
  label: string;
  hint: string;
  onPress: () => void;
}) {
  return (
    <Pressable
      onPress={onPress}
      className="flex-row items-center bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl px-4 py-4 mb-3"
    >
      <View className="w-11 h-11 rounded-full bg-[#6c5ce7]/10 items-center justify-center">
        <Ionicons name={icon as any} size={22} color="#6c5ce7" />
      </View>
      <View className="ml-3 flex-1">
        <Text className="text-[15px] font-semibold text-[#0F1018] dark:text-white">{label}</Text>
        <Text className="text-[12px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">{hint}</Text>
      </View>
      <Ionicons name="chevron-forward" size={18} color="#9CA3AF" />
    </Pressable>
  );
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
