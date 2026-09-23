import Layout from '@/Layouts/Layout'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function BadRequest() {
  return (
    <Layout
      title="Bad request"
      description="The request made was invalid."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90 text-center">
        <CardHeader>
          <p className="text-sm font-medium uppercase tracking-wide text-accent">400</p>
          <CardTitle className="text-2xl font-semibold">Bad request</CardTitle>
          <CardDescription className="text-sm leading-6">
            The request made was invalid.
          </CardDescription>
        </CardHeader>
      </Card>
    </Layout>
  )
}
