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

func GeneratePodSummaries(pods *v1.PodList) [][]string {
	//podInfoRows := []string{OutputColumnHeaders}
	var podInfoRows [][]string
	for _, pod := range pods.Items {
		elapsed := CalculatePodAge(pod)
		var podInfoRow []string
		podInfoRow = append(podInfoRow, fmt.Sprintf("%v", pod.Name))
		podInfoRow = append(podInfoRow, fmt.Sprintf("%v", pod.Status.Phase))
		podInfoRow = append(podInfoRow, fmt.Sprintf("%v", pod.Status.ContainerStatuses[0].RestartCount))
		podInfoRow = append(podInfoRow, fmt.Sprintf("%v", elapsed))
		podInfoRow = append(podInfoRow, fmt.Sprintf("%v", pod.Status.ContainerStatuses[0].Ready))
		if !pod.Status.ContainerStatuses[0].Ready {
			statusDetails, _ := json.Marshal(pod.Status.ContainerStatuses[0].State)
			podInfoRow = append(podInfoRow, string(statusDetails))
		} else {
			podInfoRow = append(podInfoRow, "None")
		}
		podInfoRows = append(podInfoRows, podInfoRow)
	}
	return podInfoRows
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
