package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Task struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	DueDate    string `json:"due_date"`
	CategoryID int    `json:"category_id"`
	IsDone     bool   `json:"is_done"`
	UserID     int    `json:"user_id"`
}

type Category struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Color  string `json:"color"`
	UserID int    `json:"user_id"`
}

var (
	userStorage       []User
	taskStorage       []Task
	categoryStorage   []Category
	authenticatedUser *User

	serializationMode string
)

const (
	userStoragePath = "user.txt"
	ManDarAvardi    = "mandaravardi"
	Json            = "json"
)

func main() {

	fmt.Println("Hello, Go!")
	command := flag.String("command", "no command", "command to run")
	serialize := flag.String("serialize", Json, "serialization mode")
	flag.Parse()

	fmt.Println("switch", *serialize)
	switch *serialize {
	case ManDarAvardi:
		serializationMode = ManDarAvardi
	default:
		serializationMode = Json
	}
	loadUserStorageFromFile()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		runCommand(*command)
		fmt.Println("Enter a command:")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()

		if authenticatedUser == nil {
			return
		}
	}
	switch command {
	case "create-task":
		createTask()
	case "create-cat":
		createCategory()
	case "login":
		login()
	case "register-user":
		registerUser()
	case "list-task":
		listTask()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Invalid command:", command)
	}

	fmt.Println("User Storage:", userStorage)
}

func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var name, duedate, category string

	fmt.Println("Please enter the task title:")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("Please enter the task category:")
	scanner.Scan()
	category = scanner.Text()
	categoryId, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("error in category,%v\n", err)
		return
	}
	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryId && c.UserID == authenticatedUser.ID {
			isFound = true
			break
		}
	}
	if !isFound {
		fmt.Printf("category is not valid")
		return
	}
	fmt.Println("Please enter the task due date:")
	scanner.Scan()
	duedate = scanner.Text()

	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      name,
		DueDate:    duedate,
		CategoryID: categoryId,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}
	taskStorage = append(taskStorage, task)
	fmt.Println("Task created:", name, category, duedate)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("Please enter the category title:")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Please enter the category color:")
	scanner.Scan()
	color = scanner.Text()

	category := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, category)

	fmt.Println("Category created:", title, color)
}

func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("Please enter your email for register:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Please enter your password :")
	scanner.Scan()
	password = scanner.Text()

	user := User{
		ID:       len(userStorage) + 1,
		Email:    email,
		Password: password,
	}

	writeToFile(user)
}

func login() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("Please enter your email:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Please enter your password:")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if user.Email == email {
			if user.Password == password {
				authenticatedUser = &user
				break
			} else {
				fmt.Println("Incorrect password")
			}
		}
	}

	if authenticatedUser == nil {
		fmt.Println("User not found")
	}
}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println("Task:")
			fmt.Println(task)
		}
	}
}

func loadUserStorageFromFile() {
	file, err := os.Open(userStoragePath)
	if err != nil {
		fmt.Println("Can't open file:", err)
		return
	}
	defer file.Close()

	var data = make([]byte, 10240)
	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("Can't read from file:", oErr)
		return
	}
	var dataStr = string(data)

	userSlice := strings.Split(dataStr, "\n")
	var userStruckt = User{}

	for _, u := range userSlice {

		if u == "" {
			continue
		}

		switch serializationMode {
		case ManDarAvardi:
			userStruckt, err = deserilizeFromManDaravardi(u)
			if err != nil {
				fmt.Println("cant deeserilize user record ")
				return
			}
		case Json:
			if u[0] != '{' && u[len(u)-1] != '}' {
				return User{}, errors.New("user string is empty")
			}
			uErr := json.Unmarshal([]byte(u), userStruckt)
			if uErr != nil {
				fmt.Println("cant deeserilize user record ", uErr)
				return
			}

		}

		userStorage = append(userStorage, userStruckt)
	}
}

func writeToFile(user User) {
	file, err := os.OpenFile(userStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't create user.txt:", err)
		return
	}
	defer file.Close()

	var data []byte
	if serializationMode == ManDarAvardi {
		data = []byte(fmt.Sprintf("id:%d,email:%s,password:%s\n", user.ID, user.Email, user.Password))
	} else if serializationMode == Json {
		data, err = json.Marshal(user)
		if err != nil {
			fmt.Println("Can't marshal user struct:", err)
			return
		}
		data = append(data, '\n')
	} else {
		fmt.Println("Invalid serialization mode")
		return
	}
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Can't write to file:", err)
	}
}
func deserilizeFromManDaravardi(userStr string) (User, error) {

	userFields := strings.Split(userStr, ",")
	var user = User{}
	for _, field := range userFields {
		values := strings.Split(field, ":")
		if len(values) != 2 {
			continue
		}
		fieldName := values[0]
		fieldValue := values[1]
		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				fmt.Println("Conversion error:", err)
				return User{}, errors.New(" strconv errors ")
			}
			user.ID = id
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue
		}
	}
	return user, nil
}
