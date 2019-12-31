pipeline {
    agent {
        label "lead-toolchain-goreleaser"
    }
    stages {
        stage('Test') {
            steps {
                sh 'go test'
            }
        }
        stage('Build and Release') {
            steps {
                sh 'goreleaser release'
            }
        }
    }
}