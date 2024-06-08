package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Service struct {
	Name      string
	Namespace string
	Port      int32
	URL       string
}

func main() {
	http.HandleFunc("/", handleRoot)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	services, err := getK8sServices()
	tmpl := ""
	if err != nil {
		// http.Error(w, "Error retrieving services", http.StatusInternalServerError)
		// return
		tmpl = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>K8S Service List</title>

	</head>
	<body>
		<h1>K8S Service List</h1>
		<ul>
			<p> Error Loading services </p>
		</ul>
	</body>
	</html>`
	} else {

	tmpl = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>K8S Service List</title>

	</head>
	<body>
		<h1>K8S Service List</h1>
		<ul>
			{{ range . }}
			<li><a href="{{ .URL }}">{{ .Name }}.{{ .Namespace }}.svc.cluster.local:{{ .Port }}</a></li>
			{{ end }}
		</ul>
	</body>
	</html>`
}

	t, err := template.New("services").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, services); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func getK8sServices() ([]Service, error) {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: clientcmd.RecommendedHomeFile},
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error creating in-cluster config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing services: %v", err)
	}

	var result []Service
	for _, svc := range services.Items {
		for _, port := range svc.Spec.Ports {
			service := Service{
				Name:      svc.Name,
				Namespace: svc.Namespace,
				Port:      port.Port,
				URL:       fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", svc.Name, svc.Namespace, port.Port),
			}
			result = append(result, service)
		}
	}

	return result, nil
}
