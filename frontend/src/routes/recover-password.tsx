import { Alert, Container, Heading, Input, Text } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import { createFileRoute, redirect, useSearch } from "@tanstack/react-router"
import { useEffect, useState } from "react"
import { type SubmitHandler, useForm } from "react-hook-form"
import { FiAlertCircle, FiMail } from "react-icons/fi"

import type { ApiError } from "@/client/user"
import { userServiceRecoverPassword } from "@/client/user"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { emailPattern, handleError } from "@/utils"

interface FormData {
  email: string
}

interface RecoverPasswordSearch {
  error?: string
}

export const Route = createFileRoute("/recover-password")({
  component: RecoverPassword,
  validateSearch: (search: Record<string, unknown>): RecoverPasswordSearch => {
    return {
      error: (search.error as string) || undefined,
    }
  },
  beforeLoad: async () => {
    if (isLoggedIn()) {
      throw redirect({
        to: "/",
      })
    }
  },
})

function RecoverPassword() {
  const search = useSearch({ from: "/recover-password" })
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>()
  const { showSuccessToast } = useCustomToast()

  useEffect(() => {
    if (search.error) {
      setErrorMessage(search.error)
    }
  }, [search.error])

  const recoverPassword = async (data: FormData) => {
    await userServiceRecoverPassword({ requestBody: { email: data.email } })
  }

  const mutation = useMutation({
    mutationFn: recoverPassword,
    onSuccess: () => {
      showSuccessToast("Password recovery email sent successfully.")
      reset()
      setErrorMessage(null)
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const onSubmit: SubmitHandler<FormData> = async (data) => {
    setErrorMessage(null)
    mutation.mutate(data)
  }

  return (
    <Container
      as="form"
      onSubmit={handleSubmit(onSubmit)}
      h="100vh"
      maxW="sm"
      alignItems="stretch"
      justifyContent="center"
      gap={4}
      centerContent
    >
      <Heading size="xl" color="ui.main" textAlign="center" mb={2}>
        Password Recovery
      </Heading>
      <Text textAlign="center">
        A password recovery email will be sent to the registered account.
      </Text>
      {errorMessage && (
        <Alert.Root status="error">
          <Alert.Icon>
            <FiAlertCircle />
          </Alert.Icon>
          <Alert.Description>{errorMessage}</Alert.Description>
        </Alert.Root>
      )}
      <Field invalid={!!errors.email} errorText={errors.email?.message}>
        <InputGroup w="100%" startElement={<FiMail />}>
          <Input
            id="email"
            {...register("email", {
              required: "Email is required",
              pattern: emailPattern,
            })}
            placeholder="Email"
            type="email"
          />
        </InputGroup>
      </Field>
      <Button variant="solid" type="submit" loading={isSubmitting}>
        Continue
      </Button>
    </Container>
  )
}
