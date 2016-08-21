package main

import (
  "github.com/hugozhu/rpi"
)

type SensorWatcher struct {
  Doors []DoorStatus
}

type DoorStatus struct {
  Name string
  Status string
  LastChanged int
}

func NewSensorWatcherWithConfig(config &Config) (error,*SensorWatcher) {
}

func (s *SensorWatcher) GetDoorStatus(doorName string) *DoorStatus {
  //TODO make this function work
  return &DoorStatus{
    Name: doorName,
    Status: "open",
    LastChanged: 0,
  }
}
