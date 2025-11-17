import { NextRequest, NextResponse } from 'next/server'

function getAuthToken(request: NextRequest): string | undefined {
  return request.cookies.get('auth-token')?.value
}

export async function GET(request: NextRequest) {
  try {
    const token = getAuthToken(request)
    if (!token) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    // Get query parameters for filtering
    const searchParams = request.nextUrl.searchParams
    const trade = searchParams.get('trade')
    const budget = searchParams.get('budget')
    const urgency = searchParams.get('urgency')

    const queryString = new URLSearchParams()
    if (trade) queryString.append('trade', trade)
    if (budget) queryString.append('budget', budget)
    if (urgency) queryString.append('urgency', urgency)

    const response = await fetch(
      `${process.env.MARKETPLACE_SERVICE_URL}/api/jobs?${queryString.toString()}`,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }
    )

    if (!response.ok) {
      const error = await response.json()
      return NextResponse.json(
        { error: error.message || 'Failed to fetch jobs' },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Get jobs error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function POST(request: NextRequest) {
  try {
    const token = getAuthToken(request)
    if (!token) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    const body = await request.json()

    const response = await fetch(`${process.env.MARKETPLACE_SERVICE_URL}/api/jobs`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const error = await response.json()
      return NextResponse.json(
        { error: error.message || 'Failed to create job' },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data, { status: 201 })
  } catch (error) {
    console.error('Create job error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
