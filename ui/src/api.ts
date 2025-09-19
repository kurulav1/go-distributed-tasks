type FetchOpts = { method?: string; body?: any; token?: string | null }
const base = (window as any).__CONFIG__?.API_BASE || import.meta.env.VITE_API_BASE || "http://localhost:8080"

export async function api(path: string, opts: FetchOpts = {}) {
  const headers: Record<string, string> = { "Content-Type": "application/json" }
  if (opts.token) headers.Authorization = `Bearer ${opts.token}`
  const res = await fetch(`${base}${path}`, {
    method: opts.method || "GET",
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined
  })
  const text = await res.text()
  const data = text ? JSON.parse(text) : {}
  if (!res.ok) throw new Error((data as any)?.error || res.statusText)
  return data ?? {}
}
