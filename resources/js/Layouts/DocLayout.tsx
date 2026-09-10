import type { ReactNode } from 'react'

import SiGithub from '@icons-pack/react-simple-icons/icons/SiGithub'
import { Link, usePage } from '@inertiajs/react'
import { ChevronDownIcon, ChevronRightIcon } from 'lucide-react'

import DocsSearch from '@/components/docs-search'
import DocsToc from '@/components/docs-toc'
import SeoHead from '@/components/seo-head'
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
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { docsJsonLd, useCanonicalUrl } from '@/lib/seo'
import { routes } from '@/routes'
import type { DocHeading, DocNavigationProps, DocPage, DocVersion } from '@/types/docs'
import type { SharedPageProps } from '@/types/page'

const DOC_ARTICLE_ID = 'doc-article'
const docsGutter = 'px-4 sm:px-6 lg:px-8'
const docsColumns =
  'grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] md:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(0,1fr)] xl:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(14rem,1fr)] xl:gap-x-24'
const headerColumns =
  'grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] gap-x-3 sm:gap-x-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(14rem,1fr)] xl:gap-x-24'

type DocLayoutProps = DocNavigationProps & {
  children: ReactNode
  headings?: DocHeading[]
  title: string
  description: string
  currentSection: string
}

function currentCatalog(versions: DocVersion[], currentVersion: string) {
  return versions.find((version) => version.name === currentVersion) ?? versions[0]
}

function pageOrDescendantIsCurrent(page: DocPage, currentSlug: string): boolean {
  if (page.slug === currentSlug) {
    return true
  }
  return page.children?.some((child) => pageOrDescendantIsCurrent(child, currentSlug)) ?? false
}

function DocNavPage({ page, currentSlug }: { page: DocPage; currentSlug: string }) {
  const children = page.children ?? []
  const expanded = children.length > 0 && pageOrDescendantIsCurrent(page, currentSlug)

  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        isActive={currentSlug === page.slug}
        aria-current={currentSlug === page.slug ? 'page' : undefined}
        aria-expanded={children.length > 0 ? expanded : undefined}
        tooltip={page.title}
        render={<Link href={page.url} />}
      >
        <span>{page.title}</span>
        {children.length > 0 ? (
          <ChevronRightIcon className={expanded ? 'ml-auto rotate-90' : 'ml-auto'} />
        ) : null}
      </SidebarMenuButton>
      {expanded ? (
        <SidebarMenuSub>
          {children.map((child) => (
            <SidebarMenuSubItem key={child.slug}>
              <SidebarMenuSubButton
                isActive={currentSlug === child.slug}
                aria-current={currentSlug === child.slug ? 'page' : undefined}
                render={<Link href={child.url} />}
              >
                <span>{child.title}</span>
              </SidebarMenuSubButton>
            </SidebarMenuSubItem>
          ))}
        </SidebarMenuSub>
      ) : null}
    </SidebarMenuItem>
  )
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
  currentSection,
  title,
  description,
  headings,
}: DocLayoutProps) {
  const catalog = currentCatalog(versions, currentVersion)
  const canonical = useCanonicalUrl()
  const { appUrl } = usePage<SharedPageProps>().props

  return (
    <SidebarProvider>
      <SeoHead
        title={title}
        description={description}
        type="article"
        suffix="Andurel Docs"
        jsonLd={docsJsonLd({
          canonical,
          title,
          description,
          version: currentVersion,
          section: currentSection,
          appUrl,
        })}
      />
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
                    <DocNavPage key={page.slug} page={page} currentSlug={currentSlug} />
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
