/*
Copyright 2018 The Kubernetes Authors.

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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DateTimeFormat      = "2006-01-02 15:04:05 -0700 PDT"
	OutputColumnHeaders = "PodName - PodState - RestartCount - Age - Ready? - ErrorDetails"
)

func DisplayPodStatuesLive() []string {
	kubeClient := GetKubeClient()

	// An empty string returns all namespaces
	namespace := "twistlock"
	pods, err := ListPods(namespace, kubeClient)
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		os.Exit(1)
	}
	return GeneratePodSummaries(pods)
}

func GeneratePodSummaries(pods *v1.PodList) []string {
	podInfoLines := []string{OutputColumnHeaders}
	for _, pod := range pods.Items {
		elapsed := CalculatePodAge(pod)
		podInfoLine := fmt.Sprintf("%v - %v - %v - %v - %v", pod.Name, pod.Status.Phase, pod.Status.ContainerStatuses[0].RestartCount, elapsed, pod.Status.ContainerStatuses[0].Ready)
		if !pod.Status.ContainerStatuses[0].Ready {
			statusDetails, _ := json.Marshal(pod.Status.ContainerStatuses[0].State)
			podInfoLine += " - " + string(statusDetails)
		}
		podInfoLines = append(podInfoLines, podInfoLine)
	}
	return podInfoLines
}

func CalculatePodAge(pod v1.Pod) time.Duration {
	startTime, err := time.Parse(DateTimeFormat, pod.Status.StartTime.String())
	if err != nil {
		fmt.Printf("Error converting time: %v\n", err)
		os.Exit(1)
	}
	cTime := time.Now()
	return cTime.Sub(startTime)
}

func GetKubeClient() *kubernetes.Clientset {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}
	return clientSet
}

func ListPods(namespace string, client *kubernetes.Clientset) (*v1.PodList, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{FieldSelector: ""})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v\n", err)
		return nil, err
	}
	return pods, nil
}
