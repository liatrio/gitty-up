pipeline {
    agent {
        label "lead-toolchain-goreleaser"
    }
    stages {
        state('Test') {
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