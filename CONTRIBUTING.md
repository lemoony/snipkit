# Contributing

Contributions to the project are highly welcome! There are several ways to help out, e.g.:

- Create an issue on GitHub, if you have found a bug
- Create a feature request, if you have something on your mind
- Create a pull request, if you
    - have fixed something
    - have added a new feature
    - want to add your own theme
- Contribute to the documentation

## Local development

### Setting up a dev environment

Setting up a test environment involves the following steps:

* Install [go](https://go.dev/doc/install)
* Install [pre-commit](https://pre-commit.com/)
* Run `pre-commit install`
* For working on the documentation:
    * Install [Node.js](https://nodejs.org/) (v18+)
    * Run `npm install` in the project root

After this, you'll be able to test any change. Alternatively, you can open the project via Gitpod: 

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/lemoony/snipkit)

### Commands

Check if everything works as expected:

```bash
make ci 
```

This command will run all tests as well as the linter to check if there are any issues.

During development, the following commands may also be beneficial to you:

```bash
make build # Build the binary files
make test # Run all tests
make lint # Run the linter to detect any issues
make mocks # (Re-)generate the mock files
pre-commit run --all-files # Run all pre-commit hooks manually
npm run docs:dev # Start local documentation server
npm run docs:build # Build documentation for production
```

## Features and bugs

Please file feature requests and bugs at the [issue tracker][tracker].

[tracker]: https://github.com/lemoony/snipkit/issues

## Submitting Changes

Push your changes to a topic branch in your fork of the repository. Submit a pull request to the repository on github, with
the correct target branch.
