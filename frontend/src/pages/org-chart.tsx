import { useEffect, useState } from 'react'
import { ChevronDown, ChevronRight, User } from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'

interface OrgNode {
  id: string
  first_name: string
  last_name: string
  position_name?: string
  department_name?: string
  email?: string
  photo_url?: string
  children: OrgNode[]
}

function TreeNode({ node, depth }: { node: OrgNode; depth: number }) {
  const [open, setOpen] = useState(depth < 1)
  const hasChildren = node.children && node.children.length > 0
  const initials = `${node.first_name} ${node.last_name}`.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)

  return (
    <div>
      <div
        className={`flex items-center gap-3 py-2 px-3 rounded-lg hover:bg-slate-50 transition-colors ${
          hasChildren ? 'cursor-pointer' : 'cursor-default'
        }`}
        style={{ marginLeft: depth * 28 }}
        onClick={() => hasChildren && setOpen(!open)}
      >
        {hasChildren ? (
          <button className="text-slate-400 shrink-0" onClick={() => setOpen(!open)}>
            {open ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
          </button>
        ) : (
          <span className="w-4 shrink-0" />
        )}
        <Avatar className="h-8 w-8">
          {node.photo_url && <AvatarImage src={node.photo_url} alt="" />}
          <AvatarFallback className="bg-brand-100 text-brand-700 text-xs">{initials}</AvatarFallback>
        </Avatar>
        <div className="min-w-0">
          <p className="text-sm font-medium text-slate-900 truncate">{node.first_name} {node.last_name}</p>
          <p className="text-xs text-slate-500 truncate">
            {[node.position_name, node.department_name].filter(Boolean).join(' · ') || 'Sin asignación'}
          </p>
        </div>
      </div>
      {open && hasChildren && (
        <div className="relative ml-[22px] border-l border-slate-200 pl-3">
          {node.children.map((child) => (
            <TreeNode key={child.id} node={child} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  )
}

export default function OrgChartPage() {
  const [tree, setTree] = useState<OrgNode[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    api.get('/organization/tree')
      .then((res) => setTree(res.data.data ?? []))
      .catch((err: any) => setError(err?.response?.data?.error || 'Error al cargar el organigrama'))
      .finally(() => setLoading(false))
  }, [])

  const count = (nodes: OrgNode[]): number => nodes.reduce((acc, n) => acc + 1 + count(n.children ?? []), 0)

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Organigrama</h1>
          <p className="text-sm text-slate-500 mt-1">
            {loading ? 'Cargando...' : `${count(tree)} empleados activos`}
          </p>
        </div>
      </div>

      {error && (
        <div className="mb-4 p-4 rounded-lg bg-red-50 text-red-700 text-sm">{error}</div>
      )}

      {loading ? (
        <div className="text-center text-slate-500 py-12">Cargando...</div>
      ) : tree.length === 0 ? (
        <Card>
          <CardContent className="p-10 text-center">
            <User className="mx-auto h-10 w-10 text-slate-300 mb-3" />
            <p className="text-slate-500">No hay empleados activos para mostrar el organigrama</p>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent className="p-4">
            {tree.map((node) => (
              <TreeNode key={node.id} node={node} depth={0} />
            ))}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
