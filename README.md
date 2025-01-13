# 🚀 Gin + Go Authentication Starter App

Welcome to the **Gin + Go Authentication Starter App**! 🎉 This project is a template to help you quickly get started with a **Go** application using the **Gin** framework for creating secure APIs, a robust **authentication mechanism**, **graceful error handling**, and server management. It also integrates **PostgreSQL** for database management using **GORM** and **Resend** for sending emails. 📧

## 🔑 Features

- **Gin Framework**: Fast and lightweight web framework for building APIs.
- **Authentication Mechanism**: User authentication with JWT using Gin and the Go standard library.
- **Graceful Server Start/Stop**: Handle server startup and shutdown smoothly.
- **Error Handling**: Graceful error responses with appropriate status codes.
- **PostgreSQL Database**: Uses GORM as the ORM for database interactions.
- **Resend**: Easily send emails via the Resend API.

## 🛠️ Getting Started

### Prerequisites

Before running the application, make sure you have the following installed on your machine:

- [Go](https://golang.org/doc/install) (version 1.18+)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Docker](https://www.docker.com/products/docker-desktop) (optional, for containerization)
- [Resend Account](https://resend.com) (for email functionality)

### 🚀 Running Locally

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/omniflare/go_gin_starter
   cd go_gin_starter
   ```

2. **Set up the `.env` file**:

   Create a `.env` file in the root of the project directory with the following variables:

   ```env
   # Database configuration
    DATABASE_URL='get from neon.tech or supabase'

    SERVER_PORT=8080
    JWT_SECRET= *random JWT *
    CLIENT_URL=http://localhost:3000 <for NEXTJS>

   # Email configuration
    RESEND_API_KEY=your_resend_api_key
   ```

3. **Start the Database**:

   You can either use a local PostgreSQL instance or run it via Docker:

   ```bash
   docker run --name postgres -e POSTGRES_PASSWORD=mypassword -d -p 5432:5432 postgres
   ```
   preferred is to use a hosted postgres database via neon.tech or supabase.com

4. **Install Dependencies**:

   Run the following command to install the required dependencies:

   ```bash
   go mod tidy
   ```

5. **Start the Application**:

   To start the server, use the following command:

   ```bash
   make start
   ```

   The application will run on `http://localhost:8080` by default.

---

## ⚙️ Key Features in Detail

### 1. Authentication Mechanism 🔐

This app implements user authentication using JWT (JSON Web Tokens). The authentication system is built with Gin and Go's standard library.

- **User Registration**: Users can register by providing a username and password.
- **Login**: After registration, users can log in to obtain a JWT token.
- **JWT Token**: The token is used for authenticated routes, such as profile viewing or updating.

### 2. Graceful Server Start and Stop 🚦

This app handles graceful server startup and shutdown using Go's `context` and `http.Server` mechanisms. This ensures the server can handle active connections during shutdown and avoids data corruption.

- **Graceful Shutdown**: The server listens for termination signals and gracefully closes the active connections when the server is stopped.

### 3. Error Handling 🛑

- **Standardized Error Responses**: The app uses consistent error formats with appropriate HTTP status codes for better client-server communication.
- **Middleware**: Custom Gin middleware ensures that error responses are handled cleanly.

### 4. PostgreSQL Database & GORM Integration 📊

- **GORM** is used as the ORM to interact with the PostgreSQL database.
- You can define models, perform CRUD operations, and query the database with GORM.
- Example model:

   ```go
   type User struct {
     ID       uint   `json:"id" gorm:"primaryKey"`
	 Name     string
	 Email    string
	 Password string
   }
   ```

### 5. Sending Emails with Resend 📧

- The app uses **Resend** to send transactional emails, such as welcome emails or password reset links.
- Resend API integration is simple and secure.

### 6. Docker Support 🐳

The project is ready to be containerized with Docker. You can use the following command to build and run the app in a Docker container:

```bash
docker-compose up --build

or

make start
```

This will start both the Go application.
Additionally the PostgreSQL database container can also be started using the docker-compose.yaml file.

---

## 📝 API Endpoints

### 1. `POST /register`
- **Description**: Register a new user.
- **Request**: `{ "name": "user", "email":"email@mail.com", "password": "password123" }`
- **Response**: `{ "message": "User registered successfully" }`

### 2. `POST /api/authlogin`
- **Description**: Log in to get a JWT token.
- **Request**: `{ "username": "user", "password": "password123" }`
- **Response**: `{ "token": "jwt_token_here" }`

### 3. `GET /api/auth/me`
- **Description**: Get the authenticated user's profile.
- **Authorization**: Bearer JWT Token
- **Response**: `{ "username": "user", "email": "user@example.com" }`

---

## ⚖️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

If you'd like to contribute, feel free to fork the repository and create a pull request. We welcome any suggestions, bug fixes, or improvements! 😊

---

## 🙏 Acknowledgements

- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Resend](https://resend.com)
- [PostgreSQL](https://www.postgresql.org/)
```
