-- there is a button with description "End" for properly ending a call, but is hidden unless the sidebar is in focus.
tell application "System Events"
	tell first window of application process "FaceTime"
		tell attribute "AXCloseButton"
			click its value
		end tell
	end tell
end tell