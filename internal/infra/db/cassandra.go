package db

import (
	"github.com/gocql/gocql"
)

func NewCassandraSession() *gocql.Session {
	cluster := gocql.NewCluster("cassandra1")
	cluster.Keyspace = "catalog"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}
