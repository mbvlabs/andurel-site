import Layout from '@/Layouts/Layout'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function InternalError() {
  return (
    <Layout
      title="Something went wrong"
      description="The application hit an unexpected error."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90 text-center">
        <CardHeader>
          <p className="text-sm font-medium uppercase tracking-wide text-accent">500</p>
          <CardTitle className="text-2xl font-semibold">Something went wrong.</CardTitle>
          <CardDescription className="text-sm leading-6">
            The application hit an unexpected error.
          </CardDescription>
        </CardHeader>
      </Card>
    </Layout>
  )
}
