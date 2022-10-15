import { Stack } from "@mui/material";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import Spectrum from "./Spectrum";
const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <BrowserRouter basename={"/app"}>
    <Stack justifyContent={"center"} alignItems="center">
      <Spectrum />
    </Stack>
  </BrowserRouter>
);
