<%=licenseText%>
package util

import cli "gopkg.in/urfave/cli.v1"

// MongoOptions are the commonly used options when connecting to a MongoDB
// server.
type MongoOptions struct {
	MongoURI      string
	MongoDatabase string
}

// ParseMongoOptions fetches the values from urfave/cli Context and returns
// them as a MongoOptions. Uses the names as specified in MongoFlags.
func ParseMongoOptions(c *cli.Context) *MongoOptions {
	return &MongoOptions{
		MongoURI:      c.String("mongo"),
		MongoDatabase: c.String("mongo-database"),
	}
}

// MongoFlags returns the flags that will be used by ParseMongoOptions.
// defaultDatabase will be used for the --mongo-database flag.
func MongoFlags(defaultDatabase string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "mongo",
			Usage:  "MongoDB connection string",
			Value:  "mongodb://localhost:27017",
			EnvVar: "MONGODB_URI",
		},
		cli.StringFlag{
			Name:   "mongo-database",
			Usage:  "MongoDB Database",
			Value:  defaultDatabase,
			EnvVar: "MONGODB_DATABASE",
		},
	}
}
