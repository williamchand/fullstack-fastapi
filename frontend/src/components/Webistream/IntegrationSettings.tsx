import type { User } from "@/types/webistream"
import {
  BellRing,
  Bot,
  CheckCircle2,
  CreditCard,
  ExternalLink,
  Loader2,
  RefreshCcw,
  Settings2,
  ShieldCheck,
} from "lucide-react"
import type React from "react"
import { useState } from "react"

interface IntegrationSettingsProps {
  currentUser: User
  onConnectStripe: () => void
}

const IntegrationSettings: React.FC<IntegrationSettingsProps> = ({
  currentUser,
  onConnectStripe,
}) => {
  const [connecting, setConnecting] = useState(false)

  const handleStripeConnect = () => {
    setConnecting(true)
    setTimeout(() => {
      onConnectStripe()
      setConnecting(false)
    }, 1500)
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in slide-in-from-bottom-4 duration-500">
      <div className="space-y-2">
        <h2 className="text-3xl font-bold text-gray-900">
          Integrations & Settings
        </h2>
        <p className="text-gray-500">
          Connect your webinar platforms and configure automated tools.
        </p>
      </div>

      <div className="bg-[#635BFF]/5 rounded-[2.5rem] p-10 border border-[#635BFF]/10 flex flex-col md:flex-row items-center gap-10">
        <div className="bg-[#635BFF] p-6 rounded-[2rem] shadow-2xl shadow-[#635BFF]/30">
          <CreditCard className="w-12 h-12 text-white" />
        </div>
        <div className="flex-1 space-y-3">
          <h3 className="text-2xl font-black text-gray-900 tracking-tight">
            Stripe Onboarding
          </h3>
          <p className="text-gray-500 font-medium leading-relaxed">
            Connect your Stripe account to start collecting payments for your
            webinars. This is required for publishing paid events.
          </p>
          <div className="pt-4">
            {currentUser.isStripeConnected ? (
              <div className="flex items-center gap-4">
                <div className="px-5 py-3 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center gap-2 font-black text-sm">
                  <CheckCircle2 className="w-5 h-5" /> Account Connected
                </div>
                <button
                  type="button"
                  className="text-gray-400 hover:text-gray-600 font-bold text-sm flex items-center gap-2"
                >
                  <RefreshCcw className="w-4 h-4" /> Refresh Status
                </button>
              </div>
            ) : (
              <button
                onClick={handleStripeConnect}
                disabled={connecting}
                className="px-8 py-4 bg-[#635BFF] text-white rounded-2xl font-black shadow-xl shadow-[#635BFF]/20 hover:bg-[#5249db] transition-all flex items-center gap-3 disabled:opacity-50"
              >
                {connecting ? (
                  <Loader2 className="w-5 h-5 animate-spin" />
                ) : (
                  <ExternalLink className="w-5 h-5" />
                )}
                Connect with Stripe
              </button>
            )}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <PlatformCard
          name="Zoom"
          logo="https://upload.wikimedia.org/wikipedia/commons/thumb/9/94/Zoom_Logo.svg/2560px-Zoom_Logo.svg.png"
          status="Connected"
          description="Manage meetings, auto-generate links, and sync attendees."
          connected
        />
        <PlatformCard
          name="Google Meet"
          logo="https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Google_Meet_icon_%282020%29.svg/1024px-Google_Meet_icon_%282020%29.svg.png"
          status="Not Connected"
          description="Schedule events directly through your Google Calendar."
          connected={false}
        />
      </div>

      <div className="bg-white rounded-[2.5rem] border border-gray-100 shadow-sm p-10 space-y-8">
        <h3 className="text-xl font-black text-gray-900 border-b border-gray-50 pb-6">
          Automated Platform Tools
        </h3>
        <div className="grid grid-cols-1 gap-8">
          <AutomationToggle
            icon={<Bot className="w-6 h-6 text-indigo-600" />}
            title="AI Webinar Assistant"
            description="Our bot automatically joins meetings as a co-host to manage entry, record sessions, and generate post-event summaries."
            enabled
          />
          <AutomationToggle
            icon={<ShieldCheck className="w-6 h-6 text-emerald-600" />}
            title="Instant Invitation Delivery"
            description="Send personalized secure meeting links to attendees immediately after successful payment verification."
            enabled
          />
          <AutomationToggle
            icon={<BellRing className="w-6 h-6 text-amber-600" />}
            title="Smart Reminders"
            description="Automated email and SMS notifications 24 hours and 1 hour before the webinar starts."
            enabled={false}
          />
        </div>
      </div>
    </div>
  )
}

const PlatformCard: React.FC<{
  name: string
  logo: string
  status: string
  description: string
  connected: boolean
}> = ({ name, logo, status, description, connected }) => (
  <div className="bg-white p-8 rounded-[2.5rem] border border-gray-100 shadow-sm hover:shadow-md transition-all flex flex-col">
    <div className="flex items-center justify-between mb-8">
      <div className="h-10 w-32 relative">
        <img src={logo} alt={name} className="h-full object-contain" />
      </div>
      <span
        className={`flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest ${
          connected
            ? "bg-emerald-50 text-emerald-600"
            : "bg-gray-100 text-gray-400"
        }`}
      >
        {connected && <CheckCircle2 className="w-3 h-3" />}
        {status}
      </span>
    </div>
    <p className="text-gray-500 text-sm mb-8 flex-1 leading-relaxed">
      {description}
    </p>
    <button
      className={`w-full py-4 rounded-2xl font-black text-sm transition-all flex items-center justify-center gap-2 ${
        connected
          ? "bg-gray-50 text-gray-700 hover:bg-gray-100"
          : "bg-indigo-600 text-white hover:bg-indigo-700"
      }`}
    >
      {connected ? (
        <Settings2 className="w-5 h-5" />
      ) : (
        <ExternalLink className="w-5 h-5" />
      )}
      {connected ? "Configure" : "Connect Account"}
    </button>
  </div>
)

const AutomationToggle: React.FC<{
  icon: React.ReactNode
  title: string
  description: string
  enabled: boolean
}> = ({ icon, title, description, enabled }) => (
  <div className="flex items-start gap-6 group">
    <div className="p-4 bg-gray-50 rounded-2xl group-hover:bg-white group-hover:shadow-lg group-hover:shadow-gray-100 transition-all">
      {icon}
    </div>
    <div className="flex-1 space-y-1">
      <div className="flex items-center justify-between">
        <h4 className="font-black text-gray-900 tracking-tight">{title}</h4>
        <button
          className={`relative w-12 h-7 transition-colors rounded-full focus:outline-none ${
            enabled ? "bg-indigo-600" : "bg-gray-200"
          }`}
        >
          <span
            className={`absolute top-1 left-1 w-5 h-5 transition-transform bg-white rounded-full shadow-sm ${
              enabled ? "translate-x-5" : "translate-x-0"
            }`}
          />
        </button>
      </div>
      <p className="text-sm text-gray-500 leading-relaxed font-medium">
        {description}
      </p>
    </div>
  </div>
)

export default IntegrationSettings
