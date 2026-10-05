export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

export interface ApiError {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

// apiErrorMessage extracts a human-readable message from any error
// that came out of the API client (axios or fetch). The standard error
// envelope is { error: { code, message, details? } }; this walks the
// usual axios.error.response.data.error.message chain plus a few
// fallbacks so toast.error(apiErrorMessage(err)) is always meaningful.
//
// Use it everywhere instead of err.message directly:
//
//   import { apiErrorMessage } from "@grit/shared/types/api";
//   ...
//   onError: (err) => toast.error(apiErrorMessage(err))
export function apiErrorMessage(err: unknown, fallback = "Something went wrong"): string {
  if (!err) return fallback;
  if (typeof err === "string") return err;

  // Standard envelope — preferred path
  const e = err as {
    response?: { data?: { error?: { message?: string; code?: string } } };
    message?: string;
  };
  const envMsg = e.response?.data?.error?.message;
  if (envMsg) return envMsg;

  // Axios error.message ("Network Error", "timeout of 5000ms exceeded")
  if (e.message) return e.message;

  // Last resort
  return fallback;
}

// apiErrorCode returns the standard envelope code (VALIDATION_ERROR,
// NOT_FOUND, ...) when present. Useful for branching on specific codes
// (e.g. "VERSION_CONFLICT" -> open conflict dialog).
export function apiErrorCode(err: unknown): string | null {
  if (!err) return null;
  const e = err as { response?: { data?: { error?: { code?: string } } } };
  return e.response?.data?.error?.code ?? null;
}

// apiErrorFields surfaces per-field validation errors when the API
// returns details: { fieldName: "error message" }. Used to highlight
// individual inputs in a form.
export function apiErrorFields(err: unknown): Record<string, string> {
  if (!err) return {};
  const e = err as { response?: { data?: { error?: { details?: Record<string, string> } } } };
  return e.response?.data?.error?.details ?? {};
}
