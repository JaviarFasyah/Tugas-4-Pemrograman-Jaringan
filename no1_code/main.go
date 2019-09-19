package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Karyawan a
type Karyawan struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nama   string
	Email  string
	NoTlp  string "no_tlp"
	Alamat string
}

var tmpl = template.Must(template.ParseGlob("view/*"))

const dbname = "mydb"
const collname = "karyawan"
const connlink = "mongodb://localhost:27017"

func index(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	querry, err := coll.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	var re []*Karyawan
	for querry.Next(context.TODO()) {
		var kar Karyawan
		err := querry.Decode(&kar)
		if err != nil {
			log.Fatal(err)
		}
		re = append(re, &kar)
	}
	if err := querry.Err(); err != nil {
		log.Fatal(err)
	}
	querry.Close(context.TODO())
	tmpl.ExecuteTemplate(w, "index", re)
}

func insert(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	if r.Method == "POST" {
		kar := Karyawan{Nama: r.FormValue("nama"), Email: r.FormValue("email"), NoTlp: r.FormValue("notlp"), Alamat: r.FormValue("alamat")}
		_, err := coll.InsertOne(context.TODO(), kar)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Redirect(w, r, "/", 301)
}

func edit(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	var kar Karyawan
	rid := r.URL.Query().Get("id")
	rid = rid[10:34]
	brid, _ := primitive.ObjectIDFromHex(rid)
	err = coll.FindOne(context.TODO(), bson.M{"_id": brid}).Decode(&kar)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "edit", kar)
}

func update(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	if r.Method == "POST" {
		rid := r.FormValue("id")
		rid = rid[10:34]
		brid, _ := primitive.ObjectIDFromHex(rid)
		ukar := bson.M{
			"$set": bson.M{
				"nama":   r.FormValue("nama"),
				"email":  r.FormValue("email"),
				"no_tlp": r.FormValue("notlp"),
				"alamat": r.FormValue("alamat"),
			},
		}
		_, err := coll.UpdateOne(context.TODO(), bson.M{"_id": brid}, ukar)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Redirect(w, r, "/", 301)
}

func del(w http.ResponseWriter, r *http.Request) {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(dbname).Collection(collname)
	rid := r.URL.Query().Get("id")
	rid = rid[10:34]
	brid, _ := primitive.ObjectIDFromHex(rid)
	_, err = coll.DeleteOne(context.TODO(), bson.M{"_id": brid})
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", 301)
}

func new(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "new", nil)
}

func main() {
	clientop := options.Client().ApplyURI(connlink)
	client, err := mongo.Connect(context.TODO(), clientop)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mongodb online...")

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/insert", insert).Methods("POST")
	r.HandleFunc("/new", new)
	r.HandleFunc("/edit", edit).Methods("GET")
	r.HandleFunc("/update", update).Methods("POST")
	r.HandleFunc("/delete", del).Methods("GET")
	http.ListenAndServe(":8000", r)
}
