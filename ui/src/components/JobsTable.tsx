import { useEffect, useState } from "react"
import { Job, listJobs } from "../api"
type Props = { token: string }
export default function JobsTable({ token }: Props) {
  const [items, setItems] = useState<Job[]>([])
  const [err, setErr] = useState("")
  useEffect(() => {
    let mounted = true
    const run = async () => {
      if (!token) return
      const r = await listJobs(token, 50)
      if (!mounted) return
      if (r.ok) {
        setItems(r.jobs)
        setErr("")
      } else {
        setErr(r.text)
      }
    }
    run()
    const id = setInterval(run, 5000)
    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [token])
  return (
    <div className="card">
      <div className="text-sm text-neutral-300 mb-2">recent jobs</div>
      {err ? <pre className="card whitespace-pre-wrap mb-3">{err}</pre> : null}
      <div className="overflow-auto">
        <table className="w-full text-sm">
          <thead className="text-neutral-300">
            <tr>
              <th className="text-left p-2">id</th>
              <th className="text-left p-2">type</th>
              <th className="text-left p-2">enqueued_at</th>
              <th className="text-left p-2">payload</th>
            </tr>
          </thead>
          <tbody>
            {items.map(j => (
              <tr key={j.id} className="border-t border-neutral-800">
                <td className="p-2">{j.id}</td>
                <td className="p-2">{j.type}</td>
                <td className="p-2">{new Date(j.enqueued_at).toLocaleString()}</td>
                <td className="p-2"><pre className="whitespace-pre-wrap">{JSON.stringify(j.payload)}</pre></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
