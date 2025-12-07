import { Container, Image, Input, Text, Tabs } from "@chakra-ui/react"
import {
  Link as RouterLink,
  createFileRoute,
  redirect,
  useNavigate,
} from "@tanstack/react-router"
import { useMutation } from "@tanstack/react-query"
import { type SubmitHandler, useForm } from "react-hook-form"
import { FiLock, FiUser, FiPhone } from "react-icons/fi"
import { useState } from "react"

import type { UserRegister } from "@/types/user"
import type { PhoneRegisterForm } from "@/types/phone"
import type { ApiError } from "@/client/user"
import { userServiceCreateUser, userServiceRegisterPhoneUser } from "@/client/user"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { PasswordInput } from "@/components/ui/password-input"
import useAuth, { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { confirmPasswordRules, emailPattern, passwordRules, handleError } from "@/utils"
import Logo from "/assets/images/fastapi-logo.svg"

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
  const { showSuccessToast } = useCustomToast()
  const [signupMethod, setSignupMethod] = useState<"email" | "phone">("email")
  const { signUpMutation } = useAuth()

  // Email signup form
  const emailForm = useForm<UserRegisterForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      email: "",
      full_name: "",
      password: "",
      confirm_password: "",
    },
  })

  // Phone signup form
  const phoneForm = useForm<PhoneRegisterForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { phone_number: "", full_name: "", region: "" },
  })

  const phoneRegisterMutation = useMutation({
    mutationFn: async (data: PhoneRegisterForm) => {
      await userServiceRegisterPhoneUser({ requestBody: { phoneNumber: data.phone_number, fullName: data.full_name, region: data.region } })
    },
    onSuccess: (_data, variables) => {
      showSuccessToast("Registration successful. Please verify your phone.")
      const search = new URLSearchParams({ phone_number: variables.phone_number, region: variables.region }).toString()
      navigate({ to: `/verify-phone?${search}` })
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
        defaultValue={signupMethod}
        onValueChange={(e) => setSignupMethod(e.value as "email" | "phone")}
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
              invalid={!!emailForm.formState.errors.full_name}
              errorText={emailForm.formState.errors.full_name?.message}
            >
              <InputGroup w="100%" startElement={<FiUser />}>
                <Input
                  id="full_name"
                  minLength={3}
                  {...emailForm.register("full_name", {
                    required: "Full Name is required",
                  })}
                  placeholder="Full Name"
                  type="text"
                />
              </InputGroup>
            </Field>
            <Field invalid={!!emailForm.formState.errors.email} errorText={emailForm.formState.errors.email?.message}>
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
              {...emailForm.register("confirm_password", confirmPasswordRules(emailForm.getValues))}
              placeholder="Confirm Password"
              errors={emailForm.formState.errors}
            />
            <Button variant="solid" type="submit" loading={emailForm.formState.isSubmitting}>
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
            <Field required invalid={!!phoneForm.formState.errors.full_name} errorText={phoneForm.formState.errors.full_name?.message}>
              <InputGroup w="100%" startElement={<FiUser />}>
                <Input
                  id="full_name"
                  {...phoneForm.register("full_name", { required: "Full Name is required" })}
                  placeholder="Full Name"
                  type="text"
                />
              </InputGroup>
            </Field>
            <Field required invalid={!!phoneForm.formState.errors.phone_number} errorText={phoneForm.formState.errors.phone_number?.message}>
              <InputGroup w="100%" startElement={<FiPhone />}>
                <Input
                  id="phone_number"
                  {...phoneForm.register("phone_number", { required: "Phone number is required" })}
                  placeholder="Phone Number"
                  type="tel"
                />
              </InputGroup>
            </Field>
            <Field required invalid={!!phoneForm.formState.errors.region} errorText={phoneForm.formState.errors.region?.message}>
              <InputGroup w="100%">
                <Input
                  id="region"
                  {...phoneForm.register("region", { required: "Region is required" })}
                  placeholder="Region (e.g., ID)"
                  type="text"
                />
              </InputGroup>
            </Field>
            <Button variant="solid" type="submit" loading={phoneForm.formState.isSubmitting || phoneRegisterMutation.isPending}>
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
