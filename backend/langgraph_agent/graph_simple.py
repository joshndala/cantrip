#!/usr/bin/env python3
"""
Simplified LangGraph Agent Flow for CanTrip
A working version that handles chat requests properly
"""

import asyncio
import json
import logging
import os
from datetime import datetime
from typing import Dict, Any, List, Optional, TypedDict
from langchain_core.messages import HumanMessage, AIMessage
from langgraph.graph import StateGraph, END
from langchain_google_vertexai import ChatVertexAI

# Import model router
from model_router import model_router

logger = logging.getLogger(__name__)

class TravelState(TypedDict):
    """State for travel planning"""
    message: str
    session_id: str
    context: Dict[str, Any]
    history: List[Any]
    response: str
    intent: str
    confidence: float
    suggestions: List[str]

def create_simple_travel_graph() -> StateGraph:
    """Create a simple travel planning graph"""
    
    # Define the state
    workflow = StateGraph(TravelState)
    
    # Add nodes
    workflow.add_node("process_chat", process_chat_node)
    workflow.add_node("finalize_response", finalize_response_node)
    
    # Define edges
    workflow.set_entry_point("process_chat")
    workflow.add_edge("process_chat", "finalize_response")
    workflow.add_edge("finalize_response", END)
    
    return workflow

async def process_chat_node(state: TravelState) -> TravelState:
    """Process conversational chat messages"""
    logger.info("Processing chat message...")
    
    # Extract message and context
    message = state.get('message', '')
    session_id = state.get('session_id', '')
    context = state.get('context', {})
    history = state.get('history', [])
    
    # Get appropriate LLM for chat processing
    try:
        chat_llm = model_router.get_model_for_agent("explore", prompt_text=message)
        
        # Create a comprehensive prompt for travel planning
        prompt = f"""
        You are a helpful Canadian travel assistant. The user has sent this message: "{message}"
        
        Context: {context}
        History: {history}
        
        Please provide a helpful, detailed response about Canadian travel. If they're asking about:
        - Trip planning: Offer specific suggestions and ask for more details
        - Destinations: Provide interesting facts and recommendations
        - Activities: Suggest popular Canadian activities
        - Weather: Mention seasonal considerations
        - General questions: Be friendly and helpful
        
        Keep your response conversational and helpful. Ask follow-up questions to better understand their needs.
        """
        
        # Get response from LLM
        response = await chat_llm.ainvoke([HumanMessage(content=prompt)])
        
        # Analyze intent (simplified)
        if any(word in message.lower() for word in ['plan', 'trip', 'visit', 'go to']):
            intent = "trip_planning"
        elif any(word in message.lower() for word in ['weather', 'temperature', 'climate']):
            intent = "weather_inquiry"
        elif any(word in message.lower() for word in ['pack', 'packing', 'bring']):
            intent = "packing_request"
        elif any(word in message.lower() for word in ['itinerary', 'schedule', 'plan']):
            intent = "itinerary_request"
        else:
            intent = "general_question"
        
        # Generate suggestions based on intent
        if intent == "trip_planning":
            suggestions = [
                "Tell me about popular Canadian destinations",
                "Help me plan a trip to Toronto", 
                "What's the weather like in Vancouver?",
                "I need a budget-friendly trip to Montreal"
            ]
        elif intent == "weather_inquiry":
            suggestions = [
                "What's the weather in Toronto?",
                "How's the weather in Vancouver?", 
                "Weather forecast for Montreal",
                "Best time to visit Canada?"
            ]
        else:
            suggestions = [
                "Tell me about popular Canadian destinations",
                "Help me plan a trip to Toronto",
                "What's the weather like in Vancouver?"
            ]
        
        # Update state
        state["response"] = response.content
        state["intent"] = intent
        state["suggestions"] = suggestions
        state["confidence"] = 0.9
        
    except Exception as e:
        logger.error(f"Error processing chat: {e}")
        state["response"] = "I'm here to help with your Canadian travel planning! What would you like to know?"
        state["intent"] = "general"
        state["suggestions"] = ["Tell me about popular Canadian destinations", "Help me plan a trip"]
        state["confidence"] = 0.5
    
    return state

async def finalize_response_node(state: TravelState) -> TravelState:
    """Finalize the response"""
    logger.info("Finalizing response...")
    
    # Ensure we have a response
    if not state.get("response"):
        state["response"] = "I'm here to help with your travel planning!"
    
    return state

class SimpleTravelPlanningGraph:
    """Simple travel planning graph"""
    
    def __init__(self):
        self.graph = create_simple_travel_graph()
        self.app = self.graph.compile()
    
    async def run(self, input_data: Dict[str, Any]) -> Dict[str, Any]:
        """Run the travel planning graph"""
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
                "suggestions": []
            }
            
            # Run the graph
            result = await self.app.ainvoke(state)
            
            return {
                "response": result.get("response", "I'm here to help with your travel planning!"),
                "session_id": result.get("session_id", ""),
                "intent": result.get("intent", "general"),
                "confidence": result.get("confidence", 0.8),
                "suggestions": result.get("suggestions", []),
                "data": {},
                "timestamp": datetime.now().isoformat()
            }
                
        except Exception as e:
            logger.error(f"Error running travel planning graph: {e}")
            return {
                "response": "I'm here to help with your travel planning! What would you like to know?",
                "session_id": input_data.get("session_id", ""),
                "intent": "general",
                "confidence": 0.5,
                "suggestions": ["Tell me about popular Canadian destinations"],
                "data": {},
                "timestamp": datetime.now().isoformat()
            }
