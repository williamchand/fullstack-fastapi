import { Container, Input, Text, Tabs, Flex, Heading } from "@chakra-ui/react"
import {
  Link as RouterLink,
  createFileRoute,
  redirect,
  useNavigate,
} from "@tanstack/react-router"
import { useMutation } from "@tanstack/react-query"
import { Controller, type SubmitHandler, useForm } from "react-hook-form"
import { FiLock, FiMail, FiPhone } from "react-icons/fi"
import { useState } from "react"

import type { v1LoginUserRequest as AccessToken, ApiError } from "@/client/user"
import { userServiceLoginWithPhone, userServiceRequestPhoneOtp } from "@/client/user"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { PasswordInput } from "@/components/ui/password-input"
import { RegionSelector } from "@/components/Common/RegionSelector"
import useAuth, { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { emailPattern, passwordRules, handleError } from "../utils"
import type { PhoneLoginForm } from "@/types/phone"

export const Route = createFileRoute("/login")({
  component: Login,
  beforeLoad: async () => {
    if (isLoggedIn()) {
      throw redirect({
        to: "/",
      })
    }
  },
})

function Login() {
  const navigate = useNavigate()
  const { showSuccessToast } = useCustomToast()
  const [loginMethod, setLoginMethod] = useState<"email" | "phone">("email")
  const { loginMutation, error, resetError } = useAuth()

  // Email login form
  const emailForm = useForm<AccessToken>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      username: "",
      password: "",
    },
  })

  // Phone login form
  const phoneForm = useForm<PhoneLoginForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { phone_number: "", region: "ID", otp_code: "" },
  })

  const requestOtp = useMutation({
    mutationFn: async (data: { phone_number: string; region: string }) => {
      await userServiceRequestPhoneOtp({ requestBody: { phoneNumber: data.phone_number, region: data.region } })
    },
    onSuccess: () => {
      showSuccessToast("OTP sent to your phone.")
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const phoneLoginMutation = useMutation({
    mutationFn: async (data: PhoneLoginForm) => {
      const res = await userServiceLoginWithPhone({ requestBody: { phoneNumber: data.phone_number, otpCode: data.otp_code || "", region: data.region } })
      if (res.accessToken) {
        localStorage.setItem("access_token", res.accessToken)
      }
    },
    onSuccess: () => {
      navigate({ to: "/" })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const onEmailSubmit: SubmitHandler<AccessToken> = async (data) => {
    if (emailForm.formState.isSubmitting) return
    resetError()
    try {
      await loginMutation.mutateAsync(data)
    } catch {
      // error is handled by useAuth hook
    }
  }

  const onRequestOtp: SubmitHandler<Pick<PhoneLoginForm, "phone_number" | "region">> = async (data) => {
    requestOtp.mutate(data)
  }

  const onPhoneLogin: SubmitHandler<PhoneLoginForm> = async (data) => {
    phoneLoginMutation.mutate(data)
  }

  return (
      <Container
        h="100vh"
        maxW="sm"
        alignItems="stretch"
        justifyContent="center"
        gap={4}
        centerContent
      >
        <Text fontSize="lg" color="brand.darkKhaki">ameno signy</Text>
        <Heading size="xl" color="ui.main" textAlign="center" mb={2}>Ameno Signy Super App</Heading>
      <Tabs.Root
        defaultValue={loginMethod}
        onValueChange={(e) => {
          setLoginMethod(e.value as "email" | "phone")
          resetError()
        }}
        w="100%"
      >
        <Tabs.List>
          <Tabs.Trigger value="email">Email</Tabs.Trigger>
          <Tabs.Trigger value="phone">Phone</Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="email">
          <Container
            as="form"
            onSubmit={emailForm.handleSubmit(onEmailSubmit)}
            p={0}
            gap={4}
            display="flex"
            flexDirection="column"
          >
        <Field
              invalid={!!emailForm.formState.errors.username}
              errorText={emailForm.formState.errors.username?.message || (error || undefined)}
        >
          <InputGroup w="100%" startElement={<FiMail />}>
            <Input
              id="username"
                  {...emailForm.register("username", {
                    required: "Email is required",
                pattern: emailPattern,
              })}
              placeholder="Email"
              type="email"
            />
          </InputGroup>
        </Field>
        <PasswordInput
          type="password"
          startElement={<FiLock />}
              {...emailForm.register("password", passwordRules())}
          placeholder="Password"
              errors={emailForm.formState.errors}
        />
        <RouterLink to="/recover-password" className="main-link">
          Forgot Password?
        </RouterLink>
            <Button variant="solid" type="submit" loading={emailForm.formState.isSubmitting} size="md">
          Log In
        </Button>
          </Container>
        </Tabs.Content>
        <Tabs.Content value="phone">
          <Container
            p={0}
            gap={4}
            display="flex"
            flexDirection="column"
          >
            <Flex gap={2} alignItems="flex-start">
              <Field required invalid={!!phoneForm.formState.errors.region} errorText={phoneForm.formState.errors.region?.message} w="140px">
                <Controller
                  control={phoneForm.control}
                  name="region"
                  rules={{ required: "Region is required" }}
                  render={({ field }) => (
                    <RegionSelector
                      value={field.value}
                      onChange={field.onChange}
                      disabled={field.disabled}
                      size="sm"
                    />
                  )}
                />
              </Field>
              <Field required invalid={!!phoneForm.formState.errors.phone_number} errorText={phoneForm.formState.errors.phone_number?.message} flex="1">
                <InputGroup w="100%" startElement={<FiPhone />}>
                  <Input
                    id="phone_number"
                    {...phoneForm.register("phone_number", { required: "Phone number is required" })}
                    placeholder="Phone Number"
                    type="tel"
                  />
                </InputGroup>
              </Field>
            </Flex>
            <Button
              variant="solid"
              onClick={phoneForm.handleSubmit(onRequestOtp)}
              loading={requestOtp.isPending}
            >
              Request OTP
            </Button>
            <Field invalid={!!phoneForm.formState.errors.otp_code} errorText={phoneForm.formState.errors.otp_code?.message}>
              <InputGroup w="100%">
                <Input
                  id="otp_code"
                  {...phoneForm.register("otp_code")}
                  placeholder="OTP Code"
                  type="text"
                />
              </InputGroup>
            </Field>
            <Button
              variant="solid"
              onClick={phoneForm.handleSubmit(onPhoneLogin)}
              loading={phoneForm.formState.isSubmitting || phoneLoginMutation.isPending}
            >
              Login
            </Button>
          </Container>
        </Tabs.Content>
      </Tabs.Root>
        <Text>
          Don't have an account?{" "}
          <RouterLink to="/signup" className="main-link">
            Sign Up
          </RouterLink>
        </Text>
      </Container>
  )
}
