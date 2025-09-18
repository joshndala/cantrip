'use client';

import React, { useState, useRef, useEffect } from 'react';
// Removed framer-motion imports for cleaner, faster chat
import { Send, User, Sparkles } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { detectCityFromMessage, getCityTheme, type CityTheme } from '@/lib/cityThemes';

interface Message {
  id: string;
  content: string;
  role: 'user' | 'assistant';
  timestamp: Date;
}

interface ChatProps {
  onCityDetected?: (city: string, theme: CityTheme) => void;
}

export default function Chat({ onCityDetected }: ChatProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId, setSessionId] = useState<string>('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  // Initialize session and load chat history
  useEffect(() => {
    const initializeSession = async () => {
      // Get or create session ID
      let currentSessionId = localStorage.getItem('cantrip_session_id');
      if (!currentSessionId) {
        currentSessionId = `session_${Date.now()}`;
        localStorage.setItem('cantrip_session_id', currentSessionId);
      }
      setSessionId(currentSessionId);

      // Load existing chat history
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
        const response = await fetch(`${apiUrl}/api/v1/chat/history/${currentSessionId}`);
        
        if (response.ok) {
          const data = await response.json();
          if (data.history && data.history.length > 0) {
            // Convert backend history format to frontend format
            const historyMessages: Message[] = data.history.map((msg: { message: string; role: string; timestamp: string }, index: number) => ({
              id: `${currentSessionId}_${index}`,
              content: msg.message,
              role: msg.role as 'user' | 'assistant',
              timestamp: new Date(msg.timestamp)
            }));
            setMessages(historyMessages);
          } else {
            // No history found, add welcome message
            setMessages([{
              id: '1',
              content: "Hello! I'm your AI Canadian travel assistant. I can help you plan trips, suggest destinations, create itineraries, and more. What would you like to know?",
              role: 'assistant',
              timestamp: new Date()
            }]);
          }
        } else {
          // API call failed, add welcome message
          setMessages([{
            id: '1',
            content: "Hello! I'm your AI Canadian travel assistant. I can help you plan trips, suggest destinations, create itineraries, and more. What would you like to know?",
            role: 'assistant',
            timestamp: new Date()
          }]);
        }
      } catch (error) {
        console.error('Failed to load chat history:', error);
        // Fallback to welcome message
        setMessages([{
          id: '1',
          content: "Hello! I'm your AI Canadian travel assistant. I can help you plan trips, suggest destinations, create itineraries, and more. What would you like to know?",
          role: 'assistant',
          timestamp: new Date()
        }]);
      }
    };

    initializeSession();
  }, []);


  const sendMessage = async () => {
    if (!inputValue.trim() || isLoading || !sessionId) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      content: inputValue.trim(),
      role: 'user',
      timestamp: new Date()
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setIsLoading(true);
    
    // Scroll to bottom after adding user message
    setTimeout(() => scrollToBottom(), 100);

    try {
      // Detect city from user message
      const detectedCity = detectCityFromMessage(userMessage.content);
      
      // Create a placeholder assistant message for streaming
      const assistantMessageId = (Date.now() + 1).toString();
      const assistantMessage: Message = {
        id: assistantMessageId,
        content: '',
        role: 'assistant',
        timestamp: new Date()
      };

      setMessages(prev => [...prev, assistantMessage]);
      
      // Call streaming backend API
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/chat/stream`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          message: userMessage.content,
          session_id: sessionId,
          user_id: 'user_1'
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to send message');
      }

      // Handle streaming response
      const reader = response.body?.getReader();
      const decoder = new TextDecoder();
      let fullResponse = '';

      if (reader) {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const chunk = decoder.decode(value);
          const lines = chunk.split('\n');

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              try {
                const data = JSON.parse(line.slice(6));
                
                if (data.type === 'token') {
                  fullResponse += data.content;
                  
                  // Update the assistant message with the new content
                  setMessages(prev => prev.map(msg => 
                    msg.id === assistantMessageId 
                      ? { ...msg, content: fullResponse }
                      : msg
                  ));
                  
                  // Scroll to bottom as content streams
                  setTimeout(() => scrollToBottom(), 10);
                } else if (data.type === 'done') {
                  // Streaming complete
                  break;
                } else if (data.type === 'error') {
                  throw new Error(data.content);
                }
              } catch (parseError) {
                console.error('Error parsing streaming data:', parseError);
              }
            }
          }
        }
      }

      // If city was detected, notify parent component
      if (detectedCity && onCityDetected) {
        const theme = getCityTheme(detectedCity);
        onCityDetected(detectedCity, theme);
      }

    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: "I'm sorry, I'm having trouble connecting right now. Please try again in a moment.",
        role: 'assistant',
        timestamp: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const startNewChat = async () => {
    // Generate new session ID
    const newSessionId = `session_${Date.now()}`;
    
    // Update localStorage
    localStorage.setItem('cantrip_session_id', newSessionId);
    setSessionId(newSessionId);
    
    // Reset to welcome message
    setMessages([{
      id: '1',
      content: "Hello! I'm your AI Canadian travel assistant. I can help you plan trips, suggest destinations, create itineraries, and more. What would you like to know?",
      role: 'assistant',
      timestamp: new Date()
    }]);
    
    // Scroll to bottom after resetting messages
    setTimeout(() => scrollToBottom(), 100);
  };

  const MessageBubble = ({ message }: { message: Message }) => {
    const isUser = message.role === 'user';
    
    return (
      <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-6 px-4`}>
        <div className={`flex items-start space-x-3 max-w-[80%] ${isUser ? 'flex-row-reverse space-x-reverse' : ''}`}>
          {/* Avatar */}
          <Avatar className="w-8 h-8 flex-shrink-0">
            {isUser ? (
              <>
                <AvatarImage src="/placeholder-user.jpg" alt="User" />
                <AvatarFallback className="bg-urban-blue text-white">
                  <User size={16} />
                </AvatarFallback>
              </>
            ) : (
              <AvatarFallback className="bg-green-500 text-white">
                <Sparkles size={16} />
              </AvatarFallback>
            )}
          </Avatar>
          
          {/* Message Bubble */}
          <div className={`rounded-2xl px-4 py-3 max-w-full ${
            isUser 
              ? 'bg-urban-blue text-white' 
              : 'bg-gray-100 text-gray-900'
          }`}>
            {isUser ? (
              <p className="text-sm leading-relaxed whitespace-pre-wrap">
                {message.content}
              </p>
            ) : (
              <div className="text-sm leading-relaxed prose prose-sm max-w-none">
                <ReactMarkdown
                  remarkPlugins={[remarkGfm]}
                  components={{
                    p: ({ children }) => <p className="mb-2 last:mb-0">{children}</p>,
                    ul: ({ children }) => <ul className="list-disc list-inside mb-2 space-y-1">{children}</ul>,
                    ol: ({ children }) => <ol className="list-decimal list-inside mb-2 space-y-1">{children}</ol>,
                    li: ({ children }) => <li className="text-sm">{children}</li>,
                    strong: ({ children }) => <strong className="font-semibold text-gray-900">{children}</strong>,
                    em: ({ children }) => <em className="italic">{children}</em>,
                    h1: ({ children }) => <h1 className="text-lg font-bold mb-2">{children}</h1>,
                    h2: ({ children }) => <h2 className="text-base font-bold mb-2">{children}</h2>,
                    h3: ({ children }) => <h3 className="text-sm font-bold mb-1">{children}</h3>,
                    blockquote: ({ children }) => <blockquote className="border-l-2 border-gray-300 pl-2 italic">{children}</blockquote>,
                    code: ({ children }) => <code className="bg-gray-200 px-1 py-0.5 rounded text-xs">{children}</code>,
                    pre: ({ children }) => <pre className="bg-gray-200 p-2 rounded text-xs overflow-x-auto">{children}</pre>,
                  }}
                >
                  {message.content}
                </ReactMarkdown>
              </div>
            )}
            <p className={`text-xs mt-2 ${
              isUser ? 'text-blue-100' : 'text-gray-500'
            }`}>
              {message.timestamp.toLocaleTimeString([], { 
                hour: '2-digit', 
                minute: '2-digit' 
              })}
            </p>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="flex flex-col h-full">
      {/* Messages Container */}
      <div className="flex-1 overflow-y-auto py-4">
        {messages.map((message) => (
          <MessageBubble key={message.id} message={message} />
        ))}
        
        {isLoading && (
          <div className="flex justify-start mb-6 px-4">
            <div className="flex items-start space-x-3 max-w-[80%]">
              <Avatar className="w-8 h-8 flex-shrink-0">
                <AvatarFallback className="bg-green-500 text-white">
                  <Sparkles size={16} />
                </AvatarFallback>
              </Avatar>
              <div className="bg-gray-100 rounded-2xl px-4 py-3">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full" />
                  <div className="w-2 h-2 bg-gray-400 rounded-full" />
                  <div className="w-2 h-2 bg-gray-400 rounded-full" />
                </div>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

      {/* Input Container */}
      <div className="border-t border-gray-200 p-4 bg-white">
        <div className="flex items-center space-x-3">
          <div className="flex-1 relative">
            <Input
              ref={inputRef}
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Type your message..."
              className="w-full rounded-full border-gray-300 focus:border-urban-blue focus:ring-urban-blue pr-12"
              disabled={isLoading}
            />
            <Button
              onClick={sendMessage}
              disabled={!inputValue.trim() || isLoading}
              className="absolute right-1 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-urban-blue hover:bg-urban-blue/90 text-white p-0"
            >
              <Send size={16} />
            </Button>
          </div>
          <Button
            onClick={startNewChat}
            disabled={isLoading}
            variant="outline"
            className="border-gray-300 text-gray-600 hover:bg-gray-50 px-4 rounded-full"
            title="Start a new chat"
          >
            New Chat
          </Button>
        </div>
      </div>
    </div>
  );
}
