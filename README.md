# UDC Codec: Industrial Asset Tagging & Classification Platform

🗂 **Universal Decimal Classification (UDC) for Industrial Automation Systems**

## Project Purpose

This project is a toolkit and platform for organizing, tagging, automatically classifying, and understanding both new and legacy industrial automation systems and their components—especially equipment and instruments. It leverages the Universal Decimal Classification (UDC) system and related standards to help engineers, asset managers, and organizations structure, migrate, and analyze their industrial asset data.

## Features

- **Automatic Classification**: Classify equipment and instruments using UDC and other standards (e.g., IEC 81346).
- **Tagging & Asset Management**: Tag and manage assets with support for both manual and automated workflows.
- **Import/Export Pipelines**: Import data from various sources and export to standardized formats.
- **Data Aggregation & Validation**: Aggregate and validate data for quality and compliance.
- **Web & CLI Interfaces**: Interact via a web portal or command-line tools.
- **Database-Backed**: Persistent storage for projects, assets, and tags.

---

## Setup

### Prerequisites

- Go 1.20+ (for building from source)
- Docker (optional, for containerized deployment)
- Make

### Building and Running (Locally)

```bash
# Build all binaries
make build

# Run the REST API server
make run-server

# Run the CLI tool
make run-cli

# Run the web portal (if separate)
go build -o bin/webserver ./cmd/webserver
./bin/webserver
```

### Docker

To build and run everything in Docker:

```bash
docker-compose up --build
```

This will expose the web portal and API on port 8080 by default.

---

## Usage

### Command-Line Tools

- **udccli**  
  - `scrape`: Recursively scrape the UDC summary and save to `data/udc_full.yaml`.
    ```bash
    ./bin/udccli scrape
    ```
  - `lookup [code]`: Lookup a UDC code in the local database.
    ```bash
    ./bin/udccli lookup 621.3
    ```

- **server**  
  REST API for tag management.
  - Start:  
    ```bash
    ./bin/server
    ```
  - Endpoints:
    - `GET /tags/{tag}`: Lookup a tag.
    - `POST /tags`: Insert a new tag (JSON body).

- **webserver**  
  Web portal for browsing, uploading BOM files, and managing tags.
  - Start:  
    ```bash
    ./bin/webserver
    ```
  - Access via browser at [http://localhost:8080](http://localhost:8080)

- **autopipeline**  
  Automated pipeline for validating and exporting tag lists from BOM and aggregated data.
  - Start:  
    ```bash
    ./bin/autopipeline
    ```
  - Expects input files in `data/` (e.g., `project_bom.yaml`, `aggregated_master.yaml`).

---

## Data & Configuration

- Place your data files (BOMs, UDC YAML, etc.) in the `data/` directory.
- The database file (`tags.db`) is used for persistent tag storage.

---

## Contributing & Next Steps

This README is a work in progress. As the codebase evolves, please update usage instructions and setup steps. See the Makefile and Dockerfile for more advanced options.
