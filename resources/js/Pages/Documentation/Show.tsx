import { Link } from '@inertiajs/react'
import { cn } from 'cn'
import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-react'

import DocLayout from '@/Layouts/DocLayout'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import type { DocHeading, DocLink, DocNavigationProps } from '@/types/docs'

type ShowProps = DocNavigationProps & {
  currentSection: string
  title: string
  description: string
  html: string
  headings: DocHeading[]
  parent: DocLink | null
  previous: DocLink | null
  next: DocLink | null
}

function VersionCrumb({ name }: { name: string }) {
  return (
    <span className="inline-flex items-center gap-1.5">
      <span>{name}</span>
      {name === 'head' ? (
        <span className="hidden text-[10px] font-medium tracking-wide text-primary/90 sm:inline">
          v2 alpha
        </span>
      ) : null}
    </span>
  )
}

function PagerLink({
  align,
  label,
  page,
}: {
  align: 'start' | 'end'
  label: string
  page: DocLink
}) {
  return (
    <Link
      href={page.url}
      className={cn(
        'group flex min-w-0 flex-col gap-1 rounded-none border border-border/70 bg-transparent px-4 py-3 transition-colors hover:border-border hover:bg-muted/30',
        align === 'end' && 'items-end text-right sm:col-start-2',
      )}
    >
      <span className="text-[11px] font-medium tracking-wide text-muted-foreground uppercase">
        {label}
      </span>
      <span className="flex min-w-0 items-center gap-1.5 text-sm font-medium text-foreground">
        {align === 'start' && (
          <ChevronLeftIcon className="size-4 shrink-0 text-muted-foreground transition-colors group-hover:text-foreground" />
        )}
        <span className="truncate">{page.title}</span>
        {align === 'end' && (
          <ChevronRightIcon className="size-4 shrink-0 text-muted-foreground transition-colors group-hover:text-foreground" />
        )}
      </span>
    </Link>
  )
}

export default function Show({
  versions,
  currentVersion,
  currentSlug,
  currentSection,
  title,
  description,
  html,
  headings,
  parent,
  previous,
  next,
}: ShowProps) {
  const versionUrl =
    versions.find((version) => version.name === currentVersion)?.url ??
    `/docs/${currentVersion}`

  return (
    <DocLayout
      versions={versions}
      currentVersion={currentVersion}
      currentSlug={currentSlug}
      currentSection={currentSection}
      title={title}
      description={description}
      headings={headings}
    >
      <Breadcrumb className="mb-5 min-w-0 sm:mb-6">
        <BreadcrumbList className="gap-1 text-[11px] tracking-wide sm:gap-1.5">
          <BreadcrumbItem>
            <BreadcrumbLink
              render={<Link href={versionUrl} />}
              className="text-muted-foreground hover:text-foreground"
            >
              <VersionCrumb name={currentVersion} />
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator className="text-muted-foreground/50" />
          <BreadcrumbItem>
            <span className="text-muted-foreground">{currentSection}</span>
          </BreadcrumbItem>
          {parent ? (
            <>
              <BreadcrumbSeparator className="text-muted-foreground/50" />
              <BreadcrumbItem className="min-w-0">
                <BreadcrumbLink
                  render={<Link href={parent.url} />}
                  className="truncate text-muted-foreground hover:text-foreground"
                >
                  {parent.title}
                </BreadcrumbLink>
              </BreadcrumbItem>
            </>
          ) : null}
          <BreadcrumbSeparator className="text-muted-foreground/50" />
          <BreadcrumbItem className="min-w-0">
            <BreadcrumbPage className="truncate font-medium text-foreground/90">
              {title}
            </BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <article className="docs-content min-w-0" dangerouslySetInnerHTML={{ __html: html }} />

      {(previous || next) && (
        <nav
          className="mt-10 grid gap-3 border-t border-border/80 pt-6 sm:mt-12 sm:grid-cols-2 sm:pt-8"
          aria-label="Documentation pagination"
        >
          {previous ? <PagerLink align="start" label="Previous" page={previous} /> : null}
          {next ? <PagerLink align="end" label="Next" page={next} /> : null}
        </nav>
      )}
    </DocLayout>
  )
}
