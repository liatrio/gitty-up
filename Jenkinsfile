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
        stage('Build & Publish') {
            steps {
                container('goreleaser') {
                    withCredentials([usernamePassword(credentialsId: 'jenkins-credential-github', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
                        script {
                            sh "goreleaser release --parallelism=1 ${BRANCH_NAME != 'master' ? '--skip-publish' : ''}"
                        }
                    }
                }
            }
        }
    }
}
