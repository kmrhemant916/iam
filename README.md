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

## Flow

* User root signup flow
  * Username
  * Password
  * Company
* Normal user signin
  * Check if it's a first login
* Add user to the Administrator group
* Invite user
  * Username
  * Default role will be reader
  * Password
* Add user to group
