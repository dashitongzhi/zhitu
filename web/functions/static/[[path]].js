const ORIGIN = 'http://zhitu-origin.kralai.tech:8880'

export async function onRequest(context) {
  const incoming = new URL(context.request.url)
  const upstream = new URL(incoming.pathname + incoming.search, ORIGIN)
  const headers = new Headers(context.request.headers)
  headers.set('host', upstream.host)
  headers.set('x-forwarded-host', incoming.host)
  headers.set('x-forwarded-proto', 'https')
  return fetch(new Request(upstream.toString(), {
    method: context.request.method,
    headers,
    body: ['GET', 'HEAD'].includes(context.request.method)
      ? undefined
      : context.request.body,
    redirect: 'manual',
  }))
}
