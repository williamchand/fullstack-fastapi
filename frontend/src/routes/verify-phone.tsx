import { Container, Flex, Heading, Input, Text } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router"
import { Controller, type SubmitHandler, useForm } from "react-hook-form"
import { FiPhone } from "react-icons/fi"

import type { ApiError } from "@/client/user"
import { userServiceVerifyPhoneOtp } from "@/client/user"
import { RegionSelector } from "@/components/Common/RegionSelector"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import type { PhoneVerifyForm } from "@/types/phone"
import { handleError } from "@/utils"

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
  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<PhoneVerifyForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      phone_number:
        new URLSearchParams(window.location.search).get("phone_number") || "",
      region: new URLSearchParams(window.location.search).get("region") || "ID",
      otp_code: "",
    },
  })

  const mutation = useMutation({
    mutationFn: async (data: PhoneVerifyForm) => {
      await userServiceVerifyPhoneOtp({
        requestBody: {
          phoneNumber: data.phone_number,
          otpCode: data.otp_code,
          region: data.region,
        },
      })
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
      <Text fontSize="lg" color="brand.darkKhaki">
        ameno signy
      </Text>
      <Heading size="xl" color="ui.main" textAlign="center" mb={2}>
        Ameno Signy Super App
      </Heading>
      <Flex gap={2} alignItems="flex-start">
        <Field
          required
          invalid={!!errors.region}
          errorText={errors.region?.message}
          w="140px"
        >
          <Controller
            control={control}
            name="region"
            rules={{ required: "Region is required" }}
            render={({ field }) => (
              <RegionSelector
                value={field.value}
                onChange={field.onChange}
                disabled
                size="sm"
              />
            )}
          />
        </Field>
        <Field
          required
          invalid={!!errors.phone_number}
          errorText={errors.phone_number?.message}
          flex="1"
        >
          <InputGroup w="100%" startElement={<FiPhone />}>
            <Input
              id="phone_number"
              {...register("phone_number", {
                required: "Phone number is required",
              })}
              placeholder="Phone Number"
              type="tel"
              readOnly
            />
          </InputGroup>
        </Field>
      </Flex>
      <Field
        required
        invalid={!!errors.otp_code}
        errorText={errors.otp_code?.message}
      >
        <InputGroup w="100%">
          <Input
            id="otp_code"
            {...register("otp_code", {
              required: "OTP is required",
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
      <Button
        variant="solid"
        type="submit"
        loading={isSubmitting || mutation.isPending}
      >
        Verify
      </Button>
    </Container>
  )
}
