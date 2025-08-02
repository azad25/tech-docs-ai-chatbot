import { useState, useEffect, useRef } from "react";
import { CssBaseline, Box } from "@mui/material";
import { motion } from "framer-motion";
import Sidebar from "./components/Sidebar";
import ChatList from "./components/ChatList";
import ChatWindow from "./components/ChatWindow";

interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

interface ChatWindowMessage {
  id: string;
  content: string;
  sender: "user" | "other";
  timestamp: string;
  isFile?: boolean;
  fileName?: string;
  fileSize?: string;
  aiLearning?: boolean;
  confidence?: number;
  category?: string;
}

export default function App() {
  const [chatHistory, setChatHistory] = useState<ChatMessage[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    // Use the same host as the current page, but with ws protocol
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    ws.current = new WebSocket(`${wsProtocol}//${wsHost}/ws`);

    ws.current.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "chat_response") {
          setChatHistory((prev) => [...prev, { role: "assistant", content: data.response || data.message }]);
          setLoading(false);
        } else if (data.type === "typing") {
          setLoading(true);
        } else if (data.type === "error") {
          setLoading(false);
          setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style="color:red">Error: ${data.error}</p>` }]);
        }
      } catch (e) {
        setLoading(false);
        setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style="color:red">Invalid response from server</p>` }]);
      }
    };

    ws.current.onerror = () => {
      setLoading(false);
      setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style="color:red">WebSocket error occurred</p>` }]);
    };

    ws.current.onclose = () => {
      setLoading(false);
      setChatHistory((prev) => [...prev, { role: "assistant", content: `<p style="color:red">Disconnected from server</p>` }]);
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

  // Convert chat history to ChatWindow format with AI learning features
  const chatWindowMessages: ChatWindowMessage[] = chatHistory.map((msg, index) => {
    const isAI = msg.role === "assistant";
    
    return {
      id: `msg-${index}`,
      content: msg.content,
      sender: msg.role === "user" ? "user" : "other",
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      aiLearning: isAI, // AI messages show learning indicators
      confidence: isAI ? 0.85 + Math.random() * 0.1 : undefined, // Simulate confidence scores
      category: isAI ? "TechDocs AI" : undefined,
    };
  });

  return (
    <>
      <CssBaseline />
      <Box sx={{ 
        height: "100vh", 
        overflow: "hidden",
        background: "linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%)",
        fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
      }}>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ 
            duration: 1, 
            ease: [0.25, 0.46, 0.45, 0.94]
          }}
        >
          <Sidebar />
          <ChatList />
          <ChatWindow 
            messages={chatWindowMessages}
            onSendMessage={sendMessage}
            loading={loading}
          />
        </motion.div>
      </Box>
    </>
  );
}