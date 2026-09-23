import { Link, useForm } from '@inertiajs/react'

import Layout from '@/Layouts/Layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'

type LoginProps = {
  errors?: Record<string, string>
}

export default function Login({ errors = {} }: LoginProps) {
  const form = useForm({ email: '', password: '' })

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.post(routes.sessionCreate())
  }

  return (
    <Layout
      title="Log in"
      description="Log in to your Andurel account."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">Login to your account</CardTitle>
          <CardDescription>Enter your details below to login to your account</CardDescription>
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
              <Field data-invalid={errors.password ? true : undefined}>
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <Input
                  id="password"
                  type="password"
                  value={form.data.password}
                  onChange={(event) => form.setData('password', event.target.value)}
                  aria-invalid={errors.password ? true : undefined}
                  required
                />
                <FieldError>{errors.password}</FieldError>
              </Field>
              <p className="text-right text-sm">
                <Link className="text-foreground/80 hover:text-foreground hover:underline" href={routes.passwordNew()}>
                  Forgot your password?
                </Link>
              </p>
              <Button type="submit" className="w-full" size="lg" disabled={form.processing}>
                {form.processing ? 'Loading' : 'Login'}
              </Button>
            </FieldGroup>
          </form>
          <p className="mt-6 text-center text-sm text-muted-foreground">
            Don&apos;t have an account?{' '}
            <Link className="text-foreground/80 hover:text-foreground hover:underline" href={routes.registrationNew()}>
              Sign up
            </Link>
          </p>
        </CardContent>
      </Card>
    </Layout>
  )
}
