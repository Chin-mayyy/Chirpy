# Chirpy

Chirpy is a social network similar to Twitter, created to showcase backend engineering skills using Go. This project demonstrates core concepts of backend development including user management, posting, following, and social feed generation.

## Features

- **User Authentication:** Secure signup and login functionality.
- **Chirp Posting:** Users can post short messages ("chirps").
- **Follow System:** Users can follow and unfollow other users.
- **Timeline Feed:** Displays a feed of chirps from users you follow.
- **Written in Go:** Built for performance and concurrency.

## Getting Started

1. **Clone the Repository**
   ```bash
   git clone https://github.com/Chin-mayyy/Chirpy.git
   cd Chirpy
   ```

2. **Install Dependencies**
   Make sure you have Go installed (version 1.18+ recommended).

   ```bash
   go mod download
   ```

3. **Run the Application**
   ```bash
   go run main.go
   ```

4. **API Usage**
   Chirpy exposes endpoints for authentication, posting, and user interaction. Explore the `main.go` and `/handlers` directory for API documentation and route definitions.

## Project Structure

- `main.go`: Application entry point, sets up routes and server.
- `handlers/`: Contains HTTP request handlers for various endpoints (authentication, chirps, users, etc).
- `models/`: Data models and database logic.
- `db/`: Database initialization and migration scripts.
- `utils/`: Utility functions for authentication, validation, etc.

## Contributing

Pull requests and issues are welcome. Please fork the repository and submit a pull request with your changes.


