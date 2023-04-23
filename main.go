package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// サンプルデータ
	users := []User{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
		{ID: "3", Name: "Charlie"},
	}

	// GraphQLのフィールドを定義
	fields := graphql.Fields{
		"user": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "User",
				Fields: graphql.Fields{
					"id":   &graphql.Field{Type: graphql.String},
					"name": &graphql.Field{Type: graphql.String},
				},
			}),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(string)
				if ok {
					for _, user := range users {
						if user.ID == id {
							return user, nil
						}
					}
				}
				return nil, nil
			},
		},
	}

	// GraphQLのスキーマを定義
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// GraphQLハンドラーを設定
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// エンドポイントを定義
	http.Handle("/graphql", h)
	fmt.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
