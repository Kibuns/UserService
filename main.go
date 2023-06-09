package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kibuns/UserService/DAL"
	"github.com/Kibuns/UserService/Models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	fmt.Println("Started UserService")
	// create a channel to receive signals to stop the application
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	// start the goroutine to receive messages from the queue
	go receiveDeleted()
	
	// start the goroutine to handle API requests
	go handleRequests()
	
	// wait for a signal to stop the application
	<-stop
}

//controllers

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage of the Search Service!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAll")
	json.NewEncoder(w).Encode(getAllUsers())
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["user"]
	DAL.DeleteAllOfUser(idParam)
	fmt.Fprintf(w, "deleted everything from user: " + idParam)
}

func storeUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Storing User")
	// parse the request body into a User struct

	var user Models.User
	err := json.NewDecoder(body).Decode(&user)
	fmt.Println(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// Encrypt user.password using bcrypt
	var hashedUserInput = sha256.Sum256([]byte(user.Password))
	var hashedString = hex.EncodeToString(hashedUserInput[:])
	user.Password = hashedString

	// insert the user into the database
	err = DAL.RegisterUser(user)
	if err != nil {
		if err.Error() == "username is not unique" {
			http.Error(w, "username is not unique", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//send message to rabbitMQ queue
	send("user", &user)

	fmt.Fprintln(w, "User stored successfully")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.Use(CORS)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAll)
	myRouter.HandleFunc("/delete/{user}", deleteUser)
	myRouter.HandleFunc("/create", storeUser)
	myRouter.HandleFunc("/get/{user}", getUser)

	log.Fatal(http.ListenAndServe(":9998", myRouter))
}

func getUser(w http.ResponseWriter, r *http.Request){
	var userParam string = mux.Vars(r)["user"]
	json.NewEncoder(w).Encode(DAL.SearchUser(userParam))
}



//service functions

func getAllUsers() (values []primitive.M) {
	result, err := DAL.ReadAllUsers()
	FailOnError(err, "Could not retrieve all users")
	return result
}

// func getUser(query string) (twoots []primitive.M){
// 	twoots, err := DAL.SearchUserByID(query);
// 	if err != nil {
// 		log.Panicf("%s: %s", "could not search for twoots", err)
// 		return
// 	}
// 	return twoots;
// }

// other
// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		//return
	})

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
