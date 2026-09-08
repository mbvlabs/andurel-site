import { useEffect, useState } from 'react'

import { router } from '@inertiajs/react'
import { cn } from 'cn'
import { SearchIcon } from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  Command,
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import type { DocVersion } from '@/types/docs'

type DocsSearchProps = {
  versions: DocVersion[]
  currentVersion: string
  className?: string
}

export default function DocsSearch({
  versions,
  currentVersion,
  className,
}: DocsSearchProps) {
  const [open, setOpen] = useState(false)
  const catalog = versions.find((version) => version.name === currentVersion) ?? versions[0]

  useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
        event.preventDefault()
        setOpen((current) => !current)
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [])

  function visit(url: string) {
    setOpen(false)
    router.visit(url)
  }

  return (
    <>
      <Button
        type="button"
        variant="outline"
        className={cn('h-8 w-full justify-start text-muted-foreground', className)}
        onClick={() => setOpen(true)}
      >
        <SearchIcon />
        <span className="flex-1 truncate text-left">Search documentation...</span>
        <kbd className="pointer-events-none hidden h-5 items-center gap-1 border border-border bg-muted px-1.5 font-mono text-[0.65rem] text-muted-foreground sm:inline-flex">
          ⌘K
        </kbd>
      </Button>
      <CommandDialog
        open={open}
        onOpenChange={setOpen}
        title="Search documentation"
        description="Find a page in the Andurel docs."
      >
        <Command>
          <CommandInput placeholder="Search documentation..." />
          <CommandList>
            <CommandEmpty>No documentation found.</CommandEmpty>
            {catalog?.sections.map((section) => (
              <CommandGroup key={section.title} heading={section.title}>
                {section.pages.map((page) => (
                  <CommandItem
                    key={page.url}
                    value={`${section.title} ${page.title} ${page.description}`}
                    onSelect={() => visit(page.url)}
                  >
                    <span className="flex min-w-0 flex-col">
                      <span>{page.title}</span>
                      <span className="truncate text-muted-foreground">{page.description}</span>
                    </span>
                  </CommandItem>
                ))}
              </CommandGroup>
            ))}
          </CommandList>
        </Command>
      </CommandDialog>
    </>
  )
}
