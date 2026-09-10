import type { ReactNode } from 'react'

import SiGithub from '@icons-pack/react-simple-icons/icons/SiGithub'
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
const docsGutter = 'px-4 sm:px-6 lg:px-8'
const docsColumns =
  'grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] md:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(0,1fr)] xl:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(14rem,1fr)] xl:gap-x-24'
const headerColumns =
  'grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] gap-x-3 sm:gap-x-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(14rem,1fr)] xl:gap-x-24'

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
    <nav className="flex shrink-0 items-center gap-1 sm:gap-2 text-sm">
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
        aria-label="GitHub"
        nativeButton={false}
        render={<a href="https://github.com/mbvlabs/andurel" />}
      >
        <SiGithub title="" color="currentColor" />
        <span className="hidden sm:inline">GitHub</span>
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

      <SidebarInset className="min-w-0 overflow-x-clip">
        <header className="sticky top-0 z-30 border-b border-sidebar-border bg-background">
          <div className={`h-14 items-center ${headerColumns} ${docsGutter}`}>
            <div className="flex shrink-0 items-center">
              <SidebarTrigger />
              <a
                className="ml-2 hidden items-center gap-2 text-sm font-semibold sm:inline-flex md:hidden"
                href={routes.homePage()}
              >
                <BrandMark />
                <span className="sr-only">Andurel Docs</span>
              </a>
            </div>
            <div className="flex min-w-0 items-center justify-end md:block">
              <DocsSearch versions={versions} currentVersion={currentVersion} />
            </div>
            <div className="flex shrink-0 items-center justify-end">
              <HeaderActions versions={versions} currentVersion={currentVersion} />
            </div>
          </div>
        </header>

        <div className={`min-w-0 flex-1 items-start py-6 sm:py-8 ${docsColumns} ${docsGutter}`}>
          <div aria-hidden="true" />
          <div id={DOC_ARTICLE_ID} className="min-w-0">
            {children}
          </div>
          <aside className="sticky top-20 hidden w-56 min-w-0 self-start xl:block">
            <DocsToc
              rootId={DOC_ARTICLE_ID}
              pageKey={`${currentVersion}:${currentSlug}`}
              headings={headings}
            />
          </aside>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
