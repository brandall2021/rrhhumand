import { useEffect, useState } from 'react'
import { Plus, Pencil, Trash2, Receipt } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface Expense {
  id: string
  description: string
  amount: string
  status: string
  category: string
  created_at: string
  employee_id: string | null
}

export default function ExpensesPage() {
  const [items, setItems] = useState<Expense[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    api.get('/expenses', { params: { limit: '100' } })
      .then(res => {
        const data = res.data.data ?? res.data ?? []
        setItems(Array.isArray(data) ? data : [])
      })
      .catch(e => setError(e.response?.data?.error || 'Error al cargar gastos'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Gastos</h1>
        <div className="flex gap-2">
          <Button size="sm" variant="outline"><Receipt size={16} className="mr-1" /> Reportes</Button>
          <Button size="sm"><Plus size={16} className="mr-1" /> Nuevo gasto</Button>
        </div>
      </div>

      {error && <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg">{error}</div>}

      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 && !error ? <div className="p-6 text-center text-slate-500">No hay gastos registrados</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Descripción</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Monto</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Estado</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(e => (
                    <tr key={e.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3 font-medium text-slate-900">{e.description}</td>
                      <td className="px-4 py-3 text-slate-600">{e.category}</td>
                      <td className="px-4 py-3 text-slate-900 font-medium">${parseFloat(e.amount || '0').toLocaleString()}</td>
                      <td className="px-4 py-3">
                        <span className="inline-flex px-2 py-0.5 rounded-full text-xs font-medium bg-amber-50 text-amber-700">{e.status}</span>
                      </td>
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
