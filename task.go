package swarm

// Task is like the "Locust object" in locust, the python version.
// When swarmer receives a start message from master, it will spawn several goroutines to run Task.Fn.
// But users can keep some information in the python version, they can't do the same things in swarmer.
// Because Task.Fn is a pure function.
type Task struct {
	// The weight is used to distribute goroutines over multiple tasks.
	Weightf int
	// Fn is called by the goroutines allocated to this task, in a loop.
	Fn    func()
	Namef string
}

func (t Task) Run() {
	t.Fn()
}

func (t Task) Name() string {
	return t.Namef
}

func (t Task) Weight() int {
	return t.Weightf
}

type Tasker interface {
	Weight() int
	Name() string
	Run()
}
