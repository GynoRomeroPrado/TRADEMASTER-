'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Play, CheckCircle, Lock, Clock, Award, BookOpen, Brain, Download } from 'lucide-react'

export default function CourseDetailPage({ params }: { params: { id: string } }) {
  const [activeTab, setActiveTab] = useState('overview')

  // Mock data
  const course = {
    id: params.id,
    title: 'Data Center Electrical Systems',
    description: 'Master the electrical systems used in modern data centers, including power distribution, backup systems, and cooling infrastructure.',
    trade: 'electrical',
    level: 'advanced',
    duration: 12,
    format: 'mixed',
    certification: {
      enabled: true,
      name: 'Data Center Electrician Certified',
      issuer: 'TRADEMASTER Academy',
    },
    price: 299,
    enrolled: 245,
    rating: 4.8,
    instructor: {
      name: 'Robert Johnson',
      title: 'Master Electrician',
      experience: '25 years',
    },
    modules: [
      {
        id: '1',
        title: 'Introduction to Data Centers',
        order: 1,
        lessons: [
          { id: '1-1', title: 'What is a Data Center?', type: 'video', duration: 15, completed: true },
          { id: '1-2', title: 'Data Center Power Requirements', type: 'video', duration: 20, completed: true },
          { id: '1-3', title: 'Industry Standards & Certifications', type: 'video', duration: 18, completed: false },
          { id: '1-4', title: 'Module 1 Quiz', type: 'quiz', duration: 10, completed: false },
        ],
      },
      {
        id: '2',
        title: 'Power Distribution Systems',
        order: 2,
        lessons: [
          { id: '2-1', title: 'Main Power Distribution', type: 'video', duration: 25, completed: false },
          { id: '2-2', title: 'PDU Installation', type: 'ar', duration: 30, completed: false },
          { id: '2-3', title: 'Cable Management', type: 'video', duration: 20, completed: false },
          { id: '2-4', title: 'Hands-on Exercise', type: 'assessment', duration: 45, completed: false },
        ],
      },
      {
        id: '3',
        title: 'Backup Power Systems',
        order: 3,
        lessons: [
          { id: '3-1', title: 'UPS Systems Overview', type: 'video', duration: 22, completed: false },
          { id: '3-2', title: 'Generator Installation', type: 'video', duration: 28, completed: false },
          { id: '3-3', title: 'Automatic Transfer Switches', type: 'video', duration: 18, completed: false },
        ],
      },
    ],
    prerequisites: ['Basic Electrical Fundamentals', 'NEC Code Basics'],
    whatYouLearn: [
      'Design and install data center power distribution systems',
      'Implement backup power solutions (UPS, generators)',
      'Follow industry standards and best practices',
      'Troubleshoot complex electrical issues',
      'Ensure high availability and redundancy',
    ],
  }

  const totalLessons = course.modules.reduce((sum, m) => sum + m.lessons.length, 0)
  const completedLessons = course.modules.reduce((sum, m) =>
    sum + m.lessons.filter(l => l.completed).length, 0
  )
  const progress = Math.round((completedLessons / totalLessons) * 100)

  const getLessonIcon = (type: string) => {
    switch (type) {
      case 'video': return <Play className="h-4 w-4" />
      case 'ar': return <Brain className="h-4 w-4" />
      case 'quiz': return <BookOpen className="h-4 w-4" />
      case 'assessment': return <Award className="h-4 w-4" />
      default: return <BookOpen className="h-4 w-4" />
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-gradient-to-r from-primary-600 to-blue-600 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Link href="/training" className="inline-flex items-center gap-2 text-white/80 hover:text-white mb-4">
            <ArrowLeft className="h-4 w-4" />
            Back to Courses
          </Link>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2">
              <div className="flex items-center gap-2 mb-3">
                <span className="px-3 py-1 bg-white/20 text-white text-sm font-medium rounded-full capitalize">
                  {course.level}
                </span>
                <span className="px-3 py-1 bg-white/20 text-white text-sm font-medium rounded-full capitalize">
                  {course.trade}
                </span>
              </div>

              <h1 className="text-3xl font-bold mb-4">{course.title}</h1>
              <p className="text-white/90 text-lg mb-6">{course.description}</p>

              <div className="flex items-center gap-6 text-sm">
                <div className="flex items-center gap-2">
                  <Clock className="h-5 w-5" />
                  {course.duration} hours
                </div>
                <div className="flex items-center gap-2">
                  <BookOpen className="h-5 w-5" />
                  {totalLessons} lessons
                </div>
                <div className="flex items-center gap-2">
                  <Award className="h-5 w-5" />
                  {course.rating} rating
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg p-6 text-gray-900">
              <div className="text-3xl font-bold mb-2">${course.price}</div>
              <button className="w-full px-6 py-3 bg-primary-600 text-white font-semibold rounded-lg hover:bg-primary-700 mb-4">
                Enroll Now
              </button>

              {course.certification.enabled && (
                <div className="p-3 bg-green-50 rounded-lg mb-4">
                  <div className="flex items-center gap-2 text-green-700 font-medium mb-1">
                    <Award className="h-5 w-5" />
                    Certification Included
                  </div>
                  <div className="text-sm text-green-600">{course.certification.name}</div>
                </div>
              )}

              <div className="text-sm text-gray-600 space-y-2">
                <div>✓ Lifetime access</div>
                <div>✓ AI-powered coaching</div>
                <div>✓ Mobile app access</div>
                <div>✓ Certificate of completion</div>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Tabs */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-8">
            {['overview', 'curriculum', 'instructor'].map((tab) => (
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
            <div className="lg:col-span-2 space-y-6">
              {/* What You'll Learn */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">What You'll Learn</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {course.whatYouLearn.map((item, index) => (
                    <div key={index} className="flex items-start gap-3">
                      <CheckCircle className="h-5 w-5 text-green-500 mt-0.5 flex-shrink-0" />
                      <span className="text-gray-700">{item}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Prerequisites */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Prerequisites</h2>
                <div className="space-y-2">
                  {course.prerequisites.map((prereq, index) => (
                    <div key={index} className="flex items-center gap-2 text-gray-700">
                      <div className="w-2 h-2 bg-primary-600 rounded-full"></div>
                      {prereq}
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              {/* Progress */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h3 className="font-semibold text-gray-900 mb-4">Your Progress</h3>
                <div className="mb-2">
                  <div className="flex justify-between text-sm mb-1">
                    <span className="text-gray-600">Completion</span>
                    <span className="font-medium text-gray-900">{progress}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-primary-600 h-2 rounded-full transition-all"
                      style={{ width: `${progress}%` }}
                    />
                  </div>
                </div>
                <div className="text-sm text-gray-600 mt-4">
                  {completedLessons} of {totalLessons} lessons completed
                </div>
              </div>

              {/* Stats */}
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <h3 className="font-semibold text-gray-900 mb-4">Course Stats</h3>
                <div className="space-y-3 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Students Enrolled</span>
                    <span className="font-medium text-gray-900">{course.enrolled}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Average Rating</span>
                    <span className="font-medium text-gray-900">{course.rating}/5.0</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Total Duration</span>
                    <span className="font-medium text-gray-900">{course.duration}h</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'curriculum' && (
          <div className="bg-white rounded-lg shadow-sm">
            <div className="p-6 border-b">
              <h2 className="text-xl font-semibold text-gray-900">Course Curriculum</h2>
              <p className="text-sm text-gray-600 mt-1">
                {course.modules.length} modules • {totalLessons} lessons • {course.duration}h total
              </p>
            </div>

            <div className="divide-y">
              {course.modules.map((module) => (
                <div key={module.id} className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-semibold text-gray-900">
                      Module {module.order}: {module.title}
                    </h3>
                    <span className="text-sm text-gray-600">
                      {module.lessons.length} lessons
                    </span>
                  </div>

                  <div className="space-y-2">
                    {module.lessons.map((lesson) => (
                      <div
                        key={lesson.id}
                        className="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 cursor-pointer"
                      >
                        <div className="flex items-center gap-3">
                          {lesson.completed ? (
                            <CheckCircle className="h-5 w-5 text-green-500" />
                          ) : (
                            <Lock className="h-5 w-5 text-gray-400" />
                          )}
                          <div className="flex items-center gap-2">
                            {getLessonIcon(lesson.type)}
                            <span className="font-medium text-gray-900">{lesson.title}</span>
                          </div>
                        </div>
                        <div className="flex items-center gap-2 text-sm text-gray-600">
                          <Clock className="h-4 w-4" />
                          {lesson.duration} min
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'instructor' && (
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-xl font-semibold text-gray-900 mb-6">Your Instructor</h2>
            <div className="flex items-start gap-6">
              <div className="w-24 h-24 bg-gray-200 rounded-full flex items-center justify-center">
                <span className="text-2xl font-bold text-gray-600">
                  {course.instructor.name.split(' ').map(n => n[0]).join('')}
                </span>
              </div>
              <div>
                <h3 className="text-xl font-semibold text-gray-900">{course.instructor.name}</h3>
                <p className="text-gray-600 mb-2">{course.instructor.title}</p>
                <p className="text-sm text-gray-600">{course.instructor.experience} of industry experience</p>
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  )
}
