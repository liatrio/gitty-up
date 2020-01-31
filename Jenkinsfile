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
metadata:
  labels:
    some-label: some-label-value
spec:
  serviceAccountName: jenkins
  containers:
  - name: gitty-up
    image: liatrio/gitty-up:latest
    command: ["/bin/sh", "-c"]
    args: ["cat"]
    tty: true
    // env:
      // - name: GITOPS_GIT_USERNAME
      //   valueFrom:
      //       secretKeyRef:
      //           name: jenkins-credential-github
      //           key: username
      // - name: GITOPS_GIT_PASSWORD
      //   valueFrom:
      //       secretKeyRef:
      //           name: jenkins-credential-github
      //           key: password
"""
        }
      }
      environment {
        GITOPS_GIT_URL = "https://github.com/liatrio/gitty-up-manifest.git"
        GITOPS_REPO_FILE = "tools.json"
        GITOPS_VALUES = "testing.gitty-up=${sh(returnStdout: true, script: 'gitty-up --version')}"
      }
      steps {
        container('gitty-up') {
          withCredentials([usernamePassword(credentialsId: 'jenkins-credential-github', usernameVariable: 'GITOPS_GIT_USERNAME', passwordVariable: 'GITOPS_GIT_PASSWORD')]) {
            script {
              sh '/gitops --dry-run'
            }
          }
        }
      }
    }
  }
}
