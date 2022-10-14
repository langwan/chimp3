import { useEffect } from "react";
import Sketch from "react-p5";
import io from "socket.io-client";
let spectrumSize = 40;
let windowWidth = 320;
let windowHeight = (320 / 16) * 9 + (((320 / 16) * 9) % 40);
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
const sio = io("ws://localhost:8000", {
  transports: ["websocket"],
  reconnect: true,
});

let t = 0;
//rain
var drop = [];

export default (props) => {
  useEffect(() => {
    sio.on("connect", () => {
      console.log("socketio connect", sio.id);
    });
    sio.on("push", (message) => {
      freqSpectrum = message;
    });
  }, []);

  const setup = (p5, canvasParentRef) => {
    p5.createCanvas(windowWidth, windowHeight).parent(canvasParentRef);
    p5.frameRate(fps);
    orange = p5.color(255, 161, 0);
    green = p5.color(0, 228, 48);
    // for (var i = 0; i < 200; i++) {
    //   drop[i] = new Drop(p5);
    // }
    p5.noStroke();
    p5.fill(0);
    // p5.pixelDensity(3 / 7); //important
  };
  const draw = (p5) => {
    p5.background(255);

    // for (var i = 0; i < 20; i++) {
    //   drop[i].show();
    //   drop[i].update();
    // }

    if (count++ % fps == 0) {
      r = p5.random(255);
      g = p5.random(255);
      b = p5.random(255);
    }

    let columnWidth = windowWidth / spectrumSize;
    for (var i = 0; i < spectrumSize; i++) {
      let height = freqSpectrum[i];
      // p5.stroke(23, 23, 23);
      p5.fill(r, g, b, 50);
      p5.rect(columnWidth * i, windowHeight / 2 - height, columnWidth, height);
      p5.fill(Math.max(0, r - 80), Math.max(g - 80), Math.max(b - 80));
      p5.rect(columnWidth * i, windowHeight / 2, columnWidth, height);
    }
    //p5.scale(7);
    // for (var i = 0; i < 100; i++) {
    //   for (var j = 0; j < 50; j++) {
    //     p5.rect(
    //       i,
    //       j,
    //       ((p5.sin(j * j + i / j - t * 7) +
    //         p5.cos(j ** 5 - (i / j) * 6 + t) ** 3) *
    //         j) /
    //         50,
    //       1
    //     );
    //   }
    // }

    // t += 0.01;
    gridWidth = columnWidth;
    for (let i = 0; i <= windowWidth / gridWidth; i++) {
      for (let j = 0; j <= windowHeight / gridWidth; j++) {
        p5.stroke(gridColor, 0, 0, 1);
        p5.line(i * gridWidth, 0, i * gridWidth, windowHeight);
        p5.line(0, j * gridWidth, windowWidth, j * gridWidth);
      }
    }
  };
  return <Sketch setup={setup} draw={draw} />;
};

function Drop(p5) {
  p5.noStroke();
  this.x = p5.random(0, windowWidth);
  this.y = p5.random(0, -windowHeight);

  this.show = function () {
    p5.fill(0);
    p5.ellipse(this.x, this.y, p5.random(1, 5), p5.random(1, 5));
  };
  this.update = function () {
    this.speed = p5.random(5, 10);
    this.gravity = 1.05;
    this.y = this.y + this.speed * this.gravity;

    if (this.y > windowHeight) {
      this.y = p5.random(0, -windowHeight);
      this.gravity = 0;
    }
  };
}
