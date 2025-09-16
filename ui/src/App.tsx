import { useEffect, useState } from "react"
import HealthBadge from "./components/HealthBadge"
import JobForm from "./components/JobForm"
import JobLog from "./components/JobLog"
import AuthPanel from "./components/AuthPanel"
import JobsTable from "./components/JobsTable"
export default function App() {
  const [lines, setLines] = useState<string[]>([])
  const [token, setToken] = useState("")
  useEffect(() => {
    const t = localStorage.getItem("token") || ""
    setToken(t)
  }, [])
  const push = (s: string) => setLines(x => [...x, s])
  const logout = () => {
    localStorage.removeItem("token")
    setToken("")
  }
  return (
    <div className="min-h-screen p-6 md:p-10">
      <div className="max-w-6xl mx-auto grid gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold">distrib-jobs UI</h1>
          <div className="flex items-center gap-3">
            <HealthBadge />
            {token ? <button className="button" onClick={logout}>logout</button> : null}
          </div>
        </div>
        {!token ? (
          <AuthPanel onToken={setToken} />
        ) : (
          <div className="grid md:grid-cols-2 gap-6">
            <div className="card">
              <JobForm token={token} onResult={push} />
            </div>
            <div className="card">
              <div className="text-sm text-neutral-300 mb-2">log</div>
              <JobLog lines={lines} />
            </div>
            <div className="md:col-span-2">
              <JobsTable token={token} />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
