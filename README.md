# secret-creator-secrets-manager
The purpose of this project is to enable the creation of multiple secrets in AWS Secrets Manager. It is a AWS Serverless Application Model (SAM) based app. The secrets creation process performs rate limiting to adhere to the published [Secrets Manager Rate Quota](https://docs.aws.amazon.com/secretsmanager/latest/userguide/reference_limits.html) for the CreateRequest request type.

## Requirements
* [AWS SAM CLI](https://github.com/awslabs/aws-sam-cli) 
* [Docker installed](https://www.docker.com/community-edition) - If testing locally

## Building the application

```bash
$ sam build
```

<details><summary><b>Show sample output</b></summary>

```
Building codeuri: secrets-creator/ runtime: go1.x metadata: {} functions: ['SecretsCreatorFunction']
Running GoModulesBuilder:Build

Build Succeeded

Built Artifacts  : .aws-sam/build
Built Template   : .aws-sam/build/template.yaml

Commands you can use next
=========================
[*] Invoke Function: sam local invoke
[*] Deploy: sam deploy --guided
```

</details>


### (Optional) Testing the application locally

When using this application, you might find it useful to test locally. The AWS SAM CLI provides the `sam local` command to run your application using Docker containers that simulate the execution environment of Lambda.

<details><summary><b>Show instructions</b></summary>

1. Run application:

    ```bash
    $ sam local start-api
    ```

2. If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://127.0.0.1:3000/createsecret`

3. Create a file with sample data for the test e.g. `testsecretdata.json` :

    ```json
        [{"name":"Secret1", "username":"username1", "password" : "password1"},
        {"name":"Secret2", "username":"username2", "password" : "password2"},
        {"name":"Secret3", "username":"username3", "password" : "password3"},
        {"name":"Secret4", "username":"username4", "password" : "password4"}]
    ```

4. Now, invoke the function with the test data as follows:

    ```bash
    curl -X POST http://127.0.0.1:3000/createsecret -d @testsecretdata.json --header "Content-Type: application/json"
    ```

5. If the command ran successfully then the json output should have the arn, name and versionid fields populated, and the error field should be empty


    <details><summary><b>Show sample output</b></summary>

    ```json
        [{"arn":"arn:aws:secretsmanager:<region>:<account-id>:secret:Secret1-IVIXy3","name":"Secret1","versionid":"<uuid>","error":""},{"arn":"arn:aws:secretsmanager:<region>:<account-id>:secret:Secret2-0c2jUG","name":"Secret2","versionid":"<uuid>","error":""},{"arn":"arn:aws:secretsmanager:<region>:<account-id>:secret:Secret3-gPGgiv","name":"Secret3","versionid":"<uuid>","error":""},{"arn":"arn:aws:secretsmanager:<region>:<account-id>:secret:Secret4-LRDGhu","name":"Secret4","versionid":"<uuid>", "error":""}]
    ```

    </details>



6. The test secrets can be deleted using the [AWS CLI](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/secretsmanager/delete-secret.html) or using the AWS Management Console 

</details>


### Deploy application

To deploy your application for the first time, run the following in your shell:


```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modified IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment. Create a file with the secrets to be imported, using the JSON format shown in the local testing steps. Call the API Gateway Endpoint URL, passing the file contents for the secrets to be created.

```bash
curl -X POST <API Gateway Endpoint URL> -d @<secretsimportfile> --header "Content-Type: application/json"
```

    e.g. If the file name is `secrets2beimported.json`, then:

```bash
curl -X POST https://<somerandomstring>.execute-api.<region>.amazonaws.com/Prod/createsecret/ -d @secrets2beimported.json --header "Content-Type: application/json"
```


## API Gateway Security
Refer to the Identity and access management for API Gateway page to securely control access to the API Gateway resource created in this project.

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the [LICENSE](/LICENSE) file.

