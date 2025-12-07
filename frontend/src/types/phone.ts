export type PhoneLoginForm = {
  phone_number: string
  region: string
  otp_code?: string
}

export type PhoneRegisterForm = {
  phone_number: string
  full_name: string
  region: string
}

export type PhoneVerifyForm = {
  phone_number: string
  region: string
  otp_code: string
}
