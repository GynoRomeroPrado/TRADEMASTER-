'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Search, GraduationCap, Clock, Award, Play, BookOpen } from 'lucide-react'

export default function TrainingPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [tradeFilter, setTradeFilter] = useState('all')
  const [levelFilter, setLevelFilter] = useState('all')

  const courses = [
    {
      id: '1',
      title: 'Data Center Electrical Systems',
      description: 'Master the electrical systems used in modern data centers, including power distribution, backup systems, and cooling infrastructure.',
      trade: 'electrical',
      level: 'advanced',
      duration: 12,
      format: 'mixed',
      thumbnail: '/courses/datacenter.jpg',
      certification: true,
      certificationName: 'Data Center Electrician Certified',
      enrolled: 245,
      rating: 4.8,
      modules: 8,
      lessons: 42,
      price: 299,
    },
    {
      id: '2',
      title: 'Solar Panel Installation Fundamentals',
      description: 'Learn the complete process of solar panel installation from site assessment to final inspection and grid connection.',
      trade: 'electrical',
      level: 'intermediate',
      duration: 8,
      format: 'video',
      thumbnail: '/courses/solar.jpg',
      certification: true,
      certificationName: 'Solar Installation Professional',
      enrolled: 432,
      rating: 4.9,
      modules: 6,
      lessons: 28,
      price: 199,
    },
    {
      id: '3',
      title: 'HVAC System Design for Commercial Buildings',
      description: 'Advanced course on designing efficient HVAC systems for commercial applications with focus on energy efficiency.',
      trade: 'hvac',
      level: 'advanced',
      duration: 15,
      format: 'mixed',
      thumbnail: '/courses/hvac-design.jpg',
      certification: true,
      certificationName: 'HVAC Design Specialist',
      enrolled: 178,
      rating: 4.7,
      modules: 10,
      lessons: 52,
      price: 349,
    },
    {
      id: '4',
      title: 'Precision Welding with AR Training',
      description: 'Interactive AR-based welding training for precision structural welding with real-time feedback and coaching.',
      trade: 'welding',
      level: 'intermediate',
      duration: 10,
      format: 'ar',
      thumbnail: '/courses/welding-ar.jpg',
      certification: true,
      certificationName: 'Precision Welder Certified',
      enrolled: 312,
      rating: 4.9,
      modules: 7,
      lessons: 35,
      price: 249,
    },
    {
      id: '5',
      title: 'Electrical Safety & Code Compliance',
      description: 'Essential safety procedures and NEC code compliance for all electrical work. Required for certification.',
      trade: 'electrical',
      level: 'beginner',
      duration: 4,
      format: 'video',
      thumbnail: '/courses/safety.jpg',
      certification: true,
      certificationName: 'Electrical Safety Certified',
      enrolled: 892,
      rating: 4.6,
      modules: 4,
      lessons: 18,
      price: 99,
    },
  ]

  const filteredCourses = courses.filter((course) => {
    const matchesSearch = course.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         course.description.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesTrade = tradeFilter === 'all' || course.trade === tradeFilter
    const matchesLevel = levelFilter === 'all' || course.level === levelFilter
    return matchesSearch && matchesTrade && matchesLevel
  })

  const getLevelColor = (level: string) => {
    switch (level) {
      case 'beginner': return 'bg-green-100 text-green-700'
      case 'intermediate': return 'bg-blue-100 text-blue-700'
      case 'advanced': return 'bg-purple-100 text-purple-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  const getFormatIcon = (format: string) => {
    switch (format) {
      case 'video': return <Play className="h-4 w-4" />
      case 'ar': return <GraduationCap className="h-4 w-4" />
      case 'vr': return <GraduationCap className="h-4 w-4" />
      default: return <BookOpen className="h-4 w-4" />
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Training & Certifications</h1>
              <p className="text-sm text-gray-600">Master your trade with AI-powered courses</p>
            </div>
            <Link
              href="/training/my-courses"
              className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
            >
              <BookOpen className="h-5 w-5" />
              My Courses
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
                placeholder="Search courses..."
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
                value={levelFilter}
                onChange={(e) => setLevelFilter(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="all">All Levels</option>
                <option value="beginner">Beginner</option>
                <option value="intermediate">Intermediate</option>
                <option value="advanced">Advanced</option>
              </select>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Available Courses</div>
            <div className="text-2xl font-bold text-gray-900">{courses.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Students</div>
            <div className="text-2xl font-bold text-primary-600">
              {courses.reduce((sum, c) => sum + c.enrolled, 0).toLocaleString()}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Certifications</div>
            <div className="text-2xl font-bold text-green-600">
              {courses.filter(c => c.certification).length}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm">
            <div className="text-sm text-gray-600">Total Hours</div>
            <div className="text-2xl font-bold text-gray-900">
              {courses.reduce((sum, c) => sum + c.duration, 0)}+
            </div>
          </div>
        </div>

        {/* Courses Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCourses.map((course) => (
            <Link
              key={course.id}
              href={`/training/${course.id}`}
              className="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden"
            >
              {/* Thumbnail */}
              <div className="h-48 bg-gradient-to-br from-primary-500 to-blue-600 flex items-center justify-center">
                <GraduationCap className="h-20 w-20 text-white opacity-50" />
              </div>

              {/* Content */}
              <div className="p-6">
                <div className="flex items-center gap-2 mb-2">
                  <span className={`px-2 py-1 text-xs font-medium rounded-full capitalize ${getLevelColor(course.level)}`}>
                    {course.level}
                  </span>
                  <span className="text-xs text-gray-500 capitalize">{course.trade}</span>
                </div>

                <h3 className="text-lg font-semibold text-gray-900 mb-2">{course.title}</h3>
                <p className="text-sm text-gray-600 mb-4 line-clamp-2">{course.description}</p>

                <div className="flex items-center gap-4 text-sm text-gray-600 mb-4">
                  <div className="flex items-center gap-1">
                    <Clock className="h-4 w-4" />
                    {course.duration}h
                  </div>
                  <div className="flex items-center gap-1">
                    {getFormatIcon(course.format)}
                    {course.modules} modules
                  </div>
                  <div className="flex items-center gap-1">
                    <Award className="h-4 w-4" />
                    {course.rating}
                  </div>
                </div>

                {course.certification && (
                  <div className="flex items-center gap-2 mb-4 p-2 bg-green-50 rounded">
                    <Award className="h-4 w-4 text-green-600" />
                    <span className="text-xs text-green-700 font-medium">
                      Certification Available
                    </span>
                  </div>
                )}

                <div className="flex justify-between items-center pt-4 border-t">
                  <div>
                    <div className="text-2xl font-bold text-gray-900">${course.price}</div>
                    <div className="text-xs text-gray-500">{course.enrolled} enrolled</div>
                  </div>
                  <button className="px-4 py-2 bg-primary-600 text-white text-sm font-medium rounded-lg hover:bg-primary-700">
                    Enroll Now
                  </button>
                </div>
              </div>
            </Link>
          ))}
        </div>

        {filteredCourses.length === 0 && (
          <div className="bg-white p-12 rounded-lg shadow-sm text-center">
            <GraduationCap className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No courses found</h3>
            <p className="text-gray-600">Try adjusting your search or filters</p>
          </div>
        )}
      </div>
    </div>
  )
}
