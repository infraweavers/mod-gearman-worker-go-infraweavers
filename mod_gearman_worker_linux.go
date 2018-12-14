package modgearman

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func setupUsr1Channel(osSignalUsrChannel chan os.Signal) {
	signal.Notify(osSignalUsrChannel, syscall.SIGUSR1)
}

func mainSignalHandler(sig os.Signal) MainStateType {
	switch sig {
	case syscall.SIGTERM:
		logger.Infof("got sigterm, quiting gracefully")
		return ShutdownGraceFully
	case syscall.SIGINT:
		fallthrough
	case os.Interrupt:
		logger.Infof("got sigint, quitting")
		return Shutdown
	case syscall.SIGHUP:
		logger.Infof("got sighup, reloading configuration...")
		return Reload
	case syscall.SIGUSR1:
		logger.Errorf("requested thread dump via signal %s", sig)
		logThreaddump()
		return Resume
	default:
		logger.Warnf("Signal not handled: %v", sig)
	}
	return Resume
}

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
}

func processKill(p *os.Process) {
	p.Kill()
	// make sure hole process group is being killed
	syscall.Kill(-p.Pid, syscall.SIGKILL)
}
