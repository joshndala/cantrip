#!/usr/bin/env python3
"""
Enhanced LangGraph Agent Flow for CanTrip
Integrates all available tools with proper system prompts and tool calling
"""

import asyncio
import json
import logging
import os
import re
from datetime import datetime, timedelta
from typing import Dict, Any, List, Optional, TypedDict
from langchain_core.messages import HumanMessage, AIMessage, SystemMessage
from langgraph.graph import StateGraph, END
from langchain_google_vertexai import ChatVertexAI

# Import model router
from model_router import model_router

# Import all tools
from tools.events import EventsTool
from tools.attractions import AttractionsTool
from tools.planner import PlanningTool
from tools.recommend import RecommendationTool
from tools.weather import WeatherTool

logger = logging.getLogger(__name__)

class EnhancedTravelState(TypedDict):
    """Enhanced state for travel planning with tool integration"""
    message: str
    session_id: str
    context: Dict[str, Any]
    history: List[Any]
    response: str
    intent: str
    confidence: float
    suggestions: List[str]
    data: Dict[str, Any]
    tools_used: List[str]
    city_detected: Optional[str]

def create_enhanced_travel_graph() -> StateGraph:
    """Create an enhanced travel planning graph with tool integration"""
    
    # Define the state
    workflow = StateGraph(EnhancedTravelState)
    
    # Add nodes
    workflow.add_node("analyze_intent", analyze_intent_node)
    workflow.add_node("call_tools", call_tools_node)
    workflow.add_node("generate_response", generate_response_node)
    workflow.add_node("finalize_response", finalize_response_node)
    
    # Define edges
    workflow.set_entry_point("analyze_intent")
    workflow.add_edge("analyze_intent", "call_tools")
    workflow.add_edge("call_tools", "generate_response")
    workflow.add_edge("generate_response", "finalize_response")
    workflow.add_edge("finalize_response", END)
    
    return workflow

async def analyze_intent_node(state: EnhancedTravelState) -> EnhancedTravelState:
    """Analyze user intent and determine which tools to use"""
    logger.info("Analyzing user intent...")
    
    message = state.get('message', '').lower()
    history = state.get('history', [])
    
    # Detect city from message
    city = detect_city_from_message(state.get('message', ''))
    state["city_detected"] = city
    
    # Analyze intent and determine tools needed
    tools_needed = []
    intent = "general_question"
    
    # Event-related queries
    if any(word in message for word in ['event', 'events', 'concert', 'show', 'game', 'sports', 'festival', 'happening', 'this weekend', 'tonight']):
        tools_needed.append("events")
        intent = "events_inquiry"
    
    # Weather-related queries
    if any(word in message for word in ['weather', 'temperature', 'climate', 'rain', 'snow', 'sunny', 'forecast', 'hot', 'cold']):
        tools_needed.append("weather")
        intent = "weather_inquiry"
    
    # Attraction-related queries
    if any(word in message for word in ['attraction', 'attractions', 'museum', 'gallery', 'landmark', 'sightseeing', 'visit', 'see', 'explore']):
        tools_needed.append("attractions")
        intent = "attractions_inquiry"
    
    # Planning-related queries
    if any(word in message for word in ['plan', 'itinerary', 'schedule', 'trip', 'visit', 'go to', 'travel', 'day', 'weekend']):
        tools_needed.append("planner")
        intent = "planning_inquiry"
    
    # Recommendation queries
    if any(word in message for word in ['recommend', 'suggest', 'best', 'popular', 'good', 'where to', 'what to do']):
        tools_needed.append("recommendations")
        intent = "recommendations_inquiry"
    
    # If no specific intent detected but city is mentioned, get general recommendations
    if city and not tools_needed:
        tools_needed = ["recommendations", "attractions"]
        intent = "general_city_inquiry"
    
    state["intent"] = intent
    state["tools_used"] = tools_needed
    
    logger.info(f"Detected intent: {intent}, tools needed: {tools_needed}, city: {city}")
    
    return state

async def call_tools_node(state: EnhancedTravelState) -> EnhancedTravelState:
    """Call the appropriate tools based on intent analysis"""
    logger.info("Calling tools...")
    
    tools_needed = state.get("tools_used", [])
    city = state.get("city_detected")
    message = state.get("message", "")
    tool_results = {}
    
    try:
        # Initialize tools
        events_tool = EventsTool()
        attractions_tool = AttractionsTool()
        planner_tool = PlanningTool()
        recommend_tool = RecommendationTool()
        weather_tool = WeatherTool()
        
        # Call events tool
        if "events" in tools_needed and city:
            logger.info(f"Calling events tool for {city}")
            try:
                # Extract event type and date from message
                event_type = extract_event_type(message)
                event_date = extract_date(message)
                events = await events_tool.get_events(city, date=event_date, category=event_type)
                tool_results["events"] = events
                logger.info(f"Found {len(events)} events")
            except Exception as e:
                logger.error(f"Error calling events tool: {e}")
                tool_results["events"] = []
        
        # Call weather tool
        if "weather" in tools_needed and city:
            logger.info(f"Calling weather tool for {city}")
            try:
                weather = await weather_tool.get_current_weather(city)
                forecast = await weather_tool.get_weather_forecast(city, 3)
                tool_results["weather"] = {
                    "current": weather,
                    "forecast": forecast
                }
                logger.info("Weather data retrieved")
            except Exception as e:
                logger.error(f"Error calling weather tool: {e}")
                tool_results["weather"] = {}
        
        # Call attractions tool
        if "attractions" in tools_needed and city:
            logger.info(f"Calling attractions tool for {city}")
            try:
                attractions = await attractions_tool.get_attractions(city)
                tool_results["attractions"] = attractions
                logger.info(f"Found {len(attractions)} attractions")
            except Exception as e:
                logger.error(f"Error calling attractions tool: {e}")
                tool_results["attractions"] = []
        
        # Call recommendations tool
        if "recommendations" in tools_needed and city:
            logger.info(f"Calling recommendations tool for {city}")
            try:
                recommendations = await recommend_tool.get_recommendations(city)
                tool_results["recommendations"] = recommendations
                logger.info(f"Found {len(recommendations)} recommendations")
            except Exception as e:
                logger.error(f"Error calling recommendations tool: {e}")
                tool_results["recommendations"] = []
        
        # Call planner tool
        if "planner" in tools_needed and city:
            logger.info(f"Calling planner tool for {city}")
            try:
                # Extract planning parameters from message
                duration = extract_duration(message)
                budget = extract_budget(message)
                interests = extract_interests(message)
                
                plan = await planner_tool.create_itinerary(
                    city=city,
                    duration=duration,
                    budget=budget,
                    interests=interests
                )
                tool_results["plan"] = plan
                logger.info("Itinerary created")
            except Exception as e:
                logger.error(f"Error calling planner tool: {e}")
                tool_results["plan"] = {}
        
        state["data"] = tool_results
        
    except Exception as e:
        logger.error(f"Error in tool calling: {e}")
        state["data"] = {}
    
    return state

async def generate_response_node(state: EnhancedTravelState) -> EnhancedTravelState:
    """Generate response using LLM with tool results"""
    logger.info("Generating response...")
    
    message = state.get('message', '')
    session_id = state.get('session_id', '')
    context = state.get('context', {})
    history = state.get('history', [])
    intent = state.get('intent', '')
    city = state.get('city_detected', '')
    tool_results = state.get('data', {})
    tools_used = state.get('tools_used', [])
    
    try:
        # Get appropriate LLM
        chat_llm = model_router.get_model_for_agent("explore", prompt_text=message)
        
        # Create comprehensive system prompt
        system_prompt = create_system_prompt(intent, city, tools_used, tool_results)
        
        # Create user prompt with context
        user_prompt = f"""
        User message: "{message}"
        
        Context: {context}
        History: {history}
        
        Please provide a helpful, detailed response based on the available data and tools.
        """
        
        # Get response from LLM
        messages = [
            SystemMessage(content=system_prompt),
            HumanMessage(content=user_prompt)
        ]
        
        response = await chat_llm.ainvoke(messages)
        
        # Generate suggestions based on intent and results
        suggestions = generate_suggestions(intent, city, tool_results)
        
        state["response"] = response.content
        state["suggestions"] = suggestions
        state["confidence"] = 0.9
        
    except Exception as e:
        logger.error(f"Error generating response: {e}")
        state["response"] = "I'm here to help with your Canadian travel planning! What would you like to know?"
        state["suggestions"] = ["Tell me about popular Canadian destinations", "Help me plan a trip"]
        state["confidence"] = 0.5
    
    return state

async def finalize_response_node(state: EnhancedTravelState) -> EnhancedTravelState:
    """Finalize the response"""
    logger.info("Finalizing response...")
    
    # Ensure we have a response
    if not state.get("response"):
        state["response"] = "I'm here to help with your travel planning!"
    
    return state

def create_system_prompt(intent: str, city: str, tools_used: List[str], tool_results: Dict[str, Any]) -> str:
    """Create a comprehensive system prompt based on intent and available data"""
    
    base_prompt = """You are CanTrip, an AI Canadian travel assistant. You help users plan trips, discover events, find attractions, and get weather information for Canadian destinations.

You have access to real-time data from multiple sources including Ticketmaster, Eventbrite, OpenWeather, and more. Always use this data to provide accurate, helpful responses.

Guidelines:
- Be friendly, helpful, and informative
- Use real data when available
- Provide specific recommendations with details
- Ask follow-up questions to better understand user needs
- Format responses clearly with bullet points and sections when appropriate
- Include practical information like prices, dates, and booking details when available
"""
    
    # Add intent-specific instructions
    if intent == "events_inquiry":
        events_data = tool_results.get('events', [])
        base_prompt += f"""
EVENTS INQUIRY: The user is asking about events in {city or 'a Canadian city'}.

IMPORTANT: You MUST use ONLY the real event data provided below. Do NOT make up or invent events. Do NOT use hardcoded dates from previous years.

Available event data: {json.dumps(events_data, indent=2)}

CRITICAL INSTRUCTIONS:
- If the event data is empty ([]), you MUST say "No events found for the specified criteria" or "No events are scheduled for this weekend"
- Do NOT generate fictional events like "Toronto Symphony Orchestra", "Blue Jays games", or "Luminato Festival"
- Do NOT use hardcoded dates from previous years (like "June 8th & 9th")
- Use ONLY the real events from the data above
- Present events with their actual dates, times, prices, and venues from the data
- Include booking links when available in the data
- If no events are found, suggest alternative activities like visiting attractions, parks, or markets
"""
    
    elif intent == "weather_inquiry":
        base_prompt += f"""
WEATHER INQUIRY: The user is asking about weather in {city or 'a Canadian city'}.

Available weather data: {json.dumps(tool_results.get('weather', {}), indent=2)}

Instructions:
- Provide current weather conditions with specific temperatures
- Include forecast information if available
- Give travel advice based on weather conditions
- Suggest appropriate clothing and activities
"""
    
    elif intent == "attractions_inquiry":
        base_prompt += f"""
ATTRACTIONS INQUIRY: The user is asking about attractions in {city or 'a Canadian city'}.

Available attraction data: {json.dumps(tool_results.get('attractions', []), indent=2)}

Instructions:
- List specific attractions with descriptions
- Include practical information (hours, prices, locations)
- Suggest the best attractions to visit
- Provide tips for visiting
"""
    
    elif intent == "planning_inquiry":
        base_prompt += f"""
PLANNING INQUIRY: The user is asking about trip planning for {city or 'a Canadian city'}.

Available planning data: {json.dumps(tool_results.get('plan', {}), indent=2)}

Instructions:
- Provide a structured itinerary or plan
- Include timing, locations, and activities
- Consider weather and events in recommendations
- Suggest budget-friendly options
"""
    
    elif intent == "recommendations_inquiry":
        base_prompt += f"""
RECOMMENDATIONS INQUIRY: The user is asking for recommendations for {city or 'a Canadian city'}.

Available recommendation data: {json.dumps(tool_results.get('recommendations', []), indent=2)}

Instructions:
- Provide specific, actionable recommendations
- Include a mix of popular and unique suggestions
- Consider the user's interests and preferences
- Provide practical details for each recommendation
"""
    
    return base_prompt

def generate_suggestions(intent: str, city: str, tool_results: Dict[str, Any]) -> List[str]:
    """Generate contextual suggestions based on intent and results"""
    
    suggestions = []
    
    if intent == "events_inquiry":
        suggestions = [
            f"What concerts are happening in {city}?",
            f"Are there any sports events in {city}?",
            f"What festivals are coming up in {city}?",
            f"Show me family-friendly events in {city}"
        ]
    elif intent == "weather_inquiry":
        suggestions = [
            f"What's the weather forecast for {city}?",
            f"What should I pack for {city}?",
            f"Is it good weather for outdoor activities in {city}?",
            f"What's the best time to visit {city}?"
        ]
    elif intent == "attractions_inquiry":
        suggestions = [
            f"What are the must-see attractions in {city}?",
            f"Show me free attractions in {city}",
            f"What museums are in {city}?",
            f"What outdoor attractions are in {city}?"
        ]
    elif intent == "planning_inquiry":
        suggestions = [
            f"Create a 3-day itinerary for {city}",
            f"Plan a budget trip to {city}",
            f"What should I do in {city} this weekend?",
            f"Plan a family trip to {city}"
        ]
    else:
        suggestions = [
            f"Tell me about {city}",
            f"What's the weather like in {city}?",
            f"What events are happening in {city}?",
            f"Plan a trip to {city}"
        ]
    
    return suggestions

def detect_city_from_message(message: str) -> Optional[str]:
    """Detect city from user message"""
    message_lower = message.lower()
    
    canadian_cities = [
        'toronto', 'vancouver', 'montreal', 'calgary', 'ottawa', 'edmonton',
        'winnipeg', 'quebec city', 'hamilton', 'kitchener', 'london', 'victoria',
        'halifax', 'oshawa', 'windsor', 'saskatoon', 'regina', 'sherbrooke',
        'barrie', 'kelowna', 'abbotsford', 'kingston', 'trois-rivières',
        'guelph', 'cambridge', 'whitby', 'ajax', 'milton', 'st. catharines',
        'brantford', 'thunder bay', 'saint john', 'peterborough', 'red deer',
        'lethbridge', 'kamloops', 'nanaimo', 'prince george', 'chilliwack',
        'vernon', 'fort mcmurray', 'sarnia', 'belleville', 'charlottetown',
        'fredericton', 'moncton', 'saint john', 'yellowknife', 'whitehorse',
        'iqaluit', 'banff', 'whistler', 'niagara falls', 'jasper', 'victoria'
    ]
    
    for city in canadian_cities:
        if city in message_lower:
            return city.title()
    
    return None

def extract_event_type(message: str) -> Optional[str]:
    """Extract event type from message"""
    message_lower = message.lower()
    
    # Skip common verbs that might be confused with event types
    skip_words = ['show me', 'show you', 'show us', 'tell me', 'give me', 'find me']
    if any(skip_word in message_lower for skip_word in skip_words):
        return None
    
    if any(word in message_lower for word in ['concert', 'music', 'band', 'singer']):
        return 'music'
    elif any(word in message_lower for word in ['sports', 'game', 'match', 'hockey', 'basketball', 'baseball']):
        return 'sports'
    elif any(word in message_lower for word in ['festival', 'fair', 'celebration']):
        return 'festival'
    elif any(word in message_lower for word in ['theater', 'theatre', 'play', 'musical', 'drama']):
        return 'theater'
    elif any(word in message_lower for word in ['comedy', 'standup', 'joke']):
        return 'comedy'
    
    return None

def extract_date(message: str) -> Optional[str]:
    """Extract date from message"""
    message_lower = message.lower()
    
    # Get current date
    today = datetime.now()
    
    if 'this weekend' in message_lower:
        # Find the next Saturday and Sunday
        days_until_saturday = (5 - today.weekday()) % 7
        if days_until_saturday == 0:  # Today is Saturday
            days_until_saturday = 7
        saturday = today + timedelta(days=days_until_saturday)
        return saturday.strftime("%Y-%m-%d")
    
    elif 'next weekend' in message_lower:
        # Find the Saturday and Sunday of next week
        days_until_saturday = (5 - today.weekday()) % 7 + 7
        saturday = today + timedelta(days=days_until_saturday)
        return saturday.strftime("%Y-%m-%d")
    
    elif 'tonight' in message_lower or 'this evening' in message_lower:
        return today.strftime("%Y-%m-%d")
    
    elif 'tomorrow' in message_lower:
        tomorrow = today + timedelta(days=1)
        return tomorrow.strftime("%Y-%m-%d")
    
    elif 'today' in message_lower:
        return today.strftime("%Y-%m-%d")
    
    # Try to extract specific dates (MM/DD, MM-DD, etc.)
    import re
    date_patterns = [
        r'(\d{1,2})/(\d{1,2})',  # MM/DD
        r'(\d{1,2})-(\d{1,2})',  # MM-DD
        r'(\d{4})-(\d{1,2})-(\d{1,2})',  # YYYY-MM-DD
    ]
    
    for pattern in date_patterns:
        match = re.search(pattern, message)
        if match:
            if len(match.groups()) == 2:  # MM/DD or MM-DD
                month, day = match.groups()
                # Assume current year
                try:
                    date_obj = datetime(today.year, int(month), int(day))
                    return date_obj.strftime("%Y-%m-%d")
                except ValueError:
                    continue
            elif len(match.groups()) == 3:  # YYYY-MM-DD
                year, month, day = match.groups()
                try:
                    date_obj = datetime(int(year), int(month), int(day))
                    return date_obj.strftime("%Y-%m-%d")
                except ValueError:
                    continue
    
    return None

def extract_duration(message: str) -> int:
    """Extract trip duration from message"""
    message_lower = message.lower()
    
    if 'weekend' in message_lower or '2 days' in message_lower:
        return 2
    elif 'week' in message_lower or '7 days' in message_lower:
        return 7
    elif 'month' in message_lower or '30 days' in message_lower:
        return 30
    elif any(word in message_lower for word in ['day', 'days']):
        # Try to extract number
        import re
        numbers = re.findall(r'\d+', message)
        if numbers:
            return int(numbers[0])
    
    return 3  # Default

def extract_budget(message: str) -> float:
    """Extract budget from message"""
    message_lower = message.lower()
    
    if 'budget' in message_lower or 'cheap' in message_lower:
        return 500.0
    elif 'luxury' in message_lower or 'expensive' in message_lower:
        return 3000.0
    elif 'mid' in message_lower or 'moderate' in message_lower:
        return 1500.0
    
    # Try to extract dollar amount
    import re
    amounts = re.findall(r'\$?(\d+)', message)
    if amounts:
        return float(amounts[0])
    
    return 1000.0  # Default

def extract_interests(message: str) -> List[str]:
    """Extract interests from message"""
    message_lower = message.lower()
    interests = []
    
    interest_keywords = {
        'music': ['music', 'concert', 'band', 'singer'],
        'sports': ['sports', 'game', 'hockey', 'basketball', 'baseball'],
        'arts': ['art', 'museum', 'gallery', 'theater', 'theatre'],
        'food': ['food', 'restaurant', 'dining', 'cuisine'],
        'outdoor': ['outdoor', 'hiking', 'park', 'nature', 'beach'],
        'culture': ['culture', 'cultural', 'history', 'heritage'],
        'family': ['family', 'kids', 'children'],
        'nightlife': ['nightlife', 'bar', 'club', 'drinks']
    }
    
    for interest, keywords in interest_keywords.items():
        if any(keyword in message_lower for keyword in keywords):
            interests.append(interest)
    
    return interests

class EnhancedTravelPlanningGraph:
    """Enhanced travel planning graph with tool integration"""
    
    def __init__(self):
        self.graph = create_enhanced_travel_graph()
        self.app = self.graph.compile()
    
    async def run(self, input_data: Dict[str, Any]) -> Dict[str, Any]:
        """Run the enhanced travel planning graph"""
        try:
            # Create initial state with defaults
            state = {
                "message": input_data.get("message", ""),
                "session_id": input_data.get("session_id", ""),
                "context": input_data.get("context", {}),
                "history": input_data.get("history", []),
                "response": "",
                "intent": "",
                "confidence": 0.0,
                "suggestions": [],
                "data": {},
                "tools_used": [],
                "city_detected": None
            }
            
            # Run the graph
            result = await self.app.ainvoke(state)
            
            return {
                "response": result.get("response", "I'm here to help with your travel planning!"),
                "session_id": result.get("session_id", ""),
                "intent": result.get("intent", "general"),
                "confidence": result.get("confidence", 0.8),
                "suggestions": result.get("suggestions", []),
                "data": result.get("data", {}),
                "tools_used": result.get("tools_used", []),
                "city_detected": result.get("city_detected"),
                "timestamp": datetime.now().isoformat()
            }
                
        except Exception as e:
            logger.error(f"Error running enhanced travel planning graph: {e}")
            return {
                "response": "I'm here to help with your travel planning! What would you like to know?",
                "session_id": input_data.get("session_id", ""),
                "intent": "general",
                "confidence": 0.5,
                "suggestions": ["Tell me about popular Canadian destinations"],
                "data": {},
                "tools_used": [],
                "city_detected": None,
                "timestamp": datetime.now().isoformat()
            }
