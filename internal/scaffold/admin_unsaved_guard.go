package scaffold

// The unsaved-changes guard.
//
// The admin had none. A resource form with twenty fields, a rich text body and
// three uploads could be thrown away by a stray click on the sidebar, and
// nothing anywhere said so: no prompt, no warning on tab close, no indication
// that there was anything to lose. The row went back to how it was and the
// typing was gone.
//
// Two halves, and both are needed because they catch different accidents.
//
// `beforeunload` catches the tab being closed, the page being reloaded and the
// back button leaving the site. It does not catch a click on a Next <Link>,
// which is a client-side navigation the browser never hears about, and in an
// admin panel that is the common one: you click Users in the sidebar with a
// half-filled form behind you.
//
// So the second half is a capture-phase click listener that finds the anchor a
// click is heading for and asks first. It uses window.confirm, which is ugly,
// and the uglier alternative was considered and rejected: the admin's own
// confirm modal is asynchronous, and by the time it resolves the navigation it
// was meant to stop has already happened. A native confirm is the only thing
// that can block a click synchronously. Looking right is worth less than
// working.
//
// Idea from github.com/jubayer910/Admin-panel (MIT), whose settings form keeps
// a sticky save bar reading "Unsaved changes" or "All changes saved" with a
// Discard beside it. Grit's forms are modals and pages rather than one long
// settings screen, so this is the same guarantee in a smaller shape.

// adminUnsavedGuard emits components/forms/unsaved-guard.tsx.
func adminUnsavedGuard() string {
	return `"use client";

import { useEffect } from "react";

const MESSAGE = "You have unsaved changes. Leave this page and lose them?";

/**
 * Warns before unsaved edits are thrown away.
 *
 * Pass whether the form is dirty. While it is true the hook holds two
 * listeners, and takes them off the moment it goes false, so a saved form
 * stops asking.
 *
 * It deliberately does not try to block the browser's own back button on a
 * client-side navigation. The History API gives no way to ask first without
 * pushing a decoy entry onto the stack, and a back button that needs pressing
 * twice is a worse bug than the one being fixed.
 */
export function useUnsavedGuard(dirty: boolean) {
  useEffect(() => {
    if (!dirty) return;

    // Closing the tab, reloading, or leaving the site.
    const onLeave = (event: BeforeUnloadEvent) => {
      event.preventDefault();
      // Chrome still wants returnValue set. No browser has shown this string
      // for years: they all display their own wording.
      event.returnValue = "";
    };

    // A click on a link, which beforeunload never hears because Next handles
    // it in the client. Capture phase, so this runs before the router does.
    const onClick = (event: MouseEvent) => {
      if (event.defaultPrevented || event.button !== 0) return;
      // A modified click opens a new tab and leaves this one alone.
      if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;

      const link = (event.target as HTMLElement | null)?.closest("a");
      if (!link) return;

      const href = link.getAttribute("href");
      if (!href || href.startsWith("#")) return;
      if (link.target && link.target !== "_self") return;
      if (link.hasAttribute("download")) return;

      let destination: URL;
      try {
        destination = new URL(link.href, window.location.href);
      } catch {
        return;
      }
      // Another origin leaves the app, which beforeunload already covers.
      if (destination.origin !== window.location.origin) return;
      // Going where we already are is not leaving.
      if (destination.pathname === window.location.pathname) return;

      if (window.confirm(MESSAGE)) return;
      event.preventDefault();
      event.stopPropagation();
    };

    window.addEventListener("beforeunload", onLeave);
    document.addEventListener("click", onClick, true);
    return () => {
      window.removeEventListener("beforeunload", onLeave);
      document.removeEventListener("click", onClick, true);
    };
  }, [dirty]);
}

/**
 * The line beside a save button that says whether there is anything to lose.
 *
 * Renders nothing when the form is untouched. A permanent "All changes saved"
 * is noise on a form nobody has typed in yet, and it trains people to stop
 * reading the place the warning will appear.
 */
export function UnsavedBadge({ dirty }: { dirty: boolean }) {
  if (!dirty) return null;
  return (
    <span role="status" className="mr-auto flex items-center gap-2 text-xs text-warning">
      <span className="h-1.5 w-1.5 rounded-full bg-warning" aria-hidden="true" />
      Unsaved changes
    </span>
  );
}
`
}
