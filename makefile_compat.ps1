$Makefile = Get-Content .\Makefile -Raw
$Output = ""

foreach ($Character in $Makefile) {
    if ($Character -eq "/") {
        $Character = "\"
    }
    $Output = $Output + $Character
}

Write-Host $Output
