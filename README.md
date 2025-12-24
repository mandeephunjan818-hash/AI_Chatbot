# AI Chatbot System Architecture  
**Date**: December 24, 2025  
**Version**: 1.0  

## 1. Overview  
An embeddable, multi-tenant AI chatbot that:  
- Understands user **intent** and **mood** using MINILMv2  
- Generates **empathetic, context-aware replies** via rule-based templates  
- Supports **multiple host applications** via iframe/widget  
- Ensures **data isolation**, **low latency**, and **high accuracy**

## 2. Tech Stack

| Layer | Technology | Purpose |
|------|-----------|--------|
| **Frontend** | React (SPA) | Chat widget UI, embeddable via `<iframe>` |
| **API Gateway** | Go | Auth, rate limiting, orchestration |
| **NLU Engine** | Python + MINILMv2 (6×384) → ONNX | Intent + mood classification |
| **Cache / Session** | Redis | Rate limiting, temporary tokens |
| **Database** | MongoDB | Conversation history, metadata |
| **Hosting** | Vercel (Frontend), Kubernetes (Backend) | Scalable, secure deployment |

## 3. Data Flow

1. **User** types message in embedded widget (`app_id=APP123`, `user_token=usr_abc`)
2. **React SPA** sends to `POST /api/chat`
3. **Go API**:
   - Validates `app_id` + `user_token`
   - Checks rate limit (Redis)
   - Generates `session_id = hash(app_id + user_token)`
4. **Go → Python** (internal): sends `{"text": "...", "session_id": "..."}`  
5. **Python (MINILMv2 ONNX)**:
   - Runs inference → returns `{"intent": "...", "mood": "...", "confidences": ...}`
6. **Go**:
   - Matches `(intent, mood)` to **reply template**
   - Generates empathetic response
   - Logs to **MongoDB** with full context
7. **Response** sent to frontend → rendered in chat bubble

## 4. Model Strategy

### 4.1. Understanding (MINILMv2)
- **Base**: `MINILMv2 6×384` (Microsoft)
- **Fine-tuned on**:
  - **GoEmotions** (27 emotions)
  - **Custom intent labels** (e.g., billing, login, greeting)
- **Output**: Structured `(intent, mood)` with confidence scores
- **Optimized**: Converted to **ONNX** for 2–3× faster CPU inference

### 4.2. Response Generation
- **Method**: **Template-based** (YAML/JSON rules)
- **Input**: `(intent, mood)`
- **Example**:
  ```yaml
  - intent: billing_complaint
    mood: angry
    reply: "I'm truly sorry you're upset. Let me fix this right away."
    escalate: true
