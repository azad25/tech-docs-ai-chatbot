import * as React from "react";
import { 
  Box, 
  Typography, 
  Avatar, 
  IconButton, 
  TextField,
  Paper,
  LinearProgress
} from "@mui/material";
import { motion, AnimatePresence } from "framer-motion";
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import {
  AttachFile as AttachIcon,
  MoreVert as MoreIcon,
  Add as AddIcon,
  EmojiEmotions as EmojiIcon,
  Send as SendIcon,
  Videocam as VideoIcon,
  Photo as PhotoIcon,
  Description as DocumentIcon,
  Add as PlusIcon,
  Psychology as AIIcon,
  AutoAwesome as SparkleIcon
} from "@mui/icons-material";

interface ChatMessage {
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

interface ChatWindowProps {
  messages: ChatMessage[];
  onSendMessage: (message: string) => void;
  loading?: boolean;
}

const ChatWindow: React.FC<ChatWindowProps> = ({ messages, onSendMessage, loading = false }) => {
  const [inputValue, setInputValue] = React.useState("");
  const messagesEndRef = React.useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  React.useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = () => {
    if (inputValue.trim() && !loading) {
      onSendMessage(inputValue.trim());
      setInputValue("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ 
        duration: 0.8, 
        delay: 0.5,
        ease: [0.25, 0.46, 0.45, 0.94]
      }}
    >
      <Box
        sx={{
          position: "fixed",
          left: 620,
          top: 0,
          right: 0,
          height: "100vh",
          background: "linear-gradient(180deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
          backdropFilter: "blur(20px)",
          display: "flex",
          flexDirection: "column",
          zIndex: 1000,
          fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
        }}
      >
        {/* Header */}
        <motion.div
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.6 }}
        >
          <Box
            sx={{
              p: 3,
              borderBottom: "1px solid rgba(0,0,0,0.08)",
              display: "flex",
              alignItems: "center",
              gap: 2,
            }}
          >
            <motion.div
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <Avatar
                sx={{
                  width: 56,
                  height: 56,
                  background: "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)",
                  fontSize: "1.4rem",
                  fontWeight: "bold",
                  boxShadow: "0 8px 32px rgba(0,122,255,0.3)",
                }}
              >
                AI
              </Avatar>
            </motion.div>
            <Box sx={{ flexGrow: 1 }}>
              <Typography variant="h6" sx={{ fontWeight: 700, color: "#1D1D1F", mb: 0.5 }}>
                TechDocs AI Assistant
              </Typography>
              <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                  <Box
                    sx={{
                      width: 8,
                      height: 8,
                      borderRadius: "50%",
                      backgroundColor: "#34C759",
                      boxShadow: "0 0 8px rgba(52,199,89,0.5)",
                    }}
                  />
                  <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.8rem" }}>
                    {loading ? "Learning & Processing..." : "Online & Learning"}
                  </Typography>
                </Box>
                {loading && (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.3 }}
                  >
                    <SparkleIcon sx={{ color: "#34C759", fontSize: "1rem" }} />
                  </motion.div>
                )}
              </Box>
            </Box>
            <motion.div
              whileHover={{ rotate: 5 }}
              transition={{ duration: 0.2 }}
            >
              <IconButton size="small">
                <AttachIcon sx={{ color: "#86868B" }} />
              </IconButton>
            </motion.div>
            <motion.div
              whileHover={{ rotate: 5 }}
              transition={{ duration: 0.2 }}
            >
              <IconButton size="small">
                <MoreIcon sx={{ color: "#86868B" }} />
              </IconButton>
            </motion.div>
          </Box>
        </motion.div>

        {/* Messages */}
        <Box
          sx={{
            flexGrow: 1,
            overflowY: "auto",
            p: 3,
            display: "flex",
            flexDirection: "column",
            gap: 2,
          }}
        >
          <AnimatePresence>
            {messages.map((message, index) => (
              <motion.div
                key={message.id}
                initial={{ opacity: 0, y: 20, scale: 0.95 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ 
                  duration: 0.5, 
                  delay: index * 0.1,
                  ease: [0.25, 0.46, 0.45, 0.94]
                }}
                exit={{ opacity: 0, y: -20, scale: 0.95 }}
              >
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: message.sender === "user" ? "flex-end" : "flex-start",
                    mb: 2,
                  }}
                >
                  {message.sender === "other" && (
                    <motion.div
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                    >
                      <Avatar
                        sx={{
                          width: 36,
                          height: 36,
                          background: "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)",
                          fontSize: "0.9rem",
                          fontWeight: "bold",
                          mr: 1,
                          mt: 0.5,
                          boxShadow: "0 4px 16px rgba(0,122,255,0.2)",
                        }}
                      >
                        AI
                      </Avatar>
                    </motion.div>
                  )}
                  
                  <Box
                    sx={{
                      maxWidth: "75%",
                      display: "flex",
                      flexDirection: "column",
                      alignItems: message.sender === "user" ? "flex-end" : "flex-start",
                    }}
                  >
                    <motion.div
                      whileHover={{ scale: 1.02 }}
                      transition={{ duration: 0.2 }}
                    >
                      <Paper
                        elevation={0}
                        sx={{
                          p: 2.5,
                          background: message.sender === "user" 
                            ? "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)" 
                            : "linear-gradient(135deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
                          color: message.sender === "user" ? "white" : "#1D1D1F",
                          borderRadius: 4,
                          position: "relative",
                          backdropFilter: "blur(10px)",
                          border: message.sender === "user" 
                            ? "none" 
                            : "1px solid rgba(0,0,0,0.08)",
                          boxShadow: message.sender === "user"
                            ? "0 8px 32px rgba(0,122,255,0.3)"
                            : "0 4px 20px rgba(0,0,0,0.05)",
                          "&::after": message.sender === "user" ? {
                            content: '""',
                            position: "absolute",
                            right: -8,
                            top: 16,
                            width: 0,
                            height: 0,
                            borderLeft: "8px solid #007AFF",
                            borderTop: "8px solid transparent",
                            borderBottom: "8px solid transparent",
                          } : {
                            content: '""',
                            position: "absolute",
                            left: -8,
                            top: 16,
                            width: 0,
                            height: 0,
                            borderRight: "8px solid rgba(255,255,255,0.95)",
                            borderTop: "8px solid transparent",
                            borderBottom: "8px solid transparent",
                          }
                        }}
                      >
                        {message.isFile ? (
                          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                            <AttachIcon sx={{ fontSize: "1.2rem" }} />
                            <Box>
                              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                                {message.fileName}
                              </Typography>
                              <Typography variant="caption" sx={{ opacity: 0.8 }}>
                                {message.fileSize}
                              </Typography>
                            </Box>
                          </Box>
                        ) : (
                          <Box
                            sx={{
                              lineHeight: 1.6,
                              fontSize: "0.9rem",
                              color: message.sender === "user" ? "white" : "#1D1D1F",
                              "& p": { 
                                margin: "0 0 12px 0",
                                fontSize: "0.9rem",
                                lineHeight: 1.6
                              },
                              "& p:last-child": { mb: 0 },
                              "& h1, & h2, & h3, & h4, & h5, & h6": { 
                                margin: "20px 0 12px 0",
                                fontWeight: 700,
                                color: message.sender === "user" ? "white" : "#1D1D1F",
                                lineHeight: 1.3
                              },
                              "& h1": { 
                                fontSize: "1.4rem",
                                borderBottom: message.sender === "user" ? "2px solid rgba(255,255,255,0.3)" : "2px solid rgba(0,122,255,0.3)",
                                paddingBottom: "8px",
                                marginBottom: "16px"
                              },
                              "& h2": { 
                                fontSize: "1.2rem",
                                color: message.sender === "user" ? "rgba(255,255,255,0.95)" : "#007AFF"
                              },
                              "& h3": { 
                                fontSize: "1.1rem",
                                color: message.sender === "user" ? "rgba(255,255,255,0.9)" : "#5856D6"
                              },
                              "& code": { 
                                backgroundColor: message.sender === "user" ? "rgba(255,255,255,0.25)" : "rgba(0,122,255,0.1)",
                                color: message.sender === "user" ? "rgba(255,255,255,0.95)" : "#007AFF",
                                padding: "3px 6px",
                                borderRadius: "4px",
                                fontFamily: "'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace",
                                fontSize: "0.85rem",
                                fontWeight: 500
                              },
                              "& pre": {
                                backgroundColor: message.sender === "user" ? "rgba(255,255,255,0.15)" : "rgba(248,249,250,0.8)",
                                border: message.sender === "user" ? "1px solid rgba(255,255,255,0.2)" : "1px solid rgba(0,0,0,0.1)",
                                padding: "16px",
                                borderRadius: "8px",
                                overflow: "auto",
                                margin: "16px 0",
                                boxShadow: "inset 0 2px 4px rgba(0,0,0,0.1)"
                              },
                              "& pre code": {
                                backgroundColor: "transparent",
                                padding: 0,
                                color: message.sender === "user" ? "rgba(255,255,255,0.9)" : "#1D1D1F",
                                fontSize: "0.8rem",
                                lineHeight: 1.5
                              },
                              "& ul, & ol": { 
                                margin: "12px 0",
                                paddingLeft: "24px"
                              },
                              "& li": { 
                                margin: "6px 0",
                                lineHeight: 1.5
                              },
                              "& li::marker": {
                                color: message.sender === "user" ? "rgba(255,255,255,0.7)" : "#007AFF"
                              },
                              "& blockquote": {
                                borderLeft: `4px solid ${message.sender === "user" ? "rgba(255,255,255,0.5)" : "#007AFF"}`,
                                paddingLeft: "16px",
                                margin: "16px 0",
                                fontStyle: "italic",
                                opacity: 0.9,
                                backgroundColor: message.sender === "user" ? "rgba(255,255,255,0.1)" : "rgba(0,122,255,0.05)",
                                padding: "12px 16px",
                                borderRadius: "0 8px 8px 0"
                              },
                              "& strong": { 
                                fontWeight: 700,
                                color: message.sender === "user" ? "white" : "#1D1D1F"
                              },
                              "& em": { 
                                fontStyle: "italic",
                                color: message.sender === "user" ? "rgba(255,255,255,0.9)" : "#5856D6"
                              },
                              "& table": {
                                width: "100%",
                                borderCollapse: "collapse",
                                margin: "16px 0",
                                fontSize: "0.85rem"
                              },
                              "& th, & td": {
                                border: message.sender === "user" ? "1px solid rgba(255,255,255,0.3)" : "1px solid rgba(0,0,0,0.1)",
                                padding: "8px 12px",
                                textAlign: "left"
                              },
                              "& th": {
                                backgroundColor: message.sender === "user" ? "rgba(255,255,255,0.2)" : "rgba(0,122,255,0.1)",
                                fontWeight: 600
                              },
                              "& a": {
                                color: message.sender === "user" ? "rgba(255,255,255,0.9)" : "#007AFF",
                                textDecoration: "underline",
                                "&:hover": {
                                  opacity: 0.8
                                }
                              }
                            }}
                          >
                            <ReactMarkdown
                              remarkPlugins={[remarkGfm]}
                              components={{
                                // Custom components for better styling
                                h1: ({ children }) => (
                                  <Typography variant="h4" component="h1" sx={{ 
                                    fontSize: "1.4rem",
                                    fontWeight: 700,
                                    margin: "20px 0 12px 0",
                                    color: message.sender === "user" ? "white" : "#1D1D1F",
                                    borderBottom: message.sender === "user" ? "2px solid rgba(255,255,255,0.3)" : "2px solid rgba(0,122,255,0.3)",
                                    paddingBottom: "8px"
                                  }}>
                                    {children}
                                  </Typography>
                                ),
                                h2: ({ children }) => (
                                  <Typography variant="h5" component="h2" sx={{ 
                                    fontSize: "1.2rem",
                                    fontWeight: 700,
                                    margin: "18px 0 10px 0",
                                    color: message.sender === "user" ? "rgba(255,255,255,0.95)" : "#007AFF"
                                  }}>
                                    {children}
                                  </Typography>
                                ),
                                h3: ({ children }) => (
                                  <Typography variant="h6" component="h3" sx={{ 
                                    fontSize: "1.1rem",
                                    fontWeight: 600,
                                    margin: "16px 0 8px 0",
                                    color: message.sender === "user" ? "rgba(255,255,255,0.9)" : "#5856D6"
                                  }}>
                                    {children}
                                  </Typography>
                                ),
                                p: ({ children }) => (
                                  <Typography variant="body2" sx={{ 
                                    margin: "0 0 12px 0",
                                    fontSize: "0.9rem",
                                    lineHeight: 1.6,
                                    color: message.sender === "user" ? "white" : "#1D1D1F"
                                  }}>
                                    {children}
                                  </Typography>
                                )
                              }}
                            >
                              {message.content}
                            </ReactMarkdown>
                          </Box>
                        )}
                        
                        {/* AI Learning Indicators */}
                        {message.aiLearning && (
                          <motion.div
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ duration: 0.3, delay: 0.5 }}
                          >
                            <Box sx={{ display: "flex", alignItems: "center", gap: 1, mt: 2, pt: 2, borderTop: "1px solid rgba(255,255,255,0.2)" }}>
                              <AIIcon sx={{ fontSize: "0.8rem", color: message.sender === "user" ? "rgba(255,255,255,0.8)" : "#007AFF" }} />
                              <Typography variant="caption" sx={{ 
                                color: message.sender === "user" ? "rgba(255,255,255,0.8)" : "#86868B",
                                fontSize: "0.7rem"
                              }}>
                                Learning from this interaction
                              </Typography>
                              {message.confidence && (
                                <Typography variant="caption" sx={{ 
                                  color: message.sender === "user" ? "rgba(255,255,255,0.8)" : "#34C759",
                                  fontWeight: 600,
                                  fontSize: "0.7rem"
                                }}>
                                  {Math.round(message.confidence * 100)}% confidence
                                </Typography>
                              )}
                            </Box>
                          </motion.div>
                        )}
                      </Paper>
                    </motion.div>
                    
                    <Typography 
                      variant="caption" 
                      sx={{ 
                        color: "#86868B", 
                        mt: 1, 
                        fontSize: "0.7rem",
                        textAlign: message.sender === "user" ? "right" : "left"
                      }}
                    >
                      {message.timestamp}
                    </Typography>
                  </Box>
                </Box>
              </motion.div>
            ))}
          </AnimatePresence>
          
          {loading && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5 }}
            >
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "flex-start",
                  mb: 2,
                }}
              >
                <Avatar
                  sx={{
                    width: 36,
                    height: 36,
                    background: "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)",
                    fontSize: "0.9rem",
                    fontWeight: "bold",
                    mr: 1,
                    mt: 0.5,
                    boxShadow: "0 4px 16px rgba(0,122,255,0.2)",
                  }}
                >
                  AI
                </Avatar>
                <Paper
                  elevation={0}
                  sx={{
                    p: 2.5,
                    background: "linear-gradient(135deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
                    color: "#1D1D1F",
                    borderRadius: 4,
                    position: "relative",
                    backdropFilter: "blur(10px)",
                    border: "1px solid rgba(0,0,0,0.08)",
                    boxShadow: "0 4px 20px rgba(0,0,0,0.05)",
                    minWidth: "200px",
                    "&::after": {
                      content: '""',
                      position: "absolute",
                      left: -8,
                      top: 16,
                      width: 0,
                      height: 0,
                      borderRight: "8px solid rgba(255,255,255,0.95)",
                      borderTop: "8px solid transparent",
                      borderBottom: "8px solid transparent",
                    }
                  }}
                >
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 1 }}>
                    <SparkleIcon sx={{ color: "#34C759", fontSize: "1rem" }} />
                    <Typography variant="body2" sx={{ color: "#34C759", fontWeight: 600 }}>
                      AI Learning & Processing
                    </Typography>
                  </Box>
                  <LinearProgress 
                    sx={{ 
                      height: 3, 
                      borderRadius: 2,
                      backgroundColor: "rgba(52,199,89,0.2)",
                      "& .MuiLinearProgress-bar": {
                        background: "linear-gradient(90deg, #34C759 0%, #30D158 100%)",
                        borderRadius: 2,
                      }
                    }} 
                  />
                  <Typography variant="caption" sx={{ color: "#86868B", mt: 1, display: "block" }}>
                    Analyzing context and generating response...
                  </Typography>
                </Paper>
              </Box>
            </motion.div>
          )}
          
          <div ref={messagesEndRef} />
        </Box>

        {/* Input Area */}
        <motion.div
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.8 }}
        >
          <Box sx={{ p: 3, borderTop: "1px solid rgba(0,0,0,0.08)" }}>
            <Box sx={{ display: "flex", alignItems: "flex-end", gap: 2, mb: 2 }}>
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <IconButton
                  sx={{
                    background: "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)",
                    color: "white",
                    width: 44,
                    height: 44,
                    boxShadow: "0 8px 32px rgba(0,122,255,0.3)",
                    "&:hover": {
                      background: "linear-gradient(135deg, #0056CC 0%, #4A4AC8 100%)",
                      transform: "translateY(-2px)",
                      boxShadow: "0 12px 40px rgba(0,122,255,0.4)",
                    },
                  }}
                >
                  <AddIcon />
                </IconButton>
              </motion.div>
              
              <TextField
                fullWidth
                multiline
                maxRows={4}
                placeholder="Ask TechDocs AI about programming, tutorials, or documentation..."
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyDown={handleKeyDown}
                variant="outlined"
                size="small"
                disabled={loading}
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 4,
                    backgroundColor: "rgba(255,255,255,0.8)",
                    backdropFilter: "blur(10px)",
                    border: "1px solid rgba(0,0,0,0.1)",
                    "& fieldset": {
                      borderColor: "transparent",
                    },
                    "&:hover fieldset": {
                      borderColor: "rgba(0,122,255,0.3)",
                    },
                    "&.Mui-focused fieldset": {
                      borderColor: "#007AFF",
                    },
                  },
                }}
              />
              
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <IconButton size="small">
                  <EmojiIcon sx={{ color: "#86868B" }} />
                </IconButton>
              </motion.div>
              
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <IconButton
                  onClick={handleSend}
                  disabled={!inputValue.trim() || loading}
                  sx={{
                    background: "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)",
                    color: "white",
                    width: 44,
                    height: 44,
                    boxShadow: "0 8px 32px rgba(0,122,255,0.3)",
                    "&:hover": {
                      background: "linear-gradient(135deg, #0056CC 0%, #4A4AC8 100%)",
                      transform: "translateY(-2px)",
                      boxShadow: "0 12px 40px rgba(0,122,255,0.4)",
                    },
                    "&:disabled": {
                      background: "linear-gradient(135deg, #8E8E93 0%, #AEAEB2 100%)",
                      color: "rgba(255,255,255,0.5)",
                      boxShadow: "none",
                    },
                  }}
                >
                  <SendIcon />
                </IconButton>
              </motion.div>
            </Box>
            
            {/* Quick Action Icons */}
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 1 }}
            >
              <Box sx={{ display: "flex", gap: 3, justifyContent: "center" }}>
                <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                  <IconButton size="small" sx={{ color: "#86868B" }}>
                    <VideoIcon />
                  </IconButton>
                </motion.div>
                <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                  <IconButton size="small" sx={{ color: "#86868B" }}>
                    <PhotoIcon />
                  </IconButton>
                </motion.div>
                <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                  <IconButton size="small" sx={{ color: "#86868B" }}>
                    <DocumentIcon />
                  </IconButton>
                </motion.div>
                <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
                  <IconButton size="small" sx={{ color: "#86868B" }}>
                    <PlusIcon />
                  </IconButton>
                </motion.div>
              </Box>
            </motion.div>
          </Box>
        </motion.div>
      </Box>
    </motion.div>
  );
};

export default ChatWindow; 