import SubPage from "@/webistream/components/SubPage"
import { Container } from "@chakra-ui/react"
import { createFileRoute, useNavigate } from "@tanstack/react-router"

export const Route = createFileRoute("/webistream/page/$pageKey" as any)({
  component: WebiStreamSubPage,
})

function WebiStreamSubPage() {
  const navigate = useNavigate()
  const pathParts = (
    typeof window !== "undefined"
      ? window.location.pathname
      : "/webistream/page"
  )?.split("/")
  const pageKey = decodeURIComponent(pathParts[pathParts.length - 1] || "About")

  return (
    <Container maxW="full" p={0}>
      <SubPage
        pageKey={pageKey}
        onBack={() => navigate({ to: "/webistream" as any })}
        onHostLogin={() => navigate({ to: "/webinars/manage" as any })}
      />
    </Container>
  )
}
