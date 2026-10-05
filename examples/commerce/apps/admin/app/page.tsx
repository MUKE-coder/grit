import { redirect } from "next/navigation";

// The root of the panel only ever sends a visitor on, so the server does it: a
// 307 before any JavaScript loads. This used to be a client page that
// downloaded the app, asked /api/auth/me who was signed in, and only then
// navigated, so every visit to the root paid for a page it never showed.
//
// Nothing is lost by not asking: the dashboard's layout sends a signed-out
// visitor to the login, and a USER with no grants to their profile.
export default function RootPage() {
  redirect("/dashboard");
}
