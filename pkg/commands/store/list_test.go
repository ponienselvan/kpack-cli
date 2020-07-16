// Copyright 2020-2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package store_test

import (
	"testing"

	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	expv1alpha1 "github.com/pivotal/kpack/pkg/apis/experimental/v1alpha1"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/sclevine/spec"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pivotal/build-service-cli/pkg/commands/store"
	"github.com/pivotal/build-service-cli/pkg/testhelpers"
)

func TestStoreListCommand(t *testing.T) {
	spec.Run(t, "TestStoreListCommand", testStoreListCommand)
}

func testStoreListCommand(t *testing.T, when spec.G, it spec.S) {

	cmdFunc := func(clientSet *fake.Clientset) *cobra.Command {
		clientSetProvider := testhelpers.GetFakeKpackClusterProvider(clientSet)
		return store.NewListCommand(clientSetProvider)
	}

	when("stores exist", func() {
		it("returns a table of store details", func() {
			store1 := &expv1alpha1.Store{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-store-1",
				},
				Status: expv1alpha1.StoreStatus{
					Status: corev1alpha1.Status{
						Conditions: []corev1alpha1.Condition{
							{
								Type:   corev1alpha1.ConditionReady,
								Status: corev1.ConditionFalse,
							},
						},
					},
				},
			}

			store2 := &expv1alpha1.Store{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-store-2",
				},
				Status: expv1alpha1.StoreStatus{
					Status: corev1alpha1.Status{
						Conditions: []corev1alpha1.Condition{
							{
								Type:   corev1alpha1.ConditionReady,
								Status: corev1.ConditionUnknown,
							},
						},
					},
				},
			}

			store3 := &expv1alpha1.Store{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-store-3",
				},
				Status: expv1alpha1.StoreStatus{
					Status: corev1alpha1.Status{
						Conditions: []corev1alpha1.Condition{
							{
								Type:   corev1alpha1.ConditionReady,
								Status: corev1.ConditionTrue,
							},
						},
					},
				},
			}

			testhelpers.CommandTest{
				Objects: []runtime.Object{
					store1,
					store2,
					store3,
				},
				ExpectedOutput: `NAME            READY
test-store-1    False
test-store-2    Unknown
test-store-3    True

`,
			}.TestKpack(t, cmdFunc)
		})

		when("no stores exist", func() {
			it("returns a message that there are no stores", func() {
				testhelpers.CommandTest{
					ExpectErr:      true,
					ExpectedOutput: "Error: no stores found\n",
				}.TestKpack(t, cmdFunc)
			})
		})
	})
}