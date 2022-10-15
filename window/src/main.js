const { app, BrowserWindow, ipcMain } = require("electron");

const freeport = require("freeport");
const fetch = require("electron-fetch").default;

const { execFile } = require("child_process");
let backendPort = 8000;
let win = null;
app.commandLine.appendSwitch("disable-web-security");
app.whenReady().then(() => {
  freeport(function (err, port) {
    backendPort = port;
    let backendBin = "./chimp3_backend";
    if (process.platform == "darwin") {
      backendBin = "./chimp3_backend";
    } else {
      backendBin = "./chimp3_backend.exe";
    }

    try {
      const child = execFile(backendBin, ["--port", port], {
        cwd: __dirname + "/bin",
        env: { langgo_env: "production", PATH: process.env.PATH },
      });
      let url = "http://127.0.0.1:" + port + "/app";

      child.stdout.on("data", (data) => {
        if (win != null) {
          return;
        }
        let additionalArguments = [
          "NODE_TLS_REJECT_UNAUTHORIZED=0",
          "REACT_APP_BACKEND_HOST=" + "localhost:" + port,
        ];
        console.log(additionalArguments);
        win = new BrowserWindow({
          title: "chimp3 v2",
          width: 360,
          height: 300,
          maximizable: false,
          resizable: process.env.NODE_ENV == "development" ? true : false,
          maximizable: true,
          webPreferences: {
            nodeIntegration: true,
            enableRemoteModule: true,
            additionalArguments: additionalArguments,
            webSecurity: false,
          },
        });
        setTimeout(() => {
          win.loadURL(url, {
            userAgent: "App",
          });
          if (process.env.NODE_ENV == "development")
            win.webContents.openDevTools();
        }, 1000);
      });
    } catch (e) {
      console.log(e);
    }
  });
});

app.on("before-quit", function () {
  if (process.env.NODE_ENV != "development") {
    new Promise((resolve, reject) => {
      fetch("http://127.0.0.1:" + backendPort + "/rpc/Quit", {
        method: "get",
      }).then(function (response) {
        console.log(response);
        resolve(response);
      }),
        (error) => {
          reject(new Error(error.message));
        };
    });
  }
});
