import { ItemsManager } from "@/components/items-manager";

export default function Home() {
  return (
    <main className="mx-auto w-full max-w-3xl px-4 py-12 sm:px-6">
      <header className="mb-8">
        <p className="text-sm font-semibold uppercase tracking-wide text-indigo-600">
          Fullstack Starter
        </p>
        <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
          Items
        </h1>
        <p className="mt-2 text-sm leading-6 text-slate-600">
          Contoh fitur CRUD untuk memastikan frontend, backend, dan database
          sudah terhubung.
        </p>
      </header>

      <ItemsManager />
    </main>
  );
}
