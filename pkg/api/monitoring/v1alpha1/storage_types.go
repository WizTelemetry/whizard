/*
Copyright 2021 The KubeSphere authors.

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

package v1alpha1

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	MonitoringPaodinService = "monitoring.paodin.io/service"
	MonitoringPaodinStorage = "monitoring.paodin.io/storage"
	MonitoringPaodinTenant  = "monitoring.paodin.io/tenant"

	FinalizerMonitoringPaodin = "finalizers.monitoring.paodin.io"
)

type StorageSpec struct {
	ThanosStorage *ThanosStorage `json:"thanos"`
}

type ThanosStorage struct {
	S3Storage S3Storage `json:"S3,omitempty"`
	Prefix    string    `json:"omitempty"`
}

type S3Storage struct {
	Bucket    string            `json:"bucket"`
	Endpoint  string            `json:"endpoint"`
	AccessKey string            `json:"access_key"`
	SecretKey string            `json:"secret_key"`
	Region    string            `json:"region,omitempty"`
	Params    map[string]string `json:"param,omitempty"`
}

type StorageStatus struct {
	ThanosResource ThanosResource `json:"thanos,omitempty"`
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type Storage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageSpec   `json:"spec,omitempty"`
	Status StorageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type StorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Storage `json:"items"`
}

func ManagedLabelByStorage(storage metav1.Object) map[string]string {
	return map[string]string{
		MonitoringPaodinStorage: storage.GetNamespace() + "." + storage.GetName(),
	}
}

func StorageNamespacedName(managedByStorage metav1.Object) *types.NamespacedName {
	ls := managedByStorage.GetLabels()
	if len(ls) == 0 {
		return nil
	}

	namespacedName := ls[MonitoringPaodinStorage]
	arr := strings.Split(namespacedName, ".")
	if len(arr) != 2 {
		return nil
	}

	return &types.NamespacedName{
		Namespace: arr[0],
		Name:      arr[1],
	}
}

func init() {
	SchemeBuilder = SchemeBuilder.
		Register(&Storage{}, &StorageList{})
}
