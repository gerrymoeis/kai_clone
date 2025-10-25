/**
 * Cloudflare Pages Function for /counter/sync endpoint
 * 
 * This runs as part of your Pages deployment automatically!
 * No separate Worker deployment needed.
 * 
 * File location: functions/counter/sync.js
 * Handles: POST /counter/sync
 */

export async function onRequestPost(context) {
  try {
    // Parse form data
    const formData = await context.request.formData();
    const count = formData.get('count') || '0';
    
    // Parse as integer
    const num = parseInt(count, 10) || 0;
    
    // Return the count back (server echo demo)
    return new Response(num.toString(), {
      status: 200,
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
        'Access-Control-Allow-Origin': '*',
        'X-Powered-By': 'Cloudflare Pages Functions',
      },
    });
  } catch (error) {
    return new Response('bad request', {
      status: 400,
      headers: { 'Content-Type': 'text/plain; charset=utf-8' },
    });
  }
}

// Handle OPTIONS for CORS preflight (if needed)
export async function onRequestOptions(context) {
  return new Response(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type',
    },
  });
}
