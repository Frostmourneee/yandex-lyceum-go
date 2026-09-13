package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func main() {
	ex1()
}

type Person struct {
	Name         string `json:"name"` // `json:"name"` — тег поля. Без него в JSON будет ключ "Name" вместо "name".
	Age          int    `json:"age"`
	Gender       string `json:"gender"` // Поле с тегом `json:"-"` при кодировании json игнорируется
	privateNotes string // Неэкспортируемые поля так же игнорируются
}

func ex1() {
	jsonStr := `{"name": "John", "age": 30, "Gender": "male"}`
	var person Person
	err := json.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		panic(err)
	}
	fmt.Println(person)
}

func SerializeIntSlice(nums []int) ([]byte, error) {
	res, err := json.Marshal(nums)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeserializeStringMap(data string) (map[string]string, error) {
	var m map[string]string
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type Student struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

func DecodeStudentFromReader(r io.Reader) (Student, error) {
	decoder := json.NewDecoder(r)

	var student Student
	err := decoder.Decode(&student)
	if err != nil {
		return Student{}, err
	}
	return student, nil
}

func EncodeStudentsToWriter(w io.Writer, students []Student) error {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(students)
	if err != nil {
		return err
	}
	w.Write(buf.Bytes())
	return nil
}
