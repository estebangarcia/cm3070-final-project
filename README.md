# CM3070 Final Project

This repository contains the backend code for the module CM3070 Final Project of the University of London BSc Computer Science degree

The UI is available on a different github repository but also included as a submodule in this one.

## How to run?

You ideally need an [ngrok account](https://ngrok.com/docs/getting-started/) to run an https tunnel to the local server, this is needed for using the docker client as it requires registries to be exposed via https.

1. Create the `.env` file with the data specified in the report's appendix.
2. Install [Docker](https://www.docker.com/) in your computer

3. Run:
```bash
$ docker-compose up
```

4. Apply migrations
```bash
$ make migrate
```

5. Run the UI


## Project structure

- `main.go` - This is the main file that exposes the cli interface for starting server and workers
- `cmd/` - Contains the server and workers commands
- `pkg/` - Contains all the application code
    - `api/` - Contains all API-related handler code
    - `config/` - Contains the configurations structs
    - `helpers/` - Contains helper code to initialize the different AWS SDK clients needed and some small utilities
    - `middleware/` - Contains all middleware code
    - `oci_models/` - Contains structs for mapping JSON requests to OCI manifest structs
    - `repositories/` - Contains all the database logic
        - `ent/` - For the most part is automatically generated code by EntGo
            - `schema/` - Contains the code schema used to generate the entgo assets
            - `migrate/migrations/` - Contains the AtlasGO migrations
    - `requests/` - Contains structs to bind for request validations
    - `responses/` - Contains structs to map to JSON responses
    - `templates/` - Contains HTML and Email templates
    - `workers/` - Contains the scanner and user signup worker logic
