import { useForm } from '@inertiajs/react'

import Layout from '@/Layouts/Layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'

type ConfirmEmailProps = {
  errors?: Record<string, string>
}

export default function ConfirmEmail({ errors = {} }: ConfirmEmailProps) {
  const form = useForm({ code: '' })

  function submit(event: SubmitEvent) {
    event.preventDefault()
    form.post(routes.confirmationCreate())
  }

  return (
    <Layout
      title="Verify your email"
      description="Verify your email address to finish creating your Andurel account."
      robots="noindex, nofollow"
    >
      <Card className="mx-auto w-full max-w-md border-border bg-card/90">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">Verify Your Email</CardTitle>
          <CardDescription>
            Please enter the 6-digit verification code sent to your email.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={(event) => submit(event.nativeEvent as SubmitEvent)}>
            <FieldGroup>
              <Field data-invalid={errors.code ? true : undefined}>
                <FieldLabel htmlFor="code">Verification Code</FieldLabel>
                <Input
                  id="code"
                  type="text"
                  maxLength={6}
                  value={form.data.code}
                  onChange={(event) => form.setData('code', event.target.value)}
                  className="text-center tracking-[0.3em]"
                  aria-invalid={errors.code ? true : undefined}
                  required
                />
                <FieldError>{errors.code}</FieldError>
              </Field>
              <Button type="submit" className="w-full" size="lg" disabled={form.processing}>
                {form.processing ? 'Loading' : 'Verify Email'}
              </Button>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </Layout>
  )
}
