package commands_test

import (
	"bytes"
	. "github.com/starkandwayne/cf-cli/cf/commands"
	"github.com/starkandwayne/cf-cli/cf/configuration"
	"github.com/starkandwayne/cf-cli/cf/errors"
	"github.com/starkandwayne/cf-cli/cf/trace"
	testapi "github.com/starkandwayne/cf-cli/testhelpers/api"
	testcmd "github.com/starkandwayne/cf-cli/testhelpers/commands"
	testconfig "github.com/starkandwayne/cf-cli/testhelpers/configuration"
	testreq "github.com/starkandwayne/cf-cli/testhelpers/requirements"
	testterm "github.com/starkandwayne/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/starkandwayne/cf-cli/testhelpers/matchers"
)

var _ = Describe("curl command", func() {
	var deps curlDependencies

	BeforeEach(func() {
		deps = newCurlDependencies()
	})

	It("does not pass requirements when not logged in", func() {
		runCurlWithInputs(deps, []string{"/foo"})
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			deps.requirementsFactory.LoginSuccess = true
		})

		It("fails with usage when not given enough input", func() {
			runCurlWithInputs(deps, []string{})
			Expect(deps.ui.FailedWithUsage).To(BeTrue())
		})

		It("passes requirements", func() {
			runCurlWithInputs(deps, []string{"/foo"})
			Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
		})

		It("makes a get request given an endpoint", func() {
			deps.curlRepo.ResponseHeader = "Content-Size:1024"
			deps.curlRepo.ResponseBody = "response for get"
			runCurlWithInputs(deps, []string{"/foo"})

			Expect(deps.curlRepo.Method).To(Equal("GET"))
			Expect(deps.curlRepo.Path).To(Equal("/foo"))
			Expect(deps.ui.Outputs).To(ContainSubstrings([]string{"response for get"}))
			Expect(deps.ui.Outputs).ToNot(ContainSubstrings(
				[]string{"FAILED"},
				[]string{"Content-Size:1024"},
			))
		})

		It("makes a post request given -X", func() {
			runCurlWithInputs(deps, []string{"-X", "post", "/foo"})

			Expect(deps.curlRepo.Method).To(Equal("post"))
			Expect(deps.ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
		})

		It("sends headers given -H", func() {
			runCurlWithInputs(deps, []string{"-H", "Content-Type:cat", "/foo"})

			Expect(deps.curlRepo.Header).To(Equal("Content-Type:cat"))
			Expect(deps.ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
		})

		It("sends multiple headers given multiple -H flags", func() {
			runCurlWithInputs(deps, []string{"-H", "Content-Type:cat", "-H", "Content-Length:12", "/foo"})

			Expect(deps.curlRepo.Header).To(Equal("Content-Type:cat\nContent-Length:12"))
			Expect(deps.ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
		})

		It("prints out the response headers given -i", func() {
			deps.curlRepo.ResponseHeader = "Content-Size:1024"
			deps.curlRepo.ResponseBody = "response for get"
			runCurlWithInputs(deps, []string{"-i", "/foo"})

			Expect(deps.ui.Outputs).To(ContainSubstrings(
				[]string{"Content-Size:1024"},
				[]string{"response for get"},
			))
			Expect(deps.ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
		})

		It("sets the request body given -d", func() {
			runCurlWithInputs(deps, []string{"-d", "body content to upload", "/foo"})

			Expect(deps.curlRepo.Body).To(Equal("body content to upload"))
			Expect(deps.ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
		})

		It("prints verbose output given the -v flag", func() {
			output := bytes.NewBuffer(make([]byte, 1024))
			trace.SetStdout(output)

			runCurlWithInputs(deps, []string{"-v", "/foo"})
			trace.Logger.Print("logging enabled")

			Expect([]string{output.String()}).To(ContainSubstrings([]string{"logging enabled"}))
		})

		It("prints a failure message when the response is not success", func() {
			deps.curlRepo.Error = errors.New("ooops")
			runCurlWithInputs(deps, []string{"/foo"})

			Expect(deps.ui.Outputs).To(ContainSubstrings(
				[]string{"FAILED"},
				[]string{"ooops"},
			))
		})
	})
})

type curlDependencies struct {
	ui                  *testterm.FakeUI
	config              configuration.Reader
	requirementsFactory *testreq.FakeReqFactory
	curlRepo            *testapi.FakeCurlRepository
}

func newCurlDependencies() (deps curlDependencies) {
	deps.ui = &testterm.FakeUI{}
	deps.config = testconfig.NewRepository()
	deps.requirementsFactory = &testreq.FakeReqFactory{}
	deps.curlRepo = &testapi.FakeCurlRepository{}
	return
}

func runCurlWithInputs(deps curlDependencies, inputs []string) {
	ctxt := testcmd.NewContext("curl", inputs)
	cmd := NewCurl(deps.ui, deps.config, deps.curlRepo)
	testcmd.RunCommand(cmd, ctxt, deps.requirementsFactory)
}
