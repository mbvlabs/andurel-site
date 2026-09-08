export type DocPage = {
  slug: string
  title: string
  description: string
  url: string
}

export type DocSection = {
  title: string
  pages: DocPage[]
}

export type DocVersion = {
  name: string
  url: string
  sections: DocSection[]
}

export type DocNavigationProps = {
  versions: DocVersion[]
  currentVersion: string
  currentSlug: string
}
