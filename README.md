# Codeforces-Like Backend System with Code Submission and Judging

This project is a backend system that allows users to submit code solutions to programming questions and have them judged asynchronously using the [Judge0 API](https://judge0.com/). The system is designed to handle code submissions, process them in the background, and return the results after executing test cases.

## Features
- User authentication (signup, login, refresh token).
- Contest creation and management.
- Problem creation and management within contests.
- Code submission to problems.
- Asynchronous code judging using the Judge0 API.
- Submission result retrieval.

## Routes
Here is an overview of the available API routes:

### User Management
- `POST /signup`: Register a new user.
- `POST /login`: User login.
- `POST /refresh-token`: Refresh user authentication token.
- `GET /users`: Retrieve a list of all users.
- `GET /user/:id1`: Retrieve a specific user by ID.

### Contest Management
- `POST /contest`: Create a new contest (Admin only).
- `DELETE /contest/:id2`: Delete a contest by ID (Admin only).
- `GET /contests`: Retrieve a list of all contests.
- `GET /contest/:id2`: Retrieve details of a specific contest.

### Problem Management
- `POST /contest/:id2/problem`: Add a new problem to a contest (Admin only).
- `DELETE /contest/:id2/problem/:id3`: Delete a problem from a contest (Admin only).
- `GET /contest/:id2/problems`: Retrieve a list of problems in a contest.
- `GET /contest/:id2/problem/:id3`: Retrieve details of a specific problem.

### Submission Management
- `POST /contest/:id2/problem/:id3/submit/:id1`: Submit code for a problem (Authenticated users).
- `GET /contest/:id2/problem/:id3/submission/:id4`: Retrieve the result of a submission.

## Judge0 Submission Flow
1. **Send Code**: Submit code to Judge0.
2. **Receive Token**: Receive a submission token.
3. **Check Result**: Use the token after a delay to check the result.
4. **Process Result**: Process the result and store it in the system.

## Technologies Used
- **Backend Framework**: [Go Gin Framework](https://gin-gonic.com/)
- **Database**: PostgreSQL
- **Code Execution API**: [Judge0](https://judge0.com/)
- **Authentication**: JWT Tokens
- **Middleware**: Role-based access control for admin and participants.


