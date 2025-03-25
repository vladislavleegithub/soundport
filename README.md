# SoundPort CLI

## Overview

SoundPort is a command-line tool that enables users to transfer playlists from Spotify to YouTube Music seamlessly.

## Prerequisites

- Go (version 1.23.0 or higher)
- Spotify Developer Account
- Spotify Account (Free / Premium)
- YouTube Music Account (Free / Premium)

## Installation

### Requirements

- Golang 1.23.0+
- Git

### Install from Source

```bash
# Clone the repository
git clone https://github.com/Samarthbhat52/soundport.git

# Navigate to the project directory
cd soundport

# Build the application
go build -o soundport

# Optional: Install globally
go install
```

## Spotify setup

### How to Obtain Spotify Developer Credentials

1. **Create a Spotify Developer Account**

   - Visit [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)
   - Log in with your Spotify account

2. **Register a New Application**

   - Click "Create App"
   - Give a name and description for `SoundPort`
   - Add a Redirect of `http://localhost:8080/callback`
     - Make sure the URL matches exactly.
     - This is where Spotify will send authentication responses
   - Select `Web API` as api/sdk option.
   - Agree to Spotify's Developer Terms of Service
   - Save changes

3. **Retrieve Credentials**

   - On the Spotify Developer Dashboard, select your project.
   - Click "Settings" to view your app's information.
   - You'll find:
     - Client ID
     - `View client secret` link.
     - Click the `View client secret` link and copy both Client ID and Client Secret.
   - Keep these credentials secure and do not share publicly

### Configure in SoundPort

- Use `soundport spotify setup` to input your credentials.
- The tool will store and manage these tokens on your system.
