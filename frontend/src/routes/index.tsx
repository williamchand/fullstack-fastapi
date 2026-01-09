import useAuth from "@/hooks/useAuth"
import AuthModal from "@/webistream/components/AuthModal"
import LandingPage from "@/webistream/components/LandingPage"
import StripeCheckoutModal from "@/webistream/components/StripeCheckoutModal"
import type { Attendee, User, WebinarEvent } from "@/webistream/types"
import { EventStatus, PlatformType, UserRole } from "@/webistream/types"
import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useMemo, useState } from "react"

export const Route = createFileRoute("/" as any)({
  component: HomeLanding,
})

function HomeLanding() {
  const navigate = useNavigate()
  const { user: authUser, logout } = useAuth()
  const currentUser: User | null = useMemo(() => {
    return authUser
      ? {
          id: (authUser as any).id || "me",
          name: authUser.fullName || authUser.email || "User",
          email: authUser.email || "",
          role: authUser.roles?.includes("salon_owner")
            ? UserRole.BUSINESS
            : UserRole.CUSTOMER,
          avatar: "https://i.pravatar.cc/100",
          isStripeConnected: authUser.roles?.includes("salon_owner") || false,
        }
      : null
  }, [authUser])

  const [events, setEvents] = useState<WebinarEvent[]>([
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
  ])

  const [showAuth, setShowAuth] = useState(false)
  const [authRole, setAuthRole] = useState<UserRole | undefined>(undefined)
  const [checkout, setCheckout] = useState<{
    amount: number
    title: string
    eventId: string
    purpose: "REGISTRATION" | "PUBLISHING"
  } | null>(null)

  const handleAuthClick = (role?: UserRole) => {
    setAuthRole(role)
    setShowAuth(true)
  }

  const handleAuthSuccess = () => {
    setShowAuth(false)
  }

  const handleRegister = (event: WebinarEvent) => {
    setCheckout({
      amount: event.price,
      title: event.title,
      eventId: event.id,
      purpose: "REGISTRATION",
    })
  }

  const handleCheckoutSuccess = () => {
    if (!checkout) return
    const attendee: Attendee = {
      id: Math.random().toString(36).slice(2),
      name: currentUser?.name || "You",
      email: currentUser?.email || "you@example.com",
      paid: true,
      joinedAt: new Date().toISOString(),
    }
    setEvents((prev) =>
      prev.map((e) =>
        e.id === checkout.eventId
          ? { ...e, attendees: [...e.attendees, attendee] }
          : e,
      ),
    )
    setCheckout(null)
  }

  return (
    <>
      <LandingPage
        events={events}
        currentUser={currentUser}
        onHostLogin={() => navigate({ to: "/webinars/manage" as any })}
        onMyWebinars={() =>
          navigate({ to: "/webistream/page/My Webinars" as any })
        }
        onRegister={handleRegister}
        onAuthClick={handleAuthClick}
        onLogout={logout}
      />
      {showAuth && (
        <AuthModal
          onClose={() => setShowAuth(false)}
          onAuthSuccess={handleAuthSuccess}
          initialRole={authRole}
        />
      )}
      {checkout && (
        <StripeCheckoutModal
          amount={checkout.amount}
          title={checkout.title}
          purpose={checkout.purpose}
          onClose={() => setCheckout(null)}
          onSuccess={handleCheckoutSuccess}
        />
      )}
    </>
  )
}
