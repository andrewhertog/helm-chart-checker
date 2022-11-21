package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Resource struct {
	GVR       schema.GroupVersionResource
	namespace string
	context   context.Context
	dynamic   dynamic.Interface
	Content   []unstructured.Unstructured
}

func (r *Resource) setContext() {
	if r.context == nil {
		r.context = context.Background()
	}
}

func (r *Resource) setDynamic() {
	if r.dynamic == nil {
		config := ctrl.GetConfigOrDie()
		r.dynamic = dynamic.NewForConfigOrDie(config)
	}
}

func (r *Resource) Resource(group string, version string, resource string) {
	r.GVR = schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
}

func (r *Resource) GetHelmReleaseForNamespace(namespace string) {
	r.setContext()
	r.setDynamic()
	r.namespace = namespace
	r.getResourcesDynamically()
}

func (r *Resource) getResourcesDynamically() {

	list, err := r.dynamic.Resource(r.GVR).Namespace(r.namespace).
		List(r.context, metav1.ListOptions{})

	if err != nil {
		fmt.Println(err)
	}

	r.Content = list.Items
}
