package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/rymccue/grpc-kafka-example/pb"
	"time"
)

func newProducer(brokerList []string) (sarama.SyncProducer, error) {

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	return sarama.NewSyncProducer(brokerList, config)
}

func main() {
	http.HandleFunc("/add", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		p, err := newProducer([]string{"kafka:29092"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		numbersStr := req.URL.Query()["numbers"]
		total := 0
		numsList := make([]int32, 0)
		for _, numberStr := range numbersStr {
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			numsList = append(numsList, int32(number))
			total += number
		}
		message, err := proto.Marshal(&analytics.AdditionAnalytics{
			Numbers: numsList,
			Time:    time.Now().Unix(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _, err = p.SendMessage(&sarama.ProducerMessage{
			Topic: "addition-analytics",
			Value: sarama.ByteEncoder(message),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "total: %d", total)
	})

	http.ListenAndServe(":8080", nil)
}
