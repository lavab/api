package env

import (
	"github.com/Sirupsen/logrus"
	"github.com/bitly/go-nsq"
	"github.com/dancannon/gorethink"
	"github.com/willf/bloom"

	"github.com/lavab/api/cache"
	"github.com/lavab/api/db"
	"github.com/lavab/api/factor"
)

var (
	// Config contains flags passed to the API
	Config *Flags
	// Log is the API's logrus instance
	Log *logrus.Logger
	// Rethink contains the RethinkDB session used in the API
	Rethink *gorethink.Session
	// Cache is the global instance of the cache interface
	Cache cache.Cache
	// Accounts is the global instance of AccountsTable
	Accounts *db.AccountsTable
	// Addresses is the global instance of Addresses table
	Addresses *db.AddressesTable
	// Tokens is the global instance of TokensTable
	Tokens *db.TokensTable
	// Keys is the global instance of KeysTable
	Keys *db.KeysTable
	// Contacts is the global instance of ContactsTable
	Contacts *db.ContactsTable
	// Reservations is the global instance of ReservationsTable
	Reservations *db.ReservationsTable
	// Emails is the global instance of EmailsTable
	Emails *db.EmailsTable
	// Labels is the global instance of LabelsTable
	Labels *db.LabelsTable
	// Files is the global instance of FilesTable
	Files *db.FilesTable
	// Threads is the global instance of ThreadsTable
	Threads *db.ThreadsTable
	// Factors contains all currently registered factors
	Factors map[string]factor.Factor
	// Producer is the nsq producer used to send messages to other components of the system
	Producer *nsq.Producer
	// PasswordBF is the bloom filter used for leaked password matching
	PasswordBF *bloom.BloomFilter
)
