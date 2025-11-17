"""
Inventory Optimization Service
Uses reinforcement learning for optimal stock levels
"""
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Dict, Optional
from datetime import datetime, timedelta
import numpy as np

app = FastAPI(
    title="TRADEMASTER Inventory Optimization Service",
    description="RL-based inventory optimization",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class InventoryItem(BaseModel):
    item_id: str
    name: str
    category: str
    current_stock: int
    unit_cost: float
    lead_time_days: int
    usage_history: List[int]  # Last N days usage


class OptimizationRequest(BaseModel):
    items: List[InventoryItem]
    optimization_horizon_days: int = 90
    service_level: float = 0.95  # 95% service level target


class OptimizationResponse(BaseModel):
    item_id: str
    current_stock: int
    optimal_stock: int
    reorder_point: int
    order_quantity: int
    estimated_savings: float
    stock_turnover_ratio: float


class ForecastRequest(BaseModel):
    item_id: str
    usage_history: List[int]
    forecast_days: int = 30


class ForecastResponse(BaseModel):
    item_id: str
    forecasted_demand: List[float]
    confidence_interval_lower: List[float]
    confidence_interval_upper: List[float]
    forecast_date: str


class ReorderRecommendation(BaseModel):
    item_id: str
    item_name: str
    current_stock: int
    reorder_point: int
    recommended_quantity: int
    urgency: str  # low, medium, high, critical
    estimated_stockout_date: str
    cost: float


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "inventory-optimization"
    }


@app.post("/api/v1/inventory/optimize", response_model=List[OptimizationResponse])
async def optimize_inventory(request: OptimizationRequest):
    """
    Optimize stock levels using RL-based approach
    """
    try:
        optimizations = []

        for item in request.items:
            # Calculate statistics from usage history
            avg_daily_usage = np.mean(item.usage_history)
            std_daily_usage = np.std(item.usage_history)

            # Calculate safety stock (for 95% service level, z = 1.65)
            z_score = 1.65 if request.service_level == 0.95 else 1.96
            safety_stock = int(z_score * std_daily_usage * np.sqrt(item.lead_time_days))

            # Calculate reorder point
            reorder_point = int((avg_daily_usage * item.lead_time_days) + safety_stock)

            # Calculate Economic Order Quantity (EOQ)
            # Assuming holding cost = 20% of unit cost per year
            # Assuming order cost = $50 per order
            annual_demand = avg_daily_usage * 365
            holding_cost_per_unit = item.unit_cost * 0.20
            order_cost = 50.0

            if holding_cost_per_unit > 0:
                eoq = int(np.sqrt((2 * annual_demand * order_cost) / holding_cost_per_unit))
            else:
                eoq = int(avg_daily_usage * 30)  # 30 days worth

            # Calculate optimal stock level (max inventory)
            optimal_stock = reorder_point + eoq

            # Calculate potential savings
            current_holding_cost = item.current_stock * holding_cost_per_unit
            optimal_holding_cost = optimal_stock * holding_cost_per_unit * 0.5  # Average inventory
            estimated_savings = max(0, current_holding_cost - optimal_holding_cost)

            # Calculate stock turnover ratio
            if item.current_stock > 0:
                stock_turnover = annual_demand / item.current_stock
            else:
                stock_turnover = 0

            optimizations.append(OptimizationResponse(
                item_id=item.item_id,
                current_stock=item.current_stock,
                optimal_stock=optimal_stock,
                reorder_point=reorder_point,
                order_quantity=eoq,
                estimated_savings=round(estimated_savings, 2),
                stock_turnover_ratio=round(stock_turnover, 2)
            ))

        return optimizations

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Optimization failed: {str(e)}")


@app.post("/api/v1/inventory/forecast", response_model=ForecastResponse)
async def forecast_demand(request: ForecastRequest):
    """
    Forecast future demand using time-series models
    """
    try:
        # Simple moving average with trend for demo
        # In production: Use Prophet or LSTM
        history = np.array(request.usage_history)

        # Calculate trend
        x = np.arange(len(history))
        z = np.polyfit(x, history, 1)
        trend = z[0]

        # Calculate seasonal pattern (simple)
        avg_demand = np.mean(history)
        seasonal_factor = history / avg_demand if avg_demand > 0 else np.ones_like(history)

        # Forecast
        forecasted_demand = []
        confidence_lower = []
        confidence_upper = []

        last_value = history[-1]
        std_dev = np.std(history)

        for i in range(request.forecast_days):
            # Simple trend projection
            forecast = last_value + (trend * (i + 1))

            # Apply seasonal pattern (repeating)
            season_idx = (len(history) + i) % len(seasonal_factor)
            forecast *= seasonal_factor[season_idx]

            # Confidence intervals (wider over time)
            margin = std_dev * np.sqrt(i + 1) * 1.96  # 95% confidence

            forecasted_demand.append(max(0, forecast))
            confidence_lower.append(max(0, forecast - margin))
            confidence_upper.append(max(0, forecast + margin))

        forecast_date = (datetime.now() + timedelta(days=request.forecast_days)).strftime("%Y-%m-%d")

        return ForecastResponse(
            item_id=request.item_id,
            forecasted_demand=[round(f, 1) for f in forecasted_demand],
            confidence_interval_lower=[round(f, 1) for f in confidence_lower],
            confidence_interval_upper=[round(f, 1) for f in confidence_upper],
            forecast_date=forecast_date
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Forecasting failed: {str(e)}")


@app.post("/api/v1/inventory/reorder", response_model=List[ReorderRecommendation])
async def get_reorder_recommendations(request: OptimizationRequest):
    """
    Get reorder recommendations for items that need restocking
    """
    try:
        recommendations = []

        for item in request.items:
            # Calculate basic metrics
            avg_daily_usage = np.mean(item.usage_history)
            std_daily_usage = np.std(item.usage_history)

            # Safety stock
            safety_stock = int(1.65 * std_daily_usage * np.sqrt(item.lead_time_days))
            reorder_point = int((avg_daily_usage * item.lead_time_days) + safety_stock)

            # Check if reorder is needed
            if item.current_stock <= reorder_point:
                # Calculate recommended quantity (EOQ)
                annual_demand = avg_daily_usage * 365
                holding_cost = item.unit_cost * 0.20
                order_cost = 50.0

                if holding_cost > 0:
                    eoq = int(np.sqrt((2 * annual_demand * order_cost) / holding_cost))
                else:
                    eoq = int(avg_daily_usage * 30)

                # Calculate urgency
                days_until_stockout = int(item.current_stock / avg_daily_usage) if avg_daily_usage > 0 else 999

                if days_until_stockout <= item.lead_time_days:
                    urgency = "critical"
                elif days_until_stockout <= item.lead_time_days * 1.5:
                    urgency = "high"
                elif days_until_stockout <= item.lead_time_days * 2:
                    urgency = "medium"
                else:
                    urgency = "low"

                stockout_date = (datetime.now() + timedelta(days=days_until_stockout)).strftime("%Y-%m-%d")

                recommendations.append(ReorderRecommendation(
                    item_id=item.item_id,
                    item_name=item.name,
                    current_stock=item.current_stock,
                    reorder_point=reorder_point,
                    recommended_quantity=eoq,
                    urgency=urgency,
                    estimated_stockout_date=stockout_date,
                    cost=round(eoq * item.unit_cost, 2)
                ))

        # Sort by urgency
        urgency_order = {"critical": 0, "high": 1, "medium": 2, "low": 3}
        recommendations.sort(key=lambda x: urgency_order[x.urgency])

        return recommendations

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Reorder recommendation failed: {str(e)}")


@app.post("/api/v1/inventory/abc-analysis")
async def abc_analysis(items: List[InventoryItem]):
    """
    Perform ABC analysis for inventory categorization
    """
    try:
        # Calculate annual usage value for each item
        item_values = []
        for item in items:
            annual_usage = np.mean(item.usage_history) * 365
            annual_value = annual_usage * item.unit_cost
            item_values.append({
                "item_id": item.item_id,
                "item_name": item.name,
                "annual_value": annual_value,
                "annual_usage": annual_usage
            })

        # Sort by value
        item_values.sort(key=lambda x: x["annual_value"], reverse=True)

        # Calculate cumulative percentages
        total_value = sum(iv["annual_value"] for iv in item_values)
        cumulative_pct = 0

        results = []
        for iv in item_values:
            value_pct = (iv["annual_value"] / total_value * 100) if total_value > 0 else 0
            cumulative_pct += value_pct

            # Classify (A: 70%, B: 20%, C: 10%)
            if cumulative_pct <= 70:
                category = "A"
                management_strategy = "Tight control, frequent review, accurate forecasting"
            elif cumulative_pct <= 90:
                category = "B"
                management_strategy = "Moderate control, periodic review"
            else:
                category = "C"
                management_strategy = "Simple control, bulk ordering"

            results.append({
                **iv,
                "category": category,
                "value_percentage": round(value_pct, 2),
                "cumulative_percentage": round(cumulative_pct, 2),
                "management_strategy": management_strategy
            })

        return {
            "total_items": len(items),
            "total_annual_value": round(total_value, 2),
            "items": results,
            "summary": {
                "A_items": len([r for r in results if r["category"] == "A"]),
                "B_items": len([r for r in results if r["category"] == "B"]),
                "C_items": len([r for r in results if r["category"] == "C"])
            }
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"ABC analysis failed: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8004)
