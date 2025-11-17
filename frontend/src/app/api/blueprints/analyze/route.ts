import { NextRequest, NextResponse } from 'next/server'

function getAuthToken(request: NextRequest): string | undefined {
  return request.cookies.get('auth-token')?.value
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

    const formData = await request.formData()
    const file = formData.get('file') as File

    if (!file) {
      return NextResponse.json(
        { error: 'No file provided' },
        { status: 400 }
      )
    }

    // Create new FormData for the ML service
    const mlFormData = new FormData()
    mlFormData.append('file', file)

    // Call ML Computer Vision service
    const response = await fetch(`${process.env.ML_CV_SERVICE_URL}/api/blueprints/analyze`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: mlFormData,
    })

    if (!response.ok) {
      const error = await response.json()
      return NextResponse.json(
        { error: error.message || 'Failed to analyze blueprint' },
        { status: response.status }
      )
    }

    const data = await response.json()

    // Return detected components, quantities, and suggested materials
    return NextResponse.json({
      detected_components: data.components,
      area_sqft: data.area,
      suggested_materials: data.materials,
      confidence_score: data.confidence,
    })
  } catch (error) {
    console.error('Blueprint analysis error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
