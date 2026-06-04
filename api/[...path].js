const hopByHopHeaders = new Set(['connection', 'content-length', 'host'])

module.exports = async (req, res) => {
  const baseUrl = (process.env.API_BASE_URL || '').replace(/\/$/, '')
  if (!baseUrl) {
    res.statusCode = 502
    res.setHeader('Content-Type', 'application/json')
    res.end(JSON.stringify({ error: 'API_BASE_URL is not configured on Vercel' }))
    return
  }

  try {
    const response = await fetch(`${baseUrl}${req.url}`, {
      method: req.method,
      headers: requestHeaders(req.headers),
      body: hasBody(req.method) ? await readBody(req) : undefined
    })
    res.statusCode = response.status
    response.headers.forEach((value, key) => res.setHeader(key, value))
    res.end(Buffer.from(await response.arrayBuffer()))
  } catch (error) {
    res.statusCode = 502
    res.setHeader('Content-Type', 'application/json')
    res.end(JSON.stringify({ error: 'Failed to reach API service' }))
  }
}

function requestHeaders(headers) {
  return Object.fromEntries(Object.entries(headers).filter(([key]) => !hopByHopHeaders.has(key.toLowerCase())))
}

function hasBody(method) {
  return !['GET', 'HEAD'].includes(method)
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    const chunks = []
    req.on('data', chunk => chunks.push(chunk))
    req.on('end', () => resolve(Buffer.concat(chunks)))
    req.on('error', reject)
  })
}
