import useAuth, { isLoggedIn } from "@/hooks/useAuth"
import { useUIStore } from "@/stores/uiStore"
import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Input,
  Text,
} from "@chakra-ui/react"
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router"

export const Route = createFileRoute("/webinars/manage" as any)({
  component: WebinarsManage,
  beforeLoad: async () => {
    if (!isLoggedIn()) {
      throw redirect({ to: "/login", search: { redirect: "/webinars/manage" } })
    }
  },
})

function WebinarsManage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { selectedRole } = useUIStore()

  const canManage =
    user?.roles?.includes("salon_owner") ||
    user?.roles?.includes("superuser") ||
    selectedRole === "salon_owner"

  if (!canManage) {
    navigate({ to: "/choose-role" as any })
    return null
  }

  return (
    <Container maxW="4xl" py={12}>
      <Heading size="lg" mb={6}>
        Manage Webinars
      </Heading>
      <Box
        p={6}
        borderRadius="2xl"
        border="1px solid"
        borderColor="gray.subtle"
        background="white"
      >
        <Heading size="md" mb={4}>
          Create New Event
        </Heading>
        <Flex direction="column" gap={3}>
          <Input placeholder="Title" />
          <Input placeholder="Date (YYYY-MM-DD)" />
          <Input placeholder="Description" />
          <Button
            variant="solid"
            onClick={() => navigate({ to: "/webinars/manage" as any })}
          >
            Create Event
          </Button>
        </Flex>
      </Box>
      <Box
        mt={8}
        p={6}
        borderRadius="2xl"
        border="1px solid"
        borderColor="gray.subtle"
        background="white"
      >
        <Heading size="md" mb={4}>
          Your Events
        </Heading>
        <Text color="gray.500">No events yet. Create your first above.</Text>
      </Box>
    </Container>
  )
}
