# MyGram API Endpoint

## Users

### Register a new user

- Method: POST
- Endpoint: /users/register

### Login user

- Method: POST
- Endpoint: /users/login

### Update current user information

- Method: PUT
- Endpoint: /users

### Delete current user

- Method: DELETE
- Endpoint: /users

## Photos

### Create a new photo

- Method: POST
- Endpoint: /photos

### Get all photos from all users

- Method: GET
- Endpoint: /photos

### Get photo by ID

- Method: GET
- Endpoint: /photos/{photoId}

### Update photo by ID

- Method: PUT
- Endpoint: /photos/{photoId}

### Delete photo by ID

- Method: DELETE
- Endpoint: /photos/{photoId}

## Comments

### Create a new comment

- Method: POST
- Endpoint: /comments

### Get all comments from all users

- Method: GET
- Endpoint: /comments

### Get comment by ID

- Method: GET
- Endpoint: /comments/{commentId}

### Update comment by ID

- Method: PUT
- Endpoint: /comments/{commentId}

### Delete comment by ID

- Method: DELETE
- Endpoint: /comments/{commentId}

## Social Medias

### Create a new social media entry

- Method: POST
- Endpoint: /socialmedias

### Get all social media entries only from logged-in user

- Method: GET
- Endpoint: /socialmedias

### Get specific social media entries only from logged-in user

- Method: GET
- Endpoint: /socialmedias/{socialMediaId}

### Update social media entry by ID

- Method: PUT
- Endpoint: /socialmedias/{socialMediaId}

### Delete social media entry by ID

- Method: DELETE
- Endpoint: /socialmedias/{socialMediaId}
