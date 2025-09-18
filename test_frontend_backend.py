#!/usr/bin/env python3
"""
Test script to verify frontend-backend connection
"""

import requests
import json

def test_backend_health():
    """Test if backend is running and healthy"""
    try:
        response = requests.get("http://localhost:8080/api/v1/health")
        print(f"✅ Backend Health: {response.status_code}")
        print(f"   Response: {response.json()}")
        return True
    except Exception as e:
        print(f"❌ Backend Health Failed: {e}")
        return False

def test_chat_endpoint():
    """Test the chat endpoint that frontend uses"""
    try:
        data = {
            "message": "Hello, I want to visit Toronto",
            "session_id": "test-session-123",
            "user_id": "test-user"
        }
        
        response = requests.post(
            "http://localhost:8080/api/v1/chat",
            json=data,
            headers={"Content-Type": "application/json"}
        )
        
        print(f"✅ Chat Endpoint: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Response: {result.get('response', 'No response')[:100]}...")
            print(f"   Session ID: {result.get('session_id', 'No session')}")
        else:
            print(f"   Error: {response.text}")
        return response.status_code == 200
        
    except Exception as e:
        print(f"❌ Chat Endpoint Failed: {e}")
        return False

def test_cors_headers():
    """Test CORS headers for frontend compatibility"""
    try:
        # Test OPTIONS request (preflight)
        response = requests.options("http://localhost:8080/api/v1/chat")
        print(f"✅ CORS Preflight: {response.status_code}")
        
        # Check CORS headers
        cors_headers = {
            'Access-Control-Allow-Origin': response.headers.get('Access-Control-Allow-Origin'),
            'Access-Control-Allow-Methods': response.headers.get('Access-Control-Allow-Methods'),
            'Access-Control-Allow-Headers': response.headers.get('Access-Control-Allow-Headers')
        }
        
        print(f"   CORS Headers: {cors_headers}")
        return cors_headers['Access-Control-Allow-Origin'] == '*'
        
    except Exception as e:
        print(f"❌ CORS Test Failed: {e}")
        return False

def main():
    print("🧪 Testing Frontend-Backend Connection")
    print("=" * 50)
    
    # Test backend health
    backend_healthy = test_backend_health()
    print()
    
    # Test chat endpoint
    chat_working = test_chat_endpoint()
    print()
    
    # Test CORS
    cors_working = test_cors_headers()
    print()
    
    # Summary
    print("📊 Connection Summary:")
    print(f"   Backend Health: {'✅' if backend_healthy else '❌'}")
    print(f"   Chat Endpoint: {'✅' if chat_working else '❌'}")
    print(f"   CORS Headers: {'✅' if cors_working else '❌'}")
    
    if backend_healthy and chat_working and cors_working:
        print("\n🎉 Frontend-Backend connection is working!")
        print("   The frontend should be able to communicate with the backend.")
    else:
        print("\n⚠️  Some issues detected. Check the backend service.")
        print("   Make sure the backend is running on port 8080.")

if __name__ == "__main__":
    main()
