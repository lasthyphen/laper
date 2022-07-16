// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package process

import (
	"fmt"
	"sync"

	"github.com/lasthyphen/beacongo/app"
	"github.com/lasthyphen/beacongo/nat"
	"github.com/lasthyphen/beacongo/node"
	"github.com/lasthyphen/beacongo/utils/constants"
	"github.com/lasthyphen/beacongo/utils/dynamicip"
	"github.com/lasthyphen/beacongo/utils/logging"
	"github.com/lasthyphen/beacongo/utils/perms"
	"github.com/lasthyphen/beacongo/utils/ulimit"
)

const (
	Header = `
    ____  _   _      __          _   __     __                      __
   / __ \(_) (_)__  / /______   / | / /__  / /__      ______  _____/ /__
  / / / / / / / _ \/ __/ ___/  /  |/ / _ \/ __/ | /| / / __ \/ ___/ //_/
 / /_/ / / / /  __/ /_(__  )  / /|  /  __/ /_ | |/ |/ / /_/ / /  / ,<
/_____/_/_/ /\___/\__/____/  /_/ |_/\___/\__/ |__/|__/\____/_/  /_/|_|
       /___/
`
)

var (
	stakingPortName = fmt.Sprintf("%s-staking", constants.AppName)
	httpPortName    = fmt.Sprintf("%s-http", constants.AppName)

	_ app.App = &process{}
)

// process is a wrapper around a node that runs in this process
type process struct {
	config node.Config
	node   *node.Node
	exitWG sync.WaitGroup
}

func NewApp(config node.Config) app.App {
	return &process{
		config: config,
		node:   &node.Node{},
	}
}

// Start the business logic of the node (as opposed to config reading, etc).
// Does not block until the node is done. Errors returned from this method
// are not logged.
func (p *process) Start() error {
	// Set the data directory permissions to be read write.
	if err := perms.ChmodR(p.config.DatabaseConfig.Path, true, perms.ReadWriteExecute); err != nil {
		return fmt.Errorf("failed to restrict the permissions of the database directory with: %w", err)
	}
	if err := perms.ChmodR(p.config.LoggingConfig.Directory, true, perms.ReadWriteExecute); err != nil {
		return fmt.Errorf("failed to restrict the permissions of the log directory with: %w", err)
	}

	// we want to create the logger after the plugin has started the app
	logFactory := logging.NewFactory(p.config.LoggingConfig)
	log, err := logFactory.Make("main")
	if err != nil {
		logFactory.Close()
		return err
	}

	// update fd limit
	fdLimit := p.config.FdLimit
	if err := ulimit.Set(fdLimit, log); err != nil {
		log.Fatal("failed to set fd-limit: %s", err)
		logFactory.Close()
		return err
	}

	// Track if sybil control is enforced
	if !p.config.EnableStaking {
		log.Warn("Staking is disabled. Sybil control is not enforced.")
	}

	// Check if transaction signatures should be checked
	if !p.config.EnableCrypto {
		// TODO: actually disable crypto verification
		log.Warn("transaction signatures are not being checked")
	}

	// Track if assertions should be executed
	if p.config.LoggingConfig.Assertions {
		log.Debug("assertions are enabled. This may slow down execution")
	}

	// TODO move this to config
	// SupportsNAT() for NoRouter is false.
	// Which means we tried to perform a NAT activity but we were not successful.
	if p.config.AttemptedNATTraversal && !p.config.Nat.SupportsNAT() {
		log.Warn("UPnP or NAT-PMP router attach failed, you may not be listening publicly. " +
			"Please confirm the settings in your router")
	}

	mapper := nat.NewPortMapper(log, p.config.Nat)

	// Open staking port we want for NAT Traversal to have the external port
	// (config.IP.Port) to connect to our internal listening port
	// (config.InternalStakingPort) which should be the same in most cases.
	if p.config.IPPort.IPPort().Port != 0 {
		mapper.Map(
			"TCP",
			p.config.IPPort.IPPort().Port,
			p.config.IPPort.IPPort().Port,
			stakingPortName,
			p.config.IPPort,
			p.config.DynamicUpdateDuration,
		)
	}

	// Open the HTTP port iff the HTTP server is not listening on localhost
	if p.config.HTTPHost != "127.0.0.1" && p.config.HTTPHost != "localhost" && p.config.HTTPPort != 0 {
		// For NAT Traversal we want to route from the external port
		// (config.ExternalHTTPPort) to our internal port (config.HTTPPort)
		mapper.Map(
			"TCP",
			p.config.HTTPPort,
			p.config.HTTPPort,
			httpPortName,
			nil,
			p.config.DynamicUpdateDuration,
		)
	}

	// Regularly updates our public IP (or does nothing, if configured that way)
	externalIPUpdater := dynamicip.NewDynamicIPManager(
		p.config.DynamicPublicIPResolver,
		p.config.DynamicUpdateDuration,
		log,
		p.config.IPPort,
	)

	if err := p.node.Initialize(&p.config, log, logFactory); err != nil {
		log.Fatal("error initializing node: %s", err)
		mapper.UnmapAllPorts()
		externalIPUpdater.Stop()
		log.Stop()
		logFactory.Close()
		return err
	}

	// [p.ExitCode] will block until [p.exitWG.Done] is called
	p.exitWG.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("caught panic", r)
			}
			log.Stop()
			logFactory.Close()
			p.exitWG.Done()
		}()
		defer func() {
			mapper.UnmapAllPorts()
			externalIPUpdater.Stop()

			// If [p.node.Dispatch()] panics, then we should log the panic and
			// then re-raise the panic. This is why the above defer is broken
			// into two parts.
			log.StopOnPanic()
		}()

		err := p.node.Dispatch()
		log.Debug("dispatch returned with: %s", err)
	}()
	return nil
}

// Stop attempts to shutdown the currently running node. This function will
// return immediately.
func (p *process) Stop() error {
	p.node.Shutdown(0)
	return nil
}

// ExitCode returns the exit code that the node is reporting. This function
// blocks until the node has been shut down.
func (p *process) ExitCode() (int, error) {
	p.exitWG.Wait()
	return p.node.ExitCode(), nil
}
