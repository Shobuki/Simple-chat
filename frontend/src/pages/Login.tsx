import { useState } from 'react'
import { login } from '../api'
import { useNavigate, Link } from 'react-router-dom'


export default function Login(){
const [email, setEmail] = useState('')
const [password, setPassword] = useState('')
const [err, setErr] = useState('')
const nav = useNavigate()


async function onSubmit(e: React.FormEvent){
e.preventDefault()
setErr('')
try {
const { token } = await login(email, password)
localStorage.setItem('token', token)
nav('/chat')
} catch (e:any) { setErr(e.message) }
}


return (
<div className="min-h-screen flex items-center justify-center">
<form onSubmit={onSubmit} className="p-8 rounded-2xl shadow w-80">
<h1 className="text-xl font-bold mb-4">Login</h1>
{err && <p className="text-red-600 mb-2">{err}</p>}
<input className="border p-2 w-full mb-2" placeholder="Email" value={email} onChange={e=>setEmail(e.target.value)} />
<input className="border p-2 w-full mb-4" placeholder="Password" type="password" value={password} onChange={e=>setPassword(e.target.value)} />
<button className="w-full bg-black text-white rounded p-2">Login</button>
<p className="mt-2 text-sm">Belum punya akun? <Link className="underline" to="/register">Register</Link></p>
</form>
</div>
)
}