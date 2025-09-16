import { useState } from "react"
import { login, signup } from "../api"
type Props = { onToken: (t: string) => void }
export default function AuthPanel({ onToken }: Props) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [busy, setBusy] = useState(false)
  const [msg, setMsg] = useState("")
  const doSignup = async () => {
    setBusy(true)
    const r = await signup(email, password)
    setMsg(r.text)
    setBusy(false)
  }
  const doLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    const r = await login(email, password)
    if (r.ok && r.token) {
      localStorage.setItem("token", r.token)
      onToken(r.token)
    }
    setMsg(r.text)
    setBusy(false)
  }
  return (
    <div className="card grid gap-3 max-w-xl">
      <div className="text-lg font-semibold">auth</div>
      <form className="grid gap-3" onSubmit={doLogin}>
        <input className="input" placeholder="email" value={email} onChange={e => setEmail(e.target.value)} />
        <input className="input" placeholder="password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
        <div className="flex gap-2">
          <button className="button" disabled={busy} type="submit">login</button>
          <button className="button" disabled={busy} type="button" onClick={doSignup}>signup</button>
        </div>
      </form>
      {msg ? <pre className="card whitespace-pre-wrap">{msg}</pre> : null}
    </div>
  )
}
