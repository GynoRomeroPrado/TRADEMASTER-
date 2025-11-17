'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Calendar, DollarSign, MapPin, Users, CheckCircle, Clock, Camera, Plus } from 'lucide-react'

export default function ProjectDetailPage({ params }: { params: { id: string } }) {
  const [activeTab, setActiveTab] = useState('overview')

  // Mock data - in production, fetch from API
  const project = {
    id: params.id,
    name: 'Data Center Wiring - Phase 2',
    description: 'Complete electrical wiring for new data center facility including power distribution, backup systems, and network cabling.',
    client: {
      name: 'Tech Corp',
      email: 'contact@techcorp.com',
      phone: '(555) 123-4567',
      address: '123 Tech Street, Northern Virginia',
    },
    status: 'in_progress',
    progress: 65,
    budget: 125000,
    actualCost: 78000,
    estimatedHours: 800,
    actualHours: 520,
    startDate: '2025-01-15',
    endDate: '2025-03-30',
    location: 'Northern Virginia',
    locationAddress: '456 Data Center Blvd, Ashburn, VA 20147',
    trade: 'electrical',
    teamMembers: [
      { id: '1', name: 'John Smith', role: 'lead', avatar: null },
      { id: '2', name: 'Jane Doe', role: 'technician', avatar: null },
      { id: '3', name: 'Mike Johnson', role: 'helper', avatar: null },
    ],
    tasks: [
      { id: '1', title: 'Main power distribution', status: 'completed', dueDate: '2025-01-25' },
      { id: '2', title: 'Backup power systems', status: 'in_progress', dueDate: '2025-02-10' },
      { id: '3', title: 'Network cabling', status: 'pending', dueDate: '2025-02-25' },
      { id: '4', title: 'Final inspection', status: 'pending', dueDate: '2025-03-28' },
    ],
    photos: [
      { id: '1', url: '/placeholder.jpg', caption: 'Initial site setup', date: '2025-01-15' },
      { id: '2', url: '/placeholder.jpg', caption: 'Power distribution panel', date: '2025-01-20' },
    ],
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'bg-green-100 text-green-700'
      case 'in_progress': return 'bg-blue-100 text-blue-700'
      case 'pending': return 'bg-yellow-100 text-yellow-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link href="/projects" className="p-2 hover:bg-gray-100 rounded-lg">
                <ArrowLeft className="h-5 w-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{project.name}</h1>
                <div className="flex items-center gap-2 mt-1">
                  <span className={`px-3 py-1 text-xs font-medium rounded-full ${getStatusColor(project.status)}`}>
                    {project.status.replace('_', ' ')}
                  </span>
                  <span className="text-sm text-gray-600">{project.client.name}</span>
                </div>
              </div>
            </div>
            <div className="flex gap-2">
              <button className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50">
                Edit
              </button>
              <button className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
                Add Task
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Tabs */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-8">
            {['overview', 'tasks', 'team', 'photos', 'materials'].map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`py-4 border-b-2 font-medium text-sm capitalize ${
                  activeTab === tab
                    ? 'border-primary-600 text-primary-600'
                    : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
              >
                {tab}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'overview' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-6">
              {/* Progress */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Project Progress</h2>
                <div className="mb-2">
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-gray-600">Overall Completion</span>
                    <span className="font-medium text-gray-900">{project.progress}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-3">
                    <div
                      className="bg-primary-600 h-3 rounded-full transition-all"
                      style={{ width: `${project.progress}%` }}
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4 mt-4">
                  <div>
                    <div className="text-sm text-gray-600">Budget Usage</div>
                    <div className="text-xl font-bold text-gray-900">
                      ${project.actualCost.toLocaleString()} / ${project.budget.toLocaleString()}
                    </div>
                    <div className="text-sm text-gray-600">
                      {Math.round((project.actualCost / project.budget) * 100)}% spent
                    </div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600">Hours Logged</div>
                    <div className="text-xl font-bold text-gray-900">
                      {project.actualHours} / {project.estimatedHours}
                    </div>
                    <div className="text-sm text-gray-600">
                      {Math.round((project.actualHours / project.estimatedHours) * 100)}% complete
                    </div>
                  </div>
                </div>
              </div>

              {/* Description */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Description</h2>
                <p className="text-gray-700">{project.description}</p>
              </div>

              {/* Recent Activity */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h2>
                <div className="space-y-4">
                  <div className="flex gap-3">
                    <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
                    <div>
                      <div className="font-medium text-gray-900">Task completed</div>
                      <div className="text-sm text-gray-600">Main power distribution completed by John Smith</div>
                      <div className="text-xs text-gray-500 mt-1">2 hours ago</div>
                    </div>
                  </div>
                  <div className="flex gap-3">
                    <Camera className="h-5 w-5 text-blue-500 mt-0.5" />
                    <div>
                      <div className="font-medium text-gray-900">Photo uploaded</div>
                      <div className="text-sm text-gray-600">Power distribution panel</div>
                      <div className="text-xs text-gray-500 mt-1">5 hours ago</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              {/* Details Card */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Details</h2>
                <div className="space-y-3">
                  <div className="flex items-start gap-3">
                    <Calendar className="h-5 w-5 text-gray-400 mt-0.5" />
                    <div>
                      <div className="text-sm text-gray-600">Timeline</div>
                      <div className="font-medium text-gray-900">
                        {new Date(project.startDate).toLocaleDateString()} - {new Date(project.endDate).toLocaleDateString()}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <MapPin className="h-5 w-5 text-gray-400 mt-0.5" />
                    <div>
                      <div className="text-sm text-gray-600">Location</div>
                      <div className="font-medium text-gray-900">{project.locationAddress}</div>
                    </div>
                  </div>
                  <div className="flex items-start gap-3">
                    <DollarSign className="h-5 w-5 text-gray-400 mt-0.5" />
                    <div>
                      <div className="text-sm text-gray-600">Budget</div>
                      <div className="font-medium text-gray-900">${project.budget.toLocaleString()}</div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Client Info */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Client</h2>
                <div className="space-y-2">
                  <div className="font-medium text-gray-900">{project.client.name}</div>
                  <div className="text-sm text-gray-600">{project.client.email}</div>
                  <div className="text-sm text-gray-600">{project.client.phone}</div>
                  <div className="text-sm text-gray-600">{project.client.address}</div>
                </div>
              </div>

              {/* Team */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Team</h2>
                <div className="space-y-3">
                  {project.teamMembers.map((member) => (
                    <div key={member.id} className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                        <Users className="h-5 w-5 text-gray-600" />
                      </div>
                      <div>
                        <div className="font-medium text-gray-900">{member.name}</div>
                        <div className="text-sm text-gray-600 capitalize">{member.role}</div>
                      </div>
                    </div>
                  ))}
                  <button className="w-full mt-2 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 flex items-center justify-center gap-2">
                    <Plus className="h-4 w-4" />
                    Add Member
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'tasks' && (
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-semibold text-gray-900">Tasks</h2>
              <button className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700">
                Add Task
              </button>
            </div>
            <div className="space-y-3">
              {project.tasks.map((task) => (
                <div key={task.id} className="flex items-center gap-4 p-4 border border-gray-200 rounded-lg">
                  <input type="checkbox" checked={task.status === 'completed'} className="h-5 w-5" />
                  <div className="flex-1">
                    <div className="font-medium text-gray-900">{task.title}</div>
                    <div className="text-sm text-gray-600 flex items-center gap-2">
                      <Clock className="h-4 w-4" />
                      Due {new Date(task.dueDate).toLocaleDateString()}
                    </div>
                  </div>
                  <span className={`px-3 py-1 text-xs font-medium rounded-full ${getStatusColor(task.status)}`}>
                    {task.status.replace('_', ' ')}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  )
}
