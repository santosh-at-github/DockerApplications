This is a sample app which is designed to run on an EC2 instance and show some basic Info (such as ins id, container id, IP, client visit count, etc.) about the instance and container it is running from.
Before running this APP, create a temporary mysql DB with same user and password, set needed environment variable and create mysql DB and Schema using below command (the same volume will be used by the app after deployment):
```
$ docker exec -i mysql sh -c 'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD"' < mysql.sql
```
