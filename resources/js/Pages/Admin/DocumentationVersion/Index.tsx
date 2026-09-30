import { Link, router } from '@inertiajs/react'
import { BookOpenIcon, GripVerticalIcon } from 'lucide-react'
import { useEffect, useState } from 'react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { routes } from '@/routes'
import type { DocumentationVersionData, DocumentationVersionIndexProps } from '@/types/payloads'

function moveItem<T>(items: T[], from: number, to: number): T[] {
  if (from === to || from < 0 || to < 0 || from >= items.length || to >= items.length) {
    return items
  }
  const copy = [...items]
  const [removed] = copy.splice(from, 1)
  copy.splice(to, 0, removed)
  return copy
}

function sameOrder(left: DocumentationVersionData[], right: DocumentationVersionData[]) {
  return left.length === right.length && left.every((item, index) => item.id === right[index]?.id)
}

export default function Index({ versions }: DocumentationVersionIndexProps) {
  const [ordered, setOrdered] = useState(versions)
  const [dragIndex, setDragIndex] = useState<number | null>(null)
  const [overIndex, setOverIndex] = useState<number | null>(null)

  useEffect(() => {
    setOrdered(versions)
  }, [versions])

  function persist(next: DocumentationVersionData[]) {
    if (sameOrder(next, versions)) {
      return
    }
    setOrdered(next)
    router.put(
      routes.adminDocumentationVersionReorder.path(),
      { ids: next.map((version) => version.id) },
      {
        preserveScroll: true,
        only: ['versions'],
        onError: () => setOrdered(versions),
      },
    )
  }

  return (
    <DashboardLayout title="Docs">
      <div className="mb-6 flex items-center justify-between gap-3">
        <p className="max-w-xl text-sm text-muted-foreground">
          Versions, pages, and sidebar drafts. Drag the handle to change public version order.
        </p>
        <Button render={<Link href={routes.adminDocumentationVersionNew.path()} />}>New version</Button>
      </div>

      {ordered.length === 0 ? (
        <Card>
          <CardHeader>
            <CardTitle>No versions yet</CardTitle>
            <CardDescription>Create head, 1.5.5, and 1.5.2 here, then port Markdown into each page.</CardDescription>
          </CardHeader>
        </Card>
      ) : (
        <div className="grid gap-3">
          {ordered.map((version, index) => (
            <Card
              key={version.id}
              onDragOver={(event) => {
                event.preventDefault()
                event.dataTransfer.dropEffect = 'move'
                if (overIndex !== index) {
                  setOverIndex(index)
                }
              }}
              onDrop={(event) => {
                event.preventDefault()
                const from = dragIndex
                setDragIndex(null)
                setOverIndex(null)
                if (from == null) {
                  return
                }
                persist(moveItem(ordered, from, index))
              }}
              onDragLeave={(event) => {
                if (event.currentTarget.contains(event.relatedTarget as Node | null)) {
                  return
                }
                if (overIndex === index) {
                  setOverIndex(null)
                }
              }}
              className={cn(overIndex === index && dragIndex !== null && dragIndex !== index && 'ring-primary')}
            >
              <CardHeader className="flex flex-row items-start justify-between gap-3">
                <div className="flex min-w-0 items-start gap-3">
                  <div
                    draggable
                    role="button"
                    tabIndex={0}
                    aria-label={`Reorder ${version.label}`}
                    className="mt-0.5 shrink-0 cursor-grab touch-none text-muted-foreground active:cursor-grabbing"
                    onDragStart={(event) => {
                      event.dataTransfer.effectAllowed = 'move'
                      event.dataTransfer.setData('text/plain', String(version.id))
                      setDragIndex(index)
                    }}
                    onDragEnd={() => {
                      setDragIndex(null)
                      setOverIndex(null)
                    }}
                  >
                    <GripVerticalIcon className="size-4" />
                  </div>
                  <div className="min-w-0">
                    <CardTitle className="flex items-center gap-2">
                      <BookOpenIcon className="size-4" />
                      {version.label}
                    </CardTitle>
                    <CardDescription>
                      /docs/{version.slug}
                      {version.publishedNavId ? ' · nav released' : ' · nav unpublished'}
                    </CardDescription>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {version.isLatest ? <Badge>Latest</Badge> : null}
                  <Button
                    variant="outline"
                    size="sm"
                    render={<Link href={routes.adminDocumentationVersionShow.path(version.id)} />}
                  >
                    Open
                  </Button>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}
    </DashboardLayout>
  )
}
