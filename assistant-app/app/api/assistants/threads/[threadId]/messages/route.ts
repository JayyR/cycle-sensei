import { assistantId } from "@/app/assistant-config";
import { openai } from "@/app/openai";

export const runtime = "nodejs";

// Send a new message to a thread
export async function POST(request, { params: { threadId } }) {
  const { content } = await request.json();

  await openai.beta.threads.messages.create(threadId, {
    role: "user",
    content: content,
  });

  const stream = openai.beta.threads.runs.stream(threadId, {
    assistant_id: assistantId,
  });

  return new Response(stream.toReadableStream());
}

// Get messages from a thread
export async function GET(request, { params: { threadId } }) {
  const messages = await openai.beta.threads.messages.list(threadId, {
    order: "asc",
    limit: 5,
  });

  return Response.json(messages);
}