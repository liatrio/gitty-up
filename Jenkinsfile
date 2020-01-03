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
        stage('Fetch Tags') {
            steps {
                container('goreleaser') {
                    sh 'git fetch --tag'
                }
            }
        }
        stage('Build') {
            steps {
                container('goreleaser') {
                    withCredentials([usernamePassword(credentialsId: 'jenkins-credential-github', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
                        sh 'goreleaser release --parallelism=1 --skip-publish'
                    }
                }
            }
        }
        stage('Publish') {
            steps {
                container('goreleaser') {
                    withCredentials([usernamePassword(credentialsId: 'jenkins-credential-github', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
                        sh 'goreleaser release --parallelism=1'
                    }
                }
            }
        }
    }
}
