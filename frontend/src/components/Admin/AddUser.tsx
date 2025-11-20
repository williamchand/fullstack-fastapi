import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query"
import { Controller, type SubmitHandler, useForm } from "react-hook-form"

import { type RolesPublic, type UserCreate, UsersService, RolesService } from "@/client"
import type { ApiError } from "@/client/core/ApiError"
import useCustomToast from "@/hooks/useCustomToast"
import { emailPattern, handleError } from "@/utils"
import {
  Button,
  DialogActionTrigger,
  DialogTitle,
  Flex,
  Input,
  Text,
  VStack,
  Listbox,
  Popover,
  Portal,
  Checkmark,
  createListCollection,
} from "@chakra-ui/react"
import { LuChevronDown } from "react-icons/lu"
import { useRef, useState } from "react"
import { FaPlus } from "react-icons/fa"
import { Checkbox } from "../ui/checkbox"
import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTrigger,
} from "../ui/dialog"
import { Field } from "../ui/field"

interface UserCreateForm extends UserCreate {
  confirm_password: string
}


function getRolesQueryOptions() {
  return {
    queryFn: () =>
      RolesService.listRoles(),
    queryKey: ["roles"],
    staleTime: 1000 * 60 * 5, // optional: cache for 5 min
  }
}

const AddUser = () => {
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
    formState: { errors, isValid, isSubmitting },
  } = useForm<UserCreateForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      email: "",
      full_name: "",
      password: "",
      confirm_password: "",
      roles: [],
      is_active: false,
    },
  })

  const mutation = useMutation({
    mutationFn: (data: UserCreate) =>
      UsersService.createUser({ requestBody: data }),
    onSuccess: () => {
      showSuccessToast("User created successfully.")
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

  const onSubmit: SubmitHandler<UserCreateForm> = (data) => {
    mutation.mutate(data)
  }

  return (
    <DialogRoot
      size={{ base: "xs", md: "md" }}
      placement="center"
      open={isOpen}
      onOpenChange={({ open }) => setIsOpen(open)}
    >
      <DialogTrigger asChild>
        <Button value="add-user" my={4}>
          <FaPlus fontSize="16px" />
          Add User
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>Add User</DialogTitle>
          </DialogHeader>
          <DialogBody>
            <Text mb={4}>
              Fill in the form below to add a new user to the system.
            </Text>
            <VStack gap={4}>
              <Field
                required
                invalid={!!errors.email}
                errorText={errors.email?.message}
                label="Email"
              >
                <Input
                  id="email"
                  {...register("email", {
                    required: "Email is required",
                    pattern: emailPattern,
                  })}
                  placeholder="Email"
                  type="email"
                />
              </Field>

              <Field
                invalid={!!errors.full_name}
                errorText={errors.full_name?.message}
                label="Full Name"
              >
                <Input
                  id="name"
                  {...register("full_name")}
                  placeholder="Full name"
                  type="text"
                />
              </Field>

              <Field
                required
                invalid={!!errors.password}
                errorText={errors.password?.message}
                label="Set Password"
              >
                <Input
                  id="password"
                  {...register("password", {
                    required: "Password is required",
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
                required
                invalid={!!errors.confirm_password}
                errorText={errors.confirm_password?.message}
                label="Confirm Password"
              >
                <Input
                  id="confirm_password"
                  {...register("confirm_password", {
                    required: "Please confirm your password",
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
                  const rolesData = queryClient.getQueryData<RolesPublic>(["roles"])

                  const [open, setOpen] = useState(false)
                  const selectedRoles = field.value || []
                  const triggerRef = useRef<HTMLButtonElement | null>(null)

                  const collection = createListCollection({
                    items:
                      rolesData?.data.map((r) => ({ value: r.name, label: r.name })) || [],
                  })
                  const isAllSelected = field.value?.length === collection.items.length
                  const isSomeSelected =
                    field.value && field.value.length > 0 && field.value.length < collection.items.length

                  const handleSelectAll = () => {
                    if (isAllSelected) {
                      field.onChange([])
                    } else {
                      field.onChange(collection.items.map((item) => item.value))
                    }
                  }
                  return (
                    <Field disabled={field.disabled} label="Roles" colorPalette="teal">
                      <Listbox.Root
                        collection={collection}
                        selectionMode="multiple"
                        value={selectedRoles}
                        onValueChange={(details) => field.onChange(details.value)}
                        maxW="320px"
                      >
                        <Popover.Root open={open} onOpenChange={(e) => setOpen(e.open)}>
                          <Popover.Trigger asChild>
                            <Button size="sm" ref={triggerRef} variant="outline" alignItems="center" justifyContent="flex-start">
                              <Checkmark
                                onClick={handleSelectAll}
                                filled
                                size="sm"
                                checked={isAllSelected}
                                indeterminate={isSomeSelected}
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
                                      <Listbox.Item key={item.value} item={item}>
                                        <Listbox.ItemText>{item.label}</Listbox.ItemText>
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
                      checked={field.value}
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
            <Button
              variant="solid"
              type="submit"
              disabled={!isValid}
              loading={isSubmitting}
            >
              Save
            </Button>
          </DialogFooter>
        </form>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  )
}

export default AddUser
