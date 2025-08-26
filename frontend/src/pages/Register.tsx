import { useState } from 'react'
import { register } from '../api'
import { Link, useNavigate } from 'react-router-dom'


export default function Register() {
    const [email, setEmail] = useState('')
    const [displayName, setDisplayName] = useState('')
    const [password, setPassword] = useState('')
    const [err, setErr] = useState('')
    const nav = useNavigate()


    async function onSubmit(e: React.FormEvent) {
        e.preventDefault()
        setErr('')
        try {
            await register(email, password, displayName)
            nav('/login')
        } catch (e: any) { setErr(e.message) }
    }


    return (
        <div className="min-h-screen flex items-center justify-center">
            <form onSubmit={onSubmit} className="p-8 rounded-2xl shadow w-80">
                <h1 className="text-xl font-bold mb-4">Register</h1>
                {err && <p className="text-red-600 mb-2">{err}</p>}
                <input className="border p-2 w-full mb-2" placeholder="Display name" value={displayName} onChange={e => setDisplayName(e.target.value)} />
                <input className="border p-2 w-full mb-2" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} />
                <input className="border p-2 w-full mb-4" placeholder="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
                <button className="w-full bg-black text-white rounded p-2">Register</button>
                <p className="mt-2 text-sm">Sudah punya akun? <Link className="underline" to="/login">Login</Link></p>
            </form>
        </div>
    )
}