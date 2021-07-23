package main

import (
	"flag"
	"log"
	"os"
	"plugin"
	"strings"

	"github.com/joshcarp/swarm"
)

// Trying to implement boomer-cli without any test scenarios
// Users can write test scenarios as go plugins, like plugin/demo.go

var plugins string

func createTask(pluginPath string) (task *swarm.Task, err error) {
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, err
	}
	loadedPlugin, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	task = &swarm.Task{}
	getName, err := loadedPlugin.Lookup("GetName")
	if err != nil {
		log.Println(err)
	} else {
		task.Name = getName.(func() string)()
	}

	getWeight, err := loadedPlugin.Lookup("GetWeight")
	if err != nil {
		log.Println(err)
	} else {
		task.Weight = getWeight.(func() int)()
	}

	execute, err := loadedPlugin.Lookup("Execute")
	if err != nil {
		return nil, err
	}

	task.Fn = execute.(func())
	return task, nil
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	plugins := strings.Split(plugins, ",")
	tasks := make([]*swarm.Task, 0)
	for _, plugin := range plugins {
		task, err := createTask(plugin)
		if err != nil {
			log.Printf("Ignored plugin %s, Error: %v", plugin, err)
			continue
		}
		log.Println("Loaded task", task.Name, "with weight", task.Weight)
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		log.Fatalln("No valid plugin found, exit now.")
	}

	swarm.Run(tasks...)
}

func init() {
	flag.StringVar(&plugins, "load-plugins", "plugin/demo.so", "Plugin list, separated by comma. Defaults to plugin/demo.so.")
}
