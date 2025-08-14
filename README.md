# CanTrip - AI-Powered Canadian Travel Planning Platform

CanTrip is a comprehensive travel planning platform that uses AI to generate personalized itineraries, packing lists, and travel recommendations for destinations across Canada. Built with Go (backend API) and Python (LangGraph agent), it provides intelligent travel planning based on mood, interests, budget, and Canadian destinations.

## 🚀 Features

### Core Functionality
- **Mood-based Exploration**: Get travel suggestions for Canadian destinations based on your mood and interests
- **AI Itinerary Generation**: Create detailed day-by-day itineraries for Canadian cities using LangGraph
- **Smart Packing Lists**: Generate personalized packing lists based on Canadian climate, destination, and activities
- **Cultural Tips & Advice**: Access destination-specific cultural information and travel tips for Canadian destinations
- **PDF Export**: Download itineraries and packing lists as beautiful PDFs
- **Conversational AI**: Natural language chat interface for Canadian travel planning

### Technical Features
- **LangGraph Agent**: Advanced AI agent for intelligent travel planning
- **Multi-API Integration**: Weather, events, attractions, and recommendations
- **Caching Layer**: Optimized performance with intelligent caching
- **Phoenix Evaluation**: Monitor and evaluate agent performance
- **RESTful API**: Clean, documented API endpoints

## 🏗️ Architecture

```
cantrip/
├── backend/                    # Go API server
│   ├── main.go                # Application entry point
│   ├── router/                # API route definitions
│   ├── handlers/              # HTTP request handlers
│   ├── services/              # Business logic and external API integration
│   ├── data/                  # Static data and configuration
│   ├── templates/             # HTML templates for PDF generation
│   ├── utils/                 # Utility functions
│   └── langgraph_agent/       # Python LangGraph agent
│       ├── main.py            # Agent entry point
│       ├── graph.py           # LangGraph flow definition
│       ├── tools/             # Agent tools and utilities
│       └── eval/              # Phoenix evaluation
└── README.md
```

## 🛠️ Technology Stack

### Backend (Go)
- **Framework**: Gin (HTTP web framework)
- **Language**: Go 1.25+
- **Dependencies**: See `go.mod`

### AI Agent (Python)
- **Framework**: LangGraph + LangChain
- **LLM**: OpenAI GPT-4
- **API**: FastAPI
- **Evaluation**: Phoenix
- **Dependencies**: See `langgraph_agent/requirements.txt`

### External APIs
- **Weather**: OpenWeatherMap
- **Places**: Geoapify, TripAdvisor, Google Places
- **Events**: Ticketmaster, Eventbrite
- **Storage**: Google Cloud Storage (PDFs)

## 📦 Installation

### Prerequisites
- Go 1.25+
- Python 3.9+
- Git

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd cantrip/backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export OPENWEATHER_API_KEY="your_openweather_api_key"
   export GEOAPIFY_API_KEY="your_geoapify_api_key"
   export GOOGLE_CLOUD_PROJECT="your_gcp_project"
   export GOOGLE_APPLICATION_CREDENTIALS="path/to/service-account.json"
   ```

4. **Run the server**
   ```bash
   go run main.go
   ```

### LangGraph Agent Setup

1. **Navigate to agent directory**
   ```bash
   cd langgraph_agent
   ```

2. **Create virtual environment**
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

3. **Install dependencies**
   ```bash
   pip install -r requirements.txt
   ```

4. **Set Python environment variables**
   ```bash
   export OPENAI_API_KEY="your_openai_api_key"
   export PHOENIX_ENABLED="true"  # Optional
   ```

5. **Run the agent**
   ```bash
   python main.py
   ```

## 🚀 Quick Start

### 1. Start the Backend Server
```bash
cd backend
go run main.go
```
Server will start on `http://localhost:8080`

### 2. Start the LangGraph Agent
```bash
cd backend/langgraph_agent
python main.py
```
Agent will start on `http://localhost:8001`

### 3. Test the API

**Health Check**
```bash
curl http://localhost:8080/api/v1/health
```

**Explore a Destination**
```bash
curl -X POST http://localhost:8080/api/v1/explore \
  -H "Content-Type: application/json" \
  -d '{
    "mood": "adventurous",
    "city": "Toronto",
    "budget": 1000,
    "duration": 7,
    "interests": ["outdoor", "culture"]
  }'
```

**Generate Itinerary**
```bash
curl -X POST http://localhost:8080/api/v1/itinerary \
  -H "Content-Type: application/json" \
  -d '{
    "city": "Toronto",
    "start_date": "2024-07-01",
    "end_date": "2024-07-07",
    "interests": ["culture", "food"],
    "budget": 1500,
    "group_size": 2,
    "pace": "moderate"
  }'
```

## 📚 API Documentation

### Core Endpoints

#### Explore
- `POST /api/v1/explore` - Get mood-based travel suggestions
- `GET /api/v1/explore/mood/:mood` - Get suggestions for specific mood

#### Itinerary
- `POST /api/v1/itinerary` - Generate new itinerary
- `GET /api/v1/itinerary/:id` - Get specific itinerary
- `PUT /api/v1/itinerary/:id` - Update itinerary
- `DELETE /api/v1/itinerary/:id` - Delete itinerary

#### Packing
- `POST /api/v1/packing` - Generate packing list
- `GET /api/v1/packing/:id` - Get packing list
- `PUT /api/v1/packing/:id` - Update packing list
- `GET /api/v1/packing/suggestions` - Get packing suggestions

#### Tips
- `POST /api/v1/tips` - Get travel tips
- `GET /api/v1/tips/cultural/:destination` - Cultural tips
- `GET /api/v1/tips/tipping/:destination` - Tipping guide
- `GET /api/v1/tips/safety/:destination` - Safety tips

#### PDF
- `POST /api/v1/pdf/generate` - Generate PDF
- `GET /api/v1/pdf/download/:id` - Download PDF
- `GET /api/v1/pdf/status/:id` - Check PDF status

### LangGraph Agent Endpoints

- `GET /health` - Health check
- `POST /generate-itinerary` - Generate itinerary via agent
- `POST /explore-destination` - Explore destination via agent
- `POST /generate-packing-list` - Generate packing list via agent

## 🔧 Configuration

### Environment Variables

#### Backend
```bash
# API Keys
OPENWEATHER_API_KEY=your_key
GEOAPIFY_API_KEY=your_key
TRIPADVISOR_API_KEY=your_key
TICKETMASTER_API_KEY=your_key
EVENTBRITE_API_KEY=your_key

# Google Cloud
GOOGLE_CLOUD_PROJECT=your_project
GOOGLE_APPLICATION_CREDENTIALS=path/to/credentials.json

# Server
PORT=8080
GIN_MODE=release
```

#### LangGraph Agent
```bash
# OpenAI
OPENAI_API_KEY=your_key

# Phoenix Evaluation
PHOENIX_ENABLED=true
PHOENIX_ENDPOINT=your_phoenix_endpoint

# Agent
AGENT_PORT=8001
AGENT_HOST=0.0.0.0
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Agent Tests
```bash
cd backend/langgraph_agent
pytest tests/
```

### Integration Tests
```bash
# Test full flow
curl -X POST http://localhost:8080/api/v1/explore \
  -H "Content-Type: application/json" \
  -d '{"mood": "relaxed", "city": "Vancouver", "budget": 800}'
```

## 📊 Monitoring & Evaluation

### Phoenix Evaluation
The LangGraph agent includes Phoenix integration for performance monitoring:

```bash
# Enable evaluation
export PHOENIX_ENABLED=true

# View evaluation summary
curl http://localhost:8001/evaluation/summary
```

### Logs
- Backend logs: Standard Go logging
- Agent logs: Structured logging with timestamps
- Evaluation logs: JSONL format in `phoenix_evaluations.jsonl`

## 🚀 Deployment

### Docker Deployment

1. **Build images**
   ```bash
   # Backend
   docker build -t cantrip-backend ./backend
   
   # Agent
   docker build -t cantrip-agent ./backend/langgraph_agent
   ```

2. **Run containers**
   ```bash
   docker run -p 8080:8080 cantrip-backend
   docker run -p 8001:8001 cantrip-agent
   ```

### Cloud Deployment

#### Google Cloud Run
```bash
# Deploy backend
gcloud run deploy cantrip-backend --source ./backend

# Deploy agent
gcloud run deploy cantrip-agent --source ./backend/langgraph_agent
```

#### AWS ECS
```bash
# Use provided docker-compose.yml
docker-compose up -d
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go and Python style guides
- Add tests for new features
- Update documentation
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [LangGraph](https://github.com/langchain-ai/langgraph) for the agent framework
- [Gin](https://github.com/gin-gonic/gin) for the Go web framework
- [Phoenix](https://github.com/Arize-ai/phoenix) for evaluation and monitoring
- All the external APIs that make this platform possible

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-repo/cantrip/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-repo/cantrip/discussions)
- **Email**: support@cantrip.com

---

**Happy Traveling! ✈️🌍**