# Application configurations

## Introduction

We use [Koanf]("github.com/knadh/koanf") for managing configurations in our project. Koanf allows us to define default configuration values in the `default.go` file. If the necessary environment variables are not found, Koanf will fall back to these default values, ensuring that the application has sensible defaults in place.


## Setting up configurations

### Server Configuration

- **`SERVER_PORT`**: The port number where the server listens for incoming requests.
    - **Default**: `8080`

- **`SERVER_PRODUCTION`**: Indicates the environment in which the application is running (`true` for production, `false` for development).
    - **Default**: `false`

- **`SERVER_READ_TIMEOUT`**: Duration for reading requests.
    - **Default**: `5s`

- **`SERVER_WRITE_TIMEOUT`**: Duration for writing responses.
    - **Default**: `10s`

- **`SERVER_GRACEFUL_SHUTDOWN`**: Time to wait before forcefully terminating ongoing requests during shutdown.
    - **Default**: `30s`

- **`SERVER_DOMAIN`**: The domain on which the server is accessible.
    - **Default**: `http://localhost:4000`

## OAuth Configuration

### Google OAuth

- **`OAUTH_GOOGLE_CLIENT_ID`**: Client ID for Google OAuth.
    - **Default**: `"client-id"`

- **`OAUTH_GOOGLE_CLIENT_SECRET`**: Client Secret for Google OAuth.
    - **Default**: `"secret"`

- **`OAUTH_GOOGLE_REDIRECT_URL`**: URL for redirecting users after successful authentication with Google.
    - **Default**: `http://localhost:4000/api/v1/oauth/google/callback`

- **`OAUTH_GOOGLE_SCOPES`**: Scopes for Google OAuth permissions.
    - **Default**: `"email,profile"`

### Microsoft OAuth

- **`OAUTH_MICROSOFT_CLIENT_ID`**: Client ID for Microsoft OAuth.
    - **Default**: `"client-id"`

- **`OAUTH_MICROSOFT_CLIENT_SECRET`**: Client Secret for Microsoft OAuth.
    - **Default**: `"secret"`

- **`OAUTH_MICROSOFT_REDIRECT_URL`**: URL for redirecting users after successful authentication with Microsoft.
    - **Default**: `http://localhost:4000/api/v1/oauth/microsoft/callback`

- **`OAUTH_MICROSOFT_SCOPES`**: Scopes for Microsoft OAuth permissions.
    - **Default**: `"User.Read,openid"`

## Database Configuration

- **`DB_HOST`**: Hostname or IP address of the database server.
    - **Default**: `localhost`

- **`DB_PORT`**: Port number on which the database server is listening.
    - **Default**: `5432`

- **`DB_USER`**: Username for database authentication.
    - **Default**: `root`

- **`DB_PASSWORD`**: Password for database authentication.
    - **Default**: `root`

- **`DB_NAME`**: Name of the database.
    - **Default**: `test`

- **`DB_MIGRATIONS`**: Whether to apply database migrations automatically on startup.
    - **Default**: `false`

- **`DB_LOG_LEVEL`**: Log level for database operations.
    - **Default**: `2`

- **`DB_POOL_MAX_OPEN`**: Maximum number of open connections to the database.
    - **Default**: `10`

- **`DB_POOL_MAX_IDLE`**: Maximum number of idle connections in the pool.
    - **Default**: `5`

- **`DB_POOL_MAX_LIFETIME`**: Maximum time a connection may remain open.
    - **Default**: `5m`

## JWT Configuration

- **`JWT_SECRET`**: Secret key for signing and verifying JSON Web Tokens (JWT).
    - **Default**: `"secret"`

- **`JWT_ACCESS_TOKEN_EXP`**: Expiration time for access tokens.
    - **Default**: `900s` (15 minutes)

- **`JWT_REFRESH_TOKEN_EXP`**: Expiration time for refresh tokens.
    - **Default**: `604800s` (7 days)

## Logging Configuration

- **`LOGGING_LEVEL`**: Verbosity level for logging.
    - **Default**: `-1`

- **`LOGGING_ENCODING`**: Format of log output.
    - **Default**: `console`

## AWS Configuration

- **`AWS_REGION`**: AWS region for cloud resources.
    - **Default**: `eu-west-2`

- **`AWS_SES_FROM_EMAIL`**: Default sender email address for AWS Simple Email Service (SES).
    - **Default**: `"example.com"`

## Setting Environment Variables

To configure the application, set the environment variables as described above. You can set these variables in your environment or by using a `.env` file.

### Example using `.env` file:

```bash
SERVER_PORT=8080
SERVER_PRODUCTION=false
OAUTH_GOOGLE_CLIENT_ID="your-google-client-id"
OAUTH_GOOGLE_CLIENT_SECRET="your-google-client-secret"
DB_HOST="localhost"
DB_PORT=5432
JWT_SECRET="your-secret-key"
```
