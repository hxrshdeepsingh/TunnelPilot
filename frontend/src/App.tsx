import { useEffect, useRef, useState } from "react";
import {
  CheckIsLoggedIn,
  KillAllTunnels,
  RunCommands,
  RunTunnelCreate,
  RunTunnelLogin,
  RunTunnelLogout,
  StopAnyPID,
} from "../bindings/changeme/tunnelservice";

import { Events } from "@wailsio/runtime";

interface ActiveTunnel {
  url: string;
  pid: number;
  port: string;
}

function App() {
  const [cloudflaredVersion, setCloudflaredVersion] = useState("---");
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [port, setPort] = useState("");
  const [logs, setLogs] = useState("");
  const [activeTunnels, setActiveTunnels] = useState<ActiveTunnel[]>([]);

  const currentPid = useRef(0);
  const currentPort = useRef("");

  useEffect(() => {
    checkCloudflaredVersion();
    checkLogin();

    // Wails v3 event listener
    const unsubscribe = Events.On("tunnel-log", (event) => {
      const line = String(event.data ?? "");

      handleTunnelLog(line);
    });

    return () => {
      unsubscribe();
    };
  }, []);

  async function checkCloudflaredVersion() {
    try {
      const result = await RunCommands("--version");

      setCloudflaredVersion(result);
    } catch (error) {
      appendLog(
        "Failed to get cloudflared version: " + String(error)
      );
    }
  }

  async function checkLogin() {
    try {
      const loggedIn = await CheckIsLoggedIn();

      setIsLoggedIn(loggedIn);
    } catch (error) {
      appendLog(
        "Failed to check login status: " + String(error)
      );
    }
  }

  function appendLog(text: string) {
    setLogs((previous) => previous + text + "\n");
  }

  function handleTunnelLog(line: string) {
    appendLog(line);

    const urls = extractUrls(line);

    if (urls.length > 0) {
      addActiveUrl(urls[0]);
    }
  }

  function extractUrls(text: string): string[] {
    const ignoreList = new Set([
      "https://www.cloudflare.com/website-terms/",
      "https://developers.cloudflare.com/cloudflare-one/connections/connect-apps",
    ]);

    const matches =
      text.match(/https?:\/\/[^\s<>"']+/g) || [];

    return matches
      .map((url) =>
        url.replace(/[)\],.!?;:"']+$/g, "")
      )
      .filter((url) => !ignoreList.has(url))
      .filter((url) => !url.includes("localhost"));
  }

  function addActiveUrl(url: string) {
    try {
      const newDomain = new URL(url).hostname;

      const exists = activeTunnels.some((item) => {
        try {
          return (
            new URL(item.url).hostname === newDomain
          );
        } catch {
          return false;
        }
      });

      if (exists) {
        return;
      }
    } catch {
      console.error("Invalid URL:", url);
      return;
    }

    setActiveTunnels((previous) => [
      ...previous,
      {
        url,
        pid: currentPid.current,
        port: currentPort.current,
      },
    ]);
  }

  async function handleLogin() {
    try {
      const result = await RunTunnelLogin();

      appendLog(result);

      await checkLogin();
    } catch (error) {
      appendLog(
        "Login failed: " + String(error)
      );
    }
  }

  async function handleLogout() {
    try {
      const result = await RunTunnelLogout();

      appendLog(result);

      await checkLogin();
    } catch (error) {
      appendLog(
        "Logout failed: " + String(error)
      );
    }
  }

  async function handleCreateTunnel() {
    const portNum = parseInt(port, 10);

    if (
      !port ||
      Number.isNaN(portNum) ||
      portNum < 1 ||
      portNum > 65535
    ) {
      appendLog(
        "Invalid port number. Please enter a port between 1 and 65535."
      );

      return;
    }

    const alreadyRunning = activeTunnels.some(
      (tunnel) =>
        String(tunnel.port) === String(port)
    );

    if (alreadyRunning) {
      appendLog(
        `A tunnel for port ${port} is already running.`
      );

      return;
    }

    appendLog(
      `Starting tunnel for port ${port}...`
    );

    try {
      const pid = await RunTunnelCreate(port);

      currentPid.current = pid;
      currentPort.current = port;

      appendLog(
        `Tunnel process started with PID ${pid}.`
      );
    } catch (error) {
      appendLog(
        "Failed to start tunnel: " + String(error)
      );
    }
  }

  async function handleStopTunnel(pid: number) {
    try {
      const success = await StopAnyPID(pid);

      if (success) {
        setActiveTunnels((previous) =>
          previous.filter(
            (tunnel) => tunnel.pid !== pid
          )
        );

        appendLog(
          `Tunnel ${pid} stopped.`
        );
      } else {
        appendLog(
          `Failed to stop tunnel ${pid}.`
        );
      }
    } catch (error) {
      appendLog(
        "Failed to stop tunnel: " + String(error)
      );
    }
  }

  async function handleKillAll() {
    appendLog(
      "Killing all active tunnels..."
    );

    try {
      const result = await KillAllTunnels();

      appendLog(result);

      setActiveTunnels([]);
    } catch (error) {
      appendLog(
        "Failed to kill tunnels: " + String(error)
      );
    }
  }

  async function copyUrl(url: string) {
    try {
      await navigator.clipboard.writeText(url);

      appendLog(
        "Copied URL: " + url
      );
    } catch {
      appendLog(
        "Failed to copy URL."
      );
    }
  }

  return (
    <div className="container">
      <div className="shell">

        {/* Header */}

        <header className="topbar">
          <div className="brand">
            <div>
              <div className="title">
                TunnelPilot
              </div>

              <div className="subtitle">
                Turn Localhost Into a Public Link in Seconds
              </div>
            </div>
          </div>

          <div className="meta">

            <div className="pill">
              <span className="dot" />

              <span>
                {cloudflaredVersion}
              </span>
            </div>

            {!isLoggedIn ? (
              <button
                className="btn primary"
                onClick={handleLogin}
              >
                Login
              </button>
            ) : (
              <button
                className="btn ghost"
                onClick={handleLogout}
                style={{
                  color: "var(--danger)",
                  borderColor:
                    "rgba(239, 68, 68, 0.45)",
                }}
              >
                Logout
              </button>
            )}

            <button
              className="btn ghost"
              title="Refresh Login Status"
              style={{
                padding: "10px 14px",
              }}
              onClick={checkLogin}
            >
              ↻
            </button>

          </div>
        </header>

        {/* Main */}

        <main className="grid">

          {/* Create Tunnel */}

          <section className="card">
            <div className="cardHeader">

              <div className="cardTitle">
                Create Tunnel
              </div>

              <div className="cardHint">
                Enter Port Number to Expose
              </div>

            </div>

            <div className="formRow">

              <label
                className="label"
                htmlFor="port"
              >
                Port Number
              </label>

              <div className="inputRow">

                <input
                  type="number"
                  id="port"
                  placeholder="3000"
                  min="1"
                  max="65535"
                  value={port}
                  onChange={(event) =>
                    setPort(event.target.value)
                  }
                />

                <button
                  className="btn primary"
                  onClick={handleCreateTunnel}
                >
                  Create
                </button>

              </div>

              <div className="help">
                Example: 3000, 5173, 8080
              </div>

            </div>
          </section>

          {/* Logs */}

          <section className="card">

            <div className="cardHeader">

              <div className="cardTitle">
                Logs
              </div>

              <div className="cardHint">
                Latest output from cloudflared
              </div>

            </div>

            <pre
              className="console"
              id="result"
            >
              {logs}
            </pre>

          </section>

          {/* Active Tunnels */}

          <section className="card full">

            <div
              className="cardHeader"
              style={{
                alignItems: "center",
              }}
            >

              <div
                style={{
                  display: "flex",
                  alignItems: "baseline",
                  gap: "10px",
                }}
              >

                <div className="cardTitle">
                  Active Tunnels
                </div>

                <div className="cardHint">
                  All Your Tunnels
                </div>

              </div>

              <button
                className="btn ghost"
                onClick={handleKillAll}
                style={{
                  color: "var(--danger)",
                  borderColor:
                    "rgba(239, 68, 68, 0.45)",
                  fontSize: "12px",
                  padding: "6px 10px",
                }}
              >
                Kill All Tunnels
              </button>

            </div>

            <div className="tableContainer">

              <table className="tunnelTable">

                <thead>
                  <tr>
                    <th>Local Port</th>
                    <th>Public URL</th>
                    <th>Action</th>
                  </tr>
                </thead>

                <tbody>

                  {activeTunnels.map((tunnel) => (
                    <tr
                      key={tunnel.pid}
                    >

                      <td>
                        localhost:{tunnel.port}
                      </td>

                      <td>

                        <a
                          href={tunnel.url}
                          target="_blank"
                          rel="noreferrer"
                          style={{
                            color: "var(--text)",
                          }}
                        >
                          {tunnel.url}
                        </a>

                        <span
                          style={{
                            cursor: "pointer",
                            color: "var(--primary)",
                            marginLeft: "5px",
                          }}
                          onClick={() =>
                            copyUrl(tunnel.url)
                          }
                        >
                          (copy!)
                        </span>

                      </td>

                      <td
                        style={{
                          textAlign: "right",
                        }}
                      >

                        <button
                          className="stopTunnelBtn"
                          onClick={() =>
                            handleStopTunnel(
                              tunnel.pid
                            )
                          }
                        >
                          Stop
                        </button>

                      </td>

                    </tr>
                  ))}

                </tbody>

              </table>

            </div>

          </section>

        </main>

        {/* Footer */}

        <footer
          className="footer"
          style={{
            textAlign: "center",
            marginTop: "20px",
            fontSize: "12px",
            color: "var(--muted)",
          }}
        >
          Built with love by{" "}

          <a
            href="https://github.com/hxrshdeepsingh"
            target="_blank"
            rel="noreferrer"
            style={{
              color: "var(--primary)",
              textDecoration: "none",
            }}
          >
            hxrshdeepsingh
          </a>
        </footer>

      </div>
    </div>
  );
}

export default App;