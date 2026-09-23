import { Head } from '@inertiajs/react'

import {
  defaultOgAlt,
  defaultOgImage,
  documentTitle,
  useCanonicalUrl,
} from '@/lib/seo'

type SeoHeadProps = {
  title: string
  description: string
  type?: 'website' | 'article'
  robots?: string
  suffix?: string
  jsonLd?: unknown
}

export default function SeoHead({
  title,
  description,
  type = 'website',
  robots = 'index, follow',
  suffix,
  jsonLd,
}: SeoHeadProps) {
  const canonical = useCanonicalUrl()
  const fullTitle = documentTitle(title, suffix)

  return (
    <Head title={fullTitle}>
      <meta head-key="robots" name="robots" content={robots} />
      <meta head-key="description" name="description" content={description} />
      <link head-key="canonical" rel="canonical" href={canonical} />
      <meta head-key="og:type" property="og:type" content={type} />
      <meta head-key="og:title" property="og:title" content={fullTitle} />
      <meta head-key="og:description" property="og:description" content={description} />
      <meta head-key="og:url" property="og:url" content={canonical} />
      <meta head-key="og:image" property="og:image" content={defaultOgImage} />
      <meta head-key="og:image:alt" property="og:image:alt" content={defaultOgAlt} />
      <meta head-key="twitter:title" name="twitter:title" content={fullTitle} />
      <meta head-key="twitter:description" name="twitter:description" content={description} />
      <meta head-key="twitter:image" name="twitter:image" content={defaultOgImage} />
      <meta head-key="twitter:image:alt" name="twitter:image:alt" content={defaultOgAlt} />
      {jsonLd ? (
        <script type="application/ld+json" head-key="json-ld">
          {JSON.stringify(jsonLd)}
        </script>
      ) : null}
    </Head>
  )
}
