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
	"fmt"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-cli-plugin/internal/pkg"
	"os"
	"path/filepath"
	"time"
)

func TailPodStatuesLive() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "mook [flags]",
		Short:        "Live Tail of Pod Information",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			DisplayPodStatuesLive(args)
			return nil
		},
	}

	return cmd
}

func DisplayPodStatuesLive(args []string) {
	kubeClient := GetKubeClient()

	// An empty string returns all namespaces
	namespace := args[0]
	if namespace == "" {
		namespace = "default"
	}

	writer := uilive.New()
	writer.Start()

	for index := 0; true; index++ {
		pods, err := ListPods(namespace, kubeClient)
		if err != nil {
			fmt.Printf("Error getting pods: %v\n", err)
			os.Exit(1)
		}
		lines := pkg.GeneratePodSummaries(pods)

		pkg.ColorPrintLines(lines, writer)
		time.Sleep(2 * time.Second)
	}

	writer.Stop()
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
