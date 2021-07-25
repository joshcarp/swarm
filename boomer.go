package swarm

import (
	"log"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
)

// Mode is the running mode of swarmer, both standalone and distributed are supported.
type Mode int

const (
	// DistributedMode requires connecting to a master.
	DistributedMode Mode = iota
	// StandaloneMode will run without a master.
	StandaloneMode
)

// A Swarmer is used to run tasks.
// This type is exposed, so users can create and control a Swarmer instance programmatically.
type Swarmer struct {
	masterHost  string
	masterPort  int
	mode        Mode
	rateLimiter RateLimiter
	slaveRunner *slaveRunner

	localRunner *localRunner
	spawnCount  int
	spawnRate   float64

	cpuProfile         string
	cpuProfileDuration time.Duration

	memoryProfile         string
	memoryProfileDuration time.Duration

	outputs []Output

	Events EventBus.Bus
}

// NewSwarmer returns a new Swarmer.
func NewSwarmer(masterHost string, masterPort int) *Swarmer {
	return &Swarmer{
		masterHost: masterHost,
		masterPort: masterPort,
		mode:       DistributedMode,
		Events:     EventBus.New(),
	}
}

// NewStandaloneSwarmer returns a new Swarmer, which can run without master.
func NewStandaloneSwarmer(spawnCount int, spawnRate float64) *Swarmer {
	return &Swarmer{
		spawnCount: spawnCount,
		spawnRate:  spawnRate,
		mode:       StandaloneMode,
		Events:     EventBus.New(),
	}
}

// SetRateLimiter allows user to use their own rate limiter.
// It must be called before the test is started.
func (b *Swarmer) SetRateLimiter(rateLimiter RateLimiter) {
	b.rateLimiter = rateLimiter
}

// SetMode only accepts swarm.DistributedMode and swarm.StandaloneMode.
func (b *Swarmer) SetMode(mode Mode) {
	switch mode {
	case DistributedMode:
		b.mode = DistributedMode
	case StandaloneMode:
		b.mode = StandaloneMode
	default:
		log.Println("Invalid mode, ignored!")
	}
}

// AddOutput accepts outputs which implements the swarm.Output interface.
func (b *Swarmer) AddOutput(o Output) {
	b.outputs = append(b.outputs, o)
}

// EnableCPUProfile will start cpu profiling after run.
func (b *Swarmer) EnableCPUProfile(cpuProfile string, duration time.Duration) {
	b.cpuProfile = cpuProfile
	b.cpuProfileDuration = duration
}

// EnableMemoryProfile will start memory profiling after run.
func (b *Swarmer) EnableMemoryProfile(memoryProfile string, duration time.Duration) {
	b.memoryProfile = memoryProfile
	b.memoryProfileDuration = duration
}

// Run accepts a slice of Task and connects to the locust master.
func (b *Swarmer) Run(tasks ...Tasker) {
	if b.cpuProfile != "" {
		err := StartCPUProfile(b.cpuProfile, b.cpuProfileDuration)
		if err != nil {
			log.Printf("Error starting cpu profiling, %v", err)
		}
	}
	if b.memoryProfile != "" {
		err := StartMemoryProfile(b.memoryProfile, b.memoryProfileDuration)
		if err != nil {
			log.Printf("Error starting memory profiling, %v", err)
		}
	}

	switch b.mode {
	case DistributedMode:
		b.slaveRunner = newSlaveRunner(b.Events, b.masterHost, b.masterPort, tasks, b.rateLimiter)
		for _, o := range b.outputs {
			b.slaveRunner.addOutput(o)
		}
		b.slaveRunner.run()
	case StandaloneMode:
		b.localRunner = newLocalRunner(b.Events, tasks, b.rateLimiter, b.spawnCount, b.spawnRate)
		for _, o := range b.outputs {
			b.localRunner.addOutput(o)
		}
		b.localRunner.run()
	default:
		log.Println("Invalid mode, expected swarmer.DistributedMode or swarmer.StandaloneMode")
	}
}

// RecordSuccess reports a success.
func (b *Swarmer) RecordSuccess(requestType, name string, responseTime int64, responseLength int64) {
	if b.localRunner == nil && b.slaveRunner == nil {
		return
	}
	switch b.mode {
	case DistributedMode:
		b.slaveRunner.stats.requestSuccessChan <- &requestSuccess{
			requestType:    requestType,
			name:           name,
			responseTime:   responseTime,
			responseLength: responseLength,
		}
	case StandaloneMode:
		b.localRunner.stats.requestSuccessChan <- &requestSuccess{
			requestType:    requestType,
			name:           name,
			responseTime:   responseTime,
			responseLength: responseLength,
		}
	}
}

// RecordFailure reports a failure.
func (b *Swarmer) RecordFailure(requestType, name string, responseTime int64, exception string) {
	if b.localRunner == nil && b.slaveRunner == nil {
		return
	}
	switch b.mode {
	case DistributedMode:
		b.slaveRunner.stats.requestFailureChan <- &requestFailure{
			requestType:  requestType,
			name:         name,
			responseTime: responseTime,
			error:        exception,
		}
	case StandaloneMode:
		b.localRunner.stats.requestFailureChan <- &requestFailure{
			requestType:  requestType,
			name:         name,
			responseTime: responseTime,
			error:        exception,
		}
	}
}

// Quit will send a quit message to the master.
func (b *Swarmer) Quit() {
	b.Events.Publish("swarmer:quit")
	var ticker = time.NewTicker(3 * time.Second)

	switch b.mode {
	case DistributedMode:
		// wait for quit message is sent to master
		select {
		case <-b.slaveRunner.client.disconnectedChannel():
			break
		case <-ticker.C:
			log.Println("Timeout waiting for sending quit message to master, swarmer will quit any way.")
			break
		}
		b.slaveRunner.close()
	case StandaloneMode:
		b.localRunner.close()
	}
}

// Run tasks without connecting to the master.
func runTasksForTest(runTasks string, tasks ...Tasker) {
	taskNames := strings.Split(runTasks, ",")
	for _, task := range tasks {
		if task.Name() == "" {
			continue
		} else {
			for _, name := range taskNames {
				if name == task.Name() {
					log.Println("Running " + task.Name())
					task.Run()
				}
			}
		}
	}
}
