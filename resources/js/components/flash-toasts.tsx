import { useEffect } from 'react'

import { router } from '@inertiajs/react'
import { toast } from 'sonner'

import { Toaster } from '@/components/ui/sonner'

type FlashMessage = {
  Type: string
  Message: string
}

function isFlashMessage(value: unknown): value is FlashMessage {
  if (!value || typeof value !== 'object') {
    return false
  }

  const flash = value as Record<string, unknown>
  return typeof flash.Type === 'string' && typeof flash.Message === 'string'
}

export function pageFlashes(value: unknown): FlashMessage[] | undefined {
  if (!Array.isArray(value)) {
    return undefined
  }

  return value.filter(isFlashMessage)
}

function pushFlashes(flashes?: FlashMessage[]) {
  if (!flashes || flashes.length === 0) {
    return
  }

  for (const flash of flashes) {
    const type = flash.Type.toLowerCase()
    if (type === 'success') {
      toast.success(flash.Message)
    } else if (type === 'error') {
      toast.error(flash.Message)
    } else if (type === 'warning') {
      toast.warning(flash.Message)
    } else {
      toast(flash.Message)
    }
  }
}

export function FlashToasts({ initialFlashes }: { initialFlashes?: FlashMessage[] }) {
  useEffect(() => {
    pushFlashes(initialFlashes)

    const removeListener = router.on('success', (event) => {
      pushFlashes(pageFlashes(event.detail.page.flash))
    })

    return () => {
      removeListener()
    }
  }, [initialFlashes])

  return <Toaster position="bottom-right" richColors closeButton />
}
