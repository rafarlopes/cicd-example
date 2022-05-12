# cicd-example
Sample CI/CD repository with a hello world Go application using GitHub Actions and ArgoCD

## Presentation

Attached to this repository you can find the presentation in the PDF format [here](Presentation.pdf)

## Structure

This repository contains a hello world web app written in Go, just for the sake of presenting something that is deployable.

We also have 2 workflows: CI and CICD. 

- CI runs for branches and pull requests, building our code, testing and linting
- CICD runs the same process as CI + pushes the image to DockerHub, deploys to environments (dev and stage) and then run E2E tests after it.
  - For production, it opens a PR with the changes and waits for the merge.

The workflow uses 3 repositories:
- the app repository (this repo)
- the deployment repository (where manifests are) [here](https://github.com/rafarlopes/example-deployment)
- the e2e tests repository (to validate our application) [here](https://github.com/rafarlopes/cicd-example-e2e-tests)

There's a folder named `actions` under the `.github` with some composite actions for this workflow.
Feel free to reuse them and create your own repository to share them across your workflows.
For simplicity this was done here at the same repository.
The rest of the actions used are public and you can also benefit from it.

### Deployment

The deployment happens based on a commit into the deployment repository in one of the environment branches (development, stage or production).
From that point, ArgoCD comes in and deploys the changes in the Kubernetes cluster.

The CICD workflow creates a `pre-release` in this repository at the same time it creates the pull request at the deployment repository.
There's an additional workflow in the deployment repository at the production branch that publishes this release once the pull request is merged into production.
