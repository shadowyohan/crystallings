param(
    [string]$Url = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQHalmNIOZsUBY-dMOAJQVnn1-af_4E5o9Aq62Z68ZziQgvhc8ZHtsp0MM&s=10",   # сюда свою прямую ссылку на картинку
    [int]$Interval = 10                                # интервал в секундах
)

if (-not ([System.Management.Automation.PSTypeName]'Wallpaper').Type) {
    Add-Type @"
using System.Runtime.InteropServices;
public class Wallpaper {
    [DllImport("user32.dll", CharSet = CharSet.Auto)]
    public static extern int SystemParametersInfo(int uAction, int uParam, string lpvParam, int fuWinIni);
}
"@
}

$dir  = "$env:USERPROFILE\Pictures\WallpaperRotator"
$path = Join-Path $dir "wallpaper.jpg"
New-Item -ItemType Directory -Force -Path $dir | Out-Null

Set-ItemProperty "HKCU:\Control Panel\Desktop" -Name WallpaperStyle -Value 10
Set-ItemProperty "HKCU:\Control Panel\Desktop" -Name TileWallpaper -Value 0

while ($true) {
    try {
        # качаем только если файла ещё нет (или он повреждён)
        if (-not (Test-Path $path) -or (Get-Item $path).Length -eq 0) {
            Invoke-WebRequest -Uri $Url -OutFile $path -UseBasicParsing
        }
        [Wallpaper]::SystemParametersInfo(0x0014, 0, $path, 0x03) | Out-Null
    }
    catch {
        # нет интернета: попробуем в следующий раз
    }
    Start-Sleep -Seconds $Interval
}
