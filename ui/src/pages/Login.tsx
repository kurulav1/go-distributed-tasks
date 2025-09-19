import React, { useState } from "react"
import { Link, useLocation, useNavigate } from "react-router-dom"
import { api } from "../api"
import { useAuth } from "../auth"

export default function Login() {
  const { setToken } = useAuth()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [err, setErr] = useState<string | null>(null)
  const nav = useNavigate()
  const loc = useLocation()
  async function submit(e: React.FormEvent) {
    e.preventDefault()
    setErr(null)
    try {
      const r = await api("/login", { method: "POST", body: { email, password } })
      setToken(r.token)
      const to = (loc.state as any)?.from?.pathname || "/dashboard"
      nav(to, { replace: true })
    } catch (e: any) {
      setErr(e.message)
    }
  }
  return (
    <div className="auth">
      <h1>Login</h1>
      <form onSubmit={submit}>
        <input value={email} onChange={e=>setEmail(e.target.value)} placeholder="email" />
        <input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="password" />
        <button type="submit">Login</button>
      </form>
      {err && <p>{err}</p>}
      <p>New here? <Link to="/signup">Create an account</Link></p>
    </div>
  )
}
