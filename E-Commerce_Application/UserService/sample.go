package UserService

/*
To create a user in cassandra

cassandra:
  image: cassandra:latest
  restart: always
  environment:
    CASSANDRA_CLUSTER_NAME: "TestCluster"
    CASSANDRA_NUM_TOKENS: 256
    CASSANDRA_AUTHENTICATOR: "PasswordAuthenticator"
  ports:
    - "9042:9042"


docker exec -it <container_id> cqlsh
CREATE ROLE admin WITH PASSWORD = 'yourpassword' AND SUPERUSER = true AND LOGIN = true;


*/
