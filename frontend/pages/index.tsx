import * as React from "react";
import { useState, useEffect, useRef } from "react";
import { Container, CssBaseline, Typography, Box } from "@mui/material";
import TopBar from "../components/TopBar";
import ChatInput from "../components/ChatInput";
import ResponseDisplay from "../components/ResponseDisplay";
import FloatingInputBox from "../components/FloatingInputBox";

interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

export default function Home() {
  const [chatHistory, setChatHistory] = useState<ChatMessage[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    ws.current = new WebSocket("ws://localhost:8080/ws");

    ws.current.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "chat_response") {
          setChatHistory((prev) => [...prev, { role: "assistant", content: data.message }]);
          setLoading(false);
        } else if (data.type === "typing") {
          setLoading(true);
        } else if (data.type === "error") {
          setLoading(false);
          setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style=\"color:red\">Error: ${data.error}</p>` }]);
        }
      } catch (e) {
        setLoading(false);
        setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style=\"color:red\">Invalid response from server</p>` }]);
      }
    };

    ws.current.onerror = (event) => {
      setLoading(false);
      setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style=\"color:red\">WebSocket error occurred</p>` }]);
    };

    ws.current.onclose = () => {
      setLoading(false);
      setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style=\"color:gray\">Disconnected from server</p>` }]);
    };

    return () => {
      ws.current?.close();
    };
  }, []);

  const sendMessage = (msg: string) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      setChatHistory((prev) => [...prev, { role: "user", content: msg }]);
      setLoading(true);
      ws.current.send(JSON.stringify({ type: "chat", message: msg }));
    }
  };

  return (
    <>
      <CssBaseline />
      <TopBar />
      <Container maxWidth="md" sx={{ mt: 4, mb: 10, minHeight: '70vh', display: 'flex', flexDirection: 'column' }}>
        <Box sx={{ flexGrow: 1, overflowY: 'auto', mb: 2, fontFamily: 'Georgia, serif', fontSize: '1.1rem', lineHeight: 1.6, border: '1px solid #ddd', borderRadius: 2, p: 2 }}>
          {chatHistory.map((msg, idx) => (
            <Box key={idx} sx={{ mb: 2, color: msg.role === 'user' ? 'primary.main' : 'text.primary' }}>
              <strong>{msg.role === 'user' ? 'You:' : 'TechDocs AI:'}</strong>
              <div dangerouslySetInnerHTML={{ __html: msg.content }} />
            </Box>
          ))}
          {loading && <Typography sx={{ color: 'text.secondary' }}>AI is typing...</Typography>}
        </Box>
        <FloatingInputBox>
          <ChatInput onSend={sendMessage} disabled={loading} />
        </FloatingInputBox>
      </Container>
    </>
  );
}