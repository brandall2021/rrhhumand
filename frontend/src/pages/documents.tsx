import { useEffect, useState } from 'react'
import { Upload, Download, Trash2, FileText } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface Document {
  id: string
  name: string
  category: string
  file_size: number
  created_at: string
  employee_id: string | null
}

export default function DocumentsPage() {
  const [items, setItems] = useState<Document[]>([])
  const [loading, setLoading] = useState(true)

  const formatSize = (bytes: number) => {
    if (!bytes) return '-'
    const mb = bytes / (1024 * 1024)
    return mb >= 1 ? `${mb.toFixed(1)} MB` : `${(bytes / 1024).toFixed(0)} KB`
  }

  useEffect(() => {
    api.get('/documents', { params: { limit: '100' } })
      .then(res => setItems(res.data.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">Documentos</h1>
        <Button size="sm"><Upload size={16} className="mr-1" /> Subir</Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? <div className="p-6 text-center text-slate-500">Cargando...</div>
          : items.length === 0 ? <div className="p-6 text-center text-slate-500">No hay documentos subidos</div>
          : <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 bg-slate-50">
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Nombre</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Categoría</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Tamaño</th>
                    <th className="text-left px-4 py-3 font-medium text-slate-600">Subido</th>
                    <th className="text-right px-4 py-3 font-medium text-slate-600">Acciones</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map(d => (
                    <tr key={d.id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <FileText size={16} className="text-slate-400" />
                          <span className="font-medium text-slate-900">{d.name}</span>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-slate-600">{d.category}</td>
                      <td className="px-4 py-3 text-slate-600">{formatSize(d.file_size)}</td>
                      <td className="px-4 py-3 text-slate-600">{d.created_at?.slice(0, 10)}</td>
                      <td className="px-4 py-3 text-right">
                        <Button variant="ghost" size="sm"><Download size={14} /></Button>
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
