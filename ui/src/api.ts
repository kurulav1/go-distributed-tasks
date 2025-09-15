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
export async function enqueueJob(body: EnqueueRequest): Promise<{ status: number; text: string }> {
  const r = await fetch(base + "/jobs", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(body)
  })
  const t = await r.text()
  return { status: r.status, text: t }
}
