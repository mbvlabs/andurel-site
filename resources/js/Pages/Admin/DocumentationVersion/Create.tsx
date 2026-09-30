import { useForm } from '@inertiajs/react'
import { useState } from 'react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { AdminBreadcrumb } from '@/components/docs/AdminBreadcrumb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { slugify } from '@/lib/slug'
import { routes } from '@/routes'
import type { CreateDocumentationVersionFormPayload } from '@/types/payloads'

export default function Create() {
  const form = useForm<CreateDocumentationVersionFormPayload>({
    slug: '',
    label: '',
    isLatest: false,
  })
  const [slugIsAutomatic, setSlugIsAutomatic] = useState(true)

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.post(routes.adminDocumentationVersionCreate.path())
  }

  return (
    <DashboardLayout title="New version">
      <AdminBreadcrumb
        items={[
          { label: 'Docs', href: routes.adminDocumentationVersionIndex.path() },
          { label: 'New version' },
        ]}
      />
      <Card className="max-w-xl">
        <CardHeader>
          <CardTitle>Create a documentation version</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={(event) => submit(event.nativeEvent as SubmitEvent)}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="label">Label</FieldLabel>
                <Input
                  id="label"
                  value={form.data.label}
                  onChange={(event) => {
                    const label = event.currentTarget.value
                    form.setData((data) => ({
                      ...data,
                      label,
                      slug: slugIsAutomatic ? slugify(label) : data.slug,
                    }))
                  }}
                  required
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="slug">Slug</FieldLabel>
                <Input
                  id="slug"
                  value={form.data.slug}
                  onChange={(event) => {
                    const slug = event.currentTarget.value
                    form.setData('slug', slug)
                    setSlugIsAutomatic(slug === slugify(form.data.label))
                  }}
                  required
                />
                <FieldDescription>Generated from the label until you customize it.</FieldDescription>
              </Field>
              <Field>
                <label className="flex items-center gap-2 text-xs">
                  <input
                    type="checkbox"
                    checked={form.data.isLatest}
                    onChange={(event) => form.setData('isLatest', event.currentTarget.checked)}
                  />
                  Mark as latest
                </label>
              </Field>
              <Button type="submit" disabled={form.processing}>
                {form.processing ? 'Creating…' : 'Create version'}
              </Button>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </DashboardLayout>
  )
}
