import { useEffect, useState } from 'react'

import { cn } from 'cn'
import { AlignLeftIcon } from 'lucide-react'

type Heading = {
  id: string
  text: string
  level: number
}

type DocsTocProps = {
  rootId: string
  pageKey: string
  headings?: Heading[]
}

function collectHeadings(rootId: string): Heading[] {
  const root = document.getElementById(rootId)
  if (!root) {
    return []
  }

  return [...root.querySelectorAll<HTMLElement>('h1, h2, h3')]
    .map((element) => {
      const title = element.textContent?.trim() ?? ''
      if (!title) {
        return null
      }

      if (!element.id) {
        element.id = title
          .toLowerCase()
          .replace(/[^a-z0-9]+/g, '-')
          .replace(/^-|-$/g, '')
      }

      return {
        id: element.id,
        text: title,
        level: Number(element.tagName[1]),
      }
    })
    .filter((heading): heading is Heading => heading !== null)
}

function preferSectionHeadings(headings: Heading[]) {
  const sections = headings.filter((heading) => heading.level >= 2)
  return sections.length > 0 ? sections : headings
}

const HEADER_OFFSET = 96

export default function DocsToc({ rootId, pageKey, headings: providedHeadings }: DocsTocProps) {
  const [headings, setHeadings] = useState<Heading[]>(providedHeadings ?? [])
  const [activeId, setActiveId] = useState(providedHeadings?.[0]?.id ?? '')

  useEffect(() => {
    const nextHeadings =
      providedHeadings && providedHeadings.length > 0
        ? providedHeadings
        : preferSectionHeadings(collectHeadings(rootId))
    setHeadings(nextHeadings)

    if (nextHeadings.length === 0) {
      return
    }

    let frame = 0

    function syncActive() {
      const focusY = Math.max(HEADER_OFFSET, Math.round(window.innerHeight * 0.3))
      let current = nextHeadings[0].id
      for (const heading of nextHeadings) {
        const element = document.getElementById(heading.id)
        if (element && element.getBoundingClientRect().top <= focusY) {
          current = heading.id
        }
      }
      const scrolledToEnd =
        window.innerHeight + window.scrollY >= document.documentElement.scrollHeight - 8
      if (scrolledToEnd) {
        current = nextHeadings[nextHeadings.length - 1].id
      }
      setActiveId((prev) => (prev === current ? prev : current))
    }

    function onScroll() {
      if (frame) {
        return
      }
      frame = window.requestAnimationFrame(() => {
        frame = 0
        syncActive()
      })
    }

    syncActive()
    window.addEventListener('scroll', onScroll, { passive: true })
    window.addEventListener('scrollend', syncActive)
    window.addEventListener('hashchange', syncActive)
    return () => {
      window.cancelAnimationFrame(frame)
      window.removeEventListener('scroll', onScroll)
      window.removeEventListener('scrollend', syncActive)
      window.removeEventListener('hashchange', syncActive)
    }
  }, [pageKey, providedHeadings, rootId])

  if (headings.length === 0) {
    return null
  }

  return (
    <nav aria-label="On this page">
      <p className="mb-3 flex items-center gap-2 text-xs font-semibold text-foreground">
        <AlignLeftIcon className="size-3.5 text-muted-foreground" />
        On this page
      </p>
      <ul className="border-l border-sidebar-border">
        {headings.map((heading) => (
          <li key={heading.id}>
            <a
              href={`#${heading.id}`}
              aria-current={activeId === heading.id ? 'location' : undefined}
              onClick={() => setActiveId(heading.id)}
              className={cn(
                '-ml-px block border-l-2 py-1 text-xs leading-5 transition-colors',
                heading.level > 2 ? 'pl-6' : 'pl-3',
                activeId === heading.id
                  ? 'border-accent text-foreground'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              )}
            >
              {heading.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  )
}
