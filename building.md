## 1. Prerequisites

Make sure these commands work:

```powershell
wails3 version
go version
npm --version
upx --version
garble version
```
---

### 2. Recommended Obfuscated Build

For a release build, use:

```powershell
wails3 build -obfuscated
```

### 3. Compress the EXE With UPX

Use:

```powershell
upx --best --lzma .\tunnelpilotv3.exe
```

Explanation:

- `--best` asks UPX to use its strongest compression level.
- `--lzma` uses LZMA compression.
- `.\tunnelpilotv3.exe` is the executable being compressed.

### 4. Verify the UPX Compression

Always test the packed executable:

```powershell
upx -t .\tunnelpilotv3.exe
```