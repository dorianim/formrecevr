<h1 align="center">
    Formrecevr
</h1>

<p align="center">
    <a href="https://www.gnu.org/licenses/agpl-3.0">
        <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg" />
    </a>
    <a href="https://github.com/dorianim/formrecevr/actions/workflows/tests.yml">
        <img src="https://github.com/dorianim/formrecevr/actions/workflows/tests.yml/badge.svg" alt="Badge tests">
    </a>
    <a href="https://hub.docker.com/r/dorianim/formrecevr">
        <img src="https://img.shields.io/docker/pulls/dorianim/formrecevr.svg" alt="Docker pulls" />
    </a>
    <a href="https://goreportcard.com/report/github.com/dorianim/formrecevr">
        <img src="https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat" alt="Go report" />
    </a>
    <a href="https://codecov.io/gh/dorianim/formrecevr">
        <img src="https://codecov.io/gh/dorianim/formrecevr/branch/main/graph/badge.svg?token=LJLPYEELOP"/>
    </a>
    <a href="https://codeclimate.com/github/dorianim/formrecevr/maintainability">
        <img src="https://api.codeclimate.com/v1/badges/72af8e8e7fa64952c338/maintainability" />
    </a>
</p>

Formrecevr (pronunced "Form receiver") is a simple and lightweight from receiver backend primarily designed for (but not limited to) static websites. It is inspired by [formspree.io](formspree.io) but it is simpler, self-hosted and doesn't have a frontend.

# Features
- **Easy setup:** Just create your docker-compose.yml and you're good to go!
- **Powerful templating:** Use Go templating (just like in Hugo) to compose customized content
- **Shoutrrr integration:** You can use a [wide variety of services](https://containrrr.dev/shoutrrr/v0.5/services/overview/) to reveive your form submissions
- **Flexible delivery:** Thanks to the use of shoutrrr and templating, you can use data from the form as a receiver

# Setup
The official installation method is using Docker:
1. Create a folder for installation:
    ```bash
    mkdir /opt/formrecevr && cd /opt/formrecevr
    ```
2. Create the file docker-compose.yml with this content:
    ```yaml
    version: "3.7"
    services:
      formrecevr:
        image: dorianim/formrecevr
        restart: always
        ports:
          - 5081:8088
        volumes:
          - ./config:/config
    ```
3. Adjust the port (default `5081`) to your needs
4. Start the formrecevr:
    ```bash
    docker-compose up -d
    ```
5. Done! You can reach your formrecevr on `localhost:80`
6. [Adjust your `config.yml`](#configuration) in `/opt/formrecevr/config/config.yml`
7. [OPTIONAL] Adjust your templates in `/opt/formrecevr/config/templates`
8. [OPTIONAL] To setup ssl/https, please use a reverse proxy like nginx

# Configuration
The configuration is stored in `/config/config.yml` in the container by default. It is reloaded live, so changes do not require a container restart to become effective.  
A fully populated config could look like this:
```yaml
forms:
  - id: "Example"
    enabled: true
    targets:
    - enabled: true
      template: default.html
      shoutrrrurl: 'telegram://someToken@telegram/?channels={{ .params.someChannel }}'
      params:
        someKey: someChannel
    - enabled: false
      template: email.html
      shoutrrrurl: 'smtp://username:password@host:port/?from=fromAddress&to=recipient1&to={{ .mailFromFrom }}'
      
  - id: "Example2"
    enabled: true
    targets:
    - enabled: true
      template: default.html
      shoutrrrurl: 'slack://token:token@channel/'
listen:
  host: 0.0.0.0
  port: 8088
```
### Fields
- `listen` contains settings for the webserver
- `forms` contains a list of [form blocks](#the-form-block)

#### The form block
- `id`: The ID of the form. I suggest using something like `uuidgen` to create a random ID
- `enabled`: If the form is enabled. It will not work when this is set to false
- `targets`: Contains a list of [target blocks](#the-target-block)

#### The target block
- `enabled`: If the target is enabled. It will be skipped when set to false
- `template`: The name of the template which is used for the body. See [templates](#templates)
- `shoutrrrurl`: The [shoutrrr](https://containrrr.dev/shoutrrr/v0.5/) URL to use for form submissions. Additional information on how to get the URLs can be found [here](https://containrrr.dev/shoutrrr/v0.5/services/overview/)
- `params`: Additional parameters which are also passed to the templating engine as `.params.<key>`. This can be useful when using the same template for multiple targets.

# Templates
Templates have to be stored in `/config/templates` by default. 
The templates are processed by the go templating engine and have access to all of its functionality, like range loops and if conditions.

### Variables
Templates can use all submited form fields in the root context. In addition to that, they have access to the params defined for their target inside the `.params` key.  
For testing, you can use `{{ . }}` to see all avaiable data.

### Functions
There are two additional functions wich can be used:
- `join`: Joins a list of strings with a delimiter, eg. `{{ join .someListParameter "," }}`
- `print`: Prints one or more strings or string lists without any delimiter, eg. `{{ print .someString .someStringList }}`

# Routes / API
There are currently two routes available
- `/api/v1/healthcheck` (GET): Should return 200 if everything is OK
- `/f/:formID` (POST): Accepts form data for configured formIDs. It currently supports `application/x-www-form-urlencoded` and `multipart/form-data`

# Implementation
I recommend to use the javascript [`fetch()` API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) to submit your form. An example of this can be found in [`examples/form-xhr-fetch.html`](https://github.com/dorianim/formrecevr/blob/main/examples/form-xhr-fetch.html)
