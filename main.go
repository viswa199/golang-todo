package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"todo/stringif"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todo struct {
	Todos string
}

var (
	client  *mongo.Client
	err     error
	results []*todo
)

func main() {
	router := mux.NewRouter()
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Disconnect(ctx)
	router.HandleFunc("/", Home)
	router.HandleFunc("/todos", Todos)
	router.HandleFunc("/delete/{_id}", DeleteTodo)
	http.ListenAndServe(":8080", router)
}

func Home(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tmpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		log.Fatal(err)
		return
	}
	GoCollection := client.Database("golang").Collection("todos")
	FindOptions := options.Find()
	cursor, err := GoCollection.Find(ctx, bson.D{}, FindOptions)
	if err != nil {
		log.Fatal(err)
		http.Error(rw, "unable to find documents in database", http.StatusForbidden)
		return
	}
	var results []bson.M
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Fatal(err)
		return
	}
	tmpl.Execute(rw, results)
}

func Todos(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := todo{}
	t.Todos = r.Form.Get("todo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	GoDatabase := client.Database("golang")
	GoCollection := GoDatabase.Collection("todos")
	result, err := GoCollection.InsertOne(ctx, bson.D{
		{Key: "Todo", Value: t.Todos},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(result.InsertedID)
	http.Redirect(rw,r,"/",301)
}

func DeleteTodo(rw http.ResponseWriter, r *http.Request) {
	id:=stringif.Substrings(r.URL.Path)
	fmt.Println(id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	GoCollection := client.Database("golang").Collection("todos")
	_, err := GoCollection.DeleteOne(ctx, bson.D{{Key: "Todo", Value: id}})
	if err != nil {
		log.Fatal(err)
		return
	}
	http.Redirect(rw,r,"/",302)
}