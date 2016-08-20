package main

import (
    "flag"
    "os"
    "os/signal"
    "github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// Declare exit values
const (
    Success = iota
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
    /////////////////////// Parse Command Line Args /////////////////////////
    // Parse the logging level
    var logLevel string
    flag.StringVar(&logLevel, "l", "INFO", "The logging level string. {DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL}")

    // Port for the broker
    port := flag.Int("p", 1883, "The port number for the MQTT broker")

    // MQTT Host
    var hostname string
    flag.StringVar(&hostname, "h", "localhost", "The hostname of the MQTT broker")

    // Update topic. This is the topic where garage status updates are received
    var updateTopic string
    flag.StringVar(&updateTopic, "u", "home/garage/door/update",
                    "The topic to listen for garage door updates on")

    // Metadata topic
    var metadataTopic string
    flag.StringVar(&metadataTopic, "m", "home/garage/door/metadata", "The topic "+
                    "that metadata is sent out on")

    // Garage Control topic
    var controlTopic string
    flag.StringVar(&controlTopic, "c", "home/garage/door/control", "The topic "+
                    "that is used for garage control")

    // Parse additional arguments here...

    // Do the actual parsing
    flag.Parse()

    // Try to get the log level from the cmd line
    level, err := logging.LogLevel(logLevel)

    // If the log level can't be parsed, default to info
    if err != nil {
        level = logging.INFO
        SetupLogging(level)
        log.Warningf("The command line log level '%s' could not be parsed. "+
                     "Defaulting to INFO. Use the -h option for help.", logLevel)
    }else {
        // Setup logging with the parsed level
        SetupLogging(level)
    }
    log.Debug("Logging setup properly")
    ///////////////////// Done With Command Line Args /////////////////////////

    //// Connect to the Broker
    cli, err := ConnectToBroker(hostname, *port)

    // Make sure the connection went smoothly
    if err != nil {
        log.Criticalf("Fatal error connecting to MQTT Broker: %s", err)
        os.Exit(ErrorConnecting)
    }
    log.Debugf("Successfully connected to MQTT broker %s:%i", hostname, *port)

    // Subscribe to the request and update topics that we need to listen to
    err = SubscribeToTopics(cli, metadataTopic, controlTopic)

    // Make sure we subscribed to topics okay
    if err != nil {
        log.Criticalf("Fatal Error subscribing to topics: %s", err)
        os.Exit(ErrorSubscribing)
    }
    log.Debugf("Subscribed to topic '%s'", controlTopic)

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
