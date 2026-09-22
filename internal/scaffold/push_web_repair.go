package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// grit plugin add push wrote an onNotificationTap that asked expo-notifications
// for the notification that opened the app, which it cannot answer on the web
// and throws instead: an Expo app with the push plugin crashed on its first
// screen in a browser. Found making mockups of the WhatsApp blueprint.
const (
	pushTapOld = `export function onNotificationTap(onOpen: (data: Record<string, string>) => void): () => void {
  const open = (response: Notifications.NotificationResponse | null) => {
    const data = response?.notification.request.content.data;
    if (data) onOpen(data as Record<string, string>);
  };
  void Notifications.getLastNotificationResponseAsync().then(open);
`
	pushTapNew = `export function onNotificationTap(onOpen: (data: Record<string, string>) => void): () => void {
  // The web has no push notifications, and expo-notifications throws there
  // rather than answering "none": the app crashed on its first screen when
  // opened in a browser.
  if (Platform.OS === "web") return () => {};
  const open = (response: Notifications.NotificationResponse | null) => {
    const data = response?.notification.request.content.data;
    if (data) onOpen(data as Record<string, string>);
  };
  Notifications.getLastNotificationResponseAsync().then(open, () => undefined);
`
)

func repairPushWebSource(src string) (string, []string, []string) {
	if strings.Count(src, pushTapOld) != 1 {
		return src, nil, nil
	}
	return strings.Replace(src, pushTapOld, pushTapNew, 1),
		[]string{"the push tap listener does nothing on the web instead of crashing the app"}, nil
}

// repairPushWeb applies it to an Expo app that has the push plugin.
func repairPushWeb(root string) error {
	path := filepath.Join(root, "apps", "expo", "lib", "push.ts")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairTextFile(root, m, path, repairPushWebSource)
}
