import { Container, Flex, Image, Input, Tabs, Text } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import {
  Link as RouterLink,
  createFileRoute,
  redirect,
  useNavigate,
  useSearch,
} from "@tanstack/react-router"
import { useEffect, useState } from "react"
import { Controller, type SubmitHandler, useForm } from "react-hook-form"
import { FiLock, FiPhone, FiUser } from "react-icons/fi"

import type { ApiError } from "@/client/user"
import {
  userServiceCreateUser,
  userServiceRegisterPhoneUser,
} from "@/client/user"
import { RegionSelector } from "@/components/Common/RegionSelector"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { PasswordInput } from "@/components/ui/password-input"
import useAuth, { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import type { PhoneRegisterForm } from "@/types/phone"
import type { UserRegister } from "@/types/user"
import {
  confirmPasswordRules,
  emailPattern,
  handleError,
  passwordRules,
} from "@/utils"
import Logo from "/assets/images/fastapi-logo.svg"
import { useUIStore } from "@/stores/uiStore"

export const Route = createFileRoute("/signup")({
  component: SignUp,
  beforeLoad: async () => {
    if (isLoggedIn()) {
      throw redirect({
        to: "/",
      })
    }
  },
})

interface UserRegisterForm extends UserRegister {
  confirm_password: string
}

function SignUp() {
  const navigate = useNavigate()
  const search = useSearch({ from: "/signup" }) as {
    method?: "email" | "phone"
    redirect?: string
  }
  const { showSuccessToast } = useCustomToast()
  const { method, setMethod, setVerifyData } = useUIStore()

  // Sync tab with search param and clean URL
  useEffect(() => {
    const m = search?.method
    const { method: _method, ...rest } = (search || {}) as Record<string, unknown>
    if (m === "phone") {
      setMethod("phone")
      navigate({ to: "/signup", search: rest as any, replace: true })
    } else if (m === "email") {
      setMethod("email")
      navigate({ to: "/signup", search: rest as any, replace: true })
    } else if (_method !== undefined) {
      // Invalid method, clean URL
      navigate({ to: "/signup", search: rest as any, replace: true })
    }
  }, [search?.method, navigate, setMethod])
  const { signUpMutation } = useAuth()

  // Email signup form
  const emailForm = useForm<UserRegisterForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      email: "",
      fullName: "",
      password: "",
      confirm_password: "",
    },
  })

  // Phone signup form
  const phoneForm = useForm<PhoneRegisterForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { phone_number: "", fullName: "", region: "ID" },
  })

  const phoneRegisterMutation = useMutation({
    mutationFn: async (data: PhoneRegisterForm) => {
      await userServiceRegisterPhoneUser({
        requestBody: {
          phoneNumber: data.phone_number,
          fullName: data.fullName,
          region: data.region,
        },
      })
    },
    onSuccess: (_data, variables) => {
      showSuccessToast("Registration successful. Please verify your phone.")
      setVerifyData(variables.phone_number, variables.region)
      navigate({ to: "/verify-phone" })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const onEmailSubmit: SubmitHandler<UserRegisterForm> = (data) => {
    signUpMutation.mutate(data)
  }

  const onPhoneSubmit: SubmitHandler<PhoneRegisterForm> = async (data) => {
    phoneRegisterMutation.mutate(data)
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
      <Image
        src={Logo}
        alt="FastAPI logo"
        height="auto"
        maxW="2xs"
        alignSelf="center"
        mb={4}
      />
      <Tabs.Root
        value={method}
        onValueChange={(e) => setMethod(e.value as "email" | "phone")}
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
              invalid={!!emailForm.formState.errors.fullName}
              errorText={emailForm.formState.errors.fullName?.message}
            >
              <InputGroup w="100%" startElement={<FiUser />}>
                <Input
                  id="fullName"
                  minLength={3}
                  {...emailForm.register("fullName", {
                    required: "Full Name is required",
                  })}
                  placeholder="Full Name"
                  type="text"
                />
              </InputGroup>
            </Field>
            <Field
              invalid={!!emailForm.formState.errors.email}
              errorText={emailForm.formState.errors.email?.message}
            >
              <InputGroup w="100%" startElement={<FiUser />}>
                <Input
                  id="email"
                  {...emailForm.register("email", {
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
            <PasswordInput
              type="confirm_password"
              startElement={<FiLock />}
              {...emailForm.register(
                "confirm_password",
                confirmPasswordRules(emailForm.getValues),
              )}
              placeholder="Confirm Password"
              errors={emailForm.formState.errors}
            />
            <Button
              variant="solid"
              type="submit"
              loading={emailForm.formState.isSubmitting}
            >
              Sign Up
            </Button>
          </Container>
        </Tabs.Content>
        <Tabs.Content value="phone">
          <Container
            as="form"
            onSubmit={phoneForm.handleSubmit(onPhoneSubmit)}
            p={0}
            gap={4}
            display="flex"
            flexDirection="column"
          >
            <Field
              required
              invalid={!!phoneForm.formState.errors.fullName}
              errorText={phoneForm.formState.errors.fullName?.message}
            >
              <InputGroup w="100%" startElement={<FiUser />}>
                <Input
                  id="fullName"
                  {...phoneForm.register("fullName", {
                    required: "Full Name is required",
                  })}
                  placeholder="Full Name"
                  type="text"
                />
              </InputGroup>
            </Field>
            <Flex gap={2} alignItems="flex-start">
              <Field
                required
                invalid={!!phoneForm.formState.errors.region}
                errorText={phoneForm.formState.errors.region?.message}
                w="140px"
              >
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
              <Field
                required
                invalid={!!phoneForm.formState.errors.phone_number}
                errorText={phoneForm.formState.errors.phone_number?.message}
                flex="1"
              >
                <InputGroup w="100%" startElement={<FiPhone />}>
                  <Input
                    id="phone_number"
                    {...phoneForm.register("phone_number", {
                      required: "Phone number is required",
                      setValueAs: (v: string) => (v ? v.replace(/\D/g, "") : v),
                      validate: {
                        digitsOnly: (v) =>
                          /^\d+$/.test(v)
                            ? true
                            : "Phone must contain digits only",
                        length: (v) =>
                          v.length >= 6 && v.length <= 15
                            ? true
                            : "Phone length must be 6–15 digits",
                      },
                    })}
                    placeholder="Phone Number"
                    type="tel"
                    inputMode="numeric"
                  />
                </InputGroup>
              </Field>
            </Flex>
            <Button
              variant="solid"
              type="submit"
              loading={
                phoneForm.formState.isSubmitting ||
                phoneRegisterMutation.isPending
              }
            >
              Register
            </Button>
          </Container>
        </Tabs.Content>
      </Tabs.Root>
      <Text>
        Already have an account?{" "}
        <RouterLink to="/login" className="main-link">
          Log In
        </RouterLink>
      </Text>
    </Container>
  )
}

export default SignUp
