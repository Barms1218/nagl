# NAGL

**Fantasy HR:** A backend utility for managing a simulated guild of adventurers. Architecturally mirrors real-world contractor management systems — adventurers map to contractors, contracts to jobs, and guilds to organizations.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go |
| Database | PostgreSQL |
| Router | Chi |
| Query Generation | SQLC |
| Migrations | Goose |
| Auth | JWT with ECDSA (ES256) signing |
| LLM Integration | Anthropic Go SDK |
| Environment | Neovim, Wezterm |

## Features

- **Domain-Driven Design:** Organized into modular domains (guild, adventurer, contract, party) for scalable backend architecture.
- **Dependency Injection:** App pattern with Chi mounts domain routers and injects services into handlers via closures.
- **JWT Authentication:** Secure access using ECDSA-signed tokens; middleware extracts guild identity from verified claims.
- **Procedural Generation:** Anthropic SDK drives background creation of adventurers, contracts, and party names on a cron schedule.
- **Contract Resolution:** Probabilistic outcome engine weighing party rank (50%) and individual member ranks (50% split across party), clamped to [0, 1] and resolved concurrently via background cron job.
- **Rank Progression:** Party and individual adventurer ranks advance through ACID-compliant transactions — rank-up triggers on every fifth completed contract at or above current rank.
- **Advanced Filtering:** Multi-parameter search for adventurers and contracts using optional SQL filters without dynamic query building.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL (running instance)
- Anthropic API key

## Installation & Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Barms1218/nagl.git
   cd nagl
   ```

2. **Configure environment variables:**

   ```bash
   PORT=8080
   DB_URL=postgres://user:password@localhost:5432/nagl
   ANTHROPIC_API_KEY=your_key_here
   ```

3. **Generate security keys:**

   This project uses ECDSA for JWT signing. Generate your key pair and keep the private key out of version control.

   ```bash
   # Generate the private key
   openssl ecparam -name prime256v1 -genkey -noout -out ec_private.pem

   # Derive the public key
   openssl ec -in ec_private.pem -pubout -out ec_public.pem
   ```

   > **Note:** `ec_private.pem` is already included in `.gitignore`.

4. **Install dependencies:**

   ```bash
   go mod tidy
   ```

5. **Run migrations:**

   ```bash
   goose -dir internal/sql/schema postgres "$DB_URL" up
   ```

6. **Start the server:**

   ```bash
   go run cmd/api/main.go
   ```

## Engineering Decisions

**PostgreSQL** was chosen for its reliability and the ability to define type-safe enums at the schema level, which enforces domain constraints without application-layer validation overhead.

**Chi** was selected because the Go standard library HTTP primitives are already strong. Chi adds routing and middleware without obscuring them.

**SQLC** enables writing raw SQL while getting type-safe Go structs generated from it — no ORM abstraction, no runtime query building, full control over what hits the database.

**Goose** handles migrations through numbered SQL files, making the schema history explicit and easy to follow.

**Manual DTO mapping** between database models and response types enforces strict separation of concerns and prevents accidental data leaks (e.g. password hashes) at the serialization boundary.

**Anthropic SDK** drives procedural generation because creative writing is a genuine strength of the model. Each generated adventurer and contract is checked against existing records before insertion to avoid duplicates and a structured JSON response is returned..

**Example prompt (adventurer generation):**

```go
systemPrompt := fmt.Sprintf(`You are a fantasy adventurer record keeper. Generate an adventurer profile as JSON.

The following adventurers already exist. Do not repeat their names,
and avoid generating a duplicate combination of role and rank: %s`, string(exclusionJSON))
```

## Roadmap

- [ ] Seed database with starting adventurers and contracts
- [ ] Enchanted item management — expand schema to support magical inventory and guild economy
- [ ] Artificer service — guild NPC that recharges enchanted items
