import type { User, WebinarEvent } from "@/types/webistream"
import { BookOpen, Calendar, Clock, PlayCircle, Video } from "lucide-react"
import type React from "react"

interface CustomerDashboardProps {
  user: User
  events: WebinarEvent[]
  onBrowse: () => void
}

const CustomerDashboard: React.FC<CustomerDashboardProps> = ({
  user,
  events,
  onBrowse,
}) => {
  const myEvents = events.filter((e) =>
    e.attendees.some((a) => a.email === user.email),
  )

  return (
    <div className="max-w-6xl mx-auto space-y-10 p-10 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div>
          <h2 className="text-3xl font-black text-gray-900 tracking-tight">
            My Learning Shelf
          </h2>
          <p className="text-gray-500 font-medium mt-1">
            You are enrolled in {myEvents.length} upcoming sessions.
          </p>
        </div>
        <button
          onClick={onBrowse}
          className="bg-indigo-600 text-white px-8 py-4 rounded-2xl font-bold shadow-xl shadow-indigo-100 hover:scale-105 transition-transform flex items-center gap-2"
        >
          <BookOpen className="w-5 h-5" /> Browse More
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {myEvents.map((event) => (
          <div
            key={event.id}
            className="bg-white rounded-[2.5rem] border border-gray-100 overflow-hidden shadow-sm hover:shadow-2xl transition-all group flex flex-col"
          >
            <div className="relative h-48 bg-gray-900 overflow-hidden">
              <img
                src={`https://picsum.photos/seed/${event.id}/800/600`}
                className="w-full h-full object-cover opacity-60 group-hover:scale-110 transition-transform duration-700"
                alt={event.title}
              />
              <div className="absolute top-4 left-4 px-3 py-1 bg-white/20 backdrop-blur-md rounded-full text-white text-[10px] font-black uppercase tracking-widest">
                {event.platform}
              </div>
            </div>

            <div className="p-8 space-y-4 flex-1 flex flex-col">
              <div className="flex items-center gap-2 text-indigo-600 font-bold text-xs uppercase tracking-widest">
                <Calendar className="w-3 h-3" />{" "}
                {new Date(event.date).toLocaleDateString()}
              </div>
              <h3 className="text-xl font-bold leading-tight line-clamp-2">
                {event.title}
              </h3>

              <div className="flex items-center gap-4 text-gray-500 text-sm font-medium pt-2">
                <div className="flex items-center gap-1">
                  <Clock className="w-4 h-4" /> {event.time}
                </div>
                <div className="flex items-center gap-1">
                  <Video className="w-4 h-4" /> {event.duration}m
                </div>
              </div>

              <div className="pt-6 mt-auto">
                <button className="w-full py-4 bg-gray-900 text-white font-bold rounded-2xl hover:bg-indigo-600 transition-all flex items-center justify-center gap-2">
                  <PlayCircle className="w-5 h-5" /> Join Session
                </button>
              </div>
            </div>
          </div>
        ))}

        {myEvents.length === 0 && (
          <div className="col-span-full py-24 text-center bg-white rounded-[3rem] border-2 border-dashed border-gray-200">
            <BookOpen className="w-16 h-16 text-gray-200 mx-auto mb-6" />
            <h3 className="text-xl font-bold text-gray-900">
              Your shelf is empty
            </h3>
            <p className="text-gray-500 mt-2 max-w-sm mx-auto">
              You haven't registered for any webinars yet. Start browsing to
              find your next favorite topic!
            </p>
            <button
              onClick={onBrowse}
              className="mt-8 text-indigo-600 font-bold hover:underline"
            >
              Discover upcoming sessions
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

export default CustomerDashboard
