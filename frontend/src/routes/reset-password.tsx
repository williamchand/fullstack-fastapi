import { Container, Heading, Text, Alert } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import { createFileRoute, redirect, useNavigate, useSearch } from "@tanstack/react-router"
import { type SubmitHandler, useForm } from "react-hook-form"
import { FiLock, FiAlertCircle } from "react-icons/fi"
import { useEffect, useState } from "react"

import { type ApiError } from "@/client/user"
import { userServiceResetPassword } from "@/client/user"
import { Button } from "@/components/ui/button"
import { PasswordInput } from "@/components/ui/password-input"
import { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { confirmPasswordRules, handleError, passwordRules } from "@/utils"

interface NewPasswordForm {
  confirm_password: string
  new_password: string
}

interface ResetPasswordSearch {
  token?: string
}

export const Route = createFileRoute("/reset-password")({
  component: ResetPassword,
  validateSearch: (search: Record<string, unknown>): ResetPasswordSearch => {
    return {
      token: (search.token as string) || undefined,
    }
  },
  beforeLoad: async ({ search }) => {
    if (isLoggedIn()) {
      throw redirect({
        to: "/",
      })
    }
    // Validate token exists
    if (!search.token) {
      throw redirect({
        to: "/recover-password",
        search: {
          error: "Invalid or missing reset token. Please request a new password reset link.",
        },
      })
    }
  },
})

function ResetPassword() {
  const search = useSearch({ from: "/reset-password" })
  const [tokenError, setTokenError] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    getValues,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<NewPasswordForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      new_password: "",
      confirm_password: "",
    },
  })
  const { showSuccessToast } = useCustomToast()
  const navigate = useNavigate()

  useEffect(() => {
    // Validate token on mount
    if (!search.token) {
      setTokenError("Invalid or missing reset token. Please request a new password reset link.")
    }
  }, [search.token])

  const resetPassword = async (data: { new_password: string }) => {
    if (!search.token) {
      setTokenError("Invalid or missing reset token.")
      return
    }
    await userServiceResetPassword({ requestBody: { newPassword: data.new_password, token: search.token } })
  }

  const mutation = useMutation({
    mutationFn: resetPassword,
    onSuccess: () => {
      showSuccessToast("Password updated successfully.")
      reset()
      navigate({ to: "/login" })
    },
    onError: (err: ApiError) => {
      const errDetail = (err.body as any)?.detail
      if (errDetail?.includes("Invalid token") || errDetail?.includes("expired")) {
        setTokenError("This reset link is invalid or has expired. Please request a new password reset link.")
      } else {
      handleError(err)
      }
    },
  })

  const onSubmit: SubmitHandler<NewPasswordForm> = async (data) => {
    setTokenError(null)
    mutation.mutate(data)
  }

  if (tokenError && !search.token) {
    return (
      <Container
        h="100vh"
        maxW="sm"
        alignItems="stretch"
        justifyContent="center"
        gap={4}
        centerContent
      >
        <Alert.Root status="error">
          <Alert.Icon>
            <FiAlertCircle />
          </Alert.Icon>
          <Alert.Title>Invalid Reset Link</Alert.Title>
          <Alert.Description>{tokenError}</Alert.Description>
        </Alert.Root>
        <Button
          variant="solid"
          onClick={() => navigate({ to: "/recover-password" })}
        >
          Request New Reset Link
        </Button>
      </Container>
    )
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
        Reset Password
      </Heading>
      <Text textAlign="center">
        Please enter your new password and confirm it to reset your password.
      </Text>
      {tokenError && (
        <Alert.Root status="error">
          <Alert.Icon>
            <FiAlertCircle />
          </Alert.Icon>
          <Alert.Description>{tokenError}</Alert.Description>
        </Alert.Root>
      )}
      <PasswordInput
        startElement={<FiLock />}
        type="new_password"
        errors={errors}
        {...register("new_password", passwordRules())}
        placeholder="New Password"
      />
      <PasswordInput
        startElement={<FiLock />}
        type="confirm_password"
        errors={errors}
        {...register("confirm_password", confirmPasswordRules(getValues))}
        placeholder="Confirm Password"
      />
      <Button variant="solid" type="submit" loading={isSubmitting || mutation.isPending}>
        Reset Password
      </Button>
    </Container>
  )
}
