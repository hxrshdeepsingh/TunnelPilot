#define MyAppName "Tunnelpilot"
#define MyAppVersion "3.0.0.0"
#define MyAppExeName "tunnelpilot.exe"

[Setup]
AppId={{D8D4D0E8-7C9F-4B95-9C0E-123456789ABC}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
DefaultDirName={autopf}{#MyAppName}
DefaultGroupName={#MyAppName}
OutputDir=installer
OutputBaseFilename=TunnelPilot-setup
Compression=lzma
SolidCompression=yes
SetupIconFile=build\appicon.ico

[Files]
Source: "bin\tunnelpilot.exe"; DestDir: "{app}"; Flags: ignoreversion

[Tasks]
Name: "desktopicon"; Description: "Create a desktop shortcut"; Flags: unchecked

[Icons]
Name: "{autoprograms}\Tunnelpilot"; Filename: "{app}\tunnelpilot.exe"; IconFilename: "{app}\tunnelpilot.exe"
Name: "{autodesktop}\Tunnelpilot"; Filename: "{app}\tunnelpilot.exe"; IconFilename: "{app}\tunnelpilot.exe"; Tasks: desktopicon

[Run]
Filename: "{app}\tunnelpilot.exe"; Description: "Launch Tunnelpilot"; Flags: nowait postinstall skipifsilent
