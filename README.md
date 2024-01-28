# `ghcr.io/amp-buildpacks/aptos`

A Cloud Native Buildpack that provides the Aptos Tool Suite

## Configuration

| Environment Variable      | Description                                                                                                                                                                                                                                                                                       |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `$BP_APTOS_VERSION` | Configure the version of Aptos to install. It can be a specific version or a wildcard like `1.*`. It defaults to the latest `2.4.0` version. |
| `$BP_ENABLE_APTOS_DEPLOY` | Enable the Aptos deploy. It defaults to `aptos move publish --skip-fetch-latest-git-deps --assume-yes`. |
| `$BP_APTOS_DEPLOY_PRIVATE_KEY` | Configure the wallet private key for Aptos deploy. `It defaults to must be specified.` |
| `$BP_APTOS_DEPLOY_NETWORK` | Configure the network for Aptos deploy. It defaults to `devnet`. |

## Usage

### 1. To use this buildpack, simply run:

```shell
pack build <image-name> \
    --path <aptos-samples-path> \
    --buildpack ghcr.io/amp-buildpacks/aptos \
    --builder paketobuildpacks/builder-jammy-base
```

For example:

```shell
pack build aptos-sample \
    --path ./samples/aptos \
    --buildpack ghcr.io/amp-buildpacks/aptos \
    --builder paketobuildpacks/builder-jammy-base
```

### 2. To run the image, simply run:

```shell
docker run -u <uid>:<gid> -it <image-name>
```

For example:

```shell
docker run -u 1001:cnb -it aptos-sample
```

## Contributing

If anything feels off, or if you feel that some functionality is missing, please
check out the [contributing
page](https://docs.amphitheatre.app/contributing/). There you will find
instructions for sharing your feedback, building the tool locally, and
submitting pull requests to the project.

## License

Copyright (c) The Amphitheatre Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Credits

Heavily inspired by https://buildpacks.io/docs/buildpack-author-guide/create-buildpack/
