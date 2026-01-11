import {
  EventStatus,
  PlatformType,
  type WebinarEvent,
} from "@/types/webistream"
import {
  Calendar,
  Clock,
  DollarSign,
  ExternalLink as LinkIcon,
  Loader2,
  Sparkles,
  Users,
  Video,
  X,
} from "lucide-react"
import type React from "react"
import { useState } from "react"

interface CreateEventModalProps {
  onClose: () => void
  onSubmit: (event: WebinarEvent) => void
}

const CreateEventModal: React.FC<CreateEventModalProps> = ({
  onClose,
  onSubmit,
}) => {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    date: "",
    time: "",
    duration: 60,
    price: 0,
    platform: PlatformType.ZOOM,
    maxAttendees: 100,
    meetingLink: "",
  })

  const handleAIHelp = async () => {
    if (!formData.title) {
      alert("Please enter a title first!")
      return
    }
    setLoading(true)
    await new Promise((r) => setTimeout(r, 600))
    const desc = `Join a deep dive into "${formData.title}" with practical insights and live Q&A.`
    const suggestedPrice = Math.max(
      15,
      Math.min(99, Math.round((formData.maxAttendees / 10) * 3)),
    )
    setFormData((prev) => ({
      ...prev,
      description: desc,
      price: suggestedPrice,
    }))
    setLoading(false)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const newEvent: WebinarEvent = {
      id: Math.random().toString(36).substr(2, 9),
      ...formData,
      status: EventStatus.DRAFT,
      attendees: [],
      meetingLink:
        formData.meetingLink || "https://auto-generated-by-integration.com",
      aiAssistant: true,
    }
    onSubmit(newEvent)
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md animate-in fade-in duration-300">
      <div className="bg-white w-full max-w-3xl rounded-[2.5rem] shadow-2xl overflow-hidden flex flex-col max-h-[95vh]">
        <div className="p-8 border-b border-gray-100 flex items-center justify-between bg-indigo-50/50">
          <div>
            <h3 className="text-2xl font-black text-gray-900 tracking-tight">
              Create Webinar Draft
            </h3>
            <p className="text-sm text-indigo-600 font-bold uppercase tracking-widest mt-1">
              Design your masterclass
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-3 hover:bg-white rounded-2xl transition-all shadow-sm"
          >
            <X className="w-6 h-6 text-gray-500" />
          </button>
        </div>

        <form
          onSubmit={handleSubmit}
          className="p-10 overflow-y-auto space-y-8"
        >
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="md:col-span-2 space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest pl-1">
                Webinar Title
              </label>
              <input
                required
                type="text"
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none transition-all font-bold text-lg"
                placeholder="e.g., Advanced Quantum Computing"
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
              />
            </div>

            <div className="md:col-span-2 space-y-2">
              <div className="flex items-center justify-between mb-1 pl-1">
                <label className="text-sm font-black text-gray-900 uppercase tracking-widest">
                  Event Description
                </label>
                <button
                  type="button"
                  onClick={handleAIHelp}
                  disabled={loading}
                  className="text-xs font-black text-indigo-600 flex items-center gap-2 px-4 py-2 bg-indigo-100 rounded-full hover:bg-indigo-200 disabled:opacity-50 transition-all uppercase tracking-widest"
                >
                  {loading ? (
                    <Loader2 className="w-3 h-3 animate-spin" />
                  ) : (
                    <Sparkles className="w-3 h-3" />
                  )}
                  AI Writing Assistant
                </button>
              </div>
              <textarea
                rows={4}
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none transition-all resize-none font-medium leading-relaxed"
                placeholder="What value will your participants get?"
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest flex items-center gap-2 pl-1">
                <Calendar className="w-4 h-4 text-indigo-500" /> Date
              </label>
              <input
                required
                type="date"
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none font-bold"
                value={formData.date}
                onChange={(e) =>
                  setFormData({ ...formData, date: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest flex items-center gap-2 pl-1">
                <Clock className="w-4 h-4 text-indigo-500" /> Time
              </label>
              <input
                required
                type="time"
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none font-bold"
                value={formData.time}
                onChange={(e) =>
                  setFormData({ ...formData, time: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest flex items-center gap-2 pl-1">
                <Video className="w-4 h-4 text-indigo-500" /> Platform
              </label>
              <select
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none font-bold appearance-none"
                value={formData.platform}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    platform: e.target.value as PlatformType,
                  })
                }
              >
                <option value={PlatformType.ZOOM}>Zoom Meetings</option>
                <option value={PlatformType.GOOGLE_MEET}>Google Meet</option>
                <option value={PlatformType.MICROSOFT_TEAMS}>
                  Microsoft Teams
                </option>
              </select>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest flex items-center gap-2 pl-1">
                <LinkIcon className="w-4 h-4 text-indigo-500" /> Meeting Link
                (Optional)
              </label>
              <input
                type="url"
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none transition-all font-medium text-sm"
                placeholder="https://zoom.us/j/..."
                value={formData.meetingLink}
                onChange={(e) =>
                  setFormData({ ...formData, meetingLink: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest flex items-center gap-2 pl-1">
                <DollarSign className="w-4 h-4 text-indigo-500" /> Ticket Price
                ($)
              </label>
              <input
                type="number"
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none font-bold"
                value={formData.price}
                onChange={(e) =>
                  setFormData({ ...formData, price: Number(e.target.value) })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-black text-gray-900 uppercase tracking-widest flex items-center gap-2 pl-1">
                <Users className="w-4 h-4 text-indigo-500" /> Max Attendees
              </label>
              <input
                type="number"
                className="w-full px-6 py-4 bg-gray-50 border-none rounded-2xl focus:ring-4 focus:ring-indigo-100 outline-none font-bold"
                value={formData.maxAttendees}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    maxAttendees: Number(e.target.value),
                  })
                }
              />
            </div>
          </div>

          <div className="pt-8 flex gap-6">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 py-5 text-gray-500 font-bold border-none rounded-3xl hover:bg-gray-100 transition-all"
            >
              Save as Draft
            </button>
            <button
              type="submit"
              className="flex-1 py-5 bg-indigo-600 text-white font-black rounded-3xl hover:bg-indigo-700 shadow-2xl shadow-indigo-100 transition-all transform active:scale-95"
            >
              Continue to Publishing
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default CreateEventModal
