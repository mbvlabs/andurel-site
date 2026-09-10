export type PalantirEventData = Record<string, string | number | boolean | null>

declare global {
  interface Window {
    palantir?: {
      track: (name: string, data?: PalantirEventData) => void
    }
  }
}

export function trackEvent(name: string, data?: PalantirEventData) {
  window.palantir?.track(name, data)
}
