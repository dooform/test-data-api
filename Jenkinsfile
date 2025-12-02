pipeline {
    agent any

    environment {
        GO111MODULE = 'on'
        SERVICE_NAME = 'test-data-api'
        BINARY_NAME = 'test-data-api'
        REMOTE_DIR = '/opt/test-data-api'
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
                sh "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${BINARY_NAME}"
            }
        }

        stage('Deploy') {
            steps {
                script {
                    if (!env.SERVER_IP || !env.SSH_CRED_ID || !env.SERVER_USER || !env.REMOTE_DIR) {
                        error "Missing required environment variables: SERVER_IP, SSH_CRED_ID, SERVER_USER, REMOTE_DIR"
                    }

                    sshagent([SSH_CRED_ID]) {
                        // Create remote directory
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'mkdir -p ${REMOTE_DIR}'"

                        // Stop service if running (ignore errors if not running)
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo systemctl stop ${SERVICE_NAME} || true'"

                        // Prepare service file
                        sh "sed 's/REPLACE_ME_USER/${SERVER_USER}/g' ${SERVICE_NAME}.service > ${SERVICE_NAME}.service.tmp"

                        // Upload service file
                        sh "scp -o StrictHostKeyChecking=no ${SERVICE_NAME}.service.tmp ${SERVER_USER}@${SERVER_IP}:/tmp/${SERVICE_NAME}.service"
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo mv /tmp/${SERVICE_NAME}.service /etc/systemd/system/${SERVICE_NAME}.service && sudo systemctl daemon-reload'"

                        // Upload binary
                        sh "scp -o StrictHostKeyChecking=no ${BINARY_NAME} ${SERVER_USER}@${SERVER_IP}:${REMOTE_DIR}/"

                        // Inject env vars from Jenkins credentials
                        withCredentials([
                            string(credentialsId: 'shared-db-host', variable: 'DB_HOST'),
                            string(credentialsId: 'shared-db-port', variable: 'DB_PORT'),
                            string(credentialsId: 'shared-db-user', variable: 'DB_USER'),
                            string(credentialsId: 'shared-db-password', variable: 'DB_PASSWORD'),
                            string(credentialsId: 'test-data-api-db-name', variable: 'DB_NAME'),
                            string(credentialsId: 'shared-allowed-origin-1', variable: 'ALLOWED_ORIGIN_1'),
                            string(credentialsId: 'shared-allowed-origin-2', variable: 'ALLOWED_ORIGIN_2')
                        ]) {
                            sh """
                                ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'cat > ${REMOTE_DIR}/.env << EOF
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}
CORS_ORIGINS=${ALLOWED_ORIGIN_1},${ALLOWED_ORIGIN_2}
EOF'
                            """
                        }

                        // Restart service
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP} 'sudo systemctl restart ${SERVICE_NAME}'"

                        // Cleanup
                        sh "rm ${SERVICE_NAME}.service.tmp"
                    }
                }
            }
        }
    }

    post {
        success { echo 'Pipeline completed successfully.' }
        failure { echo 'Pipeline failed.' }
    }
}
