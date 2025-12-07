import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Input,
  Text,
  Badge,
} from "@chakra-ui/react"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useState } from "react"
import { type SubmitHandler, useForm } from "react-hook-form"
import { FiMail, FiPhone } from "react-icons/fi"

import type { ApiError, v1User as UserPublic, v1UpdateUserRequest as UserUpdateMe } from "@/client/user"
import {
  userServiceUpdateUser,
  userServiceAddEmail,
  userServiceAddPhoneNumber,
  userServiceResendEmailVerification,
  userServiceVerifyAddEmailOtp,
  userServiceVerifyAddPhoneOtp,
} from "@/client/user"
import useAuth from "@/hooks/useAuth"
import useCustomToast from "@/hooks/useCustomToast"
import { emailPattern, handleError } from "@/utils"
import { Field } from "../ui/field"
import { InputGroup } from "../ui/input-group"

const UserInformation = () => {
  const queryClient = useQueryClient()
  const { showSuccessToast } = useCustomToast()
  const [editMode, setEditMode] = useState(false)
  const [addEmailMode, setAddEmailMode] = useState(false)
  const [addPhoneMode, setAddPhoneMode] = useState(false)
  const [verifyEmailMode, setVerifyEmailMode] = useState(false)
  const [verifyPhoneMode, setVerifyPhoneMode] = useState(false)
  const { user: currentUser } = useAuth()

  const nameForm = useForm<UserPublic>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      fullName: currentUser?.fullName,
    },
  })

  const addEmailForm = useForm<{ email: string }>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { email: "" },
  })

  const verifyEmailForm = useForm<{ otp_code: string }>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { otp_code: "" },
  })

  const addPhoneForm = useForm<{ phone_number: string; region: string }>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { phone_number: "", region: "" },
  })

  const verifyPhoneForm = useForm<{ otp_code: string }>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: { otp_code: "" },
  })

  const updateNameMutation = useMutation({
    mutationFn: (data: UserUpdateMe) =>
      userServiceUpdateUser({ requestBody: { fullName: data.fullName } }),
    onSuccess: () => {
      showSuccessToast("Name updated successfully.")
      setEditMode(false)
      queryClient.invalidateQueries({ queryKey: ["currentUser"] })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const addEmailMutation = useMutation({
    mutationFn: (data: { email: string }) =>
      userServiceAddEmail({ requestBody: { email: data.email } }),
    onSuccess: () => {
      showSuccessToast("Verification email sent. Please check your inbox.")
      setAddEmailMode(false)
      setVerifyEmailMode(true)
      addEmailForm.reset()
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const verifyEmailMutation = useMutation({
    mutationFn: (data: { otp_code: string }) =>
      userServiceVerifyAddEmailOtp({ requestBody: { otpCode: data.otp_code } }),
    onSuccess: () => {
      showSuccessToast("Email verified successfully.")
      setVerifyEmailMode(false)
      queryClient.invalidateQueries({ queryKey: ["currentUser"] })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const resendEmailMutation = useMutation({
    mutationFn: (email: string) =>
      userServiceResendEmailVerification({ requestBody: { email } }),
    onSuccess: () => {
      showSuccessToast("Verification email sent successfully.")
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const addPhoneMutation = useMutation({
    mutationFn: (data: { phone_number: string; region: string }) =>
      userServiceAddPhoneNumber({ requestBody: { phoneNumber: data.phone_number, region: data.region } }),
    onSuccess: () => {
      showSuccessToast("Verification code sent to your phone.")
      setAddPhoneMode(false)
      setVerifyPhoneMode(true)
      addPhoneForm.reset()
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const verifyPhoneMutation = useMutation({
    mutationFn: (data: { otp_code: string }) =>
      userServiceVerifyAddPhoneOtp({ requestBody: { otpCode: data.otp_code } }),
    onSuccess: () => {
      showSuccessToast("Phone verified successfully.")
      setVerifyPhoneMode(false)
      queryClient.invalidateQueries({ queryKey: ["currentUser"] })
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
  })

  const onNameSubmit: SubmitHandler<UserUpdateMe> = async (data) => {
    updateNameMutation.mutate(data)
  }

  const onAddEmailSubmit: SubmitHandler<{ email: string }> = async (data) => {
    addEmailMutation.mutate(data)
  }

  const onVerifyEmailSubmit: SubmitHandler<{ otp_code: string }> = async (data) => {
    verifyEmailMutation.mutate(data)
  }

  const onAddPhoneSubmit: SubmitHandler<{ phone_number: string; region: string }> = async (data) => {
    addPhoneMutation.mutate(data)
  }

  const onVerifyPhoneSubmit: SubmitHandler<{ otp_code: string }> = async (data) => {
    verifyPhoneMutation.mutate(data)
  }

  const onResendEmail = () => {
    if (currentUser?.email) {
      resendEmailMutation.mutate(currentUser.email)
    }
  }

  return (
      <Container maxW="full">
        <Heading size="sm" py={4}>
          User Information
        </Heading>
      <Box w={{ sm: "full", md: "50%" }}>
        {/* Full Name */}
        <Box as="form" onSubmit={nameForm.handleSubmit(onNameSubmit)} mb={6}>
          <Field label="Full name">
            {editMode ? (
              <Input
                {...nameForm.register("fullName", { maxLength: 30 })}
                type="text"
                size="md"
                w="auto"
              />
            ) : (
              <Text
                fontSize="md"
                py={2}
                color={!currentUser?.fullName ? "gray" : "inherit"}
                truncate
                maxWidth="250px"
              >
                {currentUser?.fullName || "N/A"}
              </Text>
            )}
          </Field>
          <Flex mt={4} gap={3}>
            <Button
              variant="solid"
              onClick={() => setEditMode(!editMode)}
              type={editMode ? "button" : "submit"}
              loading={editMode ? updateNameMutation.isPending : false}
              disabled={editMode ? !nameForm.formState.isDirty : false}
            >
              {editMode ? "Save" : "Edit"}
            </Button>
            {editMode && (
              <Button
                variant="subtle"
                colorPalette="gray"
                onClick={() => {
                  nameForm.reset()
                  setEditMode(false)
                }}
                disabled={updateNameMutation.isPending}
              >
                Cancel
              </Button>
            )}
          </Flex>
        </Box>

        {/* Email Section */}
        <Box mb={6}>
          <Flex alignItems="center" gap={2} mb={2}>
            <Text fontWeight="medium">Email</Text>
            {currentUser?.email && (
              <Badge colorPalette={currentUser?.isEmailVerified ? "green" : "orange"}>
                {currentUser.isEmailVerified ? "Verified" : "Unverified"}
              </Badge>
            )}
          </Flex>
          {currentUser?.email ? (
            <>
              <Text fontSize="md" py={2} truncate maxWidth="250px">
                {currentUser.email}
              </Text>
              {!currentUser.isEmailVerified && (
                <Flex mt={2} gap={2}>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={onResendEmail}
                    loading={resendEmailMutation.isPending}
                  >
                    Resend Verification
                  </Button>
                </Flex>
              )}
            </>
          ) : (
            <>
              {!addEmailMode && !verifyEmailMode && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setAddEmailMode(true)}
                  leftIcon={<FiMail />}
                >
                  Add Email
                </Button>
              )}
              {addEmailMode && (
                <Box as="form" onSubmit={addEmailForm.handleSubmit(onAddEmailSubmit)}>
                  <Field invalid={!!addEmailForm.formState.errors.email} errorText={addEmailForm.formState.errors.email?.message}>
                    <InputGroup w="100%" startElement={<FiMail />}>
                      <Input
                        {...addEmailForm.register("email", {
                          required: "Email is required",
                          pattern: emailPattern,
                        })}
                        placeholder="Email"
                        type="email"
                        size="md"
                      />
                    </InputGroup>
                  </Field>
                  <Flex mt={2} gap={2}>
                    <Button variant="solid" type="submit" loading={addEmailMutation.isPending} size="sm">
                      Send Verification
                    </Button>
                    <Button
                      variant="subtle"
                      colorPalette="gray"
                      onClick={() => {
                        setAddEmailMode(false)
                        addEmailForm.reset()
                      }}
                      size="sm"
                    >
                      Cancel
                    </Button>
                  </Flex>
                </Box>
              )}
              {verifyEmailMode && (
                <Box as="form" onSubmit={verifyEmailForm.handleSubmit(onVerifyEmailSubmit)}>
                  <Field invalid={!!verifyEmailForm.formState.errors.otp_code} errorText={verifyEmailForm.formState.errors.otp_code?.message}>
                    <InputGroup w="100%">
                      <Input
                        {...verifyEmailForm.register("otp_code", { required: "OTP code is required" })}
                        placeholder="OTP Code"
                        type="text"
                        size="md"
                      />
                    </InputGroup>
                  </Field>
                  <Flex mt={2} gap={2}>
                    <Button variant="solid" type="submit" loading={verifyEmailMutation.isPending} size="sm">
                      Verify
                    </Button>
                    <Button
                      variant="subtle"
                      colorPalette="gray"
                      onClick={() => {
                        setVerifyEmailMode(false)
                        verifyEmailForm.reset()
                      }}
                      size="sm"
                    >
                      Cancel
                    </Button>
                  </Flex>
                </Box>
              )}
            </>
          )}
        </Box>

        {/* Phone Section */}
        <Box mb={6}>
          <Flex alignItems="center" gap={2} mb={2}>
            <Text fontWeight="medium">Phone Number</Text>
            {currentUser?.phoneNumber && (
              <Badge colorPalette={currentUser?.isPhoneVerified ? "green" : "orange"}>
                {currentUser.isPhoneVerified ? "Verified" : "Unverified"}
              </Badge>
            )}
          </Flex>
          {currentUser?.phoneNumber ? (
            <>
              <Text fontSize="md" py={2} truncate maxWidth="250px">
                {currentUser.phoneNumber}
              </Text>
            </>
          ) : (
            <>
              {!addPhoneMode && !verifyPhoneMode && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setAddPhoneMode(true)}
                  leftIcon={<FiPhone />}
                >
                  Add Phone Number
                </Button>
              )}
              {addPhoneMode && (
                <Box as="form" onSubmit={addPhoneForm.handleSubmit(onAddPhoneSubmit)}>
                  <Field invalid={!!addPhoneForm.formState.errors.phone_number} errorText={addPhoneForm.formState.errors.phone_number?.message}>
                    <InputGroup w="100%" startElement={<FiPhone />}>
                      <Input
                        {...addPhoneForm.register("phone_number", { required: "Phone number is required" })}
                        placeholder="Phone Number"
                        type="tel"
                        size="md"
                      />
                    </InputGroup>
                  </Field>
                  <Field mt={2} invalid={!!addPhoneForm.formState.errors.region} errorText={addPhoneForm.formState.errors.region?.message}>
                    <InputGroup w="100%">
                      <Input
                        {...addPhoneForm.register("region", { required: "Region is required" })}
                        placeholder="Region (e.g., ID)"
                        type="text"
                        size="md"
                      />
                    </InputGroup>
                  </Field>
                  <Flex mt={2} gap={2}>
                    <Button variant="solid" type="submit" loading={addPhoneMutation.isPending} size="sm">
                      Send Verification
                    </Button>
                    <Button
                      variant="subtle"
                      colorPalette="gray"
                      onClick={() => {
                        setAddPhoneMode(false)
                        addPhoneForm.reset()
                      }}
                      size="sm"
                    >
                      Cancel
                    </Button>
                  </Flex>
                </Box>
              )}
              {verifyPhoneMode && (
                <Box as="form" onSubmit={verifyPhoneForm.handleSubmit(onVerifyPhoneSubmit)}>
                  <Field invalid={!!verifyPhoneForm.formState.errors.otp_code} errorText={verifyPhoneForm.formState.errors.otp_code?.message}>
                    <InputGroup w="100%">
                      <Input
                        {...verifyPhoneForm.register("otp_code", { required: "OTP code is required" })}
                        placeholder="OTP Code"
                        type="text"
                        size="md"
                      />
                    </InputGroup>
                  </Field>
                  <Flex mt={2} gap={2}>
                    <Button variant="solid" type="submit" loading={verifyPhoneMutation.isPending} size="sm">
                      Verify
                    </Button>
                    <Button
                      variant="subtle"
                      colorPalette="gray"
                      onClick={() => {
                        setVerifyPhoneMode(false)
                        verifyPhoneForm.reset()
                      }}
                      size="sm"
                    >
                      Cancel
                    </Button>
                  </Flex>
                </Box>
              )}
            </>
          )}
        </Box>
        </Box>
      </Container>
  )
}

export default UserInformation
