"""
Computer Vision Service for Blueprint Analysis
Uses YOLO v8 for object detection and OCR for text extraction
"""
import time
from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Dict, Optional
import cv2
import numpy as np
from ultralytics import YOLO
import pytesseract
from PIL import Image
import io

app = FastAPI(
    title="TRADEMASTER CV Service",
    description="Computer Vision for Blueprint Analysis",
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

# Load YOLO model (pre-trained or custom)
# For production, replace with custom-trained model
try:
    model = YOLO("yolov8n.pt")  # Nano model for speed
    print("YOLO model loaded successfully")
except Exception as e:
    print(f"Warning: Could not load YOLO model: {e}")
    model = None


class DetectedItem(BaseModel):
    type: str
    label: str
    confidence: float
    bounding_box: Dict[str, float]
    quantity: int = 1


class BlueprintAnalysis(BaseModel):
    room_count: int
    total_area: float
    dimensions: Dict[str, float]
    detected_items: List[DetectedItem]
    extracted_text: List[str]
    confidence: float
    processing_time: float


class BlueprintAnalysisRequest(BaseModel):
    file_url: str


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "cv-blueprint-analysis",
        "model_loaded": model is not None
    }


@app.post("/api/v1/cv/analyze-blueprint", response_model=BlueprintAnalysis)
async def analyze_blueprint(file: UploadFile = File(...)):
    """
    Analyze blueprint image and extract components
    """
    start_time = time.time()

    try:
        # Read image
        contents = await file.read()
        nparr = np.frombuffer(contents, np.uint8)
        img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

        if img is None:
            raise HTTPException(status_code=400, detail="Invalid image file")

        # Detect objects using YOLO
        detected_items = []
        if model:
            results = model(img)
            for result in results:
                boxes = result.boxes
                for box in boxes:
                    x1, y1, x2, y2 = box.xyxy[0].tolist()
                    conf = float(box.conf[0])
                    cls = int(box.cls[0])
                    label = model.names[cls]

                    detected_items.append(DetectedItem(
                        type="component",
                        label=label,
                        confidence=conf,
                        bounding_box={
                            "x1": x1,
                            "y1": y1,
                            "x2": x2,
                            "y2": y2
                        }
                    ))

        # Extract text using OCR
        pil_img = Image.open(io.BytesIO(contents))
        extracted_text = pytesseract.image_to_string(pil_img).split('\n')
        extracted_text = [text.strip() for text in extracted_text if text.strip()]

        # Calculate dimensions (example logic)
        height, width = img.shape[:2]
        dimensions = {
            "width_pixels": float(width),
            "height_pixels": float(height),
            "estimated_area_sqft": float(width * height) / 10000  # Example conversion
        }

        # Estimate room count based on detected walls/spaces
        room_count = max(1, len(detected_items) // 10)  # Simplified logic

        # Calculate overall confidence
        avg_confidence = np.mean([item.confidence for item in detected_items]) if detected_items else 0.5

        processing_time = time.time() - start_time

        return BlueprintAnalysis(
            room_count=room_count,
            total_area=dimensions["estimated_area_sqft"],
            dimensions=dimensions,
            detected_items=detected_items,
            extracted_text=extracted_text,
            confidence=float(avg_confidence),
            processing_time=processing_time
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Analysis failed: {str(e)}")


@app.post("/api/v1/cv/extract-text")
async def extract_text(file: UploadFile = File(...)):
    """
    Extract text and dimensions from blueprint using OCR
    """
    try:
        contents = await file.read()
        pil_img = Image.open(io.BytesIO(contents))

        # Extract text
        text = pytesseract.image_to_string(pil_img)

        # Extract dimensions (look for patterns like "10'x12'")
        import re
        dimension_pattern = r"(\d+(?:\.\d+)?)\s*['\"]?\s*[xX×]\s*(\d+(?:\.\d+)?)\s*['\"]?"
        dimensions = re.findall(dimension_pattern, text)

        return {
            "text": text,
            "dimensions": dimensions,
            "line_count": len(text.split('\n'))
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Text extraction failed: {str(e)}")


@app.post("/api/v1/cv/detect-components")
async def detect_components(file: UploadFile = File(...), component_type: str = "electrical"):
    """
    Detect specific components (electrical, HVAC, etc.)
    """
    try:
        if not model:
            return {"detected_components": [], "message": "Model not loaded"}

        contents = await file.read()
        nparr = np.frombuffer(contents, np.uint8)
        img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

        # Run detection
        results = model(img)
        components = []

        for result in results:
            boxes = result.boxes
            for box in boxes:
                x1, y1, x2, y2 = box.xyxy[0].tolist()
                conf = float(box.conf[0])
                cls = int(box.cls[0])
                label = model.names[cls]

                # Filter by component type (simplified)
                if component_type.lower() in label.lower() or conf > 0.5:
                    components.append({
                        "label": label,
                        "confidence": conf,
                        "position": {
                            "x": (x1 + x2) / 2,
                            "y": (y1 + y2) / 2
                        },
                        "size": {
                            "width": x2 - x1,
                            "height": y2 - y1
                        }
                    })

        return {
            "component_type": component_type,
            "detected_components": components,
            "count": len(components)
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Component detection failed: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
