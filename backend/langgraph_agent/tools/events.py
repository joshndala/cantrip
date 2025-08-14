#!/usr/bin/env python3
"""
Events Tool for CanTrip LangGraph Agent
Handles event discovery and recommendations
"""

import asyncio
import json
import logging
import os
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional
import aiohttp

logger = logging.getLogger(__name__)

class EventsTool:
    """Tool for discovering and managing events"""
    
    def __init__(self):
        self.api_keys = {
            "ticketmaster": os.getenv("TICKETMASTER_API_KEY"),
            "eventbrite": os.getenv("EVENTBRITE_API_KEY"),
            "seatgeek": os.getenv("SEATGEEK_API_KEY")
        }
        
        # Sample event data (in production, this would come from APIs)
        self.sample_events = self._load_sample_events()
    
    def _load_sample_events(self) -> Dict[str, List[Dict]]:
        """Load sample event data for different cities"""
        return {
            "toronto": [
                {
                    "name": "Toronto International Film Festival",
                    "type": "festival",
                    "category": "arts",
                    "description": "World-renowned film festival showcasing international cinema",
                    "date": "2024-09-05",
                    "end_date": "2024-09-15",
                    "time": "Various times",
                    "location": "Various venues across Toronto",
                    "price_range": "$$$",
                    "tickets_available": True,
                    "booking_url": "https://tiff.net",
                    "rating": 4.8,
                    "tags": ["film", "culture", "international"]
                },
                {
                    "name": "Toronto Blue Jays vs New York Yankees",
                    "type": "sports",
                    "category": "baseball",
                    "description": "MLB game at Rogers Centre",
                    "date": "2024-07-15",
                    "time": "7:07 PM",
                    "location": "Rogers Centre, Toronto",
                    "price_range": "$$",
                    "tickets_available": True,
                    "booking_url": "https://mlb.com/bluejays",
                    "rating": 4.5,
                    "tags": ["sports", "baseball", "mlb"]
                },
                {
                    "name": "CNE (Canadian National Exhibition)",
                    "type": "fair",
                    "category": "entertainment",
                    "description": "Annual fair with rides, food, and entertainment",
                    "date": "2024-08-16",
                    "end_date": "2024-09-02",
                    "time": "10:00 AM - 10:00 PM",
                    "location": "Exhibition Place, Toronto",
                    "price_range": "$$",
                    "tickets_available": True,
                    "booking_url": "https://theex.com",
                    "rating": 4.3,
                    "tags": ["fair", "entertainment", "family"]
                }
            ],
            "vancouver": [
                {
                    "name": "Vancouver International Jazz Festival",
                    "type": "festival",
                    "category": "music",
                    "description": "Annual jazz festival featuring local and international artists",
                    "date": "2024-06-21",
                    "end_date": "2024-07-01",
                    "time": "Various times",
                    "location": "Various venues across Vancouver",
                    "price_range": "$$",
                    "tickets_available": True,
                    "booking_url": "https://coastaljazz.ca",
                    "rating": 4.6,
                    "tags": ["music", "jazz", "festival"]
                },
                {
                    "name": "Vancouver Canucks vs Edmonton Oilers",
                    "type": "sports",
                    "category": "hockey",
                    "description": "NHL game at Rogers Arena",
                    "date": "2024-03-20",
                    "time": "7:00 PM",
                    "location": "Rogers Arena, Vancouver",
                    "price_range": "$$$",
                    "tickets_available": True,
                    "booking_url": "https://nhl.com/canucks",
                    "rating": 4.7,
                    "tags": ["sports", "hockey", "nhl"]
                }
            ],
            "montreal": [
                {
                    "name": "Just for Laughs Festival",
                    "type": "festival",
                    "category": "comedy",
                    "description": "World's largest comedy festival",
                    "date": "2024-07-10",
                    "end_date": "2024-07-29",
                    "time": "Various times",
                    "location": "Various venues across Montreal",
                    "price_range": "$$",
                    "tickets_available": True,
                    "booking_url": "https://hahaha.com",
                    "rating": 4.7,
                    "tags": ["comedy", "festival", "entertainment"]
                },
                {
                    "name": "Montreal International Jazz Festival",
                    "type": "festival",
                    "category": "music",
                    "description": "Largest jazz festival in the world",
                    "date": "2024-06-27",
                    "end_date": "2024-07-06",
                    "time": "Various times",
                    "location": "Various venues across Montreal",
                    "price_range": "$$",
                    "tickets_available": True,
                    "booking_url": "https://montrealjazzfest.com",
                    "rating": 4.8,
                    "tags": ["music", "jazz", "festival"]
                }
            ]
        }
    
    async def get_events(self, city: str, date: Optional[str] = None, category: Optional[str] = None) -> List[Dict]:
        """Get events for a specific city and date"""
        logger.info(f"Getting events for {city}, date: {date}, category: {category}")
        
        try:
            # Get events from sample data
            city_events = self.sample_events.get(city.lower(), [])
            
            # Filter by date if provided
            if date:
                city_events = self._filter_events_by_date(city_events, date)
            
            # Filter by category if provided
            if category:
                city_events = self._filter_events_by_category(city_events, category)
            
            # Add events from APIs (in production)
            api_events = await self._get_api_events(city, date, category)
            city_events.extend(api_events)
            
            return city_events
            
        except Exception as e:
            logger.error(f"Error getting events: {e}")
            return []
    
    def _filter_events_by_date(self, events: List[Dict], target_date: str) -> List[Dict]:
        """Filter events by date"""
        try:
            target = datetime.strptime(target_date, "%Y-%m-%d")
            filtered_events = []
            
            for event in events:
                event_date = datetime.strptime(event["date"], "%Y-%m-%d")
                
                # Check if event is on the target date or within a range
                if "end_date" in event:
                    end_date = datetime.strptime(event["end_date"], "%Y-%m-%d")
                    if target >= event_date and target <= end_date:
                        filtered_events.append(event)
                else:
                    if target == event_date:
                        filtered_events.append(event)
            
            return filtered_events
            
        except Exception as e:
            logger.error(f"Error filtering events by date: {e}")
            return events
    
    def _filter_events_by_category(self, events: List[Dict], category: str) -> List[Dict]:
        """Filter events by category"""
        category_mapping = {
            "music": ["music", "concert", "festival"],
            "sports": ["sports", "game", "match"],
            "arts": ["arts", "theater", "museum", "gallery"],
            "comedy": ["comedy", "standup"],
            "family": ["family", "kids", "children"],
            "food": ["food", "culinary", "wine"],
            "culture": ["culture", "cultural", "heritage"]
        }
        
        target_categories = category_mapping.get(category.lower(), [category.lower()])
        
        filtered_events = []
        for event in events:
            event_category = event.get("category", "").lower()
            event_type = event.get("type", "").lower()
            
            if any(cat in event_category or cat in event_type for cat in target_categories):
                filtered_events.append(event)
        
        return filtered_events
    
    async def _get_api_events(self, city: str, date: Optional[str], category: Optional[str]) -> List[Dict]:
        """Get events from external APIs"""
        events = []
        
        # Ticketmaster API
        if self.api_keys["ticketmaster"]:
            ticketmaster_events = await self._get_ticketmaster_events(city, date, category)
            events.extend(ticketmaster_events)
        
        # Eventbrite API
        if self.api_keys["eventbrite"]:
            eventbrite_events = await self._get_eventbrite_events(city, date, category)
            events.extend(eventbrite_events)
        
        return events
    
    async def _get_ticketmaster_events(self, city: str, date: Optional[str], category: Optional[str]) -> List[Dict]:
        """Get events from Ticketmaster API"""
        try:
            # This would make actual API calls to Ticketmaster
            # For now, return empty list
            return []
        except Exception as e:
            logger.error(f"Error getting Ticketmaster events: {e}")
            return []
    
    async def _get_eventbrite_events(self, city: str, date: Optional[str], category: Optional[str]) -> List[Dict]:
        """Get events from Eventbrite API"""
        try:
            # This would make actual API calls to Eventbrite
            # For now, return empty list
            return []
        except Exception as e:
            logger.error(f"Error getting Eventbrite events: {e}")
            return []
    
    async def get_upcoming_events(self, city: str, days: int = 30) -> List[Dict]:
        """Get upcoming events within specified days"""
        logger.info(f"Getting upcoming events for {city} in next {days} days")
        
        try:
            all_events = await self.get_events(city)
            upcoming_events = []
            
            today = datetime.now()
            end_date = today + timedelta(days=days)
            
            for event in all_events:
                event_date = datetime.strptime(event["date"], "%Y-%m-%d")
                if event_date >= today and event_date <= end_date:
                    upcoming_events.append(event)
            
            # Sort by date
            upcoming_events.sort(key=lambda x: x["date"])
            
            return upcoming_events
            
        except Exception as e:
            logger.error(f"Error getting upcoming events: {e}")
            return []
    
    async def get_popular_events(self, city: str) -> List[Dict]:
        """Get popular events based on ratings and reviews"""
        logger.info(f"Getting popular events for {city}")
        
        try:
            all_events = await self.get_events(city)
            
            # Sort by rating (descending)
            popular_events = sorted(all_events, key=lambda x: x.get("rating", 0), reverse=True)
            
            return popular_events[:10]  # Return top 10
            
        except Exception as e:
            logger.error(f"Error getting popular events: {e}")
            return []
    
    async def search_events(self, city: str, query: str) -> List[Dict]:
        """Search events by keyword"""
        logger.info(f"Searching events in {city} for: {query}")
        
        try:
            all_events = await self.get_events(city)
            matching_events = []
            
            query_lower = query.lower()
            
            for event in all_events:
                # Search in name, description, and tags
                if (query_lower in event["name"].lower() or
                    query_lower in event.get("description", "").lower() or
                    any(query_lower in tag.lower() for tag in event.get("tags", []))):
                    matching_events.append(event)
            
            return matching_events
            
        except Exception as e:
            logger.error(f"Error searching events: {e}")
            return []
    
    async def get_event_details(self, event_id: str) -> Optional[Dict]:
        """Get detailed information about a specific event"""
        logger.info(f"Getting details for event: {event_id}")
        
        try:
            # This would fetch from a database or API
            # For now, return None
            return None
            
        except Exception as e:
            logger.error(f"Error getting event details: {e}")
            return None
    
    async def check_ticket_availability(self, event_id: str) -> Dict:
        """Check ticket availability for an event"""
        logger.info(f"Checking ticket availability for event: {event_id}")
        
        try:
            # This would check with ticket vendors
            # For now, return mock data
            return {
                "available": True,
                "price_range": "$50-$200",
                "seats_remaining": 150,
                "last_updated": datetime.now().isoformat()
            }
            
        except Exception as e:
            logger.error(f"Error checking ticket availability: {e}")
            return {
                "available": False,
                "error": str(e)
            } 