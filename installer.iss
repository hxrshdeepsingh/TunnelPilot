#define MyAppName "Tunnelpilot"
#define MyAppVersion "1.0.0"
#define MyAppExeName "TunnelPilot.exe"

[Setup]
AppId={{D8D4D0E8-7C9F-4B95-9C0E-123456789ABC}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
OutputDir=installer
OutputBaseFilename=TunnelPilot-setup
Compression=lzma
SolidCompression=yes
SetupIconFile=build\appicon.ico

[Files]
Source: "bin\tunnelpilot.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\cloudflared.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#MyAppName}"; \
    Filename: "{app}\TunnelPilot.exe"; \
    IconFilename: "{app}\TunnelPilot.exe"

Name: "{autodesktop}\{#MyAppName}"; \
    Filename: "{app}\TunnelPilot.exe"; \
    IconFilename: "{app}\TunnelPilot.exe"; \
    Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "Create a desktop shortcut"

[Run]
Filename: "{app}\TunnelPilot.exe"; Description: "Launch {#MyAppName}"; Flags: nowait postinstall skipifsilent