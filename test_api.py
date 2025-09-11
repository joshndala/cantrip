#!/usr/bin/env python3
"""
CanTrip API Testing Script
Test your CanTrip API endpoints from the terminal
"""

import requests
import json
import time
from datetime import datetime, timedelta

class CanTripTester:
    def __init__(self):
        self.backend_url = "http://localhost:8080"
        self.agent_url = "http://localhost:8001"
        self.session_id = f"test-session-{int(time.time())}"
    
    def test_health(self):
        """Test health endpoints"""
        print("🏥 Testing Health Endpoints...")
        
        try:
            # Backend health
            response = requests.get(f"{self.backend_url}/api/v1/health")
            print(f"Backend Health: {response.status_code} - {response.json()}")
            
            # Agent health
            response = requests.get(f"{self.agent_url}/health")
            print(f"Agent Health: {response.status_code} - {response.json()}")
            
        except Exception as e:
            print(f"❌ Health check failed: {e}")
    
    def test_explore(self):
        """Test explore destination endpoint"""
        print("\n🗺️ Testing Explore Destination...")
        
        data = {
            "mood": "adventurous",
            "city": "Toronto",
            "budget": 1000,
            "duration": 7,
            "interests": ["outdoor", "culture"]
        }
        
        try:
            response = requests.post(
                f"{self.backend_url}/api/v1/explore",
                json=data,
                headers={"Content-Type": "application/json"}
            )
            
            print(f"Status: {response.status_code}")
            if response.status_code == 200:
                result = response.json()
                print(f"✅ Success! Found {len(result.get('suggestions', []))} suggestions")
                print(f"Weather: {result.get('weather', {})}")
                print(f"Events: {len(result.get('events', []))} events")
            else:
                print(f"❌ Error: {response.text}")
                
        except Exception as e:
            print(f"❌ Explore test failed: {e}")
    
    def test_itinerary(self):
        """Test itinerary generation"""
        print("\n📅 Testing Itinerary Generation...")
        
        start_date = datetime.now() + timedelta(days=30)
        end_date = start_date + timedelta(days=6)
        
        data = {
            "city": "Toronto",
            "start_date": start_date.strftime("%Y-%m-%d"),
            "end_date": end_date.strftime("%Y-%m-%d"),
            "interests": ["culture", "food"],
            "budget": 1500,
            "group_size": 2,
            "pace": "moderate"
        }
        
        try:
            response = requests.post(
                f"{self.backend_url}/api/v1/itinerary",
                json=data,
                headers={"Content-Type": "application/json"}
            )
            
            print(f"Status: {response.status_code}")
            if response.status_code == 200:
                result = response.json()
                print(f"✅ Success! Generated itinerary for {result.get('metadata', {}).get('duration', 0)} days")
                print(f"Total Cost: ${result.get('metadata', {}).get('total_cost', 0)}")
            else:
                print(f"❌ Error: {response.text}")
                
        except Exception as e:
            print(f"❌ Itinerary test failed: {e}")
    
    def test_chat(self):
        """Test chat functionality"""
        print("\n💬 Testing Chat...")
        
        messages = [
            "I want to plan a trip to Toronto",
            "Next month for 5 days",
            "My budget is $1000",
            "I like food and culture"
        ]
        
        for i, message in enumerate(messages, 1):
            print(f"\n--- Message {i}: {message} ---")
            
            data = {
                "message": message,
                "session_id": self.session_id,
                "context": {},
                "history": []
            }
            
            try:
                response = requests.post(
                    f"{self.agent_url}/chat",
                    json=data,
                    headers={"Content-Type": "application/json"}
                )
                
                print(f"Status: {response.status_code}")
                if response.status_code == 200:
                    result = response.json()
                    ai_response = result.get('response', 'No response')
                    print(f"AI: {ai_response[:200]}...")
                else:
                    print(f"❌ Error: {response.text}")
                    
            except Exception as e:
                print(f"❌ Chat test failed: {e}")
    
    def test_phoenix(self):
        """Test Phoenix evaluation"""
        print("\n📊 Testing Phoenix Evaluation...")
        
        try:
            # Check status
            response = requests.get(f"{self.agent_url}/evaluation/status")
            print(f"Phoenix Status: {response.status_code} - {response.json()}")
            
            # Get summary
            response = requests.get(f"{self.agent_url}/evaluation/summary")
            print(f"Evaluation Summary: {response.status_code} - {response.json()}")
            
        except Exception as e:
            print(f"❌ Phoenix test failed: {e}")
    
    def run_all_tests(self):
        """Run all tests"""
        print("🧪 CanTrip API Testing Suite")
        print("=" * 50)
        
        self.test_health()
        self.test_explore()
        self.test_itinerary()
        self.test_chat()
        self.test_phoenix()
        
        print("\n✅ Testing complete!")

if __name__ == "__main__":
    tester = CanTripTester()
    tester.run_all_tests() 