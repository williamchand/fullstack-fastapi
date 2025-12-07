import { Container, Heading, Input } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router"
import { type SubmitHandler, useForm } from "react-hook-form"

import type { ApiError } from "@/client/user"
import { userServiceVerifyPhoneOtp } from "@/client/user"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { handleError } from "@/utils"
import type { PhoneVerifyForm } from "@/types/phone"

export const Route = createFileRoute("/verify-phone")({
  component: VerifyPhone,
  beforeLoad: async () => {
    if (isLoggedIn()) {
      throw redirect({ to: "/" })
    }
  },
})

function VerifyPhone() {
  const navigate = useNavigate()
  const { showSuccessToast } = useCustomToast()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<PhoneVerifyForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      phone_number: new URLSearchParams(window.location.search).get("phone_number") || "",
      region: new URLSearchParams(window.location.search).get("region") || "",
      otp_code: "",
    },
  })

  const mutation = useMutation({
    mutationFn: async (data: PhoneVerifyForm) => {
      await userServiceVerifyPhoneOtp({ requestBody: { phoneNumber: data.phone_number, otpCode: data.otp_code, region: data.region } })
    },
    onSuccess: () => {
      showSuccessToast("Phone verified successfully.")
      navigate({ to: "/login" })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const onSubmit: SubmitHandler<PhoneVerifyForm> = async (data) => {
    mutation.mutate(data)
  }

  return (
    <Container as="form" onSubmit={handleSubmit(onSubmit)} h="100vh" maxW="sm" alignItems="stretch" justifyContent="center" gap={4} centerContent>
      <Heading size="xl" color="ui.main" textAlign="center" mb={2}>Verify Phone</Heading>
      <Field required invalid={!!errors.phone_number} errorText={errors.phone_number?.message}>
        <InputGroup w="100%">
          <Input id="phone_number" {...register("phone_number", { required: "Phone number is required" })} placeholder="Phone Number" type="tel" />
        </InputGroup>
      </Field>
      <Field required invalid={!!errors.region} errorText={errors.region?.message}>
        <InputGroup w="100%">
          <Input id="region" {...register("region", { required: "Region is required" })} placeholder="Region (e.g., IDR)" type="text" />
        </InputGroup>
      </Field>
      <Field required invalid={!!errors.otp_code} errorText={errors.otp_code?.message}>
        <InputGroup w="100%">
          <Input id="otp_code" {...register("otp_code", { required: "OTP is required" })} placeholder="OTP Code" type="text" />
        </InputGroup>
      </Field>
      <Button variant="solid" type="submit" loading={isSubmitting || mutation.isPending}>Verify</Button>
    </Container>
  )
}
