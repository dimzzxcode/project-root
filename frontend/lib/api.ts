export type Item = {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
};

export type ItemInput = {
  name: string;
  description: string;
};

type ApiErrorBody = {
  code: string;
  message: string;
  details?: Record<string, string>;
};

type ApiEnvelope<T> = {
  success: boolean;
  data: T;
  error?: ApiErrorBody;
};

type RequestOptions = {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  signal?: AbortSignal;
};

/**
 * Error thrown for failed API calls. `details` contains per-field validation
 * messages when the backend responds with a validation error.
 */
export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  readonly details?: Record<string, string>;

  constructor(
    message: string,
    status: number,
    code: string,
    details?: Record<string, string>,
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

function apiBaseUrl(): string {
  const value = process.env.NEXT_PUBLIC_API_URL;

  if (!value) {
    throw new ApiError(
      "NEXT_PUBLIC_API_URL is not configured.",
      0,
      "CONFIG_ERROR",
    );
  }

  return value.replace(/\/+$/, "");
}

async function request<T>(
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const { method = "GET", body, signal } = options;
  const headers: Record<string, string> = { Accept: "application/json" };

  if (body !== undefined) {
    headers["Content-Type"] = "application/json";
  }

  let response: Response;
  try {
    response = await fetch(`${apiBaseUrl()}${path}`, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
      signal,
      cache: "no-store",
    });
  } catch (error) {
    if (error instanceof DOMException && error.name === "AbortError") {
      throw error;
    }
    throw new ApiError(
      "Could not reach the API. Make sure the backend is running.",
      0,
      "NETWORK_ERROR",
    );
  }

  const text = await response.text();
  let payload: ApiEnvelope<T> | null = null;

  if (text) {
    try {
      payload = JSON.parse(text) as ApiEnvelope<T>;
    } catch {
      throw new ApiError(
        `The API returned an invalid response (HTTP ${response.status}).`,
        response.status,
        "INVALID_RESPONSE",
      );
    }
  }

  if (!response.ok || !payload?.success) {
    throw new ApiError(
      payload?.error?.message ?? `Request failed with status ${response.status}.`,
      response.status,
      payload?.error?.code ?? "HTTP_ERROR",
      payload?.error?.details,
    );
  }

  return payload.data;
}

export const itemsApi = {
  list: (signal?: AbortSignal) => request<Item[]>("/api/items", { signal }),
  create: (input: ItemInput) =>
    request<Item>("/api/items", { method: "POST", body: input }),
  update: (id: number, input: ItemInput) =>
    request<Item>(`/api/items/${id}`, { method: "PUT", body: input }),
  remove: (id: number) =>
    request<null>(`/api/items/${id}`, { method: "DELETE" }),
};
