import { useState, useEffect } from 'react'
import { Estimate } from '@/types'

export function useEstimates(autoFetch = true) {
  const [estimates, setEstimates] = useState<Estimate[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchEstimates = async () => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/estimates')

      if (!response.ok) {
        throw new Error('Failed to fetch estimates')
      }

      const data = await response.json()
      setEstimates(data.estimates || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  const createEstimate = async (estimateData: Partial<Estimate>) => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/estimates', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(estimateData),
      })

      if (!response.ok) {
        throw new Error('Failed to create estimate')
      }

      const data = await response.json()
      setEstimates(prev => [data, ...prev])
      return data
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const analyzeBlueprint = async (file: File) => {
    setLoading(true)
    setError(null)

    try {
      const formData = new FormData()
      formData.append('file', file)

      const response = await fetch('/api/blueprints/analyze', {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) {
        throw new Error('Failed to analyze blueprint')
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
      fetchEstimates()
    }
  }, [autoFetch])

  return {
    estimates,
    loading,
    error,
    fetchEstimates,
    createEstimate,
    analyzeBlueprint,
  }
}
