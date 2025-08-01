# Tech Docs AI - Documentation and Tutorial Chat Service

Welcome to Tech Docs AI, a comprehensive documentation and tutorial chat service that provides AI-powered responses using RAG (Retrieval-Augmented Generation) with TinyLlama and vector embeddings. The system learns from user interactions, maintains conversation history, and provides well-formatted Markdown responses for display in Next.js frontends.

## 🚀 Features

- **Web Scraping**: Automated scraping of documentation and tutorials from various sources
- **Event-Driven Architecture**: Kafka-based job queue for scalable scraping operations
- **Worker Pools**: Concurrent processing of scraping jobs with configurable worker pools
- **Vector Database**: Qdrant vector store for semantic search and similarity matching
- **RAG System**: Retrieval-Augmented Generation using TinyLlama and Nomic embeddings
- **PostgreSQL Storage**: Reliable document storage with full-text search capabilities
- **Redis Caching**: High-performance caching for documents, embeddings, search results, and chat sessions
- **Conversation History**: Persistent chat sessions with context awareness
- **AI Learning**: Continuous learning from user interactions and responses
- **Intelligent Content Management**: Avoids unnecessary scraping when content already exists
- **RESTful API**: Complete API for chat, document management, and scraping operations
- **Containerized**: Full Docker Compose setup with all dependencies
- **Markdown Responses**: All AI responses are formatted in Markdown for frontend display
- **Smart Scraping**: Only scrapes when necessary, prioritizes existing knowledge base

## 🧠 AI Learning & Intelligence

### Continuous Learning System

The application implements a sophisticated learning system that continuously improves its knowledge base:

1. **Response Storage**: Every AI response is stored in the vector database for future reference
2. **User Interaction Learning**: The system learns from user questions and improves responses over time
3. **Context Awareness**: Maintains conversation history for better contextual responses
4. **Smart Content Management**: Avoids redundant scraping by checking existing content first

### Conversation History

- **Persistent Sessions**: Chat sessions are maintained across multiple interactions
- **Context Preservation**: Previous conversations are used to provide better context-aware responses
- **Session Analytics**: Insights into conversation patterns and topics
- **History Retrieval**: Access to complete conversation history for analysis

### Intelligent Content Management

- **Duplicate Prevention**: Checks for existing content before scraping
- **Relevance Scoring**: Uses vector similarity to determine content relevance
- **Adaptive Responses**: Learns from user interactions to improve future responses
- **Knowledge Base Growth**: Continuously expands the knowledge base with new information

### Markdown Formatting

All AI responses are structured as comprehensive Markdown tutorials with:
- **Proper Headers**: Clear section organization with H1, H2, H3 headers
- **Code Blocks**: Syntax-highlighted code examples
- **Bullet Points**: Well-organized lists and key points
- **Step-by-Step Guides**: Structured tutorials with numbered steps
- **Best Practices**: Highlighted tips and recommendations
- **Next Steps**: Suggested learning paths

## 🏗️ Architecture

The application follows a microservices architecture with the following components:

- **API Server**: HTTP server handling chat requests and document management
- **Worker Service**: Kafka consumer processing scraping jobs with worker pools
- **PostgreSQL**: Document storage with JSONB support for metadata
- **Qdrant**: Vector database for semantic search
- **Kafka**: Message queue for scraping job distribution
- **Redis**: High-performance caching layer for improved response times
- **Ollama**: Local LLM serving TinyLlama and Nomic embeddings
- **NGINX**: Reverse proxy for load balancing

## 🎯 Performance Optimizations

### Redis Caching Strategy

The application implements a comprehensive Redis caching strategy for maximum efficiency:

1. **Document Caching**: Frequently accessed documents are cached to reduce database queries
2. **Embedding Caching**: Vector embeddings are cached to avoid expensive re-computation
3. **Search Result Caching**: Search queries and their results are cached for faster responses
4. **Chat Session Caching**: Chat history and sessions are cached for seamless conversations
5. **Category-based Invalidation**: Smart cache invalidation when documents are updated

### Cache Benefits

- **Reduced Latency**: Cached responses are served in milliseconds
- **Lower Database Load**: Fewer queries to PostgreSQL for frequently accessed data
- **Cost Optimization**: Reduced embedding API calls through intelligent caching
- **Scalability**: Redis handles high concurrent access efficiently
- **User Experience**: Faster chat responses and search results

## 🛠️ Prerequisites

Before you begin, ensure you have the following installed:

- **Docker**: Containerization platform
- **Docker Compose**: Multi-container orchestration
- **Go 1.22+**: For building the application

## 📦 Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd tech-docs-ai
```

### 2. Start the Application

The Makefile simplifies the entire process:

```bash
make up
```

This command:
- Pulls the required Ollama models (tinyllama, nomic-embed-text)
- Builds the Go application Docker images
- Starts all services (API server, worker, NGINX, PostgreSQL, Qdrant, Kafka, Redis, Ollama)

### 3. Access the Application

Once all services are running, your application will be available at:

- **Main Application**: http://localhost
- **API Documentation**: http://localhost/api/v1/health

## 🗣️ API Usage

### Chat with AI

Ask questions about programming and get AI-powered responses in Markdown format:

```bash
curl -X POST http://localhost/api/v1/chat \
  -H 'Content-Type: application/json' \
  -d '{
    "message": "What is HTML and how do I create a basic webpage?"
  }'
```

### Chat with History

Maintain conversation context across multiple interactions:

```bash
curl -X POST http://localhost/api/v1/chat/history \
  -H 'Content-Type: application/json' \
  -d '{
    "session_id": "user123",
    "message": "Can you show me how to add CSS styling to that HTML?"
  }'
```

### Scrape Documentation

Queue a scraping job for a specific URL:

```bash
curl -X POST http://localhost/api/v1/scrape \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://www.w3schools.com/html/html_intro.asp",
    "category": "HTML",
    "tags": ["tutorial", "basics", "introduction"]
  }'
```

### Search Documents

Search through scraped documentation:

```bash
curl 'http://localhost/api/v1/documents/search?q=CSS&limit=5'
```

### Add Custom Documents

Add your own documentation:

```bash
curl -X POST http://localhost/api/v1/documents \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Custom Tutorial",
    "content": "Your tutorial content here...",
    "category": "Custom",
    "tags": ["custom", "tutorial"],
    "author": "Your Name"
  }'
```

### Generate Tutorials

Generate tutorials from existing content:

```bash
curl -X POST http://localhost/api/v1/tutorials/generate \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://example.com",
    "topic": "JavaScript Promises"
  }'
```

### Get Chat History

Retrieve conversation history:

```bash
curl 'http://localhost/api/v1/chat/history?session_id=user123&limit=10'
```

### Get Conversation Insights

Analyze conversation patterns:

```bash
curl 'http://localhost/api/v1/chat/insights?session_id=user123'
```

## 🧑‍💻 Code Structure

The project follows a clean, layered architecture:

```
.
├── cmd/
│   ├── server/
│   │   └── main.go           # API server entry point
│   └── worker/
│       └── main.go           # Worker service entry point
├── internal/
│   ├── app/
│   │   ├── handler.go        # HTTP request handlers
│   │   ├── service.go        # Business logic and RAG implementation
│   │   └── seeder.go         # Database initialization
│   ├── cache/
│   │   └── redis.go          # Redis caching implementation
│   ├── emb/
│   │   ├── ollama.go         # Ollama client for embeddings and chat
│   │   └── fake.go           # Mock client for testing
│   ├── kafka/
│   │   ├── producer.go       # Kafka message producer
│   │   └── consumer.go       # Kafka consumer with worker pools
│   ├── repo/
│   │   └── postgres.go       # PostgreSQL document storage
│   ├── scraper/
│   │   └── w3schools.go      # Web scraper for documentation
│   ├── types/
│   │   └── types.go          # Shared data types
│   └── vec/
│       └── qdrant.go         # Qdrant vector database client
├── Dockerfile                # Multi-stage Docker build
├── docker-compose.yml        # Service orchestration
├── Makefile                  # Build and deployment commands
└── nginx.conf               # NGINX reverse proxy configuration
```

## 🔧 Configuration

### Environment Variables

The application uses the following environment variables:

```bash
# Ollama Configuration
OLLAMA_API_URL=http://ollama:11434
OLLAMA_MODEL=nomic-embed-text
OLLAMA_CHAT_MODEL=tinyllama

# Database Configuration
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=tech-docs-ai-chat

# Vector Database Configuration
QDRANT_API_URL=http://qdrant:6333
QDRANT_COLLECTION=tech_docs_knowledge

# Kafka Configuration
KAFKA_URL=kafka:9092

# Redis Configuration
REDIS_URL=redis://redis:6379
```

## 🚀 Deployment

### Development

```bash
# Build and start all services
make up

# View logs
docker-compose logs -f

# Stop services
make down
```

### Production

For production deployment, consider:

1. **Security**: Use proper authentication and authorization
2. **Scaling**: Deploy multiple worker instances
3. **Monitoring**: Add Prometheus/Grafana for metrics
4. **Backup**: Set up automated database backups
5. **SSL**: Configure HTTPS with proper certificates
6. **Redis Persistence**: Configure Redis persistence for data durability
7. **Cache Warming**: Implement cache warming strategies for critical data

## 📊 Performance Monitoring

Monitor the application performance using:

- **Redis Metrics**: Track cache hit/miss ratios
- **Response Times**: Monitor API response times
- **Database Queries**: Track query performance
- **Memory Usage**: Monitor Redis and application memory usage

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Ollama](https://ollama.ai/) for local LLM inference
- [Qdrant](https://qdrant.tech/) for vector database
- [Redis](https://redis.io/) for high-performance caching
- [TinyLlama](https://github.com/jzhang38/TinyLlama) for the efficient language model