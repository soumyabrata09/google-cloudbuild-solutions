# googleCloudBuildSolutions
Continious Integration in Google Cloud Platform
whenever a github chek in happens , this will trigger the GCP trigger set with this particular repository and eventualy changes will 
be introduced to GCP in the form of Google Cloud Registry with the Docker image itself, as we are creating a docker image
see dockerfile for image creation config process
and see cloudbuild.yaml for extended configuration

 currently this source code is only exploiting the POST hhtp methods with the underlying the postgres sql