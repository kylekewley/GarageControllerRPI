package main

import (
    "fmt"
    "encoding/json"
    "github.com/yosssi/gmq/mqtt"
    "github.com/yosssi/gmq/mqtt/client"
)

type Metadata struct {
  Name string
  Timezone string
  IsController bool
  Doors []DoorStatus
}



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

/* Send a full metadata topic on the metadata topic. The retain flag will be set to true */
func PublishMetadata(cli *client.Client, config *Config, sensorWatcher *SensorWatcher,
metadataTopic string) error {
  // Gather the required data
  var metadata Metadata
  metadata.Name = config.Controller.Name
  metadata.IsController = config.Controller.IsController
  metadata.Timezone = config.Controller.Timezone

  // Get the door info from the sensorWatcher
  for _,door := range config.Doors {
    doorName := door.Name
    status := sensorWatcher.GetDoorStatus(doorName)
    metadata.Doors = append(metadata.Doors, *status)
  }

  // Convert the struct to JSON
  metadataJSON, err := json.Marshal(metadata)

  if err != nil {
    return err
  }

  err = cli.Publish(&client.PublishOptions{
    QoS:       mqtt.QoS2,
    Retain:    true,
    TopicName: []byte(metadataTopic),
    Message:   []byte(metadataJSON),
  })

  log.Infof("Published metadata: %s", metadataJSON);

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
