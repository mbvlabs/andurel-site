import { Link, useForm } from '@inertiajs/react'

import Layout from '@/Layouts/Layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'

type RegistrationProps = {
  errors?: Record<string, string>
}

export default function Registration({ errors = {} }: RegistrationProps) {
  const form = useForm({ email: '', password: '', confirmPassword: '' })

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.post(routes.registrationCreate())
  }

  return (
    <Layout
      title="Create an account"
      description="Create an Andurel account."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">Create an account</CardTitle>
          <CardDescription>Enter your details below to create your account</CardDescription>
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
              <Field data-invalid={errors.confirmPassword ? true : undefined}>
                <FieldLabel htmlFor="confirmPassword">Confirm Password</FieldLabel>
                <Input
                  id="confirmPassword"
                  type="password"
                  value={form.data.confirmPassword}
                  onChange={(event) => form.setData('confirmPassword', event.target.value)}
                  aria-invalid={errors.confirmPassword ? true : undefined}
                  required
                />
                <FieldError>{errors.confirmPassword}</FieldError>
              </Field>
              <Button type="submit" className="w-full" size="lg" disabled={form.processing}>
                {form.processing ? 'Loading' : 'Sign Up'}
              </Button>
            </FieldGroup>
          </form>
          <p className="mt-6 text-center text-sm text-muted-foreground">
            Already have an account?{' '}
            <Link className="text-foreground/80 hover:text-foreground hover:underline" href={routes.sessionNew()}>
              Login
            </Link>
          </p>
        </CardContent>
      </Card>
    </Layout>
  )
}
