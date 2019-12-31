pipeline {
    agent any
    stages {
        stage('Test') {
            agent {
                label "lead-toolchain-goreleaser"
            }
            steps {
                sh 'go test'
            }
        }
        stage('Build and Release') {
            agent {
                label "lead-toolchain-goreleaser"
            }
            steps {
                sh 'goreleaser release'
            }
        }
    }
}