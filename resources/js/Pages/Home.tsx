import { useState } from 'react'

import { ArrowRightIcon, CheckIcon, CopyIcon, TerminalIcon } from 'lucide-react'

import Layout from '@/Layouts/Layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupText,
} from '@/components/ui/input-group'
import { Kbd } from '@/components/ui/kbd'
import { Separator } from '@/components/ui/separator'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { routes } from '@/routes'

const INSTALL_COMMAND = 'go install github.com/mbvlabs/andurel@latest'

const telemetry = [
  { label: 'Vehicle', value: 'Andurel v2' },
  { label: 'Propulsion', value: 'Go 1.26' },
  { label: 'Payload', value: 'Inertia v3' },
  { label: 'Adapters', value: 'React / Vue / Svelte' },
  { label: 'Safety', value: 'Nominal' },
]

const hangarBays = [
  {
    bay: '01',
    name: 'Airframe',
    status: 'GO',
    detail: 'Generated models, factories, controllers, routes, and pages that belong to the application.',
  },
  {
    bay: '02',
    name: 'GNC',
    status: 'GO',
    detail: 'Fx for dependency injection and lifecycle. Explicit wiring, not a hidden runtime.',
  },
  {
    bay: '03',
    name: 'Propellant',
    status: 'GO',
    detail: 'PostgreSQL with Bun for ordinary persistence and sqlc for the queries that matter.',
  },
  {
    bay: '04',
    name: 'Payload',
    status: 'GO',
    detail: 'Inertia v3 with React, Vue, or Svelte. Same Go backend, official Inertia adapters, optional SSR.',
  },
  {
    bay: '05',
    name: 'Range ops',
    status: 'GO',
    detail: 'River jobs, email, and queues on the same stack that serves the request path.',
  },
  {
    bay: '06',
    name: 'Ground support',
    status: 'GO',
    detail: 'Agent-ready CLI, JSON discovery, and AGENTS.md so machines can operate the pad without a second stack.',
  },
]

const countdown = [
  {
    mark: 'T-3',
    title: 'Install the CLI',
    command: 'go install github.com/mbvlabs/andurel@latest',
    href: routes.documentationShow('latest', 'installation'),
  },
  {
    mark: 'T-2',
    title: 'Stand up a vehicle',
    command: 'andurel new orbit --inertia react',
    href: routes.documentationShow('latest', 'installation'),
  },
  {
    mark: 'T-1',
    title: 'Arm the agents',
    command: 'andurel skill install',
    href: routes.documentationShow('latest', 'agent-workflows'),
  },
]

const manifest = [
  { system: 'Language', article: 'Go 1.26+', note: 'Required' },
  { system: 'Targets', article: 'Linux & macOS, amd64 / arm64', note: 'Flight range' },
  { system: 'Database', article: 'PostgreSQL · uptrace/bun · sqlc · goose', note: 'Supported' },
  { system: 'UI', article: 'Inertia v3 · React, Vue, Svelte', note: 'Primary' },
  { system: 'Agents', article: 'CLI + AGENTS.md', note: 'On pad' },
  { system: 'Jobs', article: 'River', note: 'On pad' },
]

function SectionLabel({ pad, title }: { pad: string; title: string }) {
  return (
    <div className="mb-6 flex items-end justify-between gap-4">
      <div>
        <p className="font-mono text-[0.65rem] uppercase tracking-[0.22em] text-accent">{pad}</p>
        <h2 className="mt-1 text-2xl font-semibold tracking-tight text-card-foreground">{title}</h2>
      </div>
      <span className="hidden h-px flex-1 bg-border sm:block" />
    </div>
  )
}

function InstallCommand() {
  const [copied, setCopied] = useState(false)

  async function copy() {
    await navigator.clipboard.writeText(INSTALL_COMMAND)
    setCopied(true)
    window.setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Card className="w-[85%] border-border">
      <CardHeader className="border-b border-border pb-3">
        <CardTitle className="flex items-center justify-between gap-6 text-xs uppercase tracking-widest text-muted-foreground">
          <span className="flex items-center gap-2">
            <TerminalIcon className="size-3.5 text-accent" />
            install
          </span>
          <Badge variant="outline" className="border-accent text-accent">
            cli
          </Badge>
        </CardTitle>
        <CardDescription>Get the Andurel CLI on your machine.</CardDescription>
      </CardHeader>
      <CardContent>
        <InputGroup className="h-10 bg-background" aria-label="Install command">
          <InputGroupAddon>
            <InputGroupText>
              <Kbd className="mr-4 bg-transparent px-0 font-mono text-accent">$</Kbd>
            </InputGroupText>
          </InputGroupAddon>
          <InputGroupText className="font-mono text-xs whitespace-nowrap text-card-foreground">
            {INSTALL_COMMAND}
          </InputGroupText>
          <InputGroupAddon className="mx-6" align="inline-end">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger
                  render={
                    <InputGroupButton
                      size="icon-sm"
                      aria-label={copied ? 'Copied' : 'Copy install command'}
                      data-palantir-event="copy-install-command"
                      onClick={copy}
                    />
                  }
                >
                  {copied ? <CheckIcon className="text-accent" /> : <CopyIcon />}
                </TooltipTrigger>
                <TooltipContent>{copied ? 'Copied' : 'Copy'}</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </InputGroupAddon>
        </InputGroup>
      </CardContent>
    </Card>
  )
}

export default function Home() {
  return (
    <Layout>
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-20 px-6 py-4 pb-20">
        <section className="flex min-h-[34rem] flex-col justify-center gap-8 py-20 lg:min-h-[40rem] lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col">
            <p className="mb-3 font-mono text-[0.65rem] uppercase tracking-[0.22em] text-accent">
              Pad A · Clear for launch
            </p>
            <h1 className="text-4xl font-semibold text-card-foreground sm:text-5xl lg:text-5xl lg:leading-[1.1]">
              Space-grade Go framework <br /> For humans and{' '}
              <span className="italic underline decoration-accent">agents</span>
            </h1>
            <p className="mt-5 max-w-lg text-lg leading-7 text-neutral-foreground">
              Everything you and your agent(s) need to build robust and performant applications,
              that will scale to the far-side of the moon.
            </p>
            <div className="mt-8 flex max-w-xl justify-between gap-x-4">
              <Button
                size="lg"
                className="h-12 w-1/2 px-6 text-sm"
                nativeButton={false}
                render={<a href={routes.documentationShow('latest', 'introduction')} />}
              >
                Read The Docs
                <ArrowRightIcon data-icon="inline-end" />
              </Button>
              <Button
                size="lg"
                className="h-12 w-1/2 px-6 text-sm"
                variant="outline"
                nativeButton={false}
                render={<a href="https://github.com/mbvlabs/andurel" />}
              >
                See The code
              </Button>
            </div>
          </div>
          <div className="flex flex-1 items-center justify-end">
            <InstallCommand />
          </div>
        </section>

        <section>
          <div className="grid grid-cols-2 gap-px border border-border bg-border sm:grid-cols-3 lg:grid-cols-5">
            {telemetry.map((item) => (
              <div key={item.label} className="bg-background px-4 py-4">
                <p className="font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
                  {item.label}
                </p>
                <p className="mt-2 flex items-center gap-2 font-mono text-sm text-card-foreground">
                  <span className="size-1.5 bg-accent" />
                  {item.value}
                </p>
              </div>
            ))}
          </div>
        </section>

        <section>
          <SectionLabel pad="Hangar 01" title="Flight hardware on the pad" />
          <div className="grid gap-px border border-border bg-border md:grid-cols-2 lg:grid-cols-3">
            {hangarBays.map((bay) => (
              <article key={bay.bay} className="flex flex-col gap-3 bg-background p-5">
                <div className="flex items-center justify-between gap-3 font-mono text-[0.65rem] uppercase tracking-[0.18em]">
                  <span className="text-muted-foreground">Bay {bay.bay}</span>
                  <span className="text-accent">{bay.status}</span>
                </div>
                <h3 className="text-base font-semibold text-card-foreground">{bay.name}</h3>
                <p className="text-sm leading-6 text-muted-foreground">{bay.detail}</p>
              </article>
            ))}
          </div>
        </section>

        <section>
          <SectionLabel pad="Procedure 02" title="Countdown to first metal" />
          <div className="grid gap-px bg-border lg:grid-cols-3">
            {countdown.map((step) => (
              <a
                key={step.mark}
                href={step.href}
                className="group bg-background p-5 transition-colors hover:bg-card"
              >
                <p className="font-mono text-[0.65rem] uppercase tracking-[0.22em] text-accent">{step.mark}</p>
                <h3 className="mt-3 text-base font-semibold text-card-foreground">{step.title}</h3>
                <p className="mt-3 font-mono text-xs text-muted-foreground group-hover:text-accent">
                  $ {step.command}
                </p>
              </a>
            ))}
          </div>
        </section>

        <section className="grid gap-10 lg:grid-cols-2">
          <div>
            <SectionLabel pad="Crew 03" title="Two consoles, one vehicle" />
            <div className="space-y-6">
              <div className="border border-border p-5">
                <p className="font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
                  Flight crew
                </p>
                <h3 className="mt-2 text-lg font-semibold text-card-foreground">Humans</h3>
                <p className="mt-2 text-sm leading-6 text-muted-foreground">
                  One-time generation. The scaffold writes Go and Inertia pages you can read, edit,
                  and ship. No hidden runtime owning your models.
                </p>
              </div>
              <div className="border border-border p-5">
                <p className="font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
                  Autopilot
                </p>
                <h3 className="mt-2 text-lg font-semibold text-card-foreground">Agents</h3>
                <p className="mt-2 text-sm leading-6 text-muted-foreground">
                  An agent-ready CLI with JSON discovery, plus AGENTS.md in the project. Machines
                  get the same pad as humans, without inventing a second stack.
                </p>
              </div>
            </div>
          </div>
          <div>
            <SectionLabel pad="Manifest 04" title="Vehicle configuration" />
            <Card className="border-border">
              <CardContent className="p-0">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="px-4 font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
                        System
                      </TableHead>
                      <TableHead className="px-4 font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
                        Flight article
                      </TableHead>
                      <TableHead className="px-4 font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
                        Note
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {manifest.map((row) => (
                      <TableRow key={row.system}>
                        <TableCell className="px-4 font-medium text-card-foreground">{row.system}</TableCell>
                        <TableCell className="px-4 font-mono text-muted-foreground">{row.article}</TableCell>
                        <TableCell className="px-4 font-mono text-accent">{row.note}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </div>
        </section>

        <Separator />

        <section className="flex flex-col items-start justify-between gap-6 sm:flex-row sm:items-end">
          <div>
            <p className="font-mono text-[0.65rem] uppercase tracking-[0.22em] text-accent">
              T-0 · Commit to flight
            </p>
            <h2 className="mt-2 text-2xl font-semibold tracking-tight text-card-foreground">
              The pad is hot. Take the vehicle.
            </h2>
            <p className="mt-2 max-w-xl text-sm leading-6 text-muted-foreground">
              Install the CLI, scaffold an Inertia app, and keep the metal you built.
            </p>
          </div>
          <Button
            size="lg"
            className="h-12 px-6 text-sm"
            nativeButton={false}
            render={<a href={routes.documentationShow('latest', 'introduction')} />}
          >
            Read The Docs
            <ArrowRightIcon data-icon="inline-end" />
          </Button>
        </section>
      </div>
    </Layout>
  )
}
