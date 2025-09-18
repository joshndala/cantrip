#!/usr/bin/env python3
"""
Weather Tool for CanTrip LangGraph Agent
Handles weather information and forecasts
"""

import asyncio
import json
import logging
import os
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional
import aiohttp

logger = logging.getLogger(__name__)

class WeatherTool:
    """Tool for getting weather information and forecasts"""
    
    def __init__(self):
        self.api_keys = {
            "openweather": os.getenv("OPENWEATHER_API_KEY"),
            "weather_api": os.getenv("WEATHER_API_KEY")
        }
        
        # Backend API URL for weather service
        self.backend_url = os.getenv("BACKEND_URL", "http://cantrip-backend:8080")
    
    async def get_current_weather(self, city: str) -> Dict[str, Any]:
        """Get current weather for a city"""
        logger.info(f"Getting current weather for {city}")
        
        try:
            # Call the Go backend weather service
            async with aiohttp.ClientSession() as session:
                url = f"{self.backend_url}/api/v1/weather/current"
                params = {"city": city}
                
                async with session.get(url, params=params) as response:
                    if response.status == 200:
                        weather_data = await response.json()
                        
                        return {
                            "city": city,
                            "temperature": weather_data.get("temperature", 0),
                            "condition": weather_data.get("condition", "Unknown"),
                            "humidity": weather_data.get("humidity", 0),
                            "wind_speed": weather_data.get("wind_speed", 0),
                            "timestamp": datetime.now().isoformat(),
                            "source": "OpenWeather API"
                        }
                    else:
                        logger.error(f"Weather API returned status {response.status}")
                        return self._get_fallback_weather(city)
                        
        except Exception as e:
            logger.error(f"Error getting current weather: {e}")
            return self._get_fallback_weather(city)
    
    async def get_weather_forecast(self, city: str, days: int = 5) -> List[Dict[str, Any]]:
        """Get weather forecast for a city"""
        logger.info(f"Getting {days}-day weather forecast for {city}")
        
        try:
            # Call the Go backend weather forecast service
            async with aiohttp.ClientSession() as session:
                url = f"{self.backend_url}/api/v1/weather/forecast"
                params = {"city": city, "days": str(days)}
                
                async with session.get(url, params=params) as response:
                    if response.status == 200:
                        forecast_data = await response.json()
                        
                        # Convert backend format to our format
                        forecasts = []
                        for day in forecast_data.get("forecast", []):
                            forecasts.append({
                                "date": day.get("date", ""),
                                "high_temp": day.get("high_temp", 0),
                                "low_temp": day.get("low_temp", 0),
                                "condition": day.get("condition", "Unknown"),
                                "humidity": day.get("humidity", 0),
                                "wind_speed": day.get("wind_speed", 0),
                                "precipitation": day.get("precipitation", 0)
                            })
                        
                        return forecasts
                    else:
                        logger.error(f"Weather forecast API returned status {response.status}")
                        return self._get_fallback_forecast(city, days)
                        
        except Exception as e:
            logger.error(f"Error getting weather forecast: {e}")
            return self._get_fallback_forecast(city, days)
    
    async def get_weather_with_notes(self, city: str) -> Dict[str, Any]:
        """Get weather with travel notes and recommendations"""
        logger.info(f"Getting weather with travel notes for {city}")
        
        try:
            # Call the Go backend weather with notes service
            async with aiohttp.ClientSession() as session:
                url = f"{self.backend_url}/api/v1/weather/forecast/with-notes"
                params = {"city": city}
                
                async with session.get(url, params=params) as response:
                    if response.status == 200:
                        weather_notes = await response.json()
                        
                        return {
                            "city": city,
                            "current_weather": weather_notes.get("current", {}),
                            "forecast": weather_notes.get("forecast", []),
                            "travel_notes": weather_notes.get("travel_notes", []),
                            "recommendations": weather_notes.get("recommendations", []),
                            "timestamp": datetime.now().isoformat()
                        }
                    else:
                        logger.error(f"Weather notes API returned status {response.status}")
                        return self._get_fallback_weather_notes(city)
                        
        except Exception as e:
            logger.error(f"Error getting weather with notes: {e}")
            return self._get_fallback_weather_notes(city)
    
    def _get_fallback_weather(self, city: str) -> Dict[str, Any]:
        """Get fallback weather data when API fails"""
        return {
            "city": city,
            "temperature": 20.0,
            "condition": "Partly Cloudy",
            "humidity": 60,
            "wind_speed": 10.0,
            "timestamp": datetime.now().isoformat(),
            "source": "Fallback data",
            "note": "Weather data temporarily unavailable"
        }
    
    def _get_fallback_forecast(self, city: str, days: int) -> List[Dict[str, Any]]:
        """Get fallback forecast data when API fails"""
        forecasts = []
        today = datetime.now()
        
        for i in range(days):
            date = today + timedelta(days=i)
            forecasts.append({
                "date": date.strftime("%Y-%m-%d"),
                "high_temp": 22.0 + (i * 2),
                "low_temp": 15.0 + (i * 1),
                "condition": "Partly Cloudy",
                "humidity": 60,
                "wind_speed": 10.0,
                "precipitation": 0.0
            })
        
        return forecasts
    
    def _get_fallback_weather_notes(self, city: str) -> Dict[str, Any]:
        """Get fallback weather notes when API fails"""
        return {
            "city": city,
            "current_weather": self._get_fallback_weather(city),
            "forecast": self._get_fallback_forecast(city, 5),
            "travel_notes": [
                "Weather data temporarily unavailable",
                "Please check local weather sources for current conditions"
            ],
            "recommendations": [
                "Pack layers for variable weather",
                "Check weather updates before outdoor activities"
            ],
            "timestamp": datetime.now().isoformat()
        }
    
    async def get_weather_summary(self, city: str) -> str:
        """Get a human-readable weather summary"""
        try:
            weather = await self.get_current_weather(city)
            
            temp = weather.get("temperature", 0)
            condition = weather.get("condition", "Unknown")
            humidity = weather.get("humidity", 0)
            wind = weather.get("wind_speed", 0)
            
            summary = f"Current weather in {city}: {temp:.1f}°C, {condition}"
            if humidity > 0:
                summary += f", {humidity}% humidity"
            if wind > 0:
                summary += f", {wind:.1f} km/h winds"
            
            return summary
            
        except Exception as e:
            logger.error(f"Error getting weather summary: {e}")
            return f"Weather information for {city} is currently unavailable."
    
    async def get_travel_weather_advice(self, city: str, activities: List[str] = None) -> Dict[str, Any]:
        """Get weather-based travel advice for specific activities"""
        try:
            weather = await self.get_current_weather(city)
            forecast = await self.get_weather_forecast(city, 3)
            
            advice = {
                "city": city,
                "current_weather": weather,
                "forecast": forecast,
                "activity_advice": {},
                "general_advice": []
            }
            
            # Generate activity-specific advice
            if activities:
                for activity in activities:
                    activity_lower = activity.lower()
                    
                    if any(word in activity_lower for word in ["outdoor", "hiking", "walking", "park"]):
                        if weather.get("condition", "").lower() in ["rain", "storm", "snow"]:
                            advice["activity_advice"][activity] = "Consider indoor alternatives due to weather"
                        else:
                            advice["activity_advice"][activity] = "Great weather for outdoor activities!"
                    
                    elif any(word in activity_lower for word in ["beach", "swimming", "water"]):
                        if weather.get("temperature", 0) < 20:
                            advice["activity_advice"][activity] = "Water activities may be too cold"
                        else:
                            advice["activity_advice"][activity] = "Good weather for water activities"
                    
                    elif any(word in activity_lower for word in ["museum", "indoor", "shopping"]):
                        advice["activity_advice"][activity] = "Perfect for any weather conditions"
            
            # Generate general advice
            temp = weather.get("temperature", 0)
            condition = weather.get("condition", "").lower()
            
            if temp < 10:
                advice["general_advice"].append("Pack warm clothing - it's quite cold")
            elif temp > 25:
                advice["general_advice"].append("Pack light clothing - it's warm")
            
            if "rain" in condition:
                advice["general_advice"].append("Bring an umbrella or rain jacket")
            elif "sun" in condition:
                advice["general_advice"].append("Don't forget sunscreen and sunglasses")
            
            return advice
            
        except Exception as e:
            logger.error(f"Error getting travel weather advice: {e}")
            return {
                "city": city,
                "error": "Unable to provide weather-based travel advice",
                "general_advice": ["Check local weather sources for current conditions"]
            }
