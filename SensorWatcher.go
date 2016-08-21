package main

import (
    "time"
  . "github.com/cyoung/rpi"
)

type SensorWatcher struct {
  Doors map[string]Door
  DoorStatuses map[string]DoorStatus
}

type DoorStatus struct {
  Name string
  Status string
  LastChanged int64
}

func NewSensorWatcherWithConfig(config *Config) *SensorWatcher {
    watcher := &SensorWatcher{
        Doors: make(map[string]Door),
        DoorStatuses: make(map[string]DoorStatus),
    }


    // Loop through the doors in the config file and set them up
    for _,door := range config.Doors {
        // Store the doors in the internal structres
        watcher.Doors[door.Name] = door
        watcher.DoorStatuses[door.Name] = DoorStatus{
            Name: door.Name,
            Status: "initializing",
            LastChanged: 0,
        }

        // Do the WiringPI GPIO setup
        PinMode(BoardToPin(door.SensorPin), INPUT)
        if config.Controller.IsController {
            PinMode(BoardToPin(door.ControlPin), OUTPUT)
            DigitalWrite(BoardToPin(door.ControlPin), HIGH)
        }
    }

    watcher.UpdateValues(nil)

    return watcher
}

// This will loop forever updating the internal DoorStatus structure for
// all doors that have been initialized.
// The interval parameter is the number of ms to wait between updates
func (s *SensorWatcher) CheckValuesForever(interval int, changeHandler func(*DoorStatus)) {
    for {
        s.UpdateValues(changeHandler)
        Delay(interval)
    }
}

// Update the internal structure with new values
func (s *SensorWatcher) UpdateValues(changeHandler func(*DoorStatus)) {
    for key, door := range s.Doors {
        status := s.DoorStatuses[key]
        pin := BoardToPin(door.SensorPin)

        // Read the value using WiringPI
        statusString := "open"
        if DigitalRead(pin) == HIGH {
            statusString = "closed"
        }

        // Check if there was a change
        if statusString != status.Status {
            log.Infof("Door %s changed from '%s' to '%s'", key, status.Status, statusString)
            newStatus := DoorStatus{
                Name: key,
                Status: statusString,
                LastChanged: time.Now().Unix(),
            }

            // Broadcast the message. The changeHandler might be nil if this is
            // being called for initialization
            if changeHandler != nil {
                changeHandler(&newStatus)
            }

            s.DoorStatuses[key] = newStatus
        }
    }
}

func (s *SensorWatcher) GetDoorStatus(doorName string) *DoorStatus {
    status, ok := s.DoorStatuses[doorName]
    if !ok {
        // The door doesn't exist. Return dummy DoorStatus
        return &DoorStatus{
            Name: doorName,
            Status: "",
            LastChanged: 0,
        }
    }

    return &status
}
