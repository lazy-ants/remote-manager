package runner

// Result holds a successful command output for one server.
type Result struct {
	Name  string
	Value string
}

// ErrorResult holds an error from one server.
type ErrorResult struct {
	Name    string
	Message string
}

// RunResult holds the combined results of running a command across servers.
type RunResult struct {
	Results  []Result
	Errors   []ErrorResult
	Timeouts []string
}
