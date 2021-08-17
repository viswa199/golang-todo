package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todo struct{
	todos string
}

var (
	client *mongo.Client
	err    error
	results []*todo
)

func main() {
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
	http.HandleFunc("/", Home)
	http.HandleFunc("/todos",Todos)
	http.ListenAndServe(":8080",nil)
}

func Home(rw http.ResponseWriter, r *http.Request) {
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	GoCollection:=client.Database("golang").Collection("todos")
	FindOptions:=options.Find()
	cursor,err:=GoCollection.Find(ctx,bson.D{},FindOptions)
	if err!=nil{
		log.Fatal(err)
		http.Error(rw,"unable to find documents in database",http.StatusForbidden)
		return 
	}
	for cursor.Next(ctx){
		var data *todo
		err:=cursor.Decode(&data)
		if err!=nil{
			log.Fatal(err)
			http.Error(rw,"Unable to recover data from database",http.StatusInternalServerError)
			return
		}
		results=append(results, data)
	}
	tmpl.Execute(rw, results)
}

func Todos(rw http.ResponseWriter,r *http.Request){
	r.ParseForm()
	t:=todo{}
	t.todos=r.Form.Get("todo")
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	GoDatabase:=client.Database("golang")
	GoCollection:=GoDatabase.Collection("todos")
	result,err:=GoCollection.InsertOne(ctx,bson.D{
		{Key:"todo",Value:t.todos},
	})
	if err!=nil{
		log.Fatal(err)
		return 
	}
	fmt.Println(result.InsertedID)
}