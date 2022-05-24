## How tests the changes
1. `make docker`
2. ```shell
   docker tag eu.gcr.io/sbat-gcr-develop/securebanking/securebanking-test-data-initializer:latest eu.gcr.io/sbat-gcr-release/securebanking/securebanking-test-data-initializer:latest
   ```
3. ```shell
   docker push !$[+TAB]
   ```
   Or
   ```shell
   docker push eu.gcr.io/sbat-gcr-release/securebanking/securebanking-test-data-initializer:latest
   ```
4. Change kubernetes context to `sbat-master-dev`
5. Set the namespace to the developer namespace
6. Run `docker delete pod rs-xxxxxxxx`

>The new rs pod will run the latest image pushed in the step 3.

>Check your changes on the [platform](https://iam.dev.forgerock.financial/platform)