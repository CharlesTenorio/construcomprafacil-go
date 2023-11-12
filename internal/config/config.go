package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DEVELOPER    = "developer"
	HOMOLOGATION = "homologation"
	PRODUCTION   = "production"
)

type Config struct {
	PORT           string `json:"port"`
	Mode           string `json:"mode"`
	MongoDBConfig  `json:"mongo_config"`
	RedisConfig    RedisDBConfig `json:"redis_config"`
	RMQConfig      RMQConfig     `json:"rmq_config"`
	ConsumerConfig `json:"cconfig"`
}

type MongoDBConfig struct {
	MDB_URI              string `json:"mdb_uri"`
	MDB_NAME             string `json:"mdb_name"`
	MDB_CLIENT           string `json:"mdb_client"`
	MDB_DELIVERY_ADDRESS string `json:"mdb_delivery_address"`
	MDB_GRIFE            string `json:"mdb_grife"`
	MDB_ORDER            string `json:"mdb_order"`
	MDB_DET_ORDER        string `json:"mdb_det_order"`
	MDB_PAYMENT          string `json:"mdb_payment"`
	MDB_PRODUTC          string `json:"mdb_product"`
	MDB_SUPPLIER         string `json:"mdb_supplier"`
	MDB_USER             string `json:"mdb_user"`

	MDB_COLLECTION string `json:"mdb_collection"`
}

type RedisDBConfig struct {
	RDB_HOST string `json:"rdb_host"`
	RDB_PORT string `json:"rdb_port"`
	RDB_USER string `json:"rdb_user"`
	RDB_PASS string `json:"rdb_pass"`
	RDB_DB   int64  `json:"rdb_db"`
	RDB_DSN  string `json:"-"`
}

type RMQConfig struct {
	RMQ_URI                  string `json:"rmq_uri"`
	RMQ_MAXX_RECONNECT_TIMES int    `json:"rmq_maxx_reconnect_times"`
}

type ConsumerConfig struct {
	ExchangeName  string `json:"exchange_name"`
	ExchangeType  string `json:"exchange_type"`
	RoutingKey    string `json:"routing_key"`
	QueueName     string `json:"queue_name"`
	ConsumerName  string `json:"consumer_name"`
	ConsumerCount int    `json:"consumer_count"`
	PrefetchCount int    `json:"prefetch_count"`
	Reconnect     struct {
		MaxAttempt int `json:"max_attempt"`
		Interval   int `json:"interval"`
	}
}

func NewConfig() *Config {
	conf := defaultConf()

	SRV_PORT := os.Getenv("SRV_PORT")
	if SRV_PORT != "" {
		conf.PORT = SRV_PORT
	}

	SRV_MODE := os.Getenv("SRV_MODE")
	if SRV_MODE != "" {
		conf.Mode = SRV_MODE
	}

	SRV_RDB_HOST := os.Getenv("SRV_RDB_HOST")
	if SRV_RDB_HOST != "" {
		conf.RedisConfig.RDB_HOST = SRV_RDB_HOST
	}

	SRV_RDB_PORT := os.Getenv("SRV_RDB_PORT")
	if SRV_RDB_PORT != "" {
		conf.RedisConfig.RDB_PORT = SRV_RDB_PORT
	}

	SRV_RDB_USER := os.Getenv("SRV_RDB_USER")
	if SRV_RDB_USER != "" {
		conf.RedisConfig.RDB_USER = SRV_RDB_USER
	}

	SRV_RDB_PASS := os.Getenv("SRV_R_PASS")
	if SRV_RDB_PASS != "" {
		conf.RedisConfig.RDB_PASS = SRV_RDB_PASS
	}

	SRV_RDB_DB := os.Getenv("SRV_R_DB")
	if SRV_RDB_DB != "" {
		conf.RedisConfig.RDB_DB, _ = strconv.ParseInt(SRV_RDB_DB, 10, 64)
	}

	SRV_MDB_URI := os.Getenv("SRV_MDB_URI")
	if SRV_MDB_URI != "" {
		conf.MDB_URI = SRV_MDB_URI
	}

	SRV_MDB_NAME := os.Getenv("SRV_MDB_NAME")
	if SRV_MDB_NAME != "" {
		conf.MDB_NAME = SRV_MDB_NAME
	}

	SRV_MDB_COLLECTION := os.Getenv("SRV_MDB_COLLECTION")
	if SRV_MDB_COLLECTION != "" {
		conf.MDB_COLLECTION = SRV_MDB_COLLECTION
	}

	SRV_RDB_DSN := os.Getenv("SRV_RDB_DSN")
	if SRV_RDB_DSN != "" {
		conf.RedisConfig.RDB_DSN = SRV_MDB_COLLECTION
	}

	if len(conf.RedisConfig.RDB_HOST) > 3 {

		// "redis://<user>:<pass>@localhost:6379/<db>"
		// https://redis.uptrace.dev/guide/go-redis.html#connecting-to-redis-server

		conf.RedisConfig.RDB_DSN = fmt.Sprintf("redis://%s:%s@%s:%s/%v", conf.RedisConfig.RDB_USER, conf.RedisConfig.RDB_PASS, conf.RedisConfig.RDB_HOST, conf.RedisConfig.RDB_PORT, conf.RedisConfig.RDB_DB)
	}

	SRV_RMQ_URI := os.Getenv("SRV_RMQ_URI")
	if SRV_RMQ_URI != "" {
		conf.RMQConfig.RMQ_URI = SRV_RMQ_URI
	}

	CC_EX_NAME := os.Getenv("CC_EX_NAME")
	if CC_EX_NAME != "" {
		conf.ConsumerConfig.ExchangeName = CC_EX_NAME
	}

	CC_EX_TYPE := os.Getenv("CC_EX_TYPE")
	if CC_EX_TYPE != "" {
		conf.ConsumerConfig.ExchangeType = CC_EX_TYPE
	}

	CC_RT_KEY := os.Getenv("CC_RT_KEY")
	if CC_RT_KEY != "" {
		conf.ConsumerConfig.RoutingKey = CC_RT_KEY
	}

	CC_QU_NAME := os.Getenv("CC_QU_NAME")
	if CC_QU_NAME != "" {
		conf.ConsumerConfig.QueueName = CC_QU_NAME
	}

	CC_C_NAME := os.Getenv("CC_C_NAME")
	if CC_C_NAME != "" {
		conf.ConsumerConfig.ConsumerName = CC_C_NAME
	}

	CC_C_COUNT := os.Getenv("CC_C_COUNT")
	if CC_C_COUNT != "" {
		conf.ConsumerConfig.ConsumerCount, _ = strconv.Atoi(CC_C_COUNT)
	}

	CC_PREF_COUNT := os.Getenv("CC_PREF_COUNT")
	if CC_PREF_COUNT != "" {
		conf.ConsumerConfig.PrefetchCount, _ = strconv.Atoi(CC_PREF_COUNT)
	}

	CC_MAX_ATTEMPT := os.Getenv("CC_MAX_ATTEMPT")
	if CC_MAX_ATTEMPT != "" {
		conf.ConsumerConfig.Reconnect.MaxAttempt, _ = strconv.Atoi(CC_MAX_ATTEMPT)
	}

	CC_INTERVAL := os.Getenv("CC_INTERVAL")
	if CC_INTERVAL != "" {
		conf.ConsumerConfig.Reconnect.Interval, _ = strconv.Atoi(CC_INTERVAL)
	}

	return conf
}

func defaultConf() *Config {
	default_conf := Config{
		PORT: "8080",
		MongoDBConfig: MongoDBConfig{
			MDB_URI:        "mongodb://admin:supersenha@localhost:27017/",
			MDB_NAME:       "teste_db",
			MDB_COLLECTION: "hoodid",
		},

		Mode: DEVELOPER,

		RedisConfig: RedisDBConfig{
			RDB_HOST: "localhost",
			RDB_PORT: "6379",
			RDB_DB:   0,
			RDB_DSN:  "redis://localhost:6379/0",
		},
		RMQConfig: RMQConfig{
			RMQ_URI: "amqp://admin:supersenha@localhost:5672/",
		},

		ConsumerConfig: ConsumerConfig{
			ExchangeName:  "message_teams",
			ExchangeType:  "direct",
			RoutingKey:    "create",
			QueueName:     "SEND_MESSAGE_TEAMS",
			ConsumerName:  "CONSUMER_MESSAGE_TEAMS",
			ConsumerCount: 3,
			PrefetchCount: 1,
		},
	}

	return &default_conf
}
