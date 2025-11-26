pipeline {
    agent any

    environment {
        GO111MODULE = 'on'
        // --- Configuration ---
        // Replace with your actual server details or set as Jenkins Global Env Vars
        SERVER_IP = '192.168.1.100' 
        SERVER_USER = 'ubuntu'
        REMOTE_DIR = '/opt/test-data-api'
        // The ID of the credentials stored in Jenkins (Manage Jenkins -> Credentials)
        SSH_CRED_ID = 'my-ssh-key-id' 
    }

    triggers {
        pollSCM('* * * * *')
    }

    stages {
        stage('Clean Workspace') {
            steps {
                cleanWs()
            }
        }

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install & Test') {
            steps {
                sh 'go mod download'
                sh 'go test -v ./...'
            }
        }

        stage('Build') {
            steps {
                // Build for Linux (Cross-compile if Jenkins is not on Linux)
                sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o test-data-api'
            }
        }

        stage('Deploy') {
            // Only deploy when changes are pushed to the 'main' branch
            when {
                branch 'main'
            }
            steps {
                script {
                    echo "Deploying to ${SERVER_IP}..."
                    
                    // Use sshagent to handle authentication securely
                    sshagent([SSH_CRED_ID]) {
                        // 1. Stop the service (Optional, ensures binary isn't locked)
                        // sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo systemctl stop test-data-api'"

                        // 2. Upload the new binary
                        sh "scp -o StrictHostKeyChecking=no test-data-api ${SERVER_USER}@${SERVER_IP}:${REMOTE_DIR}/"

                        // 3. Restart the service to pick up changes
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo systemctl restart test-data-api'"
                    }
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline completed successfully.'
        }
        failure {
            echo 'Pipeline failed.'
        }
    }
}