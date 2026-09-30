import { router, useForm } from '@inertiajs/react'
import { useEffect, useState } from 'react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { MarkdownEditor } from '@/components/MarkdownEditor'
import { AdminBreadcrumb } from '@/components/docs/AdminBreadcrumb'
import DocsToc from '@/components/docs-toc'
import { Badge } from '@/components/ui/badge'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { Textarea } from '@/components/ui/textarea'
import { slugify } from '@/lib/slug'
import { routes } from '@/routes'
import type { DocumentationHeadingData, DocumentationPageItemProps, UpdateDocumentationPageFormPayload } from '@/types/payloads'

const PREVIEW_ARTICLE_ID = 'doc-preview-article'
const docsGutter = 'px-4 sm:px-6 lg:px-8'
const docsColumns =
  'grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] md:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(0,1fr)] xl:grid-cols-[minmax(0,1fr)_minmax(0,48rem)_minmax(14rem,1fr)] xl:gap-x-24'

function lengthTone(length: number, recommendedMin: number, recommendedMax: number) {
  if (length === 0) {
    return 'text-muted-foreground'
  }
  if (length < recommendedMin || length > recommendedMax) {
    return 'text-amber-600 dark:text-amber-400'
  }
  return 'text-muted-foreground'
}

function LengthHint({
  value,
  recommendedMin,
  recommendedMax,
}: {
  value: string
  recommendedMin: number
  recommendedMax: number
}) {
  const length = value.length
  return (
    <p className={`text-xs ${lengthTone(length, recommendedMin, recommendedMax)}`}>
      {length} characters · {recommendedMin}–{recommendedMax} recommended
    </p>
  )
}

function ModeSwitcher({
  mode,
  onWrite,
  onPreview,
}: {
  mode: 'write' | 'preview'
  onWrite: () => void
  onPreview: () => void
}) {
  return (
    <div className="inline-flex border border-border">
      <Button type="button" size="xs" variant={mode === 'write' ? 'secondary' : 'ghost'} onClick={onWrite}>
        Write
      </Button>
      <Button type="button" size="xs" variant={mode === 'preview' ? 'secondary' : 'ghost'} onClick={onPreview}>
        Preview
      </Button>
    </div>
  )
}

export default function Show({ version, page, draft, published, revisions }: DocumentationPageItemProps) {
  const form = useForm<UpdateDocumentationPageFormPayload>({
    slug: page.slug,
    title: draft.title,
    metaTitle: draft.metaTitle,
    description: draft.description,
    bodyMarkdown: draft.bodyMarkdown,
  })
  const [slugIsAutomatic, setSlugIsAutomatic] = useState(
    () => page.slug === '' || page.slug === slugify(draft.title),
  )
  const [pageMode, setPageMode] = useState<'write' | 'preview'>('write')
  const [previewHtml, setPreviewHtml] = useState('')
  const [previewHeadings, setPreviewHeadings] = useState<DocumentationHeadingData[]>([])
  const [previewSource, setPreviewSource] = useState('')
  const [previewError, setPreviewError] = useState('')
  const [previewing, setPreviewing] = useState(false)

  useEffect(() => {
    form.setData({
      slug: page.slug,
      title: draft.title,
      metaTitle: draft.metaTitle,
      description: draft.description,
      bodyMarkdown: draft.bodyMarkdown,
    })
    setSlugIsAutomatic(page.slug === '' || page.slug === slugify(draft.title))
  }, [draft.bodyMarkdown, draft.description, draft.id, draft.metaTitle, draft.title, page.slug])

  function saveDraft(event: SubmitEvent) {
    event.preventDefault()
    form.put(routes.adminDocumentationPageUpdate.path(page.id))
  }

  function publishPage() {
    router.post(routes.adminDocumentationPagePublish.path(page.id), form.data)
  }

  function restoreRevision(revisionId: number) {
    router.post(routes.adminDocumentationPageRestore.path(String(page.id), String(revisionId)))
  }

  async function showPreview() {
    setPageMode('preview')
    if (previewSource === form.data.bodyMarkdown && previewHtml !== '') {
      return
    }

    setPreviewing(true)
    setPreviewError('')
    try {
      const response = await fetch(routes.adminDocumentationPagePreview.path(page.id), {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ bodyMarkdown: form.data.bodyMarkdown }),
      })
      if (!response.ok) {
        throw new Error('preview failed')
      }
      const payload = (await response.json()) as { html?: string; headings?: DocumentationHeadingData[] }
      setPreviewHtml(payload.html ?? '')
      setPreviewHeadings(payload.headings ?? [])
      setPreviewSource(form.data.bodyMarkdown)
    } catch {
      setPreviewHtml('')
      setPreviewHeadings([])
      setPreviewError('Could not render this draft.')
    } finally {
      setPreviewing(false)
    }
  }

  const editorChrome = (
    <div className="flex flex-wrap items-center justify-between gap-3">
      <div className="flex flex-wrap items-center gap-2">
        <ModeSwitcher mode={pageMode} onWrite={() => setPageMode('write')} onPreview={() => void showPreview()} />
        <Badge variant="outline">{draft.status}</Badge>
        {published ? <Badge variant="outline">Live</Badge> : <Badge variant="outline">Unpublished</Badge>}
      </div>
      <div className="flex flex-wrap gap-2">
        <Sheet>
          <SheetTrigger render={<Button variant="outline" />}>Revisions</SheetTrigger>
          <SheetContent className="overflow-y-auto">
            <SheetHeader>
              <SheetTitle>Revision history</SheetTitle>
              <SheetDescription>Open an old revision as a new draft. History is never mutated.</SheetDescription>
            </SheetHeader>
            <div className="space-y-3 px-4 pb-4">
              {revisions.map((revision) => (
                <div key={revision.id} className="space-y-2 border border-border p-3">
                  <div className="flex items-center justify-between gap-2">
                    <p className="text-sm font-medium">{revision.title}</p>
                    <Badge variant="outline">{revision.status}</Badge>
                  </div>
                  <p className="text-xs text-muted-foreground">{revision.createdAt}</p>
                  {revision.status !== 'draft' ? (
                    <Button size="sm" variant="outline" onClick={() => restoreRevision(revision.id)}>
                      Open as draft
                    </Button>
                  ) : null}
                </div>
              ))}
            </div>
          </SheetContent>
        </Sheet>
        <Button variant="outline" onClick={publishPage} disabled={form.processing}>
          Release page
        </Button>
      </div>
    </div>
  )

  return (
    <DashboardLayout
      title={draft.title || page.slug}
      contentClassName={pageMode === 'preview' ? 'p-0' : 'p-6'}
    >
      {pageMode === 'preview' ? (
        <div className="flex min-h-full flex-col">
          <div className="border-b border-border px-4 py-3">{editorChrome}</div>
          {previewing ? (
            <div className="flex-1 animate-pulse bg-muted/20" />
          ) : previewError ? (
            <p className={`py-8 text-sm text-destructive ${docsGutter}`}>{previewError}</p>
          ) : (
            <div className={`min-w-0 flex-1 items-start py-6 sm:py-8 ${docsColumns} ${docsGutter}`}>
              <div aria-hidden="true" />
              <div id={PREVIEW_ARTICLE_ID} className="min-w-0">
                <Breadcrumb className="mb-5 min-w-0 sm:mb-6">
                  <BreadcrumbList className="gap-1 text-[11px] tracking-wide sm:gap-1.5">
                    <BreadcrumbItem>
                      <span className="text-muted-foreground">{version.label}</span>
                    </BreadcrumbItem>
                    <BreadcrumbSeparator className="text-muted-foreground/50" />
                    <BreadcrumbItem className="min-w-0">
                      <BreadcrumbPage className="truncate font-medium text-foreground/90">
                        {form.data.title || page.slug}
                      </BreadcrumbPage>
                    </BreadcrumbItem>
                  </BreadcrumbList>
                </Breadcrumb>
                <article className="docs-content min-w-0" dangerouslySetInnerHTML={{ __html: previewHtml }} />
              </div>
              <aside className="sticky top-20 hidden w-56 min-w-0 self-start xl:block">
                <DocsToc
                  rootId={PREVIEW_ARTICLE_ID}
                  pageKey={`${page.id}:${previewSource.length}`}
                  headings={previewHeadings}
                />
              </aside>
            </div>
          )}
        </div>
      ) : (
        <div className="mx-auto max-w-5xl">
          <AdminBreadcrumb
            items={[
              { label: 'Docs', href: routes.adminDocumentationVersionIndex.path() },
              { label: version.label, href: routes.adminDocumentationVersionShow.path(version.id) },
              { label: draft.title || page.slug },
            ]}
          />
          <div className="mb-6">{editorChrome}</div>
        <form className="space-y-6" onSubmit={(event) => saveDraft(event.nativeEvent as SubmitEvent)}>
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Page</CardTitle>
                <CardDescription>Title and URL slug. The slug follows the title until you customize it.</CardDescription>
              </CardHeader>
              <CardContent>
                <FieldGroup>
                  <Field data-invalid={Boolean(form.errors.title)}>
                    <FieldLabel htmlFor="title">Title</FieldLabel>
                    <Input
                      id="title"
                      value={form.data.title}
                      onChange={(event) => {
                        const title = event.currentTarget.value
                        form.setData((data) => ({
                          ...data,
                          title,
                          slug: slugIsAutomatic ? slugify(title) : data.slug,
                        }))
                        form.clearErrors('title')
                        if (slugIsAutomatic) {
                          form.clearErrors('slug')
                        }
                      }}
                      placeholder="Getting started"
                      aria-invalid={Boolean(form.errors.title)}
                    />
                    <FieldError>{form.errors.title}</FieldError>
                  </Field>
                  <Field data-invalid={Boolean(form.errors.slug)}>
                    <FieldLabel htmlFor="slug">Slug</FieldLabel>
                    <Input
                      id="slug"
                      value={form.data.slug}
                      onChange={(event) => {
                        const slug = event.currentTarget.value
                        form.setData('slug', slug)
                        form.clearErrors('slug')
                        setSlugIsAutomatic(slug === slugify(form.data.title))
                      }}
                      placeholder="getting-started"
                      aria-invalid={Boolean(form.errors.slug)}
                    />
                    <FieldDescription>Generated from the title until you customize it.</FieldDescription>
                    <FieldError>{form.errors.slug}</FieldError>
                  </Field>
                </FieldGroup>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Metadata</CardTitle>
                <CardDescription>Used for search results when this page is public.</CardDescription>
              </CardHeader>
              <CardContent>
                <FieldGroup>
                  <Field data-invalid={Boolean(form.errors.metaTitle)}>
                    <FieldLabel htmlFor="metaTitle">Meta title</FieldLabel>
                    <Input
                      id="metaTitle"
                      value={form.data.metaTitle}
                      onChange={(event) => {
                        form.setData('metaTitle', event.currentTarget.value)
                        form.clearErrors('metaTitle')
                      }}
                      placeholder="Defaults to the page title"
                      aria-invalid={Boolean(form.errors.metaTitle)}
                    />
                    <LengthHint value={form.data.metaTitle} recommendedMin={50} recommendedMax={60} />
                    <FieldError>{form.errors.metaTitle}</FieldError>
                  </Field>
                  <Field data-invalid={Boolean(form.errors.description)}>
                    <FieldLabel htmlFor="description">Meta description</FieldLabel>
                    <Textarea
                      id="description"
                      value={form.data.description}
                      onChange={(event) => {
                        form.setData('description', event.currentTarget.value)
                        form.clearErrors('description')
                      }}
                      placeholder="A short summary shown under the title in search results"
                      className="min-h-24 resize-y"
                      aria-invalid={Boolean(form.errors.description)}
                    />
                    <LengthHint value={form.data.description} recommendedMin={70} recommendedMax={160} />
                    <FieldError>{form.errors.description}</FieldError>
                  </Field>
                </FieldGroup>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Body</CardTitle>
            </CardHeader>
            <CardContent>
              <MarkdownEditor
                value={form.data.bodyMarkdown}
                onChange={(markdown) => form.setData('bodyMarkdown', markdown)}
                disabled={form.processing}
                invalid={Boolean(form.errors.bodyMarkdown)}
              />
              <FieldError>{form.errors.bodyMarkdown}</FieldError>
            </CardContent>
          </Card>

          <Button type="submit" disabled={form.processing}>
            {form.processing ? 'Saving…' : 'Save draft'}
          </Button>
        </form>
        </div>
      )}
    </DashboardLayout>
  )
}
