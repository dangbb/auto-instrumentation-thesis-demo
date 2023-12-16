package serve

import (
	"fmt"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	sarama "github.com/helios/go-sdk/proxy-libs/heliossarama"
	jsoniter "github.com/json-iterator/go"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logrus "github.com/helios/go-sdk/proxy-libs/helioslogrus"
	"microservice/config"
)

type KafkaConnector interface {
	InsertWarehouseHandler() error
}

type kafkaConnector struct {
	producer sarama.SyncProducer
}

type Event string
type Device string

type kafkaObject struct {
	event  Event
	device Device
	date   time.Time
}

func (s *kafkaConnector) InsertWarehouseHandler() error {
	object := kafkaObject{
		event:  "create-order",
		device: "phone",
		date:   time.Now(),
	}

	objectStr, err := jsoniter.Marshal(object)
	if err != nil {
		logrus.Errorf("Error marshal object: %v", object)
		return err
	}

	// send to kafka
	msg := &sarama.ProducerMessage{
		Topic: "warehouse",
		Key:   sarama.ByteEncoder("key"),
		Value: sarama.ByteEncoder(objectStr),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("header-key"),
				Value: []byte("header-value"),
			},
			{
				Key:   []byte("header-key-2"),
				Value: []byte("header-value-2"),
			},
		},
	}

	_, _, err = s.producer.SendMessage(msg)
	if err != nil {
		logrus.Infof("Error send kafka: %s", err.Error())
		return err
	}
	return nil
}

func newKafkaConnector(config config.Config) (KafkaConnector, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Version = sarama.V2_1_0_0
	cfg.Net.MaxOpenRequests = 1

	cfg.Producer.Compression = sarama.CompressionLZ4
	cfg.Producer.Idempotent = true
	cfg.Producer.Return.Successes = true

	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}

	logrus.Infof("Connect to broker addr: %s", config.KafkaConfig.Broker)

	brokers := []string{config.KafkaConfig.Broker}

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		logrus.Fatalf("cant establish sarama producer %s", err.Error())
		return nil, err
	}

	return &kafkaConnector{
		producer: producer,
	}, err
}

func Run(cfg config.Config) {
	logrus.SetLevel(logrus.DebugLevel)
	// define connector
	kafkaConn, err := newKafkaConnector(cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// create gin-gonic server
	r := gin.New()

	// endpoint to test logrus
	r.GET("/send-log", func(c *gin.Context) {
		logrus.Debug("Useful debugging information.")
		logrus.Debug("Useful debugging information.")

		logrus.Info("Something noteworthy happened!")
		logrus.Info("Something noteworthy happened!")

		logrus.Warn("You should probably take a look at this.")
		logrus.Warn("You should probably take a look at this.")

		logrus.Error("Error when process API send-log")
		logrus.Error("Error when process API send-log")

		// return ok
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// endpoint to test kafka
	r.GET("/send-kafka", func(c *gin.Context) {
		if err := kafkaConn.InsertWarehouseHandler(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort)); err != nil {
		logrus.Fatalf("cant run server %s", err.Error())
	}
}
