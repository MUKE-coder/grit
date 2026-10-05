import { useQuery } from "@tanstack/react-query";
import type { Notification as NotificationRow } from "@repo/shared/types";
import { apiClient } from "@/lib/api-client";

// The row is the model's, from grit sync, with the two strings that are really
// enums narrowed so the bell can switch on them.
export type Notification = Omit<NotificationRow, "source" | "severity"> & {
  source: "sentinel" | "pulse" | "system";
  severity: "critical" | "high" | "medium" | "low" | "info";
};

export interface NotificationList {
  data: Notification[];
  unread: number;
}

// One key for the notification list wherever it is read. The bell and the
// dashboard's tile used to ask under keys of their own, so the dashboard
// polled the same URL twice a minute, and marking one read in the bell left
// the tile's count where it was.
export const notificationKeys = {
  all: ["notifications"] as const,
  list: () => ["notifications", "list"] as const,
};

// useNotificationList polls the list once a minute. Pass select to read a part
// of it, such as the unread count, without a request of its own.
export function useNotificationList<TSelected = NotificationList>(
  select?: (list: NotificationList) => TSelected,
) {
  return useQuery<NotificationList, Error, TSelected>({
    queryKey: notificationKeys.list(),
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<NotificationList>("/api/notifications");
        return data;
      } catch {
        // The bell is on every page, and a notifications outage should not
        // break any of them.
        return { data: [], unread: 0 };
      }
    },
    select,
    refetchInterval: 60_000,
    retry: false,
    staleTime: 30_000,
  });
}
