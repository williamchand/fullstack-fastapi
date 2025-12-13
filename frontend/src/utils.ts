import type { ApiError } from "./client/user"
import useCustomToast from "./hooks/useCustomToast"

export const emailPattern = {
  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
  message: "Invalid email address",
}

export const namePattern = {
  value: /^[A-Za-z\s\u00C0-\u017F]{1,30}$/,
  message: "Invalid name",
}

export const passwordRules = (isRequired = true) => {
  const rules: any = {
    minLength: {
      value: 8,
      message: "Password must be at least 8 characters",
    },
  }

  if (isRequired) {
    rules.required = "Password is required"
  }

  return rules
}

export const confirmPasswordRules = (
  getValues: () => any,
  isRequired = true,
) => {
  const rules: any = {
    validate: (value: string) => {
      const password = getValues().password || getValues().new_password
      return value === password ? true : "The passwords do not match"
    },
  }

  if (isRequired) {
    rules.required = "Password confirmation is required"
  }

  return rules
}

export const handleError = (err: ApiError) => {
  const { showErrorToast } = useCustomToast()
  const { message } = getAuthErrorInfo(err)
  showErrorToast(message)
}

export type AuthErrorCode =
  | "INVALID_CREDENTIALS"
  | "USER_NOT_FOUND"
  | "USER_INACTIVE"
  | "EMAIL_NOT_VERIFIED"
  | "PHONE_NOT_VERIFIED"
  | "GENERIC_ERROR"

export const getAuthErrorInfo = (
  err: ApiError,
): { code: AuthErrorCode; message: string } => {
  const body: any = err.body ?? {}
  const detail = body?.detail
  let message = "Something went wrong."

  if (typeof detail === "string") {
    message = detail
  } else if (Array.isArray(detail) && detail.length > 0) {
    message = detail[0]?.msg ?? message
  } else if (typeof body?.message === "string") {
    message = body.message
  } else if (err.statusText) {
    message = err.statusText
  }

  const normalized = (message || "").toLowerCase()

  if (normalized.includes("incorrect email or password")) {
    return { code: "INVALID_CREDENTIALS", message }
  }
  if (normalized.includes("invalid username or password")) {
    return { code: "INVALID_CREDENTIALS", message }
  }
  if (normalized.includes("not found")) {
    return { code: "USER_NOT_FOUND", message }
  }
  if (
    normalized.includes("inactive") ||
    normalized.includes("disabled") ||
    normalized.includes("user is not active")
  ) {
    return { code: "USER_INACTIVE", message }
  }
  if (
    normalized.includes("email not verified") ||
    (normalized.includes("not verified") && normalized.includes("email")) ||
    normalized.includes("verify email")
  ) {
    return { code: "EMAIL_NOT_VERIFIED", message }
  }
  if (
    normalized.includes("phone not verified") ||
    (normalized.includes("not verified") && normalized.includes("phone")) ||
    normalized.includes("verify phone")
  ) {
    return { code: "PHONE_NOT_VERIFIED", message }
  }

  return { code: "GENERIC_ERROR", message }
}
