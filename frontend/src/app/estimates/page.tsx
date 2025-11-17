'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Plus, Search, FileText, CheckCircle, XCircle, Clock } from 'lucide-react'

export default function EstimatesPage() {
  const [searchQuery, setSearchQuery] = useState('')

  const estimates = [
    {
      id: '1',
      projectName: 'Data Center Wiring - Phase 2',
      client: { name: 'Tech Corp', email: 'contact@techcorp.com' },
      status: 'accepted',
      total: 125000,
      confidenceScore: 0.92,
      createdAt: '2025-01-10',
      validUntil: '2025-02-10',
    },
    {
      id: '2',
      projectName: 'HVAC Installation - Building A',
      client: { name: 'Real Estate LLC', email: 'info@realesta

te.com' },
      status: 'sent',
      total: 85000,
      confidenceScore: 0.88,
      createdAt: '2025-01-12',
      validUntil: '2025-02-12',
    },
    {
      id: '3',
      projectName: 'Solar Panel Installation',
      client: { name: 'Green Energy Inc', email: 'solar@green.com' },
      status: 'draft',
      total: 95000,
      confidenceScore: 0.85,
      createdAt: '2025-01-14',
      validUntil: '2025-02-14',
    },
  ]

  const filteredEstimates = estimates.filter((est) =>
    est.projectName.toLowerCase().includes(searchQuery.toLowerCase()) ||
    est.client.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'accepted': return <CheckCircle className="h-5 w-5 text-green-500" />
      case 'rejected': return <XCircle className="h-5 w-5 text-red-500" />
      case 'sent': return <Clock className="h-5 w-5 text-blue-500" />
      default: return <FileText className="h-5 w-5 text-gray-400" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'accepted': return 'bg-green-100 text-green-700'
      case 'rejected': return 'bg-red-100 text-red-700'
      case 'sent': return 'bg-blue-100 text-blue-700'
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
              <h1 className="text-2xl font-bold text-gray-900">Estimates</h1>
              <p className="text-sm text-gray-600">Create and manage project estimates</p>
            </div>
            <Link
              href="/estimates/new"
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
            >
              <Plus className="h-5 w-5" />
              New Estimate
            </Link>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Search */}
        <div className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search estimates..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            />
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Estimates</div>
            <div className="text-2xl font-bold text-gray-900">{estimates.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Accepted</div>
            <div className="text-2xl font-bold text-green-600">
              {estimates.filter(e => e.status === 'accepted').length}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Pending</div>
            <div className="text-2xl font-bold text-blue-600">
              {estimates.filter(e => e.status === 'sent').length}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Value</div>
            <div className="text-2xl font-bold text-gray-900">
              ${estimates.reduce((sum, e) => sum + e.total, 0).toLocaleString()}
            </div>
          </div>
        </div>

        {/* Estimates List */}
        <div className="space-y-4">
          {filteredEstimates.map((estimate) => (
            <Link
              key={estimate.id}
              href={`/estimates/${estimate.id}`}
              className="block bg-white p-6 rounded-lg shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-start gap-3">
                  {getStatusIcon(estimate.status)}
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">{estimate.projectName}</h3>
                    <div className="flex items-center gap-2 text-sm text-gray-600 mt-1">
                      <span>{estimate.client.name}</span>
                      <span>•</span>
                      <span>{estimate.client.email}</span>
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-2xl font-bold text-gray-900">${estimate.total.toLocaleString()}</div>
                  <span className={`inline-block mt-1 px-3 py-1 text-xs font-medium rounded-full ${getStatusColor(estimate.status)}`}>
                    {estimate.status}
                  </span>
                </div>
              </div>

              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-4">
                  <div>
                    <span className="text-gray-600">Confidence: </span>
                    <span className="font-medium text-gray-900">{(estimate.confidenceScore * 100).toFixed(0)}%</span>
                  </div>
                  <div>
                    <span className="text-gray-600">Created: </span>
                    <span className="font-medium text-gray-900">{new Date(estimate.createdAt).toLocaleDateString()}</span>
                  </div>
                </div>
                <div>
                  <span className="text-gray-600">Valid until: </span>
                  <span className="font-medium text-gray-900">{new Date(estimate.validUntil).toLocaleDateString()}</span>
                </div>
              </div>
            </Link>
          ))}

          {filteredEstimates.length === 0 && (
            <div className="bg-white p-12 rounded-lg shadow-sm text-center">
              <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No estimates found</h3>
              <p className="text-gray-600 mb-4">Create your first estimate to get started</p>
              <Link
                href="/estimates/new"
                className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
              >
                <Plus className="h-5 w-5" />
                Create Estimate
              </Link>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}
