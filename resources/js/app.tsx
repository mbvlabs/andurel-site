import '../../css/base.css'

import { createRoot, hydrateRoot } from 'react-dom/client'
import { createInertiaApp, type ResolvedComponent } from '@inertiajs/react'

import { FlashToasts, pageFlashes } from '@/components/flash-toasts'

type PageProps = {
  [key: string]: unknown
}

type PageModule = {
  default: ResolvedComponent
}

type InertiaSetupProps = {
  initialPage?: {
    flash?: unknown
  }
}

function AppTree({
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

createInertiaApp<PageProps>({
  resolve: (name: string) => {
    const pages = import.meta.glob<PageModule>('./Pages/**/*.tsx', { eager: true })
    return pages[`./Pages/${name}.tsx`].default
  },
  setup({ el, App, props }) {
    const app = <AppTree App={App} props={props} />
    if (el.dataset.serverRendered === 'true') {
      hydrateRoot(el, app)
    } else {
      createRoot(el).render(app)
    }
  },
})
