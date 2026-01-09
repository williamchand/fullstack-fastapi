export enum PlatformType {
  ZOOM = "Zoom",
  GOOGLE_MEET = "Google Meet",
  MICROSOFT_TEAMS = "Teams",
}

export enum EventStatus {
  DRAFT = "Draft",
  PUBLISHED = "Published",
  UPCOMING = "Upcoming",
  LIVE = "Live",
  COMPLETED = "Completed",
}

export enum UserRole {
  CUSTOMER = "CUSTOMER",
  BUSINESS = "BUSINESS",
}

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  avatar?: string
  isStripeConnected?: boolean
}

export interface Attendee {
  id: string
  name: string
  email: string
  paid: boolean
  joinedAt?: string
}

export interface WebinarEvent {
  id: string
  title: string
  description: string
  date: string
  time: string
  duration: number
  price: number
  platform: PlatformType
  meetingLink: string
  status: EventStatus
  maxAttendees: number
  attendees: Attendee[]
  aiAssistant?: boolean
  isPaidForPublishing?: boolean
}

export interface DashboardStats {
  totalRevenue: number
  totalEvents: number
  activeAttendees: number
  upcomingEventsCount: number
}
