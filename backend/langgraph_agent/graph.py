#!/usr/bin/env python3
"""
LangGraph Agent Flow for CanTrip
Defines the graph structure and nodes for travel planning
"""

import asyncio
import json
import logging
from datetime import datetime, timedelta
from typing import Dict, Any, List, Optional
from langchain_core.messages import HumanMessage, AIMessage
from langgraph.graph import StateGraph, END
from langgraph.prebuilt import ToolNode
from langchain_openai import ChatOpenAI
from langchain_core.tools import tool

# Import tools
from tools.recommend import RecommendationTool
from tools.events import EventsTool
from tools.attractions import AttractionsTool
from tools.planner import PlanningTool

logger = logging.getLogger(__name__)

# Initialize tools
recommendation_tool = RecommendationTool()
events_tool = EventsTool()
attractions_tool = AttractionsTool()
planning_tool = PlanningTool()

# Initialize LLM
llm = ChatOpenAI(
    model="gpt-4-turbo-preview",
    temperature=0.7,
    api_key=os.getenv("OPENAI_API_KEY")
)

class TravelState:
    """State class for travel planning"""
    def __init__(self):
        self.messages: List = []
        self.city: str = ""
        self.start_date: str = ""
        self.end_date: str = ""
        self.interests: List[str] = []
        self.budget: float = 1000.0
        self.group_size: int = 1
        self.pace: str = "moderate"
        self.accommodation: str = "mid-range"
        self.task: str = ""
        self.weather: Dict = {}
        self.events: List = []
        self.attractions: List = []
        self.suggestions: List = []
        self.itinerary: Dict = {}
        self.packing_list: Dict = {}
        self.mood: str = ""
        self.duration: int = 7
        self.destination: str = ""
        self.activities: List[str] = []
        self.age_group: str = "adult"
        self.special_needs: List[str] = []
        # Chat-specific fields
        self.message: str = ""
        self.session_id: str = ""
        self.context: Dict = {}
        self.history: List = []
        self.response: str = ""
        self.intent: str = ""
        self.confidence: float = 0.0
        self.suggestions: List[str] = []

def create_travel_graph() -> StateGraph:
    """Create the travel planning graph"""
    
    # Define the state
    workflow = StateGraph(TravelState)
    
    # Add nodes
    workflow.add_node("analyze_request", analyze_request_node)
    workflow.add_node("process_chat", process_chat_node)
    workflow.add_node("get_weather", get_weather_node)
    workflow.add_node("get_events", get_events_node)
    workflow.add_node("get_attractions", get_attractions_node)
    workflow.add_node("generate_suggestions", generate_suggestions_node)
    workflow.add_node("plan_itinerary", plan_itinerary_node)
    workflow.add_node("generate_packing_list", generate_packing_list_node)
    workflow.add_node("finalize_response", finalize_response_node)
    
    # Define edges
    workflow.set_entry_point("analyze_request")
    
    # Conditional routing based on task
    workflow.add_conditional_edges(
        "analyze_request",
        route_by_task,
        {
            "chat": "process_chat",
            "explore": "get_weather",
            "itinerary": "get_weather",
            "packing": "generate_packing_list"
        }
    )
    
    # Chat flow
    workflow.add_edge("process_chat", "finalize_response")
    
    # Explore flow
    workflow.add_edge("get_weather", "get_events")
    workflow.add_edge("get_events", "get_attractions")
    workflow.add_edge("get_attractions", "generate_suggestions")
    workflow.add_edge("generate_suggestions", "finalize_response")
    
    # Itinerary flow
    workflow.add_edge("get_weather", "plan_itinerary")
    workflow.add_edge("plan_itinerary", "finalize_response")
    
    # Packing flow
    workflow.add_edge("generate_packing_list", "finalize_response")
    
    # End
    workflow.add_edge("finalize_response", END)
    
    return workflow

async def analyze_request_node(state: TravelState) -> TravelState:
    """Analyze the incoming request and set up the state"""
    logger.info("Analyzing request...")
    
    # Extract information from the request
    if state.task == "explore_destination":
        state.city = state.city
        state.mood = state.mood
        state.duration = state.duration
    elif state.task == "generate_itinerary":
        state.city = state.city
        state.duration = calculate_duration(state.start_date, state.end_date)
    elif state.task == "generate_packing_list":
        state.destination = state.destination
        state.duration = calculate_duration(state.start_date, state.end_date)
    
    # Add analysis message
    state.messages.append(
        AIMessage(content=f"Analyzed request for task: {state.task}")
    )
    
    return state

async def process_chat_node(state: TravelState) -> TravelState:
    """Process conversational chat messages"""
    logger.info("Processing chat message...")
    
    # Extract message and context
    message = getattr(state, 'message', '')
    session_id = getattr(state, 'session_id', '')
    context = getattr(state, 'context', {})
    history = getattr(state, 'history', [])
    
    # Analyze intent using LLM
    intent_prompt = f"""
    Analyze this travel-related message and determine the intent:
    Message: {message}
    Context: {context}
    
    Possible intents:
    - destination_inquiry: User wants to know about a specific destination
    - trip_planning: User wants to plan a trip
    - weather_inquiry: User wants weather information
    - activity_suggestion: User wants activity recommendations
    - itinerary_request: User wants an itinerary
    - packing_request: User wants a packing list
    - general_question: General travel question
    - greeting: Simple greeting
    
    Return only the intent category.
    """
    
    try:
        intent_response = llm.invoke([HumanMessage(content=intent_prompt)])
        intent = intent_response.content.strip().lower()
    except Exception as e:
        logger.error(f"Error analyzing intent: {e}")
        intent = "general_question"
    
    # Generate response based on intent
    if intent == "destination_inquiry":
        response = f"I'd be happy to tell you about {message.split()[-1]}! What specific information are you looking for - attractions, weather, culture, or something else?"
        suggestions = ["Tell me about attractions", "What's the weather like?", "Give me cultural tips"]
    elif intent == "trip_planning":
        response = "Great! I can help you plan your trip across Canada. To get started, could you tell me your destination, travel dates, and what interests you?"
        suggestions = ["I want to go to Toronto", "Help me plan a trip to Vancouver", "I need a budget-friendly trip to Montreal"]
    elif intent == "weather_inquiry":
        response = "I can help you check the weather for your Canadian destination. Which city or province are you interested in?"
        suggestions = ["What's the weather in Toronto?", "How's the weather in Vancouver?", "Weather forecast for Montreal"]
    elif intent == "itinerary_request":
        response = "I'd love to create an itinerary for your Canadian adventure! Please share your destination, dates, and interests."
        suggestions = ["Create a 7-day Toronto itinerary", "Plan a weekend in Vancouver", "Make an itinerary for Montreal"]
    elif intent == "packing_request":
        response = "I can help you create a packing list for your Canadian trip! Tell me your destination, travel dates, and planned activities."
        suggestions = ["Packing list for Toronto", "What to pack for Vancouver", "Winter travel packing list for Canada"]
    else:
        response = "I'm your AI Canadian travel assistant! I can help you plan trips across Canada, suggest destinations, create itineraries, check weather, and more. What would you like to know?"
        suggestions = ["Tell me about popular Canadian destinations", "Help me plan a trip to Toronto", "What's the weather like in Vancouver?"]
    
    # Update state
    state.intent = intent
    state.response = response
    state.suggestions = suggestions
    state.confidence = 0.8
    
    return state

async def get_weather_node(state: TravelState) -> TravelState:
    """Get weather information for the destination"""
    logger.info(f"Getting weather for {state.city}")
    
    try:
        # This would integrate with a weather API
        # For now, we'll simulate weather data
        weather_data = {
            "temperature": 22,
            "condition": "Sunny",
            "humidity": 65,
            "wind_speed": 10
        }
        
        state.weather = weather_data
        state.messages.append(
            AIMessage(content=f"Weather data retrieved for {state.city}")
        )
        
    except Exception as e:
        logger.error(f"Error getting weather: {e}")
        state.weather = {}
    
    return state

async def get_events_node(state: TravelState) -> TravelState:
    """Get events for the destination"""
    logger.info(f"Getting events for {state.city}")
    
    try:
        events = await events_tool.get_events(state.city)
        state.events = events
        state.messages.append(
            AIMessage(content=f"Retrieved {len(events)} events for {state.city}")
        )
        
    except Exception as e:
        logger.error(f"Error getting events: {e}")
        state.events = []
    
    return state

async def get_attractions_node(state: TravelState) -> TravelState:
    """Get attractions for the destination"""
    logger.info(f"Getting attractions for {state.city}")
    
    try:
        attractions = await attractions_tool.get_attractions(state.city)
        state.attractions = attractions
        state.messages.append(
            AIMessage(content=f"Retrieved {len(attractions)} attractions for {state.city}")
        )
        
    except Exception as e:
        logger.error(f"Error getting attractions: {e}")
        state.attractions = []
    
    return state

async def generate_suggestions_node(state: TravelState) -> TravelState:
    """Generate travel suggestions based on mood and interests"""
    logger.info(f"Generating suggestions for {state.city} with mood: {state.mood}")
    
    try:
        # Use LLM to generate suggestions
        prompt = f"""
        Generate travel suggestions for {state.city} based on:
        - Mood: {state.mood}
        - Interests: {', '.join(state.interests)}
        - Budget: ${state.budget}
        - Duration: {state.duration} days
        - Weather: {state.weather}
        - Available events: {len(state.events)} events
        - Available attractions: {len(state.attractions)} attractions
        
        Provide 5-7 detailed suggestions with activities, estimated costs, and timing.
        """
        
        response = await llm.ainvoke([HumanMessage(content=prompt)])
        
        # Parse suggestions (this would be more sophisticated in practice)
        suggestions = parse_suggestions(response.content)
        state.suggestions = suggestions
        
        state.messages.append(
            AIMessage(content=f"Generated {len(suggestions)} suggestions for {state.city}")
        )
        
    except Exception as e:
        logger.error(f"Error generating suggestions: {e}")
        state.suggestions = []
    
    return state

async def plan_itinerary_node(state: TravelState) -> TravelState:
    """Plan a detailed itinerary"""
    logger.info(f"Planning itinerary for {state.city}")
    
    try:
        # Use planning tool to create itinerary
        itinerary = await planning_tool.create_itinerary(
            city=state.city,
            start_date=state.start_date,
            end_date=state.end_date,
            interests=state.interests,
            budget=state.budget,
            group_size=state.group_size,
            pace=state.pace,
            accommodation=state.accommodation,
            weather=state.weather
        )
        
        state.itinerary = itinerary
        state.messages.append(
            AIMessage(content=f"Created itinerary for {state.city}")
        )
        
    except Exception as e:
        logger.error(f"Error planning itinerary: {e}")
        state.itinerary = {}
    
    return state

async def generate_packing_list_node(state: TravelState) -> TravelState:
    """Generate a packing list"""
    logger.info(f"Generating packing list for {state.destination}")
    
    try:
        # Use LLM to generate packing list
        prompt = f"""
        Generate a comprehensive packing list for {state.destination} based on:
        - Duration: {state.duration} days
        - Activities: {', '.join(state.activities)}
        - Group size: {state.group_size}
        - Age group: {state.age_group}
        - Special needs: {', '.join(state.special_needs)}
        - Weather: {state.weather}
        
        Organize by categories (clothing, toiletries, electronics, etc.) and include quantities.
        """
        
        response = await llm.ainvoke([HumanMessage(content=prompt)])
        
        # Parse packing list (this would be more sophisticated in practice)
        packing_list = parse_packing_list(response.content)
        state.packing_list = packing_list
        
        state.messages.append(
            AIMessage(content=f"Generated packing list for {state.destination}")
        )
        
    except Exception as e:
        logger.error(f"Error generating packing list: {e}")
        state.packing_list = {}
    
    return state

async def finalize_response_node(state: TravelState) -> TravelState:
    """Finalize the response based on the task"""
    logger.info("Finalizing response...")
    
    # Add final message
    state.messages.append(
        AIMessage(content="Response finalized successfully")
    )
    
    return state

def route_by_task(state: TravelState) -> str:
    """Route to the appropriate flow based on task"""
    if state.task == "explore_destination":
        return "explore"
    elif state.task == "generate_itinerary":
        return "itinerary"
    elif state.task == "generate_packing_list":
        return "packing"
    else:
        return "explore"  # default

def calculate_duration(start_date: str, end_date: str) -> int:
    """Calculate duration in days between two dates"""
    try:
        start = datetime.strptime(start_date, "%Y-%m-%d")
        end = datetime.strptime(end_date, "%Y-%m-%d")
        return (end - start).days
    except:
        return 7  # default

def parse_suggestions(content: str) -> List[Dict]:
    """Parse LLM response into structured suggestions"""
    # This is a simplified parser - in practice, you'd want more robust parsing
    suggestions = []
    
    # Split by lines and look for numbered items
    lines = content.split('\n')
    current_suggestion = {}
    
    for line in lines:
        line = line.strip()
        if line.startswith(('1.', '2.', '3.', '4.', '5.', '6.', '7.')):
            if current_suggestion:
                suggestions.append(current_suggestion)
            current_suggestion = {'title': line[3:].strip()}
        elif line and current_suggestion:
            if 'description' not in current_suggestion:
                current_suggestion['description'] = line
            elif 'cost' not in current_suggestion and '$' in line:
                current_suggestion['cost'] = line
    
    if current_suggestion:
        suggestions.append(current_suggestion)
    
    return suggestions

def parse_packing_list(content: str) -> Dict:
    """Parse LLM response into structured packing list"""
    # This is a simplified parser - in practice, you'd want more robust parsing
    packing_list = {
        'categories': [],
        'total_items': 0
    }
    
    # Split by categories
    categories = content.split('\n\n')
    
    for category in categories:
        if ':' in category:
            lines = category.split('\n')
            category_name = lines[0].replace(':', '').strip()
            items = []
            
            for line in lines[1:]:
                if line.strip() and not line.startswith('-'):
                    items.append(line.strip())
            
            if items:
                packing_list['categories'].append({
                    'name': category_name,
                    'items': items
                })
                packing_list['total_items'] += len(items)
    
    return packing_list

class TravelPlanningGraph:
    """Main class for the travel planning graph"""
    
    def __init__(self):
        self.graph = create_travel_graph()
        self.app = self.graph.compile()
    
    async def run(self, input_data: Dict[str, Any]) -> Dict[str, Any]:
        """Run the travel planning graph"""
        try:
            # Create initial state
            state = TravelState()
            
            # Populate state from input
            for key, value in input_data.items():
                if hasattr(state, key):
                    setattr(state, key, value)
            
            # Run the graph
            result = await self.app.ainvoke(state)
            
            # Extract results based on task
            if state.task == "explore_destination":
                return {
                    "suggestions": state.suggestions,
                    "weather": state.weather,
                    "events": state.events,
                    "generated_at": datetime.now().isoformat()
                }
            elif state.task == "generate_itinerary":
                return {
                    "itinerary": state.itinerary,
                    "duration": state.duration,
                    "total_cost": state.itinerary.get("total_cost", 0.0),
                    "generated_at": datetime.now().isoformat()
                }
            elif state.task == "generate_packing_list":
                return {
                    "packing_list": state.packing_list,
                    "weather": state.weather,
                    "total_items": state.packing_list.get("total_items", 0),
                    "generated_at": datetime.now().isoformat()
                }
            else:
                return {
                    "error": "Unknown task",
                    "generated_at": datetime.now().isoformat()
                }
                
        except Exception as e:
            logger.error(f"Error running travel planning graph: {e}")
            return {
                "error": str(e),
                "generated_at": datetime.now().isoformat()
            } 