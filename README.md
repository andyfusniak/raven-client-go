# Raven Mailer HTTP Client library and CLI tool

Raven Mailer cli tool provides a means to manage templates with the Raven Mailer system.


## Usage

### List templates
```shell
$ raven list templates
```

## Build

In the root directory run make and copy the appropriate `raven` binary to a directory on your path.

```
make
```
## Environment Variables

+ `RAVEN_ENDPOINT` (optional) used during testing to override the compiled in endpoint. e.g. `http://localhost:8080/v1`.
