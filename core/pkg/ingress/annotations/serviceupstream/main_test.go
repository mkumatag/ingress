/*
Copyright 2017 The Kubernetes Authors.

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

package serviceupstream

import (
	"testing"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	api "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func buildIngress() *extensions.Ingress {
	defaultBackend := extensions.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	return &extensions.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: extensions.IngressSpec{
			Backend: &extensions.IngressBackend{
				ServiceName: "default-backend",
				ServicePort: intstr.FromInt(80),
			},
			Rules: []extensions.IngressRule{
				{
					Host: "foo.bar.com",
					IngressRuleValue: extensions.IngressRuleValue{
						HTTP: &extensions.HTTPIngressRuleValue{
							Paths: []extensions.HTTPIngressPath{
								{
									Path:    "/foo",
									Backend: defaultBackend,
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestIngressAnnotationServiceUpstreamEnabled(t *testing.T) {
	ing := buildIngress()

	data := map[string]string{}
	data[annotationServiceUpstream] = "true"
	ing.SetAnnotations(data)

	val, _ := NewParser().Parse(ing)
	enabled, ok := val.(bool)
	if !ok {
		t.Errorf("expected a bool type")
	}

	if !enabled {
		t.Errorf("expected annotation value to be true, got false")
	}
}

func TestIngressAnnotationServiceUpstreamSetFalse(t *testing.T) {
	ing := buildIngress()

	// Test with explicitly set to false
	data := map[string]string{}
	data[annotationServiceUpstream] = "false"
	ing.SetAnnotations(data)

	val, _ := NewParser().Parse(ing)
	enabled, ok := val.(bool)
	if !ok {
		t.Errorf("expected a bool type")
	}

	if enabled {
		t.Errorf("expected annotation value to be false, got true")
	}

	// Test with no annotation specified, should default to false
	data = map[string]string{}
	ing.SetAnnotations(data)

	val, _ = NewParser().Parse(ing)
	enabled, ok = val.(bool)
	if !ok {
		t.Errorf("expected a bool type")
	}

	if enabled {
		t.Errorf("expected annotation value to be false, got true")
	}
}
