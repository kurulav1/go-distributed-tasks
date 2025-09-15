import { useState } from "react"
import HealthBadge from "./components/HealthBadge"
import JobForm from "./components/JobForm"
import JobLog from "./components/JobLog"
export default function App() {
  const [lines, setLines] = useState<string[]>([])
  const push = (s: string) => setLines(x => [...x, s])
  return (
    <div className="min-h-screen p-6 md:p-10">
      <div className="max-w-6xl mx-auto grid gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold">distrib-jobs UI</h1>
          <HealthBadge />
        </div>
        <div className="grid md:grid-cols-2 gap-6">
          <div className="card">
            <JobForm onResult={push} />
          </div>
          <div className="card">
            <div className="text-sm text-neutral-300 mb-2">log</div>
            <JobLog lines={lines} />
          </div>
        </div>
      </div>
    </div>
  )
}
