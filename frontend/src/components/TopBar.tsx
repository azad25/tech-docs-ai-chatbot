import { AppBar, Toolbar, Typography } from "@mui/material";

const TopBar = () => (
  <AppBar position="sticky" color="default" elevation={1}>
    <Toolbar>
      <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 600, fontFamily: 'Georgia, serif' }}>
        TechDocs AI
      </Typography>
    </Toolbar>
  </AppBar>
);

export default TopBar;