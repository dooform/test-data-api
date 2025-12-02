pipeline {
    agent any

    environment {
        GO111MODULE = 'on'
        // We are now using variables set in Jenkins "Global properties" or Job configuration
        // Required variables:
        // - SERVER_IP
        // - SERVER_USER
        // - REMOTE_DIR
        // - SSH_CRED_ID
        // - DB_HOST
        // - DB_USER
        // - DB_PASSWORD
        // - DB_NAME
        // - DB_PORT
        // - CORS_ORIGINS
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
            steps {
                script {
                    // Check if variables are set
                    if (!env.SERVER_IP || !env.SSH_CRED_ID || !env.SERVER_USER || !env.REMOTE_DIR) {
                        error "Missing required environment variables: SERVER_IP, SSH_CRED_ID, SERVER_USER, REMOTE_DIR"
                    }
                    if (!env.DB_HOST || !env.DB_USER || !env.DB_PASSWORD || !env.DB_NAME || !env.DB_PORT || !env.CORS_ORIGINS) {
                        error "Missing required environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, CORS_ORIGINS"
                    }

                    echo "Deploying to ${SERVER_IP}..."
                    
                    // Use sshagent to handle authentication securely
                    sshagent([SSH_CRED_ID]) {
                        // 1. Prepare the Service File dynamically
                        // Replace placeholders with actual Jenkins Environment Variables
                        sh "sed 's/REPLACE_ME_USER/${SERVER_USER}/g' test-data-api.service > test-data-api.service.tmp"

                        // 2. Upload the Service File (Requires sudo on server usually, so we copy to tmp first)
                        sh "scp -o StrictHostKeyChecking=no test-data-api.service.tmp ${SERVER_USER}@${SERVER_IP}:/tmp/test-data-api.service"
                        
                        // Move it to the correct place and set permissions (Runs on Server)
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo mv /tmp/test-data-api.service /etc/systemd/system/test-data-api.service && sudo systemctl daemon-reload'"

                        // 3. Upload the Binary
                        sh "scp -o StrictHostKeyChecking=no test-data-api ${SERVER_USER}@${SERVER_IP}:${REMOTE_DIR}/"

                        // 4. Create .env file on server
                        sh """ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'cat > ${REMOTE_DIR}/.env << EOF
DB_HOST=${DB_HOST}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}
DB_PORT=${DB_PORT}
CORS_ORIGINS=${CORS_ORIGINS}
EOF'"""

                        // 5. Restart the service
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo systemctl restart test-data-api'"
                        
                        // Cleanup local tmp file
                        sh "rm test-data-api.service.tmp"
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
