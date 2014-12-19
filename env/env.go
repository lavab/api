package env

import (
	"github.com/Sirupsen/logrus"
	"github.com/bitly/go-nsq"
	"github.com/dancannon/gorethink"

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
	// Tokens is the global instance of TokensTable
	Tokens *db.TokensTable
	// Keys is the global instance of KeysTable
	Keys *db.KeysTable
	// Contacts is the global instance of ContactsTable
	Contacts *db.ContactsTable
	// Reservations is the global instance of ReservationsTable
	Reservations *db.ReservationsTable
	// Factors contains all currently registered factors
	Factors map[string]factor.Factor
	// NSQProducer is used for email sending
	NSQProducer *nsq.Producer
	// DeliveryConsumer is the NSQ consumer used for email delivery confirmations
	DeliveryConsumer *nsq.Consumer
	// ReceiptConsumer is the NSQ consumer for new email receipt handling
	ReceiptConsumer *nsq.Consumer
)
