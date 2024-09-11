# Task: Build an API for a Simple Key-Value Store with Expiry Functionality

## Requirements:

### Create a RESTful API that:

Allows clients to store, retrieve, and delete key-value pairs.
Supports setting an optional expiration time for each key.
Expired keys should be automatically deleted when accessed after the expiration time.
Endpoints:

- `POST /store:` To store a key-value pair with an optional expiry time.
- `GET /store/:key:` To retrieve the value for a given key.
- `DELETE /store/:key`: To delete a key-value pair.

### Conditions:

The solution should use in-memory storage (no database required).
Include error handling for scenarios like expired keys, missing keys, and invalid input.
Write efficient code to ensure the expiration is handled in the background or at the time of access.
Optional (if time allows):

Implement an endpoint to get all active keys and their remaining time before expiry (`GET /store-all`).
Add basic logging for requests and key expirations.
Technical stack: He can use any backend framework he prefers (Node.js, Python/Flask, Java/Spring, etc.).

### After the task:

Ask him to record the entire process: approach, thought process, design decisions, and code explanation.

If time permits, have him reflect on how the task can be scaled up or optimized for production.
This task should be engaging for him, given his experience, and will allow him to showcase his technical abilities within a short time frame.