# Web :: Login feature


## User Interface
* Use html template from @login.html

## User Flow
1. User navigates to the login page.
2. User enters their username and password.
3. System validates the input from  the user.
4. Send a request to API to authenticate the user.
5. If the authentication is successful, redirect the user to the dashboard.
6. If the authentication fails, display an error message to the user.

## Input validation in table formats
| Field | Validation Rule | Error Message | Input data |
|-------|-----------------|---------------|------------|
| Username | Required | "Please enter a valid email address." | somkiat |
| Password | Required, minimum 8 characters | "Password must be at least 8 characters long." | 12345678 |   


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

* Response with 200 
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user_id": 1,
    "username": "somkiat",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InNvbWtpYXQiLCJleHAiOjE2MzQyNzY4MDB9.abc123"
  }
```

* Response with 400     
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

## Acceptance Criterias of UI in table formats
| Criteria | Description | Test Case | Steps with input | Expect results |
|---------|-------------|-----------|-----------------|----------------|
| Login page loads successfully | The login page should load without errors | TC-001 | Navigate to the login page | The login page is displayed with username and password fields |
| Valid login credentials | User can log in with valid credentials | TC-002 | Enter valid username and password, click the login button | The user is redirected to the dashboard |


