package main

import (
    "errors"
  . "github.com/cyoung/rpi"
)

type IOController struct {
  Doors map[string]Door
}

func NewIOControllerWithConfig(config *Config) *IOController {
    controller := &IOController{
        Doors: make(map[string]Door),
    }


    // Loop through the doors in the config file and set them up
    for _,door := range config.Doors {
        if config.Controller.IsController {
            // Store the doors in the internal structres
            controller.Doors[door.Name] = door
            PinMode(BoardToPin(door.ControlPin), OUTPUT)
            DigitalWrite(BoardToPin(door.ControlPin), HIGH)
        }
    }

    return controller;
}

func (controller *IOController)triggerDoor(doorName string) error {
    door, ok := controller.Doors[doorName]
    if !ok {
        return errors.New("Door " + doorName + " not found")
    }
    DigitalWrite(BoardToPin(door.ControlPin), LOW)
    Delay(150)
    DigitalWrite(BoardToPin(door.ControlPin), HIGH)

    return nil
}

