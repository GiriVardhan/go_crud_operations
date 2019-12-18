# go-crud-operations

This GO web application performs CRUD operation on Cassandra DB.

To run this application you should have already installed GO and setup Cassandra on docker.

Create the below keyspace and emps tables on Cassandra db.
Create keyspace in cassandra DB

cqlsh> CREATE KEYSPACE clusterbd WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
Create table in the keyspace created above

cqlsh> CREATE TABLE emps ( empid text PRIMARY KEY, first_name text, last_name text, age int );

# To create docker image run the below commands from project root directory
$ go get github.com/gocql/gocql
$ go get github.com/gorilla/mux

# Command to build GO apllication
$ go build
 
# Get the executable generated from the above command and run below command to create docer image
docker build -t executable .

# Please update IP address in the database connection section accordingly.
cluster := gocql.NewCluster("172.18.0.3")


# Setting Cassandra DB on Docker

#Below command pulls the latest Docker Cassandra Image
$ docker pull cassandra

#To check the Docker Image run the below command
$ docker images

#Below command creates a Cassandra container from the Cassandra Image
$ docker run --name cassandraDB -d cassandra:latest

#Below command will open interactive termial for Cassandra cql to create Keyspace, table and sample data.
$ docker exec -it cassandraDB bash 
cassandraDB@<continerID>:~$ cqlsh
cqlsh> 
  # Create keyspace in cassandra DB
  cqlsh> CREATE KEYSPACE clusterdb
  WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
  
  # Create table in the keyspace created above
  cqlsh> CREATE TABLE emps (
  empid text PRIMARY KEY,
  first_name text,
  last_name text,
  age int
  );

  # Insert some sample data into table emps
  cqlsh> INSERT INTO testdb.emps (empid, first_name, last_name, age) 
  VALUES ('EMP1','RAVI','T',40);

  cqlsh> INSERT INTO testdb.emps (empid, first_name, last_name, age) 
  VALUES ('EMP2','GIRI','K',36);

  cqlsh> INSERT INTO testdb.emps (empid, first_name, last_name, age) 
  VALUES ('EMP3','VENU','K',34);

  cqlsh> INSERT INTO testdb.emps (empid, first_name, last_name, age) 
  VALUES ('EMP4','HYMA','P',35);

  cqlsh> INSERT INTO testdb.emps (empid, first_name, last_name, age) 
  VALUES ('EMP5','EMMANUEL','V',37);


 
