package cli

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/themisir/databoard/config"
	"github.com/themisir/databoard/query"
	"github.com/themisir/databoard/server"
	"gopkg.in/yaml.v2"
)

type App struct {
	queries   map[string]*query.Query
	mutations map[string]*query.Query
	config    *config.DataboardConfig
	routes    []*server.Route
	router    *mux.Router
	db        *sql.DB
}

func NewApp() *App {
	return &App{
		queries:   make(map[string]*query.Query),
		mutations: make(map[string]*query.Query),
		config:    new(config.DataboardConfig),
		routes:    []*server.Route{},
		router:    mux.NewRouter().StrictSlash(true),
	}
}

func (a *App) Run() {
	a.initEnv()
	a.initConfig()
	a.initDb()
	a.initQueries("query", a.config.Queries, a.queries)
	a.initQueries("mutation", a.config.Mutations, a.mutations)
	a.initRoutes()
	a.initRouter()

	defer a.db.Close()
	http.ListenAndServe(":8000", a.router)
}

func (a *App) initEnv() {
	env := os.Getenv("ENVIRONMENT")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env.local." + env)
	if "test" != env {
		godotenv.Load(".local.env")
	}
	godotenv.Load("." + env + ".env")
	godotenv.Load()
}

func (a *App) initDb() {
	driver := a.config.Database.Driver
	connection := a.config.Database.Connection
	if "" == driver {
		driver = os.Getenv("DB_DRIVER")
	}
	if "" == connection {
		connection = os.Getenv("DB_CONN")
	}

	db, err := sql.Open(driver, connection)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	a.db = db
}

func (a *App) initConfig() {
	bytes, err := ioutil.ReadFile("databoard.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}
	if err := yaml.Unmarshal(bytes, a.config); err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}
}

func (a *App) initQueries(n string, from map[string]config.Query, to map[string]*query.Query) {
	for k, v := range from {
		q, err := query.New(k, v.Query)
		if err != nil {
			log.Fatalf("Failed to load %s: %s", n, err)
		}
		for _, param := range v.Parameters {
			switch param.Type {
			case "string":
				q.AddParam(query.String(param.Name, param.Optional))
				break
			case "int":
				q.AddParam(query.Int(param.Name, param.Optional))
				break
			default:
				log.Fatalf("Parameter type '%s' on %s parameter '%s'.'%s' is not supported", param.Type, n, k, param.Name)
				break
			}
		}
		to[k] = q
	}
}

func (a *App) initRouter() {
	for _, route := range a.routes {
		a.router.Handle(route.Path(), route)
	}
}

func (a *App) initRoutes() {
	for _, routeConfig := range a.config.Routes {
		route := server.NewRoute(routeConfig.Path)

		for k, methodConfig := range routeConfig.Methods {
			var delegate server.MethodDelegate
			if queryConfig := methodConfig.Query; queryConfig != nil {
				if query, ok := a.queries[queryConfig.Name]; ok {
					delegate = server.Query(a.db, query, methodConfig.Query.First)
				} else {
					log.Fatalf("Mutation called '%s' is not exists", queryConfig.Name)
				}
			} else if mutationConfig := methodConfig.Mutation; mutationConfig != nil {
				if mutation, ok := a.mutations[mutationConfig.Name]; ok {
					delegate = server.Mutation(a.db, mutation)
				} else {
					log.Fatalf("Mutation called '%s' is not exists", mutationConfig.Name)
				}
			} else {
				log.Fatalf("No handler configured for route (%s) %s", k, routeConfig.Path)
			}

			method := route.AddMethod(k, delegate)
			for paramName, paramConfig := range methodConfig.Parameters {
				err := method.AddParam(paramName, paramConfig.Value, paramConfig.Validation)
				if err != nil {
					log.Fatalf("Failed to parse parameter '%s' value template for route (%s) %s: %s", paramName, k, routeConfig.Path, err)
				}
			}
		}

		a.routes = append(a.routes, route)
	}
}
