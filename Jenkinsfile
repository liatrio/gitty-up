pipeline {
    agent {
        label "lead-toolchain-goreleaser"
    }
    stages {
        stage('Test') {
            steps {
                container('goreleaser') {
                    sh 'go test'
                }
            }
        }
        stage('Build and Release') {
            steps {
                container('goreleaser') {
                    sh 'git fetch --tag'
                    sh 'goreleaser release'
                }
            }
        }
    }
}
