# Web Server Setup Guide

This guide provides a structured approach to setting up a web server using Go. The server will serve static files, provide a RESTful API endpoint, and include structured file organization for scalability. All necessary code files are included.

## Project Overview

- **Purpose** : Create a lightweight web server with static file serving and a JSON API.
- **Features** :
- Serve static HTML, CSS, and JavaScript files.
- Provide a `/api` endpoint returning a JSON response.
- Handle errors (e.g., 404) gracefully.
- Support environment-based configuration.
- **Structure** :

```
  my-web-server/
  ├── cmd/
  │   └── server/
  │       └── main.go
  ├── public/
  │   ├── index.html
  │   └── css/
  │       └── styles.css
  ├── handlers/
  │   └── handlers.go
  ├── .env
  ├── go.mod
  ├── go.sum
  └── Procfile
```

## Prerequisites

- Go (version 1.16 or higher)
- Basic knowledge of Go and command-line operations

## Setup Instructions

### 1. Initialize the Project

Create a project directory and initialize a Go module:

```bash
mkdir my-web-server
cd my-web-server
go mod init my-web-server
```

### 2. Install Dependencies

Add the `godotenv` package for environment variable management:

```bash
go get github.com/joho/godotenv
```

### 3. Create Project Structure

Set up the directory structure:

```bash
mkdir -p cmd/server public/css handlers
touch cmd/server/main.go public/index.html public/css/styles.css handlers/handlers.go .env Procfile
```

### 4. Configure Environment Variables

Create the `.env` file to store configuration:

```env
PORT=3000
```

PORT=3000
