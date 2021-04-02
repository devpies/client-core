# Core Event-driven Architecture for Devpie Client

## Goal

This is an experimental project for learning.

Devpie Client is a business management tool for performing software development with clients. Features will include 
kanaban or agile style board management and auxiliary services like cost estimation, payments and more. 


### Setup 

#### Requirements
* [Docker Desktop](https://docs.docker.com/desktop/) (Kubernetes enabled)
* [Pgcli](https://www.pgcli.com/install)
* [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
* [Tilt](https://tilt.dev/)

#### Resources Being Used (no setup required)

* [An Auth0 Account](http://auth0.com/)
* [Auth0 Github Deployments Extension](https://auth0.com/docs/extensions/github-deployments)
* [3 Free 20MB Managed Database Services](elephantsql.com)
    
#### Configuration
* \__infra\__ contains the kubernetes infrastructure
* \__auth0\__ contains the auth0 configuration

`__infra__/secrets.yaml` is required. Rename `__infra__/secrets.sample.yaml` 
and provide base64 encoded credentials for all postgres databases properties (_username, password, host etc_.).

For convenience, set the full remote db urls to environment variables to your `.bashrc` or `.zshrc` file. 
You can use them to connect with pgcli for debugging.

```bash
# Use inside your .bashrc or .zshrc file

export MSG_NATS=postgres://username:password@remote-db-host:5432/dbname
export MIC_DB_IDENTITY=postgres://username:password@remote-db-host:5432/dbname
export VIEW_DB_IDENTITY=postgres://username:password@remote-db-host:5432/dbname

alias k=kubectl
```

## Architecture 

This backend uses CQRS and event sourcing sparingly. 
CQRS is not an architecture. You don't use CQRS everywhere.
 
This backend should be used with the [devpie-client-app](https://github.com/ivorscott/devpie-client-app).
It uses [devpie-client-events](https://github.com/ivorscott/devpie-client-common-module) as a shared library to generate 
message interfaces across multiple programming languages, but the Typescript definitions in the events repository are the source of truth.

## How Data Moves Through System Parts

Two architectural models are adopted: _a traditional microservices model_ and 
_an event sourcing model_ driven by CQRS.
[CQRS allows you to scale your writes and reads separately](https://medium.com/@hugo.oliveira.rocha/what-they-dont-tell-you-about-event-sourcing-6afc23c69e9a). For example, the `identity` feature makes strict use CQRS to write data in one shape and read it in one or more other shapes. This introduces eventual consistency and requires the frontend's support in handling eventual consistent data intelligently. 

In the traditional microservices model every microservice has its own database. Within in
the event sourcing model the authoritative state of truth is stored in a single message store (NATS).

In both models, all events are persisted in a message store. In the traditional microservices model,
the message store serves to promote a fault tolerant system. Microservices can have temporary downtime and return without 
the loss of messages. In the event sourcing model, the current state of an entity is achieved through 
projections on an event stream.


![two models](arch.png)

### The Event Sourcing Model
Devpie Client will need to display all kinds of screens to its users. End users send requests to Applications. Applications write messages (commands or events) to the Messaging System in response to those requests. Components pick up those messages, perform their work, and write new messages to the Messaging System. Aggregators observe all this activity and transform these messages into View Data that Applications use at a later time (eventual consistency) to send responses to users.


<details>
<summary>Read more</summary>

### Definitions

#### Applications

- Applications are not microservices.
- An Application is a feature with its own endpoints that accepts user interaction.
- Applications provide immediate responses to user input.

#### Messaging System

- A stateful msg broker plays a central role in entire architecture.
- All state transitions will be stored by NATS Streaming in streams of messages. These state transitions become the authoritative state used to make decisions.
- NATS Streaming is a durable state store as well as a transport mechanism.

#### Components

- Components don't have their own database.
- Components derive authoritative state from a message store using projections.
- Components are small and focused doing one thing well.

#### Aggregators

- Aggregators aggregate state transitions into View Data that Applications use to respond to a user.

#### View Data

- View Data are read-only models derived from state transitions.
- View Data are eventually consistent.
- View Data are not for making decisions.
- View Data are not authoritative state, but derived from authoritative state.
- View Data can be stored in any format or database that makes sense for the Application
</details>


## Developement

Run front and back ends simultaneously. For faster development don't run the [devpie-client-app](https://github.com/ivorscott/devpie-client-app) in a container/pod.

```bash
# devpie-client-app
npm start

# devpie-client-cqrs-core
tilt up
```

### Testing

Navigate to the feature folder to run tests.
```bash
cd identity/application
npm run tests
```

### Debugging
 
#### Inspecting Managed Databases
Provide `pgcli` a remote connection string.
```bash
pgcli $MIC_DB_IDENTITY 

# opts: [ $MSG_NATS | $MIC_DB_IDENTITY | $VIEW_DB_IDENTITY ... ]
```

#### Using PgAdmin 
If you prefer a UI to debug postgres you may use pgadmin. Run a pod instance and then apply port fowarding. To access pgadmin go to `localhost:8888` and enter the credentials below.
```bash
kubectl run pgadmin --env="PGADMIN_DEFAULT_EMAIL=test@example.com" --env="PGADMIN_DEFAULT_PASSWORD=SuperSecret" --image dpage/pgadmin4 
kubectl port-forward pod/pgadmin 8888:80 
```
### Migrations
Microservices and Aggregators should have remote [database services](elephantsql.com).

Migrations exist under the following paths:

- `<feature>/microservice/migrations`
- `<feature>/aggregator/migrations` 

#### Migration Flow
1. move to a feature's `microservice` or `aggregator`
2. create a `migration`
3. add sql for `up` and `down` migration files
4. `tag` an image containing the latest migrations
5. `push` image to registry

<details>
<summary>View example</summary>
<br>

```bash
cd identity/microservice

migrate create -ext sql -dir migrations -seq create_table 

docker build -t devpies/mic-db-identity-migration:v000001 ./migrations

docker push devpies/mic-db-identity-migration:v000001  
```
</details>

Then apply the latest migration with `initContainers`
<details>
<summary>View example</summary>
<br>

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mic-identity-depl
spec:
  selector:
    matchLabels:
      app: mic-identity
  template:
    metadata:
      labels:
        app: mic-identity
    spec:
      containers:
        - image: devpies/client-mic-identity
          name: mic-identity
          env:
            - name: POSTGRES_DB
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-database-name
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-password
            - name: POSTGRES_HOST
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-host
# ============================================
#  Init containers are specialized containers
#  that run before app containers in a Pod.
# ============================================
      initContainers:
        - name: schema-migration
          image: devpies/mic-db-identity-migration:v000001
          env:
            - name: DB_URL
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-url
          command: ["migrate"]
          args: ["-path", "/migrations", "-verbose", "-database", "$(DB_URL)", "up"]
```
</details>

Learn more about migrate cli [here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md). 
