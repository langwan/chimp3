const { app, BrowserWindow } = require("electron");

app.whenReady().then(() => {
  const win = new BrowserWindow({
    title: "chimp3 v2",
    width: 360,
    height: 300,
    resizable: false,
  });
  win.loadURL("http://localhost:3000/");
});
