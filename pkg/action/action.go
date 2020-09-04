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

package action

import (
	//	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	//"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	//"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Configuration injects the dependencies that all actions share.
type Configuration struct {
	// RESTClientGetter is an interface that loads Kubernetes clients.
	RESTClientGetter RESTClientGetter
}

// RESTClientGetter gets the rest client
type RESTClientGetter interface {
	ToRESTConfig() (*rest.Config, error)
	ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error)
	ToRESTMapper() (meta.RESTMapper, error)
}
