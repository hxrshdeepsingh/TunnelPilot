```text
App Name:              TunnelPilot
Package Identity Name: Harshdeepsingh.TunnelPilot
Publisher:             CN=ECC6BB33-156A-4CA5-AB9D-90DA53CEAF89
Publisher Display:     Harshdeepsingh
Store ID:              9NMFWT6PHT8T
Package Family Name:   Harshdeepsingh.TunnelPilot_mjaedyhsnsee8
```

## 1. Build Wails App

```powershell
wails3 build
```

Check the generated `.exe` in `bin`.

## 2. Compress TunnelPilot EXE

```powershell
upx --best --lzma .\tunnelpilotv3.exe
```

Verify:

```powershell
upx -t .\tunnelpilotv3.exe
```

Expected:

```text
[OK]
```

Run the app and test it.

## 3. Create MSIX

Use **Microsoft MSIX Packaging Tool** to convert the Windows installer into an `.msix`.

The MSIX must keep these Store identity values:

```text
Name:      Harshdeepsingh.TunnelPilot
Publisher: CN=ECC6BB33-156A-4CA5-AB9D-90DA53CEAF89
Display:   Harshdeepsingh
```

Set the package version to a version **higher than the currently published Store version**.

Example:

```text
1.0.0.0 → 1.0.1.0
```

Do not change the package identity.

## 4. Update Microsoft Store

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