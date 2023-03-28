# Go-with-ScyllaDB
Demo app to learn usage of Go driver with ScyllaDB cluster

By playing with this demo app you can learn how to use gocqlx driver to interact with your ScyllaDB cluster, how to select, insert delete data and object mapping with go.

### Step 1
Clone this repo and cd into it. Play around with main.go file to your need. Then to see code in action, first create prepare ScyllaDB.
```
chmode +x ./create-cluster.sh
./create-cluster.sh
```
### Step 2
Wait for about 1 minutes to let the node sync. The open up cql shell on any of the node.
```
podman exec -it scylla-node1 nodetool status
podman exec -it scylla-node1 cqlsh
```

### Step 3
Then create keyspace, use keyspace, create table and insert some data. Run these inside cqlsh prompt.
```
CREATE KEYSPACE catalog WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy','datacenter1' : 3};

USE catalog;

CREATE TABLE mutant_data (
   first_name text,
   last_name text,
   address text,
   picture_location text,
   PRIMARY KEY((first_name, last_name)));

INSERT INTO mutant_data ("first_name","last_name","address","picture_location") VALUES ('Bob','Loblaw','1313 Mockingbird Lane', 'http://www.facebook.com/bobloblaw');
INSERT INTO mutant_data ("first_name","last_name","address","picture_location") VALUES ('Bob','Zemuda','1202 Coffman Lane', 'http://www.facebook.com/bzemuda');
INSERT INTO mutant_data ("first_name","last_name","address","picture_location") VALUES ('Jim','Jeffries','1211 Hollywood Lane', 'http://www.facebook.com/jeffries');
```

### Step 4
Create a docker image for the goapp by using the Dockerfile or pull the prebuilt image from my container repository.
```
podman pull quay.io/csjoy/goscylladb
podman run --network=cluster --rm goscylladb
```
### Step 5
You can see output similar to this.
```
2023/03/28 09:34:54 Creating Cluster
2023/03/28 09:34:54 successfully connected
2023/03/28 09:34:54 Displaying Results
2023/03/28 09:34:54     Bob Loblaw 1313 Mockingbird Lane http://www.facebook.com/bobloblaw
2023/03/28 09:34:54     Jim Jeffries 1211 Hollywood Lane http://www.facebook.com/jeffries
2023/03/28 09:34:54     Bob Zemuda 1202 Coffman Lane http://www.facebook.com/bzemuda
2023/03/28 09:34:54 Inserting Prosenjit...
2023/03/28 09:34:54 Displaying Results
2023/03/28 09:34:54     Bob Loblaw 1313 Mockingbird Lane http://www.facebook.com/bobloblaw
2023/03/28 09:34:54     Jim Jeffries 1211 Hollywood Lane http://www.facebook.com/jeffries
2023/03/28 09:34:54     Prosenjit Joy 1410 Mirpur Dhaka http://www.facebook.com/prosenjit.joy
2023/03/28 09:34:54     Bob Zemuda 1202 Coffman Lane http://www.facebook.com/bzemuda
2023/03/28 09:34:54 Deleting Prosenjit...
2023/03/28 09:34:54 Displaying Results
2023/03/28 09:34:54     Bob Loblaw 1313 Mockingbird Lane http://www.facebook.com/bobloblaw
2023/03/28 09:34:54     Jim Jeffries 1211 Hollywood Lane http://www.facebook.com/jeffries
2023/03/28 09:34:54     Bob Zemuda 1202 Coffman Lane http://www.facebook.com/bzemuda
```
