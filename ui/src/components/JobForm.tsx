import { useState } from "react"
import { enqueueJob } from "../api"
type Props = { onResult: (s: string) => void }
export default function JobForm({ onResult }: Props) {
  const [type, setType] = useState("email")
  const [payload, setPayload] = useState('{"to":"user@example.com"}')
  const [busy, setBusy] = useState(false)
  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    let parsed: unknown
    try {
      parsed = JSON.parse(payload)
    } catch (e) {
      onResult("invalid JSON")
      return
    }
    setBusy(true)
    try {
      const body = { type, payload: parsed }
      onResult("enqueue -> " + JSON.stringify(body))
      const res = await enqueueJob(body)
      onResult("response (" + res.status + "): " + res.text)
    } finally {
      setBusy(false)
    }
  }
  return (
    <form onSubmit={submit} className="grid gap-3 max-w-2xl">
      <div className="grid gap-2">
        <label className="text-sm text-neutral-300">type</label>
        <input className="input" value={type} onChange={e => setType(e.target.value)} />
      </div>
      <div className="grid gap-2">
        <label className="text-sm text-neutral-300">payload (JSON)</label>
        <textarea className="input h-40" value={payload} onChange={e => setPayload(e.target.value)} />
      </div>
      <div className="flex items-center gap-3">
        <button className="button" disabled={busy} type="submit">enqueue</button>
      </div>
    </form>
  )
}
