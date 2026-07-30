import { useEffect, useState } from 'react'
import api from '@/lib/api'

export default function ExpensesPage() {
  const [data, setData] = useState<{ count: number } | null>(null)
  useEffect(() => {
    api.get('/expenses').then(r => setData(r.data)).catch(() => {})
  }, [])
  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-2">Gastos</h1>
      <p className="text-slate-500">Registros encontrados: {data?.count ?? '...'}</p>
    </div>
  )
}
