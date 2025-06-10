# ACI Integration Guide

This document provides an overview of how Casibase's server components integrate with ACI compatible clients. It focuses on the interaction between the `ModelProvider` abstraction and the agent modules, shows a recommended tag structure for personalizing chat requests, and outlines a simple flow for React applications that call the ACI APIs.

## Interaction between ModelProvider and agent modules

Casibase wraps access to large language models through the `ModelProvider` interface defined in `model/provider.go`. Each provider implementation handles a specific backend (OpenAI, local models, etc.). The agent system, found under `agent/`, exposes tools that can be called during a conversation. These tools are discovered and managed through an `AgentProvider` which returns an `AgentClients` structure containing the available MCP tool clients.

When a request is processed in `controllers/message_answer.go` the backend:

1. Resolves the desired `ModelProvider` based on the configured store settings.
2. Builds an `AgentInfo` struct containing agent clients and a buffer for `AgentMessages`.
3. Invokes `model.QueryTextWithTools()` which passes the `AgentInfo` to the model provider.
4. The model provider may generate tool calls that are executed via the agent clients. Responses from those tools are appended to the message stream and fed back into the model until no further tool calls are produced.

This loop allows models to extend their capabilities by using external agents while keeping a clean separation between language model logic and tool execution.

## Recommended tag structure

Tags are used across Casibase objects (nodes, machines, videos, etc.) and can also link to external resources via the `TagLink` model. For personalized deployments it is helpful to follow a consistent naming scheme. A common approach is:

- `org:<name>` – organizational grouping or tenant ID.
- `role:<name>` – user role such as `admin`, `student`, or `guest`.
- `topic:<keyword>` – domain specific categories like `security`, `finance` or `demo`.

Combining these tags lets you filter data and serve different prompts to distinct audiences. For example, a tag string `org:acme,role:student,topic:finance` can be attached to messages or assets to target finance training materials for students of the Acme organization.

## Example React flow using ACI APIs

Below is a simplified sequence for a React client integrating with ACI:

1. **Authenticate** – obtain an API token from your Casibase deployment.
2. **Create a chat** – POST to `/api/chat` (or your configured endpoint) with the desired model and personalization tags.
3. **Stream messages** – send user questions to `/api/chat/:id/message` using Server‑Sent Events (SSE) to receive incremental updates.
4. **Handle tool calls** – when tool call events are emitted, display intermediate status while the server executes agent tools.
5. **Render final answer** – once the stream ends, show the full assistant response along with any cited tool outputs.

This flow is typically wrapped in a custom hook so React components only need to provide the input text and react to streaming results.

---

With the above concepts you can integrate ACI compatible clients with Casibase and extend conversations with personalized tags and external agent tools.
