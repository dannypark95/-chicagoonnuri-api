# Chicago Onnuri API

**Chicago Onnuri Church API** is a Golang application that provides RESTful APIs to power the frontend church website. It manages authentication, allowing users to log in securely with JWTs. The app provides file management functionalities by interacting with AWS S3, allowing the upload, deletion, and listing of PDF files, specifically in the 'weekly_bulletin' folder of the S3 bucket. The application is designed for easy deployment on Heroku.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go version 1.16 or newer. You can download it from [here](https://golang.org/dl/)
- AWS Account with access to S3 service. You can create an account [here](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
- MongoDB Database. If you do not have a MongoDB setup, you can set one up on MongoDB Atlas [here](https://www.mongodb.com/cloud/atlas/register)

### Installation

1. Clone the repository:

```
git clone https://github.com/dannypark95/ChicagoOnnuri.git
```

2. Install the dependencies:

```shell
go mod download
```

4. Create a .env file in the project root and fill in the following details:

```shell
DATABASE_URL=<your-mongodb-connection-string>
DATABASE_NAME=<your-database-name>
JWT_SECRET=<your-secret-key>
AWS_ACCESS_KEY_ID=<your-aws-access-key>
AWS_SECRET_ACCESS_KEY=<your-aws-secret-access-key>
```

5. Run the server:

```shell
go run main.go
```

The server will start on the default port `80`, or on the port specified in your `PORT` environment variable.

## Usage

You can test the APIs using a tool such as [Postman](https://www.postman.com/) or [curl](https://curl.haxx.se/).

The available endpoints are:

- `/jubo`: (GET)
- `/login`: (POST)
- `/listPDF`: (GET)
- `/setLiveJubo`: (POST)
- `/pdf`: (POST, DELETE)

## Frontend

The frontend of the Chicago Onnuri Church website that interacts with this API is hosted [here](https://www.chicagoonnuri.com/).

## License

This project is licensed under the MIT License
