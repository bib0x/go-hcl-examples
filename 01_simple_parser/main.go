package main

import (
    "fmt"
    "github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
    Type      string      `hcl:"type"`
    Name      string      `hcl:"name"`
    Resources []Resource  `hcl:"resource,block"`
}

type Resource struct {
    Name      string `hcl:"resource_name,label"`
    Type      string `hcl:"resource_type,label"`
    Username  string `hcl:"name"`
    State     string `hcl:"state"`
    Task      Task   `hcl:"task,block"`
}

type Task struct {
    Name            string `hcl:"task_name,label"`
    ConnectorName   string `hcl:"connector"`
}

func main() {
    var config Config

    if err := hclsimple.DecodeFile("config.hcl", nil, &config); err == nil {
        fmt.Printf("Raw %#v\n\n", config)
        fmt.Printf("Config type: %v\n", config.Type)
        fmt.Printf("Config name: %v\n", config.Name)

        for _, r := range config.Resources {
            fmt.Printf("Resource name: %v\n", r.Name)
            fmt.Printf("Resource type: %v\n", r.Type)
            fmt.Printf("Resource username: %v\n", r.Username)
            fmt.Printf("Resource state: %v\n", r.State)
            fmt.Printf("Task name: %v\n", r.Task.Name)
            fmt.Printf("Task connector: %v\n", r.Task.ConnectorName)
        }
    }
}
