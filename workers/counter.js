/**
 * Cloudflare Worker for Gothic Forge Demo
 * Handles /counter/sync endpoint without backend server
 */

addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request));
});

async function handleRequest(request) {
  const url = new URL(request.url);
  
  // Handle counter sync endpoint
  if (url.pathname === '/counter/sync' && request.method === 'POST') {
    return handleCounterSync(request);
  }
  
  // Pass through everything else (static assets handled by Pages)
  return fetch(request);
}

async function handleCounterSync(request) {
  try {
    // Parse form data
    const formData = await request.formData();
    const count = formData.get('count') || '0';
    
    // Parse as integer
    const num = parseInt(count, 10) || 0;
    
    // Return the count back (server echo demo)
    return new Response(num.toString(), {
      status: 200,
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
        'Access-Control-Allow-Origin': '*',
      },
    });
  } catch (error) {
    return new Response('bad request', {
      status: 400,
      headers: { 'Content-Type': 'text/plain; charset=utf-8' },
    });
  }
}
