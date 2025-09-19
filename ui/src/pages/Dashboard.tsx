import React, { useEffect, useState } from "react"
import { api } from "../api"
import { useAuth } from "../auth"

type Job = { id: string; type: string; payload: any; enqueued_at: string }

export default function Dashboard() {
  const { token, setToken } = useAuth()
  const [type, setType] = useState("email")
  const [payload, setPayload] = useState('{"to":"user@example.com"}')
  const [jobs, setJobs] = useState<Job[]>([])
  const [err, setErr] = useState<string | null>(null)

  async function load() {
    try {
      const r = await api("/jobs", { token })
      const list = Array.isArray((r as any)?.jobs) ? (r as any).jobs : Array.isArray(r) ? (r as any) : []
      setJobs(list)
    } catch (e: any) {
      setErr(e.message)
      setJobs([])
    }
  }

  async function enqueue(e: React.FormEvent) {
    e.preventDefault()
    setErr(null)
    try {
      const body = { type, payload: JSON.parse(payload || "{}") }
      await api("/enqueue", { method: "POST", token, body })
      await load()
    } catch (e: any) {
      setErr(e.message)
    }
  }

  useEffect(() => { load() }, [])

  return (
    <div className="dash">
      <header>
        <h1>Jobs</h1>
        <button onClick={()=>{ setToken(null) }}>Logout</button>
      </header>
      <section>
        <form onSubmit={enqueue}>
          <select value={type} onChange={e=>setType(e.target.value)}>
            <option value="email">email</option>
            <option value="log">log</option>
          </select>
          <textarea value={payload} onChange={e=>setPayload(e.target.value)} rows={4} />
          <button type="submit">Enqueue</button>
        </form>
        {err && <p>{err}</p>}
      </section>
      <section>
        <table>
          <thead>
            <tr><th>ID</th><th>Type</th><th>Payload</th><th>Enqueued</th></tr>
          </thead>
          <tbody>
            {(jobs ?? []).map((j: Job) => (
              <tr key={j.id}>
                <td>{j.id}</td>
                <td>{j.type}</td>
                <td><pre>{JSON.stringify(j.payload)}</pre></td>
                <td>{j?.enqueued_at ? new Date(j.enqueued_at).toLocaleString() : ""}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </div>
  )
}
