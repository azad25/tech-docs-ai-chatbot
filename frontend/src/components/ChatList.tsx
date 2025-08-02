import * as React from "react";
import { 
  Box, 
  Typography, 
  TextField, 
  Avatar, 
  Chip,
  IconButton,
  InputAdornment,
  LinearProgress
} from "@mui/material";
import { motion, AnimatePresence } from "framer-motion";
import {
  Search as SearchIcon,
  Add as AddIcon,
  KeyboardArrowDown as ArrowDownIcon,
  FiberManualRecord as StatusIcon,
  AttachFile as FileIcon,
  Photo as PhotoIcon,
  Videocam as VideoIcon,
  Mic as MicIcon,
  Psychology as AIIcon,
  School as LearnIcon,
  AutoAwesome as SparkleIcon
} from "@mui/icons-material";
import TutorialExpanded from "./TutorialExpanded";

interface ChatPreview {
  id: string;
  name: string;
  avatar: string;
  status: string;
  lastMessage: string;
  timestamp: string;
  unreadCount: number;
  isOnline: boolean;
  hasVoiceMessage?: boolean;
  hasFiles?: boolean;
  hasPhotos?: boolean;
  hasVideos?: boolean;
  aiLearning?: boolean;
  category?: string;
  confidence?: number;
}

const mockChats: ChatPreview[] = [
  {
    id: "1",
    name: "HTML Basics Tutorial",
    avatar: "HT",
    status: "Learning from interaction",
    lastMessage: "I've learned about HTML structure and semantic elements. Would you like to explore CSS styling next?",
    timestamp: "2 minutes ago",
    unreadCount: 0,
    isOnline: true,
    category: "HTML",
    aiLearning: true,
    confidence: 0.95,
  },
  {
    id: "2",
    name: "JavaScript Promises",
    avatar: "JS",
    status: "Vector search active",
    lastMessage: "Found 3 relevant tutorials about async/await patterns. Here's a comprehensive guide...",
    timestamp: "5 minutes ago",
    unreadCount: 1,
    isOnline: true,
    hasFiles: true,
    category: "JavaScript",
    aiLearning: true,
    confidence: 0.88,
  },
  {
    id: "3",
    name: "React Hooks Deep Dive",
    avatar: "RH",
    status: "RAG processing",
    lastMessage: "Analyzing your code patterns to provide personalized React hooks examples...",
    timestamp: "10 minutes ago",
    unreadCount: 0,
    isOnline: true,
    category: "React",
    aiLearning: true,
    confidence: 0.92,
  },
  {
    id: "4",
    name: "Python Data Structures",
    avatar: "PY",
    status: "Learning from feedback",
    lastMessage: "Based on your previous questions, I've updated my knowledge about list comprehensions.",
    timestamp: "1 hour ago",
    unreadCount: 0,
    isOnline: false,
    category: "Python",
    aiLearning: true,
    confidence: 0.85,
  },
  {
    id: "5",
    name: "Docker Containerization",
    avatar: "DC",
    status: "Scraping new docs",
    lastMessage: "I'm currently learning about Docker best practices from the official documentation.",
    timestamp: "2 hours ago",
    unreadCount: 0,
    isOnline: true,
    category: "DevOps",
    aiLearning: true,
    confidence: 0.78,
  },
];

// Mock tutorial content for expansion
const mockTutorialContent = {
  "1": {
    id: "1",
    title: "HTML Basics Tutorial",
    category: "HTML",
    content: `
      <h1>HTML Basics: A Comprehensive Guide</h1>
      
      <p>HTML (HyperText Markup Language) is the standard markup language for creating web pages. It describes the structure of a web page semantically and originally included cues for the appearance of the document.</p>
      
      <h2>What is HTML?</h2>
      <p>HTML stands for HyperText Markup Language. It is the standard markup language for creating web pages. HTML describes the structure of a web page semantically and originally included cues for the appearance of the document.</p>
      
      <h2>Basic HTML Structure</h2>
      <p>Every HTML document has a basic structure that includes the following elements:</p>
      
      <pre><code>&lt;!DOCTYPE html&gt;
&lt;html&gt;
&lt;head&gt;
    &lt;title&gt;Page Title&lt;/title&gt;
&lt;/head&gt;
&lt;body&gt;
    &lt;h1&gt;This is a heading&lt;/h1&gt;
    &lt;p&gt;This is a paragraph.&lt;/p&gt;
&lt;/body&gt;
&lt;/html&gt;</code></pre>
      
      <h2>HTML Elements</h2>
      <p>HTML elements are the building blocks of HTML pages. An HTML element is defined by a start tag, some content, and an end tag.</p>
      
      <h3>Common HTML Tags</h3>
      <ul>
        <li><strong>&lt;h1&gt; to &lt;h6&gt;</strong> - Headings</li>
        <li><strong>&lt;p&gt;</strong> - Paragraphs</li>
        <li><strong>&lt;a&gt;</strong> - Links</li>
        <li><strong>&lt;img&gt;</strong> - Images</li>
        <li><strong>&lt;div&gt;</strong> - Divisions</li>
        <li><strong>&lt;span&gt;</strong> - Inline elements</li>
      </ul>
      
      <h2>Best Practices</h2>
      <blockquote>
        Always use semantic HTML elements to improve accessibility and SEO. Choose the most appropriate element for your content.
      </blockquote>
      
      <h3>Semantic Elements</h3>
      <p>Use semantic elements to give meaning to your content:</p>
      <ul>
        <li><code>&lt;header&gt;</code> - Defines a header for a document or section</li>
        <li><code>&lt;nav&gt;</code> - Defines navigation links</li>
        <li><code>&lt;main&gt;</code> - Specifies the main content of a document</li>
        <li><code>&lt;section&gt;</code> - Defines a section in a document</li>
        <li><code>&lt;article&gt;</code> - Defines an article</li>
        <li><code>&lt;footer&gt;</code> - Defines a footer for a document or section</li>
      </ul>
      
      <h2>Next Steps</h2>
      <p>Now that you understand HTML basics, you can:</p>
      <ol>
        <li>Learn CSS for styling</li>
        <li>Study JavaScript for interactivity</li>
        <li>Explore responsive design principles</li>
        <li>Practice building real projects</li>
      </ol>
    `,
    confidence: 0.95,
    timestamp: "2 minutes ago",
    avatar: "HT",
  },
  "2": {
    id: "2",
    title: "JavaScript Promises",
    category: "JavaScript",
    content: `
      <h1>JavaScript Promises: Complete Guide</h1>
      
      <p>Promises are a way to handle asynchronous operations in JavaScript. They represent a value that may not be available immediately but will be resolved at some point in the future.</p>
      
      <h2>What are Promises?</h2>
      <p>A Promise is an object representing the eventual completion or failure of an asynchronous operation. It has three states:</p>
      <ul>
        <li><strong>Pending</strong> - Initial state, neither fulfilled nor rejected</li>
        <li><strong>Fulfilled</strong> - Operation completed successfully</li>
        <li><strong>Rejected</strong> - Operation failed</li>
      </ul>
      
      <h2>Creating Promises</h2>
      <pre><code>const myPromise = new Promise((resolve, reject) => {
    // Async operation
    const success = true;
    
    if (success) {
        resolve('Operation completed successfully');
    } else {
        reject('Operation failed');
    }
});</code></pre>
      
      <h2>Using Promises</h2>
      <p>You can use promises with <code>.then()</code> and <code>.catch()</code> methods:</p>
      
      <pre><code>myPromise
    .then(result => {
        console.log('Success:', result);
    })
    .catch(error => {
        console.error('Error:', error);
    });</code></pre>
      
      <h2>Async/Await</h2>
      <p>Modern JavaScript provides async/await syntax for cleaner promise handling:</p>
      
      <pre><code>async function fetchData() {
    try {
        const response = await fetch('https://api.example.com/data');
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching data:', error);
    }
}</code></pre>
      
      <h2>Promise Methods</h2>
      <h3>Promise.all()</h3>
      <p>Waits for all promises to resolve:</p>
      <pre><code>const promises = [promise1, promise2, promise3];
Promise.all(promises)
    .then(results => {
        console.log('All promises resolved:', results);
    });</code></pre>
      
      <h3>Promise.race()</h3>
      <p>Returns the first promise to resolve or reject:</p>
      <pre><code>Promise.race([promise1, promise2])
    .then(result => {
        console.log('First to complete:', result);
    });</code></pre>
      
      <h2>Best Practices</h2>
      <blockquote>
        Always handle errors in promises. Use try-catch with async/await or .catch() with .then() to prevent unhandled promise rejections.
      </blockquote>
      
      <h2>Common Patterns</h2>
      <h3>Timeout Pattern</h3>
      <pre><code>function timeout(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function delayedOperation() {
    await timeout(1000);
    console.log('Operation completed after 1 second');
}</code></pre>
      
      <h3>Retry Pattern</h3>
      <pre><code>async function retryOperation(operation, maxRetries = 3) {
    for (let i = 0; i < maxRetries; i++) {
        try {
            return await operation();
        } catch (error) {
            if (i === maxRetries - 1) throw error;
            await timeout(1000 * Math.pow(2, i)); // Exponential backoff
        }
    }
}</code></pre>
    `,
    confidence: 0.88,
    timestamp: "5 minutes ago",
    avatar: "JS",
  },
  "3": {
    id: "3",
    title: "React Hooks Deep Dive",
    category: "React",
    content: `
      <h1>React Hooks: A Comprehensive Guide</h1>
      
      <p>React Hooks are functions that allow you to use state and other React features in functional components. They were introduced in React 16.8 and have become the standard way to write React components.</p>
      
      <h2>What are Hooks?</h2>
      <p>Hooks are functions that let you "hook into" React state and lifecycle features from function components. They don't work inside classes — they let you use React without classes.</p>
      
      <h2>Basic Hooks</h2>
      
      <h3>useState</h3>
      <p>The useState hook allows you to add state to functional components:</p>
      <pre><code>import React, { useState } from 'react';

function Counter() {
    const [count, setCount] = useState(0);
    
    return (
        &lt;div&gt;
            &lt;p&gt;You clicked {count} times&lt;/p&gt;
            &lt;button onClick={() => setCount(count + 1)}&gt;
                Click me
            &lt;/button&gt;
        &lt;/div&gt;
    );
}</code></pre>
      
      <h3>useEffect</h3>
      <p>The useEffect hook lets you perform side effects in function components:</p>
      <pre><code>import React, { useState, useEffect } from 'react';

function Example() {
    const [count, setCount] = useState(0);
    
    useEffect(() => {
        document.title = \`You clicked \${count} times\`;
    }, [count]); // Only re-run if count changes
    
    return (
        &lt;div&gt;
            &lt;p&gt;You clicked {count} times&lt;/p&gt;
            &lt;button onClick={() => setCount(count + 1)}&gt;
                Click me
            &lt;/button&gt;
        &lt;/div&gt;
    );
}</code></pre>
      
      <h2>Custom Hooks</h2>
      <p>You can create your own hooks to reuse stateful logic between components:</p>
      <pre><code>function useWindowSize() {
    const [size, setSize] = useState({
        width: window.innerWidth,
        height: window.innerHeight
    });
    
    useEffect(() => {
        const handleResize = () => {
            setSize({
                width: window.innerWidth,
                height: window.innerHeight
            });
        };
        
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);
    
    return size;
}</code></pre>
      
      <h2>Advanced Hooks</h2>
      
      <h3>useContext</h3>
      <p>useContext lets you subscribe to React context without introducing nesting:</p>
      <pre><code>const ThemeContext = React.createContext('light');

function ThemedButton() {
    const theme = useContext(ThemeContext);
    return &lt;button className={theme}&gt;Themed Button&lt;/button&gt;;
}</code></pre>
      
      <h3>useReducer</h3>
      <p>useReducer is an alternative to useState for complex state logic:</p>
      <pre><code>function reducer(state, action) {
    switch (action.type) {
        case 'increment':
            return { count: state.count + 1 };
        case 'decrement':
            return { count: state.count - 1 };
        default:
            throw new Error();
    }
}

function Counter() {
    const [state, dispatch] = useReducer(reducer, { count: 0 });
    
    return (
        &lt;div&gt;
            Count: {state.count}
            &lt;button onClick={() => dispatch({ type: 'increment' })}&gt;+&lt;/button&gt;
            &lt;button onClick={() => dispatch({ type: 'decrement' })}&gt;-&lt;/button&gt;
        &lt;/div&gt;
    );
}</code></pre>
      
      <h2>Best Practices</h2>
      <blockquote>
        Always follow the Rules of Hooks: only call hooks at the top level and only call hooks from React functions.
      </blockquote>
      
      <h3>Rules of Hooks</h3>
      <ul>
        <li>Only call hooks at the top level of your function</li>
        <li>Don't call hooks inside loops, conditions, or nested functions</li>
        <li>Only call hooks from React function components or custom hooks</li>
      </ul>
      
      <h2>Performance Optimization</h2>
      <h3>useMemo</h3>
      <p>useMemo memoizes expensive calculations:</p>
      <pre><code>const expensiveValue = useMemo(() => {
    return computeExpensiveValue(a, b);
}, [a, b]);</code></pre>
      
      <h3>useCallback</h3>
      <p>useCallback memoizes functions to prevent unnecessary re-renders:</p>
      <pre><code>const memoizedCallback = useCallback(() => {
    doSomething(a, b);
}, [a, b]);</code></pre>
    `,
    confidence: 0.92,
    timestamp: "10 minutes ago",
    avatar: "RH",
  },
};

const ChatList: React.FC = () => {
  const [expandedTutorial, setExpandedTutorial] = React.useState<string | null>(null);

  const handleTutorialClick = (chatId: string) => {
    setExpandedTutorial(chatId);
  };

  const handleCloseTutorial = () => {
    setExpandedTutorial(null);
  };

  return (
    <>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ 
          duration: 0.8, 
          delay: 0.3,
          ease: [0.25, 0.46, 0.45, 0.94]
        }}
      >
        <Box
          sx={{
            width: 380,
            height: "100vh",
            background: "linear-gradient(180deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
            backdropFilter: "blur(20px)",
            borderRight: "1px solid rgba(0,0,0,0.1)",
            display: "flex",
            flexDirection: "column",
            position: "fixed",
            left: 240,
            top: 0,
            zIndex: 1100,
            boxShadow: "0 0 40px rgba(0,0,0,0.1)",
            overflow: "hidden",
            fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
          }}
        >
          {/* Header */}
          <motion.div
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.6, delay: 0.4 }}
          >
            <Box sx={{ p: 3, pb: 2, borderBottom: "1px solid rgba(0,0,0,0.08)" }}>
              <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2, gap: 2 }}>
                <Typography variant="h5" sx={{ fontWeight: 800, color: "#1D1D1F", letterSpacing: "-0.5px", flexGrow: 1 }}>
                  AI Conversations
                </Typography>
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
              </Box>
              
              <Box sx={{ display: "flex", alignItems: "center", gap: 1, mb: 2 }}>
                <Typography variant="body2" sx={{ color: "#86868B", fontWeight: 600 }}>
                  Learning Sessions
                </Typography>
                <motion.div
                  whileHover={{ rotate: 180 }}
                  transition={{ duration: 0.3 }}
                >
                  <IconButton size="small">
                    <ArrowDownIcon sx={{ color: "#86868B" }} />
                  </IconButton>
                </motion.div>
              </Box>

              {/* Search Bar */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.6 }}
              >
                <TextField
                  fullWidth
                  placeholder="Search conversations..."
                  variant="outlined"
                  size="small"
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <SearchIcon sx={{ color: "#86868B" }} />
                      </InputAdornment>
                    ),
                  }}
                  sx={{
                    "& .MuiOutlinedInput-root": {
                      borderRadius: 3,
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
              </motion.div>

              {/* AI Learning Status */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.7 }}
              >
                <Box sx={{ mt: 2, display: "flex", alignItems: "center", gap: 2 }}>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                    <SparkleIcon sx={{ color: "#34C759", fontSize: "1rem" }} />
                    <Typography variant="caption" sx={{ color: "#34C759", fontWeight: 600 }}>
                      AI Learning Active
                    </Typography>
                  </Box>
                  <LinearProgress 
                    variant="determinate" 
                    value={75} 
                    sx={{ 
                      flexGrow: 1, 
                      height: 4, 
                      borderRadius: 2,
                      backgroundColor: "rgba(52,199,89,0.2)",
                      "& .MuiLinearProgress-bar": {
                        background: "linear-gradient(90deg, #34C759 0%, #30D158 100%)",
                        borderRadius: 2,
                      }
                    }} 
                  />
                </Box>
              </motion.div>
            </Box>
          </motion.div>

          {/* Chat List */}
          <Box sx={{ flexGrow: 1, overflowY: "auto", p: 2, pb: 3 }}>
            <AnimatePresence>
              {mockChats.map((chat, index) => (
                <motion.div
                  key={chat.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ 
                    duration: 0.5, 
                    delay: 0.8 + index * 0.1,
                    ease: [0.25, 0.46, 0.45, 0.94]
                  }}
                  whileHover={{ 
                    scale: 1.02,
                    transition: { duration: 0.2 }
                  }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Box
                    onClick={() => handleTutorialClick(chat.id)}
                    sx={{
                      p: 2.5,
                      mb: 2,
                      borderRadius: 4,
                      background: "linear-gradient(135deg, rgba(255,255,255,0.9) 0%, rgba(248,249,250,0.9) 100%)",
                      backdropFilter: "blur(10px)",
                      cursor: "pointer",
                      border: "1px solid rgba(0,0,0,0.08)",
                      boxShadow: "0 4px 20px rgba(0,0,0,0.05)",
                      "&:hover": {
                        background: "linear-gradient(135deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
                        borderColor: "rgba(0,122,255,0.3)",
                        transform: "translateY(-2px)",
                        boxShadow: "0 8px 32px rgba(0,122,255,0.15)",
                      },
                    }}
                  >
                    <Box sx={{ display: "flex", alignItems: "flex-start", gap: 2 }}>
                      <motion.div
                        whileHover={{ scale: 1.05 }}
                        whileTap={{ scale: 0.95 }}
                      >
                        <Avatar
                          sx={{
                            width: 52,
                            height: 52,
                            background: chat.isOnline 
                              ? "linear-gradient(135deg, #34C759 0%, #30D158 100%)" 
                              : "linear-gradient(135deg, #8E8E93 0%, #AEAEB2 100%)",
                            fontSize: "1.1rem",
                            fontWeight: "bold",
                            boxShadow: chat.isOnline 
                              ? "0 8px 32px rgba(52,199,89,0.3)" 
                              : "0 4px 16px rgba(142,142,147,0.2)",
                          }}
                        >
                          {chat.avatar}
                        </Avatar>
                      </motion.div>
                      
                      <Box sx={{ flexGrow: 1, minWidth: 0 }}>
                        <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", mb: 1, gap: 2 }}>
                          <Typography variant="subtitle1" sx={{ fontWeight: 700, color: "#1D1D1F", flexGrow: 1 }}>
                            {chat.name}
                          </Typography>
                          <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.7rem", flexShrink: 0 }}>
                            {chat.timestamp}
                          </Typography>
                        </Box>
                        
                        <Box sx={{ display: "flex", alignItems: "center", gap: 1, mb: 1 }}>
                          <StatusIcon 
                            sx={{ 
                              fontSize: "0.6rem", 
                              color: chat.isOnline ? "#34C759" : "#86868B" 
                            }} 
                          />
                          <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.7rem" }}>
                            {chat.status}
                          </Typography>
                          {chat.aiLearning && (
                            <motion.div
                              initial={{ scale: 0 }}
                              animate={{ scale: 1 }}
                              transition={{ duration: 0.3, delay: 1 }}
                            >
                              <AIIcon sx={{ fontSize: "0.8rem", color: "#007AFF", ml: 1 }} />
                            </motion.div>
                          )}
                        </Box>
                        
                        <Typography 
                          variant="body2" 
                          sx={{ 
                            color: "#1D1D1F", 
                            fontSize: "0.8rem",
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                            display: "-webkit-box",
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: "vertical",
                            mb: 1.5,
                            lineHeight: 1.4,
                          }}
                        >
                          {chat.lastMessage}
                        </Typography>
                        
                        {/* AI Learning Indicators */}
                        {chat.aiLearning && (
                          <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 1, flexWrap: "wrap" }}>
                            <Chip
                              icon={<LearnIcon />}
                              label={chat.category}
                              size="small"
                              sx={{ 
                                backgroundColor: "rgba(0,122,255,0.1)",
                                color: "#007AFF",
                                fontSize: "0.6rem",
                                fontWeight: 600,
                                height: 20,
                                "& .MuiChip-icon": { fontSize: "0.7rem" }
                              }}
                            />
                            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                              <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.6rem" }}>
                                Confidence:
                              </Typography>
                              <Typography variant="caption" sx={{ color: "#34C759", fontWeight: 600, fontSize: "0.6rem" }}>
                                {Math.round(chat.confidence! * 100)}%
                              </Typography>
                            </Box>
                          </Box>
                        )}
                        
                        {/* Action buttons */}
                        {(chat.hasFiles || chat.hasPhotos || chat.hasVideos || chat.hasVoiceMessage) && (
                          <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                            {chat.hasVoiceMessage && (
                              <Chip
                                icon={<MicIcon />}
                                label="Voice"
                                size="small"
                                sx={{ 
                                  backgroundColor: "rgba(0,122,255,0.1)", 
                                  color: "#007AFF",
                                  fontSize: "0.6rem",
                                  height: 20
                                }}
                              />
                            )}
                            {chat.hasFiles && (
                              <Chip
                                icon={<FileIcon />}
                                label="Files"
                                size="small"
                                sx={{ 
                                  backgroundColor: "rgba(52,199,89,0.1)", 
                                  color: "#34C759",
                                  fontSize: "0.6rem",
                                  height: 20
                                }}
                              />
                            )}
                            {chat.hasPhotos && (
                              <Chip
                                icon={<PhotoIcon />}
                                label="Photos"
                                size="small"
                                sx={{ 
                                  backgroundColor: "rgba(255,149,0,0.1)", 
                                  color: "#FF9500",
                                  fontSize: "0.6rem",
                                  height: 20
                                }}
                              />
                            )}
                            {chat.hasVideos && (
                              <Chip
                                icon={<VideoIcon />}
                                label="Videos"
                                size="small"
                                sx={{ 
                                  backgroundColor: "rgba(175,82,222,0.1)", 
                                  color: "#AF52DE",
                                  fontSize: "0.6rem",
                                  height: 20
                                }}
                              />
                            )}
                          </Box>
                        )}
                      </Box>
                      
                      {chat.unreadCount > 0 && (
                        <motion.div
                          initial={{ scale: 0 }}
                          animate={{ scale: 1 }}
                          transition={{ duration: 0.3, delay: 1.2 }}
                        >
                          <Box
                            sx={{
                              width: 24,
                              height: 24,
                              borderRadius: "50%",
                              background: "linear-gradient(135deg, #FF3B30 0%, #FF453A 100%)",
                              color: "white",
                              display: "flex",
                              alignItems: "center",
                              justifyContent: "center",
                              fontSize: "0.7rem",
                              fontWeight: "bold",
                              boxShadow: "0 4px 16px rgba(255,59,48,0.3)",
                              flexShrink: 0,
                              ml: 1,
                            }}
                          >
                            {chat.unreadCount}
                          </Box>
                        </motion.div>
                      )}
                    </Box>
                  </Box>
                </motion.div>
              ))}
            </AnimatePresence>
          </Box>
        </Box>
      </motion.div>

      {/* Tutorial Expanded View */}
      <TutorialExpanded
        isOpen={!!expandedTutorial}
        onClose={handleCloseTutorial}
        tutorial={expandedTutorial ? mockTutorialContent[expandedTutorial as keyof typeof mockTutorialContent] : null}
      />
    </>
  );
};

export default ChatList; 