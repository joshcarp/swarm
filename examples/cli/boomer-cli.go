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

func createTask(pluginPath string) (task swarm.Task, err error) {
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return swarm.Task{}, err
	}
	loadedPlugin, err := plugin.Open(pluginPath)
	if err != nil {
		return swarm.Task{}, err
	}
	task = swarm.Task{}
	getName, err := loadedPlugin.Lookup("GetName")
	if err != nil {
		log.Println(err)
	} else {
		task.Namef = getName.(func() string)()
	}

	getWeight, err := loadedPlugin.Lookup("GetWeight")
	if err != nil {
		log.Println(err)
	} else {
		task.Weightf = getWeight.(func() int)()
	}

	execute, err := loadedPlugin.Lookup("Execute")
	if err != nil {
		return swarm.Task{}, err
	}

	task.Fn = execute.(func())
	return task, nil
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	plugins := strings.Split(plugins, ",")
	tasks := make([]swarm.Tasker, 0)
	for _, plugin := range plugins {
		task, err := createTask(plugin)
		if err != nil {
			log.Printf("Ignored plugin %s, Error: %v", plugin, err)
			continue
		}
		log.Println("Loaded task", task.Namef, "with weight", task.Weightf)
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		log.Fatalln("No valid plugin found, exit now.")
	}
	bm := swarm.NewBoomer("localhost", 5557)
	bm.Run(tasks...)
}

func init() {
	flag.StringVar(&plugins, "load-plugins", "plugin/demo.so", "Plugin list, separated by comma. Defaults to plugin/demo.so.")
}
