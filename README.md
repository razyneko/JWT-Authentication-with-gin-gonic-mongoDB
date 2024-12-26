# JWT Authentication with Gin Gonic

A RESTful API built using **Go** and **Gin Gonic** that implements JSON Web Token (JWT) based authentication for secure user login, registration, and session management. The API allows users to sign up, log in, and access protected resources with token-based authentication.

## Features

- **User Authentication**: Secure login and registration via JWT tokens.
- **Protected Routes**: Token-based authorization for accessing protected routes.
- **Role-Based Access**: Implemented basic user roles for varying access levels (optional).
- **MongoDB Integration**: User data stored in MongoDB for persistent authentication.
- **Session Management**: Token expiration and refresh features (optional).
- **API Endpoints**:
  - `/signup`: Registers a new user.
  - `/login`: Authenticates a user and returns a JWT token.
  - `/users`: Returns user data for authenticated users.
  - Protected routes require JWT in the `Authorization` header.

## Technologies Used

- **Backend**: Go, Gin Gonic
- **Authentication**: JWT (JSON Web Tokens)
- **Database**: MongoDB
- **Go Packages**:
  - `github.com/dgrijalva/jwt-go`: For generating and validating JWT tokens.
  - `github.com/gin-gonic/gin`: Web framework for building the REST API.
  - `go.mongodb.org/mongo-driver/mongo`: MongoDB driver for Go.

## Installation

### Prerequisites

1. **Go**: Make sure Go is installed on your system. If not, install it from [the Go website](https://golang.org/doc/install).

2. **MongoDB**: Ensure you have MongoDB running on your local machine or have access to a cloud database.

### Steps to Set Up

1. Clone the repository:
   ```bash
   git clone <repository_url>
   ```

2. Navigate to the project directory:
   ```bash
   cd jwt-authentication-gin
   ```

3. Install Go dependencies:
   ```bash
   go mod tidy
   ```

4. Set up your MongoDB connection:
   - Open the `routes/connection.go` file and configure the MongoDB URI if necessary.
   - By default, it connects to `mongodb://localhost:27017`.

5. Run the API server:
   ```bash
   go run main.go
   ```

6. The API will be accessible at `http://localhost:8000`.

## API Endpoints

### `POST /signup`
- Registers a new user.
- Request Body:
  ```json
  {
    "username": "user123",
    "password": "password123",
    "email": "user@example.com"
  }
  ```
- Response:
  - **200 OK**: User registered successfully.
  - **400 Bad Request**: Invalid input or existing user.

### `POST /login`
- Authenticates a user and returns a JWT token.
- Request Body:
  ```json
  {
    "username": "user123",
    "password": "password123"
  }
  ```
- Response:
  - **200 OK**: Returns JWT token.
  - **401 Unauthorized**: Invalid credentials.

### `GET /users`
- Returns the authenticated user's data (requires a valid JWT token in the `Authorization` header).
- Example Request:
  ```bash
  curl -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8000/user
  ```
- Response:
  - **200 OK**: Returns user data.
  - **401 Unauthorized**: Invalid or expired token.

### Protected Routes
- All protected routes require the JWT token to be passed in the `Authorization` header.
- Example:
  ```bash
  curl -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8000/protected
  ```

## JWT Token Expiry & Refresh (Optional)
- Tokens have a limited expiry time.
- To implement token refresh, you can extend the functionality of the API to accept a refresh token and return a new access token.

## Middleware

The application uses a middleware to protect routes that require authentication. This middleware checks the validity of the JWT token in the `Authorization` header.

Example:
```go
func TokenAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
            c.Abort()
            return
        }

        // Validate token here
        // If invalid, abort the request
    }
}
```

## Testing

You can use tools like **Postman** or **cURL** to interact with the API.

### Example cURL Commands:

- **Register a User**:
  ```bash
  curl -X POST http://localhost:8000/signup -d '{"username": "user123", "password": "password123", "email": "user@example.com"}' -H "Content-Type: application/json"
  ```

- **Login to Get Token**:
  ```bash
  curl -X POST http://localhost:8000/login -d '{"username": "user123", "password": "password123"}' -H "Content-Type: application/json"
  ```

- **Access Protected Route**:
  ```bash
  curl -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8000/user
  ```
