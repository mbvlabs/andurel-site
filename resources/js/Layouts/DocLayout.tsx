import type { ReactNode } from 'react'

import { Link } from '@inertiajs/react'
import { ChevronDownIcon } from 'lucide-react'

import DocsSearch from '@/components/docs-search'
import DocsToc from '@/components/docs-toc'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { routes } from '@/routes'
import type { DocHeading, DocNavigationProps, DocVersion } from '@/types/docs'

const DOC_ARTICLE_ID = 'doc-article'
const docsPadRight = 'pr-10 sm:pr-12'

type DocLayoutProps = DocNavigationProps & {
  children: ReactNode
  headings?: DocHeading[]
}

function currentCatalog(versions: DocVersion[], currentVersion: string) {
  return versions.find((version) => version.name === currentVersion) ?? versions[0]
}

function BrandMark() {
  return (
    <span className="grid size-8 shrink-0 grid-cols-2 gap-1 border border-sidebar-border bg-sidebar p-1 shadow-sm shadow-black/40">
      <span className="border border-sidebar-primary" />
      <span className="border border-sidebar-border" />
      <span className="border border-sidebar-border" />
      <span className="bg-sidebar-primary" />
    </span>
  )
}

function DocBrand() {
  return (
    <a
      className="flex h-full items-center gap-3 text-sm font-semibold text-sidebar-foreground"
      href={routes.homePage()}
    >
      <BrandMark />
      <span>Andurel.</span>
    </a>
  )
}

function HeaderActions({
  versions,
  currentVersion,
}: {
  versions: DocVersion[]
  currentVersion: string
}) {
  return (
    <nav className="flex shrink-0 items-center gap-2 text-sm">
      <DropdownMenu>
        <DropdownMenuTrigger render={<Button variant="outline" size="sm" />}>
          {currentVersion}
          <ChevronDownIcon />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuGroup>
            <DropdownMenuLabel>Versions</DropdownMenuLabel>
            {versions.map((version) => (
              <DropdownMenuItem key={version.name} render={<Link href={version.url} />}>
                {version.name}
              </DropdownMenuItem>
            ))}
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
      <Button
        variant="ghost"
        size="sm"
        render={<a href="https://github.com/mbvlabs/andurel" />}
      >
        GitHub
      </Button>
    </nav>
  )
}

export default function DocLayout({
  children,
  versions,
  currentVersion,
  currentSlug,
  headings,
}: DocLayoutProps) {
  const catalog = currentCatalog(versions, currentVersion)

  return (
    <SidebarProvider>
      <Sidebar className="border-sidebar-border">
        <SidebarHeader className="h-14 flex-row items-center border-b border-sidebar-border px-4 py-0">
          <DocBrand />
        </SidebarHeader>
        <SidebarContent>
          {catalog?.sections.map((section) => (
            <SidebarGroup key={section.title}>
              <SidebarGroupLabel>{section.title}</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  {section.pages.map((page) => (
                    <SidebarMenuItem key={page.slug}>
                      <SidebarMenuButton
                        isActive={currentSlug === page.slug}
                        aria-current={currentSlug === page.slug ? 'page' : undefined}
                        tooltip={page.title}
                        render={<Link href={page.url} />}
                      >
                        <span>{page.title}</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          ))}
        </SidebarContent>
      </Sidebar>

      <SidebarInset>
        <header className="sticky top-0 z-30 border-b border-sidebar-border bg-background">
          <div className={`flex h-14 items-center ${docsPadRight}`}>
            <div className="flex min-w-0 flex-1 items-center pl-2">
              <SidebarTrigger />
              <a
                className="ml-3 inline-flex items-center gap-2 text-sm font-semibold md:hidden"
                href={routes.homePage()}
              >
                <BrandMark />
                <span className="sr-only">Andurel Docs</span>
              </a>
            </div>
            <div className="w-full min-w-0 max-w-3xl">
              <DocsSearch versions={versions} currentVersion={currentVersion} />
            </div>
            <div className="flex flex-1 items-center justify-end">
              <HeaderActions versions={versions} currentVersion={currentVersion} />
            </div>
          </div>
        </header>

        <div className={`flex flex-1 py-8 ${docsPadRight}`}>
          <div className="min-w-0 flex-1" aria-hidden="true" />
          <div id={DOC_ARTICLE_ID} className="w-full min-w-0 max-w-3xl">
            {children}
          </div>
          <div className="min-w-0 flex-1">
            <aside className="sticky top-20 hidden max-w-64 self-start pl-6 xl:block">
              <DocsToc
                rootId={DOC_ARTICLE_ID}
                pageKey={`${currentVersion}:${currentSlug}`}
                headings={headings}
              />
            </aside>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
