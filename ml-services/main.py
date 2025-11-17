"""
Main entry point for all ML services
Runs all services concurrently using multiprocessing
"""
import multiprocessing
import uvicorn
import sys
from pathlib import Path


def run_cv_service():
    """Run Computer Vision service on port 8001"""
    sys.path.insert(0, str(Path(__file__).parent / "cv-blueprint-analysis"))
    from cv_blueprint_analysis.app import app
    uvicorn.run(app, host="0.0.0.0", port=8001)


def run_cost_prediction_service():
    """Run Cost Prediction service on port 8002"""
    sys.path.insert(0, str(Path(__file__).parent / "cost-prediction"))
    from cost_prediction.app import app
    uvicorn.run(app, host="0.0.0.0", port=8002)


def run_training_coach_service():
    """Run Training Coach service on port 8003"""
    sys.path.insert(0, str(Path(__file__).parent / "training-coach"))
    from training_coach.app import app
    uvicorn.run(app, host="0.0.0.0", port=8003)


def run_inventory_optimization_service():
    """Run Inventory Optimization service on port 8004"""
    sys.path.insert(0, str(Path(__file__).parent / "inventory-optimization"))
    from inventory_optimization.app import app
    uvicorn.run(app, host="0.0.0.0", port=8004)


if __name__ == "__main__":
    print("Starting TRADEMASTER ML Services...")
    print("=" * 60)
    print("Computer Vision Service:      http://localhost:8001")
    print("Cost Prediction Service:      http://localhost:8002")
    print("Training Coach Service:       http://localhost:8003")
    print("Inventory Optimization:       http://localhost:8004")
    print("=" * 60)

    # Create processes for each service
    processes = [
        multiprocessing.Process(target=run_cv_service, name="CV-Service"),
        multiprocessing.Process(target=run_cost_prediction_service, name="Cost-Prediction"),
        multiprocessing.Process(target=run_training_coach_service, name="Training-Coach"),
        multiprocessing.Process(target=run_inventory_optimization_service, name="Inventory-Optimization"),
    ]

    # Start all processes
    for p in processes:
        p.start()

    # Wait for all processes
    try:
        for p in processes:
            p.join()
    except KeyboardInterrupt:
        print("\nShutting down ML services...")
        for p in processes:
            p.terminate()
            p.join()
        print("All services stopped.")
