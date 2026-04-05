package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"pb"
)

func main() {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewLabServiceClient(conn)

	// Simulasi pengiriman metrik secara kontinyu (Client-side Streaming)
	stream, _ := client.ReportMetrics(context.Background())
	
	go func() {
		for {
			metric := &pb.ServerMetric{
				ServerId: "LAB-RUM-01",
				CpuUsage: rand.Float32() * 100,
				RamUsage: rand.Float32() * 100,
			}
			stream.Send(metric)
			time.Sleep(5 * time.Second)
		}
	}()

	// Simulasi alur mahasiswa
	// 1. Request
	res, _ := client.RequestEnvironment(context.Background(), &pb.EnvRequest{
		StudentId: "NIM12345",
		EnvType:   "DataScience_Python",
	})
	log.Printf("Job Created: %s", res.JobId)

	// 2. Monitor Progress
	logStream, _ := client.MonitorProvisioning(context.Background(), &pb.ProvisionJob{JobId: res.JobId})
	for {
		update, err := logStream.Recv()
		if err != nil { break }
		log.Printf("[%d%%] %s", update.Progress, update.Status)
	}
}
