package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	pb "example.com/student/students/proto"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type server struct {
	pb.StudentServiceServer
	//students map[int32]*pb.Student
	db *sql.DB
}

func (s *server) AddStudent(ctx context.Context, student *pb.Student) (*pb.Student, error) {
	/*log.Println("AddStudent was invoked")
	s.students[student.Id] = student
	log.Printf("Adding - %v\n", student)
	return student, nil*/

	log.Println("AddStudent was invoked")
	_, err := s.db.Exec("insert into students (name, email, mobile) values (?,?,?)", student.Name, student.Email, student.Mobile)
	if err != nil {
		return nil, err
	}
	log.Printf("Adding - %v\n", student)
	return student, nil

}

func (s *server) GetStudentById(ctx context.Context, id *pb.StudentId) (*pb.Student, error) {
	/*log.Println("GetStudentById was invoked")
	log.Printf("Getting - %v\n", s.students[id.Id])
	return s.students[id.Id], nil*/

	log.Println("GetStudentById was invoked")
	var student pb.Student
	err := s.db.QueryRow("select id, name, email, mobile from students where id=?", id.Id).Scan(&student.Id, &student.Name, &student.Email, &student.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Student with Id %d not found", id.Id)
		}
		return nil, err
	}
	student.Id = id.Id
	log.Printf("Getting - %v\n", &student)
	return &student, nil

}

func (s *server) UpdateStudent(ctx context.Context, student *pb.Student) (*pb.Student, error) {
	/*log.Println("UpdateStudent was invoked")
	s.students[student.Id] = student
	log.Printf("Updating - %v\n", student)
	return student, nil*/

	log.Println("UpdateStudent was invoked")
	_, err := s.db.Exec("update students set name=?, email=?, mobile=? where id=?", student.Name, student.Email, student.Mobile, student.Id)
	if err != nil {
		return nil, err
	}
	log.Printf("Updating - %v\n", student)
	return student, nil

}

func (s *server) DeleteStudent(ctx context.Context, id *pb.StudentId) (*pb.Student, error) {
	/*log.Println("DeleteStudent was invoked")
	deletedStudent := s.students[id.Id]
	delete(s.students, id.Id)
	log.Printf("Deleting - %v\n", deletedStudent)
	return deletedStudent, nil*/

	log.Println("DeleteStudent was invoked")
	var student pb.Student
	err := s.db.QueryRow("select id, name, email, mobile from students where id=?", id.Id).Scan(&student.Id, &student.Name, &student.Email, &student.Mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Student with Id %d not found", id.Id)
		}
		return nil, err
	}

	_, err = s.db.Exec("delete from students where id=?", id.Id)
	if err != nil {
		return nil, err
	}
	student.Id = id.Id
	log.Printf("Deleting - %v\n", &student)
	return &student, nil

}

var address string = "0.0.0.0:50051"

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/studentinfo")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on %s\n", address)
	s := grpc.NewServer()
	//pb.RegisterStudentServiceServer(s, &server{students: make(map[int32]*pb.Student)})
	pb.RegisterStudentServiceServer(s, &server{db: db})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
