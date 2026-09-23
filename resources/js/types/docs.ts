export type DocHeading = {
  id: string
  text: string
  level: number
}

export type DocLink = {
  title: string
  url: string
}

export type DocPage = {
  slug: string
  title: string
  description: string
  url: string
  children?: DocPage[]
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
