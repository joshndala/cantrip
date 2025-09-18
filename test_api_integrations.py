#!/usr/bin/env python3
"""
Test script for CanTrip API integrations
Tests Ticketmaster, Eventbrite, and TripAdvisor API calls
"""

import requests
import json
import os
from datetime import datetime, timedelta

class CanTripAPITester:
    def __init__(self):
        self.backend_url = "http://localhost:8080"
        self.test_cities = ["Toronto", "Vancouver", "Montreal", "Calgary", "Ottawa"]
        self.test_moods = ["excited", "relaxed", "adventurous", "romantic", "family", "cultural"]
        self.test_interests = ["music", "sports", "arts", "food", "outdoor", "culture"]
    
    def test_events_api(self):
        """Test the events API endpoint"""
        print("🎫 Testing Events API (Ticketmaster + Eventbrite)")
        print("=" * 60)
        
        for city in self.test_cities[:2]:  # Test first 2 cities
            for mood in self.test_moods[:2]:  # Test first 2 moods
                print(f"\n📍 Testing: {city} with mood '{mood}'")
                
                # Test with different interest combinations
                test_cases = [
                    ["music", "entertainment"],
                    ["arts", "culture"],
                    ["sports", "outdoor"],
                    ["food", "dining"]
                ]
                
                for interests in test_cases:
                    try:
                        params = {
                            "city": city,
                            "mood": mood,
                            "interests": interests
                        }
                        
                        response = requests.get(
                            f"{self.backend_url}/api/v1/places/events",
                            params=params,
                            timeout=30
                        )
                        
                        print(f"  🎯 Interests: {', '.join(interests)}")
                        print(f"  📊 Status: {response.status_code}")
                        
                        if response.status_code == 200:
                            data = response.json()
                            events = data.get('events', []) if isinstance(data, dict) else data
                            
                            print(f"  ✅ Found {len(events)} events")
                            
                            # Show first few events
                            for i, event in enumerate(events[:3]):
                                print(f"    {i+1}. {event.get('name', 'Unknown Event')}")
                                print(f"       📅 {event.get('date', 'TBD')} at {event.get('location', 'TBD')}")
                                print(f"       💰 ${event.get('price', 0):.2f} - {event.get('category', 'General')}")
                                
                        else:
                            print(f"  ❌ Error: {response.text}")
                            
                    except Exception as e:
                        print(f"  ❌ Exception: {e}")
                
                print()  # Add spacing between cities
    
    def test_trip_suggestions_api(self):
        """Test the trip suggestions API endpoint"""
        print("🗺️ Testing Trip Suggestions API")
        print("=" * 60)
        
        for city in self.test_cities[:2]:  # Test first 2 cities
            for mood in self.test_moods[:2]:  # Test first 2 moods
                print(f"\n📍 Testing: {city} with mood '{mood}'")
                
                try:
                    params = {
                        "city": city,
                        "mood": mood,
                        "interests": ["culture", "food"],
                        "budget": 1000,
                        "duration": 3
                    }
                    
                    response = requests.get(
                        f"{self.backend_url}/api/v1/places/suggestions",
                        params=params,
                        timeout=30
                    )
                    
                    print(f"  📊 Status: {response.status_code}")
                    
                    if response.status_code == 200:
                        data = response.json()
                        suggestions = data.get('suggestions', []) if isinstance(data, dict) else data
                        
                        print(f"  ✅ Found {len(suggestions)} trip suggestions")
                        
                        # Show suggestions
                        for i, suggestion in enumerate(suggestions[:3]):
                            print(f"    {i+1}. {suggestion.get('title', 'Unknown Suggestion')}")
                            print(f"       📝 {suggestion.get('description', 'No description')[:100]}...")
                            print(f"       💰 ${suggestion.get('estimated_cost', 0):.2f}")
                            print(f"       🏷️ Tags: {', '.join(suggestion.get('tags', []))}")
                            
                    else:
                        print(f"  ❌ Error: {response.text}")
                        
                except Exception as e:
                    print(f"  ❌ Exception: {e}")
                
                print()  # Add spacing between cities
    
    def test_weather_api(self):
        """Test the weather API endpoint"""
        print("🌤️ Testing Weather API")
        print("=" * 60)
        
        for city in self.test_cities[:3]:  # Test first 3 cities
            print(f"\n📍 Testing weather for: {city}")
            
            try:
                params = {"city": city}
                
                response = requests.get(
                    f"{self.backend_url}/api/v1/weather/current",
                    params=params,
                    timeout=30
                )
                
                print(f"  📊 Status: {response.status_code}")
                
                if response.status_code == 200:
                    data = response.json()
                    print(f"  ✅ Weather data received")
                    print(f"  🌡️ Temperature: {data.get('temperature', 'N/A')}°C")
                    print(f"  ☁️ Condition: {data.get('condition', 'N/A')}")
                    print(f"  💨 Humidity: {data.get('humidity', 'N/A')}%")
                    print(f"  🌬️ Wind: {data.get('wind_speed', 'N/A')} km/h")
                    
                else:
                    print(f"  ❌ Error: {response.text}")
                    
            except Exception as e:
                print(f"  ❌ Exception: {e}")
    
    def test_chat_with_events(self):
        """Test chat functionality with event-related queries"""
        print("💬 Testing Chat with Event Queries")
        print("=" * 60)
        
        event_queries = [
            "What events are happening in Toronto this weekend?",
            "I want to see a concert in Vancouver",
            "Are there any cultural events in Montreal?",
            "What sports events can I attend in Calgary?",
            "I'm looking for family-friendly activities in Ottawa"
        ]
        
        for query in event_queries:
            print(f"\n💭 Query: {query}")
            
            try:
                data = {
                    "message": query,
                    "session_id": f"test-session-{datetime.now().timestamp()}",
                    "user_id": "test-user"
                }
                
                response = requests.post(
                    f"{self.backend_url}/api/v1/chat",
                    json=data,
                    headers={"Content-Type": "application/json"},
                    timeout=30
                )
                
                print(f"  📊 Status: {response.status_code}")
                
                if response.status_code == 200:
                    data = response.json()
                    response_text = data.get('response', 'No response')
                    print(f"  ✅ Response: {response_text[:200]}...")
                    
                    # Check if response mentions events or activities
                    if any(word in response_text.lower() for word in ['event', 'concert', 'show', 'activity', 'attraction']):
                        print(f"  🎯 Response contains event-related content!")
                    
                else:
                    print(f"  ❌ Error: {response.text}")
                    
            except Exception as e:
                print(f"  ❌ Exception: {e}")
    
    def test_api_keys_status(self):
        """Check if API keys are configured"""
        print("🔑 Checking API Keys Status")
        print("=" * 60)
        
        # Test if we can get environment info (this would need to be implemented)
        try:
            response = requests.get(f"{self.backend_url}/api/v1/health")
            if response.status_code == 200:
                print("✅ Backend is running")
            else:
                print("❌ Backend is not responding")
        except Exception as e:
            print(f"❌ Cannot connect to backend: {e}")
        
        print("\n📝 Note: API keys are configured in the backend environment variables:")
        print("   - TICKETMASTER_API_KEY")
        print("   - EVENTBRITE_API_KEY") 
        print("   - TRIPADVISOR_API_KEY")
        print("   - OPENWEATHER_API_KEY")
        print("   - GEOAPIFY_API_KEY")
    
    def run_all_tests(self):
        """Run all API integration tests"""
        print("🧪 CanTrip API Integration Test Suite")
        print("=" * 80)
        print(f"🕐 Started at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print()
        
        # Check backend status first
        self.test_api_keys_status()
        print()
        
        # Run API tests
        self.test_events_api()
        print()
        
        self.test_trip_suggestions_api()
        print()
        
        self.test_weather_api()
        print()
        
        self.test_chat_with_events()
        
        print("\n" + "=" * 80)
        print(f"🏁 Tests completed at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("✅ All API integration tests finished!")

if __name__ == "__main__":
    tester = CanTripAPITester()
    tester.run_all_tests()
