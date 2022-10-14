import { Stack } from "@mui/material";
import ReactDOM from "react-dom/client";
import Spectrum from "./Spectrum";
const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <Stack justifyContent={"center"} alignItems="center">
    <Spectrum />
  </Stack>
);
