import { Link } from 'react-router-dom'


export default function App() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="p-8 rounded-2xl shadow">
        <h1 className="text-2xl font-bold mb-4">Go + React Chat</h1>
        <div className="flex gap-4">
          <Link className="px-4 py-2 rounded bg-black text-white" to="/login">Login</Link>
          <Link className="px-4 py-2 rounded border" to="/register">Register</Link>
        </div>
      </div>
    </div>
  )
}