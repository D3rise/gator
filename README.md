# Gator - RSS Feed Aggregator CLI

Gator is a command-line RSS feed aggregator written in Go. It allows users to manage, follow, and browse RSS feeds from the terminal.

## Features

- User management (register, login)
- Add and manage RSS feeds
- Follow/unfollow feeds
- Browse posts from followed feeds
- Automatic feed aggregation with concurrent fetching
- PostgreSQL database for persistent storage

## Installation

### Prerequisites

- Go 1.16 or higher
- PostgreSQL database

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/D3rise/gator.git
   cd gator
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   - Create a PostgreSQL database
   - Run the migration scripts in the `sql/schema` directory (the project uses goose for migrations)

4. Build the application:
   ```bash
   go build -o gator
   ```

5. Create a configuration file at `~/.gatorconfig.json` with the following content:
   ```json
   {
     "db_url": "postgresql://username:password@localhost:5432/gator?sslmode=disable"
   }
   ```

   Alternatively, you can specify a custom config path using the `GATOR_CONFIG_PATH` environment variable.

## Usage

### Basic Commands

```bash
# Show help
./gator help

# Register a new user
./gator register <username>

# Login as a user
./gator login <username>

# List all users
./gator users

# Add a new RSS feed
./gator addfeed <feed_name> <feed_url>

# List all available feeds
./gator feeds

# Follow a feed
./gator follow <feed_id>

# List feeds you're following
./gator following

# Unfollow a feed
./gator unfollow <feed_id>

# Browse posts from followed feeds
./gator browse <page_number>

# Aggregate feeds (fetch new posts)
./gator agg
```

### Aggregating Feeds

The `agg` command starts the feed aggregation process, which:
- Fetches RSS feeds in the background
- Updates the database with new posts
- Displays feed content as it's fetched
- Continues running until interrupted with Ctrl+C

Example:
```bash
./gator agg
```

## Database Schema

The application uses the following database tables:

1. **user** - Stores user information
   - id (UUID, primary key)
   - name (VARCHAR, unique)
   - created_at, updated_at (TIMESTAMP)

2. **feed** - Stores feed information
   - id (UUID, primary key)
   - user_id (UUID, foreign key to user)
   - name (VARCHAR, unique)
   - url (VARCHAR)
   - created_at, updated_at (TIMESTAMP)
   - fetched_at (TIMESTAMP, nullable)

3. **feed_follow** - Tracks which users follow which feeds
   - id (UUID, primary key)
   - user_id (UUID, foreign key to user)
   - feed_id (UUID, foreign key to feed)
   - created_at, updated_at (TIMESTAMP)
   - Unique constraint on (user_id, feed_id)

4. **post** - Stores posts from feeds
   - id (UUID, primary key)
   - feed_id (UUID, foreign key to feed)
   - title (TEXT)
   - url (TEXT, unique)
   - description (TEXT)
   - published_at (TIMESTAMP)
   - created_at, updated_at (TIMESTAMP)

## Development

### Project Structure

- `main.go` - Application entry point
- `internal/` - Internal packages
  - `cli/` - CLI implementation
  - `commands/` - Command implementations
  - `config/` - Configuration handling
  - `database/` - Database access and models
  - `middleware/` - Command middleware (e.g., authentication)
  - `rss/` - RSS feed fetching and parsing
  - `state/` - Application state management
- `sql/` - SQL files
  - `queries/` - SQL queries (used by sqlc)
  - `schema/` - Database schema migrations

### Tools Used

- [sqlc](https://github.com/kyleconroy/sqlc) - Generate Go code from SQL
- [goose](https://github.com/pressly/goose) - Database migration tool
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver for Go

## License

[MIT License](LICENSE)