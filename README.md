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
  - `addendum list`: List all addendum files.
    ```bash
    ./bin/udccli addendum list
    ```
  - `addendum create [filename] [code] [title]`: Create a new addendum file.
    ```bash
    ./bin/udccli addendum create company 999.1 "Company Classifications"
    ```
  - `addendum add [filename] [code] [title]`: Add to existing addendum file.
    ```bash
    ./bin/udccli addendum add company 999.1.1 "Proprietary Equipment"
    ```
  - `addendum delete [filename]`: Delete an addendum file.
    ```bash
    ./bin/udccli addendum delete company
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

## CLI Usage

The `udccli` tool provides command-line access to UDC functionality:

```bash
# Scrape UDC data from the official website
./bin/udccli scrape

# Look up a classification by code
./bin/udccli lookup 621.3

# Look up a classification by title (fuzzy search)
./bin/udccli lookup "electrical engineering"

# Manage addendum files
./bin/udccli addendum list                    # List all addendum files
./bin/udccli addendum add 999.1 "Custom Code" # Add to default addendum
./bin/udccli addendum add 999.2 "Another" project # Add to specific addendum
./bin/udccli addendum delete project          # Delete an addendum file
```

### Addendum Management

Addendums allow you to add custom UDC classifications that are loaded alongside the base UDC data:

- **Default addendum**: If no filename is specified, classifications are added to `udc_addendum_default.yaml`
- **Custom addendums**: Specify a filename to create or modify a specific addendum file
- **Filename format**: Files are automatically prefixed with `udc_addendum_` and suffixed with `.yaml`
- **Validation**: Addendums cannot override existing UDC codes from the base classification

---

## Data & Configuration

- Place your data files (BOMs, UDC YAML, etc.) in the `data/` directory.
- The database file (`tags.db`) is used for persistent tag storage.

### UDC Addendum System

The platform supports local addendums to the UDC classification system:

- **`data/udc_full.yaml`**: The main UDC data file (updated only by scraping)
- **`data/udc_addendum_*.yaml`**: Local addendum files for custom classifications

#### Addendum Behavior

- **New codes only**: Addendum codes must be unique and cannot overlap with existing UDC codes
- **No overrides**: Addendums cannot modify or override existing UDC classifications
- **Children**: Addendum children are added as new codes under existing classifications
- **Multiple files**: All addendum files are loaded and merged
- **Validation**: The system will reject addendum files containing overlapping codes

#### Addendum File Format

Create files named `udc_addendum_*.yaml` in the `data/` directory:

```yaml
# Example: Adding new local classifications (valid)
- code: "621.3.001"
  title: "Local Electronics Classification"
  children:
    - code: "621.3.001.1"
      title: "Custom Circuit Design"

# Example: Adding completely new local classifications (valid)
- code: "999.1"
  title: "Local Company Classifications"
  children:
    - code: "999.1.1"
      title: "Proprietary Equipment"

# Example: Adding children to existing classifications (valid)
- code: "621.3.LOCAL"
  title: "Local Electrical Engineering Topics"
  children:
    - code: "621.3.LOCAL.1"
      title: "Company-Specific Electrical Systems"

# INVALID: Overriding existing UDC codes (will cause error)
# - code: "004"
#   title: "Computer Science and Technology (Local Enhancement)"
```

#### Best Practices

- Use descriptive filenames (e.g., `udc_addendum_company.yaml`)
- Keep addendums focused on specific domains or use cases
- Document the purpose of each addendum file
- Test addendums before deploying to production

---

## Contributing & Next Steps

This README is a work in progress. As the codebase evolves, please update usage instructions and setup steps. See the Makefile and Dockerfile for more advanced options.
