package pkg

import (
	"encoding/json"
	"fmt"
	"k8s.io/api/core/v1"
	"os"
	"time"
)

const (
	DateTimeFormat      = "2006-01-02 15:04:05 -0700 PDT"
	OutputColumnHeaders = "PodName - PodState - RestartCount - Age - Ready? - ErrorDetails"
)

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
