# Security Policy

## About TunnelPilot

TunnelPilot is a desktop application that allows users to expose local applications through a public URL.

Because tunnelling can make locally hosted services accessible over the internet, security vulnerabilities involving access control, tunnel exposure, sensitive information, or unintended network access should be reported responsibly.

## Supported Versions

Security fixes are prioritized for the latest stable release.

| Version | Supported |
|---|---|
| Latest stable release | Yes |
| Older releases | Best effort |

## Reporting a Vulnerability

If you discover a security vulnerability in TunnelPilot, please report it privately.

**Preferred reporting method:** Use GitHub's Private Vulnerability Reporting feature, if enabled for this repository.

Repository: https://github.com/hxrshdeepsingh/TunnelPilot/security/advisories/new

If private reporting is unavailable, contact the repository maintainer through a private communication channel. Do not publish sensitive vulnerability details in public issues.

When reporting a vulnerability, please include:

- A clear description of the issue.
- The affected version and operating system.
- Steps to reproduce the vulnerability.
- The potential security impact.
- A proof of concept or relevant logs, with secrets and personal information removed.

## Security Areas of Interest

Reports involving the following areas are especially welcome:

- Unauthorized access to exposed local services.
- Unexpected public exposure of local ports or applications.
- Authentication or access-control bypasses.
- Leakage of API keys, credentials, tokens, or other sensitive information.
- Insecure handling of tunnel URLs or configuration files.
- Remote code execution or other vulnerabilities affecting the host system.
- Vulnerabilities in dependencies used by TunnelPilot.

## Responsible Disclosure

Please allow the maintainer reasonable time to investigate and address a reported vulnerability before public disclosure.

Do not access, modify, or extract other people's data, disrupt services, or test vulnerabilities against systems without authorization.

## Scope

This policy applies to security vulnerabilities in TunnelPilot itself and its official distribution artifacts.

Third-party services and infrastructure are subject to their respective security policies.

## Acknowledgments

Security researchers who responsibly report valid vulnerabilities may be acknowledged after a fix is available, subject to their consent.