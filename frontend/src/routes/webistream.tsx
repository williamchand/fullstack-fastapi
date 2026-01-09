import useAuth from "@/hooks/useAuth"
import SubPage from "@/webistream/components/SubPage"
import {
  EventStatus,
  PlatformType,
  type User,
  UserRole,
  type WebinarEvent,
} from "@/webistream/types"
import { Container } from "@chakra-ui/react"
import { createFileRoute, useNavigate } from "@tanstack/react-router"

export const Route = createFileRoute("/webistream" as any)({
  component: WebiStreamLanding,
})

function WebiStreamLanding() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const currentUser: User | null = user
    ? {
        id: (user as any).id || "me",
        name: user.fullName || user.email || "User",
        email: user.email || "",
        role: user.roles?.includes("salon_owner")
          ? UserRole.BUSINESS
          : UserRole.CUSTOMER,
        avatar: "https://i.pravatar.cc/100",
      }
    : null

  const events: WebinarEvent[] = [
    {
      id: "evt1",
      title: "AI for Business",
      description: "Leverage AI to grow.",
      date: new Date().toISOString(),
      time: "10:00",
      duration: 60,
      price: 29,
      platform: PlatformType.ZOOM,
      meetingLink: "#",
      status: EventStatus.PUBLISHED,
      maxAttendees: 100,
      attendees: [],
    },
    {
      id: "evt2",
      title: "React + TypeScript Tips",
      description: "Practical patterns.",
      date: new Date().toISOString(),
      time: "14:00",
      duration: 90,
      price: 19,
      platform: PlatformType.GOOGLE_MEET,
      meetingLink: "#",
      status: EventStatus.UPCOMING,
      maxAttendees: 200,
      attendees: [],
    },
  ]

  return (
    <Container maxW="full" p={0}>
      <SubPage
        pageKey={"About Us"}
        onBack={() => navigate({ to: "/" })}
        onHostLogin={() => navigate({ to: "/webinars/manage" as any })}
      />
    </Container>
  )
}
