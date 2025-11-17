'use client'

import { BarChart, Calendar, DollarSign, FolderKanban, GraduationCap, Users } from 'lucide-react'
import Link from 'next/link'

export default function DashboardPage() {
  const stats = [
    { name: 'Active Projects', value: '12', change: '+4.75%', icon: FolderKanban, color: 'bg-blue-500' },
    { name: 'Revenue (MTD)', value: '$45,231', change: '+12.3%', icon: DollarSign, color: 'bg-green-500' },
    { name: 'Team Members', value: '8', change: '+2', icon: Users, color: 'bg-purple-500' },
    { name: 'Certifications', value: '24', change: '+3', icon: GraduationCap, color: 'bg-orange-500' },
  ]

  const recentProjects = [
    { id: 1, name: 'Data Center Wiring - Phase 2', client: 'Tech Corp', status: 'in_progress', progress: 65 },
    { id: 2, name: 'HVAC Installation - Building A', client: 'Real Estate LLC', status: 'scheduled', progress: 0 },
    { id: 3, name: 'Solar Panel Installation', client: 'Green Energy Inc', status: 'in_progress', progress: 45 },
  ]

  const upcomingTasks = [
    { id: 1, title: 'Site inspection - Tech Corp', due: 'Today, 2:00 PM', priority: 'high' },
    { id: 2, title: 'Material pickup - Supplier A', due: 'Tomorrow, 9:00 AM', priority: 'medium' },
    { id: 3, title: 'Client meeting - Green Energy', due: 'Thu, 3:00 PM', priority: 'low' },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
          <div className="flex gap-4">
            <Link href="/estimates/new" className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
              New Estimate
            </Link>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8">
          {stats.map((stat) => (
            <div key={stat.name} className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <div className="flex items-center">
                  <div className={`${stat.color} rounded-md p-3`}>
                    <stat.icon className="h-6 w-6 text-white" />
                  </div>
                  <div className="ml-5 w-0 flex-1">
                    <dl>
                      <dt className="text-sm font-medium text-gray-500 truncate">{stat.name}</dt>
                      <dd className="flex items-baseline">
                        <div className="text-2xl font-semibold text-gray-900">{stat.value}</div>
                        <div className="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                          {stat.change}
                        </div>
                      </dd>
                    </dl>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
          {/* Recent Projects */}
          <div className="bg-white shadow rounded-lg">
            <div className="p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Projects</h2>
              <div className="space-y-4">
                {recentProjects.map((project) => (
                  <div key={project.id} className="border-b pb-4 last:border-0">
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <h3 className="font-medium text-gray-900">{project.name}</h3>
                        <p className="text-sm text-gray-500">{project.client}</p>
                      </div>
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        project.status === 'in_progress' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'
                      }`}>
                        {project.status.replace('_', ' ')}
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div className="bg-primary-600 h-2 rounded-full" style={{ width: `${project.progress}%` }} />
                    </div>
                  </div>
                ))}
              </div>
              <Link href="/projects" className="mt-4 block text-center text-primary-600 hover:text-primary-700 text-sm font-medium">
                View all projects →
              </Link>
            </div>
          </div>

          {/* Upcoming Tasks */}
          <div className="bg-white shadow rounded-lg">
            <div className="p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Upcoming Tasks</h2>
              <div className="space-y-4">
                {upcomingTasks.map((task) => (
                  <div key={task.id} className="flex items-center gap-4">
                    <div className={`w-2 h-2 rounded-full ${
                      task.priority === 'high' ? 'bg-red-500' : task.priority === 'medium' ? 'bg-yellow-500' : 'bg-green-500'
                    }`} />
                    <div className="flex-1">
                      <h3 className="font-medium text-gray-900">{task.title}</h3>
                      <p className="text-sm text-gray-500 flex items-center gap-1">
                        <Calendar className="h-4 w-4" />
                        {task.due}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
              <Link href="/tasks" className="mt-4 block text-center text-primary-600 hover:text-primary-700 text-sm font-medium">
                View all tasks →
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
