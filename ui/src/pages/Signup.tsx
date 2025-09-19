import React, { useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import { api } from "../api"

export default function Signup() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [ok, setOk] = useState(false)
  const [err, setErr] = useState<string | null>(null)
  const nav = useNavigate()
  async function submit(e: React.FormEvent) {
    e.preventDefault()
    setErr(null)
    try {
      await api("/signup", { method: "POST", body: { email, password } })
      setOk(true)
      setTimeout(()=>nav("/login"), 600)
    } catch (e: any) {
      setErr(e.message)
    }
  }
  return (
    <div className="auth">
      <h1>Sign up</h1>
      <form onSubmit={submit}>
        <input value={email} onChange={e=>setEmail(e.target.value)} placeholder="email" />
        <input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="password" />
        <button type="submit">Create account</button>
      </form>
      {ok && <p>Account created</p>}
      {err && <p>{err}</p>}
      <p>Already have an account? <Link to="/login">Log in</Link></p>
    </div>
  )
}
