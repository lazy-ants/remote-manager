package ssh

import "fmt"

// ConnError indicates a connection failure.
type ConnError struct {
	Host string
	Addr string
	Err  error
}

func (e *ConnError) Error() string {
	if e.Addr != "" {
		return fmt.Sprintf("connection to %s (%s) failed: %v", e.Host, e.Addr, e.Err)
	}
	return fmt.Sprintf("connection to %s failed: %v", e.Host, e.Err)
}

func (e *ConnError) Unwrap() error { return e.Err }

// AuthError indicates an authentication failure.
type AuthError struct {
	Host string
	Err  error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("auth for %s failed: %v", e.Host, e.Err)
}

func (e *AuthError) Unwrap() error { return e.Err }

// CmdError indicates a command execution failure.
type CmdError struct {
	Host    string
	Command string
	Output  string
	Err     error
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("command on %s failed: %v", e.Host, e.Err)
}

func (e *CmdError) Unwrap() error { return e.Err }
