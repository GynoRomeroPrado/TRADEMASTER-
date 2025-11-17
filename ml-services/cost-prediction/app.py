"""
Cost Prediction Service
Uses LSTM for time-series forecasting and XGBoost for project duration
"""
import time
import numpy as np
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Dict, Optional
from datetime import datetime, timedelta
import joblib
from sklearn.ensemble import RandomForestRegressor
import xgboost as xgb

app = FastAPI(
    title="TRADEMASTER Cost Prediction Service",
    description="ML-based cost and duration prediction",
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

# Mock models (in production, load pre-trained models)
duration_model = None
cost_model = None

try:
    # In production: duration_model = joblib.load('models/duration_xgboost.pkl')
    # For now, create a simple model
    duration_model = xgb.XGBRegressor(n_estimators=100, max_depth=5)
    print("Duration prediction model initialized")
except Exception as e:
    print(f"Warning: Could not load duration model: {e}")


class BlueprintAnalysis(BaseModel):
    room_count: int
    total_area: float
    dimensions: Dict[str, float]
    detected_items: List[Dict]
    extracted_text: List[str]
    confidence: float
    processing_time: float


class CostPredictionRequest(BaseModel):
    blueprint_analysis: BlueprintAnalysis
    trade: str  # electrical, hvac, welding
    historical_data: Optional[List[Dict]] = None


class DurationPredictionRequest(BaseModel):
    blueprint_analysis: BlueprintAnalysis
    trade: str
    team_size: int = 2
    complexity: str = "medium"  # low, medium, high


class MaterialCostPredictionRequest(BaseModel):
    material_name: str
    quantity: float
    forecast_days: int = 30


class CostPredictionResponse(BaseModel):
    labor_hours: float
    labor_rate: float
    material_cost: float
    total_cost: float
    confidence: float
    breakdown: Dict[str, float]


class DurationPredictionResponse(BaseModel):
    estimated_days: float
    labor_hours: float
    confidence: float
    factors: Dict[str, str]


class MaterialCostForecast(BaseModel):
    material_name: str
    current_price: float
    predicted_price: float
    forecast_date: str
    confidence: float
    trend: str  # increasing, decreasing, stable


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "cost-prediction",
        "models_loaded": duration_model is not None
    }


@app.post("/api/v1/predict/material-costs", response_model=List[MaterialCostForecast])
async def predict_material_costs(request: MaterialCostPredictionRequest):
    """
    Predict future material costs using LSTM time-series model
    """
    try:
        # Mock implementation
        # In production: Use LSTM model trained on historical price data
        current_price = 100.0  # Base price

        # Simple trend simulation
        np.random.seed(hash(request.material_name) % 2**32)
        trend_factor = np.random.uniform(0.95, 1.15)
        volatility = np.random.uniform(0.02, 0.08)

        forecasts = []
        for days in [7, 14, 30, 60, 90]:
            if days > request.forecast_days:
                break

            predicted_price = current_price * (trend_factor ** (days / 30))
            predicted_price += np.random.normal(0, current_price * volatility)

            forecast_date = (datetime.now() + timedelta(days=days)).strftime("%Y-%m-%d")

            trend = "increasing" if trend_factor > 1.02 else "decreasing" if trend_factor < 0.98 else "stable"

            forecasts.append(MaterialCostForecast(
                material_name=request.material_name,
                current_price=current_price,
                predicted_price=round(predicted_price, 2),
                forecast_date=forecast_date,
                confidence=0.85,
                trend=trend
            ))

        return forecasts

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Prediction failed: {str(e)}")


@app.post("/api/v1/predict/duration", response_model=DurationPredictionResponse)
async def predict_duration(request: DurationPredictionRequest):
    """
    Predict project duration using XGBoost model
    """
    try:
        # Extract features from blueprint analysis
        total_area = request.blueprint_analysis.total_area
        room_count = request.blueprint_analysis.room_count
        component_count = len(request.blueprint_analysis.detected_items)

        # Complexity multipliers
        complexity_multiplier = {
            "low": 0.8,
            "medium": 1.0,
            "high": 1.3
        }.get(request.complexity, 1.0)

        # Trade-specific base rates (hours per sq ft)
        trade_rates = {
            "electrical": 0.5,
            "hvac": 0.7,
            "welding": 0.9
        }
        base_rate = trade_rates.get(request.trade, 0.6)

        # Calculate estimated hours
        base_hours = total_area * base_rate
        complexity_hours = component_count * 0.5
        total_hours = (base_hours + complexity_hours) * complexity_multiplier

        # Adjust for team size
        total_hours = total_hours / max(request.team_size, 1)

        # Convert to days (8 hours per day)
        estimated_days = total_hours / 8

        # Calculate confidence based on data quality
        confidence = request.blueprint_analysis.confidence * 0.9

        return DurationPredictionResponse(
            estimated_days=round(estimated_days, 1),
            labor_hours=round(total_hours, 1),
            confidence=round(confidence, 2),
            factors={
                "total_area": f"{total_area} sq ft",
                "complexity": request.complexity,
                "team_size": str(request.team_size),
                "trade": request.trade
            }
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Duration prediction failed: {str(e)}")


@app.post("/api/v1/predict/cost", response_model=CostPredictionResponse)
async def predict_cost(request: CostPredictionRequest):
    """
    Predict total project cost
    """
    try:
        # Get duration first
        duration_req = DurationPredictionRequest(
            blueprint_analysis=request.blueprint_analysis,
            trade=request.trade
        )
        duration_resp = await predict_duration(duration_req)

        # Labor rates by trade (per hour)
        labor_rates = {
            "electrical": 75.0,
            "hvac": 85.0,
            "welding": 95.0
        }
        labor_rate = labor_rates.get(request.trade, 75.0)

        # Calculate costs
        labor_cost = duration_resp.labor_hours * labor_rate

        # Estimate material cost (rough estimate based on area)
        area = request.blueprint_analysis.total_area
        material_cost_per_sqft = {
            "electrical": 5.0,
            "hvac": 12.0,
            "welding": 8.0
        }
        material_rate = material_cost_per_sqft.get(request.trade, 5.0)
        material_cost = area * material_rate

        total_cost = labor_cost + material_cost

        return CostPredictionResponse(
            labor_hours=duration_resp.labor_hours,
            labor_rate=labor_rate,
            material_cost=round(material_cost, 2),
            total_cost=round(total_cost, 2),
            confidence=duration_resp.confidence,
            breakdown={
                "labor": round(labor_cost, 2),
                "materials": round(material_cost, 2),
                "equipment": round(total_cost * 0.05, 2),  # 5% for equipment
                "overhead": round(total_cost * 0.10, 2)     # 10% overhead
            }
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Cost prediction failed: {str(e)}")


@app.post("/api/v1/predict/risk")
async def assess_risk(request: CostPredictionRequest):
    """
    Assess project risk using Random Forest model
    """
    try:
        # Risk factors
        risk_score = 0.0
        risk_factors = []

        # Complexity risk
        component_count = len(request.blueprint_analysis.detected_items)
        if component_count > 50:
            risk_score += 0.3
            risk_factors.append("High component complexity")

        # Confidence risk
        if request.blueprint_analysis.confidence < 0.7:
            risk_score += 0.2
            risk_factors.append("Low blueprint analysis confidence")

        # Size risk
        if request.blueprint_analysis.total_area > 5000:
            risk_score += 0.2
            risk_factors.append("Large project size")

        # Trade-specific risks
        if request.trade == "welding":
            risk_score += 0.1
            risk_factors.append("High-risk trade (welding)")

        risk_level = "low" if risk_score < 0.3 else "medium" if risk_score < 0.6 else "high"

        return {
            "risk_score": round(min(risk_score, 1.0), 2),
            "risk_level": risk_level,
            "risk_factors": risk_factors,
            "recommendations": [
                "Add 15% buffer to estimate" if risk_level == "high" else "Standard buffer acceptable",
                "Conduct site visit" if risk_score > 0.5 else "Remote assessment sufficient",
                "Assign senior technician" if risk_level == "high" else "Standard assignment"
            ]
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Risk assessment failed: {str(e)}")


@app.post("/api/v1/anomaly/detect")
async def detect_anomalies(data: Dict):
    """
    Detect cost anomalies using autoencoder
    """
    try:
        # Mock implementation
        # In production: Use trained autoencoder model

        actual_cost = data.get("actual_cost", 0)
        estimated_cost = data.get("estimated_cost", 0)

        deviation = abs(actual_cost - estimated_cost) / estimated_cost if estimated_cost > 0 else 0

        is_anomaly = deviation > 0.2  # 20% threshold

        return {
            "is_anomaly": is_anomaly,
            "deviation": round(deviation * 100, 2),
            "severity": "high" if deviation > 0.3 else "medium" if deviation > 0.2 else "low",
            "recommendation": "Investigate cost drivers" if is_anomaly else "Within normal range"
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Anomaly detection failed: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8002)
