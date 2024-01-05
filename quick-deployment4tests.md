## How tests the changes
1. `make docker`
2. ```shell
   docker tag europe-west4-docker.pkg.dev/sbat-gcr-develop/sapig-docker-artifact/securebanking/securebanking-test-data-initializer:latest europe-west4-docker.pkg.dev/sbat-gcr-release/sapig-docker-artifact/securebanking/securebanking-test-data-initializer:latest
   ```
3. ```shell
   docker push !$[+TAB]
   ```
   Or
   ```shell
   docker push europe-west4-docker.pkg.dev/sbat-gcr-release/sapig-docker-artifact/securebanking/securebanking-test-data-initializer:latest
   ```
4. Change kubernetes context to `sbat-master-dev`
5. Set the namespace to the developer namespace
6. Run `docker delete pod rs-xxxxxxxx`

>The new rs pod will run the latest image pushed in the step 3.

>Check your changes on the [platform](https://iam.dev.forgerock.financial/platform)