import { Link } from '@inertiajs/react'
import { KeyRoundIcon } from 'lucide-react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { AdminBreadcrumb } from '@/components/docs/AdminBreadcrumb'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { routes } from '@/routes'
import type { ApiTokenIndexProps } from '@/types/payloads'

function formatStamp(value: string) {
  if (!value) {
    return '—'
  }
  return value.replace('T', ' ').replace('Z', ' UTC')
}

export default function Index({ tokens }: ApiTokenIndexProps) {
  return (
    <DashboardLayout title="API tokens">
      <AdminBreadcrumb items={[{ label: 'API tokens' }]} />
      <div className="mb-6 flex items-center justify-between gap-3">
        <p className="max-w-xl text-sm text-muted-foreground">
          Bearer tokens for the documentation import API. The secret is stored hashed; you can set the value and timeout here.
        </p>
        <Button render={<Link href={routes.adminApiTokenNew.path()} />}>New token</Button>
      </div>

      {tokens.length === 0 ? (
        <Card>
          <CardHeader>
            <CardTitle>No tokens yet</CardTitle>
            <CardDescription>Create one, then call the API with Authorization: Bearer &lt;token&gt;.</CardDescription>
          </CardHeader>
        </Card>
      ) : (
        <div className="grid gap-3">
          {tokens.map((token) => (
            <Card key={token.id}>
              <CardHeader className="flex flex-row items-start justify-between gap-3">
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <KeyRoundIcon className="size-4" />
                    {token.name}
                  </CardTitle>
                  <CardDescription className="font-mono">{token.tokenPrefix}…</CardDescription>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant="outline">Expires {formatStamp(token.expiresAt)}</Badge>
                  <Button variant="outline" size="sm" render={<Link href={routes.adminApiTokenShow.path(token.id)} />}>
                    Open
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="text-xs text-muted-foreground">
                Created {formatStamp(token.createdAt)}
                {token.lastUsedAt ? ` · Last used ${formatStamp(token.lastUsedAt)}` : ' · Never used'}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </DashboardLayout>
  )
}
