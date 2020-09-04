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

package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/ToucanSoftware/cloudship/pkg/action"
)

const createHelp = `
This command consists of multiple subcommands to work with the chart cache.

The subcommands can be used to push, pull, tag, list, or remove Helm charts.
`

func newCreateCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "creates a deployment",
		Long:  createHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return action.NewCreate(cfg).Run(out)
		},
	}
}