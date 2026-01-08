export type PhoneLoginForm = {
  phone_number: string
  region: string
  otp_code?: string
  token?: string
}

export type PhoneRegisterForm = {
  phone_number: string
  fullName: string
  region: string
}

export type PhoneVerifyForm = {
  phone_number: string
  region: string
  otp_code: string
  verification_token: string
}
