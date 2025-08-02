import * as React from "react";
import { styled, Box } from "@mui/material";

const FloatingInputBox = styled(Box)(({ theme }) => ({
  position: 'sticky',
  bottom: 0,
  left: 0,
  right: 0,
  backgroundColor: theme.palette.background.default,
  padding: theme.spacing(1, 2),
  boxShadow: '0 -2px 8px rgba(0,0,0,0.1)',
  zIndex: 10,
}));

export default FloatingInputBox;