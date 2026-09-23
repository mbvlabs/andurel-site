import { usePage } from '@inertiajs/react'

import type { SharedPageProps } from '@/types/page'

export const siteName = 'Andurel'
export const siteTagline = 'Space-grade Go framework for humans and agents'
export const siteDescription =
  'Andurel is the web development framework for Go. Everything you and your agents need to build robust, performant applications.'
export const twitterSite = '@mbvlabs'
export const defaultOgImage = 'https://media.andurel.com/andurel-og.png'
export const defaultOgAlt = 'Andurel wordmark'
export const githubURL = 'https://github.com/mbvlabs/andurel'
export const organizationURL = 'https://mbvlabs.com'
export const xURL = 'https://x.com/mbvlabs'

export const homeTitle = `${siteName} · ${siteTagline}`

export function useCanonicalUrl() {
  const page = usePage<SharedPageProps>()
  return absoluteUrl(page.url, page.props.appUrl)
}

export function absoluteUrl(path: string, base = '') {
  if (/^https?:\/\//i.test(path)) {
    return path
  }
  const origin = (base || (typeof window !== 'undefined' ? window.location.origin : '')).replace(/\/$/, '')
  if (!origin) {
    return path
  }
  return `${origin}${path.startsWith('/') ? path : `/${path}`}`
}

export function documentTitle(title: string, suffix = siteName) {
  if (!suffix || title === suffix || title.endsWith(` · ${suffix}`) || title.startsWith(`${suffix} ·`)) {
    return title
  }
  return `${title} · ${suffix}`
}

function organizationNode() {
  return {
    '@type': 'Organization',
    '@id': `${organizationURL}#organization`,
    name: 'MBV Labs',
    url: organizationURL,
    sameAs: [githubURL, xURL],
  }
}

export function homeJsonLd(canonical: string) {
  const org = organizationNode()
  return {
    '@context': 'https://schema.org',
    '@graph': [
      org,
      {
        '@type': 'WebSite',
        '@id': `${canonical}#website`,
        name: siteName,
        url: canonical,
        description: siteDescription,
        publisher: { '@id': org['@id'] },
        inLanguage: 'en-US',
      },
      {
        '@type': 'SoftwareApplication',
        name: siteName,
        applicationCategory: 'DeveloperApplication',
        operatingSystem: 'Linux, macOS',
        programmingLanguage: 'Go',
        url: canonical,
        description: siteDescription,
        downloadUrl: githubURL,
        author: { '@id': org['@id'] },
        publisher: { '@id': org['@id'] },
        offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' },
      },
    ],
  }
}

export function docsJsonLd({
  canonical,
  title,
  description,
  version,
  section,
  appUrl,
}: {
  canonical: string
  title: string
  description: string
  version: string
  section: string
  appUrl: string
}) {
  const org = organizationNode()
  const elements: Array<Record<string, unknown>> = [
    { '@type': 'ListItem', position: 1, name: siteName, item: absoluteUrl('/', appUrl) },
    {
      '@type': 'ListItem',
      position: 2,
      name: 'Docs',
      item: absoluteUrl('/docs/latest/introduction', appUrl),
    },
  ]
  if (version) {
    elements.push({
      '@type': 'ListItem',
      position: 3,
      name: version,
      item: absoluteUrl(`/docs/${version}`, appUrl),
    })
  }
  if (section) {
    elements.push({ '@type': 'ListItem', position: elements.length + 1, name: section })
  }
  elements.push({
    '@type': 'ListItem',
    position: elements.length + 1,
    name: title,
    item: canonical,
  })

  return {
    '@context': 'https://schema.org',
    '@graph': [
      org,
      {
        '@type': 'TechArticle',
        headline: title,
        description,
        url: canonical,
        inLanguage: 'en-US',
        isAccessibleForFree: true,
        author: { '@id': org['@id'] },
        publisher: { '@id': org['@id'] },
        mainEntityOfPage: canonical,
      },
      { '@type': 'BreadcrumbList', itemListElement: elements },
    ],
  }
}
