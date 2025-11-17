# TRADEMASTER ML Services

Python-based machine learning and AI services for TRADEMASTER platform.

## Services

1. **cv-blueprint-analysis**: Computer vision for blueprint analysis (YOLO v8)
2. **cost-prediction**: Cost and duration prediction (LSTM, XGBoost)
3. **training-coach**: AI coaching with GPT-4
4. **inventory-optimization**: Inventory optimization with reinforcement learning

## Tech Stack

- Python 3.11+
- FastAPI (REST API)
- TensorFlow 2.x (LSTM models)
- PyTorch (YOLO models)
- OpenAI API (GPT-4)
- Scikit-learn (preprocessing)
- Pandas, NumPy (data processing)

## Development

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Run all services
python main.py

# Run specific service
cd cv-blueprint-analysis && python app.py

# Run tests
pytest

# Train models
python scripts/train_models.py
```

## Project Structure

```
ml-services/
├── cv-blueprint-analysis/
│   ├── app.py               # FastAPI app
│   ├── models/              # YOLO models
│   ├── utils/               # Image processing
│   └── tests/
├── cost-prediction/
│   ├── app.py
│   ├── models/              # LSTM, XGBoost
│   ├── training/            # Training scripts
│   └── tests/
├── training-coach/
│   ├── app.py
│   ├── prompts/             # GPT-4 prompts
│   └── tests/
├── inventory-optimization/
│   ├── app.py
│   ├── agents/              # RL agents
│   └── tests/
├── shared/
│   ├── config.py
│   └── utils.py
└── main.py                  # Main entry point
```

## API Endpoints

### Computer Vision Service (Port 8001)
- `POST /api/v1/cv/analyze-blueprint` - Analyze blueprint
- `POST /api/v1/cv/extract-text` - Extract text/dimensions
- `POST /api/v1/cv/detect-components` - Detect electrical components

### Cost Prediction Service (Port 8002)
- `POST /api/v1/predict/material-costs` - Predict material costs
- `POST /api/v1/predict/duration` - Predict project duration
- `POST /api/v1/predict/risk` - Assess risk
- `POST /api/v1/anomaly/detect` - Detect cost anomalies

### Training Coach Service (Port 8003)
- `POST /api/v1/coach/chat` - Chat with AI coach
- `POST /api/v1/coach/analyze-movement` - Analyze technique
- `POST /api/v1/coach/voice-command` - Process voice command
- `POST /api/v1/coach/feedback` - Get personalized feedback

### Inventory Optimization Service (Port 8004)
- `POST /api/v1/inventory/optimize` - Optimize stock levels
- `POST /api/v1/inventory/forecast` - Forecast demand
- `POST /api/v1/inventory/reorder` - Get reorder recommendations

## Model Training

```bash
# Train cost prediction model
python cost-prediction/training/train_lstm.py

# Train blueprint detection model
python cv-blueprint-analysis/training/train_yolo.py

# Evaluate models
python scripts/evaluate_models.py
```

## Model Deployment

Models are served using:
- TensorFlow Serving (LSTM models)
- TorchServe (YOLO models)
- Direct FastAPI integration (smaller models)

## Performance

- Blueprint analysis: < 2 seconds
- Cost prediction: < 500ms
- AI coaching response: < 1 second
- Inventory optimization: < 1 second
