# TechDocs AI Frontend

This is the Next.js frontend application for the TechDocs AI Chatbot.

## Setup

1. Install dependencies:

```bash
npm install
```

2. Run the development server:

```bash
npm run dev
```

3. Open your browser and navigate to `http://localhost:3000`.

## Features

- Responsive, centered layout inspired by Medium.com
- Top bar with "TechDocs AI" title
- Floating/sticky prompt input box at the bottom
- WebSocket connection to backend at `ws://localhost:8080/ws`
- AI tutorial-style Markdown/HTML response display
- Loading animations using Framer Motion
- Material UI components with clean typography

## Notes

- Ensure the backend server is running and accessible at the WebSocket URL.
- The frontend uses React 18, Next.js 13, Material UI v5, and Framer Motion.

## Docker

1. Build the Docker image:

```bash
docker build -t techdocs-ai-frontend .
```

2. Run the Docker container:

```bash
docker run -p 3000:3000 techdocs-ai-frontend
```

3. Access the app at `http://localhost:3000`.