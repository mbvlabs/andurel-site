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
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
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

function PagerCard({
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
      className={cn('block h-full min-w-0', align === 'end' && 'sm:col-start-2')}
    >
      <Card size="sm" className="h-full transition-colors hover:bg-muted/50">
        <CardHeader className={align === 'end' ? 'items-end text-right' : undefined}>
          <CardDescription>{label}</CardDescription>
          <CardTitle className="flex min-w-0 items-center gap-1.5">
            {align === 'start' && <ChevronLeftIcon className="size-4 shrink-0" />}
            <span className="truncate">{page.title}</span>
            {align === 'end' && <ChevronRightIcon className="size-4 shrink-0" />}
          </CardTitle>
        </CardHeader>
      </Card>
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
      <Breadcrumb className="mb-6 min-w-0 sm:mb-7">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink render={<Link href={versionUrl} />}>{currentVersion}</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <span>{currentSection}</span>
          </BreadcrumbItem>
          {parent ? (
            <>
              <BreadcrumbSeparator />
              <BreadcrumbItem className="min-w-0">
                <BreadcrumbLink render={<Link href={parent.url} />} className="truncate">
                  {parent.title}
                </BreadcrumbLink>
              </BreadcrumbItem>
            </>
          ) : null}
          <BreadcrumbSeparator />
          <BreadcrumbItem className="min-w-0">
            <BreadcrumbPage className="truncate">{title}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <article className="docs-content min-w-0" dangerouslySetInnerHTML={{ __html: html }} />
      {(previous || next) && (
        <>
          <Separator className="mt-10 sm:mt-14" />
          <nav className="mt-6 grid gap-3 sm:grid-cols-2" aria-label="Documentation pagination">
            {previous ? <PagerCard align="start" label="Previous" page={previous} /> : null}
            {next ? <PagerCard align="end" label="Next" page={next} /> : null}
          </nav>
        </>
      )}
    </DocLayout>
  )
}
