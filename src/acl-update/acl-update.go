package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
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
	delCommand := flag.NewFlagSet("delOpts", flag.ExitOnError)

	// subcommands
	getServiceName := getCommand.String("service", "", "Resource name of the service to add to the new allowed IP range to")

	addServiceName := addCommand.String("service", "", "Resource name of the service to add to the new allowed IP range to")
	addNewIP := addCommand.String("cidr", "", "New CIDR to add to the allowed IP ranges for the service")

	delServiceName := delCommand.String("service", "", "Resource name of the service to add to the new allowed IP range to")
	delIP := delCommand.String("cidr", "", "CIDR to remove from the allowed IP ranges for the service")

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

	flag.Parse()

	switch os.Args[1] {
	case "get":
		getCommand.Parse(os.Args[2:])
	case "add":
		addCommand.Parse(os.Args[2:])
	case "del":
		delCommand.Parse(os.Args[2:])
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

	if addCommand.Parsed() {
		if *addServiceName == "" {
			addCommand.PrintDefaults()
			os.Exit(1)
		}
		if *addNewIP == "" {
			addCommand.PrintDefaults()
			os.Exit(1)
		}
		res, err := patch("add", *clientset, *addServiceName, *addNewIP)
		if err != nil {
			fmt.Printf("%v", err)
		}
		println(string(res))

	}

	if delCommand.Parsed() {
		if *delServiceName == "" {
			delCommand.PrintDefaults()
			os.Exit(1)
		}
		if *delIP == "" {
			delCommand.PrintDefaults()
			os.Exit(1)
		}
		res, err := patch("remove", *clientset, *delServiceName, *delIP)
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

}
