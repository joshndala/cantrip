#!/usr/bin/env python3
"""
Model Router for CanTrip
Intelligently routes requests to appropriate Gemini models based on complexity
"""

import os
import logging
from typing import Dict, Any, Optional, List
from dataclasses import dataclass
from langchain_google_vertexai import ChatVertexAI

logger = logging.getLogger(__name__)

@dataclass
class ModelConfig:
    """Configuration for a specific model"""
    name: str
    max_tokens: int
    temperature: float
    is_pro: bool

class ModelRouter:
    """Routes requests to appropriate Gemini models based on complexity"""
    
    def __init__(self):
        # Default models from environment variables
        self.flash_model = os.getenv("GEMINI_FLASH_MODEL", "gemini-1.5-flash")
        self.pro_model = os.getenv("GEMINI_PRO_MODEL", "gemini-1.5-pro")
        self.flash_lite_model = os.getenv("GEMINI_FLASH_LITE_MODEL", "gemini-1.5-flash")
        
        # Model configurations
        self.models = {
            self.flash_model: ModelConfig(
                name=self.flash_model,
                max_tokens=8192,
                temperature=0.7,
                is_pro=False
            ),
            self.pro_model: ModelConfig(
                name=self.pro_model,
                max_tokens=32768,
                temperature=0.7,
                is_pro=True
            ),
            self.flash_lite_model: ModelConfig(
                name=self.flash_lite_model,
                max_tokens=4096,
                temperature=0.7,
                is_pro=False
            )
        }
        
        # Agent default models
        self.agent_defaults = {
            "explore": self.flash_model,
            "itinerary": self.pro_model,
            "packing": self.flash_model,
            "tips": self.flash_model,
            "events": self.flash_model,
            "formatter": self.flash_lite_model,
            "pdf": self.flash_lite_model
        }
        
        # Google Cloud settings
        self.project_id = os.getenv("GOOGLE_CLOUD_PROJECT")
        self.location = os.getenv("VERTEX_AI_LOCATION", "us-central1")
    
    def estimate_tokens(self, text: str) -> int:
        """Estimate token count from text (4 chars ≈ 1 token)"""
        return len(text) // 4
    
    def should_escalate_to_pro(self, 
                              agent_type: str,
                              prompt_text: str = "",
                              cities: int = 1,
                              date_span_days: int = 1,
                              tool_chain_length: int = 1,
                              has_images: bool = False,
                              previous_error: Optional[str] = None) -> bool:
        """Determine if request should escalate to Pro model"""
        
        # Check for previous errors that warrant escalation
        if previous_error and any(error_type in previous_error.lower() 
                                for error_type in ["context too large", "safety block", "tool schema too complex"]):
            logger.info(f"Escalating to Pro due to previous error: {previous_error}")
            return True
        
        # Check prompt length
        estimated_tokens = self.estimate_tokens(prompt_text)
        if estimated_tokens > 40000:
            logger.info(f"Escalating to Pro due to large prompt: {estimated_tokens} estimated tokens")
            return True
        
        # Check complexity factors
        if cities > 2:
            logger.info(f"Escalating to Pro due to multiple cities: {cities}")
            return True
        
        if date_span_days > 5:
            logger.info(f"Escalating to Pro due to long date span: {date_span_days} days")
            return True
        
        if tool_chain_length > 3:
            logger.info(f"Escalating to Pro due to complex tool chain: {tool_chain_length} tools")
            return True
        
        if has_images:
            logger.info("Escalating to Pro due to multimodal content")
            return True
        
        return False
    
    def get_model_for_agent(self, 
                           agent_type: str,
                           **escalation_factors) -> ChatVertexAI:
        """Get appropriate model for agent type with escalation logic"""
        
        # Get default model for agent
        default_model_name = self.agent_defaults.get(agent_type, self.flash_model)
        
        # Check if escalation is needed
        should_escalate = self.should_escalate_to_pro(agent_type, **escalation_factors)
        
        if should_escalate:
            model_name = self.pro_model
            logger.info(f"Escalating {agent_type} agent from {default_model_name} to {model_name}")
        else:
            model_name = default_model_name
            logger.info(f"Using default model {model_name} for {agent_type} agent")
        
        # Get model configuration
        model_config = self.models[model_name]
        
        # Create and return the model
        return ChatVertexAI(
            model_name=model_config.name,
            project=self.project_id,
            location=self.location,
            temperature=model_config.temperature,
            max_output_tokens=model_config.max_tokens,
        )
    
    def get_model_by_name(self, model_name: str) -> ChatVertexAI:
        """Get model by specific name"""
        if model_name not in self.models:
            logger.warning(f"Unknown model {model_name}, falling back to {self.flash_model}")
            model_name = self.flash_model
        
        model_config = self.models[model_name]
        
        return ChatVertexAI(
            model_name=model_config.name,
            project=self.project_id,
            location=self.location,
            temperature=model_config.temperature,
            max_output_tokens=model_config.max_tokens,
        )

# Global router instance
model_router = ModelRouter() 