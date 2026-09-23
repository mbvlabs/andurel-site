import { Link, useForm } from '@inertiajs/react'

import Layout from '@/Layouts/Layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'

type ResetPasswordRequestProps = {
  errors?: Record<string, string>
}

export default function ResetPasswordRequest({ errors = {} }: ResetPasswordRequestProps) {
  const form = useForm({ email: '' })

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.post(routes.passwordCreate())
  }

  return (
    <Layout
      title="Reset password"
      description="Request a password reset for your Andurel account."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">Reset Password</CardTitle>
          <CardDescription>
            Enter your email address and we&apos;ll send you a code to reset your password.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={(event) => submit(event.nativeEvent as SubmitEvent)}>
            <FieldGroup>
              <Field data-invalid={errors.email ? true : undefined}>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  id="email"
                  type="email"
                  value={form.data.email}
                  onChange={(event) => form.setData('email', event.target.value)}
                  aria-invalid={errors.email ? true : undefined}
                  required
                />
                <FieldError>{errors.email}</FieldError>
              </Field>
              <Button type="submit" className="w-full" size="lg" disabled={form.processing}>
                {form.processing ? 'Loading' : 'Send Reset Code'}
              </Button>
            </FieldGroup>
          </form>
          <p className="mt-6 text-center text-sm text-muted-foreground">
            Remember your password?{' '}
            <Link className="text-foreground/80 hover:text-foreground hover:underline" href={routes.sessionNew()}>
              Login
            </Link>
          </p>
        </CardContent>
      </Card>
    </Layout>
  )
}
