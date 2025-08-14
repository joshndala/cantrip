# CanTrip Conversational Chatbot

## Overview

CanTrip now supports conversational AI interactions through a chatbot interface. Users can ask natural language questions about travel across Canada and get intelligent responses with follow-up suggestions.

## How It Works

### 1. Conversational Flow
```
User Message → Go Handler → LangGraph Agent → Intent Analysis → Response + Suggestions
```

### 2. Session Management
- Each conversation has a unique session ID
- Context and history are preserved across messages
- No user authentication required

### 3. AI Agent Integration
- LangGraph agent handles natural language processing
- Intent recognition for different types of travel requests
- Context-aware responses with follow-up suggestions

## API Endpoints

### Chat Endpoints

#### `POST /api/v1/chat`
Main conversational endpoint.

**Request:**
```json
{
  "message": "I want to go somewhere warm in Canada",
  "session_id": "optional_session_id"
}
```

**Response:**
```json
{
  "response": "I'd recommend Vancouver, Victoria, or the Okanagan Valley for milder weather in Canada. Which interests you?",
  "session_id": "session_1234567890",
  "intent": "trip_planning",
  "confidence": 0.85,
  "suggestions": [
    "Tell me about Vancouver",
    "What's the weather like in Victoria?",
    "Help me plan a trip to the Okanagan"
  ],
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### `GET /api/v1/chat/history/:session_id`
Get conversation history for a session.

#### `DELETE /api/v1/chat/history/:session_id`
Clear conversation history.

#### `GET /api/v1/chat/suggestions/:session_id`
Get suggested follow-up questions.

## Supported Intents

The AI agent recognizes these types of requests:

- **destination_inquiry**: "Tell me about Paris"
- **trip_planning**: "I want to plan a trip"
- **weather_inquiry**: "What's the weather like in Tokyo?"
- **activity_suggestion**: "What activities are available?"
- **itinerary_request**: "Create an itinerary for me"
- **packing_request**: "Help me pack for my trip"
- **general_question**: General travel questions
- **greeting**: Simple greetings

## Example Conversations

### Planning a Trip
```
User: "I want to go somewhere warm in Canada"
Bot: "I'd recommend Vancouver, Victoria, or the Okanagan Valley for milder weather in Canada. Which interests you?"

User: "Tell me about Vancouver"
Bot: "Vancouver is a beautiful coastal city known for its stunning scenery, diverse culture, and outdoor activities. What specific information are you looking for - attractions, weather, culture, or something else?"

User: "Create a 7-day itinerary for Vancouver"
Bot: "I'd love to create an itinerary for your Canadian adventure! Please share your travel dates and interests."
```

### Weather Inquiry
```
User: "What's the weather like in Toronto?"
Bot: "I can help you check the weather for Toronto. Let me get that information for you..."
```

## Architecture

### Go Backend (API Layer)
- **`handlers/chat.go`**: Chat request handling
- **`services/conversation.go`**: Session management and AI integration
- **`router/routes.go`**: Chat endpoint routing

### Python LangGraph Agent (AI Brain)
- **`langgraph_agent/main.py`**: Chat endpoint
- **`langgraph_agent/graph.py`**: Chat processing logic
- Intent recognition and response generation

## Getting Started

1. **Start the Go backend:**
   ```bash
   cd backend
   go run main.go
   ```

2. **Start the LangGraph agent:**
   ```bash
   cd backend/langgraph_agent
   python main.py
   ```

3. **Test the chat:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/chat \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello, I need help planning a trip to Canada"}'
   ```

## Features

- ✅ Natural language processing
- ✅ Intent recognition
- ✅ Session management
- ✅ Context preservation
- ✅ Follow-up suggestions
- ✅ Integration with existing travel services
- ✅ No authentication required

## Future Enhancements

- Voice input/output
- Multi-language support
- Advanced context understanding
- Integration with booking systems
- Personalized recommendations based on history 