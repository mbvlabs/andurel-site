import { router, useForm } from '@inertiajs/react'
import { useState } from 'react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { AdminBreadcrumb } from '@/components/docs/AdminBreadcrumb'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'
import type { ApiTokenItemProps, UpdateApiTokenFormPayload } from '@/types/payloads'

function generateToken() {
  const bytes = new Uint8Array(24)
  crypto.getRandomValues(bytes)
  return Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('')
}

function toDatetimeLocal(value: string) {
  if (!value) {
    return ''
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export default function Show({ token, revealedToken = '' }: ApiTokenItemProps) {
  const form = useForm<UpdateApiTokenFormPayload>({
    name: token.name,
    token: '',
    expiresAt: toDatetimeLocal(token.expiresAt),
  })
  const [copied, setCopied] = useState(false)

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.put(routes.adminApiTokenUpdate.path(token.id))
  }

  async function copyRevealed() {
    if (!revealedToken) {
      return
    }
    await navigator.clipboard.writeText(revealedToken)
    setCopied(true)
  }

  return (
    <DashboardLayout title={token.name}>
      <AdminBreadcrumb
        items={[
          { label: 'API tokens', href: routes.adminApiTokenIndex.path() },
          { label: token.name },
        ]}
      />

      {revealedToken ? (
        <Card className="mb-6 max-w-xl border-primary/40">
          <CardHeader>
            <CardTitle>Copy this token now</CardTitle>
            <CardDescription>The secret is hashed after this. It will not be shown again.</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-wrap items-center gap-2">
            <code className="min-w-0 flex-1 break-all font-mono text-xs">{revealedToken}</code>
            <Button type="button" variant="outline" onClick={() => void copyRevealed()}>
              {copied ? 'Copied' : 'Copy'}
            </Button>
          </CardContent>
        </Card>
      ) : null}

      <Card className="max-w-xl">
        <CardHeader>
          <CardTitle>Configure token</CardTitle>
          <CardDescription>
            Prefix {token.tokenPrefix}… · Last used {token.lastUsedAt ? token.lastUsedAt.replace('T', ' ') : 'never'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-6" onSubmit={(event) => submit(event.nativeEvent as SubmitEvent)}>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="name">Name</FieldLabel>
                <Input
                  id="name"
                  value={form.data.name}
                  onChange={(event) => form.setData('name', event.currentTarget.value)}
                  required
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="token">New token value</FieldLabel>
                <div className="flex gap-2">
                  <Input
                    id="token"
                    value={form.data.token}
                    onChange={(event) => form.setData('token', event.currentTarget.value)}
                    placeholder="Leave blank to keep the current secret"
                    className="font-mono"
                  />
                  <Button type="button" variant="outline" onClick={() => form.setData('token', generateToken())}>
                    Generate
                  </Button>
                </div>
                <FieldDescription>Setting a value rotates the secret. Authorization: Bearer &lt;token&gt;</FieldDescription>
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
              <div className="flex flex-wrap gap-2">
                <Button type="submit" disabled={form.processing}>
                  {form.processing ? 'Saving…' : 'Save token'}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    if (window.confirm('Delete this token? API calls using it will fail immediately.')) {
                      router.delete(routes.adminApiTokenDestroy.path(token.id))
                    }
                  }}
                >
                  Delete
                </Button>
              </div>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </DashboardLayout>
  )
}
