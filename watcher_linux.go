// +build linux

package proctify

import (
	"os"
	"strconv"
	"time"
)

func createNewWatcher() (*Watcher, error) {
	w := Watcher{
		Events:   make(chan Event, 500),
		Errors:   make(chan error),
		currents: make(map[int]Process),
		done:     make(chan struct{}),
		doneResp: make(chan struct{}),
	}

	err := w.startWatcher()
	if err != nil {
		return &w, err
	}

	return &w, nil
}

func (w *Watcher) startWatcher() error {
	running, err := runningProcesses()
	if err != nil {
		return err
	}

	for _, p := range running {
		w.currents[p.Pid()] = p
		w.Events <- Event{Process: p, Operation: OperationCreated}
	}

	go w.watch()

	return nil
}

func (w *Watcher) watch() {
	count := 0

	for {
		switch {
		case w.isClosed():
			close(w.doneResp)
			return
		default:
		}

		pDir, err := os.Open("/proc")
		if err != nil {
			w.Errors <- err
			return
		}
		names, err := pDir.Readdirnames(0)
		if err != nil {
			w.Errors <- err
			pDir.Close()
			continue
		}
		pDir.Close()

		count++

		pids := make([]int, 0, 100)

		for _, n := range names {
			// Filter only numberic dir name
			if n[0] < '0' || n[0] > '9' {
				continue
			}

			pid, err := strconv.ParseInt(n, 10, 0)
			if err != nil {
				continue
			}

			pids = append(pids, int(pid))
		}

		i := 0
		memPIDs := make([]int, len(w.currents))
		for n := range w.currents {
			memPIDs[i] = n
			i++
		}

		created := subtraction(pids, memPIDs)
		dead := subtraction(memPIDs, pids)
		running := intersection(pids, memPIDs)

		for _, pid := range created {
			p := &linuxProcess{pid: pid}
			err := p.reload()
			if err != nil {
				w.Errors <- err
				continue
			}
			w.currents[p.Pid()] = p
			w.Events <- Event{Process: p, Operation: OperationCreated}
		}

		for _, pid := range dead {
			p := w.currents[pid]
			delete(w.currents, pid)
			w.Events <- Event{Process: p, Operation: OperationDestroyed}
		}

		for _, pid := range running {
			currentProcess := w.currents[pid]
			updateProcess := &linuxProcess{pid: pid}
			err := updateProcess.reload()
			if err != nil {
				w.Errors <- err
				continue
			}

			if currentProcess.Path() != updateProcess.Path() {
				w.currents[pid] = updateProcess

				w.Events <- Event{Process: currentProcess, Operation: OperationDestroyed}
				w.Events <- Event{Process: updateProcess, Operation: OperationCreated}
			}

			if currentProcess.State() != updateProcess.State() || currentProcess.PPid() != updateProcess.PPid() {
				w.currents[pid] = updateProcess

				w.Events <- Event{Process: updateProcess, Operation: OperationModified}
			}
		}

		time.Sleep(2000 * time.Millisecond)
	}
}

func (w *Watcher) isClosed() bool {
	select {
	case <-w.done:
		return true
	default:
	}
	return false
}

func (w *Watcher) closeWatcher() {
	if w.isClosed() {
		return
	}

	// Send 'close' signal to goroutine, and set the Watcher to closed.
	close(w.done)

	// Wait for goroutine to close
	<-w.doneResp
}
