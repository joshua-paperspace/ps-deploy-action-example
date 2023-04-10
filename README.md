# ps-deploy-action-example

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


### Changes ### 

The items in the template scripts that would need to change for each user are:

- Set the GitHub action secrets for the following variables (GitHub repo -> Settings -> Secrets and variables -> Actions -> Secrets):
    - DOCKERHUB_TOKEN (Access key or token to authorize action on the Docker Hub repository)
    - DOCKERHUB_USERNAME (Username to authorize action on the Docker Hub repository)
    - PAPERSPACE_API_KEY (Paperspace API Key with access to the team/project specified)
    - PSBOT_GITHUB_TOKEN (GitHub access token with appropriate permissions)
- .github/fixture-release-drafter.yaml
    - tag template (lines 1, 2, 3) to be named appropriately for the use case (i.e. ps-deploy-action-example@) 
- .github/workflows/workflow.yaml
    - tag name (line 5) to match the tag name set in fixture-release-drafter.yaml (i.e. ps-deploy-action-example@*)
    - tag name (line 15) to match tag name on line 5 (i.e. ps-deploy-action-example@)
    - image name (line 33) of the desired repository in the container registry (e.g. Docker Hub) being used (i.e. joshuapaperspace/fastapi-resnet)
    - projectId (line 41) of the user's project in Paperspace (i.e. p5rlnw4tcga)
    - image name (line 42) that matches the image name on line 33 (i.e. joshuapaperspace/fastapi-resnet)
