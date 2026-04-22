# NAGL
**Fantasy HR:** A backend utility designed to manage a simulated guild of adventurers.

## Tech Stack
* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Framework:** Chi Router
* **SQL Engine:** SQLC
* **Migrations:** Goose
* **Auth**: JWT with ECDSA signing
* **Environment:** Developed using Neovim and Wezterm

## Features
* **Domain-Driven Design:** Organized into modular domains for scalable backend architecture.
* **Dependency Management:** Using App pattern with Chi Router to inject services into handlers.
* **Persistent Storage:** Full PostgreSQL integration for guild data, adventurers, and contracts.
* **JWT Authentication:** Secure access using ECDSA-signed tokens
* **Procedural Generation:** Anthropic SDK for procedural generation of adventurers, contracts, and party names.
* **Complex Simulation Logic:** Manages party and individual progression through ACID-compliant PostgreSQL transactions.
* **Dynamic Contract Lifestyle:** Handles the full flow from procedural generation and claiming to completion or failure with "Party Fate" calculations.
* **Advanced Filtering:** Flexible listing service for adventurers and contracts using multi-parameter search filters.

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

## Engineering Decisions
* PostgreSQL was chosen for its flexibility and speed. The ability to create type-safe enums was also a big factor.
* Chi was selected because the Go standard library is already excellent. Chi was used because it doesn't interfere with that strength.
* SQLC was chosen for the ability to write raw sql but work with type-safe Go code.
* Goose was chosen for its simplicity in handling migrations, and the encouragement to write numbered sql files.
* The Anthropic SDK was chosen for procedural generation due to Claude's strength in creative writing.
* Manually mapped database models to response DTOs to ensure strict separation of concerns, preventing accidental data leaks.

**Example of a prompt to Anthropic SDK**
```
	systemPrompt := fmt.Sprintf(`You are a fantasy adventurer record keeper. Generate an adventurer profile as JSON. Respond with ONLY valid JSON, no markdown, no explanation. Use this exact shape:
	{
		"name": "string",
		"role": "frontliner" | "spellcaster" | "healer" | "generalist",
		"current_rank": 1-5,
		"description": "string (2-3 sentences of flavor text)",
		"upkeep_cost": 10-100,
		"recruitment_cost": 50-(100 * current_rank)
	}
	The following adventurers already exist. Do not repeat their names, and avoid generating a duplicate combination of role and rank: %s`, string(exclusionJSON))
```

## Roadmap
* [ ] **Seed Database:** Add starting adventurers and contracts.
* [ ] **Asynchronous Operations:** Implement Goroutines for background task processing.
* [ ] **Automated Guild Growth:** Add Cron jobs for automated creation of adventurers and contracts.
* [ ] **Enchanted Item Management:** Expand the schema to support magical inventory and equipment, adding more depth to guild economy.
* [ ] **Artificer Service:** Add a guild artificer that can recharge enchanted items.
