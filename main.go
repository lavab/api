package main

import (
	"net"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/lavab/flag"
	"github.com/zenazn/goji/graceful"

	"github.com/lavab/api/env"
	"github.com/lavab/api/setup"
)

// TODO: "Middleware that implements a few quick security wins"
// 		 https://github.com/unrolled/secure

var (
	// Enable namsral/flag functionality
	configFlag = flag.String("config", "", "config file to load")
	// General flags
	bindAddress      = flag.String("bind", ":5000", "Network address used to bind")
	apiVersion       = flag.String("api_version", "v0", "Shown API version")
	logFormatterType = flag.String("log", "text", "Log formatter type. Either \"json\" or \"text\"")
	forceColors      = flag.Bool("force_colors", false, "Force colored prompt?")
	emailDomain      = flag.String("email_domain", "lavaboom.io", "Domain of the default email service")
	// Registration settings
	sessionDuration = flag.Int("session_duration", 72, "Session duration expressed in hours")
	// Cache-related flags
	redisAddress = flag.String("redis_address", func() string {
		address := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
		if address == "" {
			address = "127.0.0.1"
		}
		return address + ":6379"
	}(), "Address of the redis server")
	redisDatabase = flag.Int("redis_db", 0, "Index of redis database to use")
	redisPassword = flag.String("redis_password", "", "Password of the redis server")
	// Database-related flags
	rethinkdbAddress = flag.String("rethinkdb_address", func() string {
		address := os.Getenv("RETHINKDB_PORT_28015_TCP_ADDR")
		if address == "" {
			address = "127.0.0.1"
		}
		return address + ":28015"
	}(), "Address of the RethinkDB database")
	rethinkdbKey      = flag.String("rethinkdb_key", os.Getenv("RETHINKDB_AUTHKEY"), "Authentication key of the RethinkDB database")
	rethinkdbDatabase = flag.String("rethinkdb_db", func() string {
		database := os.Getenv("RETHINKDB_DB")
		if database == "" {
			database = "dev"
		}
		return database
	}(), "Database name on the RethinkDB server")
	// nsq and lookupd addresses
	nsqdAddress = flag.String("nsqd_address", func() string {
		address := os.Getenv("NSQD_PORT_4150_TCP_ADDR")
		if address == "" {
			address = "127.0.0.1"
		}
		return address + ":4150"
	}(), "Address of the nsqd server")
	lookupdAddress = flag.String("lookupd_address", func() string {
		address := os.Getenv("NSQLOOKUPD_PORT_4160_TCP_ADDR")
		if address == "" {
			address = "127.0.0.1"
		}
		return address + ":4160"
	}(), "Address of the lookupd server")
	// YubiCloud params
	yubiCloudID  = flag.String("yubicloud_id", "", "YubiCloud API id")
	yubiCloudKey = flag.String("yubicloud_key", "", "YubiCloud API key")
	// loggly
	logglyToken = flag.String("loggly_token", "", "Loggly token")
	// etcd
	etcdAddress  = flag.String("etcd_address", "", "etcd peer addresses split by commas")
	etcdCAFile   = flag.String("etcd_ca_file", "", "etcd path to server cert's ca")
	etcdCertFile = flag.String("etcd_cert_file", "", "etcd path to client cert file")
	etcdKeyFile  = flag.String("etcd_key_file", "", "etcd path to client key file")
	etcdPath     = flag.String("etcd_path", "settings/", "Path of the keys")
	// slack
	slackURL      = flag.String("slack_url", "", "URL of the Slack Incoming webhook")
	slackLevels   = flag.String("slack_level", "warning", "minimal level required to have messages sent to slack")
	slackChannel  = flag.String("slack_channel", "#notif-api-logs", "channel to which Slack bot will send messages")
	slackIcon     = flag.String("slack_icon", ":ghost:", "emoji icon of the Slack bot")
	slackUsername = flag.String("slack_username", "API", "username of the Slack bot")
	// Password bloom filter path
	bloomFilter = flag.String("bloom_filter", "bloom.db", "Bloom filter containing passwords")
	bloomCount  = flag.Uint("bloom_count", 14522336, "Estimated count of passwords in the bloom filter")
)

func main() {
	// Parse the flags
	flag.Parse()

	// Put config into the environment package
	env.Config = &env.Flags{
		BindAddress:      *bindAddress,
		APIVersion:       *apiVersion,
		LogFormatterType: *logFormatterType,
		ForceColors:      *forceColors,
		EmailDomain:      *emailDomain,

		SessionDuration: *sessionDuration,

		RedisAddress:  *redisAddress,
		RedisDatabase: *redisDatabase,
		RedisPassword: *redisPassword,

		RethinkDBAddress:  *rethinkdbAddress,
		RethinkDBKey:      *rethinkdbKey,
		RethinkDBDatabase: *rethinkdbDatabase,

		NSQdAddress:    *nsqdAddress,
		LookupdAddress: *lookupdAddress,

		YubiCloudID:  *yubiCloudID,
		YubiCloudKey: *yubiCloudKey,

		LogglyToken: *logglyToken,

		SlackURL:      *slackURL,
		SlackLevels:   *slackLevels,
		SlackChannel:  *slackChannel,
		SlackIcon:     *slackIcon,
		SlackUsername: *slackUsername,

		BloomFilter: *bloomFilter,
		BloomCount:  *bloomCount,
	}

	// Generate a mux
	mux := setup.PrepareMux(env.Config)

	// Make the mux handle every request
	http.Handle("/", mux)

	// Log that we're starting the server
	env.Log.WithFields(logrus.Fields{
		"address": env.Config.BindAddress,
	}).Info("Starting the HTTP server")

	// Initialize the goroutine listening to signals passed to the app
	graceful.HandleSignals()

	// Pre-graceful shutdown event
	graceful.PreHook(func() {
		env.Log.Info("Received a singnal, stopping the application")
	})

	// Post-shutdown event
	graceful.PostHook(func() {
		env.Log.Info("Stopped the application")
	})

	// Listen to the passed address
	listener, err := net.Listen("tcp", env.Config.BindAddress)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error":   err,
			"address": *bindAddress,
		}).Fatal("Cannot set up a TCP listener")
	}

	// Start the listening
	err = graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		// Don't use .Fatal! We need the code to shut down properly.
		env.Log.Error(err)
	}

	// If code reaches this place, it means that it was forcefully closed.

	// Wait until open connections close.
	graceful.Wait()
}
