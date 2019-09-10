pipeline {
  agent {
    kubernetes {
      label 'jenkins-worker'
      defaultContainer 'jnlp'
      yamlFile 'JenkinsPod.yaml'
    }
  }
  stages {

    /*========================================================================*/
    stage('Deploying') {
      when {
        branch 'master'
      }
      steps {
        container('helm') {
          sh(script: "helm init -c --skip-refresh", label: "Initializing Helm Client")
          sh(script: "make deploy", label: "Atomically Deploying Helm Chart")
        }
      }
    }
    /*========================================================================*/

  }
}
