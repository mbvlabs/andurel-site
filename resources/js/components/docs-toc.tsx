import { useEffect, useState } from 'react'

import { cn } from 'cn'
import { AlignLeftIcon } from 'lucide-react'

type Heading = {
  id: string
  title: string
  level: number
}

type DocsTocProps = {
  rootId: string
  pageKey: string
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
        title,
        level: Number(element.tagName[1]),
      }
    })
    .filter((heading): heading is Heading => heading !== null)
}

function preferSectionHeadings(headings: Heading[]) {
  const sections = headings.filter((heading) => heading.level >= 2)
  return sections.length > 0 ? sections : headings
}

export default function DocsToc({ rootId, pageKey }: DocsTocProps) {
  const [headings, setHeadings] = useState<Heading[]>([])
  const [activeId, setActiveId] = useState('')

  useEffect(() => {
    const nextHeadings = preferSectionHeadings(collectHeadings(rootId))
    setHeadings(nextHeadings)
    setActiveId(nextHeadings[0]?.id ?? '')

    if (nextHeadings.length === 0) {
      return
    }

    const observed = nextHeadings
      .map((heading) => document.getElementById(heading.id))
      .filter((element): element is HTMLElement => element !== null)

    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries
          .filter((entry) => entry.isIntersecting)
          .sort((left, right) => right.intersectionRatio - left.intersectionRatio)

        if (visible[0]?.target.id) {
          setActiveId(visible[0].target.id)
        }
      },
      {
        rootMargin: '-80px 0px -60% 0px',
        threshold: [0, 1],
      }
    )

    observed.forEach((element) => observer.observe(element))
    return () => observer.disconnect()
  }, [pageKey, rootId])

  if (headings.length === 0) {
    return null
  }

  return (
    <nav aria-label="On this page" className="sticky top-14">
      <p className="mb-3 flex items-center gap-2 text-xs font-semibold text-foreground">
        <AlignLeftIcon className="size-3.5 text-muted-foreground" />
        On this page
      </p>
      <ul className="border-l border-sidebar-border">
        {headings.map((heading) => (
          <li key={heading.id}>
            <a
              href={`#${heading.id}`}
              className={cn(
                '-ml-px block border-l-2 py-1 text-xs leading-5 transition-colors',
                heading.level > 2 ? 'pl-6' : 'pl-3',
                activeId === heading.id
                  ? 'border-accent text-foreground'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              )}
            >
              {heading.title}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  )
}
