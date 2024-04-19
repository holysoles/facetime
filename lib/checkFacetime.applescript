tell application "System Events"
    set isRunning to (exists process "FaceTime")
end tell
return isRunning