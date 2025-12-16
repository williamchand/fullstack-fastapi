import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Controller, type SubmitHandler, useForm } from "react-hook-form"

import {
  Button,
  Checkmark,
  DialogActionTrigger,
  DialogRoot,
  DialogTrigger,
  Flex,
  Input,
  Listbox,
  Popover,
  Portal,
  Text,
  VStack,
  createListCollection,
} from "@chakra-ui/react"
import { useRef, useState } from "react"
import { FaExchangeAlt } from "react-icons/fa"
import { LuChevronDown } from "react-icons/lu"

import type {
  ApiError,
  v1User as UserPublic,
  v1AdminUpdateUserRequest as UserUpdate,
} from "@/client/user"
import { userServiceAdminUpdateUser } from "@/client/user"
import useCustomToast from "@/hooks/useCustomToast"
import { handleError } from "@/utils"
import { Checkbox } from "../ui/checkbox"
import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog"
import { Field } from "../ui/field"

interface EditUserProps {
  user: UserPublic
}

interface UserUpdateForm {
  email?: string
  fullName?: string
  password?: string
  confirm_password?: string
  roles?: string[]
  is_active?: boolean
}

function getRolesQueryOptions() {
  const staticRoles = [
    { name: "customer" },
    { name: "salon_owner" },
    { name: "employee" },
    { name: "superuser" },
  ]
  return {
    queryFn: async () => ({ data: staticRoles }),
    queryKey: ["roles"],
    staleTime: 1000 * 60 * 5,
  }
}

const EditUser = ({ user }: EditUserProps) => {
  const [isOpen, setIsOpen] = useState(false)
  const queryClient = useQueryClient()
  useQuery({
    ...getRolesQueryOptions(),
    placeholderData: (prevData) => prevData,
  })
  const { showSuccessToast } = useCustomToast()
  const {
    control,
    register,
    handleSubmit,
    reset,
    getValues,
    formState: { errors, isSubmitting },
  } = useForm<UserUpdateForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      email: user.email,
      fullName: user.fullName,
      password: "",
      confirm_password: "",
      roles: user.roles || [],
      is_active: user.isActive,
    },
  })

  const mutation = useMutation({
    mutationFn: (data: UserUpdate) =>
      userServiceAdminUpdateUser({ requestBody: data }),
    onSuccess: () => {
      showSuccessToast("User updated successfully.")
      reset()
      setIsOpen(false)
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] })
    },
  })

  const onSubmit: SubmitHandler<UserUpdateForm> = async (data) => {
    const updateData: UserUpdate = {
      userId: user.id,
      fullName: data.fullName,
      password: data.password === "" ? undefined : data.password,
      roles: data.roles,
      isActive: data.is_active,
    }
    mutation.mutate(updateData)
  }

  return (
    <DialogRoot
      size={{ base: "xs", md: "md" }}
      placement="center"
      open={isOpen}
      onOpenChange={({ open }) => setIsOpen(open)}
    >
      <DialogTrigger asChild>
        <Button variant="ghost" size="sm">
          <FaExchangeAlt fontSize="16px" />
          Edit User
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>Edit User</DialogTitle>
          </DialogHeader>
          <DialogBody>
            <Text mb={4}>Update the user details below.</Text>
            <VStack gap={4}>
              <Field label="Email">
                <Input
                  id="email"
                  value={user.email || ""}
                  placeholder="Email"
                  type="email"
                  readOnly
                  disabled
                />
              </Field>

              <Field
                invalid={!!errors.fullName}
                errorText={errors.fullName?.message}
                label="Full Name"
              >
                <Input
                  id="name"
                  {...register("fullName")}
                  placeholder="Full name"
                  type="text"
                />
              </Field>

              <Field
                invalid={!!errors.password}
                errorText={errors.password?.message}
                label="Set Password"
              >
                <Input
                  id="password"
                  {...register("password", {
                    minLength: {
                      value: 8,
                      message: "Password must be at least 8 characters",
                    },
                  })}
                  placeholder="Password"
                  type="password"
                />
              </Field>

              <Field
                invalid={!!errors.confirm_password}
                errorText={errors.confirm_password?.message}
                label="Confirm Password"
              >
                <Input
                  id="confirm_password"
                  {...register("confirm_password", {
                    validate: (value) =>
                      value === getValues().password ||
                      "The passwords do not match",
                  })}
                  placeholder="Password"
                  type="password"
                />
              </Field>
            </VStack>

            <Flex mt={4} direction="column" gap={4}>
              <Controller
                control={control}
                name="roles"
                render={({ field }) => {
                  const queryClient = useQueryClient()
                  const rolesData = queryClient.getQueryData<{
                    data: { name: string }[]
                  }>(["roles"])

                  const [open, setOpen] = useState(false)
                  const selectedRoles = (field.value || []) as string[]
                  const triggerRef = useRef<HTMLButtonElement | null>(null)

                  const collection = createListCollection({
                    items:
                      rolesData?.data.map((r: { name: string }) => ({
                        value: r.name,
                        label: r.name,
                      })) || [],
                  })
                  const isAllSelected =
                    selectedRoles.length === collection.items.length
                  const isSomeSelected =
                    selectedRoles.length > 0 &&
                    selectedRoles.length < collection.items.length

                  const handleSelectAll = () => {
                    if (isAllSelected) {
                      field.onChange([])
                    } else {
                      field.onChange(
                        collection.items.map((item) => item.value as string),
                      )
                    }
                  }
                  return (
                    <Field
                      disabled={field.disabled}
                      label="Roles"
                      colorPalette="teal"
                    >
                      <Listbox.Root
                        collection={collection}
                        selectionMode="multiple"
                        value={selectedRoles}
                        onValueChange={(details) =>
                          field.onChange(details.value as string[])
                        }
                        maxW="320px"
                      >
                        <Popover.Root
                          open={open}
                          onOpenChange={(e) => setOpen(e.open)}
                        >
                          <Popover.Trigger asChild>
                            <Button
                              size="sm"
                              ref={triggerRef}
                              variant="outline"
                              alignItems="center"
                              justifyContent="flex-start"
                            >
                              <Checkmark
                                onClick={handleSelectAll}
                                filled
                                size="sm"
                                checked={isAllSelected}
                                indeterminate={isSomeSelected || undefined}
                              />
                              <Listbox.Label>Select Role</Listbox.Label>
                              <LuChevronDown style={{ marginLeft: "auto" }} />
                            </Button>
                          </Popover.Trigger>
                          <Portal>
                            <Popover.Positioner>
                              <Popover.Content _closed={{ animation: "none" }}>
                                <Popover.Body p="0">
                                  <Listbox.Content maxH="300px" roundedTop="0">
                                    {collection.items.map((item) => (
                                      <Listbox.Item
                                        key={item.value as string}
                                        item={item}
                                      >
                                        <Listbox.ItemText>
                                          {item.label as string}
                                        </Listbox.ItemText>
                                        <Listbox.ItemIndicator />
                                      </Listbox.Item>
                                    ))}
                                  </Listbox.Content>
                                </Popover.Body>
                              </Popover.Content>
                            </Popover.Positioner>
                          </Portal>
                        </Popover.Root>
                      </Listbox.Root>
                    </Field>
                  )
                }}
              />
              <Controller
                control={control}
                name="is_active"
                render={({ field }) => (
                  <Field disabled={field.disabled} colorPalette="teal">
                    <Checkbox
                      checked={field.value || false}
                      onCheckedChange={({ checked }) => field.onChange(checked)}
                    >
                      Is active?
                    </Checkbox>
                  </Field>
                )}
              />
            </Flex>
          </DialogBody>

          <DialogFooter gap={2}>
            <DialogActionTrigger asChild>
              <Button
                variant="subtle"
                colorPalette="gray"
                disabled={isSubmitting}
              >
                Cancel
              </Button>
            </DialogActionTrigger>
            <Button variant="solid" type="submit" loading={isSubmitting}>
              Save
            </Button>
          </DialogFooter>
          <DialogCloseTrigger />
        </form>
      </DialogContent>
    </DialogRoot>
  )
}

export default EditUser
