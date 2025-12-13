import { Container, Flex, Image, Input, Tabs, Text } from "@chakra-ui/react"
import { useMutation } from "@tanstack/react-query"
import {
  Link as RouterLink,
  createFileRoute,
  redirect,
  useNavigate,
  useSearch,
} from "@tanstack/react-router"
import { useEffect, useRef, useState } from "react"
import { Controller, type SubmitHandler, useForm } from "react-hook-form"
import { FiLock, FiMail, FiPhone } from "react-icons/fi"
import { FcGoogle } from "react-icons/fc"

import type { v1LoginUserRequest as AccessToken, ApiError } from "@/client/user"
import {
  userServiceLoginWithPhone,
  userServiceRequestPhoneOtp,
  userServiceResendEmailVerification,
} from "@/client/user"
import { RegionSelector } from "@/components/Common/RegionSelector"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/field"
import { InputGroup } from "@/components/ui/input-group"
import { PasswordInput } from "@/components/ui/password-input"
import useAuth, { isLoggedIn } from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import type { PhoneLoginForm } from "@/types/phone"
import Logo from "/assets/images/fastapi-logo.svg"
import { emailPattern, handleError, passwordRules } from "../utils"
import {
  oauthServiceGetOauthUrl,
  oauthServiceHandleOauthCallback,
} from "@/client/oauth"

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
  const search = useSearch({ from: "/login" }) as {
    method?: "email" | "phone"
    redirect?: string
    code?: string
    state?: string
    provider?: string
  }
  const { showSuccessToast, showErrorToast } = useCustomToast()
  const [loginMethod, setLoginMethod] = useState<"email" | "phone">("email")
  const { loginMutation, error, resetError, authErrorInfo } = useAuth()
  const [otpRequested, setOtpRequested] = useState(false)
  const [secondsLeft, setSecondsLeft] = useState(0)
  const otpInputRef = useRef<HTMLInputElement | null>(null)
  const [oauthLoading, setOauthLoading] = useState(false)
  const [oauthPopup, setOauthPopup] = useState<Window | null>(null)

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
      await userServiceRequestPhoneOtp({
        requestBody: { phoneNumber: data.phone_number, region: data.region },
      })
    },
    onSuccess: () => {
      showSuccessToast("OTP sent to your phone.")
      setOtpRequested(true)
      setSecondsLeft(60)
      setTimeout(() => otpInputRef.current?.focus(), 0)
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const resendEmailMutation = useMutation({
    mutationFn: async (email: string) => {
      await userServiceResendEmailVerification({ requestBody: { email } })
    },
    onSuccess: () => {
      showSuccessToast("Verification email sent successfully.")
      setEmailResendRequested(true)
      setEmailResendSecondsLeft(60)
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const phoneLoginMutation = useMutation({
    mutationFn: async (data: PhoneLoginForm) => {
      const res = await userServiceLoginWithPhone({
        requestBody: {
          phoneNumber: data.phone_number,
          otpCode: data.otp_code || "",
          region: data.region,
        },
      })
      if (res.accessToken) {
        localStorage.setItem("access_token", res.accessToken)
      }
    },
    onSuccess: () => {
      const dest = search?.redirect || "/"
      navigate({ to: dest })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const onEmailSubmit: SubmitHandler<AccessToken> = async (data) => {
    resetError()
    console.debug("onEmailSubmit called with:", data)
    try {
      await loginMutation.mutateAsync(data)
      const dest = search?.redirect || "/"
      navigate({ to: dest })
    } catch {
      // error is handled by useAuth hook
    }
  }

  const onRequestOtp: SubmitHandler<
    Pick<PhoneLoginForm, "phone_number" | "region">
  > = async (data) => {
    requestOtp.mutate(data)
  }

  const onPhoneLogin: SubmitHandler<PhoneLoginForm> = async (data) => {
    phoneLoginMutation.mutate(data)
  }

  const [emailResendRequested, setEmailResendRequested] = useState(false)
  const [emailResendSecondsLeft, setEmailResendSecondsLeft] = useState(0)

  // Sync tab with search param (e.g., /login?method=phone). Clear only the `method` param, keep current tab
  useEffect(() => {
    const m = search?.method
    const { method: _method, ...rest } = (search || {}) as Record<string, unknown>
    if (m === "phone") {
      setLoginMethod("phone")
      navigate({ to: "/login", search: rest as any, replace: true })
    } else if (m === "email") {
      setLoginMethod("email")
      navigate({ to: "/login", search: rest as any, replace: true })
    } // if no method param, do nothing so current tab stays
  }, [search?.method, navigate])

  // Countdown timer for resend OTP availability
  useEffect(() => {
    if (!otpRequested || secondsLeft <= 0) return
    const id = setInterval(() => {
      setSecondsLeft((s) => (s > 0 ? s - 1 : 0))
    }, 1000)
    return () => clearInterval(id)
  }, [otpRequested, secondsLeft])

  useEffect(() => {
    if (!emailResendRequested || emailResendSecondsLeft <= 0) return
    const id = setInterval(() => {
      setEmailResendSecondsLeft((s) => (s > 0 ? s - 1 : 0))
    }, 1000)
    return () => clearInterval(id)
  }, [emailResendRequested, emailResendSecondsLeft])

  // Handle OAuth callback if present in URL (e.g., /login?provider=google&code=...)
  useEffect(() => {
    const code = search?.code
    const provider = (search?.provider || "google") as string
    if (!code) return
    ;(async () => {
      const processedKey = search?.state
        ? `oauth_processed:${search.state}`
        : null
      try {
        if (processedKey && sessionStorage.getItem(processedKey)) {
          const { code: _c, state: _s, provider: _p, ...rest } =
            (search || {}) as Record<string, unknown>
          navigate({ to: "/login", search: rest as any, replace: true })
          return
        }
      } catch {}
      try {
        const res = await oauthServiceHandleOauthCallback({
          provider,
          requestBody: { code },
        })
        if (res.accessToken) {
          localStorage.setItem("access_token", res.accessToken)
        }
        if (res.refreshToken) {
          localStorage.setItem("refresh_token", res.refreshToken)
        }
        if (res.refreshExpiresAt) {
          localStorage.setItem(
            "refresh_expires_at",
            res.refreshExpiresAt,
          )
        }
        const dest = search?.redirect || "/"
        navigate({ to: dest })
      } catch (err: any) {
        showErrorToast("OAuth login failed. Please try again.")
      } finally {
        // Clean up query params to avoid re-processing
        const { code: _c, state: _s, provider: _p, ...rest } =
          (search || {}) as Record<string, unknown>
        try {
          if (processedKey) sessionStorage.setItem(processedKey, "1")
        } catch {}
        navigate({ to: "/login", search: rest as any, replace: true })
      }
    })()
  }, [search?.code, search?.provider])

  // Watch OTP value for enabling login button
  const otpVal = phoneForm.watch("otp_code")

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
        value={loginMethod}
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
              errorText={
                emailForm.formState.errors.username?.message ||
                error ||
                undefined
              }
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
            <Button
              variant="solid"
              type="submit"
              loading={emailForm.formState.isSubmitting}
              size="md"
            >
              Log In
            </Button>
            <Button
              variant="outline"
              size="md"
              onClick={async () => {
                setOauthLoading(true)
                try {
                  const res = await oauthServiceGetOauthUrl({ provider: "google" })
                  if (!res.url) {
                    showErrorToast("Failed to get Google login URL.")
                    setOauthLoading(false)
                    return
                  }
                  const dest = search?.redirect || "/"
                  const w = 500
                  const h = 700
                  const y = window.top ? (window.top.outerHeight - h) / 2 : 0
                  const x = window.top ? (window.top.outerWidth - w) / 2 : 0
                  const popup = window.open(
                    res.url,
                    "google-oauth",
                    `toolbar=no,location=no,status=no,menubar=no,scrollbars=yes,resizable=yes,width=${w},height=${h},top=${y},left=${x}`,
                  )
                  if (!popup) {
                    try {
                      const key = res.state
                        ? `oauth_redirect_tab:${res.state}`
                        : "oauth_redirect_tab"
                      sessionStorage.setItem(key, dest)
                    } catch {}
                    window.location.href = res.url
                    setOauthLoading(false)
                    return
                  }
                  try {
                    const key = res.state
                      ? `oauth_redirect_popup:${res.state}`
                      : "oauth_redirect_popup"
                    sessionStorage.setItem(key, dest)
                  } catch {}
                  setOauthPopup(popup)
                  const handler = async (e: MessageEvent) => {
                    if (e.origin !== window.location.origin) return
                    const data = e.data || {}
                    if (data.type !== "oauth-callback") return
                    const processedKey = data.state
                      ? `oauth_processed:${data.state}`
                      : null
                    try {
                      if (processedKey && sessionStorage.getItem(processedKey)) {
                        window.removeEventListener("message", handler)
                        setOauthLoading(false)
                        try {
                          oauthPopup?.close()
                        } catch {}
                        setOauthPopup(null)
                        return
                      }
                    } catch {}
                    try {
                      const r = await oauthServiceHandleOauthCallback({
                        provider: data.provider || "google",
                        requestBody: { code: data.code },
                      })
                      if (r.accessToken) {
                        localStorage.setItem("access_token", r.accessToken)
                      }
                      if (r.refreshToken) {
                        localStorage.setItem("refresh_token", r.refreshToken)
                      }
                      if (r.refreshExpiresAt) {
                        localStorage.setItem(
                          "refresh_expires_at",
                          r.refreshExpiresAt,
                        )
                      }
                      let finalDest = search?.redirect || "/"
                      try {
                        const k = data.state
                          ? `oauth_redirect_popup:${data.state}`
                          : "oauth_redirect_popup"
                        finalDest = sessionStorage.getItem(k) || finalDest
                        sessionStorage.removeItem(k)
                        if (processedKey) sessionStorage.setItem(processedKey, "1")
                      } catch {}
                      navigate({ to: finalDest })
                    } catch {
                      showErrorToast("OAuth login failed. Please try again.")
                    } finally {
                      window.removeEventListener("message", handler)
                      setOauthLoading(false)
                      try {
                        oauthPopup?.close()
                      } catch {}
                      setOauthPopup(null)
                    }
                  }
                  window.addEventListener("message", handler)
                } catch {
                  showErrorToast("Failed to start Google login.")
                  setOauthLoading(false)
                }
              }}
              loading={oauthLoading}
            >
              <FcGoogle /> Login with Google
            </Button>
            {authErrorInfo?.code === "USER_NOT_FOUND" && (
              <Text color="red.500" fontSize="sm">
                User not found.{" "}
                <RouterLink to="/signup" className="main-link">
                  Sign up
                </RouterLink>
              </Text>
            )}
            {(authErrorInfo?.code === "EMAIL_NOT_VERIFIED" ||
              authErrorInfo?.code === "USER_INACTIVE") && (
              <Flex direction="column" gap={2} align="stretch">
                <Text color="red.500" fontSize="sm">
                  {authErrorInfo.message ||
                    "You need to verify your email before logging in."}
                </Text>
                <Flex direction="column" gap={2} width="100%">
                  <Button
                    variant="solid"
                    size="sm"
                    width="100%"
                    onClick={() => {
                      const email = emailForm.getValues().username
                      if (email)
                        navigate({ to: "/verify-email", search: { email } })
                    }}
                  >
                    Verify Email
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => {
                      const email = emailForm.getValues().username
                      if (email) resendEmailMutation.mutate(email)
                    }}
                    loading={resendEmailMutation.isPending}
                    size="sm"
                    width="100%"
                    disabled={
                      !emailForm.getValues().username ||
                      resendEmailMutation.isPending ||
                      (emailResendRequested && emailResendSecondsLeft > 0)
                    }
                  >
                    {emailResendRequested && emailResendSecondsLeft > 0
                      ? `Resend Email Verification (${emailResendSecondsLeft}s)`
                      : "Resend Email Verification"}
                  </Button>
                </Flex>
                <Text color="gray.600" fontSize="sm">
                  If you registered with phone, you can log in using phone.
                  <RouterLink
                    to="/login"
                    search={{ method: "phone" }}
                    className="main-link"
                  >
                    Login
                  </RouterLink>
                </Text>
              </Flex>
            )}
          </Container>
        </Tabs.Content>
        <Tabs.Content value="phone">
          <Container
            as="form"
            p={0}
            gap={4}
            display="flex"
            flexDirection="column"
          >
            <Flex gap={2} alignItems="flex-start" wrap="wrap">
              <Field
                required
                invalid={!!phoneForm.formState.errors.region}
                errorText={phoneForm.formState.errors.region?.message}
                w={{ base: "100%", md: "80px" }}
              >
                <Controller
                  control={phoneForm.control}
                  name="region"
                  rules={{ required: "Region is required" }}
                  render={({ field }) => (
                    <RegionSelector
                      value={field.value}
                      onChange={field.onChange}
                      disabled={otpRequested || field.disabled}
                      size="md"
                    />
                  )}
                />
              </Field>
              <Field
                required
                invalid={!!phoneForm.formState.errors.phone_number}
                errorText={phoneForm.formState.errors.phone_number?.message}
                flex="1"
                minW={{ base: "100%", md: "auto" }}
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
                    readOnly={otpRequested}
                  />
                </InputGroup>
              </Field>
            </Flex>
            <Button
              variant="solid"
              onClick={phoneForm.handleSubmit(onRequestOtp)}
              loading={requestOtp.isPending}
              disabled={
                requestOtp.isPending || (otpRequested && secondsLeft > 0)
              }
            >
              {otpRequested
                ? secondsLeft > 0
                  ? `Resend OTP (${secondsLeft}s)`
                  : "Resend OTP"
                : "Request OTP"}
            </Button>
            {otpRequested && (
              <>
                <Field
                  invalid={!!phoneForm.formState.errors.otp_code}
                  errorText={phoneForm.formState.errors.otp_code?.message}
                >
                  <InputGroup w="100%">
                    {(() => {
                      const otpReg = phoneForm.register("otp_code", {
                        setValueAs: (v: string) =>
                          v ? v.replace(/\D/g, "") : v,
                        required: "OTP code is required",
                        validate: {
                          digitsOnly: (v) =>
                            v === undefined ||
                            v === "" ||
                            /^\d+$/.test(v as string)
                              ? true
                              : "OTP must be digits only",
                        },
                      })
                      return (
                        <Input
                          id="otp_code"
                          {...otpReg}
                          onChange={(e) => {
                            // call react-hook-form original onChange
                            otpReg.onChange(e)
                            // clear error as user types so message disappears
                            phoneForm.clearErrors("otp_code")
                          }}
                          placeholder="OTP Code"
                          type="text"
                          inputMode="numeric"
                          ref={(el) => {
                            otpReg.ref(el)
                            otpInputRef.current = el
                          }}
                        />
                      )
                    })()}
                  </InputGroup>
                </Field>
                <Text color="gray.500" fontSize="sm">
                  Waiting for OTP…{" "}
                  {secondsLeft > 0
                    ? `Resend available in ${secondsLeft}s`
                    : "You can resend now."}
                </Text>
              </>
            )}
            {otpRequested && (
              <Button
                variant="solid"
                onClick={phoneForm.handleSubmit(onPhoneLogin)}
                loading={
                  phoneForm.formState.isSubmitting ||
                  phoneLoginMutation.isPending
                }
                disabled={!otpVal || phoneLoginMutation.isPending}
              >
                Login
              </Button>
            )}
            {(authErrorInfo?.code === "PHONE_NOT_VERIFIED" ||
              authErrorInfo?.code === "USER_INACTIVE") && (
              <Flex direction="column" gap={2}>
                <Text color="red.500" fontSize="sm">
                  {authErrorInfo.message ||
                    "You need to verify your phone before logging in."}
                </Text>
                <RouterLink
                  to="/verify-phone"
                  search={{
                    phone_number: phoneForm.getValues("phone_number"),
                    region: phoneForm.getValues("region"),
                  }}
                  className="main-link"
                >
                  Verify Phone
                </RouterLink>
              </Flex>
            )}
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
