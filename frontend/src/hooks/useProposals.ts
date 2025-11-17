import { useState, useEffect } from 'react'
import { Proposal } from '@/types'

export function useProposals(autoFetch = true) {
  const [proposals, setProposals] = useState<Proposal[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchProposals = async () => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/marketplace/proposals')

      if (!response.ok) {
        throw new Error('Failed to fetch proposals')
      }

      const data = await response.json()
      setProposals(data.proposals || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  const submitProposal = async (proposalData: {
    jobId: string
    proposedBudget: number
    timeline: string
    coverLetter: string
  }) => {
    setLoading(true)
    setError(null)

    try {
      const response = await fetch('/api/marketplace/proposals', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(proposalData),
      })

      if (!response.ok) {
        throw new Error('Failed to submit proposal')
      }

      const data = await response.json()
      setProposals(prev => [data, ...prev])
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
      fetchProposals()
    }
  }, [autoFetch])

  return {
    proposals,
    loading,
    error,
    fetchProposals,
    submitProposal,
  }
}
