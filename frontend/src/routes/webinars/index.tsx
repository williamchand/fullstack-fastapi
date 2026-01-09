import { isLoggedIn } from "@/hooks/useAuth"
import { Box, Button, Container, Flex, Heading, Text } from "@chakra-ui/react"
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router"

export const Route = createFileRoute("/webinars/" as any)({
  component: WebinarsHome,
  beforeLoad: async () => {
    if (!isLoggedIn()) {
      throw redirect({ to: "/login", search: { redirect: "/webinars" } })
    }
  },
})

function WebinarsHome() {
  const navigate = useNavigate()
  const upcoming = [
    {
      id: "w1",
      title: "Getting Started with Our Platform",
      date: "2026-01-20",
    },
    { id: "w2", title: "Advanced Tips and Q&A", date: "2026-02-05" },
  ]

  return (
    <Container maxW="4xl" py={12}>
      <Heading size="lg" mb={6}>
        Upcoming Webinars
      </Heading>
      <Flex direction="column" gap={4}>
        {upcoming.map((w) => (
          <Box
            key={w.id}
            p={6}
            borderRadius="2xl"
            border="1px solid"
            borderColor="gray.subtle"
            background="white"
          >
            <Heading size="md">{w.title}</Heading>
            <Text color="gray.500" mt={1}>
              {new Date(w.date).toLocaleDateString()}
            </Text>
            <Button
              mt={4}
              variant="solid"
              onClick={() => navigate({ to: "/webinars" as any })}
            >
              Register
            </Button>
          </Box>
        ))}
      </Flex>
    </Container>
  )
}
