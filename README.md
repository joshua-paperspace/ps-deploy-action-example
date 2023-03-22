# ps-deploy-action

### Workflow Steps ### 
- Workflow push-main.yaml (in .github/workflows) is listening for a push on the main branch of the repo to run
- When there is a push on the main branch, push-main.yaml will run and create a repo tag (e.g. ps-deploy-action@v0.1.0)
    - push-main.yaml will create a repo tag using the config specified in fixture-release-drafter.yaml 
- Workflow workflow.yaml (in .github/workflows) is listening for a tag in the format of `ps-deploy-action@*` to run
- When that tag is created by push-main.yaml, workflow.yaml is kicked off and will build and push a docker image to a container registry and then create or update a deployment on Paperspace using the GitHub action `paperspace/deploy-action@v1.0`
    - workflow.yaml pulls the following secrets to use in the workflow which are stored as GitHub Action secrets (GitHub repo -> Settings -> Secrets and variables -> Actions -> Secrets)
        - DOCKERHUB_TOKEN
        - DOCKERHUB_USERNAME
        - PAPERSPACE_API_KEY
        - PSBOT_GITHUB_TOKEN
    - workflow.yaml builds the docker image with the Dockerfile in the repo. This Dockerfile copies the application files and model files in the repo into the container
    - GitHub action `paperspace/deploy-action@v1.0` creates or updates a deployment on Paperspace using the deployment spec .paperspace/app.yaml