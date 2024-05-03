package main

import (
	"context"
	"log"
	"strings"

	pb "example.com/student/students/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewStudentServiceClient(conn)

	student := &pb.Student{
		Id:     1,
		Name:   "Shiva",
		Email:  "sss@gmail.com",
		Mobile: "9380018009",
	}

	// Add student
	addedStudent, err := client.AddStudent(context.Background(), student)
	if err != nil {
		log.Fatalf("Error in adding student: %v", err)
	}
	log.Printf("Added Student: %v", addedStudent)

	// Get student by ID
	retrievedStudent, err := client.GetStudentById(context.Background(), &pb.StudentId{Id: 1})
	if err != nil {
		if strings.Contains(err.Error(), "No Rows in result set") {
			log.Printf("Student with ID %d not found", student.Id)
		} else {
			log.Fatalf("Error in getting student: %v", err)
		}

	}
	log.Printf("Retrieved Student: %v", retrievedStudent)

	// Update student
	updatedStudent, err := client.UpdateStudent(context.Background(), &pb.Student{
		Id:     1,
		Name:   "SSShiva",
		Email:  "sss@gmail.com",
		Mobile: "9944147422",
	})
	if err != nil {
		log.Fatalf("Error in updating student: %v", err)
	}
	log.Printf("Updated Student: %v", updatedStudent)

	// Delete student
	deletedStudent, err := client.DeleteStudent(context.Background(), &pb.StudentId{Id: 1})
	if err != nil {
		log.Fatalf("Error in deleting student: %v", err)
	}
	log.Printf("Deleted Student: %v", deletedStudent)
}
