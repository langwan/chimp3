import { Box, IconButton, Stack, Typography } from "@mui/material";
import {
  IconPlayerPause,
  IconPlayerPlay,
  IconPlayerSkipBack,
  IconPlayerSkipForward,
  IconPlus,
} from "@tabler/icons";
import { useEffect, useState } from "react";
import Sketch from "react-p5";
import io from "socket.io-client";
import { backendAxios } from "./axios";
let spectrumSize = 40;
let windowWidth = 320;
const maxColor = 200;
let windowHeight = (windowWidth / 16) * 9 + (((windowWidth / 16) * 9) % 40);
let count = 0;
let r, g, b;
let fps = 60;
let orange, green;
const Y_AXIS = 1;
const X_AXIS = 2;
let gridWidth = 10;
let gridColor = 0;
let density = 3;
let freqSpectrum = [];
const sio = io(
  process.env.NODE_ENV === "development" ? "ws://127.0.0.1:8000" : "/",
  {
    transports: ["websocket"],
    reconnect: true,
  }
);

let t = 0;
//rain
var drop = [];

export default (props) => {
  const [title, setTitle] = useState("CHIMP3");
  const [isPlay, setIsPlay] = useState(false);
  const [mode, setMode] = useState(0);
  useEffect(() => {
    sio.on("connect", () => {});
    sio.on("push", (message) => {
      if ("samples" in message) {
        freqSpectrum = message.samples;
      } else {
        freqSpectrum = [];
      }
      setIsPlay(message.is_player);
      setTitle(message.name);
      setMode(message.mode);
    });
  }, []);

  const setup = (p5, canvasParentRef) => {
    p5.createCanvas(windowWidth, windowHeight).parent(canvasParentRef);
    p5.frameRate(fps);
  };

  const draw = (p5) => {
    p5.background(255);
    let columnWidth = windowWidth / spectrumSize;
    gridWidth = columnWidth;
    if (mode != 2 && mode != 3) {
      for (let i = 0; i <= windowWidth / gridWidth; i++) {
        for (let j = 0; j <= windowHeight / gridWidth; j++) {
          p5.stroke(gridColor, 0, 0, 1);
          p5.line(i * gridWidth, 0, i * gridWidth, windowHeight);
          p5.line(0, j * gridWidth, windowWidth, j * gridWidth);
        }
      }
    }

    if (!freqSpectrum || freqSpectrum.length == 0) {
      return;
    }

    if (mode == 0 || mode == 1) {
      if (count++ % fps == 0) {
        r = p5.random(maxColor);
        g = p5.random(maxColor);
        b = p5.random(maxColor);
      }
    } else {
      r = 84;
      g = 86;
      b = 93;
    }

    for (var i = 0; i < spectrumSize; i++) {
      let height = freqSpectrum[0][i];
      p5.fill(r, g, b);
      p5.rect(columnWidth * i, windowHeight / 2 - height, columnWidth, height);
      p5.fill(r, g, b, 80);
      height = freqSpectrum[1][i];
      p5.rect(columnWidth * i, windowHeight / 2, columnWidth, height);
    }
  };

  const onNext = async (event) => {
    await backendAxios.post("/rpc/Next", {});
  };
  return (
    <Stack
      direction={"column"}
      justifyContent="space-between"
      alignItems="center"
    >
      <Stack
        direction={"row"}
        width={windowWidth}
        justifyContent="space-between"
        alignItems="center"
      >
        <Stack direction={"row"} justifyContent="flex-start">
          <IconButton>
            <IconPlus
              onClick={async (event) => {
                await backendAxios.post("/rpc/FileMulti", {});
              }}
              stroke={0.5}
            />
          </IconButton>
          <IconButton
            onClick={async (event) => {
              await backendAxios.post("/rpc/Prev", {});
            }}
          >
            <IconPlayerSkipBack stroke={0.5} />
          </IconButton>
          <IconButton
            onClick={async (event) => {
              console.log("isPlay", !isPlay);
              await backendAxios.post("/rpc/Playing", { is_play: !isPlay });
            }}
          >
            {isPlay ? (
              <IconPlayerPause stroke={0.5} />
            ) : (
              <IconPlayerPlay stroke={0.5} />
            )}
          </IconButton>
          <IconButton
            onClick={async (event) => {
              await backendAxios.post("/rpc/Next", {});
            }}
          >
            <IconPlayerSkipForward stroke={0.5} />
          </IconButton>
        </Stack>

        {title == "" ? (
          <Box
            component={"img"}
            sx={{ width: 32, height: 32 }}
            src={process.env.PUBLIC_URL + "/icon.png"}
          />
        ) : (
          <Typography
            variant="subtitle1"
            align="right"
            sx={{
              flexGrow: 1,
              overflow: "hidden",
              textOverflow: "ellipsis",
              whiteSpace: "nowrap",
            }}
          >
            {title}
          </Typography>
        )}
      </Stack>
      <Box
        onClick={async (event) => {
          await backendAxios.post("/rpc/Mode", {});
        }}
      >
        <Sketch setup={setup} draw={draw} />
      </Box>
    </Stack>
  );
};
