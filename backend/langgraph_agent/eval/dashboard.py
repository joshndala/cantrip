#!/usr/bin/env python3
"""
Simple Phoenix Evaluation Dashboard
View evaluation logs in a web interface
"""

import json
import os
from datetime import datetime
from flask import Flask, render_template_string, jsonify, request
import glob

app = Flask(__name__)

# HTML template for the dashboard
DASHBOARD_HTML = """
<!DOCTYPE html>
<html>
<head>
    <title>CanTrip Phoenix Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; }
        .header { background: #2c3e50; color: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: #ecf0f1; padding: 20px; border-radius: 8px; text-align: center; }
        .stat-number { font-size: 2em; font-weight: bold; color: #2c3e50; }
        .stat-label { color: #7f8c8d; margin-top: 5px; }
        .evaluations { margin-top: 20px; }
        .evaluation-item { background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #3498db; }
        .evaluation-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
        .task-badge { background: #3498db; color: white; padding: 5px 10px; border-radius: 15px; font-size: 0.8em; }
        .success { border-left-color: #27ae60; }
        .error { border-left-color: #e74c3c; }
        .metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 10px; margin-top: 10px; }
        .metric { background: white; padding: 10px; border-radius: 5px; text-align: center; }
        .metric-value { font-weight: bold; color: #2c3e50; }
        .metric-label { font-size: 0.8em; color: #7f8c8d; }
        .refresh-btn { background: #3498db; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; }
        .refresh-btn:hover { background: #2980b9; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 CanTrip Phoenix Dashboard</h1>
            <p>AI Agent Performance Monitoring</p>
        </div>
        
        <button class="refresh-btn" onclick="location.reload()">🔄 Refresh Data</button>
        
        <div class="stats">
            <div class="stat-card">
                <div class="stat-number">{{ stats.total_evaluations }}</div>
                <div class="stat-label">Total Evaluations</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{ "%.1f"|format(stats.success_rate * 100) }}%</div>
                <div class="stat-label">Success Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{ "%.2f"|format(stats.avg_execution_time) }}s</div>
                <div class="stat-label">Avg Execution Time</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{ stats.tasks_count }}</div>
                <div class="stat-label">Task Types</div>
            </div>
        </div>
        
        <div class="evaluations">
            <h2>📊 Recent Evaluations</h2>
            {% for eval in evaluations %}
            <div class="evaluation-item {{ 'success' if eval.metrics.success else 'error' }}">
                <div class="evaluation-header">
                    <span class="task-badge">{{ eval.task }}</span>
                    <span>{{ eval.timestamp }}</span>
                </div>
                <div class="metrics">
                    <div class="metric">
                        <div class="metric-value">{{ "%.2f"|format(eval.metrics.execution_time_seconds) }}s</div>
                        <div class="metric-label">Execution Time</div>
                    </div>
                    <div class="metric">
                        <div class="metric-value">{{ "✅" if eval.metrics.success else "❌" }}</div>
                        <div class="metric-label">Success</div>
                    </div>
                    {% if eval.metrics.itinerary_quality is defined %}
                    <div class="metric">
                        <div class="metric-value">{{ "%.2f"|format(eval.metrics.itinerary_quality) }}</div>
                        <div class="metric-label">Quality Score</div>
                    </div>
                    {% endif %}
                    {% if eval.metrics.suggestion_count is defined %}
                    <div class="metric">
                        <div class="metric-value">{{ eval.metrics.suggestion_count }}</div>
                        <div class="metric-label">Suggestions</div>
                    </div>
                    {% endif %}
                    {% if eval.metrics.item_count is defined %}
                    <div class="metric">
                        <div class="metric-value">{{ eval.metrics.item_count }}</div>
                        <div class="metric-label">Items</div>
                    </div>
                    {% endif %}
                </div>
            </div>
            {% endfor %}
        </div>
    </div>
</body>
</html>
"""

def load_evaluations():
    """Load evaluation data from JSONL files"""
    evaluations = []
    log_files = glob.glob("*.jsonl") + glob.glob("../*.jsonl")
    
    for log_file in log_files:
        if os.path.exists(log_file):
            with open(log_file, 'r') as f:
                for line in f:
                    try:
                        evaluation = json.loads(line.strip())
                        evaluations.append(evaluation)
                    except json.JSONDecodeError:
                        continue
    
    return sorted(evaluations, key=lambda x: x.get('timestamp', ''), reverse=True)

def calculate_stats(evaluations):
    """Calculate dashboard statistics"""
    if not evaluations:
        return {
            'total_evaluations': 0,
            'success_rate': 0.0,
            'avg_execution_time': 0.0,
            'tasks_count': 0
        }
    
    total = len(evaluations)
    success_count = sum(1 for e in evaluations if e.get('metrics', {}).get('success', False))
    total_time = sum(e.get('metrics', {}).get('execution_time_seconds', 0) for e in evaluations)
    tasks = set(e.get('task', 'unknown') for e in evaluations)
    
    return {
        'total_evaluations': total,
        'success_rate': success_count / total if total > 0 else 0.0,
        'avg_execution_time': total_time / total if total > 0 else 0.0,
        'tasks_count': len(tasks)
    }

@app.route('/')
def dashboard():
    """Main dashboard page"""
    evaluations = load_evaluations()
    stats = calculate_stats(evaluations)
    
    return render_template_string(DASHBOARD_HTML, evaluations=evaluations[:50], stats=stats)

@app.route('/api/evaluations')
def api_evaluations():
    """API endpoint for evaluations"""
    evaluations = load_evaluations()
    return jsonify(evaluations)

@app.route('/api/stats')
def api_stats():
    """API endpoint for statistics"""
    evaluations = load_evaluations()
    stats = calculate_stats(evaluations)
    return jsonify(stats)

if __name__ == '__main__':
    print("🚀 Starting CanTrip Phoenix Dashboard...")
    print("📊 Dashboard available at: http://localhost:5000")
    print("📈 API endpoints:")
    print("   - http://localhost:5000/api/evaluations")
    print("   - http://localhost:5000/api/stats")
    
    app.run(host='0.0.0.0', port=5000, debug=True) 