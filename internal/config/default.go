package config

// defaultConfigs contains the default configuration settings for the application.
// These settings include server parameters, database connection details, JWT configuration, and logging options.
// Each key-value pair represents a configuration option with its default value.
var defaultConfigs = map[string]interface{}{
	// server.port specifies the port number where the application server will listen for incoming requests.
	// Default value is 8080.
	"server.port": 8080,

	// server.production defines the environment in which the application is running.
	// Possible values are "development" or "production".
	// Default value is "development".
	"server.production": false,

	// server.read_timeout sets the duration for which the server will wait to read the request.
	// Default value is "5s" (5 seconds).
	"server.read_timeout": "5s",

	// server.write_timeout sets the duration for which the server will wait to write the response.
	// Default value is "10s" (10 seconds).
	"server.write_timeout": "10s",

	// server.graceful_shutdown is the duration the server will wait before forcefully terminating ongoing requests during shutdown.
	// Default value is "30s" (30 seconds).
	"server.graceful_shutdown": "30s",

	// server.domain specifies the domain on which the server is accessible.
	// Default value is "http://localhost".
	"server.domain": "http://localhost",

	// db.host indicates the hostname or IP address of the database server.
	// Default value is "localhost".
	"db.host": "localhost",

	// db.port denotes the port number on which the database server is listening.
	// Default value is 5432, which is the default port for PostgreSQL.
	"db.port": 5432,

	// db.user specifies the username used to authenticate with the database server.
	// Default value is "root".
	"db.user": "root",

	// db.pass represents the password used for authentication with the database server.
	// Default value is "root".
	"db.pass": "root",

	// db.name is the name of the database to which the application will connect.
	// Default value is "test".
	"db.name": "test",

	// db.migration_enabled is a boolean flag that determines whether database migrations should be applied automatically on application startup.
	// Default value is false.
	"db.migration_enabled": false,

	// db.log_level sets the level of logging for database operations.
	// Default value is 2.
	"db.log_level": 2,

	// db.pool.max_open denotes the maximum number of open connections to the database.
	// Default value is 10.
	"db.pool.max_open": 10,

	// db.pool.max_idle represents the maximum number of idle connections in the pool.
	// Default value is 5.
	"db.pool.max_idle": 5,

	// db.pool.max_lifetime specifies the maximum amount of time a connection may remain in the pool.
	// Default value is "5m" (5 minutes).
	"db.pool.max_lifetime": "5m",

	// jwt.secret is the secret key used to sign and verify JSON Web Tokens (JWT).
	// Default value is "secret".
	"jwt.secret": "secret",

	// jwt.access_token_exp sets the duration for which an access token remains valid.
	// Default value is "900s" (15 minutes).
	"jwt.access_token_exp": "900s",

	// jwt.refresh_token_exp specifies the duration for which a refresh token remains valid.
	// Default value is "604800s" (7 days).
	"jwt.refresh_token_exp": "604800s",

	// logging.level determines the verbosity of the logging output.
	// Default value is -1
	"logging.level": -1,

	// logging.encoding defines the format of the log output.
	// Default value is "console".
	"logging.encoding": "console",

	// aws.region specifies the AWS region for cloud resources.
	// Default value is "eu-west-2".
	"aws.region": "eu-west-2",

	// aws.ses.from_email represents the default sender email address for AWS Simple Email Service (SES).
	// Default value is "example.com".
	"aws.ses.from_email": "example.com",
}
