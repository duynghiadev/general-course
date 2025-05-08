# User Management API Documentation

## Project Overview

This is a RESTful web service built with Go that provides user management functionality using MongoDB as the database. The service allows for creating, retrieving individual users, and listing all users.

## Base URL

```plaintext
http://localhost:8080
```

## Database Configuration

The service uses MongoDB with the following connection details:

- Database Name: demo-web-server ([mongodb](https://cloud.mongodb.com/v2/681c2cf7aa2cf44da4024559#/metrics/replicaSet/681c2d4b4f56965b6903b995/explorer/demo-web-server/users/find))
- Collection: users

## API Endpoints

### 1. Create User

Creates a new user in the system.

Endpoint: `/user`
Method: `POST`
Content-Type: `application/json`

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Success Response:

- Status Code: 201 (Created)

  ```json
  {
    "message": "User created successfully"
  }
  ```

Error Responses:

- Status Code: 400 (Bad Request)
- When request body is invalid
- Status Code: 500 (Internal Server Error)
- When database operation fails

### 2. Get User by Username

Retrieves a specific user by their username.

Endpoint: `/user`
Method: `GET`
Query Parameters: `username`

Example Request:

```plaintext
GET /user?username=testuser
```

Success Response:

- Status Code: 200 (OK)

  ```json
  {
    "username": "testuser",
    "password": "testuser"
  }
  ```

Error Responses:

- Status Code: 404 (Not Found)
- When user doesn't exist
- Status Code: 500 (Internal Server Error)
- When database operation fails

### 3. Get All Users

Retrieves a list of all users in the system.

Endpoint: `/user`
Method: `GET`

Success Response:

- Status Code: 200 (OK)

  ```json
  [
    {
      "username": "user1",

      "password": "hashedpassword1"
    },

    {
      "username": "user2",

      "password": "hashedpassword2"
    }
  ]
  ```

Error Response:

- Status Code: 500 (Internal Server Error)
- When database operation fails

## Error Handling

The API returns appropriate HTTP status codes and error messages in the response body when an error occurs. Common error scenarios include:

- Invalid request data (400)
- Resource not found (404)
- Method not allowed (405)
- Server errors (500)

## Project Structure

- main.go : Contains the HTTP server setup and route handlers
- database/database.go : Contains database connection and operations
- User model is defined in both files with the following structure:

  ```go
  type User struct {
    Username string  `bson:"username"`
    Password string`bson:"password"`
  }
  ```

## Security Considerations

1. Passwords are currently stored in plain text - in a production environment, these should be hashed
2. No authentication/authorization mechanism is implemented
3. CORS policies are not configured
4. Rate limiting is not implemented

## Running the Server

The server runs on port 8080 by default. To start the server:

```bash
go run main.go
```

You should see the following messages when the server starts successfully:

```plaintext
Connected to MongoDB!
Server is running on port :8080
```
