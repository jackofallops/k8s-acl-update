package main

import (
	"flag"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	//"k8s.io/api/core/v1"
	"fmt"
	"k8s.io/apimachinery/pkg/util/json"
	//jp "github.com/appscode/jsonpatch"
	"k8s.io/apimachinery/pkg/types"
)

type patchObject struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func main() {
	var kubeconfig *string
	//var myservice = "fluentd-es"

	// commands
	getCommand := flag.NewFlagSet("getOpts", flag.ExitOnError)
	addCommand := flag.NewFlagSet("addOpts", flag.ExitOnError)

	// subcommands
	getServiceName := getCommand.String("service", "", "Resource name of the service to add to the new allowed IP range to")

	addServiceName := addCommand.String("service", "", "Resource name of the service to add to the new allowed IP range to")
	addNewIP := addCommand.String("new", "", "New CIDR to add to the allowed IP ranges for the service")

	//flags
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	if len(os.Args) < 2 {
		fmt.Println("subcommand is required")
		os.Exit(1)
	}

	//newIP := flag.String("new", "", "New CIDR to add to the allowed IP ranges for the service")
	flag.Parse()
	//args := flag.Args()

	//if len(args) < 1 {
	//	panic("No operation method [get|set|add]")
	//}
	//println(string(*newIP))

	switch os.Args[1] {
	case "get":
		getCommand.Parse(os.Args[2:])
	case "add":
		addCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	// use the current context in kubeconfig
	if addCommand.Parsed() {
		if *addServiceName == "" {
			addCommand.PrintDefaults()
			os.Exit(1)
		}
		if *addNewIP == "" {
			addCommand.PrintDefaults()
			os.Exit(1)
		}
		res, err := patch(*clientset, *addServiceName, *addNewIP)
		if err != nil {
			fmt.Printf("%v", err)
		}
		println(string(res))

	}

	if getCommand.Parsed() {
		if *getServiceName == "" {
			getCommand.PrintDefaults()
			os.Exit(1)
		}
		println(string(getCurrent(*clientset, *getServiceName)))
	}

	//switch args[0] {
	//case "get":
	//	println(string("loadBalancerSourceRanges" + ""))
	//	println(string(getCurrent(*clientset, myservice)))
	//case "set":
	//	setNew(getCurrent(*clientset, myservice), "127.0.0.1/32")
	//	println(string("not implemented yet"))
	//case "add":
	//	println(string(*newIP))
	//	if *newIP != "" {
	//		patch(*clientset, myservice, *newIP)
	//	}
	//case "default":
	//	println(string("Please supply method - [get|set]"))
	//
	//}

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
	//println(string(sourcerangesNew))
	return sourceranges

}

func curJSON(clientconfig kubernetes.Clientset, serviceName string) []byte {
	res, _ := clientconfig.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
	ret, _ := json.Marshal(res)
	return ret
}

func patch(clientconfig kubernetes.Clientset, resourceName string, newSource string) ([]byte, error) {

	//cur := curJSON(clientconfig, resourceName)

	//TODO work out index to insert at end instead of top - to avoid rewriting NSG more than necessary

	var payload []patchObject

	payload = append(payload, patchObject{Op:"add", Path: "/spec/loadBalancerSourceRanges/0", Value: newSource})


	jsonload, _ := json.Marshal(payload)

	res, err := clientconfig.CoreV1().Services("default").Patch(resourceName, types.JSONPatchType, jsonload)
	sourceRanges, _ := yaml.Marshal(res.Spec.LoadBalancerSourceRanges)
	return sourceRanges, err
}

func setNew(cur []byte, newSource string) error {

	return nil
}
