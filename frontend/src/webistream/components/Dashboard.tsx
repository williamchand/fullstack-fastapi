import type { DashboardStats, WebinarEvent } from "@/webistream/types"
import {
  Activity,
  Calendar as CalendarIcon,
  DollarSign,
  TrendingUp,
  Users,
} from "lucide-react"
import type React from "react"
// Lightweight placeholder chart without external deps

interface DashboardProps {
  stats: DashboardStats
  events: WebinarEvent[]
}

const Dashboard: React.FC<DashboardProps> = ({ stats }) => {
  const sparkline = [40, 48, 45, 60, 58, 72]

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          icon={<DollarSign className="w-6 h-6 text-emerald-600" />}
          label="Total Revenue"
          value={`$${stats.totalRevenue.toLocaleString()}`}
          trend="+12.5%"
          color="bg-emerald-50"
        />
        <StatCard
          icon={<Users className="w-6 h-6 text-blue-600" />}
          label="Total Attendees"
          value={stats.activeAttendees.toString()}
          trend="+8.2%"
          color="bg-blue-50"
        />
        <StatCard
          icon={<CalendarIcon className="w-6 h-6 text-indigo-600" />}
          label="Upcoming Events"
          value={stats.upcomingEventsCount.toString()}
          trend="0"
          color="bg-indigo-50"
        />
        <StatCard
          icon={<TrendingUp className="w-6 h-6 text-amber-600" />}
          label="Conv. Rate"
          value="76%"
          trend="+4.1%"
          color="bg-amber-50"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-lg font-bold text-gray-800">Revenue Trends</h3>
            <Activity className="w-5 h-5 text-gray-400" />
          </div>
          <div className="h-64">
            <svg viewBox="0 0 300 120" className="w-full h-full">
              <defs>
                <linearGradient id="fillIndigo" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#4f46e5" stopOpacity={0.1} />
                  <stop offset="95%" stopColor="#4f46e5" stopOpacity={0} />
                </linearGradient>
              </defs>
              <polyline
                fill="url(#fillIndigo)"
                stroke="#4f46e5"
                strokeWidth="3"
                points={sparkline
                  .map(
                    (v, i) =>
                      `${i * (300 / (sparkline.length - 1))},${120 - v}`,
                  )
                  .join(" ")}
              />
            </svg>
          </div>
        </div>

        <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <h3 className="text-lg font-bold text-gray-800 mb-6">
            Latest Attendees
          </h3>
          <div className="space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div
                key={i}
                className="flex items-center justify-between py-2 border-b border-gray-50 last:border-0"
              >
                <div className="flex items-center gap-3">
                  <img
                    src={`https://picsum.photos/seed/${i}/32/32`}
                    className="w-8 h-8 rounded-full"
                    alt=""
                  />
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      Attendee {i}
                    </p>
                    <p className="text-xs text-gray-500">Joined 2h ago</p>
                  </div>
                </div>
                <span className="px-2 py-1 bg-emerald-50 text-emerald-600 text-xs font-medium rounded-full">
                  Paid
                </span>
              </div>
            ))}
          </div>
          <button className="w-full mt-6 py-2 text-indigo-600 text-sm font-semibold hover:bg-indigo-50 rounded-lg transition-colors">
            View All Attendees
          </button>
        </div>
      </div>
    </div>
  )
}

interface StatCardProps {
  icon: React.ReactNode
  label: string
  value: string
  trend: string
  color: string
}

const StatCard: React.FC<StatCardProps> = ({
  icon,
  label,
  value,
  trend,
  color,
}) => (
  <div className="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-shadow">
    <div className="flex items-start justify-between">
      <div className={`${color} p-3 rounded-xl`}>{icon}</div>
      <span
        className={`text-xs font-bold ${
          trend.startsWith("+") ? "text-emerald-600" : "text-gray-400"
        }`}
      >
        {trend}
      </span>
    </div>
    <div className="mt-4">
      <p className="text-sm text-gray-500 font-medium">{label}</p>
      <h4 className="text-2xl font-bold text-gray-900 mt-1">{value}</h4>
    </div>
  </div>
)

export default Dashboard
