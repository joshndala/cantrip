'use client';

import React, { useState } from 'react';
// Removed framer-motion imports for cleaner, faster interface
import { MapPin, Sparkles, MessageSquare, User } from 'lucide-react';
import { Sidebar, SidebarContent, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem, SidebarProvider } from '@/components/ui/sidebar';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import Chat from '@/components/Chat';
import { defaultTheme, type CityTheme } from '@/lib/cityThemes';

export default function Home() {
  const [currentTheme, setCurrentTheme] = useState<CityTheme>(defaultTheme);
  const [currentCity, setCurrentCity] = useState<string | null>(null);

  const handleCityDetected = (city: string, theme: CityTheme) => {
    if (city !== currentCity) {
      setCurrentCity(city);
      setCurrentTheme(theme);
    }
  };

  // Removed backgroundVariants - using inline styles for better performance

  return (
    <SidebarProvider>
      <div className="min-h-screen flex">
        {/* Sidebar */}
        <Sidebar className="w-64 border-r bg-white/95 backdrop-blur-sm">
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupLabel>Navigation</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  <SidebarMenuItem>
                    <SidebarMenuButton className="w-full justify-start">
                      <MessageSquare className="w-4 h-4" />
                      <span>Chat</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>
        </Sidebar>

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col relative overflow-hidden">
        {/* Static Background */}
        <div 
          className="absolute inset-0 z-0"
          style={{
            background: currentCity ? currentTheme.background : 'linear-gradient(135deg, #E0E0E0 0%, #FFFFFF 100%)'
          }}
        />

        {/* Header */}
        <header className="relative z-10 p-4 border-b bg-white/80 backdrop-blur-sm">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-urban-blue rounded-lg flex items-center justify-center">
                <Sparkles className="w-5 h-5 text-white" />
              </div>
              <h1 className="text-xl font-bold text-gray-900">CanTrip</h1>
            </div>
            
            <div className="flex items-center space-x-4">
              {currentCity && (
                <div className="flex items-center space-x-2 bg-white/80 backdrop-blur-sm rounded-full px-3 py-1.5 shadow-sm">
                  <MapPin className="w-4 h-4 text-urban-blue" />
                  <span className="text-sm font-medium text-gray-700 capitalize">
                    {currentCity.replace('-', ' ')}
                  </span>
                </div>
              )}
              
              <Avatar className="w-8 h-8">
                <AvatarImage src="/placeholder-avatar.jpg" alt="User" />
                <AvatarFallback>
                  <User className="w-4 h-4" />
                </AvatarFallback>
              </Avatar>
            </div>
          </div>
        </header>

        {/* Chat Area */}
        <main className="relative z-10 flex-1 flex flex-col">
          <div className="flex-1 bg-white/90 backdrop-blur-sm">
            <Chat onCityDetected={handleCityDetected} />
          </div>
        </main>
        </div>
      </div>
    </SidebarProvider>
  );
}