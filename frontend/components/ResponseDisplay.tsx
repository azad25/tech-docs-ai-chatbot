import * as React from "react";
import { useEffect, useRef } from "react";
import { Paper } from "@mui/material";
import { motion, AnimatePresence } from "framer-motion";

const LoadingSVG = () => (
  <svg width="64" height="64" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="32" cy="32" r="28" stroke="#1976d2" strokeWidth="4" strokeLinecap="round" strokeDasharray="44 44" strokeDashoffset="0">
      <animateTransform attributeName="transform" type="rotate" from="0 32 32" to="360 32 32" dur="1.5s" repeatCount="indefinite" />
    </circle>
  </svg>
);

interface ResponseDisplayProps {
  response: string;
  loading: boolean;
}

const ResponseDisplay: React.FC<ResponseDisplayProps> = ({ response, loading }) => {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0, behavior: 'smooth' });
    }
  }, [response]);

  return (
    <Paper elevation={3} sx={{ p: 3, mt: 3, minHeight: 300, overflowY: 'auto', fontFamily: 'Georgia, serif', fontSize: '1.1rem', lineHeight: 1.6 }} ref={containerRef}>
      <AnimatePresence>
        {loading ? (
          <motion.div
            key="loading"
            initial={{ opacity: 0 }}
            animate={{ rotate: 360, opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ repeat: Infinity, duration: 1.5, ease: "linear" }}
            style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}
          >
            <LoadingSVG />
          </motion.div>
        ) : (
          <motion.div
            key="response"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.4 }}
            dangerouslySetInnerHTML={{ __html: response }}
          />
        )}
      </AnimatePresence>
    </Paper>
  );
};

export default ResponseDisplay;