import { IconButton, Stack, Typography } from "@mui/material";
import {
  IconPlayerPlay,
  IconPlayerSkipBack,
  IconPlayerSkipForward,
  IconPlus,
} from "@tabler/icons";
import ReactDOM from "react-dom/client";
import Spectrum from "./Spectrum";
const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <Stack justifyContent={"center"} alignItems="center">
    <Stack
      direction={"column"}
      justifyContent="space-between"
      alignItems="center"
    >
      <Stack
        direction={"row"}
        width="100%"
        justifyContent="space-between"
        alignItems="center"
      >
        <Stack direction={"row"}>
          <IconButton>
            <IconPlus stroke={0.5} />
          </IconButton>
          <IconButton>
            <IconPlayerSkipBack stroke={0.5} />
          </IconButton>
          <IconButton>
            <IconPlayerPlay stroke={0.5} />
          </IconButton>
          <IconButton>
            <IconPlayerSkipForward stroke={0.5} />
          </IconButton>
        </Stack>
        <Typography align="right" sx={{ flexGrow: 1 }}>
          吻别
        </Typography>
      </Stack>
      <Spectrum />
    </Stack>
  </Stack>
);
