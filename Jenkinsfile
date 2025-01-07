pipeline {
    agent any
    environment {
        DOCKERHUB_USERNAME = "imrezaulkrm"
        APP_NAME = "bartadhara"
        IMAGE_TAG = "backend"+"${BUILD_NUMBER}"
        IMAGE_NAME = "${DOCKERHUB_USERNAME}" + "/" + "${APP_NAME}"
        REGISTRY_CREDS = 'dockerhub'
        }
    stages {
        stage('Cleanup Workspace'){
            steps {
                script {
                    cleanWs()
                }
            }
        }
        stage('source code pull from github') {
            steps {
                git branch: 'main', url: 'https://github.com/imrezaulkrm/bartadhara-backend.git'
            }
        }
        stage('Build Docker Image'){
            steps {
                sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
                sh "docker build -t ${IMAGE_NAME}:latest ."
            }
        }
        stage('Push Docker Image'){
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', passwordVariable: 'dpass', usernameVariable: 'dockeruser')]) {
                    sh "docker login -u $dockeruser --password $dpass"
                    sh "docker push ${IMAGE_NAME}:${IMAGE_TAG}"
                    sh "docker push ${IMAGE_NAME}:latest"
                }
            }
        }

        stage('Delete Docker Images'){
            steps {
                sh "docker rmi ${IMAGE_NAME}:${IMAGE_TAG}"
                sh "docker rmi ${IMAGE_NAME}:latest"
                sh "cd .."
            }
        }
        stage('Cleanup Workspace 1'){
            steps {
                script {
                    cleanWs()
                }
            }
        }
        stage('kubernetes source code pull from github') {
            steps {
                git branch: 'main', url: 'https://github.com/imrezaulkrm/bartadhara-devops.git'
            }
        }

        stage('Updating Kubernetes deployment file') {
            steps {
                sh "ls"
                sh "cat k8s/backend.yml"
                // Construct the sed command to change only line 17
                sh """sed -i '17s#image:.*#image: ${IMAGE_NAME}:${IMAGE_TAG}#' k8s/backend.yml"""
                sh "cat k8s/backend.yml"
            }
        }

        
        stage('Push the changed deployment file to Git') {
            steps {
                script {
                    sh 'git config --global user.name "imrezaulkrm"'
                    sh 'git config --global user.email "sayem010ahmed@gmail.com"'
                    sh 'git add .'
                    sh 'git commit -m "Updated the backend.yml file"'
                    withCredentials([usernamePassword(credentialsId: 'github', passwordVariable: 'gpass', usernameVariable: 'githubuser')]) {
                        sh "git push https://$githubuser:$gpass@github.com/imrezaulkrm/bartadhara-devops.git main"
                    }
                }
            }
        }
    }
}