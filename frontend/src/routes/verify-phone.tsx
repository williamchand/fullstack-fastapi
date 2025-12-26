import { Container, Flex, Image, Input } from "@chakra-ui/react"
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
import Logo from "/assets/images/fastapi-logo.svg"
import { useUIStore } from "@/stores/uiStore"

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
  const { verifyPhoneNumber, verifyRegion } = useUIStore()
  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<PhoneVerifyForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      phone_number: verifyPhoneNumber,
      region: verifyRegion,
      otp_code: "",
    },
  })

  const mutation = useMutation({
    mutationFn: async (data: PhoneVerifyForm) => {
      const res = await userServiceVerifyPhoneOtp({
        requestBody: {
          phoneNumber: data.phone_number,
          otpCode: data.otp_code,
          region: data.region,
        },
      })
      if (res.accessToken) {
        localStorage.setItem("access_token", res.accessToken)
      }
      if (res.refreshToken) {
        localStorage.setItem("refresh_token", res.refreshToken)
      }
      if (res.refreshExpiresAt) {
        localStorage.setItem("refresh_expires_at", res.refreshExpiresAt)
      }
      return res
    },
    onSuccess: () => {
      showSuccessToast("Phone verified successfully.")
      navigate({ to: "/" })
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
      <Image
        src={Logo}
        alt="FastAPI logo"
        height="auto"
        maxW="2xs"
        alignSelf="center"
        mb={4}
      />
      <Flex gap={2} alignItems="flex-start" wrap="wrap">
        <Field
          required
          invalid={!!errors.region}
          errorText={errors.region?.message}
          w={{ base: "100%", md: "80px" }}
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
                size="md"
              />
            )}
          />
        </Field>
        <Field
          required
          invalid={!!errors.phone_number}
          errorText={errors.phone_number?.message}
          flex="1"
          minW={{ base: "100%", md: "auto" }}
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
