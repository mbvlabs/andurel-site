import Layout from '@/Layouts/Layout'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function NotFound() {
  return (
    <Layout
      title="Page not found"
      description="The page you are looking for could not be found."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90 text-center">
        <CardHeader>
          <p className="text-sm font-medium uppercase tracking-wide text-accent">404</p>
          <CardTitle className="text-2xl font-semibold">Not found</CardTitle>
          <CardDescription className="text-sm leading-6">
            The page you are looking for could not be found.
          </CardDescription>
        </CardHeader>
      </Card>
    </Layout>
  )
}
