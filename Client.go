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

type ControlRequest struct {
  RequestType string
  Name string
}


/* Broker code */
func SubscribeToTopics(cli *client.Client, controlTopic string, controlHandler func(string, []byte) error) error {
    // Subscribe to topics.
    err := cli.Subscribe(&client.SubscribeOptions{
        SubReqs: []*client.SubReq{
            &client.SubReq{
                TopicFilter: []byte(controlTopic),
                QoS:         mqtt.QoS2,
                Handler: func(topicName, message []byte) {
                  err := controlHandler(string(topicName), message);
                  if err != nil {
                    log.Warningf("Issue with control message: '%s'", err)
                  }
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

  log.Debugf("Published metadata to topic %s: %s", metadataTopic, metadataJSON);

  return err
}

func PublishUpdateMessage(cli *client.Client, updateTopic string, doorStatus *DoorStatus) error {
    updateJSON,err := json.Marshal(doorStatus)

    if err != nil {
        log.Errorf("Error converting DoorStatus to JSON string. '%s'", err)
        return err
    }

    log.Debugf("Published Door update: %s", updateJSON);
    err = cli.Publish(&client.PublishOptions {
        QoS:    mqtt.QoS1,
        Retain: false,
        TopicName: []byte(updateTopic),
        Message: []byte(updateJSON),
    })

    if err != nil {
        log.Errorf("Error publishing door update to MQTT broker: '%s'", err)
    }

    return err
}

func ConnectToBroker(host string, port int, username string, password string) (*client.Client, error) {
    // Create an MQTT Client.
    cli := client.New(&client.Options{
        ErrorHandler: func(err error) {
            log.Errorf("MQTT client error: %s", err)
        },
    })

    options := &client.ConnectOptions{
            Network:  "tcp",
            Address:  fmt.Sprintf("%s:%d", host, port),
            CleanSession: true,
        }
    if len(username) > 0 {
        options.UserName = []byte(username)
        options.Password = []byte(password)
    }
    // Connect to the MQTT Server.
    err := cli.Connect(options)

    return cli, err
}
