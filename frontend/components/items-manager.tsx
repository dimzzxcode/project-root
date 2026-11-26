"use client";

import { useCallback, useEffect, useState } from "react";
import type { FormEvent } from "react";
import { ApiError, itemsApi } from "@/lib/api";
import type { Item, ItemInput } from "@/lib/api";

type FormState = ItemInput;

const emptyForm: FormState = { name: "", description: "" };

const inputClass =
  "w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm outline-none transition placeholder:text-slate-400 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-100";

const primaryButton =
  "inline-flex items-center justify-center rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:cursor-not-allowed disabled:opacity-60";

const secondaryButton =
  "inline-flex items-center justify-center rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60";

const dangerButton =
  "inline-flex items-center justify-center rounded-lg border border-red-200 bg-white px-3 py-1.5 text-sm font-medium text-red-600 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-60";

const dateFormatter = new Intl.DateTimeFormat("en", {
  dateStyle: "medium",
  timeStyle: "short",
});

export function ItemsManager() {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  const [editingId, setEditingId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<FormState>(emptyForm);
  const [editError, setEditError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  const fetchItems = useCallback(async (signal?: AbortSignal) => {
    try {
      const data = await itemsApi.list(signal);
      setItems(data);
      setLoadError(null);
    } catch (error) {
      if (isAbortError(error)) {
        return;
      }
      setLoadError(errorMessage(error));
    } finally {
      if (!signal?.aborted) {
        setLoading(false);
      }
    }
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    // Initial data load. State updates happen after the request resolves.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    void fetchItems(controller.signal);
    return () => controller.abort();
  }, [fetchItems]);

  function reload() {
    setLoading(true);
    setLoadError(null);
    void fetchItems();
  }

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormError(null);

    if (!form.name.trim()) {
      setFormError("Name is required.");
      return;
    }

    setCreating(true);
    try {
      const created = await itemsApi.create({
        name: form.name.trim(),
        description: form.description.trim(),
      });
      setItems((current) => [created, ...current]);
      setForm(emptyForm);
    } catch (error) {
      setFormError(errorMessage(error));
    } finally {
      setCreating(false);
    }
  }

  function startEdit(item: Item) {
    setEditingId(item.id);
    setEditForm({ name: item.name, description: item.description });
    setEditError(null);
    setActionError(null);
  }

  function cancelEdit() {
    setEditingId(null);
    setEditError(null);
  }

  async function handleUpdate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (editingId === null) {
      return;
    }

    setEditError(null);

    if (!editForm.name.trim()) {
      setEditError("Name is required.");
      return;
    }

    setSaving(true);
    try {
      const updated = await itemsApi.update(editingId, {
        name: editForm.name.trim(),
        description: editForm.description.trim(),
      });
      setItems((current) =>
        current.map((item) => (item.id === updated.id ? updated : item)),
      );
      setEditingId(null);
    } catch (error) {
      setEditError(errorMessage(error));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(item: Item) {
    if (!window.confirm(`Delete "${item.name}"?`)) {
      return;
    }

    setActionError(null);
    setDeletingId(item.id);

    try {
      await itemsApi.remove(item.id);
      setItems((current) => current.filter((entry) => entry.id !== item.id));
      if (editingId === item.id) {
        setEditingId(null);
      }
    } catch (error) {
      setActionError(errorMessage(error));
    } finally {
      setDeletingId(null);
    }
  }

  return (
    <section className="space-y-8">
      <form
        onSubmit={handleCreate}
        className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm"
      >
        <h2 className="text-base font-semibold text-slate-900">Add item</h2>
        <p className="mt-1 text-sm text-slate-500">
          Create a record to verify the API and database connection.
        </p>

        <div className="mt-4">
          <ItemFields
            idPrefix="create"
            values={form}
            onChange={setForm}
            disabled={creating}
          />
        </div>

        {formError ? (
          <p role="alert" className="mt-3 text-sm text-red-600">
            {formError}
          </p>
        ) : null}

        <div className="mt-4 flex justify-end">
          <button type="submit" disabled={creating} className={primaryButton}>
            {creating ? "Adding…" : "Add item"}
          </button>
        </div>
      </form>

      <div>
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-base font-semibold text-slate-900">All items</h2>
          {!loading && !loadError ? (
            <span className="text-sm text-slate-500">
              {items.length} {items.length === 1 ? "item" : "items"}
            </span>
          ) : null}
        </div>

        {actionError ? (
          <div className="mb-3">
            <ErrorBanner message={actionError} />
          </div>
        ) : null}

        {loading ? (
          <ItemListSkeleton />
        ) : loadError ? (
          <ErrorBanner message={loadError} onRetry={reload} />
        ) : items.length === 0 ? (
          <div className="rounded-xl border border-dashed border-slate-300 bg-white p-8 text-center">
            <p className="text-sm font-medium text-slate-700">No items yet</p>
            <p className="mt-1 text-sm text-slate-500">
              Add your first item using the form above.
            </p>
          </div>
        ) : (
          <ul className="space-y-3">
            {items.map((item) => (
              <li key={item.id}>
                {editingId === item.id ? (
                  <form
                    onSubmit={handleUpdate}
                    className="rounded-xl border border-indigo-200 bg-white p-5 shadow-sm"
                  >
                    <ItemFields
                      idPrefix={`edit-${item.id}`}
                      values={editForm}
                      onChange={setEditForm}
                      disabled={saving}
                    />

                    {editError ? (
                      <p role="alert" className="mt-3 text-sm text-red-600">
                        {editError}
                      </p>
                    ) : null}

                    <div className="mt-4 flex justify-end gap-2">
                      <button
                        type="button"
                        onClick={cancelEdit}
                        disabled={saving}
                        className={secondaryButton}
                      >
                        Cancel
                      </button>
                      <button
                        type="submit"
                        disabled={saving}
                        className={primaryButton}
                      >
                        {saving ? "Saving…" : "Save changes"}
                      </button>
                    </div>
                  </form>
                ) : (
                  <article className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div className="min-w-0">
                        <h3 className="font-semibold break-words text-slate-900">
                          {item.name}
                        </h3>
                        <p className="mt-1 text-sm break-words text-slate-600">
                          {item.description || "No description."}
                        </p>
                      </div>

                      <div className="flex shrink-0 gap-2">
                        <button
                          type="button"
                          onClick={() => startEdit(item)}
                          className={secondaryButton}
                        >
                          Edit
                        </button>
                        <button
                          type="button"
                          onClick={() => void handleDelete(item)}
                          disabled={deletingId === item.id}
                          className={dangerButton}
                        >
                          {deletingId === item.id ? "Deleting…" : "Delete"}
                        </button>
                      </div>
                    </div>

                    <p className="mt-3 text-xs text-slate-400">
                      Created {formatDate(item.created_at)} · Updated{" "}
                      {formatDate(item.updated_at)}
                    </p>
                  </article>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}

type ItemFieldsProps = {
  idPrefix: string;
  values: FormState;
  onChange: (values: FormState) => void;
  disabled?: boolean;
};

function ItemFields({
  idPrefix,
  values,
  onChange,
  disabled = false,
}: ItemFieldsProps) {
  return (
    <div className="space-y-3">
      <div>
        <label
          htmlFor={`${idPrefix}-name`}
          className="mb-1 block text-sm font-medium text-slate-700"
        >
          Name
        </label>
        <input
          id={`${idPrefix}-name`}
          type="text"
          value={values.name}
          onChange={(event) =>
            onChange({ ...values, name: event.target.value })
          }
          maxLength={255}
          required
          disabled={disabled}
          placeholder="e.g. First item"
          className={inputClass}
        />
      </div>

      <div>
        <label
          htmlFor={`${idPrefix}-description`}
          className="mb-1 block text-sm font-medium text-slate-700"
        >
          Description{" "}
          <span className="font-normal text-slate-400">(optional)</span>
        </label>
        <textarea
          id={`${idPrefix}-description`}
          value={values.description}
          onChange={(event) =>
            onChange({ ...values, description: event.target.value })
          }
          maxLength={2000}
          rows={2}
          disabled={disabled}
          placeholder="Short description"
          className={inputClass}
        />
      </div>
    </div>
  );
}

function ErrorBanner({
  message,
  onRetry,
}: {
  message: string;
  onRetry?: () => void;
}) {
  return (
    <div
      role="alert"
      className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
    >
      <span>{message}</span>
      {onRetry ? (
        <button
          type="button"
          onClick={onRetry}
          className="font-semibold underline underline-offset-2"
        >
          Retry
        </button>
      ) : null}
    </div>
  );
}

function ItemListSkeleton() {
  return (
    <div className="space-y-3" role="status" aria-label="Loading items">
      {[0, 1, 2].map((index) => (
        <div
          key={index}
          className="animate-pulse rounded-xl border border-slate-200 bg-white p-5"
        >
          <div className="h-4 w-1/3 rounded bg-slate-200" />
          <div className="mt-3 h-3 w-2/3 rounded bg-slate-100" />
          <div className="mt-4 h-3 w-1/4 rounded bg-slate-100" />
        </div>
      ))}
    </div>
  );
}

function isAbortError(error: unknown): boolean {
  return error instanceof DOMException && error.name === "AbortError";
}

function errorMessage(error: unknown): string {
  if (error instanceof ApiError) {
    const firstDetail = error.details
      ? Object.values(error.details)[0]
      : undefined;
    return firstDetail ?? error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return "Something went wrong. Please try again.";
}

function formatDate(value: string): string {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : dateFormatter.format(date);
}
