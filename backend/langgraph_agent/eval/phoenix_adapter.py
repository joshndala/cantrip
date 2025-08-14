#!/usr/bin/env python3
"""
Phoenix Evaluation Adapter for CanTrip LangGraph Agent
Handles evaluation and monitoring of agent performance
"""

import asyncio
import json
import logging
import os
from datetime import datetime
from typing import Dict, List, Any, Optional
from phoenix.trace import trace
from phoenix.trace.span import Span
from phoenix.trace.attributes import Attributes

logger = logging.getLogger(__name__)

class PhoenixAdapter:
    """Adapter for Phoenix evaluation and monitoring"""
    
    def __init__(self):
        self.trace_client = None
        self.evaluation_enabled = os.getenv("PHOENIX_ENABLED", "false").lower() == "true"
        
        if self.evaluation_enabled:
            self._initialize_phoenix()
    
    def _initialize_phoenix(self):
        """Initialize Phoenix trace client"""
        try:
            # Initialize Phoenix trace client
            # This would be configured based on your Phoenix setup
            logger.info("Phoenix evaluation enabled")
        except Exception as e:
            logger.error(f"Failed to initialize Phoenix: {e}")
            self.evaluation_enabled = False
    
    @trace
    async def evaluate_itinerary_generation(self, request: Dict, response: Dict, 
                                         execution_time: float) -> Dict:
        """Evaluate itinerary generation performance"""
        if not self.evaluation_enabled:
            return {}
        
        try:
            evaluation = {
                "task": "itinerary_generation",
                "request": request,
                "response": response,
                "execution_time": execution_time,
                "timestamp": datetime.now().isoformat(),
                "metrics": self._calculate_itinerary_metrics(request, response, execution_time)
            }
            
            # Log evaluation
            self._log_evaluation(evaluation)
            
            return evaluation
            
        except Exception as e:
            logger.error(f"Error evaluating itinerary generation: {e}")
            return {}
    
    @trace
    async def evaluate_exploration(self, request: Dict, response: Dict, 
                                 execution_time: float) -> Dict:
        """Evaluate exploration performance"""
        if not self.evaluation_enabled:
            return {}
        
        try:
            evaluation = {
                "task": "exploration",
                "request": request,
                "response": response,
                "execution_time": execution_time,
                "timestamp": datetime.now().isoformat(),
                "metrics": self._calculate_exploration_metrics(request, response, execution_time)
            }
            
            # Log evaluation
            self._log_evaluation(evaluation)
            
            return evaluation
            
        except Exception as e:
            logger.error(f"Error evaluating exploration: {e}")
            return {}
    
    @trace
    async def evaluate_packing_list_generation(self, request: Dict, response: Dict, 
                                             execution_time: float) -> Dict:
        """Evaluate packing list generation performance"""
        if not self.evaluation_enabled:
            return {}
        
        try:
            evaluation = {
                "task": "packing_list_generation",
                "request": request,
                "response": response,
                "execution_time": execution_time,
                "timestamp": datetime.now().isoformat(),
                "metrics": self._calculate_packing_metrics(request, response, execution_time)
            }
            
            # Log evaluation
            self._log_evaluation(evaluation)
            
            return evaluation
            
        except Exception as e:
            logger.error(f"Error evaluating packing list generation: {e}")
            return {}
    
    def _calculate_itinerary_metrics(self, request: Dict, response: Dict, 
                                   execution_time: float) -> Dict:
        """Calculate metrics for itinerary generation"""
        metrics = {
            "execution_time_seconds": execution_time,
            "success": "success" in response and response["success"],
            "itinerary_quality": 0.0,
            "cost_accuracy": 0.0,
            "activity_count": 0,
            "day_count": 0
        }
        
        try:
            if "itinerary" in response:
                itinerary = response["itinerary"]
                
                # Calculate itinerary quality
                if "days" in itinerary:
                    metrics["day_count"] = len(itinerary["days"])
                    total_activities = 0
                    for day in itinerary["days"]:
                        if "activities" in day:
                            total_activities += len(day["activities"])
                    metrics["activity_count"] = total_activities
                
                # Calculate cost accuracy
                if "total_cost" in itinerary and "budget" in request:
                    budget = request["budget"]
                    actual_cost = itinerary["total_cost"]
                    if budget > 0:
                        metrics["cost_accuracy"] = min(actual_cost / budget, 2.0)  # Cap at 200%
                
                # Calculate overall quality score
                quality_factors = []
                if metrics["day_count"] > 0:
                    quality_factors.append(min(metrics["day_count"] / 7, 1.0))  # Duration factor
                if metrics["activity_count"] > 0:
                    quality_factors.append(min(metrics["activity_count"] / (metrics["day_count"] * 5), 1.0))  # Activity density
                if metrics["cost_accuracy"] > 0:
                    quality_factors.append(max(0, 1 - abs(1 - metrics["cost_accuracy"])))  # Cost accuracy
                
                if quality_factors:
                    metrics["itinerary_quality"] = sum(quality_factors) / len(quality_factors)
            
        except Exception as e:
            logger.error(f"Error calculating itinerary metrics: {e}")
        
        return metrics
    
    def _calculate_exploration_metrics(self, request: Dict, response: Dict, 
                                     execution_time: float) -> Dict:
        """Calculate metrics for exploration"""
        metrics = {
            "execution_time_seconds": execution_time,
            "success": "success" in response and response["success"],
            "suggestion_count": 0,
            "weather_included": False,
            "events_included": False,
            "relevance_score": 0.0
        }
        
        try:
            if "suggestions" in response:
                metrics["suggestion_count"] = len(response["suggestions"])
            
            if "weather" in response and response["weather"]:
                metrics["weather_included"] = True
            
            if "events" in response and response["events"]:
                metrics["events_included"] = len(response["events"]) > 0
            
            # Calculate relevance score based on mood and interests
            if "mood" in request and "suggestions" in response:
                mood = request["mood"].lower()
                suggestions = response["suggestions"]
                
                relevant_count = 0
                for suggestion in suggestions:
                    if "category" in suggestion:
                        category = suggestion["category"].lower()
                        if mood in category or category in mood:
                            relevant_count += 1
                
                if suggestions:
                    metrics["relevance_score"] = relevant_count / len(suggestions)
            
        except Exception as e:
            logger.error(f"Error calculating exploration metrics: {e}")
        
        return metrics
    
    def _calculate_packing_metrics(self, request: Dict, response: Dict, 
                                 execution_time: float) -> Dict:
        """Calculate metrics for packing list generation"""
        metrics = {
            "execution_time_seconds": execution_time,
            "success": "success" in response and response["success"],
            "item_count": 0,
            "category_count": 0,
            "weather_considered": False,
            "completeness_score": 0.0
        }
        
        try:
            if "packing_list" in response:
                packing_list = response["packing_list"]
                
                if "categories" in packing_list:
                    metrics["category_count"] = len(packing_list["categories"])
                    
                    total_items = 0
                    for category in packing_list["categories"]:
                        if "items" in category:
                            total_items += len(category["items"])
                    metrics["item_count"] = total_items
                
                if "weather" in response and response["weather"]:
                    metrics["weather_considered"] = True
                
                # Calculate completeness score
                completeness_factors = []
                if metrics["category_count"] > 0:
                    completeness_factors.append(min(metrics["category_count"] / 8, 1.0))  # Category coverage
                if metrics["item_count"] > 0:
                    completeness_factors.append(min(metrics["item_count"] / 50, 1.0))  # Item coverage
                if metrics["weather_considered"]:
                    completeness_factors.append(1.0)  # Weather consideration
                
                if completeness_factors:
                    metrics["completeness_score"] = sum(completeness_factors) / len(completeness_factors)
            
        except Exception as e:
            logger.error(f"Error calculating packing metrics: {e}")
        
        return metrics
    
    def _log_evaluation(self, evaluation: Dict):
        """Log evaluation results"""
        try:
            # Log to file
            log_file = "phoenix_evaluations.jsonl"
            with open(log_file, "a") as f:
                f.write(json.dumps(evaluation) + "\n")
            
            # Log to console
            logger.info(f"Evaluation logged: {evaluation['task']} - "
                       f"Success: {evaluation['metrics']['success']}, "
                       f"Time: {evaluation['execution_time']:.2f}s")
            
        except Exception as e:
            logger.error(f"Error logging evaluation: {e}")
    
    async def get_evaluation_summary(self, task: str = None, 
                                   start_date: str = None, 
                                   end_date: str = None) -> Dict:
        """Get evaluation summary for specified criteria"""
        try:
            # Read evaluation logs
            evaluations = []
            log_file = "phoenix_evaluations.jsonl"
            
            if os.path.exists(log_file):
                with open(log_file, "r") as f:
                    for line in f:
                        try:
                            evaluation = json.loads(line.strip())
                            evaluations.append(evaluation)
                        except json.JSONDecodeError:
                            continue
            
            # Filter evaluations
            filtered_evaluations = []
            for evaluation in evaluations:
                # Filter by task
                if task and evaluation.get("task") != task:
                    continue
                
                # Filter by date range
                if start_date or end_date:
                    eval_date = evaluation.get("timestamp", "")
                    if start_date and eval_date < start_date:
                        continue
                    if end_date and eval_date > end_date:
                        continue
                
                filtered_evaluations.append(evaluation)
            
            # Calculate summary statistics
            summary = self._calculate_summary_statistics(filtered_evaluations)
            
            return summary
            
        except Exception as e:
            logger.error(f"Error getting evaluation summary: {e}")
            return {}
    
    def _calculate_summary_statistics(self, evaluations: List[Dict]) -> Dict:
        """Calculate summary statistics from evaluations"""
        if not evaluations:
            return {}
        
        summary = {
            "total_evaluations": len(evaluations),
            "success_rate": 0.0,
            "average_execution_time": 0.0,
            "task_breakdown": {},
            "performance_trends": {}
        }
        
        try:
            success_count = 0
            total_time = 0.0
            task_counts = {}
            
            for evaluation in evaluations:
                metrics = evaluation.get("metrics", {})
                
                # Success rate
                if metrics.get("success", False):
                    success_count += 1
                
                # Execution time
                total_time += metrics.get("execution_time_seconds", 0)
                
                # Task breakdown
                task = evaluation.get("task", "unknown")
                task_counts[task] = task_counts.get(task, 0) + 1
            
            summary["success_rate"] = success_count / len(evaluations)
            summary["average_execution_time"] = total_time / len(evaluations)
            summary["task_breakdown"] = task_counts
            
            # Calculate performance trends (simplified)
            summary["performance_trends"] = {
                "execution_time_trend": "stable",  # Would calculate actual trend
                "success_rate_trend": "stable",    # Would calculate actual trend
                "quality_trend": "stable"          # Would calculate actual trend
            }
            
        except Exception as e:
            logger.error(f"Error calculating summary statistics: {e}")
        
        return summary
    
    async def export_evaluations(self, output_file: str, 
                               task: str = None, 
                               start_date: str = None, 
                               end_date: str = None):
        """Export evaluations to a file"""
        try:
            summary = await self.get_evaluation_summary(task, start_date, end_date)
            
            with open(output_file, "w") as f:
                json.dump(summary, f, indent=2)
            
            logger.info(f"Evaluations exported to {output_file}")
            
        except Exception as e:
            logger.error(f"Error exporting evaluations: {e}")
    
    def enable_evaluation(self):
        """Enable evaluation"""
        self.evaluation_enabled = True
        logger.info("Evaluation enabled")
    
    def disable_evaluation(self):
        """Disable evaluation"""
        self.evaluation_enabled = False
        logger.info("Evaluation disabled")
    
    def is_evaluation_enabled(self) -> bool:
        """Check if evaluation is enabled"""
        return self.evaluation_enabled 