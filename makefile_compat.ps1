Write-Host "Grabbing Makefile contents..."

$Makefile = Get-Content .\Makefile -Raw
$Output = ""

Write-Host "Replacing path characters for compatibility..."

foreach ($Character in [char[]]$Makefile) {
    if ($Character -eq "/") {
        $Character = "\"
    }
    $Output = $Output + $Character
}

Write-Host "Writing output..."

Set-Content -Path ".\Makefile" -Value $Output

Write-Host "Done!"
pause
