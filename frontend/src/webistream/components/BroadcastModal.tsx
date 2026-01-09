import type { WebinarEvent } from "@/webistream/types"
import { CheckCircle, Info, Loader2, Send, X } from "lucide-react"
import type React from "react"
import { useState } from "react"

interface BroadcastModalProps {
  event: WebinarEvent
  onClose: () => void
}

const BroadcastModal: React.FC<BroadcastModalProps> = ({ event, onClose }) => {
  const [message, setMessage] = useState("")
  const [sending, setSending] = useState(false)
  const [sent, setSent] = useState(false)

  const handleSend = () => {
    if (!message) return
    setSending(true)
    setTimeout(() => {
      setSending(false)
      setSent(true)
      setTimeout(() => onClose(), 2000)
    }, 1500)
  }

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm animate-in fade-in">
      <div className="bg-white w-full max-w-lg rounded-[2.5rem] shadow-2xl overflow-hidden flex flex-col border border-gray-100">
        <div className="p-8 border-b border-gray-100 flex items-center justify-between bg-indigo-50/30">
          <div className="flex items-center gap-4">
            <div className="bg-indigo-600 p-2.5 rounded-2xl shadow-lg shadow-indigo-100">
              <Send className="w-5 h-5 text-white" />
            </div>
            <div>
              <h3 className="text-xl font-black text-gray-900 tracking-tight">
                Broadcast Message
              </h3>
              <p className="text-xs text-indigo-600 font-bold uppercase tracking-widest mt-0.5">
                Announcement System
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-3 hover:bg-white rounded-2xl transition-all shadow-sm"
          >
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        <div className="p-8 space-y-6">
          {sent ? (
            <div className="py-12 text-center animate-in zoom-in-95">
              <div className="w-20 h-20 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto mb-6">
                <CheckCircle className="w-10 h-10" />
              </div>
              <h4 className="text-2xl font-black text-gray-900">
                Message Broadcasted!
              </h4>
              <p className="text-gray-500 mt-2">
                All {event.attendees.length} participants have been notified.
              </p>
            </div>
          ) : (
            <>
              <div className="bg-gray-50 p-6 rounded-3xl flex items-start gap-4 border border-gray-100">
                <Info className="w-5 h-5 text-indigo-500 shrink-0 mt-1" />
                <p className="text-sm text-gray-600 leading-relaxed">
                  Your message will be sent to{" "}
                  <strong>
                    {event.attendees.length} enrolled participants
                  </strong>{" "}
                  via email and in-app notifications.
                </p>
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-black text-gray-400 uppercase tracking-widest pl-1">
                  Message Content
                </label>
                <textarea
                  rows={6}
                  className="w-full px-6 py-4 bg-gray-50 border-none rounded-3xl focus:ring-4 focus:ring-indigo-100 outline-none transition-all resize-none font-medium text-gray-700 leading-relaxed"
                  placeholder="e.g. 'We are starting in 5 minutes! See you there...'"
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  disabled={sending}
                />
              </div>

              <div className="pt-4 flex gap-4">
                <button
                  onClick={onClose}
                  className="flex-1 py-4 text-gray-500 font-bold hover:bg-gray-100 rounded-2xl transition-all"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSend}
                  disabled={!message || sending}
                  className="flex-1 py-4 bg-indigo-600 text-white font-black rounded-2xl hover:bg-indigo-700 shadow-xl shadow-indigo-100 transition-all flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  {sending ? (
                    <Loader2 className="w-5 h-5 animate-spin" />
                  ) : (
                    <Send className="w-5 h-5" />
                  )}
                  {sending ? "Sending..." : "Broadcast Now"}
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default BroadcastModal
