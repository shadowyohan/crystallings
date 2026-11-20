param(
    [string]$Url = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQHalmNIOZsUBY-dMOAJQVnn1-af_4E5o9Aq62Z68ZziQgvhc8ZHtsp0MM&s=10"
)

$path = "$env:USERPROFILE\Pictures\wallpaper.jpg"
Invoke-WebRequest -Uri $Url -OutFile $path

Add-Type @"
using System.Runtime.InteropServices;
public class Wallpaper {
    [DllImport("user32.dll", CharSet = CharSet.Auto)]
    public static extern int SystemParametersInfo(int uAction, int uParam, string lpvParam, int fuWinIni);
}
"@

Set-ItemProperty "HKCU:\Control Panel\Desktop" -Name WallpaperStyle -Value 10
Set-ItemProperty "HKCU:\Control Panel\Desktop" -Name TileWallpaper -Value 0
[Wallpaper]::SystemParametersInfo(0x0014, 0, $path, 0x03) | Out-Null
Write-Host "Обои установлены: $path"
