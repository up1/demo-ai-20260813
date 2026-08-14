# Login API


## API Specification
### Endpoint
* HOST = `http://localhost:8000`
* `POST /api/login`
* Body request in JSON format:
```json
{
  "username": "somkiat",
  "password": "12345678"
}
```

* Response with 200 
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user_id": 1,
    "username": "somkiat",
    "address": "123 Main St, City, Country",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InNvbWtpYXQiLCJleHAiOjE2MzQyNzY4MDB9.abc123"
  }
}
```

* Response with 400 with input validation error 
```json
{
  "status": "error",
  "message": "Invalid username or password"
}
```

* Response with 500
```json
{
  "status": "error",
  "message": "Internal server error"
}
``` 


## Business flow
1. User sends a POST request to the login API with their username and password.
2. The system validates the input data.
3. If the input data is valid, the system checks the credentials against the database.
4. If the credentials are correct, the system generates an authentication token and returns a success response
5. If the credentials are incorrect, the system returns an error response with a message indicating that the username or password is invalid.
6. If there is an internal server error, the system returns an error response with a message indicating that there was an internal server error.
7. If login fails 3 times, the system locks the user account for 15 minutes and returns an error response with a message indicating that the account is locked.

## Input validation in table formats with status code=400
| Field | Validation Rule | Error Message | Input data |
|-------|-----------------|---------------|------------|
| Username | Required | "Please enter a valid email address." | somkiat |
| Password | Required, minimum 8 characters | "Password must be at least 8 characters long." | 12345678 |


## Database schema in table formats
Table : users
| Column | Type | Description |
|--------|------|-------------|
| user_id | integer | Primary key |
| username | varchar | Username of the user |
| password | varchar | Hashed password |
| token | varchar | Authentication token |  
| created_at | timestamp | Timestamp of when the user was created |
| updated_at | timestamp | Timestamp of when the user was last updated |

## Test cases in table formats
| Test Case | Input | Expected Output |
|-----------|-------|----------------|
| Valid login | {"username": "somkiat", "password": "12345678"} | {"status": "success", "message": "Login successful", "data": {"user_id": 1, "username": "somkiat", "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InNvbWtpYXQiLCJleHAiOjE2MzQyNzY4MDB9.abc123"}} |
| Invalid login | {"username": "somkiat", "password": "wrongpassword"} | {"status": "error", "message": "Invalid username or password"} |
| Missing username | {"password": "12345678"} | {"status": "error", "message": "Please enter a valid email address."}   
```