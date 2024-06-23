# ReviewBot

A microservice responsible for initiating reviews for the products which are delivered to customer.
The customer can ask details about the product as well.

## Description

The microservice is used to manage reviews. It consists of following APIs -

- [POST] /converse - handle conversation with a bot. The bot can receive message from system(FE) or user.
- [POST] /endconverse - end conversation with a bot.

The microservice is written using clean architecture which consists of following 3 main layers -

- api - the layer is used to communicate with the application. The new APIs like grpc or graphQL can be implemented in this layer by keeping other layers as it is.
- app - core business logic which is independent of any external driver as well as api or db.
- infra - infra layer is responsible for implementation interfaces responsible for 3rd party services and DB.
- domain - domain layer consists of entities, value objects and domain services.

## Usage

### prerequisite

- go v1.22.0(for running UTs)
- docker
- docker compose

### commands

- test - to run unit test cases (app should always have 100% test coverage)  
  `make test`
- start - to start the microservice  
  `make start`  
  visit [doc link](http://localhost:5174/api/v1/swagger/index.html) in the browser
- stop - to stop the microservice  
  `make stop` or `CTR+c`

Note - go installation is not needed, if you just want to run the microservice, use command `make start`

## Improvements

- Some interfaces are just example and not real communication with 3rd party services.
- Add unit tests for all the layers.
