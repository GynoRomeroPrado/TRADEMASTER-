'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, MapPin, DollarSign, Clock, Briefcase, Star, CheckCircle, AlertCircle, Building, Calendar } from 'lucide-react'

export default function JobDetailPage({ params }: { params: { id: string } }) {
  const [activeTab, setActiveTab] = useState('details')
  const [proposalText, setProposalText] = useState('')
  const [proposedBudget, setProposedBudget] = useState('')
  const [timeline, setTimeline] = useState('')
  const [showSuccessMessage, setShowSuccessMessage] = useState(false)

  // Mock data
  const job = {
    id: params.id,
    title: 'Data Center Electrical Installation',
    description: 'Need experienced electrician for complete electrical system installation in new data center facility. Must have experience with high-voltage systems and redundant power distribution.',
    longDescription: `This is a comprehensive project requiring installation of a complete electrical infrastructure for a new 10,000 sq ft data center facility. The project includes:

- Main power distribution system installation
- Redundant UPS systems (N+1 configuration)
- Generator backup system integration
- PDU installation and configuration
- Cable management and documentation
- Testing and commissioning

The facility will house critical IT infrastructure and requires 99.99% uptime reliability. All work must meet Uptime Institute Tier III standards and comply with local electrical codes.

The selected contractor will work closely with our project management team and coordinate with other trades including HVAC and network cabling contractors.`,
    trade: 'electrical',
    budget: 45000,
    budgetType: 'fixed',
    location: {
      address: '1234 Technology Drive',
      city: 'San Francisco',
      state: 'CA',
      zip: '94105',
    },
    urgency: 'urgent',
    duration: '6 weeks',
    startDate: '2024-02-15',
    postedAt: '2 days ago',
    client: {
      name: 'TechCorp Inc.',
      rating: 4.8,
      jobsPosted: 23,
      jobsCompleted: 21,
      memberSince: '2022',
      verified: true,
    },
    requirements: [
      'Licensed Master Electrician in California',
      'Minimum 5 years data center electrical experience',
      'Available for full 6-week commitment',
      'Team of at least 3 certified electricians',
      'Liability insurance minimum $2M',
      'References from similar projects',
    ],
    scopeOfWork: [
      'Install 480V main distribution panel',
      'Install two 200kVA UPS systems',
      'Install and wire 250kW backup generator',
      'Install 40 server rack PDUs',
      'Run and terminate all power cables',
      'Label and document all circuits',
      'Perform load testing and commissioning',
    ],
    proposals: 12,
    matchScore: 95,
  }

  const similarJobs = [
    {
      id: '2',
      title: 'Commercial HVAC System Installation',
      location: 'Austin, TX',
      budget: 75000,
      matchScore: 88,
    },
    {
      id: '3',
      title: 'Industrial Electrical Upgrade',
      location: 'San Jose, CA',
      budget: 55000,
      matchScore: 91,
    },
  ]

  const handleSubmitProposal = (e: React.FormEvent) => {
    e.preventDefault()
    // Here would be API call to submit proposal
    setShowSuccessMessage(true)
    setTimeout(() => setShowSuccessMessage(false), 5000)
  }

  const getUrgencyColor = (urgency: string) => {
    switch (urgency) {
      case 'critical': return 'bg-red-100 text-red-700 border-red-300'
      case 'urgent': return 'bg-orange-100 text-orange-700 border-orange-300'
      case 'normal': return 'bg-blue-100 text-blue-700 border-blue-300'
      default: return 'bg-gray-100 text-gray-700 border-gray-300'
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-gradient-to-r from-primary-600 to-blue-600 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Link href="/marketplace" className="inline-flex items-center gap-2 text-white/80 hover:text-white mb-4">
            <ArrowLeft className="h-4 w-4" />
            Back to Marketplace
          </Link>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2">
              <div className="flex items-center gap-3 mb-3">
                <span className={`px-3 py-1 text-sm font-medium rounded-full border capitalize ${getUrgencyColor(job.urgency)}`}>
                  {job.urgency}
                </span>
                <span className="px-3 py-1 bg-white/20 text-white text-sm font-medium rounded-full capitalize">
                  {job.trade}
                </span>
              </div>

              <h1 className="text-3xl font-bold mb-4">{job.title}</h1>
              <p className="text-white/90 text-lg mb-6">{job.description}</p>

              <div className="flex items-center gap-6 text-sm">
                <div className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  {job.location.city}, {job.location.state}
                </div>
                <div className="flex items-center gap-2">
                  <DollarSign className="h-5 w-5" />
                  ${job.budget.toLocaleString()} {job.budgetType}
                </div>
                <div className="flex items-center gap-2">
                  <Clock className="h-5 w-5" />
                  {job.duration}
                </div>
                <div className="flex items-center gap-2">
                  <Calendar className="h-5 w-5" />
                  Starts {job.startDate}
                </div>
              </div>
            </div>

            {/* AI Match Score Card */}
            <div className="bg-white rounded-lg p-6 text-gray-900">
              <div className="text-center mb-4">
                <div className="text-5xl font-bold text-green-600 mb-2">{job.matchScore}%</div>
                <div className="text-sm text-gray-600">AI Match Score</div>
                <div className="text-xs text-gray-500 mt-1">Excellent match for your profile</div>
              </div>

              <div className="space-y-2 mb-6 text-sm">
                <div className="flex items-center gap-2 text-green-600">
                  <CheckCircle className="h-4 w-4" />
                  Trade expertise matches
                </div>
                <div className="flex items-center gap-2 text-green-600">
                  <CheckCircle className="h-4 w-4" />
                  Location preference matches
                </div>
                <div className="flex items-center gap-2 text-green-600">
                  <CheckCircle className="h-4 w-4" />
                  Budget range matches
                </div>
                <div className="flex items-center gap-2 text-green-600">
                  <CheckCircle className="h-4 w-4" />
                  Experience level matches
                </div>
              </div>

              <button
                onClick={() => setActiveTab('apply')}
                className="w-full px-6 py-3 bg-primary-600 text-white font-semibold rounded-lg hover:bg-primary-700"
              >
                Submit Proposal
              </button>

              <div className="mt-4 text-center text-sm text-gray-600">
                {job.proposals} contractors have already applied
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Tabs */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-8">
            {['details', 'requirements', 'client', 'apply'].map((tab) => (
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
        {showSuccessMessage && (
          <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-3">
            <CheckCircle className="h-5 w-5 text-green-600" />
            <div>
              <div className="font-medium text-green-900">Proposal submitted successfully!</div>
              <div className="text-sm text-green-700">The client will review your proposal and contact you if interested.</div>
            </div>
          </div>
        )}

        {activeTab === 'details' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 space-y-6">
              {/* Project Description */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Project Description</h2>
                <div className="prose prose-sm max-w-none text-gray-700 whitespace-pre-line">
                  {job.longDescription}
                </div>
              </div>

              {/* Scope of Work */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Scope of Work</h2>
                <div className="space-y-2">
                  {job.scopeOfWork.map((item, index) => (
                    <div key={index} className="flex items-start gap-3">
                      <CheckCircle className="h-5 w-5 text-primary-600 mt-0.5 flex-shrink-0" />
                      <span className="text-gray-700">{item}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              {/* Project Details */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h3 className="font-semibold text-gray-900 mb-4">Project Details</h3>
                <div className="space-y-3 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Budget</span>
                    <span className="font-medium text-gray-900">${job.budget.toLocaleString()}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Duration</span>
                    <span className="font-medium text-gray-900">{job.duration}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Start Date</span>
                    <span className="font-medium text-gray-900">{job.startDate}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Location</span>
                    <span className="font-medium text-gray-900">{job.location.city}, {job.location.state}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Posted</span>
                    <span className="font-medium text-gray-900">{job.postedAt}</span>
                  </div>
                </div>
              </div>

              {/* Similar Jobs */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h3 className="font-semibold text-gray-900 mb-4">Similar Jobs</h3>
                <div className="space-y-3">
                  {similarJobs.map((similarJob) => (
                    <Link
                      key={similarJob.id}
                      href={`/marketplace/${similarJob.id}`}
                      className="block p-3 rounded-lg hover:bg-gray-50 border border-gray-200"
                    >
                      <div className="font-medium text-gray-900 text-sm mb-1">{similarJob.title}</div>
                      <div className="text-xs text-gray-600 mb-2">{similarJob.location}</div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm font-medium text-gray-900">${similarJob.budget.toLocaleString()}</span>
                        <span className="text-xs text-green-600">{similarJob.matchScore}% match</span>
                      </div>
                    </Link>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'requirements' && (
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-xl font-semibold text-gray-900 mb-6">Requirements</h2>
            <div className="space-y-4">
              {job.requirements.map((req, index) => (
                <div key={index} className="flex items-start gap-3 p-4 bg-gray-50 rounded-lg">
                  <AlertCircle className="h-5 w-5 text-primary-600 mt-0.5 flex-shrink-0" />
                  <span className="text-gray-700">{req}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'client' && (
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-xl font-semibold text-gray-900 mb-6">About the Client</h2>
            <div className="flex items-start gap-6">
              <div className="w-20 h-20 bg-gray-200 rounded-full flex items-center justify-center">
                <Building className="h-10 w-10 text-gray-500" />
              </div>
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-2">
                  <h3 className="text-xl font-semibold text-gray-900">{job.client.name}</h3>
                  {job.client.verified && (
                    <CheckCircle className="h-5 w-5 text-green-500" />
                  )}
                </div>
                <div className="flex items-center gap-2 mb-4">
                  <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                  <span className="font-medium text-gray-900">{job.client.rating}</span>
                  <span className="text-gray-600">• Member since {job.client.memberSince}</span>
                </div>

                <div className="grid grid-cols-3 gap-6 mb-6">
                  <div>
                    <div className="text-2xl font-bold text-gray-900">{job.client.jobsPosted}</div>
                    <div className="text-sm text-gray-600">Jobs Posted</div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-green-600">{job.client.jobsCompleted}</div>
                    <div className="text-sm text-gray-600">Jobs Completed</div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-primary-600">
                      {Math.round((job.client.jobsCompleted / job.client.jobsPosted) * 100)}%
                    </div>
                    <div className="text-sm text-gray-600">Completion Rate</div>
                  </div>
                </div>

                <div className="p-4 bg-blue-50 rounded-lg">
                  <div className="text-sm text-blue-900">
                    This client has a strong track record of completing projects and paying on time.
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'apply' && (
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-xl font-semibold text-gray-900 mb-6">Submit Your Proposal</h2>
            <form onSubmit={handleSubmitProposal} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Cover Letter / Proposal
                </label>
                <textarea
                  value={proposalText}
                  onChange={(e) => setProposalText(e.target.value)}
                  rows={8}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Explain why you're the best fit for this project. Include your relevant experience, approach to the project, and any questions you have..."
                  required
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Your Proposed Budget ($)
                  </label>
                  <input
                    type="number"
                    value={proposedBudget}
                    onChange={(e) => setProposedBudget(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    placeholder="45000"
                    required
                  />
                  <p className="mt-1 text-sm text-gray-500">
                    Client budget: ${job.budget.toLocaleString()}
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Estimated Timeline
                  </label>
                  <input
                    type="text"
                    value={timeline}
                    onChange={(e) => setTimeline(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    placeholder="6 weeks"
                    required
                  />
                  <p className="mt-1 text-sm text-gray-500">
                    Client expects: {job.duration}
                  </p>
                </div>
              </div>

              <div className="p-4 bg-yellow-50 rounded-lg">
                <div className="flex items-start gap-3">
                  <AlertCircle className="h-5 w-5 text-yellow-600 mt-0.5" />
                  <div className="text-sm text-yellow-900">
                    <strong>Before submitting:</strong> Make sure you meet all the requirements and can commit to the timeline. The client will review your profile, ratings, and past work.
                  </div>
                </div>
              </div>

              <div className="flex gap-4">
                <button
                  type="submit"
                  className="px-6 py-3 bg-primary-600 text-white font-semibold rounded-lg hover:bg-primary-700"
                >
                  Submit Proposal
                </button>
                <button
                  type="button"
                  onClick={() => setActiveTab('details')}
                  className="px-6 py-3 border border-gray-300 text-gray-700 font-semibold rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}
      </main>
    </div>
  )
}
