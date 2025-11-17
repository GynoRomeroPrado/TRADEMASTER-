'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Briefcase, Clock, DollarSign, MapPin, CheckCircle, XCircle, AlertCircle, MessageSquare, Eye } from 'lucide-react'

export default function ApplicationsPage() {
  const [statusFilter, setStatusFilter] = useState('all')

  const applications = [
    {
      id: '1',
      jobId: '1',
      jobTitle: 'Data Center Electrical Installation',
      client: 'TechCorp Inc.',
      location: 'San Francisco, CA',
      proposedBudget: 45000,
      timeline: '6 weeks',
      status: 'under_review',
      submittedAt: '2 days ago',
      viewedByClient: true,
      messages: 2,
      proposal: 'With over 8 years of experience in data center electrical installations, I am confident I can deliver exceptional results for your project...',
    },
    {
      id: '2',
      jobId: '4',
      jobTitle: 'Residential Solar Panel Installation',
      client: 'Sarah Johnson',
      location: 'Phoenix, AZ',
      proposedBudget: 17500,
      timeline: '2 weeks',
      status: 'accepted',
      submittedAt: '1 week ago',
      viewedByClient: true,
      messages: 8,
      proposal: 'I specialize in residential solar installations and have completed over 50 similar projects in the Phoenix area...',
      acceptedAt: '5 days ago',
    },
    {
      id: '3',
      jobId: '2',
      jobTitle: 'Commercial HVAC System Installation',
      client: 'BuildCo Properties',
      location: 'Austin, TX',
      proposedBudget: 72000,
      timeline: '8 weeks',
      status: 'rejected',
      submittedAt: '2 weeks ago',
      viewedByClient: true,
      messages: 1,
      proposal: 'Our team has extensive experience with commercial HVAC installations...',
      rejectedAt: '1 week ago',
      rejectionReason: 'Client selected another contractor',
    },
    {
      id: '4',
      jobId: '5',
      jobTitle: 'Emergency HVAC Repair - Hospital',
      client: 'Memorial Hospital',
      location: 'Chicago, IL',
      proposedBudget: 14500,
      timeline: '3 days',
      status: 'under_review',
      submittedAt: '6 hours ago',
      viewedByClient: false,
      messages: 0,
      proposal: 'I can be on-site within 4 hours and have 24/7 availability for the duration of this emergency repair...',
    },
    {
      id: '5',
      jobId: '3',
      jobTitle: 'Structural Steel Welding - Bridge Repair',
      client: 'State DOT',
      location: 'Denver, CO',
      proposedBudget: 35000,
      timeline: '4 weeks',
      status: 'interview_scheduled',
      submittedAt: '3 days ago',
      viewedByClient: true,
      messages: 5,
      proposal: 'AWS D1.5 certified welder with 12 years of bridge and structural steel experience...',
      interviewDate: '2024-02-10 10:00 AM',
    },
  ]

  const filteredApplications = applications.filter((app) => {
    if (statusFilter === 'all') return true
    return app.status === statusFilter
  })

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'under_review':
        return {
          icon: <Clock className="h-4 w-4" />,
          text: 'Under Review',
          className: 'bg-yellow-100 text-yellow-700 border-yellow-300',
        }
      case 'accepted':
        return {
          icon: <CheckCircle className="h-4 w-4" />,
          text: 'Accepted',
          className: 'bg-green-100 text-green-700 border-green-300',
        }
      case 'rejected':
        return {
          icon: <XCircle className="h-4 w-4" />,
          text: 'Rejected',
          className: 'bg-red-100 text-red-700 border-red-300',
        }
      case 'interview_scheduled':
        return {
          icon: <AlertCircle className="h-4 w-4" />,
          text: 'Interview Scheduled',
          className: 'bg-blue-100 text-blue-700 border-blue-300',
        }
      default:
        return {
          icon: <Clock className="h-4 w-4" />,
          text: status,
          className: 'bg-gray-100 text-gray-700 border-gray-300',
        }
    }
  }

  const stats = {
    total: applications.length,
    underReview: applications.filter(a => a.status === 'under_review').length,
    accepted: applications.filter(a => a.status === 'accepted').length,
    rejected: applications.filter(a => a.status === 'rejected').length,
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <Link href="/marketplace" className="inline-flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-4">
            <ArrowLeft className="h-4 w-4" />
            Back to Marketplace
          </Link>
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">My Applications</h1>
              <p className="text-sm text-gray-600">Track your job proposals and applications</p>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Applications</div>
            <div className="text-2xl font-bold text-gray-900">{stats.total}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Under Review</div>
            <div className="text-2xl font-bold text-yellow-600">{stats.underReview}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Accepted</div>
            <div className="text-2xl font-bold text-green-600">{stats.accepted}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Rejected</div>
            <div className="text-2xl font-bold text-red-600">{stats.rejected}</div>
          </div>
        </div>

        {/* Filter */}
        <div className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <div className="flex items-center gap-2">
            <label className="text-sm font-medium text-gray-700">Filter by status:</label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="all">All Applications</option>
              <option value="under_review">Under Review</option>
              <option value="accepted">Accepted</option>
              <option value="rejected">Rejected</option>
              <option value="interview_scheduled">Interview Scheduled</option>
            </select>
          </div>
        </div>

        {/* Applications List */}
        <div className="space-y-4">
          {filteredApplications.map((app) => {
            const statusBadge = getStatusBadge(app.status)
            return (
              <div key={app.id} className="bg-white rounded-lg shadow-sm p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <Link
                        href={`/marketplace/${app.jobId}`}
                        className="text-xl font-semibold text-gray-900 hover:text-primary-600"
                      >
                        {app.jobTitle}
                      </Link>
                      <span className={`px-3 py-1 text-xs font-medium rounded-full border flex items-center gap-1 ${statusBadge.className}`}>
                        {statusBadge.icon}
                        {statusBadge.text}
                      </span>
                    </div>

                    <div className="flex items-center gap-6 text-sm text-gray-600 mb-3">
                      <div className="flex items-center gap-1">
                        <Briefcase className="h-4 w-4" />
                        {app.client}
                      </div>
                      <div className="flex items-center gap-1">
                        <MapPin className="h-4 w-4" />
                        {app.location}
                      </div>
                      <div className="flex items-center gap-1">
                        <DollarSign className="h-4 w-4" />
                        ${app.proposedBudget.toLocaleString()}
                      </div>
                      <div className="flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        {app.timeline}
                      </div>
                    </div>

                    <div className="bg-gray-50 p-4 rounded-lg mb-3">
                      <div className="text-sm text-gray-700 line-clamp-2">{app.proposal}</div>
                    </div>

                    <div className="flex items-center gap-4 text-sm">
                      <div className="flex items-center gap-1 text-gray-600">
                        {app.viewedByClient ? (
                          <>
                            <Eye className="h-4 w-4 text-green-500" />
                            <span>Viewed by client</span>
                          </>
                        ) : (
                          <>
                            <Eye className="h-4 w-4 text-gray-400" />
                            <span>Not viewed yet</span>
                          </>
                        )}
                      </div>
                      {app.messages > 0 && (
                        <div className="flex items-center gap-1 text-primary-600">
                          <MessageSquare className="h-4 w-4" />
                          <span>{app.messages} messages</span>
                        </div>
                      )}
                      <div className="text-gray-500">Submitted {app.submittedAt}</div>
                    </div>
                  </div>

                  <div className="ml-6 flex flex-col gap-2">
                    <Link
                      href={`/marketplace/${app.jobId}`}
                      className="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 text-center"
                    >
                      View Job
                    </Link>
                    {app.messages > 0 && (
                      <button className="px-4 py-2 bg-primary-600 text-white text-sm font-medium rounded-lg hover:bg-primary-700">
                        View Messages
                      </button>
                    )}
                  </div>
                </div>

                {/* Status-specific information */}
                {app.status === 'accepted' && (
                  <div className="p-3 bg-green-50 rounded-lg border border-green-200">
                    <div className="flex items-center gap-2 text-green-700 text-sm">
                      <CheckCircle className="h-4 w-4" />
                      <span className="font-medium">Congratulations! Your proposal was accepted {app.acceptedAt}.</span>
                    </div>
                  </div>
                )}

                {app.status === 'rejected' && app.rejectionReason && (
                  <div className="p-3 bg-red-50 rounded-lg border border-red-200">
                    <div className="text-red-700 text-sm">
                      <span className="font-medium">Reason:</span> {app.rejectionReason}
                    </div>
                  </div>
                )}

                {app.status === 'interview_scheduled' && app.interviewDate && (
                  <div className="p-3 bg-blue-50 rounded-lg border border-blue-200">
                    <div className="flex items-center gap-2 text-blue-700 text-sm">
                      <AlertCircle className="h-4 w-4" />
                      <span className="font-medium">Interview scheduled for {app.interviewDate}</span>
                    </div>
                  </div>
                )}
              </div>
            )
          })}
        </div>

        {filteredApplications.length === 0 && (
          <div className="bg-white p-12 rounded-lg shadow-sm text-center">
            <Briefcase className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No applications found</h3>
            <p className="text-gray-600 mb-6">
              {statusFilter === 'all'
                ? "You haven't submitted any proposals yet"
                : `No applications with status: ${statusFilter.replace('_', ' ')}`}
            </p>
            <Link
              href="/marketplace"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary-600 text-white font-medium rounded-lg hover:bg-primary-700"
            >
              Browse Available Jobs
            </Link>
          </div>
        )}
      </div>
    </div>
  )
}
