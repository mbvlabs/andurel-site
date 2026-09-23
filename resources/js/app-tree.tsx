import type { ResolvedComponent } from '@inertiajs/react'

import { FlashToasts, pageFlashes } from '@/components/flash-toasts'

type InertiaSetupProps = {
  initialPage?: {
    flash?: unknown
  }
}

export function AppTree({
  App,
  props,
}: {
  App: ResolvedComponent
  props: InertiaSetupProps
}) {
  return (
    <>
      <App {...props} />
      <FlashToasts initialFlashes={pageFlashes(props.initialPage?.flash)} />
    </>
  )
}
