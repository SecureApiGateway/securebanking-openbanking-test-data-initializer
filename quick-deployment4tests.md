## How tests the changes
1. `make docker`
1. ```shell
   docker tag eu.gcr.io/sbat-gcr-develop/securebanking/securebanking-test-data-initializer:latest eu.gcr.io/sbat-gcr-release/securebanking/securebanking-test-data-initializer:latest
   ```
1. ```shell
   docker push !$[+TAB]
   ```
   Or
   ```shell
   docker push eu.gcr.io/sbat-gcr-release/securebanking/securebanking-test-data-initializer:latest
   ```
1. Change kubernetes context to `sbat-master-dev`
1. Set the namespace to `ig`
1. Run `docker delete pod ig-xxxxxxxx`

>The new ig pod will run the latest image pushed in the step 3.

>Check your changes on the [platform](https://iam.dev.forgerock.financial/platform)