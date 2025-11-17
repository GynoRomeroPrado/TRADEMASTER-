'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Search, MapPin, DollarSign, Clock, Briefcase, Star, Filter } from 'lucide-react'

export default function MarketplacePage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [tradeFilter, setTradeFilter] = useState('all')
  const [budgetFilter, setBudgetFilter] = useState('all')
  const [urgencyFilter, setUrgencyFilter] = useState('all')

  const jobs = [
    {
      id: '1',
      title: 'Data Center Electrical Installation',
      description: 'Need experienced electrician for complete electrical system installation in new data center facility. Must have experience with high-voltage systems and redundant power distribution.',
      trade: 'electrical',
      budget: 45000,
      budgetType: 'fixed',
      location: {
        city: 'San Francisco',
        state: 'CA',
      },
      urgency: 'urgent',
      duration: '6 weeks',
      postedAt: '2 days ago',
      client: {
        name: 'TechCorp Inc.',
        rating: 4.8,
        jobsPosted: 23,
      },
      requirements: [
        'Licensed Master Electrician',
        'Data center experience required',
        'Available for 6-week commitment',
      ],
      proposals: 12,
      matchScore: 95,
    },
    {
      id: '2',
      title: 'Commercial HVAC System Installation',
      description: 'Installation of new HVAC system in 50,000 sq ft commercial building. Includes ductwork, controls, and integration with building management system.',
      trade: 'hvac',
      budget: 75000,
      budgetType: 'fixed',
      location: {
        city: 'Austin',
        state: 'TX',
      },
      urgency: 'normal',
      duration: '8 weeks',
      postedAt: '5 days ago',
      client: {
        name: 'BuildCo Properties',
        rating: 4.6,
        jobsPosted: 15,
      },
      requirements: [
        'EPA certification',
        'Commercial HVAC experience',
        'BMS integration experience',
      ],
      proposals: 8,
      matchScore: 88,
    },
    {
      id: '3',
      title: 'Structural Steel Welding - Bridge Repair',
      description: 'Precision welding required for structural steel repairs on highway bridge. Must meet DOT specifications and pass x-ray inspections.',
      trade: 'welding',
      budget: 35000,
      budgetType: 'fixed',
      location: {
        city: 'Denver',
        state: 'CO',
      },
      urgency: 'urgent',
      duration: '4 weeks',
      postedAt: '1 day ago',
      client: {
        name: 'State DOT',
        rating: 4.9,
        jobsPosted: 45,
      },
      requirements: [
        'AWS D1.5 certified',
        'Bridge welding experience',
        'Pass weld qualification tests',
      ],
      proposals: 5,
      matchScore: 92,
    },
    {
      id: '4',
      title: 'Residential Solar Panel Installation',
      description: 'Install 8.5kW solar panel system on residential property. Includes electrical work, mounting, and grid connection.',
      trade: 'electrical',
      budget: 18000,
      budgetType: 'fixed',
      location: {
        city: 'Phoenix',
        state: 'AZ',
      },
      urgency: 'normal',
      duration: '2 weeks',
      postedAt: '3 days ago',
      client: {
        name: 'Sarah Johnson',
        rating: 5.0,
        jobsPosted: 1,
      },
      requirements: [
        'Solar installation certified',
        'Licensed electrician',
        'Insurance required',
      ],
      proposals: 15,
      matchScore: 85,
    },
    {
      id: '5',
      title: 'Emergency HVAC Repair - Hospital',
      description: 'Urgent repair needed for hospital HVAC system. Multiple units down affecting patient care areas. Must be available immediately.',
      trade: 'hvac',
      budget: 15000,
      budgetType: 'hourly',
      location: {
        city: 'Chicago',
        state: 'IL',
      },
      urgency: 'critical',
      duration: '3 days',
      postedAt: '4 hours ago',
      client: {
        name: 'Memorial Hospital',
        rating: 4.7,
        jobsPosted: 67,
      },
      requirements: [
        'Available immediately',
        'Hospital HVAC experience',
        '24/7 availability for duration',
      ],
      proposals: 3,
      matchScore: 78,
    },
  ]

  const filteredJobs = jobs.filter((job) => {
    const matchesSearch = job.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         job.description.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesTrade = tradeFilter === 'all' || job.trade === tradeFilter
    const matchesBudget = budgetFilter === 'all' ||
                         (budgetFilter === 'under20k' && job.budget < 20000) ||
                         (budgetFilter === '20k-50k' && job.budget >= 20000 && job.budget < 50000) ||
                         (budgetFilter === 'over50k' && job.budget >= 50000)
    const matchesUrgency = urgencyFilter === 'all' || job.urgency === urgencyFilter
    return matchesSearch && matchesTrade && matchesBudget && matchesUrgency
  })

  const getUrgencyColor = (urgency: string) => {
    switch (urgency) {
      case 'critical': return 'bg-red-100 text-red-700 border-red-300'
      case 'urgent': return 'bg-orange-100 text-orange-700 border-orange-300'
      case 'normal': return 'bg-blue-100 text-blue-700 border-blue-300'
      default: return 'bg-gray-100 text-gray-700 border-gray-300'
    }
  }

  const getMatchScoreColor = (score: number) => {
    if (score >= 90) return 'text-green-600'
    if (score >= 80) return 'text-blue-600'
    if (score >= 70) return 'text-yellow-600'
    return 'text-gray-600'
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Job Marketplace</h1>
              <p className="text-sm text-gray-600">Find your next project</p>
            </div>
            <Link
              href="/marketplace/applications"
              className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
            >
              <Briefcase className="h-5 w-5" />
              My Applications
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
                placeholder="Search jobs..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
            <div className="flex gap-2">
              <select
                value={tradeFilter}
                onChange={(e) => setTradeFilter(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="all">All Trades</option>
                <option value="electrical">Electrical</option>
                <option value="hvac">HVAC</option>
                <option value="welding">Welding</option>
              </select>
              <select
                value={budgetFilter}
                onChange={(e) => setBudgetFilter(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="all">All Budgets</option>
                <option value="under20k">Under $20k</option>
                <option value="20k-50k">$20k - $50k</option>
                <option value="over50k">Over $50k</option>
              </select>
              <select
                value={urgencyFilter}
                onChange={(e) => setUrgencyFilter(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="all">All Urgency</option>
                <option value="critical">Critical</option>
                <option value="urgent">Urgent</option>
                <option value="normal">Normal</option>
              </select>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Available Jobs</div>
            <div className="text-2xl font-bold text-gray-900">{filteredJobs.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Avg Budget</div>
            <div className="text-2xl font-bold text-primary-600">
              ${Math.round(jobs.reduce((sum, j) => sum + j.budget, 0) / jobs.length).toLocaleString()}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Proposals</div>
            <div className="text-2xl font-bold text-gray-900">
              {jobs.reduce((sum, j) => sum + j.proposals, 0)}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Urgent Jobs</div>
            <div className="text-2xl font-bold text-orange-600">
              {jobs.filter(j => j.urgency === 'urgent' || j.urgency === 'critical').length}
            </div>
          </div>
        </div>

        {/* Jobs List */}
        <div className="space-y-4">
          {filteredJobs.map((job) => (
            <Link
              key={job.id}
              href={`/marketplace/${job.id}`}
              className="block bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow p-6"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-xl font-semibold text-gray-900">{job.title}</h3>
                    <span className={`px-3 py-1 text-xs font-medium rounded-full border capitalize ${getUrgencyColor(job.urgency)}`}>
                      {job.urgency}
                    </span>
                  </div>
                  <p className="text-gray-600 mb-3">{job.description}</p>

                  <div className="flex items-center gap-6 text-sm text-gray-600">
                    <div className="flex items-center gap-1">
                      <MapPin className="h-4 w-4" />
                      {job.location.city}, {job.location.state}
                    </div>
                    <div className="flex items-center gap-1">
                      <DollarSign className="h-4 w-4" />
                      ${job.budget.toLocaleString()} {job.budgetType}
                    </div>
                    <div className="flex items-center gap-1">
                      <Clock className="h-4 w-4" />
                      {job.duration}
                    </div>
                    <div className="flex items-center gap-1">
                      <Briefcase className="h-4 w-4" />
                      {job.proposals} proposals
                    </div>
                  </div>
                </div>

                <div className="ml-6 flex flex-col items-end gap-2">
                  <div className="text-right">
                    <div className={`text-3xl font-bold ${getMatchScoreColor(job.matchScore)}`}>
                      {job.matchScore}%
                    </div>
                    <div className="text-xs text-gray-500">AI Match Score</div>
                  </div>
                  <button className="px-4 py-2 bg-primary-600 text-white text-sm font-medium rounded-lg hover:bg-primary-700">
                    View Details
                  </button>
                </div>
              </div>

              <div className="flex items-center justify-between pt-4 border-t">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                    <span className="text-sm font-bold text-gray-600">
                      {job.client.name.split(' ').map(n => n[0]).join('')}
                    </span>
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">{job.client.name}</div>
                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                      {job.client.rating} • {job.client.jobsPosted} jobs posted
                    </div>
                  </div>
                </div>
                <div className="text-sm text-gray-500">Posted {job.postedAt}</div>
              </div>
            </Link>
          ))}
        </div>

        {filteredJobs.length === 0 && (
          <div className="bg-white p-12 rounded-lg shadow-sm text-center">
            <Briefcase className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No jobs found</h3>
            <p className="text-gray-600">Try adjusting your search or filters</p>
          </div>
        )}
      </div>
    </div>
  )
}
