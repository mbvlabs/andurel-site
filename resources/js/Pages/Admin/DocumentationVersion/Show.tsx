import { Link, router, useForm } from '@inertiajs/react'
import { useState } from 'react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { AdminBreadcrumb } from '@/components/docs/AdminBreadcrumb'
import { NavTreeEditor } from '@/components/docs/NavTreeEditor'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { slugify } from '@/lib/slug'
import { routes } from '@/routes'
import type {
  CreateDocumentationPageFormPayload,
  DocumentationNavNodeData,
  DocumentationVersionItemProps,
  UpdateDocumentationVersionFormPayload,
} from '@/types/payloads'

export default function Show({ version, pages, navDraft }: DocumentationVersionItemProps) {
  const versionForm = useForm<UpdateDocumentationVersionFormPayload>({
    slug: version.slug,
    label: version.label,
    isLatest: version.isLatest,
    position: version.position,
  })
  const pageForm = useForm<CreateDocumentationPageFormPayload>({
    slug: '',
    title: '',
  })
  const [pageSlugIsAutomatic, setPageSlugIsAutomatic] = useState(true)
  const [tree, setTree] = useState<DocumentationNavNodeData[]>(navDraft.tree)

  function saveVersion(event: SubmitEvent) {
    event.preventDefault()
    versionForm.put(routes.adminDocumentationVersionUpdate.path(version.id))
  }

  function createPage(event: SubmitEvent) {
    event.preventDefault()
    pageForm.post(routes.adminDocumentationVersionCreatePage.path(version.id), {
      onSuccess: () => {
        pageForm.reset()
        setPageSlugIsAutomatic(true)
      },
    })
  }

  function saveNav() {
    router.put(routes.adminDocumentationVersionUpdateNav.path(version.id), { tree })
  }

  function publishNav() {
    router.post(routes.adminDocumentationVersionPublishNav.path(version.id), { tree })
  }

  return (
    <DashboardLayout title={version.label}>
      <AdminBreadcrumb
        items={[
          { label: 'Docs', href: routes.adminDocumentationVersionIndex.path() },
          { label: version.label },
        ]}
      />
      <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          {version.isLatest ? <Badge>Latest</Badge> : null}
          {version.publishedNavId ? <Badge variant="secondary">Nav released</Badge> : <Badge variant="outline">Nav unpublished</Badge>}
        </div>
      </div>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_24rem]">
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Version</CardTitle>
              <CardDescription>Slug, label, and which version /docs/latest should resolve to later.</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={(event) => saveVersion(event.nativeEvent as SubmitEvent)}>
                <FieldGroup>
                  <Field>
                    <FieldLabel htmlFor="label">Label</FieldLabel>
                    <Input
                      id="label"
                      value={versionForm.data.label}
                      onChange={(event) => versionForm.setData('label', event.currentTarget.value)}
                    />
                  </Field>
                  <Field>
                    <FieldLabel htmlFor="slug">Slug</FieldLabel>
                    <Input
                      id="slug"
                      value={versionForm.data.slug}
                      onChange={(event) => versionForm.setData('slug', event.currentTarget.value)}
                    />
                  </Field>
                  <Field>
                    <FieldLabel htmlFor="position">Position</FieldLabel>
                    <Input
                      id="position"
                      type="number"
                      value={versionForm.data.position}
                      onChange={(event) => versionForm.setData('position', Number(event.currentTarget.value))}
                    />
                  </Field>
                  <Field>
                    <label className="flex items-center gap-2 text-xs">
                      <input
                        type="checkbox"
                        checked={versionForm.data.isLatest}
                        onChange={(event) => versionForm.setData('isLatest', event.currentTarget.checked)}
                      />
                      Mark as latest
                    </label>
                  </Field>
                  <Button type="submit" disabled={versionForm.processing}>
                    Save version
                  </Button>
                </FieldGroup>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-start justify-between gap-3">
              <div>
                <CardTitle>Sidebar draft</CardTitle>
                <CardDescription>Reorder and nest pages. Release nav separately from page bodies.</CardDescription>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={saveNav}>
                  Save nav draft
                </Button>
                <Button size="sm" onClick={publishNav}>
                  Release nav
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <NavTreeEditor tree={tree} pages={pages} onChange={setTree} />
            </CardContent>
          </Card>
        </div>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>New page</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={(event) => createPage(event.nativeEvent as SubmitEvent)}>
                <FieldGroup>
                  <Field>
                    <FieldLabel htmlFor="page-title">Title</FieldLabel>
                    <Input
                      id="page-title"
                      value={pageForm.data.title}
                      onChange={(event) => {
                        const title = event.currentTarget.value
                        pageForm.setData((data) => ({
                          ...data,
                          title,
                          slug: pageSlugIsAutomatic ? slugify(title) : data.slug,
                        }))
                      }}
                      required
                    />
                  </Field>
                  <Field>
                    <FieldLabel htmlFor="page-slug">Slug</FieldLabel>
                    <Input
                      id="page-slug"
                      value={pageForm.data.slug}
                      onChange={(event) => {
                        const slug = event.currentTarget.value
                        pageForm.setData('slug', slug)
                        setPageSlugIsAutomatic(slug === slugify(pageForm.data.title))
                      }}
                      required
                    />
                    <FieldDescription>Generated from the title until you customize it.</FieldDescription>
                  </Field>
                  <Button type="submit" disabled={pageForm.processing}>
                    Create page
                  </Button>
                </FieldGroup>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Pages</CardTitle>
              <CardDescription>{pages.length} in this version</CardDescription>
            </CardHeader>
            <CardContent className="space-y-2">
              {pages.length === 0 ? (
                <p className="text-xs text-muted-foreground">No pages yet.</p>
              ) : (
                pages.map((page) => (
                  <div key={page.id} className="flex items-center justify-between gap-3 border border-border px-3 py-2">
                    <div>
                      <p className="text-sm font-medium">{page.title}</p>
                      <p className="text-xs text-muted-foreground">{page.slug}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      {page.isPublished ? <Badge variant="secondary">Live</Badge> : <Badge variant="outline">Unpublished</Badge>}
                      {page.hasDraft ? <Badge variant="outline">Draft</Badge> : null}
                      <Button
                        size="sm"
                        variant="outline"
                        render={<Link href={routes.adminDocumentationPageShow.path(page.id)} />}
                      >
                        Edit
                      </Button>
                    </div>
                  </div>
                ))
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </DashboardLayout>
  )
}
