package main

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"log"
)

func main() {
	log.Println("Creating Cluster")
	cluster := gocql.NewCluster("scylla-node1", "scylla-node2", "scylla-node3")
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("unable to connect to scylla", err)
	}
	log.Println("successfully connected")
	defer session.Close()

	selectQuery(session)
	insertQuery(session, "Prosenjit", "Joy", "1410 Mirpur Dhaka", "http://www.facebook.com/prosenjit.joy")
	selectQuery(session)
	deleteQuery(session, "Prosenjit", "Joy")
	selectQuery(session)
}

type query struct {
	stmt  string
	names []string
}

type statements struct {
	del query
	ins query
	sel query
}

type MutantData struct {
	FirstName       string `db:"first_name"`
	LastName        string `db:"last_name"`
	Address         string `db:"address"`
	PictureLocation string `db:"picture_location"`
}

var stmts = createStatements()

func createStatements() *statements {
	mutantMetadata := table.Metadata{
		Name:    "catalog.mutant_data",
		Columns: []string{"first_name", "last_name", "address", "picture_location"},
		PartKey: []string{"first_name", "last_name"},
	}
	mutantTable := table.New(mutantMetadata)
	deleteStmt, deleteNames := mutantTable.Delete()
	insertStmt, insertNames := mutantTable.Insert()
	selectStmt, selectNames := qb.Select(mutantMetadata.Name).Columns(mutantMetadata.Columns...).ToCql()
	return &statements{
		del: query{stmt: deleteStmt, names: deleteNames},
		ins: query{stmt: insertStmt, names: insertNames},
		sel: query{stmt: selectStmt, names: selectNames},
	}
}

func deleteQuery(session gocqlx.Session, firstName string, lastName string) {
	log.Println("Deleting " + firstName + "...")
	r := MutantData{
		FirstName: firstName,
		LastName:  lastName,
	}
	err := gocqlx.Session.Query(session, stmts.del.stmt, stmts.del.names).BindStruct(r).ExecRelease()
	if err != nil {
		log.Println("delete catalog.mutant_data", err)
	}
}

func insertQuery(session gocqlx.Session, firstName, lastName, address, pictureLocation string) {
	log.Println("Inserting " + firstName + "...")
	r := MutantData{
		FirstName:       firstName,
		LastName:        lastName,
		Address:         address,
		PictureLocation: pictureLocation,
	}
	err := gocqlx.Session.Query(session, stmts.ins.stmt, stmts.ins.names).BindStruct(r).ExecRelease()
	if err != nil {
		log.Println("insert catalog.mutant_data", err)
	}
}

func selectQuery(session gocqlx.Session) {
	log.Println("Displaying Results")
	var r []MutantData
	err := gocqlx.Session.Query(session, stmts.sel.stmt, stmts.sel.names).SelectRelease(&r)
	if err != nil {
		log.Println("select catalog.mutant_data", err)
		return
	}
	for _, v := range r {
		log.Printf("\t%s %s %s %s\n", v.FirstName, v.LastName, v.Address, v.PictureLocation)
	}
}
