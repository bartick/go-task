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