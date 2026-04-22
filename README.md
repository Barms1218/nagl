# NAGL
**Fantasy HR:** A backend utility designed to manage a simulated guild of adventurers.

## Tech Stack
* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Framework:** Chi Router
* **Auth**: JWT with ECDSA signing
* **Environment:** Developed using Neovim and Wezterm

## Features
* **Domain-Driven Design:** Organized into modular domains for scalable backend architecture.
* **Dependency Management:** Using App pattern with Chi Router to inject services into handlers.
* **Persistent Storage:** Full PostgreSQL integration for guild data, adventurers, and contracts.
* **JWT Authentication:** Secure access using ECDSA-signed tokens
* **Procedural Generation:** Anthropic SDK for procedural generation of adventurers, contracts, and party names.

## Prerequisites
* **Go:** (version 1.21 or higher recommended)
* **PostgreSQL:** A running instance for data persistence

## Installation & Setup
1. **Clone and Enter the repository:**
```
git clone https://github.com/Barms1218/nagl.git
```
```
cd nagl
```

2. **Configure Enironment Variables:**
   ```
   PORT=8080
   DB_URL=postgres://user:password@localhost:5432/nagl
   ANTHROPIC_API_KEY=your_key_here
   PRIVATE_KEY_PATH=path to your private key file
   ```
   
3. **Set Up Security Keys:**
   This project uses ECDSA for JWT signing. Ensure you have your private key file in the root directory:
   ```
   # Generate the private key
   openssl ecparam -name prime256v1 -genkey -noout -out ec_private.pem

   # Derive the public key
   openssl ec -in ec_private.pem -pubout -out ec_public.pem
   ```
   **Note:** Ensure that your private key file is added to your ```.gitignore``` to keep your private key secure.
  
5. **Handle Dependencies:**
```
go mod tidy
```

5. **Running the App:**
```
go run main.go
```
**Note:** This project requires an `ANTHROPIC_API_KEY` set in your environment variables to enable procedural generation.

## Roadmap
* [ ] **Seed Database:** Add starting adventurers and contracts.
* [ ] **Asynchronous Operations:** Implement Goroutines for background task processing.
* [ ] **Automated Guild Growth:** Add Cron jobs for automated creation of adventurers and contracts.
* [ ] **Enchanted Item Management:** Expand the schema to support magical inventory and equipment.
