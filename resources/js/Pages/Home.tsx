import { Link } from '@inertiajs/react'

import Layout from '@/Layouts/Layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { routes } from '@/routes'

const features = [
  {
    title: 'Documentation',
    description: 'Start with the framework guides and learn the conventions that shape an Andurel app.',
    href: 'https://andurel.com',
  },
  {
    title: 'Authentication',
    description: 'Registration, sessions, email confirmation, and password reset are already wired.',
    href: routes.registrationNew(),
  },
  {
    title: 'Templ and Inertia',
    description: 'Render server-side Templ pages or scaffold Inertia pages when Inertia is enabled.',
    href: routes.homePage(),
  },
  {
    title: 'Command line',
    description: 'Generate models, factories, controllers, routes, jobs, emails, and views from the CLI.',
    href: 'https://andurel.com',
  },
]

const consoleLines = [
  { label: 'auth', value: 'sessions, registration, reset password' },
  { label: 'views', value: 'Templ + Inertia scaffold online' },
  { label: 'routes', value: 'slugged endpoints mounted' },
  { label: 'jobs', value: 'queue worker ready' },
]

export default function Home() {
  return (
    <Layout>
      <div className="mx-auto w-full max-w-[960px] px-4 py-4">
        <section className="grid items-center gap-6 py-6 lg:grid-cols-[minmax(0,1fr)_22rem]">
          <div className="max-w-3xl">
            <h1 className="text-4xl font-semibold text-card-foreground sm:text-5xl lg:text-6xl">
              Space-grade Go, wired locally.
            </h1>
            <p className="mt-5 max-w-2xl text-lg leading-7 text-muted-foreground">
              Andurel has generated the core app shell: routing, controllers,
              validation, authentication, email, queues, Templ, and
              Inertia-ready frontends.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Button render={<a href={routes.documentationShow("latest", "introduction")} />}>
                Read the docs
              </Button>
              <Button variant="secondary" render={<Link href={routes.registrationNew()} />}>
                Create ACCOUNT
              </Button>
            </div>
          </div>

          <Card className="border-border shadow-sm shadow-black/40">
            <CardHeader className="border-b border-border pb-3">
              <CardTitle className="flex items-center justify-between text-xs uppercase text-text-muted">
                <span>deploy console</span>
                <Badge variant="outline" className="border-accent text-accent">ready</Badge>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4 font-mono text-sm">
              {consoleLines.map((line) => (
                <div key={line.label} className="flex items-start gap-3">
                  <span className="mt-1 size-1.5 bg-accent" />
                  <div>
                    <p className="text-accent">{line.label}</p>
                    <p className="text-text-muted">{line.value}</p>
                  </div>
                </div>
              ))}
            </CardContent>
          </Card>

          <section className="grid gap-4 md:grid-cols-2 lg:col-span-2">
            {features.map((feature) => (
              <a
                key={feature.title}
                className="group block border border-border transition hover:border-ring"
                href={feature.href}
              >
                <Card className="h-full border-0 ring-0">
                  <CardHeader>
                    <div className="mb-2 flex size-8 items-center justify-center border border-ring bg-background text-accent transition group-hover:border-accent">
                      <span className="size-2 bg-accent" />
                    </div>
                    <CardTitle className="text-base">{feature.title}</CardTitle>
                    <CardDescription className="text-sm leading-5">
                      {feature.description}
                    </CardDescription>
                  </CardHeader>
                </Card>
              </a>
            ))}
          </section>
        </section>
      </div>
    </Layout>
  )
}
