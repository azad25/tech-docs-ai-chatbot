const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost/ws');

ws.on('open', function open() {
  console.log('WebSocket connected');
  
  // Send a test message
  ws.send(JSON.stringify({
    type: 'chat',
    message: 'Hello from test client'
  }));
});

ws.on('message', function message(data) {
  console.log('Received:', data.toString());
});

ws.on('error', function error(err) {
  console.error('WebSocket error:', err);
});

ws.on('close', function close() {
  console.log('WebSocket disconnected');
});

// Keep the connection alive for 30 seconds
setTimeout(() => {
  ws.close();
}, 30000);