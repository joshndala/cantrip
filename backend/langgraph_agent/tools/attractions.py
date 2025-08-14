#!/usr/bin/env python3
"""
Attractions Tool for CanTrip LangGraph Agent
Handles attraction discovery and recommendations
"""

import asyncio
import json
import logging
import os
from typing import Dict, List, Any, Optional
import aiohttp

logger = logging.getLogger(__name__)

class AttractionsTool:
    """Tool for discovering and managing attractions"""
    
    def __init__(self):
        self.api_keys = {
            "geoapify": os.getenv("GEOAPIFY_API_KEY"),
            "tripadvisor": os.getenv("TRIPADVISOR_API_KEY"),
            "google_places": os.getenv("GOOGLE_PLACES_API_KEY")
        }
        
        # Load city metadata for attractions
        self.city_data = self._load_city_data()
    
    def _load_city_data(self) -> Dict:
        """Load city metadata from JSON file"""
        try:
            with open("../../data/city_metadata.json", "r") as f:
                return json.load(f)
        except FileNotFoundError:
            logger.warning("City metadata file not found, using empty data")
            return {"cities": []}
    
    async def get_attractions(self, city: str, category: str = "all") -> List[Dict]:
        """Get attractions for a specific city and category"""
        logger.info(f"Getting attractions for {city}, category: {category}")
        
        try:
            attractions = []
            
            # Get attractions from city metadata
            city_info = self._get_city_info(city)
            if city_info:
                metadata_attractions = self._get_metadata_attractions(city_info)
                attractions.extend(metadata_attractions)
            
            # Get attractions from APIs
            api_attractions = await self._get_api_attractions(city, category)
            attractions.extend(api_attractions)
            
            # Filter by category if specified
            if category != "all":
                attractions = self._filter_attractions_by_category(attractions, category)
            
            return attractions[:20]  # Limit to 20 attractions
            
        except Exception as e:
            logger.error(f"Error getting attractions: {e}")
            return []
    
    def _get_city_info(self, city: str) -> Optional[Dict]:
        """Get city information from metadata"""
        for city_data in self.city_data.get("cities", []):
            if city_data["name"].lower() == city.lower():
                return city_data
        return None
    
    def _get_metadata_attractions(self, city_info: Dict) -> List[Dict]:
        """Get attractions from city metadata"""
        attractions = []
        
        for attraction_name in city_info.get("attractions", []):
            attractions.append({
                "name": attraction_name,
                "type": "attraction",
                "category": "landmark",
                "description": f"Famous attraction in {city_info['name']}",
                "location": city_info["name"],
                "rating": 4.0,
                "price_range": "$$",
                "hours": "9:00 AM - 6:00 PM",
                "website": f"https://example.com/{city_info['name'].lower()}-{attraction_name.lower().replace(' ', '-')}",
                "tags": ["landmark", "tourist", "popular"]
            })
        
        return attractions
    
    async def _get_api_attractions(self, city: str, category: str) -> List[Dict]:
        """Get attractions from external APIs"""
        attractions = []
        
        # Geoapify API
        if self.api_keys["geoapify"]:
            geoapify_attractions = await self._get_geoapify_attractions(city, category)
            attractions.extend(geoapify_attractions)
        
        # TripAdvisor API
        if self.api_keys["tripadvisor"]:
            tripadvisor_attractions = await self._get_tripadvisor_attractions(city, category)
            attractions.extend(tripadvisor_attractions)
        
        # Google Places API
        if self.api_keys["google_places"]:
            google_attractions = await self._get_google_places_attractions(city, category)
            attractions.extend(google_attractions)
        
        return attractions
    
    async def _get_geoapify_attractions(self, city: str, category: str) -> List[Dict]:
        """Get attractions from Geoapify API"""
        try:
            # This would make actual API calls to Geoapify
            # For now, return sample data
            sample_attractions = [
                {
                    "name": f"Historic District {city}",
                    "type": "attraction",
                    "category": "historic",
                    "description": "Beautiful historic district with preserved architecture",
                    "location": f"Old Town {city}",
                    "rating": 4.3,
                    "price_range": "$",
                    "hours": "Always open",
                    "website": f"https://example.com/{city.lower()}-historic-district",
                    "tags": ["historic", "architecture", "walking"]
                }
            ]
            return sample_attractions
        except Exception as e:
            logger.error(f"Error getting Geoapify attractions: {e}")
            return []
    
    async def _get_tripadvisor_attractions(self, city: str, category: str) -> List[Dict]:
        """Get attractions from TripAdvisor API"""
        try:
            # This would make actual API calls to TripAdvisor
            # For now, return sample data
            sample_attractions = [
                {
                    "name": f"Local Museum {city}",
                    "type": "attraction",
                    "category": "museum",
                    "description": "Interesting local museum showcasing city history",
                    "location": f"Downtown {city}",
                    "rating": 4.1,
                    "price_range": "$$",
                    "hours": "10:00 AM - 5:00 PM",
                    "website": f"https://example.com/{city.lower()}-museum",
                    "tags": ["museum", "culture", "history"]
                }
            ]
            return sample_attractions
        except Exception as e:
            logger.error(f"Error getting TripAdvisor attractions: {e}")
            return []
    
    async def _get_google_places_attractions(self, city: str, category: str) -> List[Dict]:
        """Get attractions from Google Places API"""
        try:
            # This would make actual API calls to Google Places
            # For now, return sample data
            sample_attractions = [
                {
                    "name": f"City Park {city}",
                    "type": "attraction",
                    "category": "park",
                    "description": "Beautiful city park perfect for relaxation",
                    "location": f"Central {city}",
                    "rating": 4.5,
                    "price_range": "Free",
                    "hours": "6:00 AM - 10:00 PM",
                    "website": f"https://example.com/{city.lower()}-park",
                    "tags": ["park", "nature", "free"]
                }
            ]
            return sample_attractions
        except Exception as e:
            logger.error(f"Error getting Google Places attractions: {e}")
            return []
    
    def _filter_attractions_by_category(self, attractions: List[Dict], category: str) -> List[Dict]:
        """Filter attractions by category"""
        category_mapping = {
            "museum": ["museum", "gallery", "exhibit"],
            "park": ["park", "garden", "nature"],
            "historic": ["historic", "landmark", "monument"],
            "entertainment": ["entertainment", "amusement", "fun"],
            "shopping": ["shopping", "market", "mall"],
            "outdoor": ["outdoor", "adventure", "sports"],
            "cultural": ["cultural", "arts", "theater"]
        }
        
        target_categories = category_mapping.get(category.lower(), [category.lower()])
        
        filtered_attractions = []
        for attraction in attractions:
            attraction_category = attraction.get("category", "").lower()
            attraction_type = attraction.get("type", "").lower()
            
            if any(cat in attraction_category or cat in attraction_type for cat in target_categories):
                filtered_attractions.append(attraction)
        
        return filtered_attractions
    
    async def get_popular_attractions(self, city: str) -> List[Dict]:
        """Get popular attractions based on ratings"""
        logger.info(f"Getting popular attractions for {city}")
        
        try:
            all_attractions = await self.get_attractions(city)
            
            # Sort by rating (descending)
            popular_attractions = sorted(all_attractions, key=lambda x: x.get("rating", 0), reverse=True)
            
            return popular_attractions[:10]  # Return top 10
            
        except Exception as e:
            logger.error(f"Error getting popular attractions: {e}")
            return []
    
    async def get_free_attractions(self, city: str) -> List[Dict]:
        """Get free attractions"""
        logger.info(f"Getting free attractions for {city}")
        
        try:
            all_attractions = await self.get_attractions(city)
            
            free_attractions = []
            for attraction in all_attractions:
                price_range = attraction.get("price_range", "").lower()
                if "free" in price_range or "$0" in price_range:
                    free_attractions.append(attraction)
            
            return free_attractions
            
        except Exception as e:
            logger.error(f"Error getting free attractions: {e}")
            return []
    
    async def search_attractions(self, city: str, query: str) -> List[Dict]:
        """Search attractions by keyword"""
        logger.info(f"Searching attractions in {city} for: {query}")
        
        try:
            all_attractions = await self.get_attractions(city)
            matching_attractions = []
            
            query_lower = query.lower()
            
            for attraction in all_attractions:
                # Search in name, description, and tags
                if (query_lower in attraction["name"].lower() or
                    query_lower in attraction.get("description", "").lower() or
                    any(query_lower in tag.lower() for tag in attraction.get("tags", []))):
                    matching_attractions.append(attraction)
            
            return matching_attractions
            
        except Exception as e:
            logger.error(f"Error searching attractions: {e}")
            return []
    
    async def get_attraction_details(self, attraction_id: str) -> Optional[Dict]:
        """Get detailed information about a specific attraction"""
        logger.info(f"Getting details for attraction: {attraction_id}")
        
        try:
            # This would fetch from a database or API
            # For now, return sample data
            return {
                "id": attraction_id,
                "name": "Sample Attraction",
                "description": "Detailed description of the attraction",
                "location": "Sample Location",
                "rating": 4.5,
                "reviews": 150,
                "price_range": "$$",
                "hours": "9:00 AM - 6:00 PM",
                "website": "https://example.com",
                "phone": "+1-555-123-4567",
                "address": "123 Sample Street, City, Province",
                "coordinates": {
                    "lat": 43.6532,
                    "lng": -79.3832
                },
                "photos": [
                    "https://example.com/photo1.jpg",
                    "https://example.com/photo2.jpg"
                ],
                "tags": ["landmark", "tourist", "popular"],
                "accessibility": ["wheelchair_accessible", "parking_available"],
                "best_time_to_visit": "Morning",
                "peak_hours": "2:00 PM - 4:00 PM",
                "tips": [
                    "Visit early to avoid crowds",
                    "Bring comfortable walking shoes",
                    "Check opening hours before visiting"
                ]
            }
            
        except Exception as e:
            logger.error(f"Error getting attraction details: {e}")
            return None
    
    async def get_nearby_attractions(self, lat: float, lng: float, radius: float = 5.0) -> List[Dict]:
        """Get attractions near a specific location"""
        logger.info(f"Getting nearby attractions at ({lat}, {lng}) within {radius}km")
        
        try:
            # This would use geolocation APIs to find nearby attractions
            # For now, return sample data
            nearby_attractions = [
                {
                    "name": "Nearby Park",
                    "type": "attraction",
                    "category": "park",
                    "description": "Beautiful park within walking distance",
                    "location": "Nearby",
                    "rating": 4.2,
                    "price_range": "Free",
                    "distance": "0.5 km",
                    "website": "https://example.com/nearby-park",
                    "tags": ["park", "nearby", "free"]
                }
            ]
            return nearby_attractions
            
        except Exception as e:
            logger.error(f"Error getting nearby attractions: {e}")
            return []
    
    async def get_attraction_reviews(self, attraction_id: str) -> List[Dict]:
        """Get reviews for a specific attraction"""
        logger.info(f"Getting reviews for attraction: {attraction_id}")
        
        try:
            # This would fetch reviews from APIs or database
            # For now, return sample data
            reviews = [
                {
                    "id": "review1",
                    "rating": 5,
                    "title": "Amazing experience!",
                    "content": "This attraction exceeded my expectations. Highly recommended!",
                    "author": "John D.",
                    "date": "2024-01-15",
                    "helpful_votes": 12
                },
                {
                    "id": "review2",
                    "rating": 4,
                    "title": "Good visit",
                    "content": "Nice attraction, worth visiting. Could be better with more information.",
                    "author": "Sarah M.",
                    "date": "2024-01-10",
                    "helpful_votes": 8
                }
            ]
            return reviews
            
        except Exception as e:
            logger.error(f"Error getting attraction reviews: {e}")
            return [] 