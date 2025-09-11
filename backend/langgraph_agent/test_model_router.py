#!/usr/bin/env python3
"""
Test script for Model Router
"""

import os
import asyncio
from model_router import model_router

async def test_model_router():
    """Test the model router functionality"""
    
    print("🧪 Testing Model Router")
    print("=" * 50)
    
    # Test 1: Default models for each agent
    print("\n1. Testing default model assignments:")
    agents = ["explore", "itinerary", "packing", "tips", "events", "formatter", "pdf"]
    
    for agent in agents:
        model = model_router.get_model_for_agent(agent)
        print(f"   {agent:12} → {model.model_name}")
    
    # Test 2: Escalation logic
    print("\n2. Testing escalation logic:")
    
    # Test with large prompt
    large_prompt = "x" * 200000  # ~50k tokens
    model = model_router.get_model_for_agent("explore", prompt_text=large_prompt)
    print(f"   Large prompt (50k tokens) → {model.model_name}")
    
    # Test with multiple cities
    model = model_router.get_model_for_agent("explore", cities=3)
    print(f"   Multiple cities (3) → {model.model_name}")
    
    # Test with long date span
    model = model_router.get_model_for_agent("explore", date_span_days=7)
    print(f"   Long date span (7 days) → {model.model_name}")
    
    # Test with complex tool chain
    model = model_router.get_model_for_agent("explore", tool_chain_length=4)
    print(f"   Complex tool chain (4 tools) → {model.model_name}")
    
    # Test with images
    model = model_router.get_model_for_agent("explore", has_images=True)
    print(f"   Has images → {model.model_name}")
    
    # Test with previous error
    model = model_router.get_model_for_agent("explore", previous_error="context too large")
    print(f"   Previous error → {model.model_name}")
    
    # Test 3: Environment variable configuration
    print("\n3. Current configuration:")
    print(f"   Flash Model: {model_router.flash_model}")
    print(f"   Pro Model: {model_router.pro_model}")
    print(f"   Flash Lite Model: {model_router.flash_lite_model}")
    print(f"   Project ID: {model_router.project_id}")
    print(f"   Location: {model_router.location}")
    
    print("\n✅ Model Router tests completed!")

if __name__ == "__main__":
    asyncio.run(test_model_router()) 