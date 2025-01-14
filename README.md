# SchwarzIT Probearbeit

## Description

### API Documentation
Upon starting the application, the API documentation is available at path: `swagger/index.html`.
So in a local environment (not TLS secured) it would be `http://localhost:8080/swagger/index.html`

### Prerequisites
It is expected to have PostgreSQL and Redis running (docker or local installation). Specify the connection strings via environment variables (see below). Note: docker-compose.yml file is not tested but is provided for convenience.

### Environment Variables
The following environment variables are expected to be set:
- `POSTGRES_URI` (e.g. `postgres://user:password@localhost:5432/dbname`)
- `REDIS_URI` (e.g. `redis://user:password@localhost:6379`)
- `JWT_SECRET` (e.g. `secret`)
- `LOG_FILE` (e.g. `logs.json`)
- `DEBUG` (e.g. `true`)

### Building the Application
The application can be built using Task. Run the following command in the root directory of the project `task build`. Some Linux distributions might have to use ```go-task build```.

## Important Notes and Limitations

### Dependency Injection
Dependency injection is done via uber/fx. The cmd/standalone/main.go file is the entry point of the application and the dependencies are defined there.

### Logging
uber/zap is used for logging. Depending on whether the DEBUG environment variable is set to true or not, the logger will log in development or production mode, which means development logs debug level and is line based while production logs info level and is JSON based. Logs are written to the file have a similar log level as the console logs but are always in JSON format (as specified in the task).

### Testing
Given the scope of the project is probearbeit, the DB layer is fully integration tested (with docker containers) as example. The other layers (API and cache) do not have their own integration tests, but were locally tested regardless.

### Authentication
The API is secured with JWT access and refresh token pairs. This was only done because it was requested in the task. Otherwise, I would have not used custom authentication, rather a third party service like Kinde. Authentication and Admin checking was done when felt sensible as not concretely specified in the task.

### Database, Caching and Asynchronous Processing
As mentioned in the task it is explained here that the application uses PostgreSQL for the database and Redis for caching (wrapped DB) on certain methods. The application also uses a simple asynchronous processing mechanism for cache invalidation and setting to improve the request response time. Simplicity of the API did not require more complex asynchronous processing or caching.

### OpenAPI and CRUD Endpoints
Since the requested API's were a bit vaugue, Multiple CRUD endpoints were implemented which can be categorized in the following way:
- /auth/[login/refresh/register]
- /api/v1/users/{id} (R:GET, U:PUT/PATCH, D:DELETE) (requires admin access)
- /api/v1/users/ (C:POST, R:Query [with search params and pagination]) (requires admin access)
- /api/v1/users/me (R:GET, U:PUT/PATCH, D:DELETE) (for the authenticated user)

Note that the PUT method is used for full updates and PATCH is used for partial updates.

### Limitations
There are known bugs and features that are missing in the probearbeit, such as "Checking duplicate emails on User updates and registrations", "Lack of email confirmation in registration process", "No way of adding admin users without having to use the database directly", "Lack of password confirmation on registration or user updates", and etc. That being said, the probearbeit is a good example of a simple REST API with a few features, and the mentioned features are not realistically expected in a probearbeit.
The codebase also lacks implemented usecase layer which would have been required for the aforementioned features.

### Version Control and CI/CD
Only expect basic commits with no branching or merging (and thus no PRs or code reviews). There is a basic CI/CD pipeline in the .github/workflows directory that runs the tests and builds the application for pull requests and pushes to the main branch as well as releases.