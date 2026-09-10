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
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { routes } from '@/routes'

const INSTALL_COMMAND = 'go install github.com/mbvlabs/andurel@latest'

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
      <div className="mx-auto w-full max-w-7xl px-6 py-4">
        <section className="flex justify-between">
          <div className="flex flex-col">
            <h1 className="text-4xl font-semibold text-card-foreground sm:text-5xl lg:text-5xl lg:leading-[1.1]">
              Space-grade Go framework <br /> For humans and <span className="underline italic decoration-accent">agents</span>
            </h1>
            <p className="mt-5 max-w-lg text-lg leading-7 text-neutral-foreground">
              Everything you and your agent(s) need to build robust and performant applications,
              that will scale to the far-side of the moon.
            </p>
            <div className="flex justify-between gap-x-4 mt-8 max-w-xl">
              <Button
                size="xl"
                className="w-1/2 h-12 px-6 text-sm"
                render={<a href={routes.documentationShow('latest', 'introduction')} />}
              >
                Read The Docs
                <ArrowRightIcon data-icon="inline-end" />
              </Button>
              <Button
                size="xl"
                className="w-1/2 h-12 px-6 text-sm"
				variant="outline"
                render={<a href="https://github.com/mbvlabs/andurel" />}
              >
                See The code
              </Button>
            </div>
          </div>
		  <div className="flex-1 flex items-end">
          	<InstallCommand />
          </div>
        </section>
      </div>
    </Layout>
  )
}
