import type { FormEvent, ReactNode } from 'react'

import AndurelWordmark from '@/components/andurel-wordmark'
import AsciiSky from '@/components/ascii-sky'
import SeoHead from '@/components/seo-head'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { trackEvent } from '@/lib/palantir'
import { siteDescription, siteName } from '@/lib/seo'
import { routes } from '@/routes'

type LayoutProps = {
  children: ReactNode
  title?: string
  description?: string
  robots?: string
  jsonLd?: unknown
}

const socials = [
  { href: 'https://x.com/mbvlabs', label: 'x.com/mbvlabs' },
  { href: 'https://mbvlabs.com', label: 'mbvlabs.com' },
  { href: 'https://mortenvistisen.com', label: 'mortenvistisen.com' },
  { href: 'https://deploycrate.com', label: 'deploycrate.com' },
  { href: 'https://mastergolang.com', label: 'mastergolang.com' },
  { href: 'https://linkedin.com/in/mortenvistisen', label: 'linkedin.com/in/mortenvistisen' },
  { href: 'https://youtube.com/@mbvlabs', label: 'youtube.com/@mbvlabs' },
  { href: 'https://twitch.tv/mbvlabs', label: 'twitch.tv/mbvlabs' },
]

export default function Layout({
  children,
  title = siteName,
  description = siteDescription,
  robots,
  jsonLd,
}: LayoutProps) {
  return (
    <main className="relative flex min-h-screen flex-col overflow-x-clip bg-[#090b0d] text-[#e4dfd2]">
      <SeoHead title={title} description={description} robots={robots} jsonLd={jsonLd} />
      <AsciiSky />
      <header className="relative">
        <div className="mx-auto flex w-full max-w-7xl items-center justify-between px-6 py-3">
          <a className="inline-flex items-center gap-3 text-sm font-semibold text-[#f2ead8]" href={routes.homePage()}>
            <span className="grid size-8 grid-cols-2 gap-1 border border-[#52605c] bg-[#101414] p-1 shadow-sm shadow-black/40">
              <span className="border border-[#8df7a4]" />
              <span className="border border-[#52605c]" />
              <span className="border border-[#52605c]" />
              <span className="bg-[#8df7a4]" />
            </span>
            <span>Andurel.</span>
          </a>
          <nav className="flex flex-wrap items-center justify-end gap-3 text-sm">
            <a className="px-2 py-1 text-[#aaa393] transition hover:text-[#f2ead8]" href="https://github.com/mbvlabs/andurel">GitHub</a>
            <a className="px-2 py-1 text-[#aaa393] transition hover:text-[#f2ead8]" href={routes.documentationShow('latest', 'introduction')}>Docs</a>
          </nav>
        </div>
      </header>
      <div className="relative flex flex-1 flex-col py-6">{children}</div>
      <footer className="relative border-t border-[#2f3a37] py-8">
        <div className="mx-auto flex w-full max-w-7xl flex-col gap-8 px-6">
          <div className="flex flex-col gap-10 sm:flex-row sm:items-start sm:justify-between">
            <div className="flex w-full max-w-sm flex-col gap-3">
              <p className="font-mono text-xs uppercase tracking-widest text-[#8f8a7d]">Pad Andurel</p>
              <form
                className="flex flex-col gap-2"
                onSubmit={(event: FormEvent<HTMLFormElement>) => {
                  event.preventDefault()
                  const form = new FormData(event.currentTarget)
                  trackEvent('newsletter-signup', {
                    firstName: String(form.get('firstName') ?? '').trim(),
                    lastName: String(form.get('lastName') ?? '').trim(),
                    email: String(form.get('email') ?? '').trim(),
                  })
                }}
              >
                <div className="grid grid-cols-2 gap-2">
                  <Input
                    id="footer-first-name"
                    name="firstName"
                    autoComplete="given-name"
                    aria-label="First name"
                    placeholder="First name"
                    className="h-9 border-[#2f3a37] bg-[#090c0d] px-3 text-sm text-[#e4dfd2] placeholder:text-[#8f8a7d] focus-visible:border-[#8df7a4] focus-visible:ring-[#8df7a4]/20"
                  />
                  <Input
                    id="footer-last-name"
                    name="lastName"
                    autoComplete="family-name"
                    aria-label="Last name"
                    placeholder="Last name"
                    className="h-9 border-[#2f3a37] bg-[#090c0d] px-3 text-sm text-[#e4dfd2] placeholder:text-[#8f8a7d] focus-visible:border-[#8df7a4] focus-visible:ring-[#8df7a4]/20"
                  />
                </div>
                <Input
                  id="footer-email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  aria-label="Email address"
                  placeholder="Email address"
                  className="h-9 border-[#2f3a37] bg-[#090c0d] px-3 text-sm text-[#e4dfd2] placeholder:text-[#8f8a7d] focus-visible:border-[#8df7a4] focus-visible:ring-[#8df7a4]/20"
                />
                <Button type="submit" variant="outline" className="h-9 w-full border-[#2f3a37] text-sm text-[#f2ead8] hover:bg-[#101414] hover:text-[#f2ead8]">
                  Stay updated
                </Button>
                <p className="text-xs leading-5 text-[#8f8a7d]">
                  By submitting this form, you agree to our <span className="underline">terms</span>. You can
                  opt-out anytime.
                </p>
              </form>
            </div>
            <div className="flex flex-col gap-3">
              <p className="font-mono text-xs uppercase tracking-widest text-[#8f8a7d]">Range nominal</p>
              <nav className="grid grid-cols-2 gap-x-8 gap-y-2 font-mono text-xs text-[#8f8a7d]">
                {socials.map((social) => (
                  <a
                    key={social.href}
                    href={social.href}
                    target="_blank"
                    className="transition hover:text-[#f2ead8]"
                  >
                    {social.label}
                  </a>
                ))}
              </nav>
            </div>
          </div>
          <AndurelWordmark />
          <p className="text-center font-mono text-xs uppercase tracking-widest text-[#8f8a7d]">
            &copy; {new Date().getFullYear()}
          </p>
        </div>
      </footer>
    </main>
  )
}
