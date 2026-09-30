import { Fragment } from 'react'
import { Link } from '@inertiajs/react'

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'

export type AdminBreadcrumbItem = {
  label: string
  href?: string
}

export function AdminBreadcrumb({ items }: { items: AdminBreadcrumbItem[] }) {
  return (
    <Breadcrumb className="mb-4">
      <BreadcrumbList>
        {items.map((item, index) => (
          <Fragment key={`${item.label}-${index}`}>
            {index > 0 ? <BreadcrumbSeparator className="text-muted-foreground/50" /> : null}
            <BreadcrumbItem className="min-w-0">
              {item.href ? (
                <BreadcrumbLink
                  render={<Link href={item.href} />}
                  className="truncate text-muted-foreground hover:text-foreground"
                >
                  {item.label}
                </BreadcrumbLink>
              ) : (
                <BreadcrumbPage className="truncate">{item.label}</BreadcrumbPage>
              )}
            </BreadcrumbItem>
          </Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  )
}
