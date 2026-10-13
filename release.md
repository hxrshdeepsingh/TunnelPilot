## required tools
- wails v3
- UPX
- MSIX packaging tool

```text
App Name:              TunnelPilot
Package Identity Name: Harshdeepsingh.TunnelPilot
Publisher:             CN=ECC6BB33-156A-4CA5-AB9D-90DA53CEAF89
Publisher Display:     Harshdeepsingh
Store ID:              9NMFWT6PHT8T
Package Family Name:   Harshdeepsingh.TunnelPilot_mjaedyhsnsee8
```

## 0. Update version
Edit version in build/config.yml
Run: 
```powershell
wails3 task common:update:build-assets
```

## 1. Build Wails App

```powershell
wails3 build
```

Check the generated `.exe` in `bin`.

## 2. Compress TunnelPilot EXE

```powershell
upx --best --lzma .\bin\tunnelpilot.exe
```

Verify:

```powershell
upx -t .\bin\tunnelpilot.exe
```

## 3. Create MSIX
```powershell
wails3 tool msix `
  --config ".\build\config.yml" `
  --name "TunnelPilot.exe" `
  --executable ".\bin\tunnelpilot.exe" `
  --arch "amd64" `
  --publisher "CN=ECC6BB33-156A-4CA5-AB9D-90DA53CEAF89" `
  --out ".\bin\tunnelpilot-amd64.msix" `
  --use-makeappx
  ```

```text
Package Name: Harshdeepsingh.TunnelPilot
Package Display Name: TunnelPilot
Publisher: CN=ECC6BB33-156A-4CA5-AB9D-90DA53CEAF89
Publisher Display Name: Harshdeepsingh
```

Set the package version to a version **higher than the currently published Store version**.

## 4. unpack the generated msix file and replace the assets folder with the assets folder we have in our dir

## 5. Update Microsoft Store

1. Open **Microsoft Partner Center**.
2. Open **TunnelPilot**.
3. Click **Start update**.
4. Open **Packages**.
5. Upload the new `.msix` / `.msixupload`.
6. Check identity and version.
7. Update **What's new in this version**.
8. Save the submission.
9. Click **Submit for certification**.

## 6. Quick Release Flow

```text
Wails build
    ↓
UPX TunnelPilot.exe
    ↓
UPX cloudflared.exe
    ↓
Test
    ↓
Create MSIX
    ↓
Check identity + higher version
    ↓
Partner Center → Start update
    ↓
Packages → Upload
    ↓
What's new
    ↓
Submit for certification
```