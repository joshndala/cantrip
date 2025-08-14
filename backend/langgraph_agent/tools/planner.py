#!/usr/bin/env python3
"""
Planning Tool for CanTrip LangGraph Agent
Handles itinerary planning and optimization
"""

import asyncio
import json
import logging
import os
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional
import aiohttp

logger = logging.getLogger(__name__)

class PlanningTool:
    """Tool for planning and optimizing travel itineraries"""
    
    def __init__(self):
        self.api_keys = {
            "google_maps": os.getenv("GOOGLE_MAPS_API_KEY"),
            "uber": os.getenv("UBER_API_KEY"),
            "lyft": os.getenv("LYFT_API_KEY")
        }
        
        # Load planning rules and constraints
        self.planning_rules = self._load_planning_rules()
    
    def _load_planning_rules(self) -> Dict:
        """Load planning rules and constraints"""
        return {
            "max_activities_per_day": 6,
            "min_break_time": 30,  # minutes
            "max_travel_time": 120,  # minutes
            "meal_times": {
                "breakfast": "8:00-10:00",
                "lunch": "12:00-14:00",
                "dinner": "18:00-20:00"
            },
            "activity_durations": {
                "museum": 120,
                "park": 90,
                "restaurant": 60,
                "shopping": 90,
                "tour": 180,
                "attraction": 120
            },
            "transport_modes": {
                "walking": {"speed": 5, "cost": 0},  # km/h
                "public_transit": {"speed": 25, "cost": 3.5},
                "taxi": {"speed": 30, "cost": 15},
                "uber": {"speed": 30, "cost": 12}
            }
        }
    
    async def create_itinerary(self, city: str, start_date: str, end_date: str, 
                             interests: List[str], budget: float, group_size: int,
                             pace: str, accommodation: str, weather: Dict) -> Dict:
        """Create a complete travel itinerary"""
        logger.info(f"Creating itinerary for {city} from {start_date} to {end_date}")
        
        try:
            # Calculate trip duration
            duration = self._calculate_duration(start_date, end_date)
            
            # Get available activities and attractions
            activities = await self._get_available_activities(city, interests, budget)
            
            # Plan daily itineraries
            daily_plans = []
            total_cost = 0
            
            for day in range(1, duration + 1):
                day_date = self._add_days_to_date(start_date, day - 1)
                day_plan = await self._plan_day(
                    city, day, day_date, activities, budget, group_size, 
                    pace, weather
                )
                daily_plans.append(day_plan)
                total_cost += day_plan.get("total_cost", 0)
            
            # Create summary
            summary = self._create_itinerary_summary(city, duration, total_cost, interests)
            
            itinerary = {
                "city": city,
                "start_date": start_date,
                "end_date": end_date,
                "duration": duration,
                "group_size": group_size,
                "pace": pace,
                "accommodation": accommodation,
                "total_cost": total_cost,
                "summary": summary,
                "days": daily_plans,
                "created_at": datetime.now().isoformat()
            }
            
            return itinerary
            
        except Exception as e:
            logger.error(f"Error creating itinerary: {e}")
            return {}
    
    def _calculate_duration(self, start_date: str, end_date: str) -> int:
        """Calculate duration in days between two dates"""
        try:
            start = datetime.strptime(start_date, "%Y-%m-%d")
            end = datetime.strptime(end_date, "%Y-%m-%d")
            return (end - start).days + 1
        except:
            return 7  # default
    
    def _add_days_to_date(self, date_str: str, days: int) -> str:
        """Add days to a date string"""
        try:
            date = datetime.strptime(date_str, "%Y-%m-%d")
            new_date = date + timedelta(days=days)
            return new_date.strftime("%Y-%m-%d")
        except:
            return date_str
    
    async def _get_available_activities(self, city: str, interests: List[str], budget: float) -> List[Dict]:
        """Get available activities based on interests and budget"""
        try:
            # This would integrate with other tools to get activities
            # For now, return sample activities
            activities = [
                {
                    "name": "City Museum",
                    "type": "museum",
                    "duration": 120,
                    "cost": 25,
                    "location": "Downtown",
                    "rating": 4.5,
                    "category": "cultural",
                    "hours": "9:00-17:00",
                    "description": "Interesting museum about city history"
                },
                {
                    "name": "Central Park",
                    "type": "park",
                    "duration": 90,
                    "cost": 0,
                    "location": "City Center",
                    "rating": 4.3,
                    "category": "outdoor",
                    "hours": "6:00-22:00",
                    "description": "Beautiful city park for relaxation"
                },
                {
                    "name": "Local Restaurant",
                    "type": "restaurant",
                    "duration": 60,
                    "cost": 45,
                    "location": "Downtown",
                    "rating": 4.2,
                    "category": "food",
                    "hours": "11:00-23:00",
                    "description": "Authentic local cuisine"
                },
                {
                    "name": "Shopping District",
                    "type": "shopping",
                    "duration": 90,
                    "cost": 0,
                    "location": "Downtown",
                    "rating": 4.0,
                    "category": "shopping",
                    "hours": "10:00-21:00",
                    "description": "Popular shopping area with various stores"
                }
            ]
            
            # Filter by interests and budget
            filtered_activities = []
            for activity in activities:
                if (activity["cost"] <= budget and 
                    (not interests or any(interest.lower() in activity["category"].lower() 
                                        for interest in interests))):
                    filtered_activities.append(activity)
            
            return filtered_activities
            
        except Exception as e:
            logger.error(f"Error getting available activities: {e}")
            return []
    
    async def _plan_day(self, city: str, day: int, date: str, activities: List[Dict],
                       budget: float, group_size: int, pace: str, weather: Dict) -> Dict:
        """Plan activities for a single day"""
        try:
            # Adjust activity count based on pace
            max_activities = self._get_max_activities_for_pace(pace)
            
            # Select activities for the day
            selected_activities = self._select_activities_for_day(
                activities, max_activities, budget, weather
            )
            
            # Schedule activities
            scheduled_activities = self._schedule_activities(selected_activities, date)
            
            # Add meals
            meals = self._plan_meals(city, date, budget, group_size)
            
            # Add transportation
            transport = self._plan_transportation(scheduled_activities)
            
            # Calculate day cost
            day_cost = self._calculate_day_cost(scheduled_activities, meals, transport, group_size)
            
            day_plan = {
                "day": day,
                "date": date,
                "activities": scheduled_activities,
                "meals": meals,
                "transport": transport,
                "total_cost": day_cost,
                "notes": self._generate_day_notes(weather, selected_activities)
            }
            
            return day_plan
            
        except Exception as e:
            logger.error(f"Error planning day {day}: {e}")
            return {
                "day": day,
                "date": date,
                "activities": [],
                "meals": [],
                "transport": [],
                "total_cost": 0,
                "notes": "Error planning this day"
            }
    
    def _get_max_activities_for_pace(self, pace: str) -> int:
        """Get maximum activities based on pace preference"""
        pace_limits = {
            "relaxed": 3,
            "moderate": 5,
            "intense": 8
        }
        return pace_limits.get(pace.lower(), 5)
    
    def _select_activities_for_day(self, activities: List[Dict], max_activities: int,
                                 budget: float, weather: Dict) -> List[Dict]:
        """Select activities for a day based on constraints"""
        # Simple selection algorithm - in production, this would be more sophisticated
        selected = []
        remaining_budget = budget
        
        # Sort activities by rating (descending)
        sorted_activities = sorted(activities, key=lambda x: x.get("rating", 0), reverse=True)
        
        for activity in sorted_activities:
            if len(selected) >= max_activities:
                break
            
            # Check budget
            if activity["cost"] <= remaining_budget:
                # Check weather compatibility
                if self._is_weather_compatible(activity, weather):
                    selected.append(activity)
                    remaining_budget -= activity["cost"]
        
        return selected
    
    def _is_weather_compatible(self, activity: Dict, weather: Dict) -> bool:
        """Check if activity is compatible with weather"""
        if not weather:
            return True
        
        condition = weather.get("condition", "").lower()
        temperature = weather.get("temperature", 20)
        
        # Outdoor activities in bad weather
        if activity["category"] == "outdoor":
            if "rain" in condition or "snow" in condition:
                return False
            if temperature < 5 or temperature > 35:
                return False
        
        return True
    
    def _schedule_activities(self, activities: List[Dict], date: str) -> List[Dict]:
        """Schedule activities throughout the day"""
        scheduled = []
        current_time = datetime.strptime("09:00", "%H:%M")  # Start at 9 AM
        
        for i, activity in enumerate(activities):
            # Add travel time if not first activity
            if i > 0:
                travel_time = 30  # minutes
                current_time += timedelta(minutes=travel_time)
            
            # Schedule activity
            start_time = current_time
            duration_minutes = activity["duration"]
            end_time = start_time + timedelta(minutes=duration_minutes)
            
            scheduled_activity = {
                "name": activity["name"],
                "type": activity["type"],
                "description": activity["description"],
                "location": activity["location"],
                "start_time": start_time.strftime("%H:%M"),
                "end_time": end_time.strftime("%H:%M"),
                "duration": duration_minutes,
                "cost": activity["cost"],
                "category": activity["category"],
                "rating": activity["rating"]
            }
            
            scheduled.append(scheduled_activity)
            current_time = end_time
            
            # Add break time
            if i < len(activities) - 1:
                break_time = 30  # minutes
                current_time += timedelta(minutes=break_time)
        
        return scheduled
    
    def _plan_meals(self, city: str, date: str, budget: float, group_size: int) -> List[Dict]:
        """Plan meals for the day"""
        meals = []
        
        meal_types = ["breakfast", "lunch", "dinner"]
        meal_times = ["08:00", "12:30", "19:00"]
        meal_costs = [15, 25, 35]  # per person
        
        for i, meal_type in enumerate(meal_types):
            meal = {
                "type": meal_type,
                "name": f"{meal_type.title()} at Local Restaurant",
                "location": f"Downtown {city}",
                "time": meal_times[i],
                "cost": meal_costs[i] * group_size,
                "cuisine": "Local specialties",
                "reservation": False
            }
            meals.append(meal)
        
        return meals
    
    def _plan_transportation(self, activities: List[Dict]) -> List[Dict]:
        """Plan transportation between activities"""
        transport = []
        
        for i in range(len(activities) - 1):
            current = activities[i]
            next_activity = activities[i + 1]
            
            # Simple transportation planning
            transport_mode = "walking" if self._is_walking_distance(current, next_activity) else "public_transit"
            
            transport_info = {
                "type": transport_mode,
                "from": current["location"],
                "to": next_activity["location"],
                "start_time": current["end_time"],
                "end_time": next_activity["start_time"],
                "cost": self.planning_rules["transport_modes"][transport_mode]["cost"],
                "duration": 30  # minutes
            }
            
            transport.append(transport_info)
        
        return transport
    
    def _is_walking_distance(self, activity1: Dict, activity2: Dict) -> bool:
        """Check if two activities are within walking distance"""
        # Simple heuristic - in production, would use actual distance calculation
        return activity1["location"] == activity2["location"]
    
    def _calculate_day_cost(self, activities: List[Dict], meals: List[Dict], 
                          transport: List[Dict], group_size: int) -> float:
        """Calculate total cost for the day"""
        activity_cost = sum(activity["cost"] for activity in activities)
        meal_cost = sum(meal["cost"] for meal in meals)
        transport_cost = sum(t["cost"] for t in transport) * group_size
        
        return activity_cost + meal_cost + transport_cost
    
    def _generate_day_notes(self, weather: Dict, activities: List[Dict]) -> str:
        """Generate notes for the day"""
        notes = []
        
        if weather:
            condition = weather.get("condition", "")
            temperature = weather.get("temperature", 20)
            notes.append(f"Weather: {condition}, {temperature}°C")
        
        if activities:
            outdoor_count = sum(1 for a in activities if a["category"] == "outdoor")
            if outdoor_count > 0:
                notes.append(f"{outdoor_count} outdoor activities planned")
        
        notes.append("Remember to bring comfortable walking shoes")
        
        return "; ".join(notes)
    
    def _create_itinerary_summary(self, city: str, duration: int, total_cost: float, 
                                interests: List[str]) -> str:
        """Create a summary of the itinerary"""
        summary = f"{duration}-day trip to {city}"
        
        if interests:
            summary += f" focusing on {', '.join(interests)}"
        
        summary += f". Total estimated cost: ${total_cost:.2f}"
        
        return summary
    
    async def optimize_itinerary(self, itinerary: Dict, preferences: Dict) -> Dict:
        """Optimize an existing itinerary based on preferences"""
        logger.info("Optimizing itinerary based on preferences")
        
        try:
            # This would implement optimization algorithms
            # For now, return the original itinerary
            return itinerary
            
        except Exception as e:
            logger.error(f"Error optimizing itinerary: {e}")
            return itinerary
    
    async def get_alternative_plans(self, city: str, duration: int, budget: float,
                                  interests: List[str], exclude_activities: List[str] = None) -> List[Dict]:
        """Get alternative itinerary plans"""
        logger.info(f"Getting alternative plans for {city}")
        
        try:
            # Generate 2-3 alternative plans
            alternative_plans = []
            
            for i in range(3):
                # Modify interests slightly for variety
                modified_interests = interests.copy()
                if i > 0:
                    # Add some variety
                    additional_interests = ["cultural", "outdoor", "food", "shopping"]
                    for interest in additional_interests:
                        if interest not in modified_interests:
                            modified_interests.append(interest)
                            break
                
                # Create alternative plan
                start_date = datetime.now().strftime("%Y-%m-%d")
                end_date = (datetime.now() + timedelta(days=duration-1)).strftime("%Y-%m-%d")
                
                plan = await self.create_itinerary(
                    city, start_date, end_date, modified_interests, 
                    budget, 1, "moderate", "mid-range", {}
                )
                
                if plan:
                    plan["plan_type"] = f"Alternative {i+1}"
                    alternative_plans.append(plan)
            
            return alternative_plans
            
        except Exception as e:
            logger.error(f"Error getting alternative plans: {e}")
            return [] 