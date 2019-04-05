package utils

import (
	"fmt"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"os"
	"strconv"
)

type patchObject struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// HomeDir retuns the calling user's home directory
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetCurrent returns the CIDR list active in the specified service
func GetCurrent(clientconfig kubernetes.Clientset, serviceName string) []byte {
	res, _ := clientconfig.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
	sourceranges, _ := yaml.Marshal(res.Spec.LoadBalancerSourceRanges)
	return sourceranges

}

func getCurrentList(clientconfig kubernetes.Clientset, serviceName string) []string {
	res, _ := clientconfig.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
	return res.Spec.LoadBalancerSourceRanges
}

// Patch modifies the loadBalancerSourceRanges field on the service by patching the JSON object
func Patch(method string, clientconfig kubernetes.Clientset, resourceName string, cidr string) ([]byte, error) {

	cur := getCurrentList(clientconfig, resourceName)

	var payload []patchObject
	switch method {
	case "add":
		element := sliceIndex(len(cur), func(i int) bool {
			return cur[i] == cidr
		})
		if element != -1 {
			fmt.Printf("Specified CIDR %s is already present on the service allowed list", cidr)
			os.Exit(1)
		}
		payload = append(payload, patchObject{Op: "add", Path: "/spec/loadBalancerSourceRanges/" + strconv.Itoa(len(cur)), Value: cidr})

	case "remove":
		element := sliceIndex(len(cur), func(i int) bool {
			return cur[i] == cidr
		})
		if element == -1 {
			fmt.Printf("Specified CIDR %s is not present on the service allowed list", cidr)
			os.Exit(1)
		}
		payload = append(payload, patchObject{Op: "remove", Path: "/spec/loadBalancerSourceRanges/" + strconv.Itoa(element)})
	}

	jsonload, _ := json.Marshal(payload)

	res, err := clientconfig.CoreV1().Services("default").Patch(resourceName, types.JSONPatchType, jsonload)
	sourceRanges, _ := yaml.Marshal(res.Spec.LoadBalancerSourceRanges)
	return sourceRanges, err
}
