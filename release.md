## required tools
- wails v3
- UPX (windows) / upx-ucl (linux)
- MSIX packaging tool
- wsl (ubuntu)

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
wails3 build GOOS=windows
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

## build for linux

go to wsl
```powershell
wsl -d Ubuntu
```

build the binary
```powershell
wails3 build GOOS=linux
```
compress the binary
```powershell
upx -9 bin/tunnelpilot
```

to run
```powershell
./bin/tunnelpilot
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


# TunnelPilot: Build and Publish the Snap

## 3. Build the snap

```bash
snapcraft clean
snapcraft pack
```

This produces `tunnelpilot_<version>_amd64.snap` (for example `tunnelpilot_3.0.0_amd64.snap`).

## 4. Test locally

```bash
sudo snap install ./tunnelpilot_3.0.0_amd64.snap --dangerous
tunnelpilot
```
Remove the local test install when done:

```bash
sudo snap remove tunnelpilot
```

## 5. Log in to the Snap Store

login to account using creds.txt

```bash
snapcraft login
```

if you dont have creds.txt file run this otherwise skip this step
```bash
snapcraft export-login creds.txt
export SNAPCRAFT_STORE_CREDENTIALS=$(cat creds.txt)
```

## 7. Upload and release

Upload to the `edge` channel first:

```bash
snapcraft upload --release=edge tunnelpilot_3.0.0_amd64.snap
```

Check the result:

```bash
snapcraft revisions tunnelpilot
snapcraft status tunnelpilot
```

Test it from the store:

```bash
sudo snap install tunnelpilot --edge
```

## 8. Promote to stable (public release)

Use the revision number shown by `snapcraft revisions`:

```bash
snapcraft release tunnelpilot <revision> stable
```

## Releasing a new version

1. Bump `version` in `snap/snapcraft.yaml`.
2. Replace `bin/tunnelpilot` with the new build.
3. Rebuild and upload:

```bash
snapcraft clean
snapcraft pack
snapcraft upload --release=edge tunnelpilot_<new-version>_amd64.snap
snapcraft release tunnelpilot <new-revision> stable
```