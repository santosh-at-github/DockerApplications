Install Docker-Compose:
  $ sudo curl -L "https://github.com/docker/compose/releases/download/1.25.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
Before executing docker-compose up for this app, create mysql DB Schema:
- Mysql (for db and schema creation):
  $ docker exec -i mysql sh -c 'exec mysql -urootuser -p"RootPass123$"' < mysql.sql
- Mysql login command:
  $ mysql -h localhost -u rootuser -p"$MYSQL_PASSWORD"
