# SoundPort CLI

## Overview

SoundPort is a command-line tool that enables users to transfer playlists from Spotify to YouTube Music seamlessly.

## Prerequisites

- Go (version 1.23.0 or higher)
- Spotify Developer Account
- Spotify Account (Free / Premium)
- YouTube Music Account (Free / Premium)

## Installation

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
   - Add a Redirect of `http://127.0.0.1:4214/callback`
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

## Youtube Music Setup

### The API Limitation
YouTube Music lacks an official public API for programmatic playlist management. This presents a unique authentication challenge for music transfer applications like SoundPort.

### Authentication Approach: Browser Request Mimicry
To overcome the API limitation, SoundPort mimics legitimate browser requests.

### How It Works
* The application requires a session cookie obtained directly from a web browser.
* This cookie serves as a cryptographic passport, authenticating requests to YouTube Music.
* It enables the application to interact with your personal YouTube Music account as if a human were making the requests.

### Obtaining Your YouTube Music Cookie.

- Go to  [Youtube Music homepage](music.youtube.com)
- Open developer tools and select the “Network” tab.
- To find an Authenticated `POST` request, filer the requests by `/browse` using the search bar.
![](README/Screenshot%202025-03-25%20at%206.40.48%E2%80%AFPM.jpg)

- Refresh the page to find the `POST` request. It should look something like this. If you can’t find it, click the `Library` tab on the sidebar.
![](README/Screenshot%202025-03-25%20at%206.45.03%E2%80%AFPM.jpg)

- Click on the request, Scroll till you find `Request Headers` section.
- Copy the `Cookie` property of the request header. Copy everything from 	`__Secure-ROLLOUT_TOKEN` to the end.
![](README/Screenshot%202025-03-25%20at%206.47.44%E2%80%AFPM.jpg)

### Configure in SoundPort

- Use `soundport ytmusic setup`.
- Paste the copied cookie in the textarea and press `Enter`.

> [!CAUTION]
> **Important Warning**
> * Never share your cookie with anyone.
> * Cookies can provide significant account access.
> * The cookie saved in `SoundPort` is saved on your local system and secure.

