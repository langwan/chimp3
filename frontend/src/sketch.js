let spectrumSize = 40;
let windowWidth = 800;
let windowHeight = 300;
let count = 0;
let r, g, b;
let fps = 20;
let orange, green;
const Y_AXIS = 1;
const X_AXIS = 2;
function setup() {
  console.log("setup");
  createCanvas(windowWidth, windowHeight);
  frameRate(20);
  orange = color(255, 161, 0);
  green = color(0, 228, 48);
}
function draw() {
  console.log("draw");
  background(0, 0, 0);
  if (count++ % fps == 0) {
    r = random(255);
    g = random(255);
    b = random(255);
  }
  fill(r, g, b);
  let columnWidth = windowWidth / spectrumSize;
  for (var i = 0; i < spectrumSize; i++) {
    let height = random(windowHeight);
    rect(columnWidth * i, windowHeight - height, columnWidth, height);
    // setGradient(
    //   columnWidth * i,
    //   windowHeight - height,
    //   columnWidth,
    //   height,

    //   orange,
    //   green,
    //   Y_AXIS
    // );
  }
}

function setGradient(x, y, w, h, c1, c2, axis) {
  noFill();

  if (axis === Y_AXIS) {
    // Top to bottom gradient
    for (let i = y; i <= y + h; i++) {
      let inter = map(i, y, y + h, 0, 1);
      let c = lerpColor(c1, c2, inter);
      stroke(c);
      line(x, i, x + w, i);
    }
  } else if (axis === X_AXIS) {
    // Left to right gradient
    for (let i = x; i <= x + w; i++) {
      let inter = map(i, x, x + w, 0, 1);
      let c = lerpColor(c1, c2, inter);
      stroke(c);
      line(i, y, i, y + h);
    }
  }
}
