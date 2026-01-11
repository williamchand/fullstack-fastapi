import { EventStatus, type WebinarEvent } from "@/types/webistream"
import {
  Activity,
  Cloud,
  MessageSquare,
  MicOff,
  Monitor,
  Play,
  Radio,
  Settings,
  Shield,
  Square,
  Users,
  Video,
  X,
} from "lucide-react"
import type React from "react"
import { useEffect, useState } from "react"

interface ControlPanelProps {
  event: WebinarEvent
  onClose: () => void
}

const ControlPanel: React.FC<ControlPanelProps> = ({ event, onClose }) => {
  const [isRecording, setIsRecording] = useState(false)
  const [isLive, setIsLive] = useState(event.status === EventStatus.LIVE)
  const [activeViewers, setActiveViewers] = useState(
    Math.floor(event.attendees.length * 0.8),
  )

  useEffect(() => {
    if (isLive) {
      const interval = setInterval(() => {
        setActiveViewers((v) => v + (Math.random() > 0.5 ? 1 : -1))
      }, 3000)
      return () => clearInterval(interval)
    }
  }, [isLive])

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/90 backdrop-blur-xl animate-in fade-in duration-300">
      <div className="bg-white w-full max-w-6xl h-[85vh] rounded-[3rem] shadow-2xl overflow-hidden flex flex-col border border-white/20">
        <div className="bg-gray-900 p-8 flex items-center justify-between text-white border-b border-white/10">
          <div className="flex items-center gap-6">
            <div
              className={`p-4 rounded-2xl ${
                isLive
                  ? "bg-red-500 shadow-lg shadow-red-500/20 animate-pulse"
                  : "bg-indigo-600"
              }`}
            >
              <Radio className="w-8 h-8 text-white" />
            </div>
            <div>
              <div className="flex items-center gap-3">
                <h2 className="text-2xl font-black tracking-tight">
                  {event.title}
                </h2>
                {isLive && (
                  <span className="px-3 py-1 bg-red-500 text-[10px] font-black uppercase tracking-widest rounded-full">
                    Live Session
                  </span>
                )}
              </div>
              <p className="text-gray-400 font-medium text-sm flex items-center gap-2 mt-1">
                <Cloud className="w-4 h-4" /> Connected to {event.platform}{" "}
                Platform
              </p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="text-right mr-6 hidden md:block">
              <p className="text-xs font-black text-gray-500 uppercase tracking-widest">
                Session Time
              </p>
              <p className="text-xl font-mono font-bold">01:24:45</p>
            </div>
            <button
              onClick={onClose}
              className="p-4 hover:bg-white/10 rounded-2xl transition-all border border-white/5"
            >
              <X className="w-6 h-6" />
            </button>
          </div>
        </div>

        <div className="flex-1 flex overflow-hidden">
          <div className="flex-1 overflow-y-auto p-10 bg-gray-50/50">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
              <LiveStatCard
                icon={<Users className="w-5 h-5" />}
                label="Live Viewers"
                value={activeViewers.toString()}
                color="text-indigo-600"
              />
              <LiveStatCard
                icon={<MessageSquare className="w-5 h-5" />}
                label="Active Chatters"
                value="12"
                color="text-emerald-600"
              />
              <LiveStatCard
                icon={<Activity className="w-5 h-5" />}
                label="Engagement"
                value="94%"
                color="text-amber-600"
              />
            </div>

            <h3 className="text-sm font-black text-gray-900 uppercase tracking-widest mb-6 pl-1">
              Live Controls
            </h3>
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
              <ActionButton
                onClick={() => setIsLive(!isLive)}
                active={isLive}
                icon={
                  isLive ? (
                    <Square className="w-6 h-6" />
                  ) : (
                    <Play className="w-6 h-6" />
                  )
                }
                label={isLive ? "End Webinar" : "Go Live Now"}
                primary={!isLive}
                danger={isLive}
              />
              <ActionButton
                onClick={() => setIsRecording(!isRecording)}
                active={isRecording}
                icon={
                  <Radio
                    className={`w-6 h-6 ${isRecording ? "animate-pulse" : ""}`}
                  />
                }
                label={isRecording ? "Stop Recording" : "Start Record"}
              />
              <ActionButton
                icon={<MicOff className="w-6 h-6" />}
                label="Mute All"
              />
              <ActionButton
                icon={<Shield className="w-6 h-6" />}
                label="Entry Lock"
              />
            </div>

            <div className="mt-12 bg-white rounded-[2.5rem] p-8 border border-gray-100 shadow-sm">
              <div className="flex items-center justify-between mb-8">
                <h4 className="text-lg font-black text-gray-900">
                  Platform Integration
                </h4>
                <div className="flex items-center gap-2 px-4 py-2 bg-emerald-50 text-emerald-600 rounded-full text-xs font-bold">
                  <Monitor className="w-4 h-4" /> Healthy Link
                </div>
              </div>
              <div className="flex items-center justify-between p-6 bg-gray-50 rounded-3xl border border-dashed border-gray-200">
                <div className="flex items-center gap-4">
                  <div className="bg-white p-3 rounded-2xl shadow-sm">
                    <Video className="w-6 h-6 text-indigo-600" />
                  </div>
                  <div>
                    <p className="text-xs font-bold text-gray-400 uppercase tracking-widest">
                      Direct Platform Link
                    </p>
                    <p className="text-gray-900 font-bold truncate max-w-md">
                      {event.meetingLink}
                    </p>
                  </div>
                </div>
                <button className="px-6 py-3 bg-white text-gray-900 font-bold rounded-2xl border border-gray-200 hover:bg-gray-100 transition-all shadow-sm">
                  Launch Platform
                </button>
              </div>
            </div>
          </div>

          <div className="w-96 bg-white border-l border-gray-100 flex flex-col">
            <div className="p-8 border-b border-gray-50 flex items-center justify-between">
              <h4 className="font-black text-gray-900 uppercase tracking-widest text-xs">
                Live Feed
              </h4>
              <Settings className="w-4 h-4 text-gray-400" />
            </div>
            <div className="flex-1 overflow-y-auto p-6 space-y-6">
              <FeedItem
                time="14:24"
                text="Alex joined as Co-Host AI"
                type="system"
              />
              <FeedItem
                time="14:25"
                text="Recording automatically started"
                type="record"
              />
              <FeedItem
                time="14:28"
                text="Viewer #45 raised hand"
                type="viewer"
              />
              <FeedItem
                time="14:30"
                text="Broadcast message sent: 'Welcome!'"
                type="msg"
              />
              <FeedItem
                time="14:32"
                text="Viewer #12 left meeting"
                type="viewer"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

const LiveStatCard: React.FC<{
  icon: React.ReactNode
  label: string
  value: string
  color: string
}> = ({ icon, label, value, color }) => (
  <div className="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm">
    <div className="flex items-center gap-3 mb-2">
      <div className={`p-2 rounded-xl bg-gray-50 ${color}`}>{icon}</div>
      <span className="text-xs font-bold text-gray-400 uppercase tracking-widest">
        {label}
      </span>
    </div>
    <p className="text-3xl font-black text-gray-900">{value}</p>
  </div>
)

const ActionButton: React.FC<{
  icon: React.ReactNode
  label: string
  active?: boolean
  primary?: boolean
  danger?: boolean
  onClick?: () => void
}> = ({ icon, label, active, primary, danger, onClick }) => (
  <button
    onClick={onClick}
    className={`flex flex-col items-center justify-center gap-4 p-8 rounded-[2rem] transition-all border-2 ${
      primary
        ? "bg-indigo-600 text-white border-indigo-600 shadow-xl shadow-indigo-100"
        : danger && active
          ? "bg-red-600 text-white border-red-600"
          : "bg-white text-gray-900 border-gray-100 hover:border-indigo-100 shadow-sm"
    }`}
  >
    {icon}
    <span className="text-xs font-black uppercase tracking-widest">
      {label}
    </span>
  </button>
)

const FeedItem: React.FC<{ time: string; text: string; type: string }> = ({
  time,
  text,
}) => (
  <div className="flex gap-4">
    <span className="text-[10px] font-bold text-gray-400 pt-1 font-mono">
      {time}
    </span>
    <p className="text-sm text-gray-700 leading-relaxed">
      <span className="font-bold text-indigo-600 mr-1 opacity-60">#</span>{" "}
      {text}
    </p>
  </div>
)

export default ControlPanel
