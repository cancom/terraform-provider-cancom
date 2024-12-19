resource "cancom_windows_os_deployment" "windows_os" {
  customer_environment_id = "PN20065D"
  role                    = "ApplicationServer"
  services                = ["Managed OS", "Managed Antivirus"]
  maintenance_window_id   = ["PN200438"]
}
