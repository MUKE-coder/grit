import Constants from "expo-constants";
import * as Notifications from "expo-notifications";
import { Platform } from "react-native";
import { api } from "@/lib/api";

// Show a notification that arrives while the app is open, instead of
// swallowing it. Screens that already show the thing live (an open chat, say)
// can decide otherwise in their own handler.
Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowBanner: true,
    shouldShowList: true,
    shouldPlaySound: true,
    shouldSetBadge: false,
  }),
});

let registered: string | null = null;

/**
 * Asks for permission, gets this device's Expo push token and registers it with
 * the API. Call it after sign-in, and on every launch while signed in: that is
 * how the API tells a live device from one whose app was deleted. Returns the
 * token, or null when push is not available (permission refused, a simulator,
 * the web).
 */
export async function registerForPush(): Promise<string | null> {
  if (Platform.OS === "web") return null;
  try {
    if (Platform.OS === "android") {
      await Notifications.setNotificationChannelAsync("default", {
        name: "Default",
        importance: Notifications.AndroidImportance.HIGH,
      });
    }
    const current = await Notifications.getPermissionsAsync();
    const status = current.granted ? current : await Notifications.requestPermissionsAsync();
    if (!status.granted) return null;

    // An EAS build knows its project id; Expo Go reads it from app config.
    const projectId = Constants.expoConfig?.extra?.eas?.projectId ?? Constants.easConfig?.projectId;
    const { data: token } = await Notifications.getExpoPushTokenAsync(projectId ? { projectId } : undefined);
    await api.post("/push/tokens", { token, platform: Platform.OS });
    registered = token;
    return token;
  } catch (err) {
    // A device that cannot get a token (a simulator, no network) is not an
    // error worth showing: the app works without push.
    console.warn("[push] not registered:", err);
    return null;
  }
}

/**
 * Stops notifications to this device. Call it before signing out, while the
 * session can still prove the token is this user's.
 */
export async function unregisterPush(): Promise<void> {
  if (!registered) return;
  const token = registered;
  registered = null;
  try {
    await api.post("/push/tokens/remove", { token });
  } catch {
    // Signing out goes ahead regardless. The next person to sign in on this
    // phone takes the token over, which moves it off this account anyway.
  }
}

/**
 * Calls onOpen with a notification's data when the user taps it, including the
 * tap that launched the app. Returns the unsubscribe function.
 */
export function onNotificationTap(onOpen: (data: Record<string, string>) => void): () => void {
  const open = (response: Notifications.NotificationResponse | null) => {
    const data = response?.notification.request.content.data;
    if (data) onOpen(data as Record<string, string>);
  };
  void Notifications.getLastNotificationResponseAsync().then(open);
  const sub = Notifications.addNotificationResponseReceivedListener(open);
  return () => sub.remove();
}
