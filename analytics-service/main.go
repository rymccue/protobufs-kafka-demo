package main

import (
	"github.com/Shopify/sarama"
	"fmt"
	"github.com/rymccue/grpc-kafka-example/pb"
	"net/http"
	"github.com/golang/protobuf/proto"
	"encoding/json"
)

func newConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumer([]string{"kafka:29092"}, nil)
}

func main() {
	data := make(map[int32]int)
	consumer, err := newConsumer()
	if err != nil {
		fmt.Println("Error starting consumer,", err)
	}
	partitionList, err := consumer.Partitions("addition-analytics") //get all partitions on the given topic
	if err != nil {
		fmt.Println("Error retrieving partitionList ", err)
	}
	initialOffset := sarama.OffsetOldest //get offset for the oldest message on the topic
	for _, partition := range partitionList {
		pc, _ := consumer.ConsumePartition("addition-analytics", partition, initialOffset)

		go func(pc sarama.PartitionConsumer) {
			for km := range pc.Messages() {
				var message analytics.AdditionAnalytics
				err := proto.Unmarshal(km.Value, &message)
				if err != nil {
					fmt.Println("Error unmarshaling kafka message", err)
				}
				for _, num := range message.Numbers {
					currentVal, ok := data[num]
					if !ok {
						data[num] = 0
					}
					data[num] = int(num) + currentVal
				}
			}
		}(pc)
	}

	http.HandleFunc("/analytics", func(w http.ResponseWriter, req *http.Request) {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
