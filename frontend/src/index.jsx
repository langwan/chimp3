import { Box, Button, Stack } from "@mui/material";
import ReactDOM from "react-dom/client";
import Spectrum from "./Spectrum";
const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <Stack>
    <Box>
      <Button>选择文件</Button>
    </Box>
    <Spectrum />
  </Stack>
);
