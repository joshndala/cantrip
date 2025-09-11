#!/bin/bash

# CanTrip Chat Testing Script
# Usage: ./test_chat.sh

SESSION_ID="test-session-$(date +%s)"
API_URL="http://localhost:8001"

echo "🤖 CanTrip Chat Testing"
echo "========================"
echo "Session ID: $SESSION_ID"
echo "API URL: $API_URL"
echo "Type 'quit' to exit"
echo ""

# Initialize empty history
HISTORY="[]"

while true; do
    echo -n "You: "
    read -r message
    
    if [ "$message" = "quit" ]; then
        echo "Goodbye! 👋"
        break
    fi
    
    echo "🤖 Processing..."
    
    # Send request
    response=$(curl -s -X POST "$API_URL/chat" \
        -H "Content-Type: application/json" \
        -d "{
            \"message\": \"$message\",
            \"session_id\": \"$SESSION_ID\",
            \"context\": {},
            \"history\": $HISTORY
        }")
    
    # Extract and display response
    echo "AI: $(echo "$response" | jq -r '.response // .message // "No response"')"
    echo ""
    
    # Update history (simplified - in real implementation you'd append to history)
    HISTORY="[]"
done 