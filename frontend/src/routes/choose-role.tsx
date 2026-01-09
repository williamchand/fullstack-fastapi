import useAuth, { isLoggedIn } from "@/hooks/useAuth"
import { useUIStore } from "@/stores/uiStore"
import { Box, Button, Container, Flex, Heading, Text } from "@chakra-ui/react"
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router"

export const Route = createFileRoute("/choose-role" as any)({
  component: ChooseRole,
  beforeLoad: async () => {
    if (!isLoggedIn()) {
      throw redirect({ to: "/login", search: { redirect: "/choose-role" } })
    }
  },
})

function ChooseRole() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { setSelectedRole } = useUIStore()

  if (!user) return null

  const selectAsCustomer = () => {
    setSelectedRole("customer")
    navigate({ to: "/webinars" as any })
  }
  const selectAsBusiness = () => {
    setSelectedRole("salon_owner")
    navigate({ to: "/webinars/manage" as any })
  }

  return (
    <Container maxW="lg" py={12}>
      <Heading size="lg" textAlign="center" mb={8}>
        Choose How You Will Use Webinars
      </Heading>
      <Text textAlign="center" mb={10}>
        Select a role to continue. You can switch anytime.
      </Text>
      <Flex gap={6} direction={{ base: "column", md: "row" }}>
        <Box
          flex="1"
          p={6}
          borderRadius="2xl"
          border="1px solid"
          borderColor="gray.subtle"
          background="white"
        >
          <Heading size="md" mb={2}>
            Customer
          </Heading>
          <Text mb={6}>
            Browse and join webinars. Manage your registrations.
          </Text>
          <Button variant="solid" onClick={selectAsCustomer} width="full">
            Continue as Customer
          </Button>
        </Box>
        <Box
          flex="1"
          p={6}
          borderRadius="2xl"
          border="1px solid"
          borderColor="gray.subtle"
          background="white"
        >
          <Heading size="md" mb={2}>
            Business
          </Heading>
          <Text mb={6}>
            Create and manage webinars. Track attendees and broadcasts.
          </Text>
          <Button variant="solid" onClick={selectAsBusiness} width="full">
            Continue as Business
          </Button>
        </Box>
      </Flex>
    </Container>
  )
}
