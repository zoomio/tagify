# Contributing

## Contributor License Agreement

Atlassian requires contributors to sign a Contributor License Agreement, known as a CLA. This serves as a record stating that the contributor is entitled to contribute the code/documentation/translation to the project and is willing to have it used in distributions and derivative works (or is willing to transfer ownership).

Prior to accepting your contributions we ask that you please follow the appropriate link below to digitally sign the CLA. The Corporate CLA is for those who are contributing as a member of an organisation and the individual CLA is for those contributing as an individual.

* [CLA for corporate contributors](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b)
* [CLA for individuals](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d)

## Guidelines for pull requests

- Write tests for any changes.
- Separate unrelated changes into multiple pull requests.
- For bigger changes, make sure you start a discussion first by creating an issue and explaining the intended change.
- Ensure the build is green before you open your PR. The Pipelines build won't run by default on a remote branch, so enable Pipelines.

## Build

* [Go](https://golang.org/dl/)
* To build binary run `./_scripts/build.sh` in shell, it will produce `tagify` binary.
* To install use `./_scripts/install.sh`, it will put `tagify` binary under `~/bin` directory assuming it is in your `PATH`.