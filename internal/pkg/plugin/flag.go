// Copyright (c) 2018-2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the URIs of this project regarding your
// rights to use or distribute this software.

package plugin

import (
	"github.com/spf13/pflag"
	pluginapi "github.com/sylabs/singularity/pkg/plugin"
	singularity "github.com/sylabs/singularity/pkg/runtime/engines/singularity/config"
)

type flagHook struct {
	flag     *pflag.Flag
	callback pluginapi.FlagCallbackFn
}

type flagRegistry struct {
	FlagSet *pflag.FlagSet
	Hooks   []flagHook
}

// RegisterStringFlag adds a string flag to the registry
func (r *flagRegistry) RegisterStringFlag(f pluginapi.StringFlagHook) error {
	r.FlagSet.StringP(f.Flag.Name, f.Flag.Shorthand, f.Flag.DefValue, f.Flag.Usage)

	hook := flagHook{
		flag:     r.FlagSet.Lookup(f.Flag.Name),
		callback: f.Callback,
	}

	r.Hooks = append(r.Hooks, hook)
	return nil
}

// RegisterBoolFlag adds a bool flag to the registry
func (r *flagRegistry) RegisterBoolFlag(f pluginapi.BoolFlagHook) error {
	r.FlagSet.BoolP(f.Flag.Name, f.Flag.Shorthand, false, f.Flag.Usage)

	hook := flagHook{
		flag:     r.FlagSet.Lookup(f.Flag.Name),
		callback: f.Callback,
	}
	r.Hooks = append(r.Hooks, hook)

	return nil
}

// AddFlagHooks will add the plugin defined flags to the input FlagSet
func AddFlagHooks(flagSet *pflag.FlagSet) {
	assertInitialized()

	flagSet.AddFlagSet(reg.FlagSet)
}

// FlagHookCallbacks will run the callback functions for all registered
// flag hooks
func FlagHookCallbacks(c *singularity.EngineConfig) {
	assertInitialized()

	for _, hook := range reg.flagRegistry.Hooks {
		hook.callback(hook.flag, c)
	}
}
