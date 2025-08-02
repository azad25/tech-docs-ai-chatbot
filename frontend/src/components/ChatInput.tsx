import * as React from "react";
import { useState } from "react";
import { Box, TextField, IconButton } from "@mui/material";
import SendIcon from '@mui/icons-material/Send';

interface ChatInputProps {
  onSend: (msg: string) => void;
  disabled: boolean;
}

const ChatInput: React.FC<ChatInputProps> = ({ onSend, disabled }) => {
  const [input, setInput] = useState<string>("");

  const handleSend = () => {
    if (input.trim()) {
      onSend(input.trim());
      setInput("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', p: 1, bgcolor: 'background.paper', boxShadow: 3, borderRadius: 2 }}>
      <TextField
        multiline
        maxRows={4}
        fullWidth
        placeholder="Ask TechDocs AI..."
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        variant="outlined"
        size="small"
      />
      <IconButton color="primary" onClick={handleSend} disabled={disabled || !input.trim()} aria-label="send">
        <SendIcon />
      </IconButton>
    </Box>
  );
};

export default ChatInput;