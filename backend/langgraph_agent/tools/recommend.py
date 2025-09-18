#!/usr/bin/env python3
"""
Recommendation Tool for CanTrip LangGraph Agent
Handles travel recommendations and suggestions
"""

import asyncio
import json
import logging
import os
from typing import Dict, List, Any, Optional
import aiohttp

logger = logging.getLogger(__name__)

class RecommendationTool:
    """Tool for generating travel recommendations"""
    
    def __init__(self):
        self.api_keys = {
            "geoapify": os.getenv("GEOAPIFY_API_KEY"),
            "tripadvisor": os.getenv("TRIPADVISOR_API_KEY"),
            "yelp": os.getenv("YELP_API_KEY")
        }
        
        # Load city metadata
        self.city_data = self._load_city_data()
    
    def _load_city_data(self) -> Dict:
        """Load city metadata from JSON file"""
        try:
            with open("../../data/city_metadata.json", "r") as f:
                return json.load(f)
        except FileNotFoundError:
            logger.warning("City metadata file not found, using empty data")
            return {"cities": []}
    
    async def get_recommendations(self, city: str, category: str = "all") -> List[Dict]:
        """Get recommendations for a specific city and category"""
        logger.info(f"Getting recommendations for {city}, category: {category}")
        
        try:
            recommendations = []
            
            # First try to get recommendations from Go backend
            backend_recommendations = await self._get_backend_recommendations(city, category)
            if backend_recommendations:
                recommendations.extend(backend_recommendations)
            
            # Fallback to city metadata if no backend data
            if not recommendations:
                city_info = self._get_city_info(city)
                if not city_info:
                    return []
            
            # Get recommendations based on category
            if category == "all" or category == "attractions":
                attractions = await self._get_attractions(city)
                recommendations.extend(attractions)
            
            if category == "all" or category == "restaurants":
                restaurants = await self._get_restaurants(city)
                recommendations.extend(restaurants)
            
            if category == "all" or category == "activities":
                activities = await self._get_activities(city)
                recommendations.extend(activities)
            
            if category == "all" or category == "hotels":
                hotels = await self._get_hotels(city)
                recommendations.extend(hotels)
            
            # Add seasonal recommendations
            seasonal = self._get_seasonal_recommendations(city_info)
            recommendations.extend(seasonal)
            
            return recommendations[:20]  # Limit to 20 recommendations
            
        except Exception as e:
            logger.error(f"Error getting recommendations: {e}")
            return []
    
    async def _get_backend_recommendations(self, city: str, category: str) -> List[Dict]:
        """Get recommendations from Go backend"""
        try:
            backend_url = os.getenv("BACKEND_URL", "http://cantrip-backend:8080")
            
            async with aiohttp.ClientSession() as session:
                url = f"{backend_url}/api/v1/places/suggestions"
                params = {
                    "city": city,
                    "mood": "cultural",
                    "interests": [category] if category != "all" else ["culture", "attractions", "food"],
                    "budget": 1000,
                    "duration": 3
                }
                
                async with session.get(url, params=params) as response:
                    if response.status == 200:
                        data = await response.json()
                        suggestions = data.get("suggestions", [])
                        
                        # Convert suggestions to recommendations format
                        recommendations = []
                        for suggestion in suggestions:
                            recommendation = {
                                "name": suggestion.get("title", ""),
                                "description": suggestion.get("description", ""),
                                "category": category,
                                "type": "recommendation",
                                "location": city,
                                "rating": 4.0,
                                "tags": suggestion.get("tags", []),
                                "activities": suggestion.get("activities", []),
                                "estimated_cost": suggestion.get("estimated_cost", 0),
                                "duration": suggestion.get("duration", 1)
                            }
                            recommendations.append(recommendation)
                        
                        logger.info(f"Retrieved {len(recommendations)} recommendations from backend")
                        return recommendations
                    else:
                        logger.error(f"Backend recommendations API returned status {response.status}")
                        
        except Exception as e:
            logger.error(f"Error calling backend recommendations API: {e}")
        
        return []
    
    def _get_city_info(self, city: str) -> Optional[Dict]:
        """Get city information from metadata"""
        for city_data in self.city_data.get("cities", []):
            if city_data["name"].lower() == city.lower():
                return city_data
        return None
    
    async def _get_attractions(self, city: str) -> List[Dict]:
        """Get attractions for a city"""
        try:
            # This would integrate with Geoapify or TripAdvisor API
            # For now, return sample data
            attractions = [
                {
                    "name": f"Famous Museum in {city}",
                    "type": "attraction",
                    "category": "cultural",
                    "description": "A must-visit cultural attraction",
                    "rating": 4.5,
                    "price_range": "$$",
                    "location": f"Downtown {city}",
                    "hours": "9:00 AM - 6:00 PM",
                    "website": f"https://example.com/{city.lower()}-museum"
                },
                {
                    "name": f"Historic Landmark in {city}",
                    "type": "attraction",
                    "category": "historic",
                    "description": "An important historical site",
                    "rating": 4.2,
                    "price_range": "$",
                    "location": f"Old Town {city}",
                    "hours": "10:00 AM - 5:00 PM",
                    "website": f"https://example.com/{city.lower()}-landmark"
                }
            ]
            return attractions
        except Exception as e:
            logger.error(f"Error getting attractions: {e}")
            return []
    
    async def _get_restaurants(self, city: str) -> List[Dict]:
        """Get restaurant recommendations"""
        try:
            # This would integrate with Yelp or TripAdvisor API
            restaurants = [
                {
                    "name": f"Local Cuisine Restaurant",
                    "type": "restaurant",
                    "category": "local",
                    "description": "Authentic local cuisine",
                    "rating": 4.3,
                    "price_range": "$$",
                    "cuisine": "Local specialties",
                    "location": f"Downtown {city}",
                    "hours": "11:00 AM - 10:00 PM",
                    "website": f"https://example.com/{city.lower()}-restaurant"
                },
                {
                    "name": f"Fine Dining Experience",
                    "type": "restaurant",
                    "category": "fine_dining",
                    "description": "Upscale dining experience",
                    "rating": 4.7,
                    "price_range": "$$$",
                    "cuisine": "International",
                    "location": f"Upscale district {city}",
                    "hours": "6:00 PM - 11:00 PM",
                    "website": f"https://example.com/{city.lower()}-fine-dining"
                }
            ]
            return restaurants
        except Exception as e:
            logger.error(f"Error getting restaurants: {e}")
            return []
    
    async def _get_activities(self, city: str) -> List[Dict]:
        """Get activity recommendations"""
        try:
            activities = [
                {
                    "name": f"City Walking Tour",
                    "type": "activity",
                    "category": "guided_tour",
                    "description": "Explore the city with a local guide",
                    "rating": 4.4,
                    "price_range": "$$",
                    "duration": "3 hours",
                    "location": f"Various locations in {city}",
                    "availability": "Daily",
                    "website": f"https://example.com/{city.lower()}-walking-tour"
                },
                {
                    "name": f"Adventure Activity",
                    "type": "activity",
                    "category": "adventure",
                    "description": "Thrilling outdoor adventure",
                    "rating": 4.6,
                    "price_range": "$$$",
                    "duration": "4 hours",
                    "location": f"Outdoor area near {city}",
                    "availability": "Weather dependent",
                    "website": f"https://example.com/{city.lower()}-adventure"
                }
            ]
            return activities
        except Exception as e:
            logger.error(f"Error getting activities: {e}")
            return []
    
    async def _get_hotels(self, city: str) -> List[Dict]:
        """Get hotel recommendations"""
        try:
            hotels = [
                {
                    "name": f"Luxury Hotel {city}",
                    "type": "hotel",
                    "category": "luxury",
                    "description": "5-star luxury accommodation",
                    "rating": 4.8,
                    "price_range": "$$$$",
                    "amenities": ["Spa", "Pool", "Restaurant", "Gym"],
                    "location": f"Downtown {city}",
                    "website": f"https://example.com/{city.lower()}-luxury-hotel"
                },
                {
                    "name": f"Boutique Hotel {city}",
                    "type": "hotel",
                    "category": "boutique",
                    "description": "Charming boutique hotel",
                    "rating": 4.5,
                    "price_range": "$$$",
                    "amenities": ["Restaurant", "Bar", "Free WiFi"],
                    "location": f"Historic district {city}",
                    "website": f"https://example.com/{city.lower()}-boutique-hotel"
                }
            ]
            return hotels
        except Exception as e:
            logger.error(f"Error getting hotels: {e}")
            return []
    
    def _get_seasonal_recommendations(self, city_info: Dict) -> List[Dict]:
        """Get seasonal recommendations based on city data"""
        recommendations = []
        
        if not city_info:
            return recommendations
        
        # Get current season (simplified)
        import datetime
        month = datetime.datetime.now().month
        
        if month in [12, 1, 2]:  # Winter
            season = "winter"
        elif month in [3, 4, 5]:  # Spring
            season = "spring"
        elif month in [6, 7, 8]:  # Summer
            season = "summer"
        else:  # Fall
            season = "fall"
        
        # Get seasonal activities from city metadata
        if season in city_info.get("seasons", {}):
            seasonal_activities = city_info["seasons"][season].get("activities", [])
            
            for activity in seasonal_activities:
                recommendations.append({
                    "name": activity,
                    "type": "activity",
                    "category": "seasonal",
                    "description": f"Seasonal activity for {season}",
                    "rating": 4.0,
                    "price_range": "$$",
                    "season": season,
                    "location": city_info["name"],
                    "availability": f"Best in {season}",
                    "website": f"https://example.com/{city_info['name'].lower()}-{season}-activity"
                })
        
        return recommendations
    
    async def get_mood_based_recommendations(self, city: str, mood: str, interests: List[str]) -> List[Dict]:
        """Get recommendations based on mood and interests"""
        logger.info(f"Getting mood-based recommendations for {city}, mood: {mood}")
        
        try:
            all_recommendations = await self.get_recommendations(city, "all")
            
            # Filter based on mood
            mood_filters = {
                "relaxed": ["spa", "wellness", "quiet", "peaceful"],
                "adventurous": ["adventure", "thrilling", "outdoor", "active"],
                "cultural": ["museum", "historic", "cultural", "art"],
                "romantic": ["romantic", "intimate", "fine_dining", "luxury"],
                "family": ["family", "kid", "educational", "fun"],
                "budget": ["budget", "affordable", "cheap"],
                "luxury": ["luxury", "upscale", "premium", "exclusive"]
            }
            
            mood_keywords = mood_filters.get(mood.lower(), [])
            
            # Filter recommendations
            filtered_recommendations = []
            for rec in all_recommendations:
                # Check if recommendation matches mood
                matches_mood = any(keyword in rec.get("category", "").lower() or 
                                 keyword in rec.get("description", "").lower() 
                                 for keyword in mood_keywords)
                
                # Check if recommendation matches interests
                matches_interests = any(interest.lower() in rec.get("category", "").lower() or 
                                      interest.lower() in rec.get("description", "").lower() 
                                      for interest in interests)
                
                if matches_mood or matches_interests:
                    filtered_recommendations.append(rec)
            
            return filtered_recommendations[:10]  # Return top 10
            
        except Exception as e:
            logger.error(f"Error getting mood-based recommendations: {e}")
            return []
    
    async def get_budget_recommendations(self, city: str, budget: float) -> List[Dict]:
        """Get recommendations within budget"""
        logger.info(f"Getting budget recommendations for {city}, budget: ${budget}")
        
        try:
            all_recommendations = await self.get_recommendations(city, "all")
            
            # Price range mapping
            price_ranges = {
                "$": 50,
                "$$": 150,
                "$$$": 300,
                "$$$$": 1000
            }
            
            # Filter by budget
            budget_recommendations = []
            for rec in all_recommendations:
                price_range = rec.get("price_range", "$$")
                estimated_cost = price_ranges.get(price_range, 150)
                
                if estimated_cost <= budget:
                    budget_recommendations.append(rec)
            
            return budget_recommendations[:15]  # Return top 15
            
        except Exception as e:
            logger.error(f"Error getting budget recommendations: {e}")
            return [] 