# Run the Code

To run the code, follow these steps:

1. **Clone the Repository**: If you haven't already, clone the repository to your local machine using:
```bash
git clone git@github.com:bartick/go-task.git
cd go-task
```

2. **Install Dependencies**: Make sure you have all the necessary dependencies installed. You can usually do this by running:
```bash
go install
```

3. **Start Database**: Ensure that your database is running. If you're using a local database, start it up. If you're using a cloud database, make sure you have access to it.
```bash
make db-up

make db-migrate
```

4. **Set Environment Variables**: For this example you don't need to set any environment variables if using a localdb.
```bash
DB_USER=root
DB_PASS=
DB_HOST=localhost
DB_PORT=4000
DB_NAME=tasking
SERVER_ADDRESS=localhost
SERVER_PORT=3000
LOG_LEVEL=info
```
5. **Run the Application**: You can run the application using:
```bash
go build -o ./go-task ./app/cmd/go-task
./go-task

# Or Run with air
make watch
```


## Usage

Routes:
- `GET /tasks`: Retrieve all tasks
- `GET /tasks/{id}`: Retrieve a specific task by ID
- `POST /tasks`: Create a new task
- `PATCH /tasks/{id}`: Update an existing task by ID
- `DELETE /tasks/{id}`: Delete a task by ID
- `GET /tasks/{id}/subtasks`: Retrieve all subtasks for a specific task (and all nested subtasks)

Route Body
- `POST /tasks`
```json
{
    "title": "Task Title",
    "description": "Task Description", // Optional
    "status": "pending", // or "in_progress", "completed"
    "priority": 1, // integer value for task priority, higher number means higher priority
    "due_date": "2023-12-31T23:59:59Z", // Optional, in ISO 8601 format
    "parent_id": 1, // Optional, ID of the parent task if it's a subtask
    "category_name": "Backend" // Optional, or Frontend, Bug, Feature
}
```
- `PATCH /tasks/{id}`
```json
{
    "title": "Updated Task Title", // Optional
    "description": "Updated Task Description", // Optional
    "status": "in_progress", // Optional, can be "pending", "in_progress", "completed" (the software do not have a check to change to "completed" if has completed_at < NOW())
    "priority": 2, // Optional, integer value for task priority, higher number means higher priority
    "due_date": "2024-01-15T23:59:59Z", // Optional, in ISO 8601 format
    "completed_at": "2024-01-10T12:00:00Z", // Optional, in ISO 8601 format
    "parent_id": 2, // Optional, ID of the new parent task if changing
    "category_name": "Frontend" // Optional, or Backend, Bug, Feature
}