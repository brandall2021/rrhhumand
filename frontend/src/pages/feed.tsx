import { useEffect, useState } from 'react'
import {
  Send,
  ThumbsUp,
  Heart,
  Smile,
  PartyPopper,
  MessageCircle,
  Pin,
  Pencil,
  Trash2,
  X,
  Check,
} from 'lucide-react'
import api from '@/lib/api'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Select } from '@/components/ui/select'

interface Post {
  id: string
  author_id: string
  author_name?: string
  author_photo?: string
  content: string
  visibility: string
  pinned: boolean
  reaction_counts?: Record<string, number>
  comment_count: number
  created_at: string
}

interface Comment {
  id: string
  author_id: string
  author_name?: string
  author_photo?: string
  content: string
  created_at: string
}

const reactions = [
  { type: 'LIKE', icon: ThumbsUp, label: 'Me gusta', active: 'bg-blue-50 text-blue-600' },
  { type: 'LOVE', icon: Heart, label: 'Me encanta', active: 'bg-red-50 text-red-500' },
  { type: 'CELEBRATE', icon: PartyPopper, label: 'Celebrar', active: 'bg-amber-50 text-amber-600' },
  { type: 'FUN', icon: Smile, label: 'Me divierte', active: 'bg-emerald-50 text-emerald-600' },
]

const visibilityLabels: Record<string, string> = {
  company: 'Toda la empresa',
  department: 'Departamento',
  managers: 'Jefes',
  all: 'Toda la empresa',
}

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'ahora'
  if (mins < 60) return `hace ${mins} min`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `hace ${hours} h`
  const days = Math.floor(hours / 24)
  if (days < 30) return `hace ${days} d`
  return new Date(dateStr).toLocaleDateString('es-AR')
}

export default function FeedPage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [myEmployeeId, setMyEmployeeId] = useState('')
  const [content, setContent] = useState('')
  const [visibility, setVisibility] = useState('company')
  const [pinned, setPinned] = useState(false)
  const [posting, setPosting] = useState(false)
  const [expandedComments, setExpandedComments] = useState<Record<string, boolean>>({})
  const [comments, setComments] = useState<Record<string, Comment[]>>({})
  const [commentText, setCommentText] = useState<Record<string, string>>({})
  const [commentLoading, setCommentLoading] = useState<Record<string, boolean>>({})
  const [reactionLoading, setReactionLoading] = useState<Record<string, boolean>>({})
  const [editingPost, setEditingPost] = useState<Post | null>(null)
  const [editContent, setEditContent] = useState('')
  const [myReactions, setMyReactions] = useState<Record<string, string>>({})

  const fetchPosts = async () => {
    setLoading(true)
    try {
      const res = await api.get('/feed', { params: { limit: '100' } })
      setPosts(res.data.data ?? [])
    } catch { setPosts([]) }
    finally { setLoading(false) }
  }

  useEffect(() => {
    api.get('/me/profile')
      .then((res) => setMyEmployeeId(res.data.data?.id ?? ''))
      .catch(() => {})
    fetchPosts()
  }, [])

  const handleCreate = async () => {
    if (!content.trim()) return
    setPosting(true)
    try {
      const res = await api.post('/feed', { content, visibility, pinned })
      if (res.data?.data) {
        setPosts((prev) => [res.data.data, ...prev])
      } else {
        fetchPosts()
      }
      setContent('')
      setPinned(false)
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al publicar')
    } finally { setPosting(false) }
  }

  const handleUpdate = async () => {
    if (!editingPost) return
    try {
      await api.put(`/feed/${editingPost.id}`, { content: editContent })
      setEditingPost(null)
      fetchPosts()
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al editar')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('¿Eliminar esta publicación?')) return
    try {
      await api.delete(`/feed/${id}`)
      setPosts((prev) => prev.filter((p) => p.id !== id))
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al eliminar')
    }
  }

  const toggleReaction = async (post: Post, type: string) => {
    const key = `${post.id}:${type}`
    setReactionLoading((prev) => ({ ...prev, [key]: true }))
    const alreadyReacted = myReactions[post.id] === type
    const counts = { ...(post.reaction_counts ?? {}) }
    counts[type] = Math.max(0, (counts[type] ?? 0) + (alreadyReacted ? -1 : 1))

    setPosts((prev) => prev.map((p) => (p.id === post.id ? { ...p, reaction_counts: counts } : p)))
    setMyReactions((prev) => ({ ...prev, [post.id]: alreadyReacted ? '' : type }))

    try {
      if (alreadyReacted) {
        await api.delete(`/feed/${post.id}/reactions/${type}`)
      } else {
        await api.post(`/feed/${post.id}/reactions`, { reaction_type: type })
      }
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al reaccionar')
      fetchPosts()
    } finally {
      setReactionLoading((prev) => {
        const next = { ...prev }
        delete next[key]
        return next
      })
    }
  }

  const toggleComments = async (post: Post) => {
    const expanded = expandedComments[post.id]
    setExpandedComments((prev) => ({ ...prev, [post.id]: !expanded }))
    if (!expanded && !comments[post.id]) {
      try {
        const res = await api.get(`/feed/${post.id}`)
        setComments((prev) => ({ ...prev, [post.id]: res.data.data?.comments ?? [] }))
      } catch { setComments((prev) => ({ ...prev, [post.id]: [] })) }
    }
  }

  const addComment = async (post: Post) => {
    const text = commentText[post.id]?.trim()
    if (!text) return
    setCommentLoading((prev) => ({ ...prev, [post.id]: true }))
    try {
      await api.post(`/feed/${post.id}/comments`, { content: text })
      const res = await api.get(`/feed/${post.id}`)
      setComments((prev) => ({ ...prev, [post.id]: res.data.data?.comments ?? [] }))
      setCommentText((prev) => ({ ...prev, [post.id]: '' }))
      setPosts((prev) =>
        prev.map((p) => (p.id === post.id ? { ...p, comment_count: (p.comment_count ?? 0) + 1 } : p)),
      )
    } catch (err: any) {
      alert(err?.response?.data?.error || 'Error al comentar')
    } finally {
      setCommentLoading((prev) => ({ ...prev, [post.id]: false }))
    }
  }

  const canModify = (post: Post) => myEmployeeId && post.author_id === myEmployeeId

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Feed</h1>

      <Card className="mb-6">
        <CardContent className="p-4">
          <textarea
            className="flex min-h-[100px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm mb-3"
            placeholder="¿Qué estás pensando? Mencioná a un compañero con @legajo"
            value={content}
            onChange={e => setContent(e.target.value)}
          />
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2">
              <Select
                className="w-44"
                options={Object.entries(visibilityLabels).map(([value, label]) => ({ value, label }))}
                value={visibility}
                onChange={e => setVisibility(e.target.value)}
              />
              <button
                onClick={() => setPinned(!pinned)}
                className={`flex items-center gap-1.5 text-xs px-2.5 py-1.5 rounded-lg border transition-colors ${
                  pinned ? 'bg-amber-50 text-amber-700 border-amber-200' : 'text-slate-500 border-slate-200 hover:bg-slate-50'
                }`}
              >
                <Pin size={14} /> Fijar
              </button>
            </div>
            <Button size="sm" onClick={handleCreate} disabled={posting || !content.trim()}>
              <Send size={14} /> {posting ? 'Publicando...' : 'Publicar'}
            </Button>
          </div>
        </CardContent>
      </Card>

      {loading ? (
        <div className="text-center text-slate-500 py-12">Cargando...</div>
      ) : posts.length === 0 ? (
        <Card>
          <CardContent className="p-10 text-center text-slate-500">
            Aún no hay publicaciones. ¡Sé el primero en escribir algo!
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => {
            const initials = (post.author_name ?? 'U').split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
            const counts = post.reaction_counts ?? {}
            const totalReactions = Object.values(counts).reduce((a, b) => a + b, 0)
            const isEditing = editingPost?.id === post.id

            return (
              <Card key={post.id}>
                <CardContent className="p-4">
                  <div className="flex items-start gap-3">
                    <Avatar className="h-10 w-10">
                      {post.author_photo && <AvatarImage src={post.author_photo} alt="" />}
                      <AvatarFallback className="bg-brand-100 text-brand-700">{initials}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-medium text-slate-900">{post.author_name || 'Empleado'}</span>
                        {post.pinned && (
                          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-amber-50 text-amber-700 text-xs font-medium">
                            <Pin size={11} /> Fijado
                          </span>
                        )}
                        <span className="text-xs text-slate-400">{timeAgo(post.created_at)}</span>
                      </div>
                      <p className="text-sm text-slate-800 mt-2 whitespace-pre-wrap">{post.content}</p>
                    </div>
                    {canModify(post) && (
                      <div className="flex gap-1 shrink-0">
                        <Button
                          variant="ghost" size="sm" className="text-slate-400"
                          onClick={() => { setEditingPost(post); setEditContent(post.content) }}
                        >
                          <Pencil size={14} />
                        </Button>
                        <Button variant="ghost" size="sm" className="text-red-400" onClick={() => handleDelete(post.id)}>
                          <Trash2 size={14} />
                        </Button>
                      </div>
                    )}
                  </div>

                  {isEditing && (
                    <div className="mt-3 p-3 rounded-lg bg-slate-50">
                      <textarea
                        className="flex min-h-[80px] w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                        value={editContent}
                        onChange={e => setEditContent(e.target.value)}
                      />
                      <div className="flex justify-end gap-2 mt-2">
                        <Button variant="ghost" size="sm" onClick={() => setEditingPost(null)}>
                          <X size={14} /> Cancelar
                        </Button>
                        <Button size="sm" onClick={handleUpdate}>
                          <Check size={14} /> Guardar
                        </Button>
                      </div>
                    </div>
                  )}

                  <div className="flex items-center gap-1 mt-3 pt-3 border-t border-slate-100">
                    {reactions.map((r) => {
                      const Icon = r.icon
                      const active = myReactions[post.id] === r.type
                      const count = counts[r.type] ?? 0
                      return (
                        <button
                          key={r.type}
                          onClick={() => toggleReaction(post, r.type)}
                          disabled={!!reactionLoading[`${post.id}:${r.type}`]}
                          className={`flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm transition-colors disabled:opacity-50 ${
                            active ? r.active : 'text-slate-500 hover:bg-slate-100'
                          }`}
                        >
                          <Icon size={16} />
                          {count > 0 && <span className="font-medium">{count}</span>}
                        </button>
                      )
                    })}
                    <span className="text-xs text-slate-400 ml-auto">
                      {totalReactions > 0 && `${totalReactions} reacciones`}
                      {totalReactions > 0 && (post.comment_count > 0) && ' · '}
                      {post.comment_count > 0 && `${post.comment_count} comentarios`}
                    </span>
                  </div>

                  <button
                    onClick={() => toggleComments(post)}
                    className="flex items-center gap-1.5 mt-1 text-sm text-slate-500 hover:text-slate-700"
                  >
                    <MessageCircle size={15} />
                    {expandedComments[post.id] ? 'Ocultar comentarios' : 'Comentar'}
                  </button>

                  {expandedComments[post.id] && (
                    <div className="mt-3 space-y-3">
                      <div className="flex gap-2">
                        <input
                          className="flex h-9 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm"
                          placeholder="Escribí un comentario..."
                          value={commentText[post.id] ?? ''}
                          onChange={e => setCommentText((prev) => ({ ...prev, [post.id]: e.target.value }))}
                          onKeyDown={e => { if (e.key === 'Enter') addComment(post) }}
                        />
                        <Button size="sm" onClick={() => addComment(post)} disabled={commentLoading[post.id] || !(commentText[post.id]?.trim())}>
                          <Send size={14} />
                        </Button>
                      </div>
                      {(comments[post.id] ?? []).map((c) => {
                        const cInitials = (c.author_name ?? 'U').split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
                        return (
                          <div key={c.id} className="flex items-start gap-2">
                            <Avatar className="h-7 w-7">
                              {c.author_photo && <AvatarImage src={c.author_photo} alt="" />}
                              <AvatarFallback className="bg-slate-100 text-slate-600 text-xs">{cInitials}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1 rounded-lg bg-slate-50 px-3 py-2">
                              <div className="flex items-center gap-2">
                                <span className="text-xs font-semibold text-slate-900">{c.author_name || 'Empleado'}</span>
                                <span className="text-xs text-slate-400">{timeAgo(c.created_at)}</span>
                              </div>
                              <p className="text-sm text-slate-700 mt-1">{c.content}</p>
                            </div>
                          </div>
                        )
                      })}
                    </div>
                  )}
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
