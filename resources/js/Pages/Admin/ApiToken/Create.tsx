import { useForm } from '@inertiajs/react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { AdminBreadcrumb } from '@/components/docs/AdminBreadcrumb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'
import type { CreateApiTokenFormPayload } from '@/types/payloads'

function generateToken() {
  const bytes = new Uint8Array(24)
  crypto.getRandomValues(bytes)
  return Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('')
}

export default function Create() {
  const form = useForm<CreateApiTokenFormPayload>({
    name: '',
    token: '',
    expiresAt: '',
  })

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.post(routes.adminApiTokenCreate.path())
  }

  return (
    <DashboardLayout title="New API token">
      <AdminBreadcrumb
        items={[
          { label: 'API tokens', href: routes.adminApiTokenIndex.path() },
          { label: 'New token' },
        ]}
      />
      <Card className="max-w-xl">
        <CardHeader>
          <CardTitle>Create an API token</CardTitle>
          <CardDescription>
            Set the secret yourself or generate one. It is hashed on save. Every token needs an expiry.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={(event) => submit(event.nativeEvent as SubmitEvent)}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="name">Name</FieldLabel>
                <Input
                  id="name"
                  value={form.data.name}
                  onChange={(event) => form.setData('name', event.currentTarget.value)}
                  placeholder="Docs importer"
                  required
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="token">Token value</FieldLabel>
                <div className="flex gap-2">
                  <Input
                    id="token"
                    value={form.data.token}
                    onChange={(event) => form.setData('token', event.currentTarget.value)}
                    placeholder="Leave blank to generate"
                    className="font-mono"
                  />
                  <Button type="button" variant="outline" onClick={() => form.setData('token', generateToken())}>
                    Generate
                  </Button>
                </div>
                <FieldDescription>At least 16 characters if you set it. Copy it before saving if you generate it here.</FieldDescription>
              </Field>
              <Field>
                <FieldLabel htmlFor="expiresAt">Timeout</FieldLabel>
                <Input
                  id="expiresAt"
                  type="datetime-local"
                  value={form.data.expiresAt}
                  onChange={(event) => form.setData('expiresAt', event.currentTarget.value)}
                  required
                />
                <FieldDescription>When the token stops working.</FieldDescription>
              </Field>
              <Button type="submit" disabled={form.processing}>
                {form.processing ? 'Creating…' : 'Create token'}
              </Button>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </DashboardLayout>
  )
}
