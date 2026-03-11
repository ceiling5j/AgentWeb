import { useMemo, useState } from 'react'

const apiBaseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export function App() {
  const [input, setInput] = useState('')
  const [messages, setMessages] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const endpoint = useMemo(() => `${apiBaseURL.replace(/\/$/, '')}/api/chat`, [])

  const onSend = async (event) => {
    event.preventDefault()
    if (!input.trim() || loading) return

    const userMessage = { role: 'user', content: input.trim() }
    setMessages((prev) => [...prev, userMessage])
    setLoading(true)
    setError('')
    setInput('')

    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: userMessage.content })
      })

      if (!response.ok) {
        const details = await response.text()
        throw new Error(details || `HTTP ${response.status}`)
      }

      const data = await response.json()
      setMessages((prev) => [...prev, { role: 'assistant', content: data.reply || 'No response' }])
    } catch (err) {
      setError(`请求失败: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page">
      <section className="chat-card">
        <h1>Online Agent (DeepSeek)</h1>
        <p className="hint">前端 React + 后端 Golang（Vercel 可部署）</p>

        <div className="messages">
          {messages.length === 0 ? (
            <p className="empty">开始提问吧，例如："请帮我规划今天的学习任务"</p>
          ) : (
            messages.map((message, idx) => (
              <article key={idx} className={`message ${message.role}`}>
                <strong>{message.role === 'user' ? '你' : 'Agent'}:</strong>
                <span>{message.content}</span>
              </article>
            ))
          )}
        </div>

        {error && <p className="error">{error}</p>}

        <form onSubmit={onSend} className="composer">
          <textarea
            rows="3"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="输入你的问题..."
          />
          <button type="submit" disabled={loading}>{loading ? '发送中...' : '发送'}</button>
        </form>
      </section>
    </main>
  )
}
