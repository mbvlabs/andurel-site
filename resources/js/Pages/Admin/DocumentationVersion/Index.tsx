import { Link } from '@inertiajs/react'
import { BookOpenIcon } from 'lucide-react'

import DashboardLayout from '@/Layouts/DashboardLayout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { routes } from '@/routes'
import type { DocumentationVersionIndexProps } from '@/types/payloads'

export default function Index({ versions }: DocumentationVersionIndexProps) {
  return (
    <DashboardLayout title="Docs">
      <div className="mb-6 flex items-center justify-between gap-3">
        <p className="max-w-xl text-sm text-muted-foreground">
          Versions, pages, and sidebar drafts. Public docs still serve from disk until this CMS is verified.
        </p>
        <Button render={<Link href={routes.adminDocumentationVersionNew.path()} />}>New version</Button>
      </div>

      {versions.length === 0 ? (
        <Card>
          <CardHeader>
            <CardTitle>No versions yet</CardTitle>
            <CardDescription>Create head, 1.5.5, and 1.5.2 here, then port Markdown into each page.</CardDescription>
          </CardHeader>
        </Card>
      ) : (
        <div className="grid gap-3">
          {versions.map((version) => (
            <Card key={version.id}>
              <CardHeader className="flex flex-row items-start justify-between gap-3">
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <BookOpenIcon className="size-4" />
                    {version.label}
                  </CardTitle>
                  <CardDescription>
                    /docs/{version.slug}
                    {version.publishedNavId ? ' · nav released' : ' · nav unpublished'}
                  </CardDescription>
                </div>
                <div className="flex items-center gap-2">
                  {version.isLatest ? <Badge>Latest</Badge> : null}
                  <Button
                    variant="outline"
                    size="sm"
                    render={<Link href={routes.adminDocumentationVersionShow.path(version.id)} />}
                  >
                    Open
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="text-xs text-muted-foreground">Position {version.position}</CardContent>
            </Card>
          ))}
        </div>
      )}
    </DashboardLayout>
  )
}
