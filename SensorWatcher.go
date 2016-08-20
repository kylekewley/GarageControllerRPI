package main

type SensorWatcher struct {
}

type DoorStatus struct {
  Name string
  Status string
  LastChanged int
}

func (s *SensorWatcher) GetDoorStatus(doorName string) *DoorStatus {
  //TODO make this function work
  return &DoorStatus{
    Name: doorName,
    Status: "open",
    LastChanged: 0,
  }
}
