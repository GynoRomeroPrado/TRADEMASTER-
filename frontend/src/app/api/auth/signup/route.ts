import { NextRequest, NextResponse } from 'next/server'

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()
    const { email, password, firstName, lastName, company, trade } = body

    // Call backend auth service
    const response = await fetch(`${process.env.AUTH_SERVICE_URL}/api/auth/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        password,
        profile: {
          first_name: firstName,
          last_name: lastName,
          company,
          trade,
        },
      }),
    })

    if (!response.ok) {
      const error = await response.json()
      return NextResponse.json(
        { error: error.message || 'Signup failed' },
        { status: response.status }
      )
    }

    const data = await response.json()

    // Set httpOnly cookie with JWT token
    const res = NextResponse.json({
      user: data.user,
      message: 'Signup successful',
    })

    res.cookies.set('auth-token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 7, // 7 days
      path: '/',
    })

    return res
  } catch (error) {
    console.error('Signup error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
