# NAGL
A backend utility designed to handle logic for managing a simulated guild of adventurers (Fantasy HR).

# Tech Stack
* Language: Go (Golang)
* Database: Postgresql
* Framework: Chi
* Environment: Developed using Neovim and Wezterm

# Features
* Go-based backend: Leveraged Go's efficiency for core logic
* Modular Design: Organized packages into domains for easy scaling
* Dependency Management: Uses Go modules for reliable builds
* Postgresql database for persistence
* Authentication using JWT with ECDSA signing method

# Prerequisites
* Go (version 1.21 or higher recommended)
* A code editor (Neovim or VS Code)

# Installation & Setup
If you are having issues importing the code or running the project, follow tehse steps to ensure your environment is configured correctly
1. Clone the repository:
```git clone https://github.com/Barms1218/nagl.git```
```cd nagl```

2. Handle Dependencies:
If the go.mod is present, run the following to download the necessary dependencies: ```go mod tidy```.

3. Running the App
To start the service, run: ```go run main.go```.

# Roadmap
* Add goroutines to handle asynchronous operations.
* Add cron jobs to handle automated creation of adventurers and contracts.
* Enchanted Item Management.
