import { useState, useEffect, useCallback } from 'react'
import './App.css'

const APIS = [
  { id: 'go',   label: 'Go',   base: '/go-api',   color: '#00ADD8' },
  { id: 'rust', label: 'Rust', base: '/rust-api',  color: '#CE422B' },
]

function ApiCard({ api }) {
  const [count, setCount]       = useState(null)
  const [timestamp, setTimestamp] = useState(null)
  const [calling, setCalling]   = useState(false)
  const [error, setError]       = useState(null)

  const fetchCount = useCallback(async () => {
    try {
      const res = await fetch(`${api.base}/requests`)
      if (!res.ok) throw new Error(res.statusText)
      const data = await res.json()
      setCount(data.count)
      setTimestamp(data.timestamp)
      setError(null)
    } catch (e) {
      setError('Failed to reach API')
    }
  }, [api.base])

  useEffect(() => {
    fetchCount()
  }, [fetchCount])

  const handleCall = async () => {
    setCalling(true)
    try {
      const res = await fetch(`${api.base}/requests`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: api.id }),
      })
      if (!res.ok) throw new Error(res.statusText)
      await fetchCount()
    } catch (e) {
      setError('Call failed')
    } finally {
      setCalling(false)
    }
  }

  return (
    <div className="card" style={{ '--accent': api.color }}>
      <div className="card-header">
        <span className="lang-badge" style={{ background: api.color }}>{api.label}</span>
        <span className="card-title">API</span>
      </div>

      <div className="count-block">
        <span className="count">{count ?? '—'}</span>
        <span className="count-label">calls</span>
      </div>

      {timestamp && (
        <p className="timestamp">last at {new Date(timestamp).toLocaleTimeString()}</p>
      )}

      {error && <p className="error">{error}</p>}

      <button
        className="call-btn"
        onClick={handleCall}
        disabled={calling}
        style={{ background: api.color }}
      >
        {calling ? 'Calling…' : `Call ${api.label} API`}
      </button>
    </div>
  )
}

export default function App() {
  return (
    <div className="app">
      <h1 className="title">API Dashboard</h1>
      <div className="cards">
        {APIS.map((api) => (
          <ApiCard key={api.id} api={api} />
        ))}
      </div>
    </div>
  )
}
