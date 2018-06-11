package main

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

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getCurrent(clientconfig kubernetes.Clientset, serviceName string) []byte {
	res, _ := clientconfig.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
	sourceranges, _ := yaml.Marshal(res.Spec.LoadBalancerSourceRanges)
	return sourceranges

}

func getCurrentList(clientconfig kubernetes.Clientset, serviceName string) []string {
	res, _ := clientconfig.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
	return res.Spec.LoadBalancerSourceRanges
}

func patch(method string, clientconfig kubernetes.Clientset, resourceName string, cidr string) ([]byte, error) {

	cur := getCurrentList(clientconfig, resourceName)

	var payload []patchObject
	switch method {
	case "add":
		element := SliceIndex(len(cur), func(i int) bool {
			return cur[i] == cidr
		})
		if element != -1 {
			fmt.Printf("Specified CIDR %s is already present on the service allowed list", cidr)
			os.Exit(1)
		}
		payload = append(payload, patchObject{Op: "add", Path: "/spec/loadBalancerSourceRanges/" + strconv.Itoa(len(cur)), Value: cidr})

	case "remove":
		element := SliceIndex(len(cur), func(i int) bool {
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
