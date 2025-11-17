"""
Training Coach Service
AI-powered coaching using GPT-4 and computer vision
"""
import os
from fastapi import FastAPI, HTTPException, UploadFile, File
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Dict, Optional
import openai
from datetime import datetime

app = FastAPI(
    title="TRADEMASTER Training Coach Service",
    description="AI-powered vocational training coach",
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

# Initialize OpenAI
openai.api_key = os.getenv("OPENAI_API_KEY", "")


class ChatMessage(BaseModel):
    role: str  # user, assistant, system
    content: str


class ChatRequest(BaseModel):
    messages: List[ChatMessage]
    user_id: str
    context: Optional[Dict] = None  # User's current course, skill level, etc.


class ChatResponse(BaseModel):
    message: str
    suggestions: List[str]
    resources: List[Dict]
    follow_up_questions: List[str]


class MovementAnalysisRequest(BaseModel):
    video_url: str
    task_type: str  # welding, wiring, installation
    user_id: str


class MovementAnalysisResponse(BaseModel):
    score: float
    feedback: List[str]
    corrections: List[Dict]
    best_practices: List[str]


class FeedbackRequest(BaseModel):
    user_id: str
    course_id: str
    lesson_id: str
    performance_data: Dict


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "training-coach",
        "openai_configured": bool(openai.api_key)
    }


@app.post("/api/v1/coach/chat", response_model=ChatResponse)
async def chat_with_coach(request: ChatRequest):
    """
    Chat with AI coach powered by GPT-4
    """
    try:
        # Build system prompt based on context
        trade = request.context.get("trade", "electrical") if request.context else "electrical"
        skill_level = request.context.get("skill_level", "beginner") if request.context else "beginner"

        system_prompt = f"""You are an expert {trade} instructor with 20+ years of experience.
You're helping a {skill_level} level student learn {trade} trade skills.

Your teaching style:
- Patient and encouraging
- Use real-world examples
- Emphasize safety first
- Break down complex concepts into simple steps
- Provide practical, actionable advice

Always include:
1. Clear explanations
2. Safety considerations
3. Common mistakes to avoid
4. Pro tips from experience
"""

        # Prepare messages for OpenAI
        messages = [{"role": "system", "content": system_prompt}]
        messages.extend([{"role": msg.role, "content": msg.content} for msg in request.messages])

        # Call OpenAI API
        if openai.api_key:
            response = openai.ChatCompletion.create(
                model="gpt-4-turbo-preview",
                messages=messages,
                temperature=0.7,
                max_tokens=500
            )
            assistant_message = response.choices[0].message.content
        else:
            # Mock response for development
            assistant_message = f"""I understand you're learning about {trade}. As a {skill_level} level student,
let me break this down for you step by step.

Safety First: Always ensure you're wearing proper PPE before starting any work.

Here's what you need to know:
1. Start with understanding the basics
2. Practice the fundamentals
3. Build on your knowledge gradually

Would you like me to explain any specific aspect in more detail?"""

        # Generate suggestions
        suggestions = [
            f"Learn about {trade} safety protocols",
            "Watch video demonstration",
            "Practice with AR simulation",
            "Review certification requirements"
        ]

        # Generate resources
        resources = [
            {
                "title": f"{trade.title()} Safety Guidelines",
                "type": "document",
                "url": "/resources/safety-guidelines"
            },
            {
                "title": "Video Tutorial",
                "type": "video",
                "url": "/resources/video-tutorial"
            }
        ]

        # Generate follow-up questions
        follow_up_questions = [
            "What specific task are you trying to learn?",
            "Do you have experience with basic tools?",
            "Would you like to see a video demonstration?"
        ]

        return ChatResponse(
            message=assistant_message,
            suggestions=suggestions,
            resources=resources,
            follow_up_questions=follow_up_questions
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Chat failed: {str(e)}")


@app.post("/api/v1/coach/analyze-movement", response_model=MovementAnalysisResponse)
async def analyze_movement(request: MovementAnalysisRequest):
    """
    Analyze user's movement and technique using computer vision
    """
    try:
        # In production: Use computer vision model to analyze technique
        # For now, provide mock feedback

        task_type = request.task_type.lower()

        # Mock analysis results
        score = 0.75  # 75% accuracy

        feedback = []
        corrections = []
        best_practices = []

        if task_type == "welding":
            feedback = [
                "Good torch angle maintained",
                "Steady hand movement observed",
                "Travel speed slightly fast"
            ]
            corrections = [
                {
                    "issue": "Travel speed",
                    "current": "15 inches/min",
                    "recommended": "10-12 inches/min",
                    "impact": "Better penetration and bead appearance"
                }
            ]
            best_practices = [
                "Maintain 10-15 degree torch angle",
                "Keep consistent travel speed",
                "Watch the puddle, not the arc",
                "Practice smooth, steady motion"
            ]
        elif task_type == "wiring":
            feedback = [
                "Proper wire stripping length",
                "Good termination technique",
                "Color coding followed correctly"
            ]
            corrections = [
                {
                    "issue": "Wire bend radius",
                    "current": "Too tight",
                    "recommended": "Minimum 4x wire diameter",
                    "impact": "Prevents wire damage"
                }
            ]
            best_practices = [
                "Strip exactly 1/2 inch for most connections",
                "Twist stranded wires clockwise",
                "Leave no exposed copper",
                "Use proper torque for terminals"
            ]

        return MovementAnalysisResponse(
            score=score,
            feedback=feedback,
            corrections=corrections,
            best_practices=best_practices
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Movement analysis failed: {str(e)}")


@app.post("/api/v1/coach/voice-command")
async def process_voice_command(audio: UploadFile = File(...)):
    """
    Process voice commands using speech-to-text
    """
    try:
        # In production: Use Whisper API or similar for speech-to-text
        # For now, return mock response

        return {
            "transcription": "Show me how to wire a three-way switch",
            "intent": "tutorial_request",
            "entities": {
                "task": "wiring",
                "component": "three-way switch"
            },
            "response": "I'll show you how to wire a three-way switch. Let me pull up the AR tutorial...",
            "action": "launch_ar_tutorial",
            "action_params": {
                "tutorial_id": "three-way-switch-wiring"
            }
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Voice command processing failed: {str(e)}")


@app.post("/api/v1/coach/feedback")
async def get_personalized_feedback(request: FeedbackRequest):
    """
    Get personalized feedback based on performance data
    """
    try:
        performance_data = request.performance_data

        # Analyze performance
        quiz_score = performance_data.get("quiz_score", 0)
        completion_time = performance_data.get("completion_time", 0)
        attempts = performance_data.get("attempts", 1)

        # Generate personalized feedback
        feedback = {
            "overall_performance": "good" if quiz_score >= 80 else "needs_improvement",
            "strengths": [],
            "areas_for_improvement": [],
            "recommendations": [],
            "next_steps": []
        }

        if quiz_score >= 90:
            feedback["strengths"].append("Excellent understanding of concepts")
            feedback["next_steps"].append("Ready for advanced modules")
        elif quiz_score >= 70:
            feedback["strengths"].append("Good grasp of fundamentals")
            feedback["recommendations"].append("Review missed concepts")
        else:
            feedback["areas_for_improvement"].append("Core concepts need reinforcement")
            feedback["recommendations"].append("Retake lesson and practice more")

        if completion_time < performance_data.get("expected_time", 3600):
            feedback["strengths"].append("Efficient learning pace")
        else:
            feedback["recommendations"].append("Take your time to understand each concept")

        if attempts > 2:
            feedback["recommendations"].append("Consider one-on-one tutoring session")

        # Add specific recommendations
        feedback["next_steps"].extend([
            "Practice hands-on exercises",
            "Review safety procedures",
            "Take certification prep quiz"
        ])

        return feedback

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Feedback generation failed: {str(e)}")


@app.post("/api/v1/coach/generate-quiz")
async def generate_quiz(trade: str, topic: str, difficulty: str = "medium"):
    """
    Generate custom quiz questions using GPT-4
    """
    try:
        prompt = f"""Generate 5 {difficulty} difficulty quiz questions about {topic} in {trade} trade.

Format as JSON array with structure:
[
  {{
    "question": "Question text",
    "options": ["A", "B", "C", "D"],
    "correct_answer": 0,
    "explanation": "Why this is correct"
  }}
]

Focus on practical, real-world scenarios and safety considerations."""

        if openai.api_key:
            response = openai.ChatCompletion.create(
                model="gpt-4-turbo-preview",
                messages=[{"role": "user", "content": prompt}],
                temperature=0.7
            )
            quiz_json = response.choices[0].message.content
        else:
            # Mock quiz
            quiz_json = """[
                {
                    "question": "What is the proper wire gauge for a 20A circuit?",
                    "options": ["18 AWG", "16 AWG", "12 AWG", "10 AWG"],
                    "correct_answer": 2,
                    "explanation": "12 AWG wire is rated for 20A circuits according to NEC"
                }
            ]"""

        return {"quiz": quiz_json, "trade": trade, "topic": topic, "difficulty": difficulty}

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Quiz generation failed: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8003)
