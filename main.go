package main

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "time"
    
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

func main() {
    // 로그 파일 설정
    logFile, err := os.OpenFile("/var/log/events.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("Error opening log file: %v\n", err)
        os.Exit(1)
    }
    defer logFile.Close()

    // 버퍼된 writer 생성
    stdout := bufio.NewWriter(os.Stdout)
    fileWriter := bufio.NewWriter(logFile)

    // 현재 시간을 가져오는 함수 정의
    getTimeString := func() string {
        return time.Now().Format("2006/01/02 15:04:05.000000")
    }

    // 로그를 파일과 stdout에 동시에 쓰는 함수
    writeLog := func(format string, args ...interface{}) {
        logMessage := fmt.Sprintf("%s %s", getTimeString(), fmt.Sprintf(format, args...))
        stdout.WriteString(logMessage)
        fileWriter.WriteString(logMessage)
        stdout.Flush()    // 즉시 stdout으로 플러시
        fileWriter.Flush() // 즉시 파일로 플러시
    }

    // In-cluster config 로드
    config, err := rest.InClusterConfig()
    if err != nil {
        writeLog("Error loading in-cluster config: %v\n", err)
        os.Exit(1)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        writeLog("Error creating Kubernetes client: %v\n", err)
        os.Exit(1)
    }

    // 이벤트 생성 함수
    createEvents := func() {
        writeLog("Starting event creation cycle...\n")  // 디버깅을 위한 로그 추가
        
        // 실제 노드 이름 가져오기
        nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            writeLog("Error listing nodes: %v\n", err)
            return
        }
        if len(nodes.Items) == 0 {
            writeLog("No nodes found\n")
            return
        }

        // 실제 서비스 이름 가져오기
        services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            writeLog("Error listing services: %v\n", err)
            return
        }
        if len(services.Items) == 0 {
            writeLog("No services found in the default namespace\n")
            return
        }

        // Node 이벤트 생성
        for _, node := range nodes.Items {
            nodeEvent := &corev1.Event{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      fmt.Sprintf("test-node-event-%s", node.Name),
                    Namespace: "default",
                },
                InvolvedObject: corev1.ObjectReference{
                    Kind:       "Node",
                    Name:       node.Name,
                    APIVersion: "v1",
                },
                Reason:         "NodeStatusTest",
                Message:        "Node status test event",
                Type:           "Normal",
                FirstTimestamp: metav1.NewTime(time.Now()),
                LastTimestamp:  metav1.NewTime(time.Now()),
            }

            createdEvent, err := clientset.CoreV1().Events("default").Create(context.TODO(), nodeEvent, metav1.CreateOptions{})
            if err != nil {
                writeLog("Error creating node event for node %s: %v\n", node.Name, err)
                os.Exit(1)
            }
            
            // 노드 이벤트 생성 로그 기록
            writeLog("Node event created - Node: %s, Event: %s, Type: %s, Reason: %s, Message: %s, Time: %s\n",
                createdEvent.InvolvedObject.Name,
                createdEvent.Name,
                createdEvent.Type,
                createdEvent.Reason,
                createdEvent.Message,
                createdEvent.FirstTimestamp.Time)
        }

        // Service 이벤트 생성
        for _, service := range services.Items {
            serviceEvent := &corev1.Event{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      fmt.Sprintf("test-service-event-%s", service.Name),
                    Namespace: "default",
                },
                InvolvedObject: corev1.ObjectReference{
                    Kind:       "Service",
                    Name:       service.Name,
                    APIVersion: "v1",
                },
                Reason:         "ServiceStatusTest",
                Message:        "Service status test event",
                Type:           "Normal",
                FirstTimestamp: metav1.NewTime(time.Now()),
                LastTimestamp:  metav1.NewTime(time.Now()),
            }

            createdEvent, err := clientset.CoreV1().Events("default").Create(context.TODO(), serviceEvent, metav1.CreateOptions{})
            if err != nil {
                writeLog("Error creating service event for service %s: %v\n", service.Name, err)
                os.Exit(1)
            }

            // 서비스 이벤트 생성 로그 기록
            writeLog("Service event created - Service: %s, Event: %s, Type: %s, Reason: %s, Message: %s, Time: %s\n",
                createdEvent.InvolvedObject.Name,
                createdEvent.Name,
                createdEvent.Type,
                createdEvent.Reason,
                createdEvent.Message,
                createdEvent.FirstTimestamp.Time)
        }

        writeLog("Event creation cycle completed\n")    // 디버깅을 위한 로그 추가
    }

    // 무한 루프로 실행
    for {
        createEvents()  // 첫 실행
        
        // 다음 실행까지 1시간 대기
        writeLog("Waiting for next event creation cycle...\n")
        time.Sleep(1 * time.Hour)
    }
}