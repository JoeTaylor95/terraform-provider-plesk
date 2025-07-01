# Terraform Provider for Plesk

This Terraform provider allows you to manage Plesk resources such as domains, FTP accounts, clients, resellers, and mailboxes through the Plesk REST API.

---

## Features

- Create, read, update, and delete Plesk **domains**
- Manage **FTP accounts**
- Manage Plesk **clients/users**
- Manage **resellers**
- Manage **mailboxes**
- (More resources coming soon!)

---

## Prerequisites

- A Plesk server with the REST API enabled and accessible.
- An API token with appropriate permissions.
- Go 1.18+ to build the provider locally (optional).
- Terraform 1.0+ to use the provider.

---

## Installation

### Build from source

```bash
git clone https://taylor.am/terraform-git.home.provider/terraform-provider-plesk.git
cd plesk/provider
go build -o terraform-provider-plesk
