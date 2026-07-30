import { useEffect, useState } from 'react'
import { Plus, Pencil, Trash2, Users, FileText, BarChart3 } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface Offer {
  id: string
  title: string
  status: string
  candidate_count?: number
  created_at?: string
}

export default function RecruitmentPage() {
  const [items, setItems] = useState<Offer[]>([])
  const [loading, setLoading] = useState(true)
  const [tab, setTab] = useState('offers')

  useEffect(() => {
    const endpoint = tab === 'offers' ? '/recruitment/offers' : '/recruitment/postings'
    api.get(endpoint)
      .then(res => setItems(res.data.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [tab])

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Reclutamiento</h1>
        <Button size="sm"><Plus size={16} className="mr-1" /> Nueva oferta</Button>
      </div>

      <div className="flex gap-2 mb-4">
        {[
          { key: 'offers', label: 'Ofertas', icon: FileText },
          { key: 'postings', label: 'Publicaciones', icon: BarChart3 },
          { key: 'candidates', label: 'Candidatos', icon: Users },
        ].map(t => (
          <button key={t.key} onClick={() => setTab(t.key)}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${tab === t.key ? 'bg-brand-50 text-brand-700' : 'text-slate-600 hover:bg-slate-100'}`}>
            <t.icon size={16} /> {t.label}
          </button>
        ))}
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 ? <div className="p-6 text-center text-slate-500">No hay registros</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Título</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Candidatos</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(o => (
                    <tr key={o.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 font-medium text-slate-900">{o.title}</td>
                      <td className="px-4 py-3">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700">{o.status}</span>
                      </td>
                      <td className="px-4 py-3 text-slate-600">{o.candidate_count ?? 0}</td>
                      <td className="px-4 py-3 text-right">
                        <Button variant="ghost" size="sm"><Pencil size={14} /></Button>
                        <Button variant="ghost" size="sm" className="text-red-500"><Trash2 size={14} /></Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>}
        </CardContent>
      </Card>
    </div>
  )
}
