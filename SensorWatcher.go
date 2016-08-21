package main

import (
  . "github.com/cyoung/rpi"
)

type SensorWatcher struct {
  Doors map[string]Door
  DoorStatuses map[string]DoorStatus
}

type DoorStatus struct {
  Name string
  Status string
  LastChanged int
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
            DigitalWrite(BoardToPin(door.ControlPin), LOW)
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
