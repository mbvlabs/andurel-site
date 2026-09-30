import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import type { DocumentationNavNodeData, DocumentationPageSummaryData } from '@/types/payloads'

type NavTreeEditorProps = {
  tree: DocumentationNavNodeData[]
  pages: DocumentationPageSummaryData[]
  onChange: (tree: DocumentationNavNodeData[]) => void
}

function emptyNode(page?: DocumentationPageSummaryData): DocumentationNavNodeData {
  return {
    title: page?.title ?? 'New section',
    pageId: page?.id ?? null,
    children: [],
  }
}

function moveItem<T>(items: T[], index: number, offset: number): T[] {
  const next = index + offset
  if (next < 0 || next >= items.length) {
    return items
  }
  const copy = [...items]
  const [removed] = copy.splice(index, 1)
  copy.splice(next, 0, removed)
  return copy
}

type NodeEditorProps = {
  node: DocumentationNavNodeData
  pages: DocumentationPageSummaryData[]
  onChange: (node: DocumentationNavNodeData) => void
  onRemove: () => void
  onMove: (offset: number) => void
}

function NodeEditor({ node, pages, onChange, onRemove, onMove }: NodeEditorProps) {
  function selectPage(pageId: number | null) {
    const page = pages.find((item) => item.id === pageId)
    onChange({
      ...node,
      pageId,
      title: page?.title ?? node.title,
    })
  }

  return (
    <div className="space-y-3 border border-border bg-background p-3">
      <div className="grid gap-3 sm:grid-cols-2">
        <div className="space-y-1">
          <Label>Label</Label>
          <Input
            value={node.title}
            onChange={(event) => onChange({ ...node, title: event.currentTarget.value })}
          />
        </div>
        <div className="space-y-1">
          <Label>Page</Label>
          <Select
            value={node.pageId === null ? 'section' : String(node.pageId)}
            onValueChange={(value) => {
              selectPage(value === 'section' || value == null ? null : Number(value))
            }}
            items={[
              { value: 'section', label: 'Section' },
              ...pages.map((page) => ({
                value: String(page.id),
                label: `${page.title} (${page.slug})`,
              })),
            ]}
          >
            <SelectTrigger className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent align="start" alignItemWithTrigger={false}>
              <SelectItem value="section">Section</SelectItem>
              {pages.map((page) => (
                <SelectItem key={page.id} value={String(page.id)}>
                  {page.title} ({page.slug})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
      <div className="flex flex-wrap gap-2">
        <Button type="button" size="sm" variant="outline" onClick={() => onMove(-1)}>
          Up
        </Button>
        <Button type="button" size="sm" variant="outline" onClick={() => onMove(1)}>
          Down
        </Button>
        <Button
          type="button"
          size="sm"
          variant="outline"
          onClick={() =>
            onChange({
              ...node,
              children: [...node.children, emptyNode(pages[0])],
            })
          }
        >
          Add child
        </Button>
        <Button type="button" size="sm" variant="ghost" onClick={onRemove}>
          Remove
        </Button>
      </div>
      {node.children.length > 0 ? (
        <div className="space-y-3 border-l border-border pl-3">
          {node.children.map((child, index) => (
            <NodeEditor
              key={`${child.pageId ?? 'node'}-${index}`}
              node={child}
              pages={pages}
              onChange={(next) => {
                const children = [...node.children]
                children[index] = next
                onChange({ ...node, children })
              }}
              onRemove={() => {
                const children = node.children.filter((_, childIndex) => childIndex !== index)
                onChange({ ...node, children })
              }}
              onMove={(offset) => onChange({ ...node, children: moveItem(node.children, index, offset) })}
            />
          ))}
        </div>
      ) : null}
    </div>
  )
}

export function NavTreeEditor({ tree, pages, onChange }: NavTreeEditorProps) {
  return (
    <div className="space-y-3">
      {tree.map((node, index) => (
        <NodeEditor
          key={`${node.pageId ?? 'node'}-${index}`}
          node={node}
          pages={pages}
          onChange={(next) => {
            const copy = [...tree]
            copy[index] = next
            onChange(copy)
          }}
          onRemove={() => onChange(tree.filter((_, nodeIndex) => nodeIndex !== index))}
          onMove={(offset) => onChange(moveItem(tree, index, offset))}
        />
      ))}
      <Button type="button" variant="outline" onClick={() => onChange([...tree, emptyNode()])}>
        Add section
      </Button>
    </div>
  )
}
