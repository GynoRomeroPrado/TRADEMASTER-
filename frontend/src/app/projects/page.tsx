'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Plus, Search, Filter, FolderKanban } from 'lucide-react'

export default function ProjectsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')

  const projects = [
    {
      id: '1',
      name: 'Data Center Wiring - Phase 2',
      client: 'Tech Corp',
      status: 'in_progress',
      progress: 65,
      budget: 125000,
      actualCost: 78000,
      startDate: '2025-01-15',
      endDate: '2025-03-30',
      location: 'Northern Virginia',
      trade: 'electrical',
    },
    {
      id: '2',
      name: 'HVAC Installation - Building A',
      client: 'Real Estate LLC',
      status: 'scheduled',
      progress: 0,
      budget: 85000,
      actualCost: 0,
      startDate: '2025-02-01',
      endDate: '2025-04-15',
      location: 'Phoenix, AZ',
      trade: 'hvac',
    },
    {
      id: '3',
      name: 'Solar Panel Installation',
      client: 'Green Energy Inc',
      status: 'in_progress',
      progress: 45,
      budget: 95000,
      actualCost: 42000,
      startDate: '2025-01-10',
      endDate: '2025-02-28',
      location: 'Dallas, TX',
      trade: 'electrical',
    },
    {
      id: '4',
      name: 'Industrial Welding Project',
      client: 'Manufacturing Co',
      status: 'completed',
      progress: 100,
      budget: 150000,
      actualCost: 145000,
      startDate: '2024-11-01',
      endDate: '2025-01-15',
      location: 'Chicago, IL',
      trade: 'welding',
    },
  ]

  const filteredProjects = projects.filter((project) => {
    const matchesSearch = project.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         project.client.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesStatus = statusFilter === 'all' || project.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'bg-green-100 text-green-700'
      case 'in_progress': return 'bg-blue-100 text-blue-700'
      case 'scheduled': return 'bg-yellow-100 text-yellow-700'
      case 'draft': return 'bg-gray-100 text-gray-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Projects</h1>
              <p className="text-sm text-gray-600">Manage all your active and completed projects</p>
            </div>
            <Link
              href="/projects/new"
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
            >
              <Plus className="h-5 w-5" />
              New Project
            </Link>
          </div>
        </div>
      </header>

      {/* Filters */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search projects..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
            <div className="flex gap-2">
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="all">All Status</option>
                <option value="draft">Draft</option>
                <option value="scheduled">Scheduled</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
              </select>
              <button className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">
                <Filter className="h-5 w-5" />
                More Filters
              </button>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Projects</div>
            <div className="text-2xl font-bold text-gray-900">{projects.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">In Progress</div>
            <div className="text-2xl font-bold text-blue-600">
              {projects.filter(p => p.status === 'in_progress').length}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Completed</div>
            <div className="text-2xl font-bold text-green-600">
              {projects.filter(p => p.status === 'completed').length}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Budget</div>
            <div className="text-2xl font-bold text-gray-900">
              ${projects.reduce((sum, p) => sum + p.budget, 0).toLocaleString()}
            </div>
          </div>
        </div>

        {/* Projects List */}
        <div className="space-y-4">
          {filteredProjects.map((project) => (
            <Link
              key={project.id}
              href={`/projects/${project.id}`}
              className="block bg-white p-6 rounded-lg shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <FolderKanban className="h-5 w-5 text-gray-400" />
                    <h3 className="text-lg font-semibold text-gray-900">{project.name}</h3>
                    <span className={`px-3 py-1 text-xs font-medium rounded-full ${getStatusColor(project.status)}`}>
                      {project.status.replace('_', ' ')}
                    </span>
                  </div>
                  <div className="flex items-center gap-4 text-sm text-gray-600">
                    <span>{project.client}</span>
                    <span>•</span>
                    <span>{project.location}</span>
                    <span>•</span>
                    <span className="capitalize">{project.trade}</span>
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-sm text-gray-600">Budget</div>
                  <div className="text-lg font-semibold text-gray-900">${project.budget.toLocaleString()}</div>
                  <div className="text-xs text-gray-500">
                    Spent: ${project.actualCost.toLocaleString()}
                  </div>
                </div>
              </div>

              <div className="mb-2">
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-gray-600">Progress</span>
                  <span className="font-medium text-gray-900">{project.progress}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-primary-600 h-2 rounded-full transition-all"
                    style={{ width: `${project.progress}%` }}
                  />
                </div>
              </div>

              <div className="flex justify-between text-sm text-gray-600">
                <span>Start: {new Date(project.startDate).toLocaleDateString()}</span>
                <span>End: {new Date(project.endDate).toLocaleDateString()}</span>
              </div>
            </Link>
          ))}

          {filteredProjects.length === 0 && (
            <div className="bg-white p-12 rounded-lg shadow-sm text-center">
              <FolderKanban className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No projects found</h3>
              <p className="text-gray-600 mb-4">Try adjusting your search or filters</p>
              <Link
                href="/projects/new"
                className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
              >
                <Plus className="h-5 w-5" />
                Create Your First Project
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
