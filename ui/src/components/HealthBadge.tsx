import { useEffect, useState } from "react"
import { apiHealth } from "../api"
export default function HealthBadge() {
  const [status, setStatus] = useState<{ ok: boolean; text: string } | null>(null)
  useEffect(() => {
    let mounted = true
    const run = async () => {
      const res = await apiHealth()
      if (mounted) setStatus(res)
    }
    run()
    const id = setInterval(run, 5000)
    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [])
  if (!status) return <div className="badge-bad">checking</div>
  return status.ok ? <div className="badge-ok">api healthy</div> : <div className="badge-bad">api unreachable</div>
}
