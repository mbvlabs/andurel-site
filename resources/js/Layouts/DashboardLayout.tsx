import type { ReactNode } from 'react'

import { Head, Link, router, usePage } from '@inertiajs/react'
import { BookOpenIcon, KeyRoundIcon, LayoutDashboardIcon, LogOutIcon } from 'lucide-react'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { routes } from '@/routes'

type DashboardLayoutProps = {
  children: ReactNode
  title: string
  contentClassName?: string
}

export default function DashboardLayout({
  children,
  title,
  contentClassName = 'p-6',
}: DashboardLayoutProps) {
  const currentPath = usePage().url.split('?')[0]
  const overviewActive = currentPath === '/'
  const docsActive = currentPath === '/docs' || currentPath.startsWith('/docs/')
  const tokensActive = currentPath === '/tokens' || currentPath.startsWith('/tokens/')

  return (
    <SidebarProvider>
      <Head title={title} />
      <Sidebar>
        <SidebarHeader className="px-3 py-3">
          <p className="font-mono text-[0.65rem] uppercase tracking-[0.18em] text-sidebar-foreground/70">
            Andurel admin
          </p>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Console</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    isActive={overviewActive}
                    tooltip="Overview"
                    render={<Link href={routes.adminDashboardHome.path()} />}
                  >
                    <LayoutDashboardIcon />
                    <span>Overview</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    isActive={docsActive}
                    tooltip="Docs"
                    render={<Link href={routes.adminDocumentationVersionIndex.path()} />}
                  >
                    <BookOpenIcon />
                    <span>Docs</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    isActive={tokensActive}
                    tooltip="API tokens"
                    render={<Link href={routes.adminApiTokenIndex.path()} />}
                  >
                    <KeyRoundIcon />
                    <span>API tokens</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                tooltip="Sign out"
                onClick={() => router.delete(routes.sessionDestroy.path())}
              >
                <LogOutIcon />
                <span>Sign out</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>
      <SidebarInset>
        <header className="flex h-14 items-center gap-2 border-b border-border px-4">
          <SidebarTrigger />
          <h1 className="text-sm font-medium">{title}</h1>
        </header>
        <div className={contentClassName}>{children}</div>
      </SidebarInset>
    </SidebarProvider>
  )
}
