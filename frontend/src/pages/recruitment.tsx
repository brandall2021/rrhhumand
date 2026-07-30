import { useEffect, useState } from 'react'
import api from '@/lib/api'

export default function RecruitmentPage() {
  const [data, setData] = useState<{ count: number } | null>(null)
  useEffect(() => {
    api.get('/recruitment/offers').then(r => setData(r.data)).catch(() => {})
  }, [])
  return (
    <div>
      <h1 className="text-2xl font-bold text-slate-900 mb-2">Reclutamiento</h1>
      <p className="text-slate-500">Registros encontrados: {data?.count ?? '...'}</p>
    </div>
  )
}
