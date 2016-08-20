package main

import (
    "fmt"
    "github.com/yosssi/gmq/mqtt"
    "github.com/yosssi/gmq/mqtt/client"
)

/* Broker code */
func SubscribeToTopics(cli *client.Client, metadataTopic string, controlTopic string) error {
    // Subscribe to topics.
    err := cli.Subscribe(&client.SubscribeOptions{
        SubReqs: []*client.SubReq{
            &client.SubReq{
                TopicFilter: []byte(controlTopic),
                QoS:         mqtt.QoS2,
                Handler: func(topicName, message []byte) {
                    HandleControlRequest(metadataTopic, cli, string(topicName), message)
                },
            },
        },
    })

    return err
}

func ConnectToBroker(host string, port int) (*client.Client, error) {
    // Create an MQTT Client.
    cli := client.New(&client.Options{
        ErrorHandler: func(err error) {
            log.Errorf("MQTT client error: %s", err)
        },
    })

    // Connect to the MQTT Server.
    err := cli.Connect(&client.ConnectOptions{
        Network:  "tcp",
        Address:  fmt.Sprintf("%s:%d", host, port),
        ClientID: []byte("GarageHistoryServer"),
    })


    return cli, err
}

func HandleControlRequest(metadataTopic string, cli *client.Client, topicName string, message []byte) {
}
