/*
Copyright © 2020 ToucanSoftware

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/klog"

	"github.com/ToucanSoftware/cloudship/cmd"
	"github.com/ToucanSoftware/cloudship/pkg/action"
	"github.com/ToucanSoftware/cloudship/pkg/cli"
)

var settings = cli.New()

func init() {
	log.SetFlags(log.Lshortfile)
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}

func initKubeLogs() {
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	gofs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(gofs)
	pflag.CommandLine.AddGoFlagSet(gofs)
	pflag.CommandLine.Set("logtostderr", "true")
}

func main() {
	initKubeLogs()

	actionConfig := new(action.Configuration)
	cmd, err := cmd.NewRootCmd(actionConfig, os.Stdout, os.Args[1:])
	if err != nil {
		debug("%+v", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		debug("%+v", err)
		os.Exit(1)
	}
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
}
