import { useState, useEffect } from 'react'
import { Job } from '@/types'

interface UseJobsOptions {
  trade?: string
  budget?: string
  urgency?: string
  autoFetch?: boolean
}

export function useJobs(options: UseJobsOptions = {}) {
  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const { trade, budget, urgency, autoFetch = true } = options

  const fetchJobs = async () => {
    setLoading(true)
    setError(null)

    try {
      const params = new URLSearchParams()
      if (trade) params.append('trade', trade)
      if (budget) params.append('budget', budget)
      if (urgency) params.append('urgency', urgency)

      const response = await fetch(`/api/marketplace/jobs?${params.toString()}`)

      if (!response.ok) {
        throw new Error('Failed to fetch jobs')
      }

      const data = await response.json()
      setJobs(data.jobs || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  const getJobById = async (id: string) => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch(`/api/marketplace/jobs/${id}`)

      if (!response.ok) {
        throw new Error('Failed to fetch job')
      }

      const data = await response.json()
      return data
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
      throw err
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (autoFetch) {
      fetchJobs()
    }
  }, [trade, budget, urgency, autoFetch])

  return {
    jobs,
    loading,
    error,
    fetchJobs,
    getJobById,
  }
}
