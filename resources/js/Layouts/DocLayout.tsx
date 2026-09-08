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
import type { DocNavigationProps, DocVersion } from '@/types/docs'

const DOC_ARTICLE_ID = 'doc-article'
const docsPad = 'px-14 sm:px-16'

type DocLayoutProps = DocNavigationProps & {
  children: ReactNode
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
      <a
        className="px-2 py-1 text-muted-foreground transition hover:text-foreground"
        href="https://github.com/mbvlabs/andurel"
      >
        GitHub
      </a>
    </nav>
  )
}

export default function DocLayout({
  children,
  versions,
  currentVersion,
  currentSlug,
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
          <div className="relative flex h-14 items-center">
            <div className={`absolute inset-y-0 left-0 z-10 flex items-center ${docsPad}`}>
              <SidebarTrigger />
              <a
                className="ml-3 inline-flex items-center gap-2 text-sm font-semibold md:hidden"
                href={routes.homePage()}
              >
                <BrandMark />
                <span className="sr-only">Andurel Docs</span>
              </a>
            </div>
            <div className={`flex w-full items-center justify-center ${docsPad}`}>
              <div className="w-full max-w-3xl">
                <DocsSearch versions={versions} currentVersion={currentVersion} />
              </div>
            </div>
            <div className={`absolute inset-y-0 right-0 z-10 flex items-center ${docsPad}`}>
              <HeaderActions versions={versions} currentVersion={currentVersion} />
            </div>
          </div>
        </header>

        <div className={`flex flex-1 py-8 ${docsPad}`}>
          <div className="hidden w-64 shrink-0 xl:block" aria-hidden="true" />
          <div className="flex min-w-0 flex-1 justify-center">
            <div id={DOC_ARTICLE_ID} className="w-full max-w-3xl">
              {children}
            </div>
          </div>
          <aside className="hidden w-64 shrink-0 xl:block">
            <DocsToc rootId={DOC_ARTICLE_ID} pageKey={`${currentVersion}:${currentSlug}`} />
          </aside>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
