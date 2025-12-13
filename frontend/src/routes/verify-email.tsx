import { Container, Flex, Heading, Input, Text } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router"
import { type SubmitHandler, useForm } from "react-hook-form"
import { useEffect, useState } from "react"
import { FiMail } from "react-icons/fi"

import type { ApiError } from "@/client/user"
import {
  userServiceResendEmailVerification,
  userServiceVerifyEmailOtp,
} from "@/client/user"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { emailPattern, handleError } from "@/utils"

export const Route = createFileRoute("/verify-email")({
  component: VerifyEmail,
  beforeLoad: async () => {
    if (isLoggedIn()) {
      throw redirect({ to: "/" })
    }
  },
})

interface EmailVerifyForm {
  email: string
  otp_code: string
}

function VerifyEmail() {
  const navigate = useNavigate()
  const { showSuccessToast } = useCustomToast()
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<EmailVerifyForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      email: new URLSearchParams(window.location.search).get("email") || "",
      otp_code: "",
    },
  })

  const email = watch("email")
  const [resendRequested, setResendRequested] = useState(false)
  const [secondsLeft, setSecondsLeft] = useState(0)

  const verifyMutation = useMutation({
    mutationFn: async (data: EmailVerifyForm) => {
      await userServiceVerifyEmailOtp({
        requestBody: { email: data.email, otpCode: data.otp_code },
      })
    },
    onSuccess: () => {
      showSuccessToast("Email verified successfully.")
      navigate({ to: "/login" })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const resendMutation = useMutation({
    mutationFn: async (email: string) => {
      await userServiceResendEmailVerification({ requestBody: { email } })
    },
    onSuccess: () => {
      showSuccessToast("Verification email sent successfully.")
      setResendRequested(true)
      setSecondsLeft(60)
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  useEffect(() => {
    if (!resendRequested || secondsLeft <= 0) return
    const id = setInterval(() => {
      setSecondsLeft((s) => (s > 0 ? s - 1 : 0))
    }, 1000)
    return () => clearInterval(id)
  }, [resendRequested, secondsLeft])

  const onSubmit: SubmitHandler<EmailVerifyForm> = async (data) => {
    verifyMutation.mutate(data)
  }

  const onResend = () => {
    if (email) {
      resendMutation.mutate(email)
    }
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
      display="flex"
      flexDirection="column"
    >
      <Heading size="xl" color="ui.main" textAlign="center" mb={2}>
        Verify Email
      </Heading>
      <Text textAlign="center" mb={4}>
        Please enter the verification code sent to your email address.
      </Text>
      <Flex direction="column" gap={3} width="100%">
        <Field
          required
          invalid={!!errors.email}
          errorText={errors.email?.message}
        >
          <InputGroup w="100%" startElement={<FiMail />}>
            <Input
              id="email"
              {...register("email", {
                required: "Email is required",
                pattern: emailPattern,
              })}
              placeholder="Email"
              type="email"
              readOnly
            />
          </InputGroup>
        </Field>
        <Field
          required
          invalid={!!errors.otp_code}
          errorText={errors.otp_code?.message}
        >
          <InputGroup w="100%">
            <Input
              id="otp_code"
              {...register("otp_code", {
                required: "OTP code is required",
                setValueAs: (v: string) => (v ? v.replace(/\D/g, "") : v),
                validate: {
                  digitsOnly: (v) =>
                    /^\d+$/.test(v) ? true : "OTP must be digits only",
                },
              })}
              placeholder="OTP Code"
              type="text"
              inputMode="numeric"
            />
          </InputGroup>
        </Field>
      </Flex>
      <Button
        variant="solid"
        type="submit"
        loading={isSubmitting || verifyMutation.isPending}
        width="100%"
      >
        Verify Email
      </Button>
      <Flex direction="column" align="stretch" gap={2} width="100%">
        <Text color="gray.600" textAlign="left">Didn't receive the code?</Text>
        <Button
          variant="outline"
          onClick={onResend}
          loading={resendMutation.isPending}
          disabled={!email || resendMutation.isPending || (resendRequested && secondsLeft > 0)}
          size="sm"
          width="100%"
        >
          {resendRequested && secondsLeft > 0
            ? `Resend Verification Email (${secondsLeft}s)`
            : "Resend Verification Email"}
        </Button>
        {resendRequested && (
          <Text color="gray.500" fontSize="sm">
            {secondsLeft > 0
              ? `You can resend in ${secondsLeft}s`
              : "You can resend now."}
          </Text>
        )}
      </Flex>
    </Container>
  )
}
