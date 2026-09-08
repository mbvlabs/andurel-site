import { Head } from '@inertiajs/react'

import DocLayout from '@/Layouts/DocLayout'
import type { DocNavigationProps } from '@/types/docs'

type ShowProps = DocNavigationProps & {
  currentSection: string
  title: string
  description: string
}

export default function Show({
  versions,
  currentVersion,
  currentSlug,
  currentSection,
  title,
  description,
}: ShowProps) {
  return (
    <DocLayout
      versions={versions}
      currentVersion={currentVersion}
      currentSlug={currentSlug}
    >
      <Head title={title} />
      <p className="mb-7 font-mono text-xs font-semibold uppercase tracking-[0.14em] text-accent">
        {currentSection} / {currentVersion}
      </p>
      <h1 id={currentSlug} className="scroll-mt-20 text-3xl font-semibold text-foreground">
        {title}
      </h1>
      <p className="mt-4 text-base leading-7 text-muted-foreground">{description}</p>
    </DocLayout>
  )
}
