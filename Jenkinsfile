pipeline {
  agent {
    label "lead-toolchain-goreleaser"
  }
  environment {
    CGO_ENABLED = 1
  }
  stages {
    stage('Test') {
      steps {
        container('goreleaser') {
          sh 'go test'
        }
      }
    }
    stage('Build & Publish') {
      steps {
        container('goreleaser') {
          withCredentials([usernamePassword(credentialsId: 'jenkins-credential-github', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
            script {
              sh "goreleaser release ${BRANCH_NAME != 'master' ? '--skip-publish' : ''}"
            }
          }
        }
      }
    }
    stage('Gitty Up') {
      agent {
        kubernetes {
          yaml """
apiVersion: v1
kind: Pod
spec:
  serviceAccountName: jenkins
  containers:
  - name: gitty-up
    image: liatrio/gitty-up:0.1.3
    command: ["/bin/sh", "-c"]
    args: ["cat"]
    tty: true
"""
        }
      }
      environment {
        GITOPS_GIT_URL = "https://github.com/liatrio/gitty-up-manifest.git"
        GITOPS_REPO_FILE = "tools.json"
        // GITOPS_VALUES = "testing.gitty-up=${sh(returnStdout: true, script: 'gitty-up --version')}"
      }
      steps {
        container('gitty-up') {
          withCredentials([usernamePassword(credentialsId: 'jenkins-credential-github', usernameVariable: 'GITOPS_GIT_USERNAME', passwordVariable: 'GITOPS_GIT_PASSWORD')]) {
            script {
              sh 'gitty-up --values=testing.gitty-up=$(gitty-up --version) --dry-run'
            }
          }
        }
      }
    }
  }
}
