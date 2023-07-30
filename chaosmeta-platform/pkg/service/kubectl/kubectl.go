package kubectl

import (
	"bytes"
	"chaosmeta-platform/util/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	yamlV2 "gopkg.in/yaml.v2"
	"io/ioutil"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	yamlUtil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type OptionsType string

const (
	GetOption    = OptionsType("GetOption")
	DeleteOption = OptionsType("DeleteOption")
	CreateOption = OptionsType("CreateOption")
	ListOption   = OptionsType("ListOption")
)

type KubectlService struct {
	DynamicClient dynamic.Interface
	ClientSet     *kubernetes.Clientset
	RESTMapper    meta.RESTMapper
}

type Kubectl struct {
	Kind        string    `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	APIVersion  string    `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Name        string    `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Namespace   string    `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	UID         types.UID `json:"uid,omitempty" protobuf:"bytes,5,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`
	YamlContext string    `json:"yamlContext"`
	Success     bool      `json:"success,omitempty"`
	Explain     string    `json:"explain,omitempty"`
}

type KubectlServiceInterface interface {
	ApplyByFile(ctx context.Context, file string, isReturnAfterFailure bool) (error, []Kubectl)
	ApplyByContent(ctx context.Context, content []byte, isReturnAfterFailure bool) (error, []Kubectl)
	DeleteByContent(ctx context.Context, content string) error
	GetByContent(ctx context.Context, content []byte) (error, []Kubectl)
}

type KubectlServiceFactory func(config *rest.Config) (KubectlServiceInterface, error)

func NewkubectlService(config *rest.Config) (KubectlServiceInterface, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	gr, err := restmapper.GetAPIGroupResources(clientSet.Discovery())
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	return &KubectlService{DynamicClient: dd, ClientSet: clientSet, RESTMapper: mapper}, nil
}

func (k *KubectlService) JSONToYAML(ctx context.Context, j []byte) ([]byte, error) {
	var jsonObj interface{}
	err := yamlV2.Unmarshal(j, &jsonObj)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return yamlV2.Marshal(jsonObj)
}

func (k *KubectlService) MargePatch(ctx context.Context, originalObj, updatedObj interface{}) ([]byte, error) {
	originalJSON, err := json.Marshal(originalObj)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	updatedJSON, err := json.Marshal(updatedObj)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}

	data, err := jsonpatch.CreateMergePatch(originalJSON, updatedJSON)
	if err != nil {
		log.Error(ctx, err)
		return nil, fmt.Errorf("failed to marge patch data, error: %s", err)
	}

	return data, nil
}

func (k *KubectlService) addKubectlStruct(ctx context.Context, KubectlList *[]Kubectl, isSuccess bool, KubectlStruct *unstructured.Unstructured, explain string) error {
	if KubectlList == nil {
		log.Error(ctx, "KubectlList is nil")
		return errors.New("KubectlList is nil")
	}

	if KubectlStruct == nil {
		log.Error(ctx, "KubectlStruct is nil")
		return errors.New("KubectlStruct is nil")
	}

	Kubectl := Kubectl{
		Kind:       KubectlStruct.GetKind(),
		APIVersion: KubectlStruct.GetAPIVersion(),
		Name:       KubectlStruct.GetName(),
		Namespace:  KubectlStruct.GetNamespace(),
		UID:        KubectlStruct.GetUID(),
		Success:    isSuccess,
		Explain:    explain,
	}

	unstructuredObjBytes, err := KubectlStruct.MarshalJSON()
	if err != nil {
		log.Error(ctx, err)
		return err
	}

	unstructuredObjYamlBytes, err := k.JSONToYAML(ctx, unstructuredObjBytes)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	Kubectl.YamlContext = string(unstructuredObjYamlBytes)
	*KubectlList = append(*KubectlList, Kubectl)
	return nil
}

func (k *KubectlService) option(ctx context.Context, fileBytes []byte, option OptionsType, isReturnAfterFailure bool) (error, []Kubectl) {
	var (
		Kubectl         []Kubectl
		errOptionDetail string
	)
	decoder := yamlUtil.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), len(fileBytes))
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			if err.Error() == "EOF" {
				if len(errOptionDetail) > 0 {
					return errors.New(errOptionDetail), Kubectl
				}
				break
			}
			log.Error(ctx, err)
			return err, Kubectl
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Error(ctx, err)
			return err, Kubectl
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		mapping, err := k.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Error(ctx, err)
			return err, Kubectl
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = k.DynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = k.DynamicClient.Resource(mapping.Resource)
		}

		switch option {
		case CreateOption:
			objGet, err := dri.Get(context.Background(), unstructuredObj.GetName(), metav1.GetOptions{})
			if err != nil {
				if !k8sErrors.IsAlreadyExists(err) {
					objCreate, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
					if err == nil {
						k.addKubectlStruct(ctx, &Kubectl, true, objCreate, "created successfully")
						break
					}
					errOptionDetail += ";" + err.Error()
					log.Error(ctx, err)
					k.addKubectlStruct(ctx, &Kubectl, false, unstructuredObj, "can not create")
					if isReturnAfterFailure {
						return err, Kubectl
					}
				}
			}

			unstructuredObj.SetResourceVersion(objGet.GetResourceVersion())
			data, err := k.MargePatch(ctx, objGet, unstructuredObj)
			if err != nil {
				log.Error(ctx, err)
				errOptionDetail += ";" + err.Error()
				if isReturnAfterFailure {
					return err, Kubectl
				}
				break
			}

			objUpdate, err := dri.Patch(context.Background(), unstructuredObj.GetName(), types.StrategicMergePatchType, data, metav1.PatchOptions{})
			if err != nil {
				k.addKubectlStruct(ctx, &Kubectl, false, unstructuredObj, "patch failed")
				log.Error(ctx, err)
				errOptionDetail += ";" + err.Error()
				if isReturnAfterFailure {
					return err, Kubectl
				}
				break
			}
			k.addKubectlStruct(ctx, &Kubectl, true, objUpdate, "apply successfully")
		case DeleteOption:
			err := dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
			if err != nil {
				errOptionDetail += ";" + err.Error()
				log.Error(ctx, err)
				if isReturnAfterFailure {
					return err, nil
				}
				break
			}
			k.addKubectlStruct(ctx, &Kubectl, true, unstructuredObj, "delete successfully")
		case GetOption:
			objGet, err := dri.Get(context.Background(), unstructuredObj.GetName(), metav1.GetOptions{})
			if err != nil {
				errOptionDetail += ";" + err.Error()
				log.Error(err)
				k.addKubectlStruct(ctx, &Kubectl, false, unstructuredObj, "can not get")
				if isReturnAfterFailure {
					return err, Kubectl
				}
				break
			}
			k.addKubectlStruct(ctx, &Kubectl, true, objGet, "get successfully")
		}
	}
	return nil, Kubectl
}

func (k *KubectlService) ApplyByFile(ctx context.Context, file string, isReturnAfterFailure bool) (error, []Kubectl) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err, nil
	}

	return k.option(ctx, fileBytes, CreateOption, isReturnAfterFailure)
}

func (k *KubectlService) ApplyByContent(ctx context.Context, content []byte, isReturnAfterFailure bool) (error, []Kubectl) {
	return k.option(ctx, content, CreateOption, isReturnAfterFailure)
}

func (k *KubectlService) DeleteByContent(ctx context.Context, content string) error {
	err, _ := k.option(ctx, []byte(content), DeleteOption, false)
	return err
}

func (k *KubectlService) GetByContent(ctx context.Context, content []byte) (error, []Kubectl) {
	return k.option(ctx, content, GetOption, false)
}
