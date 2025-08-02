import { 
  Box, 
  Avatar, 
  Typography, 
  IconButton,
  Chip
} from "@mui/material";
import { motion } from "framer-motion";
import {
  Home as HomeIcon,
  Chat as ChatIcon,
  Search as SearchIcon,
  Psychology as AIIcon,
  School as LearnIcon,
  History as HistoryIcon,
  Settings as SettingsIcon,
  KeyboardArrowDown as ArrowDownIcon,
  AutoAwesome as SparkleIcon
} from "@mui/icons-material";

const Sidebar: React.FC = () => {
  return (
    <motion.div
      initial={{ x: -300, opacity: 0 }}
      animate={{ x: 0, opacity: 1 }}
      transition={{ 
        duration: 0.8, 
        ease: [0.25, 0.46, 0.45, 0.94],
        staggerChildren: 0.1
      }}
    >
      <Box
        sx={{
          width: 240,
          height: "100vh",
          background: "linear-gradient(180deg, rgba(255,255,255,0.95) 0%, rgba(248,249,250,0.95) 100%)",
          backdropFilter: "blur(20px)",
          borderRight: "1px solid rgba(0,0,0,0.1)",
          display: "flex",
          flexDirection: "column",
          position: "fixed",
          left: 0,
          top: 0,
          zIndex: 1200,
          boxShadow: "0 0 40px rgba(0,0,0,0.1)",
          overflow: "hidden",
          fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
        }}
      >
        {/* User Profile Section */}
        <motion.div
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          <Box sx={{ p: 3, pb: 2, borderBottom: "1px solid rgba(0,0,0,0.08)" }}>
            <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
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
              <Box sx={{ flexGrow: 1, minWidth: 0 }}>
                <Typography variant="h6" sx={{ fontWeight: 700, color: "#1D1D1F", mb: 0.5 }}>
                  ChatBot
                </Typography>
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
                    Online & Learning
                  </Typography>
                </Box>
              </Box>
              <motion.div
                whileHover={{ rotate: 180 }}
                transition={{ duration: 0.3 }}
              >
                <IconButton size="small">
                  <ArrowDownIcon sx={{ color: "#86868B" }} />
                </IconButton>
              </motion.div>
            </Box>
            
            {/* AI Status */}
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.4 }}
            >
              <Box sx={{ mt: 2, display: "flex", gap: 1, flexWrap: "wrap" }}>
                <Chip
                  icon={<SparkleIcon />}
                  label="RAG Enabled"
                  size="small"
                  sx={{
                    backgroundColor: "rgba(52,199,89,0.1)",
                    color: "#34C759",
                    fontSize: "0.7rem",
                    fontWeight: 600,
                    "& .MuiChip-icon": { fontSize: "0.8rem" }
                  }}
                />
                <Chip
                  icon={<AIIcon />}
                  label="Vector Search"
                  size="small"
                  sx={{
                    backgroundColor: "rgba(0,122,255,0.1)",
                    color: "#007AFF",
                    fontSize: "0.7rem",
                    fontWeight: 600,
                    "& .MuiChip-icon": { fontSize: "0.8rem" }
                  }}
                />
              </Box>
            </motion.div>
          </Box>
        </motion.div>

        {/* Navigation Menu */}
        <Box sx={{ flexGrow: 1, p: 2, overflowY: "auto" }}>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
            {[
              { text: "HOME", icon: HomeIcon, active: false },
              { text: "AI CHAT", icon: ChatIcon, active: true },
              { text: "SEARCH", icon: SearchIcon, active: false },
              { text: "LEARNING", icon: LearnIcon, active: false },
              { text: "HISTORY", icon: HistoryIcon, active: false },
              { text: "SETTINGS", icon: SettingsIcon, active: false },
            ].map((item, index) => (
              <motion.div
                key={item.text}
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
                  sx={{
                    p: 2,
                    borderRadius: 3,
                    cursor: "pointer",
                    background: item.active 
                      ? "linear-gradient(135deg, #34C759 0%, #30D158 100%)" 
                      : "transparent",
                    color: item.active ? "white" : "#1D1D1F",
                    display: "flex",
                    alignItems: "center",
                    gap: 2,
                    transition: "all 0.3s ease",
                    "&:hover": {
                      background: item.active 
                        ? "linear-gradient(135deg, #30D158 0%, #2ECC71 100%)" 
                        : "rgba(0,0,0,0.05)",
                      transform: "translateX(4px)",
                    },
                  }}
                >
                  <item.icon sx={{ fontSize: "1.2rem" }} />
                  <Typography 
                    variant="body2" 
                    sx={{ 
                      fontWeight: item.active ? 700 : 600,
                      fontSize: "0.9rem",
                      letterSpacing: "0.5px",
                    }}
                  >
                    {item.text}
                  </Typography>
                </Box>
              </motion.div>
            ))}
          </Box>
        </Box>

        {/* Footer */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 1.2 }}
        >
          <Box sx={{ p: 3, borderTop: "1px solid rgba(0,0,0,0.08)" }}>
            <Typography 
              variant="caption" 
              sx={{ 
                color: "#86868B", 
                fontSize: "0.7rem",
                fontWeight: 500,
                fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
              }}
            >
              Powered by llama3.2:1b & RAG
            </Typography>
          </Box>
        </motion.div>
      </Box>
    </motion.div>
  );
};

export default Sidebar; 