import { NextRequest, NextResponse } from 'next/server'

export async function POST(request: NextRequest) {
  const res = NextResponse.json({
    message: 'Logout successful',
  })

  // Clear the auth token cookie
  res.cookies.delete('auth-token')

  return res
}
