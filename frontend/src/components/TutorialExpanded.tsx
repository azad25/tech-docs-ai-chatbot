import * as React from "react";
import { 
  Box, 
  Typography, 
  Avatar, 
  IconButton, 
  Paper,
  Chip,
  Divider,
  LinearProgress
} from "@mui/material";
import { motion, AnimatePresence } from "framer-motion";
import {
  Close as CloseIcon,
  Psychology as AIIcon,
  AutoAwesome as SparkleIcon,
  School as LearnIcon,
  Book as BookIcon,
  Lightbulb as TipIcon
} from "@mui/icons-material";

interface TutorialExpandedProps {
  isOpen: boolean;
  onClose: () => void;
  tutorial: {
    id: string;
    title: string;
    category: string;
    content: string;
    confidence: number;
    timestamp: string;
    avatar: string;
  } | null;
}

const TutorialExpanded: React.FC<TutorialExpandedProps> = ({ isOpen, onClose, tutorial }) => {
  if (!tutorial) return null;

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0, x: 300 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: 300 }}
          transition={{ 
            duration: 0.6, 
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
              background: "linear-gradient(180deg, rgba(255,255,255,0.98) 0%, rgba(248,249,250,0.98) 100%)",
              backdropFilter: "blur(20px)",
              display: "flex",
              flexDirection: "column",
              boxShadow: "0 0 40px rgba(0,0,0,0.1)",
              zIndex: 1300,
              fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
            }}
          >
            {/* Header */}
            <motion.div
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ duration: 0.5, delay: 0.2 }}
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
                      width: 48,
                      height: 48,
                      background: "linear-gradient(135deg, #007AFF 0%, #5856D6 100%)",
                      fontSize: "1.2rem",
                      fontWeight: "bold",
                      boxShadow: "0 8px 32px rgba(0,122,255,0.3)",
                    }}
                  >
                    {tutorial.avatar}
                  </Avatar>
                </motion.div>
                
                <Box sx={{ flexGrow: 1 }}>
                  <Typography variant="h5" sx={{ fontWeight: 700, color: "#1D1D1F", mb: 0.5 }}>
                    {tutorial.title}
                  </Typography>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <Chip
                      icon={<LearnIcon />}
                      label={tutorial.category}
                      size="small"
                      sx={{
                        backgroundColor: "rgba(0,122,255,0.1)",
                        color: "#007AFF",
                        fontSize: "0.7rem",
                        fontWeight: 600,
                      }}
                    />
                    <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.8rem" }}>
                      {tutorial.timestamp}
                    </Typography>
                  </Box>
                </Box>
                
                <motion.div
                  whileHover={{ rotate: 90 }}
                  transition={{ duration: 0.3 }}
                >
                  <IconButton 
                    onClick={onClose}
                    sx={{
                      backgroundColor: "rgba(0,0,0,0.05)",
                      "&:hover": {
                        backgroundColor: "rgba(0,0,0,0.1)",
                      }
                    }}
                  >
                    <CloseIcon sx={{ color: "#86868B" }} />
                  </IconButton>
                </motion.div>
              </Box>
            </motion.div>

            {/* Tutorial Content */}
            <Box
              sx={{
                flexGrow: 1,
                overflowY: "auto",
                p: 4,
              }}
            >
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.4 }}
              >
                {/* AI Learning Status */}
                <Paper
                  elevation={0}
                  sx={{
                    p: 3,
                    mb: 4,
                    background: "linear-gradient(135deg, rgba(52,199,89,0.1) 0%, rgba(48,209,88,0.1) 100%)",
                    border: "1px solid rgba(52,199,89,0.2)",
                    borderRadius: 3,
                  }}
                >
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 2 }}>
                    <SparkleIcon sx={{ color: "#34C759", fontSize: "1.2rem" }} />
                    <Typography variant="h6" sx={{ color: "#34C759", fontWeight: 600 }}>
                      AI Learning Summary
                    </Typography>
                  </Box>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 3, mb: 2 }}>
                    <Box>
                      <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.8rem" }}>
                        Confidence Level
                      </Typography>
                      <Typography variant="h6" sx={{ color: "#34C759", fontWeight: 700 }}>
                        {Math.round(tutorial.confidence * 100)}%
                      </Typography>
                    </Box>
                    <Box>
                      <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.8rem" }}>
                        Category
                      </Typography>
                      <Typography variant="h6" sx={{ color: "#007AFF", fontWeight: 700 }}>
                        {tutorial.category}
                      </Typography>
                    </Box>
                    <Box>
                      <Typography variant="caption" sx={{ color: "#86868B", fontSize: "0.8rem" }}>
                        Learning Progress
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={tutorial.confidence * 100} 
                        sx={{ 
                          width: 100,
                          height: 6, 
                          borderRadius: 3,
                          backgroundColor: "rgba(52,199,89,0.2)",
                          "& .MuiLinearProgress-bar": {
                            background: "linear-gradient(90deg, #34C759 0%, #30D158 100%)",
                            borderRadius: 3,
                          }
                        }} 
                      />
                    </Box>
                  </Box>
                </Paper>

                {/* Tutorial Content */}
                <Paper
                  elevation={0}
                  sx={{
                    p: 4,
                    background: "linear-gradient(135deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
                    backdropFilter: "blur(10px)",
                    border: "1px solid rgba(0,0,0,0.08)",
                    borderRadius: 4,
                    boxShadow: "0 8px 32px rgba(0,0,0,0.05)",
                  }}
                >
                  <Box sx={{ mb: 4 }}>
                    <Typography variant="h4" sx={{ fontWeight: 800, color: "#1D1D1F", mb: 2 }}>
                      {tutorial.title}
                    </Typography>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}>
                      <Chip
                        icon={<BookIcon />}
                        label="Tutorial"
                        size="small"
                        sx={{
                          backgroundColor: "rgba(0,122,255,0.1)",
                          color: "#007AFF",
                          fontWeight: 600,
                        }}
                      />
                      <Chip
                        icon={<AIIcon />}
                        label="AI Generated"
                        size="small"
                        sx={{
                          backgroundColor: "rgba(52,199,89,0.1)",
                          color: "#34C759",
                          fontWeight: 600,
                        }}
                      />
                      <Chip
                        icon={<TipIcon />}
                        label="Best Practices"
                        size="small"
                        sx={{
                          backgroundColor: "rgba(255,149,0,0.1)",
                          color: "#FF9500",
                          fontWeight: 600,
                        }}
                      />
                    </Box>
                    <Divider sx={{ mb: 3 }} />
                  </Box>

                  {/* Tutorial Body */}
                  <Box
                    sx={{
                      "& h1, & h2, & h3, & h4, & h5, & h6": {
                        color: "#1D1D1F",
                        fontWeight: 700,
                        mb: 2,
                        mt: 4,
                      },
                      "& h1": { fontSize: "2rem" },
                      "& h2": { fontSize: "1.5rem" },
                      "& h3": { fontSize: "1.25rem" },
                      "& p": {
                        fontSize: "1rem",
                        lineHeight: 1.7,
                        color: "#333",
                        mb: 2,
                      },
                      "& code": {
                        backgroundColor: "rgba(0,0,0,0.1)",
                        padding: "2px 6px",
                        borderRadius: "4px",
                        fontFamily: "monospace",
                        fontSize: "0.9rem",
                        color: "#007AFF",
                      },
                      "& pre": {
                        backgroundColor: "rgba(0,0,0,0.05)",
                        padding: "16px",
                        borderRadius: "8px",
                        overflow: "auto",
                        margin: "16px 0",
                        border: "1px solid rgba(0,0,0,0.1)",
                        "& code": {
                          backgroundColor: "transparent",
                          padding: 0,
                          color: "#333",
                        }
                      },
                      "& ul, & ol": {
                        margin: "16px 0",
                        paddingLeft: "24px",
                        "& li": {
                          margin: "8px 0",
                          lineHeight: 1.6,
                        }
                      },
                      "& blockquote": {
                        borderLeft: "4px solid #007AFF",
                        paddingLeft: "16px",
                        margin: "16px 0",
                        fontStyle: "italic",
                        color: "#666",
                        backgroundColor: "rgba(0,122,255,0.05)",
                        padding: "12px 16px",
                        borderRadius: "0 8px 8px 0",
                      },
                      "& strong": { fontWeight: 700 },
                      "& em": { fontStyle: "italic" },
                      "& a": {
                        color: "#007AFF",
                        textDecoration: "none",
                        "&:hover": {
                          textDecoration: "underline",
                        }
                      },
                      "& table": {
                        width: "100%",
                        borderCollapse: "collapse",
                        margin: "16px 0",
                        "& th, & td": {
                          border: "1px solid rgba(0,0,0,0.1)",
                          padding: "8px 12px",
                          textAlign: "left",
                        },
                        "& th": {
                          backgroundColor: "rgba(0,122,255,0.1)",
                          fontWeight: 600,
                        }
                      }
                    }}
                    dangerouslySetInnerHTML={{ __html: tutorial.content }}
                  />
                </Paper>
              </motion.div>
            </Box>
          </Box>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default TutorialExpanded; 