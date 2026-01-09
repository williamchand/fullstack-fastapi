import {
  EventStatus,
  PlatformType,
  type WebinarEvent,
} from "@/webistream/types"
import {
  Calendar as CalendarIcon,
  Globe,
  Layers,
  Lock,
  MoreVertical,
  Users,
  Video,
  Zap,
} from "lucide-react"
import type React from "react"

interface EventListProps {
  events: WebinarEvent[]
  onManageParticipants: (event: WebinarEvent) => void
  onPublish: (event: WebinarEvent) => void
  onOpenControlPanel: (event: WebinarEvent) => void
}

const EventList: React.FC<EventListProps> = ({
  events,
  onManageParticipants,
  onPublish,
  onOpenControlPanel,
}) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
      {events.map((event) => (
        <EventCard
          key={event.id}
          event={event}
          onManageParticipants={() => onManageParticipants(event)}
          onPublish={() => onPublish(event)}
          onControlPanel={() => onOpenControlPanel(event)}
        />
      ))}
    </div>
  )
}

const EventCard: React.FC<{
  event: WebinarEvent
  onManageParticipants: () => void
  onPublish: () => void
  onControlPanel: () => void
}> = ({ event, onManageParticipants, onPublish, onControlPanel }) => {
  const isDraft = event.status === EventStatus.DRAFT

  return (
    <div
      className={`bg-white rounded-[2rem] border transition-all duration-500 overflow-hidden flex flex-col group ${
        isDraft
          ? "border-amber-100 hover:border-amber-300"
          : "border-gray-100 hover:border-indigo-300 shadow-sm hover:shadow-2xl hover:-translate-y-1"
      }`}
    >
      <div
        className={`relative h-40 flex flex-col p-6 justify-between ${
          isDraft
            ? "bg-amber-500"
            : event.platform === PlatformType.ZOOM
              ? "bg-indigo-600"
              : "bg-emerald-600"
        }`}
      >
        <div className="flex justify-between items-start">
          <div className="bg-white/20 backdrop-blur-md p-2.5 rounded-2xl text-white">
            <Video className="w-6 h-6" />
          </div>
          <div className="flex flex-col items-end gap-2">
            <span className="px-3 py-1 bg-white/20 backdrop-blur-md text-white text-[10px] font-black uppercase tracking-widest rounded-full">
              {event.platform}
            </span>
            <span
              className={`px-3 py-1 bg-white text-[10px] font-black uppercase tracking-widest rounded-full shadow-sm ${
                isDraft ? "text-amber-600" : "text-indigo-600"
              }`}
            >
              {event.status}
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {event.meetingLink ? (
            <div className="flex items-center gap-1.5 px-3 py-1 bg-white/20 backdrop-blur-md text-white rounded-full text-xs font-medium">
              <Globe className="w-3 h-3" /> Link Ready
            </div>
          ) : (
            <div className="flex items-center gap-1.5 px-3 py-1 bg-black/20 backdrop-blur-md text-white/70 rounded-full text-xs font-medium">
              <Lock className="w-3 h-3" /> Setup Needed
            </div>
          )}
        </div>
      </div>

      <div className="p-8 flex-1 flex flex-col">
        <div className="flex justify-between items-start mb-4">
          <h3 className="font-bold text-gray-900 text-xl line-clamp-1">
            {event.title}
          </h3>
          <button className="text-gray-300 hover:text-gray-600 p-1 transition-colors">
            <MoreVertical className="w-5 h-5" />
          </button>
        </div>

        <p className="text-gray-500 text-sm line-clamp-2 mb-6 flex-1 leading-relaxed">
          {event.description}
        </p>

        <div className="grid grid-cols-2 gap-4 mb-8 text-xs font-bold text-gray-400 uppercase tracking-widest">
          <div className="flex items-center gap-2">
            <CalendarIcon className="w-4 h-4 text-indigo-400" />
            <span className="text-gray-600">
              {new Date(event.date).toLocaleDateString()}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <Users className="w-4 h-4 text-indigo-400" />
            <span className="text-gray-600">
              {event.attendees.length} Users
            </span>
          </div>
        </div>

        <div className="pt-6 border-t border-gray-100 flex items-center justify-between gap-4">
          {isDraft ? (
            <button
              onClick={onPublish}
              className="flex-1 bg-amber-600 text-white px-6 py-3 rounded-2xl text-sm font-black flex items-center justify-center gap-2 hover:bg-amber-700 transition-all shadow-lg shadow-amber-100"
            >
              <Zap className="w-4 h-4" /> Publish Now ($10)
            </button>
          ) : (
            <>
              <button
                onClick={onManageParticipants}
                className="flex-1 bg-gray-50 text-gray-900 px-4 py-3 rounded-2xl text-sm font-black flex items-center justify-center gap-2 hover:bg-indigo-50 hover:text-indigo-600 transition-all border border-transparent hover:border-indigo-100"
              >
                <Users className="w-4 h-4" /> Manage
              </button>
              <button
                onClick={onControlPanel}
                className="bg-indigo-600 text-white p-3 rounded-2xl hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
                title="Open Control Panel"
              >
                <Layers className="w-5 h-5" />
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default EventList
