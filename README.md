# JWT Auth with Gin

A simple **JWT-based authentication system** built using **Go** and **Gin Gonic** for user registration, login, and secure authentication.

## Features

- **User Registration**: Allows users to sign up with their details (first name, last name, email, phone, password, and user type).
- **User Login**: Users can log in with their email and password to receive a **JWT token**.
- **JWT Authentication**: Secures private routes with JWT tokens.
- **Password Hashing**: Secure password storage using bcrypt hashing.
- **Refresh Tokens**: Allows users to refresh their session without re-authenticating.

## Technologies

- **Go**: Backend language for API development.
- **MongoDB**: NoSQL database for storing user data.
- **Gin Gonic**: Go web framework for routing.
- **JWT**: For managing authentication tokens.
- **bcrypt**: For hashing passwords securely.

## Setup Instructions

### Prerequisites

- [Go](https://golang.org/dl/) (v1.18+)
- [MongoDB](https://www.mongodb.com/try/download/community) (local or cloud)

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/jwt-auth-with-gin.git
cd jwt-auth-with-gin
```

### 2. Create `.env` File

In the root directory, create a `.env` file and add:

```env
PORT=8000
MONGO_URI="mongodb://localhost:27017/go-auth"
JWT_SECRET="your-secret-key"
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Application

```bash
go run main.go
```

The app will be available at `http://localhost:8000`.

### 5. API Endpoints

#### `POST /users/signup`
- **Description**: Register a new user.
- **Request**:
  ```json
  {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "phone": "1234567890",
    "user_type": "USER"
  }
  ```
- **Response**: 
  - `201 Created` on success
  - `400 Bad Request` on validation failure

#### `POST /users/login`
- **Description**: Log in with email and password to receive a JWT.
- **Request**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```
- **Response**: 
  - `200 OK` with JWT token
  - `400 Bad Request` on invalid credentials

#### `GET /users/profile`
- **Description**: Get user profile (requires valid JWT).
- **Headers**: `Authorization: Bearer <JWT_TOKEN>`
- **Response**: 
  - `200 OK` with user details
  - `401 Unauthorized` on missing or invalid token

Here’s how you can update the README to mention that you can either use MongoDB locally on your PC or run it in a Docker container:

---

### 6. MongoDB Setup

You have two options to run MongoDB for this project:

#### Option 1: Use MongoDB Locally on Your PC

1. **Install MongoDB** on your machine if you haven't already. You can find installation instructions on the [MongoDB official website](https://www.mongodb.com/docs/manual/installation/).

2. **Run MongoDB**:
   - Start the MongoDB server locally by running:
     ```bash
     mongod
     ```
   - This will start MongoDB on the default port `27017`.

#### Option 2: Use MongoDB in a Docker Container

1. **Start MongoDB container**:
   If you prefer to use Docker, you can start a MongoDB container by running:
   ```bash
   docker run -d --name mongodb -p 27017:27017 mongo
   ```
   This command will pull the official MongoDB image, start the container, and expose the default MongoDB port (`27017`) on your machine.
