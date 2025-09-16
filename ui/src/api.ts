const base = "/api"
export async function apiHealth(): Promise<{ ok: boolean; text: string }> {
  try {
    const r = await fetch(base + "/healthz")
    const t = await r.text()
    return { ok: r.ok, text: t }
  } catch (e) {
    return { ok: false, text: "unreachable" }
  }
}
export type EnqueueRequest = { type: string; payload: unknown }
export async function enqueueJob(token: string, body: EnqueueRequest): Promise<{ status: number; text: string }> {
  const r = await fetch(base + "/jobs", {
    method: "POST",
    headers: { "content-type": "application/json", "authorization": "Bearer " + token },
    body: JSON.stringify(body)
  })
  const t = await r.text()
  return { status: r.status, text: t }
}
export async function signup(email: string, password: string): Promise<{ ok: boolean; text: string }> {
  const r = await fetch(base + "/signup", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email, password })
  })
  const t = await r.text()
  return { ok: r.ok, text: t }
}
export async function login(email: string, password: string): Promise<{ ok: boolean; token?: string; text: string }> {
  const r = await fetch(base + "/login", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email, password })
  })
  const t = await r.text()
  try {
    const j = JSON.parse(t)
    return { ok: r.ok, token: j.token, text: t }
  } catch {
    return { ok: r.ok, text: t }
  }
}
export type Job = { id: string; user_id: string; type: string; payload: unknown; enqueued_at: string }
export async function listJobs(token: string, limit = 50): Promise<{ ok: boolean; jobs: Job[]; text: string }> {
  const r = await fetch(base + "/jobs?limit=" + limit, {
    headers: { "authorization": "Bearer " + token }
  })
  const t = await r.text()
  try {
    const j = JSON.parse(t)
    return { ok: r.ok, jobs: j as Job[], text: t }
  } catch {
    return { ok: r.ok, jobs: [], text: t }
  }
}
