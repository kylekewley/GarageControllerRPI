package main

import (
    "os"
    "os/signal"
    "encoding/json"
    "github.com/op/go-logging"
    "github.com/spf13/viper"
)

var log = logging.MustGetLogger("log")

// Declare exit values
const (
    Success = iota
    ErrorReadingConfig
    ErrorConnecting
    ErrorSubscribing
    ErrorDisconnecting
)

type Status string
const (
    Opened Status = "open"
    Closed Status = "closed"
)

func main() {
    SetupLogging(logging.INFO)
    /////////////////////// Parse Config File /////////////////////////
    viper.SetConfigName("app")
    viper.AddConfigPath("./")

    // Set the defaults
    config := NewConfig()
    err := viper.ReadInConfig()

    // Exit immediately if there is an error reading the config file
    if err != nil {
      log.Criticalf("Fatal error reading config file: '%s'", err)
      os.Exit(ErrorReadingConfig)
    }

    err = viper.Unmarshal(config)

    // Exit immediately if there is an error reading the config file
    if err != nil {
      log.Criticalf("Fatal error reading config file: '%s'", err)
      os.Exit(ErrorReadingConfig)
    }

    // Extract the config data
    broker := config.Broker
    controller := config.Controller

    // Convert back to a string so we can log it
    configJSON, _ := json.MarshalIndent(config, "", "  ")

    // Make sure the config log level is valid
    level, err := logging.LogLevel(controller.LogLevel)

    // If the log level can't be parsed, default to info
    if err != nil {
        log.Warningf("The log level '%s' could not be parsed. Defaulting to INFO. " +
        "Should be (DEBUG|INFO|NOTICE|WARNING|ERROR|CRITICAL)", controller.LogLevel)
    }else {
        // Setup logging with the parsed level
        SetupLogging(level)
    }
    log.Debug("Logging setup properly")
    log.Infof("Config file read: %s", configJSON)
    ///////////////////// Done With Command Line Args /////////////////////////


    //// Connect to the Broker
    cli, err := ConnectToBroker(broker.Hostname, broker.Port)

    // Make sure the connection went smoothly
    if err != nil {
        log.Criticalf("Fatal error connecting to MQTT Broker: %s", err)
        os.Exit(ErrorConnecting)
    }
    log.Debugf("Successfully connected to MQTT broker %s:%i", broker.Hostname, broker.Port)

    // Subscribe to the request and update topics that we need to listen to
    err = SubscribeToTopics(cli, broker.MetadataTopic, broker.ControlTopic)


    // Make sure we subscribed to topics okay
    if err != nil {
        log.Criticalf("Fatal Error subscribing to topics: %s", err)
        os.Exit(ErrorSubscribing)
    }
    log.Debugf("Subscribed to topic '%s'", broker.ControlTopic)

    // Create the sensor watcher
    // TODO: setup the watcher
    sensorWatcher := new(SensorWatcher)

    // Initial publish metadata
    err = PublishMetadata(cli, config, sensorWatcher, config.Broker.MetadataTopic)

    ////////////////////////////////////////////////////////
    // Set up channel on which to send signal notifications.
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, os.Interrupt, os.Kill)

    log.Info("Initialization successful. Waiting for requests or updates")

    // Wait for receiving a signal.
    <-sigc

    // Disconnect the Network Connection.
    if err := cli.Disconnect(); err != nil {
        log.Errorf("Error while disconnecting: %s", err)
        os.Exit(ErrorDisconnecting)
    }
}
