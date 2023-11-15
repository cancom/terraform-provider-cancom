data "cancom_windows_os_deployment_progress" "progress" {
    deployment_id = cancom_windows_os_windows_deployment.windows_os.id
}