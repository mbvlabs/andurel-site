import type { ReactNode } from 'react'

import { Head } from '@inertiajs/react'

type AuthLayoutProps = {
  children: ReactNode
  title: string
}

export default function AuthLayout({ children, title }: AuthLayoutProps) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-background px-4">
      <Head title={title} />
      {children}
    </main>
  )
}
