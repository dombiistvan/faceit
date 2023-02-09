<h2>faceit REST - User CRUD</h2>

<p>before start the application, you can change the environmental variables in <code>mysql.conn.env</code> file if you want, but you don't need to change anything by default</p>

````
MYSQL_ROOT_PASSWORD: gC38a8YrZqq4ZVtD
MYSQL_DATABASE: faceit
MYSQL_USER: faceit
MYSQL_PASSWORD: NootNoot
````

<p>the application will be reachable in your <code>http://localhost:80</code> after you started the app by entering <code>docker compose up</code> command</p>

<ul>
    <li>you can check the api is running by type in your browser address bar <code>http://localhost/test</code>
        <ul>
            <li>if everything is fine, you will see an <code>it's alive!</code> message</li>
        </ul>
    </li>
</ul>

<p>For each request, the expected response is success and the object/list of objects, or if no content needed, empty response body with appropriate http.StatusCode</p>

<h3>Create new User</h3>
Method: <code>POST</code><br />
Content-type: <code>Application/json</code><br />
URL: <code> http://localhost/user </code> <br />
Request-Body: <cite>(example)</cite><br />

````
{
	"first_name": "Firstname",
	"last_name": "Lastname",
	"nickname": "MyCoolNickName",
	"email": "myemail@gmail.com",
	"country": "HU",
	"password": "MySecretPassword_1234"
}
````

Response if succeeded: <cite>(example)</cite><br />
Response Code: <cite>201 - http.StatusCreated</cite><br />

````
{
	"id": "5c188809-1b2c-41b9-b937-543e9efece79",
	"first_name": "Firstname",
	"last_name": "Lastname",
	"nickname": "MyCoolNickName",
	"email": "myemail@gmail.com",
	"country": "HU",
	"CreatedAt": "2023-02-07T16:01:39Z",
	"UpdatedAt": "2023-02-07T16:01:39Z"
}
````

<h3>Update Existing User</h3>
Method: <code>PUT</code><br />
Content-type: <code>Application/json</code><br />
URL format: <code> http://localhost/user/:UUID </code> <br />
URL: <code> http://localhost/user/5c188809-1b2c-41b9-b937-543e9efece79 </code> <br />
Request-Body: <cite>(example)</cite><br />

````
{
	"first_name": "Firstname2",
	"last_name": "Lastname2",
	"email": "myemailnew@gmail.com",
	"country": "UK",
}
````

Response Code: <cite>200 - http.StatusOK</cite><br />
Response if succeeded: <cite>(example)</cite><br />

````
{
	"id": "5c188809-1b2c-41b9-b937-543e9efece79",
	"first_name": "Firstname2",
	"last_name": "Lastname2",
	"nickname": "MyCoolNickName",
	"email": "myemailnew@gmail.com",
	"country": "UK",
	"CreatedAt": "2023-02-07T17:01:39Z",
	"UpdatedAt": "2023-02-07T17:01:39Z"
}
````

<h3>Delete Existing User</h3>
Method: <code>DELETE</code><br />
URL format: <code> http://localhost/user/:UUID </code> <br />
URL example: <code> http://localhost/user/5c188809-1b2c-41b9-b937-543e9efece79 </code> <br />

Response Code: <cite>200 - http.StatusOK</cite><br />
Response Body: <cite>no response body</cite><br />

<h3>Get Filtered OR unfiltered list of users</h3>
Method: <code>GET</code><br />
URL
format: <code> http://localhost/user/list?page=3&limit=10&query_param=value&query_param2=value2....query_paramn=valuen </code> <br />
URL
example: <code> http://localhost/user/list?first_name=Dombi&last_name=istvan&nickname=cool&country=hu&email=gmail&id=5c188809-1b2c&page=2&limit=3 </code> <br />
Available Filter query params: ```first_name, last_name, nickname, email, country, id```<br />
Available System query params: ```page```(default: 1), ```limit```(default: 10)<br />
Request-Body: <cite>no request body</cite><br />

Response Code: <cite>200 - http.StatusOK</cite><br />
Response if succeeded: <cite>(example)</cite><br />

````
[
    {
        "id": "5c188809-1b2c-41b9-b937-543e9efece79",
        "first_name": "Firstname2",
        "last_name": "Lastname2",
        "nickname": "MyCoolNickName",
        "email": "myemailnew@gmail.com",
        "country": "UK",
        "CreatedAt": "2023-02-07T17:01:39Z",
        "UpdatedAt": "2023-02-07T17:01:39Z"
    },
    {
        "id": "68ddbcf8-be80-423e-b3f2-4749f7a4ba16",
        "first_name": "Test",
        "last_name": "Test2",
        "nickname": "Nicktest",
        "email": "Nicktest@gmail.com",
        "country": "US",
        "CreatedAt": "2023-02-07T12:01:39Z",
        "UpdatedAt": "2023-02-07T13:01:39Z"
    },
    {
        "id": "a9a1b605-6a8e-448f-b51d-62c68389101c",
        "first_name": "Kovács",
        "last_name": "András",
        "nickname": "Kovand",
        "email": "kovand@example.com",
        "country": "CA",
        "CreatedAt": "2023-02-07T14:11:39Z",
        "UpdatedAt": "2023-02-07T15:14:39Z"
    }
]
````
