// Copyright 2023 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/distribution/distribution/v3/configuration"
	_ "github.com/distribution/distribution/v3/registry/storage/driver/inmemory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/falcosecurity/falcoctl/cmd"
	commonoptions "github.com/falcosecurity/falcoctl/pkg/options"
	testutils "github.com/falcosecurity/falcoctl/pkg/test"
)

const (
	rulesfiletgz  = "../pkg/test/data/rules.tar.gz"
	rulesfileyaml = "../pkg/test/data/rules.yaml"
	plugintgz     = "../pkg/test/data/plugin.tar.gz"
)

var (
	registry     string
	ctx          = context.Background()
	output       = gbytes.NewBuffer()
	rootCmd      *cobra.Command
	opt          *commonoptions.CommonOptions
	port         int
	orasRegistry *remote.Registry
	configFile   string
)

func TestPush(t *testing.T) {
	var err error
	RegisterFailHandler(Fail)
	port, err = testutils.FreePort()
	Expect(err).ToNot(HaveOccurred())
	registry = fmt.Sprintf("localhost:%d", port)
	RunSpecs(t, "root suite")
}

var _ = BeforeSuite(func() {
	config := &configuration.Configuration{}
	config.HTTP.Addr = fmt.Sprintf("localhost:%d", port)
	// Create and configure the common options.
	opt = commonoptions.NewOptions()
	opt.Initialize(commonoptions.WithWriter(output))
	opt.Printer.DisableStylingf()

	// Create the oras registry.
	orasRegistry, err = testutils.NewOrasRegistry(registry, true)
	Expect(err).ToNot(HaveOccurred())

	// Start the local registry.
	go func() {
		err := testutils.StartRegistry(context.Background(), config)
		Expect(err).ToNot(BeNil())
	}()

	// Create temporary directory used to save the configuration file.
	configFile, err = testutils.CreateEmptyFile("falcoctl.yaml")
	Expect(err).Should(Succeed())

})

var _ = AfterSuite(func() {
	configDir := filepath.Dir(configFile)
	Expect(os.RemoveAll(configDir)).Should(Succeed())
})

func executeRoot(args []string) error {
	rootCmd.SetArgs(args)
	rootCmd.SetOut(output)
	return cmd.Execute(rootCmd, opt.Printer)
}
