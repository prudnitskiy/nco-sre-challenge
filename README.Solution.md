# NCO test assignment
## Components
This solution contains 3 parts:
- the application to containerise, root folder.
- Kubernetes deployment, see the terraform and chart folders
- log scan script (coding assessment), see alertscan folder

## Webapp
I usually try to be as non-intrusive as possible, but this case required some changes in requirements, code and tests to have the application running.
Changes made:
- Requirements cleaned up from dependencies not in use
- Some of the dependencies were updated to satisfy requirements
- Tests and code updated to have the application working. I started from the code and fixed tests according to the code as I understand the application and code (e.g. opposite to TDD).

Usually, I would ask developers to justify their decisions instead of making changes to avoid accidentally breaking the code I do not understand.

### Running the app
To run the code, please use the docker-compose command: `docker compose up --build`. The application will be reachable through http://127.0.0.1:8000.
Please avoid running code without Docker because the container contains dependencies, and this is the way the application runs in production. Container usage allows for repeatable builds.

### Kubernetes deployment
Production setup is running in Kubernetes. The Helm chart is used to create the necessary manifests. Kubernetes deployment is done via Terraform. To provide better reliability, 2 replicas are deployed. Other reasonable defaults are provided as well:
- Container uses a non-privileged user to run
- The root filesystem in the container is read-only
- Resource capping is enabled for the deployment

You can access this deployment with this URL: https://nco.prudnitskiy.pro

Problems and limitations:
- The container contains unit tests and test dependencies. Normal production setup should not have debug-related dependencies.
- The deployment doesn't have any metrics as the application doesn't support metrics scraping. External container metrics are available as part of the environment observability solution.
- Current deployments don't have a certificate definition, and the default certificate is used.
- Ingress doesn't have any authentication enabled.

## CI
### Code validations
CI runs a set of tests on each push to ensure code is correct:
- Python: run unit tests
- Dockerfile: validate + check for best practices (using Checkov)
- terraform: validate, check formatting, check for best practices (using Checkov)
- helm: template validation, check for best practices (using Checkov). Helm security check validation allowed to fail; this behaviour should be changed for production-grade infrastructures.

### Image building
CI builds an image on MR to master or version tag creation. This is a reasonable default to provide developers the ability to check their containers but avoid extensive registry usage.
An image security scan is performed before pushing to ensure the image doesn't have any security vulnerabilities (using trivy). The image is published to the GitHub artefact registry. As this repository is public, the image is also available without any authentication. To provide better security, SBOM is included in each build.

Issues and limitations:
- Trivy DB is not cached, causing VLDB to download each run.
- The image contains tests and unnecessary dependencies.
- The image was built twice (once for the scan and once to publish it). It may be a problem from a security perspective, and also increases build time and resources.
- The image is not signed because of a lack of necessary infrastructure. There is no way to prove image provenance according to SLSA.

### Deployment to Kubernetes
CD pipeline deploys the application to Kubernetes. The intended approach:
- on MR to master - deploy to staging env.
- On release tag - deploy to production.

Each deployment consists of 3 steps:
- terraform init
- terraform plan
- terraform apply as a separate step

Terraform code validation is included to push check, and there is no way to get an MR with new code without doing a push first.

Helm is used to manage deployment manifests. Terraform is used to deploy Helm to Kubernetes cluster(s)

Current issues and limitations:
- The plan step should have access to the old version already deployed. Without it, Terraform will plan to update the helm chart on each MR.
- Terraform uses Azure to store the state file.
  - OIDC federation is enabled to make this access secure. However, Azure login action is required for each Terraform interaction.
  - Azure service principal deployment is out of the scope for this task, but it is important for a secure setup.
- It is impossible to validate the current helm setup before deployment using Chekov/Kubebench because values provided by terraform on plan and not reachable outside of the plan.
- Current Kubernetes access uses a static configuration file, which is a security risk. Even the CI has a dedicated service account with limited permissions, static credentials should not be used.
- The Terraform plan file is saved as an artefact to be passed to the apply step. This may cause a security implication and requires security alignment.
- Rollback procedure requires access to the previous version deployed. Neither the previous git tag nor `helm rollback` procedures will work for this case. However, this setup is relatively stable as the old version will be available, and the pipeline will fail in case of a failed deployment.
- As the container doesn't have a signature, provenance validation is not available on deploy.
- The application doesn't have an internal metrics endpoint. Only cluster-wide app metrics are available.

# Coding assignment

The alert scan script takes a JSON file as input and provides some output based on alert parsing. The script is written using golang, with as few external dependencies as possible.

There are no clear requirements for script output, so the script provides some minimal stats. This can be changed relatively easily.

Current issues and limitations:
- The script is not optimised and may cause excessive memory use. For example, the alerts file is loaded to memory during unmarshalling.
- No tests included, making changes potentially dangerous.
- No CI included for the script as there are no tests.
- `calculateWeightedPriority` function included as a requirement but not used anywhere. The design of this function is a bit intuitive and requires justification on both the concept and the coefficients.

# AI assistance used
- Claude 3 "Sonnet": typing hints (as a bundle to the editor)
- Phind 70B (as a part of Phind search engine): search engine, research.

No code was written using AI "as is".
