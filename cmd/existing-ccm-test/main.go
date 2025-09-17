/*
Copyright 2024 The Kubernetes Authors.

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
	"os"
	"testing"

	"k8s.io/klog/v2"
)

var (
	verbose = flag.Bool("verbose", false, "Enable verbose output")
)

func TestMain(m *testing.M) {
	// Parse flags before running tests
	flag.Parse()

	// Configure verbose logging if requested
	if *verbose {
		klog.InitFlags(nil)
		if err := flag.Set("v", "4"); err != nil {
			klog.Warningf("Failed to set verbose flag: %v", err)
		}
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// main function for standalone binary usage
func main() {
	// This binary is primarily designed to be run with 'go test' for Ginkgo tests
	// For standalone usage, it will just show help information
	flag.Parse()

	if len(os.Args) == 1 {
		// No arguments provided, show usage
		flag.Usage()
		os.Exit(1)
	}

	// If run as standalone binary, show information about Ginkgo usage
	klog.Info("This binary is designed to be run with 'go test' for Ginkgo-based tests.")
	klog.Info("For standalone testing, use: go test -v --kubeconfig <path>")
	klog.Info("For more options, see the Makefile targets: make test-ginkgo-*")
}
