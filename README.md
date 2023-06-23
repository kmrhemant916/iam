# IAM

* Permission
  * CREATE
  * UPDATE
  * DELETE
  * READ
* Roles
  * Administrator
    * Permission
      * CREATE
      * UPDATE
      * DELETE
      * READ
  * Reader
    * Permission
      * READ
* Organization

```
docker run --name mysql -v $(pwd)/Documents/github/mysql:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=admin -p 3306:3306 -d mysql:latest
```

## FLow

* User signup
  * Username
  * Password
  * Company
* Add user to the Admin
* Invite user
  * Username
  * Default role will be reader
* Add user to group
