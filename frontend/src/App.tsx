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
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 5;
    const reconnectDelay = 2000;

    const connectWebSocket = () => {
      // Use the same host as the current page, but with ws protocol
      const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsHost = window.location.host;
      ws.current = new WebSocket(`${wsProtocol}//${wsHost}/ws`);

      ws.current.onopen = () => {
        console.log("WebSocket connected");
        reconnectAttempts = 0; // Reset reconnect attempts on successful connection
      };

      ws.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log("WebSocket message received:", data);
          
          if (data.type === "connection") {
            console.log("WebSocket connection confirmed:", data.response);
          } else if (data.type === "chat_response") {
            setChatHistory((prev) => [...prev, { role: "assistant", content: data.response || data.message }]);
            setLoading(false);
          } else if (data.type === "typing") {
            setLoading(true);
          } else if (data.type === "error") {
            setLoading(false);
            setChatHistory((prev) => [...prev, { role: "assistant", content: `❌ **Error occurred**\n\n${data.error}` }]);
          }
        } catch (e) {
          console.error("Failed to parse WebSocket message:", e);
          setLoading(false);
          setChatHistory((prev) => [...prev, { role: "assistant", content: "⚠️ **Invalid response from server**\n\nReceived an invalid response format. Please try again." }]);
        }
      };

      ws.current.onerror = (error) => {
        console.error("WebSocket error:", error);
        setLoading(false);
      };

      ws.current.onclose = (event) => {
        console.log("WebSocket closed:", event.code, event.reason);
        setLoading(false);
        
        // Only show disconnection message if it's not a normal closure and we're not reconnecting
        if (event.code !== 1000 && reconnectAttempts < maxReconnectAttempts) {
          console.log(`Attempting to reconnect... (${reconnectAttempts + 1}/${maxReconnectAttempts})`);
          reconnectAttempts++;
          setTimeout(connectWebSocket, reconnectDelay);
        } else if (reconnectAttempts >= maxReconnectAttempts) {
          setChatHistory((prev) => [...prev, { role: "assistant", content: "🔌 **Connection lost**\n\nUnable to reconnect to the server after multiple attempts. Please refresh the page." }]);
        }
      };
    };

    connectWebSocket();

    return () => {
      if (ws.current) {
        ws.current.close(1000, "Component unmounting"); // Normal closure
      }
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