import { useEffect, useMemo, useRef, useState } from 'react'
import { getMessages, wsUrl } from '../api'
import { useNavigate } from 'react-router-dom'

type Message = { id: number; userId: number; displayName: string; content: string; createdAt: string }

export default function Chat() {
  const nav = useNavigate()
  const token = useMemo(() => localStorage.getItem('token'), [])
  const [messages, setMessages] = useState<Message[]>([])
  const [text, setText] = useState('')
  const [loading, setLoading] = useState(true)
  const [err, setErr] = useState<string>('')

  const wsRef = useRef<WebSocket | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!token) nav('/login', { replace: true })
  }, [token, nav])

  useEffect(() => {
    if (!token) return
    let cancelled = false
    setLoading(true)
    setErr('')

    getMessages()
      .then((data) => {
        if (cancelled) return
        setMessages(Array.isArray(data) ? (data as Message[]) : [])
      })
      .catch(() => {
        if (cancelled) return
        setErr('Gagal memuat pesan awal.')
        setMessages([])
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [token])

  useEffect(() => {
    if (!token) return
    const ws = new WebSocket(wsUrl(token))
    wsRef.current = ws

    ws.onopen = () => setErr('')
    ws.onmessage = (e) => {
      try {
        const msg: Message = JSON.parse(e.data)
        setMessages((prev) => (Array.isArray(prev) ? [...prev, msg] : [msg]))
      } catch {}
    }
    ws.onerror = () => setErr('Koneksi realtime error. Mencoba kembali...')
    ws.onclose = () => {}
    return () => ws.close()
  }, [token])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  function send() {
    if (!text.trim()) return
    const ws = wsRef.current
    if (!ws || ws.readyState !== WebSocket.OPEN) {
      setErr('Tidak terhubung ke server chat.')
      return
    }
    ws.send(JSON.stringify({ content: text.trim() }))
    setText('')
  }

  return (
    <div className="h-screen flex flex-col bg-blue-50">
      {/* Header */}
      <header className="p-4 bg-blue-600 text-white font-semibold flex items-center justify-between shadow">
        <span>💬 Realtime Chat</span>
        <span className="text-sm opacity-80">
          {wsRef.current?.readyState === WebSocket.OPEN ? '🟢 Connected' : '🔴 Disconnected'}
        </span>
      </header>

      {/* Chat Area */}
      <main className="flex-1 overflow-y-auto p-4 space-y-3">
        {err && <div className="text-red-600 border border-red-200 bg-red-50 rounded p-2">{err}</div>}
        {loading && messages.length === 0 && <div className="opacity-70">Memuat pesan…</div>}

        {(messages ?? []).map((m) => {
          const isMine = m.displayName === localStorage.getItem('displayName')
          return (
            <div
              key={m.id}
              className={`max-w-xs md:max-w-md p-3 rounded-2xl shadow ${
                isMine
                  ? 'ml-auto bg-blue-500 text-white rounded-br-none'
                  : 'mr-auto bg-gray-200 text-gray-900 rounded-bl-none'
              }`}
            >
              <div className="text-xs mb-1 opacity-70">
                {m.displayName} • {new Date(m.createdAt).toLocaleTimeString()}
              </div>
              <div>{m.content}</div>
            </div>
          )
        })}
        <div ref={bottomRef} />
      </main>

      {/* Input box */}
      <div className="p-4 border-t bg-white flex gap-2">
        <input
          className="flex-1 border-2 border-blue-300 rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
          placeholder="Tulis pesan..."
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') send()
          }}
        />
        <button
          className="px-5 py-2 rounded-full bg-blue-600 text-white font-medium hover:bg-blue-700 transition"
          onClick={send}
        >
          ➤
        </button>
      </div>
    </div>
  )
}
