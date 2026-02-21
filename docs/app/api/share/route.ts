import { NextRequest, NextResponse } from 'next/server'

export async function POST(req: NextRequest) {
  const body = await req.text()

  const res = await fetch('https://go.dev/_/share', {
    method: 'POST',
    body,
  })

  const id = await res.text()
  return NextResponse.json({ url: `https://go.dev/play/p/${id}` })
}
