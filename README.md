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

┌─────────────────────────────────────────────────────────────┐
│                    CLIENT APPLICATIONS                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌──────────────────────────────────┐   │
│  │   Host      │    │     Dashboard SPA                │   │
│  │  Website    │◄───►  (React + Vercel)               │   │
│  │             │    │  • Login/Signup                  │   │
│  │  ┌──────────┴─┐  │  • Widget Management            │   │
│  │  │  Chat      │  │  • Key Generation               │   │
│  │  │  Widget    │  │  • Placement Settings           │   │
│  │  │  (iframe)  │  │                                  │   │
│  │  └────────────┘  └──────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                   API GATEWAY (Go)                          │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ • JWT Authentication                                 │  │
│  │ • Rate Limiting                                      │  │
│  │ • Request Routing                                    │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────────┐
│   AUTH SERVICE  │ │  WIDGET SERVICE │ │    CHAT SERVICE     │
│   (Go)          │ │   (Go)          │ │     (Go + Python)   │
│ • Login/Reg     │ │ • Widget CRUD   │ │ • MINILMv2 Inference│
│ • JWT Issuance  │ │ • Key Generation│ │ • Template Matching │
│ • Session Mgmt  │ │ • Placement     │ │ • Conversation Flow │
└─────────────────┘ └─────────────────┘ └─────────────────────┘
         │                     │                     │
         └─────────────────────┼─────────────────────┘
                               │
                     ┌─────────┴─────────┐
                     ▼         ▼         ▼
           ┌─────────────┐ ┌──────┐ ┌──────────┐
           │   MongoDB   │ │ Redis│ │   ONNX   │
           │ Collections │ │ Cache│ │  Models  │
           │ • users     │ │ • JWT│ │ • MINILM │
           │ • widgets   │ │ • RL │ │          │
           │ • convos    │ │ • ...│ └──────────┘
           └─────────────┘ └──────┘

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
