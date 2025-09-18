#!/usr/bin/env python3
"""
LangGraph Agent for CanTrip - Main Entry Point
Handles itinerary generation and travel planning using LangGraph
"""

import asyncio
import json
import logging
import os
import sys
from typing import Dict, Any, List, Optional
from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import uvicorn
import asyncio

# Add the current directory to Python path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from graph_enhanced import EnhancedTravelPlanningGraph
from tools.recommend import RecommendationTool
from tools.events import EventsTool
from tools.attractions import AttractionsTool
from tools.planner import PlanningTool
from eval.phoenix_adapter import PhoenixAdapter

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(
    title="CanTrip LangGraph Agent",
    description="AI-powered travel planning agent using LangGraph",
    version="1.0.0"
)

# Initialize the enhanced travel planning graph
travel_graph = EnhancedTravelPlanningGraph()

# Initialize Phoenix evaluation
phoenix_adapter = PhoenixAdapter()

class ItineraryRequest(BaseModel):
    city: str
    start_date: str
    end_date: str
    interests: List[str] = []
    budget: float = 1000.0
    group_size: int = 1
    pace: str = "moderate"  # relaxed, moderate, intense
    accommodation: str = "mid-range"  # budget, mid-range, luxury

class ExploreRequest(BaseModel):
    mood: str
    city: str
    budget: float = 1000.0
    duration: int = 7
    interests: List[str] = []

class ChatRequest(BaseModel):
    message: str
    session_id: str
    context: Optional[Dict[str, Any]] = {}
    history: Optional[List[Dict[str, Any]]] = []

class PackingRequest(BaseModel):
    destination: str
    start_date: str
    end_date: str
    activities: List[str] = []
    weather: str = ""
    group_size: int = 1
    age_group: str = "adult"
    special_needs: List[str] = []

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "CanTrip LangGraph Agent"}

@app.post("/generate-itinerary")
async def generate_itinerary(request: ItineraryRequest):
    """Generate a complete travel itinerary"""
    import time
    start_time = time.time()
    
    try:
        logger.info(f"Generating itinerary for {request.city}")
        
        # Prepare the input for the graph
        graph_input = {
            "city": request.city,
            "start_date": request.start_date,
            "end_date": request.end_date,
            "interests": request.interests,
            "budget": request.budget,
            "group_size": request.group_size,
            "pace": request.pace,
            "accommodation": request.accommodation,
            "task": "generate_itinerary"
        }
        
        # Run the graph
        result = await travel_graph.run(graph_input)
        
        response = {
            "success": True,
            "itinerary": result.get("itinerary", {}),
            "metadata": {
                "city": request.city,
                "duration": result.get("duration", 0),
                "total_cost": result.get("total_cost", 0.0),
                "generated_at": result.get("generated_at", "")
            }
        }
        
        # Evaluate performance with Phoenix
        execution_time = time.time() - start_time
        await phoenix_adapter.evaluate_itinerary_generation(
            request=graph_input,
            response=response,
            execution_time=execution_time
        )
        
        return response
        
    except Exception as e:
        logger.error(f"Error generating itinerary: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/explore-destination")
async def explore_destination(request: ExploreRequest):
    """Generate travel suggestions based on mood and interests"""
    import time
    start_time = time.time()
    
    try:
        logger.info(f"Exploring {request.city} for mood: {request.mood}")
        
        # Prepare the input for the graph
        graph_input = {
            "city": request.city,
            "mood": request.mood,
            "budget": request.budget,
            "duration": request.duration,
            "interests": request.interests,
            "task": "explore_destination"
        }
        
        # Run the graph
        result = await travel_graph.run(graph_input)
        
        response = {
            "success": True,
            "suggestions": result.get("suggestions", []),
            "weather": result.get("weather", {}),
            "events": result.get("events", []),
            "metadata": {
                "city": request.city,
                "mood": request.mood,
                "generated_at": result.get("generated_at", "")
            }
        }
        
        # Evaluate performance with Phoenix
        execution_time = time.time() - start_time
        await phoenix_adapter.evaluate_exploration(
            request=graph_input,
            response=response,
            execution_time=execution_time
        )
        
        return response
        
    except Exception as e:
        logger.error(f"Error exploring destination: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/chat")
async def chat_endpoint(request: ChatRequest):
    """Handle conversational chat with the travel agent"""
    try:
        logger.info(f"Processing chat message for session {request.session_id}")
        
        # Prepare the input for the graph
        graph_input = {
            "message": request.message,
            "session_id": request.session_id,
            "context": request.context,
            "history": request.history,
            "task": "chat"
        }
        
        # Run the graph
        result = await travel_graph.run(graph_input)
        
        return {
            "response": result.get("response", "I'm here to help with your travel planning!"),
            "session_id": request.session_id,
            "intent": result.get("intent", "general"),
            "confidence": result.get("confidence", 0.8),
            "suggestions": result.get("suggestions", []),
            "data": result.get("data", {}),
            "timestamp": result.get("timestamp")
        }
    except Exception as e:
        logger.error(f"Error processing chat message: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/chat/stream")
async def chat_stream_endpoint(request: ChatRequest):
    """Handle streaming conversational chat with the travel agent"""
    async def generate_stream():
        try:
            logger.info(f"Processing streaming chat message for session {request.session_id}")
            
            # Prepare the input for the graph
            graph_input = {
                "message": request.message,
                "session_id": request.session_id,
                "context": request.context,
                "history": request.history,
                "task": "chat"
            }
            
            # Run the graph
            result = await travel_graph.run(graph_input)
            response_text = result.get("response", "I'm here to help with your travel planning!")
            
            # Stream the response word by word
            words = response_text.split()
            for i, word in enumerate(words):
                # Create a chunk with the word and metadata
                chunk = {
                    "type": "token",
                    "content": word + (" " if i < len(words) - 1 else ""),
                    "session_id": request.session_id,
                    "intent": result.get("intent", "general"),
                    "confidence": result.get("confidence", 0.8),
                    "timestamp": result.get("timestamp")
                }
                
                yield f"data: {json.dumps(chunk)}\n\n"
                await asyncio.sleep(0.05)  # Small delay to simulate streaming
            
            # Send final metadata
            final_chunk = {
                "type": "done",
                "suggestions": result.get("suggestions", []),
                "data": result.get("data", {}),
                "session_id": request.session_id
            }
            yield f"data: {json.dumps(final_chunk)}\n\n"
            
        except Exception as e:
            logger.error(f"Error processing streaming chat message: {str(e)}")
            error_chunk = {
                "type": "error",
                "content": "I'm sorry, I'm having trouble connecting right now. Please try again in a moment.",
                "session_id": request.session_id
            }
            yield f"data: {json.dumps(error_chunk)}\n\n"
    
    return StreamingResponse(
        generate_stream(),
        media_type="text/plain",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "Content-Type": "text/event-stream"
        }
    )

@app.post("/generate-packing-list")
async def generate_packing_list(request: PackingRequest):
    """Generate a personalized packing list"""
    import time
    start_time = time.time()
    
    try:
        logger.info(f"Generating packing list for {request.destination}")
        
        # Prepare the input for the graph
        graph_input = {
            "destination": request.destination,
            "start_date": request.start_date,
            "end_date": request.end_date,
            "activities": request.activities,
            "weather": request.weather,
            "group_size": request.group_size,
            "age_group": request.age_group,
            "special_needs": request.special_needs,
            "task": "generate_packing_list"
        }
        
        # Run the graph
        result = await travel_graph.run(graph_input)
        
        response = {
            "success": True,
            "packing_list": result.get("packing_list", {}),
            "weather": result.get("weather", {}),
            "metadata": {
                "destination": request.destination,
                "total_items": result.get("total_items", 0),
                "generated_at": result.get("generated_at", "")
            }
        }
        
        # Evaluate performance with Phoenix
        execution_time = time.time() - start_time
        await phoenix_adapter.evaluate_packing_list_generation(
            request=graph_input,
            response=response,
            execution_time=execution_time
        )
        
        return response
        
    except Exception as e:
        logger.error(f"Error generating packing list: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/tools/recommendations")
async def get_recommendations(city: str, category: str = "all"):
    """Get recommendations for a specific city and category"""
    try:
        tool = RecommendationTool()
        recommendations = await tool.get_recommendations(city, category)
        
        return {
            "success": True,
            "city": city,
            "category": category,
            "recommendations": recommendations
        }
        
    except Exception as e:
        logger.error(f"Error getting recommendations: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/tools/events")
async def get_events(city: str, date: str = None):
    """Get events for a specific city and date"""
    try:
        tool = EventsTool()
        events = await tool.get_events(city, date)
        
        return {
            "success": True,
            "city": city,
            "date": date,
            "events": events
        }
        
    except Exception as e:
        logger.error(f"Error getting events: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/tools/attractions")
async def get_attractions(city: str, category: str = "all"):
    """Get attractions for a specific city and category"""
    try:
        tool = AttractionsTool()
        attractions = await tool.get_attractions(city, category)
        
        return {
            "success": True,
            "city": city,
            "category": category,
            "attractions": attractions
        }
        
    except Exception as e:
        logger.error(f"Error getting attractions: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

# Phoenix Evaluation Endpoints
@app.get("/evaluation/summary")
async def get_evaluation_summary(task: str = None, start_date: str = None, end_date: str = None):
    """Get evaluation summary for specified criteria"""
    try:
        summary = await phoenix_adapter.get_evaluation_summary(task, start_date, end_date)
        return {
            "success": True,
            "summary": summary
        }
    except Exception as e:
        logger.error(f"Error getting evaluation summary: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/evaluation/export")
async def export_evaluations(output_file: str = "evaluations_export.json", 
                           task: str = None, start_date: str = None, end_date: str = None):
    """Export evaluations to a file"""
    try:
        await phoenix_adapter.export_evaluations(output_file, task, start_date, end_date)
        return {
            "success": True,
            "message": f"Evaluations exported to {output_file}"
        }
    except Exception as e:
        logger.error(f"Error exporting evaluations: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/evaluation/enable")
async def enable_evaluation():
    """Enable Phoenix evaluation"""
    try:
        phoenix_adapter.enable_evaluation()
        return {
            "success": True,
            "message": "Evaluation enabled"
        }
    except Exception as e:
        logger.error(f"Error enabling evaluation: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/evaluation/disable")
async def disable_evaluation():
    """Disable Phoenix evaluation"""
    try:
        phoenix_adapter.disable_evaluation()
        return {
            "success": True,
            "message": "Evaluation disabled"
        }
    except Exception as e:
        logger.error(f"Error disabling evaluation: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/evaluation/status")
async def get_evaluation_status():
    """Get evaluation status"""
    try:
        return {
            "success": True,
            "enabled": phoenix_adapter.is_evaluation_enabled()
        }
    except Exception as e:
        logger.error(f"Error getting evaluation status: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/evaluation/upload")
async def upload_evaluations_to_phoenix():
    """Upload all evaluations to Phoenix server"""
    try:
        session = await phoenix_adapter.upload_to_phoenix()
        if session:
            return {
                "success": True,
                "message": "Evaluations uploaded to Phoenix successfully",
                "phoenix_url": "http://localhost:6006"
            }
        else:
            return {
                "success": False,
                "message": "Failed to upload evaluations to Phoenix"
            }
    except Exception as e:
        logger.error(f"Error uploading evaluations to Phoenix: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    # Run the FastAPI server
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8001,
        reload=True,
        log_level="info"
    ) 