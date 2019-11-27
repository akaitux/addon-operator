package main

import (
	"fmt"
	"os"

	. "github.com/flant/libjq-go"
	"gopkg.in/alecthomas/kingpin.v2"
	log "github.com/sirupsen/logrus"

	shell_operator_app "github.com/flant/shell-operator/pkg/app"
	"github.com/flant/shell-operator/pkg/executor"
	utils_signal "github.com/flant/shell-operator/pkg/utils/signal"

	operator "github.com/flant/addon-operator/pkg/addon-operator"
	"github.com/flant/addon-operator/pkg/app"
)

func main() {
	kpApp := kingpin.New(app.AppName, fmt.Sprintf("%s %s: %s", app.AppName, app.Version, app.AppDescription))

	// global defaults
	app.SetupGlobalSettings(kpApp)
	shell_operator_app.SetupGlobalSettings(kpApp)

	// print version
	kpApp.Command("version", "Show version.").Action(func(c *kingpin.ParseContext) error {
		fmt.Printf("%s %s\n", app.AppName, app.Version)
		return nil
	})

	// start main loop
	kpApp.Command("start", "Start events processing.").
		Default().
		Action(func(c *kingpin.ParseContext) error {
			shell_operator_app.SetupLogging()
			log.Infof("%s %s, shell-operator %s", app.AppName, app.Version, shell_operator_app.Version)

			// Be a good parent - clean up after the child processes
			// in case if addon-operator is a PID 1 process.
			go executor.Reap()

			jqDone := make(chan struct{})
			go JqCallLoop(jqDone)

			operator.Start()

			// Block action by waiting signals from OS.
			utils_signal.WaitForProcessInterruption()

			return nil
		})

	kingpin.MustParse(kpApp.Parse(os.Args[1:]))

	return
}
