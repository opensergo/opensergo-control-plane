package client

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/opensergo/opensergo-control-plane/constant"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/dynamic"

	appv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
)

var dynamicClient dynamic.Interface
var client *kubernetes.Clientset

func Init() error {
	home, err := homeDir()
	if err != nil {
		log.Printf("Find empty home directory, %v", err)
		return err
	}
	kubeConfig := filepath.Join(home, ".kube", "config")

	// uses the current context get restConfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Panic(err)
	}
	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("Create dynamicClient for kubernetes failed %v", err)
		return err
	}
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Create client set for kubernetes failed %v", err)
		return err
	}

	return nil
}

func GetDeployment(ctx context.Context, deploymentName, namespace string) (deployment *appv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Get Deployment %s from namespace %s failed, err is %v", deploymentName, namespace, err)
	}
	return
}

func ApplyCRD(ctx context.Context, namespace, name string, unstructuredMap map[string]interface{}, gvr schema.GroupVersionResource) (un *unstructured.Unstructured, err error) {
	meta, err := prepareMeta(unstructuredMap)
	if err != nil {
		return nil, err
	}
	vs, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		meta[constant.META_UID], meta[constant.META_RESOURCE_VERSION], unstructuredMap[constant.CRD_METADATA] = "", vs.GetResourceVersion(), meta
		un, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(ctx, &unstructured.Unstructured{Object: unstructuredMap}, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Apply VirtualService in namespace %s failed, err is %v", namespace, err)
		}
	} else {
		un, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, &unstructured.Unstructured{Object: unstructuredMap}, metav1.CreateOptions{})
	}
	return
}

func DeleteCRD(ctx context.Context, namespace, name string, gvr schema.GroupVersionResource) error {
	return dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func DynamicClient() dynamic.Interface {
	return dynamicClient
}

func homeDir() (h string, err error) {
	// for linux
	if h = os.Getenv("HOME"); h != "" {
		return h, nil
	}
	// for windows
	if h = os.Getenv("USERPROFILE"); h != "" {
		return h, nil
	}
	return "", errors.New("No home directory")
}

func prepareMeta(unstructuredMap map[string]interface{}) (map[string]interface{}, error) {
	meta, ok := unstructuredMap[constant.CRD_METADATA].(map[string]interface{})
	if !ok {
		return nil, errors.New("Unknown exception, No metadata in crd")
	}
	meta[constant.META_RESOURCE_VERSION], unstructuredMap[constant.CRD_METADATA] = "", meta
	return meta, nil
}
