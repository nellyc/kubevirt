package vm_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"kubevirt.io/kubevirt/pkg/kubecli"
	"kubevirt.io/kubevirt/tests"
)

var _ = Describe("VirtualMachine", func() {

	const unknownVM = "unknown-vm"
	const vmName = "testvm"
	var vmInterface *kubecli.MockVirtualMachineInterface
	var ctrl *gomock.Controller
	vm := kubecli.NewMinimalVM("testvm")

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		kubecli.GetKubevirtClientFromClientConfig = kubecli.GetMockKubevirtClientFromClientConfig
		kubecli.MockKubevirtClientInstance = kubecli.NewMockKubevirtClient(ctrl)
		// create mock interface
		vmInterface = kubecli.NewMockVirtualMachineInterface(ctrl)
		// set up mock client behavior
		kubecli.MockKubevirtClientInstance.EXPECT().VirtualMachine(k8smetav1.NamespaceDefault).Return(vmInterface).AnyTimes()
		// set up mock interface behavior
		vmInterface.EXPECT().Get(vm.Name, gomock.Any()).Return(vm, nil).AnyTimes()
		vmInterface.EXPECT().Get(unknownVM, gomock.Any()).Return(nil, errors.New("unknonw VM")).AnyTimes()
		vmInterface.EXPECT().Patch(vmName, gomock.Any(), gomock.Any()).Return(vm, nil).AnyTimes()
	})

	Context("should patch VM", func() {
		It("with spec:running:true", func() {
			cmd := tests.NewRepeatableVirtctlCommand("start", vmName)
			Expect(cmd()).To(BeNil())
		})

		It("with spec:running:false", func() {
			vm.Spec.Running = true
			cmd := tests.NewRepeatableVirtctlCommand("stop", vmName)
			Expect(cmd()).To(BeNil())
		})

		It("with spec:running:false when it's false already ", func() {
			vm.Spec.Running = false
			cmd := tests.NewRepeatableVirtctlCommand("stop", vmName)
			Expect(cmd()).NotTo(BeNil())
		})

	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
