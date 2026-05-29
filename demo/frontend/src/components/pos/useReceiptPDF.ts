import { useCallback, useRef, useState } from 'react'
import { pdf } from '@react-pdf/renderer'
import type { ReactElement } from 'react'

/**
 * Imperative PDF generator + download + print, using @react-pdf/renderer's
 * `pdf()` factory. We avoid <PDFDownloadLink> because in React 19 + suspense
 * mode it's been flaky — sometimes never resolves the loading state. The
 * imperative API is straightforward: build → blob → done.
 *
 * Print uses a hidden iframe (cross-browser reliable) instead of
 * `window.open()` which is increasingly blocked by popup-blockers.
 */
export function useReceiptPDF() {
  const [busy, setBusy] = useState<null | 'download' | 'print' | 'preview'>(null)
  const printFrame = useRef<HTMLIFrameElement | null>(null)
  const lastUrl = useRef<string | null>(null)

  const buildBlob = useCallback(async (doc: ReactElement) => {
    // pdf().toBlob() → returns a Blob we can hand to URL.createObjectURL.
    return await pdf(doc).toBlob()
  }, [])

  const download = useCallback(async (doc: ReactElement, filename: string) => {
    try {
      setBusy('download')
      const blob = await buildBlob(doc)
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename.endsWith('.pdf') ? filename : `${filename}.pdf`
      document.body.appendChild(a)
      a.click()
      a.remove()
      // give the browser a beat to start the download before revoking
      setTimeout(() => URL.revokeObjectURL(url), 1500)
    } finally {
      setBusy(null)
    }
  }, [buildBlob])

  const print = useCallback(async (doc: ReactElement) => {
    try {
      setBusy('print')
      const blob = await buildBlob(doc)
      const url = URL.createObjectURL(blob)

      // Re-use the same hidden iframe across prints so we don't leak.
      if (printFrame.current) {
        document.body.removeChild(printFrame.current)
        printFrame.current = null
      }
      if (lastUrl.current) URL.revokeObjectURL(lastUrl.current)
      lastUrl.current = url

      const iframe = document.createElement('iframe')
      iframe.style.position = 'fixed'
      iframe.style.right = '0'
      iframe.style.bottom = '0'
      iframe.style.width = '0'
      iframe.style.height = '0'
      iframe.style.border = '0'
      iframe.src = url
      iframe.onload = () => {
        try {
          iframe.contentWindow?.focus()
          iframe.contentWindow?.print()
        } catch {
          // Some browsers won't allow printing a blob iframe — fall back to a
          // new-tab open so the user can hit Cmd/Ctrl+P themselves.
          window.open(url, '_blank', 'noopener,noreferrer')
        }
      }
      document.body.appendChild(iframe)
      printFrame.current = iframe
    } finally {
      setBusy(null)
    }
  }, [buildBlob])

  /** Returns a blob URL the caller can stick into an <iframe src=…/> for in-app preview. */
  const preview = useCallback(async (doc: ReactElement): Promise<string> => {
    setBusy('preview')
    try {
      const blob = await buildBlob(doc)
      const url = URL.createObjectURL(blob)
      if (lastUrl.current) URL.revokeObjectURL(lastUrl.current)
      lastUrl.current = url
      return url
    } finally {
      setBusy(null)
    }
  }, [buildBlob])

  return { download, print, preview, busy }
}
