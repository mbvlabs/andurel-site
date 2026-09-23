import { useForm } from '@inertiajs/react'

import Layout from '@/Layouts/Layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'

type ResetPasswordProps = {
  token: string
  errors?: Record<string, string>
}

export default function ResetPassword({ token, errors = {} }: ResetPasswordProps) {
  const form = useForm({ resetPasswordToken: token, password: '', confirmPassword: '' })

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.put(routes.passwordUpdate())
  }

  return (
    <Layout
      title="Choose a new password"
      description="Choose a new password for your Andurel account."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">Reset Your Password</CardTitle>
          <CardDescription>Enter your new password below.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={(event) => submit(event.nativeEvent as SubmitEvent)}>
            <FieldGroup>
              {errors.resetPasswordToken ? (
                <FieldError>{errors.resetPasswordToken}</FieldError>
              ) : null}
              <Field data-invalid={errors.password ? true : undefined}>
                <FieldLabel htmlFor="password">New Password</FieldLabel>
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
                <FieldLabel htmlFor="confirmPassword">Confirm New Password</FieldLabel>
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
                {form.processing ? 'Loading' : 'Reset Password'}
              </Button>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </Layout>
  )
}
