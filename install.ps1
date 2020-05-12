Write-Output "Starting installation..."
$shellexec = $PSScriptRoot+"\red.exe --path=%1"
$contextMenuName = "Redify"
Write-Output "Working directory "$shellexec
#Set-ItemProperty -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\$contextMenuName -Name command -Value $shellexec
if (Test-Path -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\shell) {
      Write-Output "Shell Registry key already exists...checking redify."
} else {
    Write-Output "Create shell key"
    New-Item -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\shell
}
if (Test-Path -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\shell\$contextMenuName) {
    Write-Output "Registry key already exists...exiting."
    exit 0
}
Write-Output "Setting up registry key..."
New-Item -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\shell\$contextMenuName
New-Item -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\shell\$contextMenuName\command
Set-ItemProperty -Path Registry::HKEY_CURRENT_USER\Software\Classes\directory\shell\$contextMenuName\command -Name '(Default)' -Value $shellexec
